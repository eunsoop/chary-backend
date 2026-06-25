package controller

import (
	"chary/internal/modules/user/model"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          string    `bun:"type:uuid,pk" json:"id"`
	Username    string    `bun:"type:varchar(50),unique,not null" json:"username"`
	Email       string    `bun:"type:varchar(255),unique,not null" json:"email"`
	DisplayName string    `bun:"type:varchar(100)" json:"display_name"`
	AvatarURL   string    `bun:"type:text" json:"avatar_url"`
	IsBot       bool      `bun:"type:boolean,not null,default:false" json:"is_bot"`
	CreatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"type:timestamptz,nullzero,default:current_timestamp" json:"updated_at"`

	Credentials []*model.UserCredential `bun:"rel:has-many,join:id=user_id" json:"-"`
}

func NewUser(username, email, displayName, avatarURL string, isBot bool) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:          id.String(),
		Username:    username,
		Email:       email,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
		IsBot:       isBot,
	}, nil
}
