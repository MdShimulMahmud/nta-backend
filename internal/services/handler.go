package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/logger"
)

// Handler handles HTTP requests for system services
type Handler struct {
	service *Service
}

// NewHandler creates a new services handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Common service endpoints

// GetServiceStatus godoc
// @Summary      Get service status
// @Description  Get the current status of a system service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "Service name"
// @Success      200   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /services/{name}/status [get]
func (h *Handler) GetServiceStatus(c *gin.Context) {
	serviceName := c.Param("name")

	status, err := h.service.GetServiceStatus(serviceName)
	if err != nil {
		logger.Errorf("Failed to get service status for %s: %v", serviceName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get service status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  status,
	})
}

// PerformServiceAction godoc
// @Summary      Perform service action
// @Description  Perform an action (start, stop, restart, reload) on a service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name    path      string                 true  "Service name"
// @Param        request body      ServiceActionRequest   true  "Action request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /services/{name}/action [post]
func (h *Handler) PerformServiceAction(c *gin.Context) {
	serviceName := c.Param("name")

	var req ServiceActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.PerformServiceAction(serviceName, req.Action); err != nil {
		logger.Errorf("Failed to perform action %s on %s: %v", req.Action, serviceName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to perform service action",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action performed successfully",
		"action":  req.Action,
		"service": serviceName,
	})
}

// Logstash endpoints

// CheckLogstashConfig godoc
// @Summary      Check Logstash configuration
// @Description  Test Logstash configuration syntax
// @Tags         logstash
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/logstash/config/check [get]
func (h *Handler) CheckLogstashConfig(c *gin.Context) {
	valid, message, err := h.service.CheckLogstashConfig()
	if err != nil {
		logger.Error("Failed to check Logstash config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"valid":   valid,
		"message": message,
	})
}

// GetLogstashPipelines godoc
// @Summary      Get Logstash pipelines
// @Description  List all Logstash pipelines and their status
// @Tags         logstash
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/logstash/pipelines [get]
func (h *Handler) GetLogstashPipelines(c *gin.Context) {
	pipelines, err := h.service.GetLogstashPipelines()
	if err != nil {
		logger.Error("Failed to get Logstash pipelines:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pipelines",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"count":     len(pipelines),
		"pipelines": pipelines,
	})
}

// ToggleLogstashPipeline godoc
// @Summary      Toggle Logstash pipeline
// @Description  Enable or disable a Logstash pipeline
// @Tags         logstash
// @Accept       json
// @Produce      json
// @Param        request body      PipelineToggleRequest  true  "Pipeline toggle request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /services/logstash/pipelines/toggle [post]
func (h *Handler) ToggleLogstashPipeline(c *gin.Context) {
	var req PipelineToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.ToggleLogstashPipeline(req.PipelineID, req.Enabled); err != nil {
		logger.Errorf("Failed to toggle pipeline %s: %v", req.PipelineID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to toggle pipeline",
		})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pipeline " + status,
	})
}

// GetLogstashLogs godoc
// @Summary      Get Logstash logs
// @Description  Get recent Logstash service logs
// @Tags         logstash
// @Accept       json
// @Produce      json
// @Param        lines  query     int  false  "Number of lines (default: 50)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/logstash/logs [get]
func (h *Handler) GetLogstashLogs(c *gin.Context) {
	lines := 50
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil {
			lines = l
		}
	}

	logs, err := h.service.GetLogstashLogs(lines)
	if err != nil {
		logger.Error("Failed to get Logstash logs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}

// Suricata endpoints

// CheckSuricataConfig godoc
// @Summary      Check Suricata configuration
// @Description  Test Suricata configuration syntax
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/suricata/config/check [get]
func (h *Handler) CheckSuricataConfig(c *gin.Context) {
	valid, message, err := h.service.CheckSuricataConfig()
	if err != nil {
		logger.Error("Failed to check Suricata config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"valid":   valid,
		"message": message,
	})
}

// GetSuricataInterface godoc
// @Summary      Get Suricata monitoring interface
// @Description  Get the current network interface being monitored by Suricata
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/suricata/interface [get]
func (h *Handler) GetSuricataInterface(c *gin.Context) {
	iface, err := h.service.GetSuricataInterface()
	if err != nil {
		logger.Error("Failed to get Suricata interface:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get interface",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"interface": iface,
	})
}

// UpdateSuricataInterface godoc
// @Summary      Update Suricata monitoring interface
// @Description  Change the network interface monitored by Suricata
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Param        request body      UpdateSuricataInterfaceRequest  true  "Interface update request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /services/suricata/interface [put]
func (h *Handler) UpdateSuricataInterface(c *gin.Context) {
	var req UpdateSuricataInterfaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateSuricataInterface(req.Interface); err != nil {
		logger.Error("Failed to update Suricata interface:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update interface",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Interface updated successfully",
	})
}

// GetSuricataHomeNet godoc
// @Summary      Get Suricata HOME_NET configuration
// @Description  Get the current HOME_NET configuration
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/suricata/home-net [get]
func (h *Handler) GetSuricataHomeNet(c *gin.Context) {
	networks, err := h.service.GetSuricataHomeNet()
	if err != nil {
		logger.Error("Failed to get HOME_NET:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get HOME_NET",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"networks": networks,
	})
}

// UpdateSuricataHomeNet godoc
// @Summary      Update Suricata HOME_NET configuration
// @Description  Update the HOME_NET configuration
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Param        request body      UpdateSuricataHomeNetRequest  true  "HOME_NET update request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /services/suricata/home-net [put]
func (h *Handler) UpdateSuricataHomeNet(c *gin.Context) {
	var req UpdateSuricataHomeNetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateSuricataHomeNet(req.Networks); err != nil {
		logger.Error("Failed to update HOME_NET:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update HOME_NET",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "HOME_NET updated successfully",
	})
}

// UpdateSuricataRules godoc
// @Summary      Update Suricata rules
// @Description  Update Suricata detection rules
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/suricata/rules/update [post]
func (h *Handler) UpdateSuricataRules(c *gin.Context) {
	if err := h.service.UpdateSuricataRules(); err != nil {
		logger.Error("Failed to update Suricata rules:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update rules",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rules updated successfully",
	})
}

// GetSuricataAlerts godoc
// @Summary      Get Suricata alerts
// @Description  Get recent Suricata alerts
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "Number of alerts (default: 100)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/suricata/alerts [get]
func (h *Handler) GetSuricataAlerts(c *gin.Context) {
	limit := 100
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	alerts, err := h.service.GetSuricataAlerts(limit)
	if err != nil {
		logger.Error("Failed to get Suricata alerts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get alerts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(alerts),
		"alerts":  alerts,
	})
}

// GetSuricataLogs godoc
// @Summary      Get Suricata logs
// @Description  Get recent Suricata service logs
// @Tags         suricata
// @Accept       json
// @Produce      json
// @Param        lines  query     int  false  "Number of lines (default: 50)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/suricata/logs [get]
func (h *Handler) GetSuricataLogs(c *gin.Context) {
	lines := 50
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil {
			lines = l
		}
	}

	logs, err := h.service.GetSuricataLogs(lines)
	if err != nil {
		logger.Error("Failed to get Suricata logs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}

// RPort endpoints

// CheckRPortConnection godoc
// @Summary      Check RPort connection
// @Description  Check RPort connection status to remote server
// @Tags         rport
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/rport/connection [get]
func (h *Handler) CheckRPortConnection(c *gin.Context) {
	status, err := h.service.CheckRPortConnection()
	if err != nil {
		logger.Error("Failed to check RPort connection:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check connection",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  status,
	})
}

// GetRPortLogs godoc
// @Summary      Get RPort logs
// @Description  Get recent RPort service logs
// @Tags         rport
// @Accept       json
// @Produce      json
// @Param        lines  query     int  false  "Number of lines (default: 50)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/rport/logs [get]
func (h *Handler) GetRPortLogs(c *gin.Context) {
	lines := 50
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil {
			lines = l
		}
	}

	logs, err := h.service.GetRPortLogs(lines)
	if err != nil {
		logger.Error("Failed to get RPort logs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}

// Rsyslog endpoints

// CheckRsyslogConfig godoc
// @Summary      Check Rsyslog configuration
// @Description  Check Rsyslog firewall configuration status
// @Tags         rsyslog
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /services/rsyslog/config/check [get]
func (h *Handler) CheckRsyslogConfig(c *gin.Context) {
	status, err := h.service.CheckRsyslogConfig()
	if err != nil {
		logger.Error("Failed to check Rsyslog config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check configuration",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  status,
	})
}

// ToggleRsyslogConfig godoc
// @Summary      Toggle Rsyslog firewall configuration
// @Description  Enable or disable Rsyslog firewall log collection
// @Tags         rsyslog
// @Accept       json
// @Produce      json
// @Param        request body      map[string]bool  true  "Enable/Disable request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /services/rsyslog/config/toggle [post]
func (h *Handler) ToggleRsyslogConfig(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.ToggleRsyslogFirewallConfig(req.Enabled); err != nil {
		logger.Error("Failed to toggle Rsyslog config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to toggle configuration",
		})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Firewall log collection " + status,
	})
}

// GetRsyslogFirewallLogs godoc
// @Summary      Get Rsyslog firewall logs
// @Description  Get recent firewall logs
// @Tags         rsyslog
// @Accept       json
// @Produce      json
// @Param        lines  query     int  false  "Number of lines (default: 50)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/rsyslog/firewall-logs [get]
func (h *Handler) GetRsyslogFirewallLogs(c *gin.Context) {
	lines := 50
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil {
			lines = l
		}
	}

	logs, err := h.service.GetRsyslogFirewallLogs(lines)
	if err != nil {
		logger.Error("Failed to get firewall logs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}

// GetRsyslogLogs godoc
// @Summary      Get Rsyslog service logs
// @Description  Get recent Rsyslog service logs
// @Tags         rsyslog
// @Accept       json
// @Produce      json
// @Param        lines  query     int  false  "Number of lines (default: 50)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/rsyslog/logs [get]
func (h *Handler) GetRsyslogLogs(c *gin.Context) {
	lines := 50
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil {
			lines = l
		}
	}

	logs, err := h.service.GetRsyslogLogs(lines)
	if err != nil {
		logger.Error("Failed to get Rsyslog logs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}

// Wazuh endpoints

// GetWazuhLogs godoc
// @Summary      Get Wazuh agent logs
// @Description  Get recent Wazuh agent logs
// @Tags         wazuh
// @Accept       json
// @Produce      json
// @Param        lines  query     int  false  "Number of lines (default: 50)"
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /services/wazuh/logs [get]
func (h *Handler) GetWazuhLogs(c *gin.Context) {
	lines := 50
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil {
			lines = l
		}
	}

	logs, err := h.service.GetWazuhLogs(lines)
	if err != nil {
		logger.Error("Failed to get Wazuh logs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(logs),
		"logs":    logs,
	})
}
