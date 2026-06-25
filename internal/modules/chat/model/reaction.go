package model

import (
	"chary/internal/modules/user/model"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Reaction struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`

	ID        string    `bun:"type:uuid,pk" json:"id"`
	MessageID string    `bun:"type:uuid,not null,on_delete:cascade" json:"message_id"`
	UserID    string    `bun:"type:uuid,not null" json:"user_id"`
	Emoji     string    `bun:"type:varchar(50),not null" json:"emoji"`
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`

	User *model.User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

func NewReaction(messageID, userID, emoji string) (*Reaction, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Reaction{
		ID:        id.String(),
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	}, nil
}
