package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/logger"
)

// ErrorHandler handles errors and panics
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors from handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Errorf("Request error: %v", err.Err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}
}
