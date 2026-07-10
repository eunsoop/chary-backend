package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GlobalScope struct {
	bun.BaseModel `bun:"table:global_scopes,alias:gs"`

	ID          string    `bun:"type:uuid,pk" json:"id"`
	WorkspaceID string    `bun:"type:uuid,not null,unique" json:"workspace_id"`
	Domain      string    `bun:"type:varchar(255),not null,unique" json:"domain"`
	CreatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Workspace *Workspace `bun:"rel:belongs-to,join:workspace_id=id" json:"workspace,omitempty"`
}

func NewGlobalScope(workspaceID string, domain string) (*GlobalScope, error) {
	if workspaceID == "" {
		return nil, errors.New("workspace ID cannot be empty")
	}
	if domain == "" {
		return nil, errors.New("domain cannot be empty")
	}

	wsUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return nil, errors.New("invalid workspace ID format: " + err.Error())
	}

	// Generate deterministic UUIDv5 scoped under the workspace UUID namespace to isolate domains across workspaces.
	idV5 := uuid.NewSHA1(wsUUID, []byte(domain))
	if idV5 == uuid.Nil {
		return nil, errors.New("generated GlobalScope UUIDv5 cannot be the Nil UUID")
	}

	now := time.Now()
	return &GlobalScope{
		ID:          idV5.String(),
		WorkspaceID: workspaceID,
		Domain:      domain,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
