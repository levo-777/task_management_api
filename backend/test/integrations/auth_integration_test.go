package integrations

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"task-manager/backend/internal/handlers"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate
	db.AutoMigrate(&models.User{}, &models.Token{}, &models.Role{}, &models.UserRole{}, &models.Permission{}, &models.RolePermission{})

	// Initialize services
	authService := services.NewAuthService()
	registerService := services.NewRegisterService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, authService)
	registerHandler := handlers.NewRegisterHandler(db, registerService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", registerHandler.Registration)
			authRoutes.POST("/login", authHandler.Login)
		}
	}

	return router, db
}

func TestAuthIntegration_RegisterAndLogin(t *testing.T) {
	router, db := setupTestRouter()

	// Test user registration
	registerReq := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}

	registerJSON, _ := json.Marshal(registerReq)
	registerReqBody := bytes.NewBuffer(registerJSON)

	registerResp := httptest.NewRecorder()
	registerHTTPReq, _ := http.NewRequest("POST", "/api/v1/auth/register", registerReqBody)
	registerHTTPReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(registerResp, registerHTTPReq)

	assert.Equal(t, http.StatusCreated, registerResp.Code)

	// Verify user was created
	var user models.User
	result := db.Where("username = ?", "testuser").First(&user)
	assert.NoError(t, result.Error)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)

	// Test user login
	loginReq := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	loginJSON, _ := json.Marshal(loginReq)
	loginReqBody := bytes.NewBuffer(loginJSON)

	loginResp := httptest.NewRecorder()
	loginHTTPReq, _ := http.NewRequest("POST", "/api/v1/auth/login", loginReqBody)
	loginHTTPReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(loginResp, loginHTTPReq)

	assert.Equal(t, http.StatusOK, loginResp.Code)

	// Verify response contains tokens
	var loginResponse map[string]interface{}
	err := json.Unmarshal(loginResp.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	assert.Contains(t, loginResponse, "access_token")
	assert.Contains(t, loginResponse, "refresh_token")
	assert.Contains(t, loginResponse, "expires_in")
}

func TestAuthIntegration_InvalidCredentials(t *testing.T) {
	router, _ := setupTestRouter()

	// Test login with invalid credentials
	loginReq := map[string]string{
		"username": "nonexistent",
		"password": "wrongpassword",
	}

	loginJSON, _ := json.Marshal(loginReq)
	loginReqBody := bytes.NewBuffer(loginJSON)

	loginResp := httptest.NewRecorder()
	loginHTTPReq, _ := http.NewRequest("POST", "/api/v1/auth/login", loginReqBody)
	loginHTTPReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(loginResp, loginHTTPReq)

	assert.Equal(t, http.StatusUnauthorized, loginResp.Code)
}

func TestAuthIntegration_DuplicateRegistration(t *testing.T) {
	router, db := setupTestRouter()

	// Create initial user
	user := models.User{
		ID:       uuid.Must(uuid.NewV4()),
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "hashedpassword",
	}
	db.Create(&user)

	// Try to register with same username
	registerReq := map[string]string{
		"username": "existinguser",
		"email":    "new@example.com",
		"password": "password123",
	}

	registerJSON, _ := json.Marshal(registerReq)
	registerReqBody := bytes.NewBuffer(registerJSON)

	registerResp := httptest.NewRecorder()
	registerHTTPReq, _ := http.NewRequest("POST", "/api/v1/auth/register", registerReqBody)
	registerHTTPReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(registerResp, registerHTTPReq)

	assert.Equal(t, http.StatusConflict, registerResp.Code)
}
