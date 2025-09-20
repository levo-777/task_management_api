package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Token struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"not null"`
	RefreshToken uuid.UUID `json:"refresh_token" gorm:"not null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt    *time.Time `json:"-" gorm:"index"`
	
	User User `json:"user" gorm:"foreignKey:UserID"`
}
