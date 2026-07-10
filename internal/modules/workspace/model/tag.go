package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`

	ID          string    `bun:"type:uuid,pk" json:"id"`
	TeamScopeID string    `bun:"type:uuid,not null" json:"team_scope_id"`
	Name        string    `bun:"type:varchar(255),not null" json:"name"`
	Color       string    `bun:"type:varchar(10),nullzero" json:"color,omitempty"`
	CreatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`

	TeamScope *TeamScope `bun:"rel:belongs-to,join:team_scope_id=id" json:"team_scope,omitempty"`
}

func NewTag(teamScopeID string, name string, color string) (*Tag, error) {
	if teamScopeID == "" {
		return nil, errors.New("team scope ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("tag name cannot be empty")
	}

	if _, err := uuid.Parse(teamScopeID); err != nil {
		return nil, errors.New("invalid team scope ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Tag{
		ID:          id.String(),
		TeamScopeID: teamScopeID,
		Name:        name,
		Color:       color,
		CreatedAt:   time.Now(),
	}, nil
}
