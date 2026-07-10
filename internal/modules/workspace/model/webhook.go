package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type WebhookType string

const (
	WebhookIncoming WebhookType = "INCOMING"
	WebhookOutgoing WebhookType = "OUTGOING"
)

type Webhook struct {
	bun.BaseModel `bun:"table:webhooks,alias:wh"`

	ID            string      `bun:"type:uuid,pk" json:"id"`
	IntegrationID string      `bun:"type:uuid,not null,unique" json:"integration_id"`
	Type          WebhookType `bun:"type:varchar(20),not null" json:"type"` // INCOMING or OUTGOING

	// SecretToken is used for Incoming Webhooks to validate the origin path / token authentication
	SecretToken string `bun:"type:varchar(255),nullzero,unique" json:"secret_token,omitempty"`

	// TargetURL is used for Outgoing Webhooks to specify the external callback endpoint
	TargetURL string `bun:"type:text,nullzero" json:"target_url,omitempty"`

	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Integration *Integration `bun:"rel:belongs-to,join:integration_id=id" json:"integration,omitempty"`
}

func GenerateSecretToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewIncomingWebhook(integrationID string) (*Webhook, error) {
	if integrationID == "" {
		return nil, errors.New("integration ID cannot be empty")
	}
	if _, err := uuid.Parse(integrationID); err != nil {
		return nil, errors.New("invalid integration ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	token, err := GenerateSecretToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Webhook{
		ID:            id.String(),
		IntegrationID: integrationID,
		Type:          WebhookIncoming,
		SecretToken:   token,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func NewOutgoingWebhook(integrationID string, targetURL string) (*Webhook, error) {
	if integrationID == "" {
		return nil, errors.New("integration ID cannot be empty")
	}
	if targetURL == "" {
		return nil, errors.New("target URL cannot be empty")
	}
	if _, err := uuid.Parse(integrationID); err != nil {
		return nil, errors.New("invalid integration ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Webhook{
		ID:            id.String(),
		IntegrationID: integrationID,
		Type:          WebhookOutgoing,
		TargetURL:     targetURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
