package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"default:pending"`
	Priority    string     `json:"priority" gorm:"default:medium"`
	UserID      uuid.UUID  `json:"user_id" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`
	
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type TaskCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

type TaskUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
}
