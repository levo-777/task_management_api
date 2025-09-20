package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter represents a rate limiter for a specific IP
type RateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitManager manages rate limiters for different IPs
type RateLimitManager struct {
	limiters map[string]*RateLimiter
	mutex    sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimitManager creates a new rate limit manager
func NewRateLimitManager(r rate.Limit, b int) *RateLimitManager {
	rl := &RateLimitManager{
		limiters: make(map[string]*RateLimiter),
		rate:     r,
		burst:    b,
	}

	// Clean up old limiters every minute
	go rl.cleanupLimiters()

	return rl
}

// GetLimiter returns a rate limiter for the given IP
func (rlm *RateLimitManager) GetLimiter(ip string) *rate.Limiter {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	limiter, exists := rlm.limiters[ip]
	if !exists {
		limiter = &RateLimiter{
			limiter:  rate.NewLimiter(rlm.rate, rlm.burst),
			lastSeen: time.Now(),
		}
		rlm.limiters[ip] = limiter
	}

	limiter.lastSeen = time.Now()
	return limiter.limiter
}

// cleanupLimiters removes old limiters to prevent memory leaks
func (rlm *RateLimitManager) cleanupLimiters() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rlm.mutex.Lock()
		for ip, limiter := range rlm.limiters {
			if time.Since(limiter.lastSeen) > 3*time.Minute {
				delete(rlm.limiters, ip)
			}
		}
		rlm.mutex.Unlock()
	}
}

// Global rate limit manager
var globalRateLimiter = NewRateLimitManager(rate.Every(time.Second), 10) // 10 requests per second

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := globalRateLimiter.GetLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware creates a stricter rate limiter for auth endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	authRateLimiter := NewRateLimitManager(rate.Every(2*time.Second), 3) // 3 requests per 2 seconds
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := authRateLimiter.GetLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many authentication attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
