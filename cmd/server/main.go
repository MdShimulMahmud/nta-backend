package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/api"
	"github.com/sharetrip/nta-backend/internal/config"
	"github.com/sharetrip/nta-backend/internal/interfaces"
	"github.com/sharetrip/nta-backend/internal/logger"
	"github.com/sharetrip/nta-backend/internal/middleware"
	"github.com/sharetrip/nta-backend/internal/services"
)

// @title           NTA Backend API
// @version         1.0
// @description     Network Threat Analysis Backend API for managing network interfaces and system services
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@amplifysec.net

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg)
	logger.Infof("Starting %s in %s mode", cfg.App.Name, cfg.App.Env)

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin engine
	engine := gin.New()

	// Recovery middleware
	engine.Use(gin.Recovery())

	// CORS middleware
	engine.Use(middleware.CORS(cfg))

	// Rate limiting (optional)
	if cfg.App.Env == "production" {
		rateLimiter := middleware.NewRateLimiter(100, time.Minute)
		engine.Use(rateLimiter.Middleware())
	}

	// Initialize repositories
	interfacesRepo := interfaces.NewRepository()
	servicesRepo := services.NewRepository()

	// Initialize services
	interfacesService := interfaces.NewService(interfacesRepo, cfg)
	servicesService := services.NewService(servicesRepo, cfg)

	// Initialize handlers
	interfacesHandler := interfaces.NewHandler(interfacesService)
	servicesHandler := services.NewHandler(servicesService)

	// Setup routes
	router := api.NewRouter(cfg, interfacesHandler, servicesHandler)
	router.SetupRoutes(engine)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on %s", addr)
		logger.Infof("API documentation available at http://%s/swagger/index.html", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped")
}