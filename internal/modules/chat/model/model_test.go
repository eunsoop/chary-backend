package model

import (
	"crypto/rand"
	"testing"

	"github.com/google/uuid"
)

func TestEnvelopeEncryption(t *testing.T) {
	// Generate random 32-byte Key Encryption Key (KEK)
	kek := make([]byte, 32)
	if _, err := rand.Read(kek); err != nil {
		t.Fatalf("failed to generate random KEK: %v", err)
	}

	channelID := uuid.New().String()
	senderID := uuid.New().String()
	msg, err := NewMessage(channelID, senderID, nil)
	if err != nil {
		t.Fatalf("unexpected error creating message: %v", err)
	}

	originalContent := "Top-secret enterprise communication content."

	// Test successful encryption
	err = msg.EncryptContent(originalContent, kek)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if msg.EncryptedContent == "" {
		t.Error("expected EncryptedContent to be non-empty")
	}
	if msg.EncryptedDEK == "" {
		t.Error("expected EncryptedDEK to be non-empty")
	}
	if msg.EncryptedContent == originalContent {
		t.Error("expected EncryptedContent to be ciphertext, not plaintext")
	}

	// Test successful decryption
	decryptedContent, err := msg.DecryptContent(kek)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decryptedContent != originalContent {
		t.Errorf("expected decrypted content to be '%s', got '%s'", originalContent, decryptedContent)
	}

	// Test decryption failure with incorrect KEK
	badKek := make([]byte, 32)
	copy(badKek, kek)
	badKek[0] ^= 0xFF // Mutate the first byte
	_, err = msg.DecryptContent(badKek)
	if err == nil {
		t.Error("expected decryption to fail with incorrect KEK, but it succeeded")
	}

	// Test invalid KEK size validation
	shortKek := make([]byte, 16)
	err = msg.EncryptContent(originalContent, shortKek)
	if err == nil {
		t.Error("expected encryption to fail with short KEK, but it succeeded")
	}
	_, err = msg.DecryptContent(shortKek)
	if err == nil {
		t.Error("expected decryption to fail with short KEK, but it succeeded")
	}
}

func TestThreadAndExpression(t *testing.T) {
	channelID := uuid.New().String()
	senderID := uuid.New().String()
	userID := uuid.New().String()

	parentMsg, err := NewMessage(channelID, senderID, nil)
	if err != nil {
		t.Fatalf("failed to create parent message: %v", err)
	}

	// Create thread pointing to parent message
	thread, err := NewThread(parentMsg.ID)
	if err != nil {
		t.Fatalf("failed to create thread: %v", err)
	}

	if thread.ParentMessageID != parentMsg.ID {
		t.Errorf("expected parent message ID '%s', got '%s'", parentMsg.ID, thread.ParentMessageID)
	}

	// Create child message in the thread
	childMsg, err := NewMessage(channelID, senderID, &thread.ID)
	if err != nil {
		t.Fatalf("failed to create child message: %v", err)
	}

	if childMsg.ThreadID == nil || *childMsg.ThreadID != thread.ID {
		t.Errorf("expected child message to reference thread ID '%s'", thread.ID)
	}

	// Create expression (reaction) on both parent message and child thread message
	exp1, err := NewExpression(parentMsg.ID, userID, "emoji_thumbsup")
	if err != nil {
		t.Fatalf("failed to create expression on parent: %v", err)
	}
	if exp1.MessageID != parentMsg.ID || exp1.UserID != userID || exp1.EmojiID != "emoji_thumbsup" {
		t.Errorf("expression field mismatch")
	}

	exp2, err := NewExpression(childMsg.ID, userID, "emoji_heart")
	if err != nil {
		t.Fatalf("failed to create expression on thread child: %v", err)
	}
	if exp2.MessageID != childMsg.ID || exp2.UserID != userID || exp2.EmojiID != "emoji_heart" {
		t.Errorf("expression field mismatch")
	}
}
