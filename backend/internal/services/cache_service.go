package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/gofrs/uuid"
)

type CacheService interface {
	Set(key string, value interface{}, cost int64) bool
	Get(key string) (interface{}, bool)
	Del(key string)
	GetUserProfile(userID uuid.UUID) (interface{}, bool)
	SetUserProfile(userID uuid.UUID, profile interface{}) bool
	GetTask(taskID uuid.UUID) (interface{}, bool)
	SetTask(taskID uuid.UUID, task interface{}) bool
	GetUserTasks(userID uuid.UUID) (interface{}, bool)
	SetUserTasks(userID uuid.UUID, tasks interface{}) bool
	InvalidateUserCache(userID uuid.UUID)
	InvalidateTaskCache(taskID uuid.UUID)
}

type CacheServiceImpl struct {
	cache *ristretto.Cache[string, []byte]
}

func NewCacheService() (CacheService, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}

	return &CacheServiceImpl{
		cache: cache,
	}, nil
}

func (c *CacheServiceImpl) Set(key string, value interface{}, cost int64) bool {
	data, err := json.Marshal(value)
	if err != nil {
		return false
	}
	return c.cache.Set(key, data, cost)
}

func (c *CacheServiceImpl) Get(key string) (interface{}, bool) {
	data, found := c.cache.Get(key)
	if !found {
		return nil, false
	}

	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, false
	}

	return result, true
}

func (c *CacheServiceImpl) Del(key string) {
	c.cache.Del(key)
}

// User Profile Cache Methods
func (c *CacheServiceImpl) GetUserProfile(userID uuid.UUID) (interface{}, bool) {
	key := fmt.Sprintf("user_profile:%s", userID.String())
	return c.Get(key)
}

func (c *CacheServiceImpl) SetUserProfile(userID uuid.UUID, profile interface{}) bool {
	key := fmt.Sprintf("user_profile:%s", userID.String())
	return c.Set(key, profile, 1024) // 1KB cost
}

// Task Cache Methods
func (c *CacheServiceImpl) GetTask(taskID uuid.UUID) (interface{}, bool) {
	key := fmt.Sprintf("task:%s", taskID.String())
	return c.Get(key)
}

func (c *CacheServiceImpl) SetTask(taskID uuid.UUID, task interface{}) bool {
	key := fmt.Sprintf("task:%s", taskID.String())
	return c.Set(key, task, 2048) // 2KB cost
}

// User Tasks Cache Methods
func (c *CacheServiceImpl) GetUserTasks(userID uuid.UUID) (interface{}, bool) {
	key := fmt.Sprintf("user_tasks:%s", userID.String())
	return c.Get(key)
}

func (c *CacheServiceImpl) SetUserTasks(userID uuid.UUID, tasks interface{}) bool {
	key := fmt.Sprintf("user_tasks:%s", userID.String())
	return c.Set(key, tasks, 5120) // 5KB cost
}

// Cache Invalidation Methods
func (c *CacheServiceImpl) InvalidateUserCache(userID uuid.UUID) {
	profileKey := fmt.Sprintf("user_profile:%s", userID.String())
	tasksKey := fmt.Sprintf("user_tasks:%s", userID.String())

	c.Del(profileKey)
	c.Del(tasksKey)
}

func (c *CacheServiceImpl) InvalidateTaskCache(taskID uuid.UUID) {
	taskKey := fmt.Sprintf("task:%s", taskID.String())
	c.Del(taskKey)

	// Also invalidate any user task lists since they might contain this task
	// In a real application, you might want to track which user lists contain this task
}

// Cache warming and TTL utilities
func (c *CacheServiceImpl) SetWithTTL(key string, value interface{}, cost int64, ttl time.Duration) bool {
	// Store with timestamp for TTL calculation
	cacheItem := struct {
		Value     interface{} `json:"value"`
		ExpiresAt time.Time   `json:"expires_at"`
	}{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	itemData, err := json.Marshal(cacheItem)
	if err != nil {
		return false
	}

	return c.cache.Set(key, itemData, cost)
}

func (c *CacheServiceImpl) GetWithTTL(key string) (interface{}, bool) {
	data, found := c.cache.Get(key)
	if !found {
		return nil, false
	}

	var cacheItem struct {
		Value     interface{} `json:"value"`
		ExpiresAt time.Time   `json:"expires_at"`
	}

	err := json.Unmarshal(data, &cacheItem)
	if err != nil {
		return nil, false
	}

	// Check if expired
	if time.Now().After(cacheItem.ExpiresAt) {
		c.Del(key)
		return nil, false
	}

	return cacheItem.Value, true
}
