package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Channel struct {
	bun.BaseModel `bun:"table:channels,alias:ch"`

	ID         string    `bun:"type:uuid,pk" json:"id"`
	CategoryID string    `bun:"type:uuid,not null" json:"category_id"`
	Name       string    `bun:"type:varchar(255),not null" json:"name"`
	IsPrivate  bool      `bun:"type:boolean,not null,default:false" json:"is_private"`
	CreatedAt  time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Category     *Category      `bun:"rel:belongs-to,join:category_id=id" json:"category,omitempty"`
	Integrations []*Integration `bun:"rel:has-many,join:id=channel_id" json:"integrations,omitempty"`
}

func NewChannel(categoryID string, name string, isPrivate bool) (*Channel, error) {
	if categoryID == "" {
		return nil, errors.New("category ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("channel name cannot be empty")
	}

	if _, err := uuid.Parse(categoryID); err != nil {
		return nil, errors.New("invalid category ID format: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Channel{
		ID:         id.String(),
		CategoryID: categoryID,
		Name:       name,
		IsPrivate:  isPrivate,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
