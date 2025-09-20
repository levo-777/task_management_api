package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomLoggerConfig defines the configuration for custom logging
type CustomLoggerConfig struct {
	SkipPaths []string
}

// CustomLogger creates a custom logging middleware
func CustomLogger(config CustomLoggerConfig) gin.HandlerFunc {
	skipMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipMap[path] = true
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Custom log format with additional information
			return fmt.Sprintf("[%s] %s %s %d %s %s %s %s %s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.ClientIP,
				param.Method,
				param.StatusCode,
				param.Latency,
				param.Path,
				param.Request.UserAgent(),
				param.ErrorMessage,
				param.BodySize,
			)
		},
		Output: gin.DefaultWriter,
		SkipPaths: config.SkipPaths,
	})
}

// RequestLogger logs detailed request information
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Log request start
		gin.DefaultWriter.Write([]byte(
			fmt.Sprintf("[REQUEST START] %s %s %s %s\n",
				start.Format("2006/01/02 15:04:05"),
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
			),
		))

		c.Next()

		// Log request completion
		latency := time.Since(start)
		gin.DefaultWriter.Write([]byte(
			fmt.Sprintf("[REQUEST END] %s %s %s %d %v %s\n",
				time.Now().Format("2006/01/02 15:04:05"),
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				latency,
				c.ClientIP(),
			),
		))
	}
}

// ErrorLogger logs errors with stack traces
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors in the response
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				gin.DefaultErrorWriter.Write([]byte(
					fmt.Sprintf("[ERROR] %s %s %s %s\n",
						time.Now().Format("2006/01/02 15:04:05"),
						c.Request.Method,
						c.Request.URL.Path,
						err.Error(),
					),
				))
			}
		}
	}
}
