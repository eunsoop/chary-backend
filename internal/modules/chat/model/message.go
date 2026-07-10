package model

import (
	usermodel "chary/internal/modules/user/model"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`

	ID               string    `bun:"type:uuid,pk" json:"id"`
	ChannelID        string    `bun:"type:uuid,not null" json:"channel_id"`
	SenderID         string    `bun:"type:uuid,not null" json:"sender_id"`
	EncryptedContent string    `bun:"type:text,not null" json:"encrypted_content"`
	EncryptedDEK     string    `bun:"type:text,not null" json:"encrypted_dek"`
	ThreadID         *string   `bun:"type:uuid,nullzero" json:"thread_id,omitempty"`
	CreatedAt        time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt        time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Sender *usermodel.User `bun:"rel:belongs-to,join:sender_id=id" json:"sender,omitempty"`
	Thread *Thread         `bun:"rel:belongs-to,join:thread_id=id" json:"thread,omitempty"`
}

func NewMessage(channelID, senderID string, threadID *string) (*Message, error) {
	if channelID == "" {
		return nil, errors.New("channel ID cannot be empty")
	}
	if senderID == "" {
		return nil, errors.New("sender ID cannot be empty")
	}

	if _, err := uuid.Parse(channelID); err != nil {
		return nil, errors.New("invalid channel ID format: " + err.Error())
	}
	if _, err := uuid.Parse(senderID); err != nil {
		return nil, errors.New("invalid sender ID format: " + err.Error())
	}
	if threadID != nil {
		if _, err := uuid.Parse(*threadID); err != nil {
			return nil, errors.New("invalid thread ID format: " + err.Error())
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Message{
		ID:        id.String(),
		ChannelID: channelID,
		SenderID:  senderID,
		ThreadID:  threadID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GenerateRandomKey generates a random cryptographically secure key of the specified size (in bytes).
func GenerateRandomKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptAES_GCM encrypts plaintext using the given key and AES-256-GCM.
// The nonce is prepended to the ciphertext, and the final result is base64 encoded.
func EncryptAES_GCM(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES_GCM decrypts a base64 encoded ciphertext (which has the nonce prepended) using the given key.
func DecryptAES_GCM(base64Ciphertext string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptContent encrypts the raw message content using a newly generated DEK,
// and encrypts that DEK using the provided Key Encryption Key (KEK).
func (m *Message) EncryptContent(rawContent string, kek []byte) error {
	if len(kek) != 32 {
		return errors.New("key encryption key (KEK) must be exactly 32 bytes for AES-256")
	}

	// 1. Generate DEK (32 bytes for AES-256)
	dek, err := GenerateRandomKey(32)
	if err != nil {
		return err
	}

	// 2. Encrypt raw content with DEK
	encContent, err := EncryptAES_GCM([]byte(rawContent), dek)
	if err != nil {
		return err
	}

	// 3. Encrypt DEK with KEK
	encDEK, err := EncryptAES_GCM(dek, kek)
	if err != nil {
		return err
	}

	m.EncryptedContent = encContent
	m.EncryptedDEK = encDEK
	return nil
}

// DecryptContent decrypts the message content by first decrypting the DEK using the KEK,
// and then decrypting the content using the decrypted DEK.
func (m *Message) DecryptContent(kek []byte) (string, error) {
	if len(kek) != 32 {
		return "", errors.New("key encryption key (KEK) must be exactly 32 bytes for AES-256")
	}

	// 1. Decrypt DEK using KEK
	dek, err := DecryptAES_GCM(m.EncryptedDEK, kek)
	if err != nil {
		return "", errors.New("failed to decrypt DEK: " + err.Error())
	}

	// 2. Decrypt content using decrypted DEK
	rawContent, err := DecryptAES_GCM(m.EncryptedContent, dek)
	if err != nil {
		return "", errors.New("failed to decrypt content: " + err.Error())
	}

	return string(rawContent), nil
}
