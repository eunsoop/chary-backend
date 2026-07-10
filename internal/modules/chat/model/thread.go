package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Thread struct {
	bun.BaseModel `bun:"table:threads,alias:th"`

	ID              string    `bun:"type:uuid,pk" json:"id"`
	ParentMessageID string    `bun:"type:uuid,not null,unique" json:"parent_message_id"`
	CreatedAt       time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	ParentMessage *Message   `bun:"rel:belongs-to,join:parent_message_id=id" json:"parent_message,omitempty"`
	Messages      []*Message `bun:"rel:has-many,join:id=thread_id" json:"messages,omitempty"`
}

func NewThread(parentMessageID string) (*Thread, error) {
	if parentMessageID == "" {
		return nil, errors.New("parent message ID cannot be empty")
	}

	if _, err := uuid.Parse(parentMessageID); err != nil {
		return nil, errors.New("invalid parent message ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Thread{
		ID:              id.String(),
		ParentMessageID: parentMessageID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}
