package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeamScope struct {
	bun.BaseModel `bun:"table:team_scopes,alias:ts"`

	ID          string    `bun:"type:uuid,pk" json:"id"`
	WorkspaceID string    `bun:"type:uuid,not null" json:"workspace_id"`
	Name        string    `bun:"type:varchar(255),not null" json:"name"`
	CreatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Workspace *Workspace `bun:"rel:belongs-to,join:workspace_id=id" json:"workspace,omitempty"`
}

func NewTeamScope(workspaceID string, name string) (*TeamScope, error) {
	if workspaceID == "" {
		return nil, errors.New("workspace ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("team name cannot be empty")
	}

	if _, err := uuid.Parse(workspaceID); err != nil {
		return nil, errors.New("invalid workspace ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &TeamScope{
		ID:          id.String(),
		WorkspaceID: workspaceID,
		Name:        name,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
