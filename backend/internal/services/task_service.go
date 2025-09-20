package services

import (
	"errors"
	"fmt"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/utils"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(db *gorm.DB, task models.Task, cacheService CacheService) (*models.Task, error)
	UpdateTask(db *gorm.DB, taskID uuid.UUID, updateReq models.TaskUpdateRequest, userID uuid.UUID, cacheService CacheService) (*models.Task, error)
	DeleteTask(db *gorm.DB, taskID uuid.UUID, userID uuid.UUID, isAdmin bool, cacheService CacheService) error
	GetTaskByID(db *gorm.DB, taskID uuid.UUID, userID uuid.UUID, isAdmin bool, cacheService CacheService) (*models.Task, error)
	GetTasksByUser(db *gorm.DB, userID uuid.UUID, pagination utils.PaginationParams, filters utils.FilterParams, cacheService CacheService) (utils.PaginationResponse, error)
	GetTasks(db *gorm.DB, userID uuid.UUID, isAdmin bool, pagination utils.PaginationParams, filters utils.FilterParams, cacheService CacheService) (utils.PaginationResponse, error)
}

type TaskServiceImpl struct{}

func NewTaskService() *TaskServiceImpl {
	return &TaskServiceImpl{}
}

func (s *TaskServiceImpl) CreateTask(db *gorm.DB, task models.Task, cacheService CacheService) (*models.Task, error) {
	task.ID = uuid.Must(uuid.NewV4())
	
	result := db.Create(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	// Cache the new task
	cacheService.SetTask(task.ID, task)

	// Invalidate user tasks cache
	cacheService.InvalidateUserCache(task.UserID)

	return &task, nil
}

func (s *TaskServiceImpl) UpdateTask(db *gorm.DB, taskID uuid.UUID, updateReq models.TaskUpdateRequest, userID uuid.UUID, cacheService CacheService) (*models.Task, error) {
	// Try to get from cache first
	if cachedTask, found := cacheService.GetTask(taskID); found {
		if task, ok := cachedTask.(*models.Task); ok {
			// Check if user owns the task (unless admin)
			if task.UserID != userID {
				return nil, errors.New("unauthorized: cannot update task owned by another user")
			}
		}
	}

	var task models.Task
	
	// Find the task
	result := db.Where("id = ?", taskID).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, result.Error
	}

	// Check if user owns the task (unless admin)
	if task.UserID != userID {
		return nil, errors.New("unauthorized: cannot update task owned by another user")
	}

	// Update fields if provided
	if updateReq.Title != nil {
		task.Title = *updateReq.Title
	}
	if updateReq.Description != nil {
		task.Description = *updateReq.Description
	}
	if updateReq.Status != nil {
		task.Status = *updateReq.Status
	}
	if updateReq.Priority != nil {
		task.Priority = *updateReq.Priority
	}

	result = db.Save(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update cache
	cacheService.SetTask(task.ID, task)
	
	// Invalidate user tasks cache
	cacheService.InvalidateUserCache(task.UserID)

	return &task, nil
}

func (s *TaskServiceImpl) DeleteTask(db *gorm.DB, taskID uuid.UUID, userID uuid.UUID, isAdmin bool, cacheService CacheService) error {
	var task models.Task
	
	// Find the task
	result := db.Where("id = ?", taskID).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("task not found")
		}
		return result.Error
	}

	// Check if user owns the task (unless admin)
	if !isAdmin && task.UserID != userID {
		return errors.New("unauthorized: cannot delete task owned by another user")
	}

	result = db.Delete(&task)
	if result.Error != nil {
		return result.Error
	}

	// Invalidate caches
	cacheService.InvalidateTaskCache(taskID)
	cacheService.InvalidateUserCache(task.UserID)

	return nil
}

func (s *TaskServiceImpl) GetTaskByID(db *gorm.DB, taskID uuid.UUID, userID uuid.UUID, isAdmin bool, cacheService CacheService) (*models.Task, error) {
	// Try to get from cache first
	if cachedTask, found := cacheService.GetTask(taskID); found {
		if task, ok := cachedTask.(models.Task); ok {
			// Check if user owns the task (unless admin)
			if !isAdmin && task.UserID != userID {
				return nil, errors.New("unauthorized: cannot view task owned by another user")
			}
			return &task, nil
		}
	}

	var task models.Task
	
	result := db.Where("id = ?", taskID).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, result.Error
	}

	// Check if user owns the task (unless admin)
	if !isAdmin && task.UserID != userID {
		return nil, errors.New("unauthorized: cannot view task owned by another user")
	}

	// Cache the task
	cacheService.SetTask(taskID, task)

	return &task, nil
}

func (s *TaskServiceImpl) GetTasksByUser(db *gorm.DB, userID uuid.UUID, pagination utils.PaginationParams, filters utils.FilterParams, cacheService CacheService) (utils.PaginationResponse, error) {
	// Create cache key based on parameters
	cacheKey := fmt.Sprintf("user_tasks:%s:page:%d:size:%d:search:%s:sort:%s:%s", 
		userID.String(), pagination.Page, pagination.PageSize, filters.Search, filters.SortBy, filters.SortOrder)
	
	// Try to get from cache first
	if cachedTasks, found := cacheService.Get(cacheKey); found {
		if response, ok := cachedTasks.(utils.PaginationResponse); ok {
			return response, nil
		}
	}

	var tasks []models.Task
	var total int64

	// Build query
	query := db.Model(&models.Task{}).Where("user_id = ?", userID)
	
	// Apply search
	query = utils.ApplySearch(query, filters.Search, []string{"title", "description"})
	
	// Apply filters
	allowedFilters := []string{"status", "priority"}
	query = utils.ApplyFilters(query, filters.Filters, allowedFilters)
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return utils.PaginationResponse{}, err
	}
	
	// Apply sorting and pagination
	allowedSortFields := []string{"title", "status", "priority", "created_at", "updated_at"}
	query = utils.ApplySorting(query, filters.SortBy, filters.SortOrder, allowedSortFields)
	query = query.Offset(pagination.Offset).Limit(pagination.Limit)
	
	result := query.Find(&tasks)
	if result.Error != nil {
		return utils.PaginationResponse{}, result.Error
	}

	// Create pagination response
	response := utils.CreatePaginationResponse(tasks, total, pagination)
	
	// Cache the response
	cacheService.Set(cacheKey, response, 1024)

	return response, nil
}

func (s *TaskServiceImpl) GetTasks(db *gorm.DB, userID uuid.UUID, isAdmin bool, pagination utils.PaginationParams, filters utils.FilterParams, cacheService CacheService) (utils.PaginationResponse, error) {
	// Create cache key based on parameters
	cacheKey := fmt.Sprintf("tasks:user:%s:admin:%t:page:%d:size:%d:search:%s:sort:%s:%s", 
		userID.String(), isAdmin, pagination.Page, pagination.PageSize, filters.Search, filters.SortBy, filters.SortOrder)
	
	// Try to get from cache first
	if cachedTasks, found := cacheService.Get(cacheKey); found {
		if response, ok := cachedTasks.(utils.PaginationResponse); ok {
			return response, nil
		}
	}

	var tasks []models.Task
	var total int64

	// Build query
	var query *gorm.DB
	if isAdmin {
		// Admin can see all tasks
		query = db.Model(&models.Task{})
	} else {
		// Regular user can only see their own tasks
		query = db.Model(&models.Task{}).Where("user_id = ?", userID)
	}
	
	// Apply search
	query = utils.ApplySearch(query, filters.Search, []string{"title", "description"})
	
	// Apply filters
	allowedFilters := []string{"status", "priority", "user_id"}
	query = utils.ApplyFilters(query, filters.Filters, allowedFilters)
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return utils.PaginationResponse{}, err
	}
	
	// Apply sorting and pagination
	allowedSortFields := []string{"title", "status", "priority", "created_at", "updated_at", "user_id"}
	query = utils.ApplySorting(query, filters.SortBy, filters.SortOrder, allowedSortFields)
	query = query.Offset(pagination.Offset).Limit(pagination.Limit)
	
	result := query.Find(&tasks)
	if result.Error != nil {
		return utils.PaginationResponse{}, result.Error
	}

	// Create pagination response
	response := utils.CreatePaginationResponse(tasks, total, pagination)
	
	// Cache the response
	cacheService.Set(cacheKey, response, 2048)

	return response, nil
}
