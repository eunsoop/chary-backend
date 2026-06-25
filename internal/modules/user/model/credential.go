package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProviderType string

const (
	ProviderLocal ProviderType = "LOCAL"
	ProviderOIDC  ProviderType = "OIDC"
	ProviderSAML  ProviderType = "SAML"
)

type UserCredential struct {
	bun.BaseModel `bun:"table:user_credentials,alias:uc"`

	ID       string       `bun:"type:uuid,pk" json:"id"`
	UserID   string       `bun:"type:uuid,not null" json:"user_id"`
	Provider ProviderType `bun:"type:varchar(20),not null" json:"provider"` // LOCAL, OIDC, SAML

	ProviderKey string `bun:"type:varchar(255),not null,unique:idx_provider_key" json:"provider_key"`

	PasswordHash string `bun:"type:text,nullzero" json:"-"`

	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`

	_ struct{} `bun:"unique:idx_provider_key,cb:provider,provider_key"`
}

func NewLocalCredential(userID, email, passwordHash string) (*UserCredential, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &UserCredential{
		ID:           id.String(),
		UserID:       userID,
		Provider:     ProviderLocal,
		ProviderKey:  email,
		PasswordHash: passwordHash,
	}, nil
}
