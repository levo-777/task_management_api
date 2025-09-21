package main

import (
	"log"
	"task-manager/backend/internal/handlers"
	"task-manager/backend/internal/middleware"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/repositories"
	"task-manager/backend/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @title Taskify API
// @version 1.0
// @description A comprehensive task management system with role-based access control and JWT authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer <token>

func main() {
	// Initialize database connection
	dbCfg := repositories.NewDatabaseConfig()
	db, err := dbCfg.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}
	defer sqlDB.Close()

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Role{},
		&models.UserRole{},
		&models.Permission{},
		&models.RolePermission{},
		&models.Task{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Initialize database with default data
	err = initializeDatabase(db)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	// Initialize cache service
	cacheService, err := services.NewCacheService()
	if err != nil {
		log.Fatal("Failed to initialize cache service: ", err)
	}

	// Initialize services
	authService := services.NewAuthService()
	registerService := services.NewRegisterService()
	userService := services.NewUserService()
	taskService := services.NewTaskService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, authService)
	registerHandler := handlers.NewRegisterHandler(db, registerService)
	userHandler := handlers.NewUserHandler(db, userService)
	taskHandler := handlers.NewTaskHandler(db, taskService, cacheService)
	refreshHandler := handlers.NewRefreshHandler(db, authService)

	// Initialize Gin router
	r := gin.Default()

	// Add custom logging middleware
	r.Use(middleware.CustomLogger(middleware.CustomLoggerConfig{
		SkipPaths: []string{"/health", "/metrics"},
	}))
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ErrorLogger())

	// Rate limiting middleware
	r.Use(middleware.RateLimitMiddleware())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://host.docker.internal"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "taskify-api",
		})
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes with stricter rate limiting
		authRoutes := v1.Group("/auth")
		authRoutes.Use(middleware.AuthRateLimitMiddleware())
		{
			authRoutes.POST("/register", registerHandler.Registration)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", refreshHandler.Refresh)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Task routes
			taskRoutes := protected.Group("/tasks")
			{
				taskRoutes.POST("", middleware.RequirePermission("task", "create"), taskHandler.CreateTask)
				taskRoutes.PUT("/:id", middleware.RequirePermission("task", "write"), taskHandler.UpdateTask)
				taskRoutes.DELETE("/:id", middleware.RequirePermission("task", "delete"), taskHandler.DeleteTask)
				taskRoutes.GET("/:id", middleware.RequirePermission("task", "read"), taskHandler.GetTaskByID)
				taskRoutes.GET("", middleware.RequirePermission("task", "read"), taskHandler.GetTasks)
			}

			// User routes
			userRoutes := protected.Group("/users")
			{
				userRoutes.GET("/profile", middleware.RequirePermission("profile", "read"), userHandler.GetUserProfile)
				userRoutes.GET("/profile/:user_id", middleware.RequirePermission("profile", "read"), userHandler.GetUserProfileByUserId)
				userRoutes.GET("/:user_id/tasks", middleware.RequirePermission("task", "read"), taskHandler.GetTasksByUser)

				// Admin only routes
				userRoutes.GET("", middleware.RequireAdmin(), userHandler.GetUsers)
				userRoutes.DELETE("/:user_id", middleware.RequireAdmin(), userHandler.DeleteUser)
			}
		}
	}

	// Start server
	log.Println("Starting server on :8080")
	r.Run(":8080")
}

func initializeDatabase(db *gorm.DB) error {
	// Check if roles already exist
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)

	if roleCount == 0 {
		// Create roles
		roles := []models.Role{
			{ID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440001"), Name: "user"},
			{ID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), Name: "admin"},
		}

		for _, role := range roles {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}

		// Create permissions
		permissions := []models.Permission{
			{ID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440001"), Resource: "profile", Action: "read"},
			{ID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440002"), Resource: "profile", Action: "write"},
			{ID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440003"), Resource: "task", Action: "create"},
			{ID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440004"), Resource: "task", Action: "read"},
			{ID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440005"), Resource: "task", Action: "write"},
		}

		for _, permission := range permissions {
			if err := db.Create(&permission).Error; err != nil {
				return err
			}
		}

		// Create role-permission mappings
		rolePermissions := []models.RolePermission{
			// Admin permissions
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440001")}, // profile:read
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440002")}, // profile:write
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440003")}, // task:create
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440004")}, // task:read
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440005")}, // task:write
			// User permissions
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440001"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440001")}, // profile:read
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440001"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440003")}, // task:create
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440001"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440004")}, // task:read
			{RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440001"), PermissionID: uuid.FromStringOrNil("750e8400-e29b-41d4-a716-446655440005")}, // task:write
		}

		for _, rp := range rolePermissions {
			if err := db.Create(&rp).Error; err != nil {
				return err
			}
		}

		// Create admin user
		adminUser := models.User{
			ID:       uuid.FromStringOrNil("bd006d41-aded-4040-9934-2ba4e909ef9a"),
			Username: "admin",
			Email:    "admin@gmail.com",
			Password: "$2a$10$ZsE5IA/WUP2CKVd0rND/rum3DKp6lLfjTGSHAmqONb9eGCLY4GYD6", // admin123
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return err
		}

		// Assign admin role to admin user
		adminUserRole := models.UserRole{
			UserID: adminUser.ID,
			RoleID: uuid.FromStringOrNil("550e8400-e29b-41d4-a716-446655440002"), // admin role
		}

		if err := db.Create(&adminUserRole).Error; err != nil {
			return err
		}

		log.Println("Database initialized with default data")
	}

	return nil
}
