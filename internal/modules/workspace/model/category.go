package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ScopeType string

const (
	ScopeGlobal ScopeType = "GLOBAL"
	ScopeTeam   ScopeType = "TEAM"
)

type Category struct {
	bun.BaseModel `bun:"table:categories,alias:c"`

	ID        string    `bun:"type:uuid,pk" json:"id"`
	ScopeID   string    `bun:"type:uuid,not null" json:"scope_id"`
	ScopeType ScopeType `bun:"type:varchar(20),not null" json:"scope_type"` // GLOBAL or TEAM
	Name      string    `bun:"type:varchar(255),not null" json:"name"`
	Order     int       `bun:"type:integer,not null,default:0" json:"order"`
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`
}

func NewCategory(scopeID string, scopeType ScopeType, name string, order int) (*Category, error) {
	if scopeID == "" {
		return nil, errors.New("scope ID cannot be empty")
	}
	if scopeType != ScopeGlobal && scopeType != ScopeTeam {
		return nil, errors.New("invalid scope type: must be GLOBAL or TEAM")
	}
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	if _, err := uuid.Parse(scopeID); err != nil {
		return nil, errors.New("invalid scope ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Category{
		ID:        id.String(),
		ScopeID:   scopeID,
		ScopeType: scopeType,
		Name:      name,
		Order:     order,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
