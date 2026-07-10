package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)
type ProviderType string

const (
	ProviderLocal ProviderType = "LOCAL"
	ProviderOIDC  ProviderType = "OIDC"
	ProviderSAML  ProviderType = "SAML"
	ProviderLDAP  ProviderType = "LDAP"
)

type UserCredential struct {
	bun.BaseModel `bun:"table:user_credentials,alias:uc"`

	ID       string       `bun:"type:uuid,pk" json:"id"`
	UserID   string       `bun:"type:uuid,not null" json:"user_id"`
	Provider ProviderType `bun:"type:varchar(20),not null" json:"provider"` // LOCAL, OIDC, SAML, LDAP

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

func NewExternalCredential(userID string, provider ProviderType, providerKey string) (*UserCredential, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if provider != ProviderOIDC && provider != ProviderSAML && provider != ProviderLDAP {
		return nil, errors.New("invalid external provider type")
	}
	if providerKey == "" {
		return nil, errors.New("provider key cannot be empty")
	}

	if _, err := uuid.Parse(userID); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &UserCredential{
		ID:          id.String(),
		UserID:      userID,
		Provider:    provider,
		ProviderKey: providerKey,
	}, nil
}
