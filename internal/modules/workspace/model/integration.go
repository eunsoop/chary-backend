package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type IntegrationType string

const (
	IntegrationWebhook IntegrationType = "WEBHOOK"
)

type Integration struct {
	bun.BaseModel `bun:"table:integrations,alias:intg"`

	ID        string          `bun:"type:uuid,pk" json:"id"`
	ChannelID string          `bun:"type:uuid,not null" json:"channel_id"`
	Name      string          `bun:"type:varchar(255),not null" json:"name"`
	Type      IntegrationType `bun:"type:varchar(50),not null" json:"type"` // e.g., WEBHOOK
	IsEnabled bool            `bun:"type:boolean,not null,default:true" json:"is_enabled"`
	CreatedAt time.Time       `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time       `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Channel *Channel `bun:"rel:belongs-to,join:channel_id=id" json:"channel,omitempty"`
}

func NewIntegration(channelID string, name string, itype IntegrationType) (*Integration, error) {
	if channelID == "" {
		return nil, errors.New("channel ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("integration name cannot be empty")
	}
	if itype != IntegrationWebhook {
		return nil, errors.New("unsupported integration type: only WEBHOOK is supported currently")
	}

	if _, err := uuid.Parse(channelID); err != nil {
		return nil, errors.New("invalid channel ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Integration{
		ID:        id.String(),
		ChannelID: channelID,
		Name:      name,
		Type:      itype,
		IsEnabled: true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
