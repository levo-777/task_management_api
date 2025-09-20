package services

import (
	"testing"
	"task-manager/backend/internal/models"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the schema
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}, &models.Permission{}, &models.RolePermission{})

	return db
}

func TestAuthService_LoginUser(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService()

	// Create test user with hashed password
	hashedPassword, _ := HashPassword("password123")
	user := models.User{
		ID:       uuid.Must(uuid.NewV4()),
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	db.Create(&user)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid credentials",
			username: "testuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid username",
			username: "wronguser",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "invalid password",
			username: "testuser",
			password: "wrongpassword",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.LoginUser(db, tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Username)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	// Test password hashing
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Test password verification
	isValid := VerifyPassword(hashedPassword, password)
	assert.True(t, isValid)

	// Test invalid password
	isValid = VerifyPassword(hashedPassword, "wrongpassword")
	assert.False(t, isValid)
}

func TestAuthService_GenerateToken(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService()

	// Create test user
	userID := uuid.Must(uuid.NewV4())
	user := models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	db.Create(&user)

	// Create role and user role
	role := models.Role{
		ID:   uuid.Must(uuid.NewV4()),
		Name: "user",
	}
	db.Create(&role)

	userRole := models.UserRole{
		ID:     uuid.Must(uuid.NewV4()),
		UserID: userID,
		RoleID: role.ID,
	}
	db.Create(&userRole)

	// Test token generation
	accessToken, refreshToken, err := authService.GenerateToken(db, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestAuthService_ValidateRefreshToken(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService()

	// Create test user
	userID := uuid.Must(uuid.NewV4())
	refreshTokenUUID := uuid.Must(uuid.NewV4())

	// Create token
	token := models.Token{
		ID:           uuid.Must(uuid.NewV4()),
		UserID:       userID,
		RefreshToken: refreshTokenUUID,
	}
	db.Create(&token)

	// Test valid refresh token
	validToken, err := authService.ValidateRefreshToken(db, refreshTokenUUID)
	assert.NoError(t, err)
	assert.NotNil(t, validToken)
	assert.Equal(t, userID, validToken.UserID)

	// Test invalid refresh token
	invalidTokenUUID := uuid.Must(uuid.NewV4())
	_, err = authService.ValidateRefreshToken(db, invalidTokenUUID)
	assert.Error(t, err)
}
