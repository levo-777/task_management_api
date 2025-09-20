package services

import (
	"errors"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/utils"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	LoginUser(db *gorm.DB, username, password string) (*models.User, error)
	GenerateToken(db *gorm.DB, userID uuid.UUID) (string, string, error)
	ValidateRefreshToken(db *gorm.DB, refreshToken uuid.UUID) (*models.Token, error)
	InvalidateRefreshToken(db *gorm.DB, refreshToken uuid.UUID) error
	GetUserRolesAndPermissions(db *gorm.DB, userID uuid.UUID) ([]string, bool, []utils.Permission, error)
}

type AuthServiceImpl struct{}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{}
}

func (s *AuthServiceImpl) LoginUser(db *gorm.DB, username, password string) (*models.User, error) {
	var user models.User
	
	result := db.Where("username = ? OR email = ?", username, username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, result.Error
	}

	if !VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *AuthServiceImpl) GenerateToken(db *gorm.DB, userID uuid.UUID) (string, string, error) {
	// Get user roles and permissions
	roles, isAdmin, permissions, err := s.GetUserRolesAndPermissions(db, userID)
	if err != nil {
		return "", "", err
	}

	// Get user details for JWT
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return "", "", err
	}

	// Generate access token
	accessToken, err := utils.GenerateJWT(userID, user.Username, roles, isAdmin, permissions)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := uuid.Must(uuid.NewV4())
	expiresAt := time.Now().Add(time.Hour)

	// Store refresh token in database
	token := models.Token{
		ID:           uuid.Must(uuid.NewV4()),
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := db.Create(&token).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.String(), nil
}

func (s *AuthServiceImpl) ValidateRefreshToken(db *gorm.DB, refreshToken uuid.UUID) (*models.Token, error) {
	var token models.Token
	
	result := db.Where("refresh_token = ? AND expires_at > ?", refreshToken, time.Now()).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired refresh token")
		}
		return nil, result.Error
	}

	return &token, nil
}

func (s *AuthServiceImpl) InvalidateRefreshToken(db *gorm.DB, refreshToken uuid.UUID) error {
	return db.Where("refresh_token = ?", refreshToken).Delete(&models.Token{}).Error
}

func (s *AuthServiceImpl) GetUserRolesAndPermissions(db *gorm.DB, userID uuid.UUID) ([]string, bool, []utils.Permission, error) {
	var userRoles []models.UserRole
	
	// Get user roles
	result := db.Preload("Role").Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, false, nil, result.Error
	}

	var roles []string
	isAdmin := false
	
	for _, ur := range userRoles {
		roles = append(roles, ur.Role.Name)
		if ur.Role.Name == "admin" {
			isAdmin = true
		}
	}

	// Get permissions for all roles
	var rolePermissions []models.RolePermission
	result = db.Preload("Permission").Where("role_id IN (SELECT role_id FROM user_roles WHERE user_id = ?)", userID).Find(&rolePermissions)
	if result.Error != nil {
		return nil, false, nil, result.Error
	}

	// Group permissions by resource
	permissionMap := make(map[string][]string)
	for _, rp := range rolePermissions {
		resource := rp.Permission.Resource
		action := rp.Permission.Action
		
		if actions, exists := permissionMap[resource]; exists {
			permissionMap[resource] = append(actions, action)
		} else {
			permissionMap[resource] = []string{action}
		}
	}

	// Convert to utils.Permission format
	var permissions []utils.Permission
	for resource, actions := range permissionMap {
		permissions = append(permissions, utils.Permission{
			Resource: resource,
			Actions:  actions,
		})
	}

	return roles, isAdmin, permissions, nil
}

func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
