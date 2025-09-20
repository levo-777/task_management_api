package handlers

import (
	"errors"
	"net/http"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/services"
	"task-manager/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db          *gorm.DB
	taskService services.TaskService
	cacheService services.CacheService
}

func NewTaskHandler(db *gorm.DB, taskService services.TaskService, cacheService services.CacheService) *TaskHandler {
	return &TaskHandler{db: db, taskService: taskService, cacheService: cacheService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID.(uuid.UUID),
	}

	if req.Status != "" {
		task.Status = req.Status
	} else {
		task.Status = "pending"
	}

	if req.Priority != "" {
		task.Priority = req.Priority
	} else {
		task.Priority = "medium"
	}

	createdTask, err := h.taskService.CreateTask(h.db, task, h.cacheService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "task created successfully", "task": createdTask})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.FromString(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	updatedTask, err := h.taskService.UpdateTask(h.db, taskID, req, userID.(uuid.UUID), h.cacheService)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: cannot update task owned by another user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully", "task": updatedTask})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.FromString(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isAdmin, _ := c.Get("is_admin")

	err = h.taskService.DeleteTask(h.db, taskID, userID.(uuid.UUID), isAdmin.(bool), h.cacheService)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: cannot delete task owned by another user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.FromString(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isAdmin, _ := c.Get("is_admin")

	task, err := h.taskService.GetTaskByID(h.db, taskID, userID.(uuid.UUID), isAdmin.(bool), h.cacheService)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: cannot view task owned by another user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *TaskHandler) GetTasksByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is accessing their own tasks or if they're admin
	authenticatedUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isAdmin, _ := c.Get("is_admin")
	if !isAdmin.(bool) && userID != authenticatedUserID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to view other user's tasks"})
		return
	}

	// Get pagination and filter parameters
	pagination := utils.GetPaginationParams(c)
	filters := utils.GetFilterParams(c)

	response, err := h.taskService.GetTasksByUser(h.db, userID, pagination, filters, h.cacheService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isAdmin, _ := c.Get("is_admin")

	// Get pagination and filter parameters
	pagination := utils.GetPaginationParams(c)
	filters := utils.GetFilterParams(c)

	response, err := h.taskService.GetTasks(h.db, userID.(uuid.UUID), isAdmin.(bool), pagination, filters, h.cacheService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func handleTaskError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process task request",
		})
	}
}
