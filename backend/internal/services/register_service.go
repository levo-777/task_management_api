package services

import (
	"errors"
	"task-manager/backend/internal/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type RegisterService interface {
	RegisterUser(db *gorm.DB, user models.User) error
}

type RegisterServiceImpl struct{}

func NewRegisterService() *RegisterServiceImpl {
	return &RegisterServiceImpl{}
}

func (s *RegisterServiceImpl) RegisterUser(db *gorm.DB, user models.User) error {
	// Check if username or email already exists
	var existingUser models.User
	result := db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser)
	if result.Error == nil {
		return errors.New("username or email already exists")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// Hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Generate UUID for user
	user.ID = uuid.Must(uuid.NewV4())

	// Create user
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// Assign default user role
	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440001"), // user role ID
	}

	if err := db.Create(&userRole).Error; err != nil {
		return err
	}

	return nil
}
