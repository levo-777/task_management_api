package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	gorm.Model
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RoleID       uuid.UUID `json:"role_id" gorm:"not null"`
	PermissionID uuid.UUID `json:"permission_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	Role       Role       `json:"role" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID"`
}
