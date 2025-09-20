package services

import (
	"fmt"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/utils"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_CreateTask(t *testing.T) {
	db := setupTestDB()
	taskService := NewTaskService()
	cacheService, _ := NewCacheService()

	userID := uuid.Must(uuid.NewV4())

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		Priority:    "high",
		UserID:      userID,
	}

	// Test task creation
	createdTask, err := taskService.CreateTask(db, task, cacheService)
	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.Equal(t, "Test Task", createdTask.Title)
	assert.Equal(t, userID, createdTask.UserID)
	assert.NotEqual(t, uuid.Nil, createdTask.ID)
}

func TestTaskService_UpdateTask(t *testing.T) {
	db := setupTestDB()
	taskService := NewTaskService()
	cacheService, _ := NewCacheService()

	userID := uuid.Must(uuid.NewV4())
	taskID := uuid.Must(uuid.NewV4())

	// Create initial task
	task := models.Task{
		ID:          taskID,
		Title:       "Original Task",
		Description: "Original Description",
		Status:      "pending",
		Priority:    "medium",
		UserID:      userID,
	}
	db.Create(&task)

	// Test task update
	newTitle := "Updated Task"
	newStatus := "in_progress"
	updateReq := models.TaskUpdateRequest{
		Title:  &newTitle,
		Status: &newStatus,
	}

	updatedTask, err := taskService.UpdateTask(db, taskID, updateReq, userID, cacheService)
	assert.NoError(t, err)
	assert.NotNil(t, updatedTask)
	assert.Equal(t, "Updated Task", updatedTask.Title)
	assert.Equal(t, "in_progress", updatedTask.Status)
	assert.Equal(t, "Original Description", updatedTask.Description) // Should remain unchanged
}

func TestTaskService_GetTaskByID(t *testing.T) {
	db := setupTestDB()
	taskService := NewTaskService()
	cacheService, _ := NewCacheService()

	userID := uuid.Must(uuid.NewV4())
	taskID := uuid.Must(uuid.NewV4())

	// Create test task
	task := models.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		Priority:    "high",
		UserID:      userID,
	}
	db.Create(&task)

	// Test getting task by ID
	retrievedTask, err := taskService.GetTaskByID(db, taskID, userID, false, cacheService)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTask)
	assert.Equal(t, taskID, retrievedTask.ID)
	assert.Equal(t, "Test Task", retrievedTask.Title)
}

func TestTaskService_DeleteTask(t *testing.T) {
	db := setupTestDB()
	taskService := NewTaskService()
	cacheService, _ := NewCacheService()

	userID := uuid.Must(uuid.NewV4())
	taskID := uuid.Must(uuid.NewV4())

	// Create test task
	task := models.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		Priority:    "high",
		UserID:      userID,
	}
	db.Create(&task)

	// Test task deletion
	err := taskService.DeleteTask(db, taskID, userID, false, cacheService)
	assert.NoError(t, err)

	// Verify task is deleted
	var deletedTask models.Task
	result := db.First(&deletedTask, taskID)
	assert.Error(t, result.Error)
}

func TestTaskService_GetTasksWithPagination(t *testing.T) {
	db := setupTestDB()
	taskService := NewTaskService()
	cacheService, _ := NewCacheService()

	userID := uuid.Must(uuid.NewV4())

	// Create multiple test tasks
	for i := 0; i < 15; i++ {
		task := models.Task{
			ID:     uuid.Must(uuid.NewV4()),
			Title:  fmt.Sprintf("Task %d", i+1),
			Status: "pending",
			UserID: userID,
		}
		db.Create(&task)
	}

	// Test pagination
	pagination := utils.PaginationParams{
		Page:     1,
		PageSize: 10,
		Offset:   0,
		Limit:    10,
	}

	filters := utils.FilterParams{
		Search:    "",
		SortBy:    "created_at",
		SortOrder: "desc",
		Filters:   make(map[string]string),
	}

	response, err := taskService.GetTasks(db, userID, false, pagination, filters, cacheService)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 10)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.PageSize)
	assert.Equal(t, int64(15), response.Pagination.Total)
	assert.True(t, response.Pagination.HasNext)
	assert.False(t, response.Pagination.HasPrev)
}
