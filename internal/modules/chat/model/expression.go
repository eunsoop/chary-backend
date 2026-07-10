package model

import (
	usermodel "chary/internal/modules/user/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Expression struct {
	bun.BaseModel `bun:"table:expressions,alias:exp"`

	ID        string    `bun:"type:uuid,pk" json:"id"`
	MessageID string    `bun:"type:uuid,not null,on_delete:cascade" json:"message_id"`
	UserID    string    `bun:"type:uuid,not null" json:"user_id"`
	EmojiID   string    `bun:"type:varchar(255),not null" json:"emoji_id"`
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`

	Message *Message        `bun:"rel:belongs-to,join:message_id=id" json:"message,omitempty"`
	User    *usermodel.User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

func NewExpression(messageID, userID, emojiID string) (*Expression, error) {
	if messageID == "" {
		return nil, errors.New("message ID cannot be empty")
	}
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if emojiID == "" {
		return nil, errors.New("emoji ID cannot be empty")
	}

	if _, err := uuid.Parse(messageID); err != nil {
		return nil, errors.New("invalid message ID format: " + err.Error())
	}
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Expression{
		ID:        id.String(),
		MessageID: messageID,
		UserID:    userID,
		EmojiID:   emojiID,
		CreatedAt: time.Now(),
	}, nil
}
