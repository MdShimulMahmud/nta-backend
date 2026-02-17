package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/logger"
)

// RequestLogger logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.Infof(
			"[%s] %s %s - Status: %d - Duration: %v - IP: %s",
			method,
			path,
			c.Request.Proto,
			statusCode,
			duration,
			clientIP,
		)
	}
}
