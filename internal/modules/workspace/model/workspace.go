package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Workspace struct {
	bun.BaseModel `bun:"table:workspaces,alias:w"`

	ID        string    `bun:"type:uuid,pk" json:"id"`
	Name      string    `bun:"type:varchar(255),not null" json:"name"`
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`
}

func NewWorkspace(name string) (*Workspace, error) {
	if name == "" {
		return nil, errors.New("workspace name cannot be empty")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Workspace{
		ID:        id.String(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
