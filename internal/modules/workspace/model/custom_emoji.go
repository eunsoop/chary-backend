package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CustomEmoji struct {
	bun.BaseModel `bun:"table:custom_emojis,alias:ce"`

	ID        string    `bun:"type:uuid,pk" json:"id"`
	ScopeID   string    `bun:"type:uuid,not null" json:"scope_id"`
	ScopeType ScopeType `bun:"type:varchar(20),not null" json:"scope_type"` // GLOBAL or TEAM
	Name      string    `bun:"type:varchar(255),not null" json:"name"`
	ImageURL  string    `bun:"type:text,not null" json:"image_url"`
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
}

func NewCustomEmoji(scopeID string, scopeType ScopeType, name string, imageURL string) (*CustomEmoji, error) {
	if scopeID == "" {
		return nil, errors.New("scope ID cannot be empty")
	}
	if scopeType != ScopeGlobal && scopeType != ScopeTeam {
		return nil, errors.New("invalid scope type: must be GLOBAL or TEAM")
	}
	if name == "" {
		return nil, errors.New("emoji name cannot be empty")
	}
	if imageURL == "" {
		return nil, errors.New("emoji image URL cannot be empty")
	}

	if _, err := uuid.Parse(scopeID); err != nil {
		return nil, errors.New("invalid scope ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &CustomEmoji{
		ID:        id.String(),
		ScopeID:   scopeID,
		ScopeType: scopeType,
		Name:      name,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}, nil
}
