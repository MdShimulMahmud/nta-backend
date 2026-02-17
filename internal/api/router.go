package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/config"
	"github.com/sharetrip/nta-backend/internal/interfaces"
	"github.com/sharetrip/nta-backend/internal/middleware"
	"github.com/sharetrip/nta-backend/internal/services"

	// Swagger
	docs "github.com/sharetrip/nta-backend/docs/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router sets up the API routes
type Router struct {
	config             *config.Config
	interfacesHandler  *interfaces.Handler
	servicesHandler    *services.Handler
}

// NewRouter creates a new API router
func NewRouter(
	cfg *config.Config,
	interfacesHandler *interfaces.Handler,
	servicesHandler *services.Handler,
) *Router {
	return &Router{
		config:            cfg,
		interfacesHandler: interfacesHandler,
		servicesHandler:   servicesHandler,
	}
}

// SetupRoutes sets up all API routes
func (r *Router) SetupRoutes(engine *gin.Engine) {
	// Swagger docs info (optional, but good for dynamic host)
	docs.SwaggerInfo.Host = r.config.App.Host + ":" + r.config.App.Port
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Root endpoint
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "NTA Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Swagger UI
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := engine.Group(r.config.API.Prefix)
	{
		// Health check (now under /api/v1/health)
		v1.GET("/health", HealthCheck)
		// Apply middlewares
		v1.Use(middleware.RequestLogger())
		v1.Use(middleware.ErrorHandler())

		// Network interfaces routes
		interfacesGroup := v1.Group("/interfaces")
		{
			interfacesGroup.GET("", r.interfacesHandler.ListInterfaces)
			interfacesGroup.GET("/:name", r.interfacesHandler.GetInterface)
			interfacesGroup.PUT("/:name/state", r.interfacesHandler.UpdateInterfaceState)
			interfacesGroup.PUT("/:name/promiscuous", r.interfacesHandler.UpdatePromiscuousMode)
			interfacesGroup.GET("/:name/check-network-logs", r.interfacesHandler.CheckNetworkLogs)
			interfacesGroup.GET("/:name/check-firewall-logs", r.interfacesHandler.CheckFirewallLogs)
			
			// Netplan routes
			interfacesGroup.GET("/netplan", r.interfacesHandler.GetNetplanConfig)
			interfacesGroup.PUT("/netplan", r.interfacesHandler.UpdateNetplanConfig)
		}

		// System services routes
		servicesGroup := v1.Group("/services")
		{
			// Common service endpoints
			servicesGroup.GET("/:name/status", r.servicesHandler.GetServiceStatus)
			servicesGroup.POST("/:name/action", r.servicesHandler.PerformServiceAction)

			// Logstash routes
			logstashGroup := servicesGroup.Group("/logstash")
			{
				logstashGroup.GET("/config/check", r.servicesHandler.CheckLogstashConfig)
				logstashGroup.GET("/pipelines", r.servicesHandler.GetLogstashPipelines)
				logstashGroup.POST("/pipelines/toggle", r.servicesHandler.ToggleLogstashPipeline)
				logstashGroup.GET("/logs", r.servicesHandler.GetLogstashLogs)
			}

			// Suricata routes
			suricataGroup := servicesGroup.Group("/suricata")
			{
				suricataGroup.GET("/config/check", r.servicesHandler.CheckSuricataConfig)
				suricataGroup.GET("/interface", r.servicesHandler.GetSuricataInterface)
				suricataGroup.PUT("/interface", r.servicesHandler.UpdateSuricataInterface)
				suricataGroup.GET("/home-net", r.servicesHandler.GetSuricataHomeNet)
				suricataGroup.PUT("/home-net", r.servicesHandler.UpdateSuricataHomeNet)
				suricataGroup.POST("/rules/update", r.servicesHandler.UpdateSuricataRules)
				suricataGroup.GET("/alerts", r.servicesHandler.GetSuricataAlerts)
				suricataGroup.GET("/logs", r.servicesHandler.GetSuricataLogs)
			}

			// RPort routes
			rportGroup := servicesGroup.Group("/rport")
			{
				rportGroup.GET("/connection", r.servicesHandler.CheckRPortConnection)
				rportGroup.GET("/logs", r.servicesHandler.GetRPortLogs)
			}

			// Rsyslog routes
			rsyslogGroup := servicesGroup.Group("/rsyslog")
			{
				rsyslogGroup.GET("/config/check", r.servicesHandler.CheckRsyslogConfig)
				rsyslogGroup.POST("/config/toggle", r.servicesHandler.ToggleRsyslogConfig)
				rsyslogGroup.GET("/firewall-logs", r.servicesHandler.GetRsyslogFirewallLogs)
				rsyslogGroup.GET("/logs", r.servicesHandler.GetRsyslogLogs)
			}

			// Wazuh routes
			wazuhGroup := servicesGroup.Group("/wazuh")
			{
				wazuhGroup.GET("/logs", r.servicesHandler.GetWazuhLogs)
			}
		}
	}
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the API is running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"message": "NTA Backend is running",
	})
}
