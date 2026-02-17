package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/config"
)

// CORS configures CORS middleware
func CORS(cfg *config.Config) gin.HandlerFunc {
	// Debug log for AllowedOrigins
	println("CORS_ALLOWED_ORIGINS:", cfg.CORS.AllowedOrigins)

	corsConfig := cors.DefaultConfig()

	// If any value in AllowedOrigins is '*', allow all origins
	allowAll := false
	for _, origin := range cfg.CORS.AllowedOrigins {
		if origin == "*" {
			allowAll = true
			break
		}
	}
	if allowAll {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = false
	} else {
		corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins
		corsConfig.AllowCredentials = true
	}

	if len(cfg.CORS.AllowedMethods) > 0 {
		corsConfig.AllowMethods = cfg.CORS.AllowedMethods
	}

	if len(cfg.CORS.AllowedHeaders) > 0 && cfg.CORS.AllowedHeaders[0] == "*" {
		corsConfig.AllowHeaders = []string{"*"}
	} else {
		corsConfig.AllowHeaders = cfg.CORS.AllowedHeaders
	}

	return cors.New(corsConfig)
}
