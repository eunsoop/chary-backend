package model

import (
	"chary/internal/modules/user/model"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`

	ID        string `bun:"type:uuid,pk" json:"id"`
	ChannelID string `bun:"type:uuid,not null,index:idx_channel_created" json:"channel_id"`
	SenderID  string `bun:"type:uuid,not null" json:"sender_id"`
	Content   string `bun:"type:text,not null" json:"content"`

	ParentID  *string `bun:"type:uuid,nullzero" json:"parent_id,omitempty"`   // Thread
	ReplyToID *string `bun:"type:uuid,nullzero" json:"reply_to_id,omitempty"` // Reply

	Attachments map[string]interface{} `bun:"type:jsonb,nullzero" json:"attachments,omitempty"`

	CreatedAt time.Time `bun:"type:timestamptz,nullzero,index:idx_channel_created,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Sender *model.User `bun:"rel:belongs-to,join:sender_id=id" json:"sender,omitempty"`
}

func NewMessage(channelID, senderID, content string, parentID, replyToID *string, attachments map[string]interface{}) (*Message, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:          id.String(),
		ChannelID:   channelID,
		SenderID:    senderID,
		Content:     content,
		ParentID:    parentID,
		ReplyToID:   replyToID,
		Attachments: attachments,
	}, nil
}
