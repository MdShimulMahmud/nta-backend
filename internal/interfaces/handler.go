package interfaces

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sharetrip/nta-backend/internal/logger"
)

// Handler handles HTTP requests for network interfaces
type Handler struct {
	service *Service
}

// NewHandler creates a new interface handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ListInterfaces godoc
// @Summary      List all network interfaces
// @Description  Get a list of all available network interfaces with their details
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /interfaces [get]
func (h *Handler) ListInterfaces(c *gin.Context) {
	interfaces, err := h.service.GetAllInterfaces()
	if err != nil {
		logger.Error("Failed to list interfaces:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve network interfaces",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"count":      len(interfaces),
		"interfaces": interfaces,
	})
}

// GetInterface godoc
// @Summary      Get network interface details
// @Description  Get detailed information about a specific network interface
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "Interface name"
// @Success      200   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /interfaces/{name} [get]
func (h *Handler) GetInterface(c *gin.Context) {
	name := c.Param("name")

	iface, err := h.service.GetInterface(name)
	if err != nil {
		logger.Errorf("Failed to get interface %s: %v", name, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Interface not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"interface": iface,
	})
}

// UpdateInterfaceState godoc
// @Summary      Update interface state
// @Description  Bring an interface up or down
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Param        name    path      string                       true  "Interface name"
// @Param        request body      UpdateInterfaceStateRequest  true  "State update request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /interfaces/{name}/state [put]
func (h *Handler) UpdateInterfaceState(c *gin.Context) {
	name := c.Param("name")

	var req UpdateInterfaceStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateInterfaceState(name, req.State); err != nil {
		logger.Errorf("Failed to update interface %s state: %v", name, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update interface state",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Interface %s is now %s", name, req.State),
	})
}

// UpdatePromiscuousMode godoc
// @Summary      Update promiscuous mode
// @Description  Enable or disable promiscuous mode on an interface
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Param        name    path      string                     true  "Interface name"
// @Param        request body      UpdatePromiscModeRequest   true  "Promiscuous mode update request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /interfaces/{name}/promiscuous [put]
func (h *Handler) UpdatePromiscuousMode(c *gin.Context) {
	name := c.Param("name")

	var req UpdatePromiscModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdatePromiscuousMode(name, req.Enabled); err != nil {
		logger.Errorf("Failed to update promiscuous mode on %s: %v", name, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update promiscuous mode",
		})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Promiscuous mode %s on %s", status, name),
	})
}

// GetNetplanConfig godoc
// @Summary      Get netplan configuration
// @Description  Retrieve the current netplan configuration
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /interfaces/netplan [get]
func (h *Handler) GetNetplanConfig(c *gin.Context) {
	config, err := h.service.GetNetplanConfig()
	if err != nil {
		logger.Error("Failed to get netplan config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve netplan configuration",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"config":  config,
	})
}

// UpdateNetplanConfig godoc
// @Summary      Update netplan configuration
// @Description  Update netplan configuration for a specific interface
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Param        request body      UpdateNetplanRequest  true  "Netplan update request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /interfaces/netplan [put]
func (h *Handler) UpdateNetplanConfig(c *gin.Context) {
	var req UpdateNetplanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateNetplanInterface(&req); err != nil {
		logger.Error("Failed to update netplan config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update netplan configuration",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Netplan configuration updated successfully",
	})
}

// CheckNetworkLogs godoc
// @Summary      Check network logs with tcpdump
// @Description  Check if network logs are being received on an interface using tcpdump
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "Interface name"
// @Success      200   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /interfaces/{name}/check-network-logs [get]
func (h *Handler) CheckNetworkLogs(c *gin.Context) {
	name := c.Param("name")

	result, err := h.service.CheckTCPDumpNetwork(name)
	if err != nil {
		logger.Errorf("Failed to check network logs on %s: %v", name, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check network logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// CheckFirewallLogs godoc
// @Summary      Check firewall logs with tcpdump
// @Description  Check if firewall logs are being received on an interface using tcpdump (port 514)
// @Tags         interfaces
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "Interface name"
// @Success      200   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /interfaces/{name}/check-firewall-logs [get]
func (h *Handler) CheckFirewallLogs(c *gin.Context) {
	name := c.Param("name")

	result, err := h.service.CheckTCPDumpFirewall(name)
	if err != nil {
		logger.Errorf("Failed to check firewall logs on %s: %v", name, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check firewall logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}
