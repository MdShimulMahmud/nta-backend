package services

import (
	"fmt"

	"github.com/sharetrip/nta-backend/internal/config"
	"github.com/sharetrip/nta-backend/internal/logger"
)

// Service handles business logic for system services
type Service struct {
	repo   *Repository
	config *config.Config
}

// NewService creates a new services service
func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
	}
}

// GetServiceStatus gets the status of a service
func (s *Service) GetServiceStatus(serviceName string) (*ServiceStatus, error) {
	logger.Debugf("Getting status for service: %s", serviceName)
	return s.repo.GetServiceStatus(serviceName)
}

// PerformServiceAction performs an action on a service
func (s *Service) PerformServiceAction(serviceName string, action ServiceAction) error {
	logger.Infof("Performing action %s on service %s", action, serviceName)

	switch action {
	case ActionStart:
		return s.repo.StartService(serviceName)
	case ActionStop:
		return s.repo.StopService(serviceName)
	case ActionRestart:
		return s.repo.RestartService(serviceName)
	case ActionReload:
		return s.repo.ReloadService(serviceName)
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

// Logstash-specific methods

// CheckLogstashConfig checks logstash configuration
func (s *Service) CheckLogstashConfig() (bool, string, error) {
	logger.Debug("Checking Logstash configuration")
	return s.repo.CheckLogstashConfig(s.config.Services.LogstashPath + "/conf.d/")
}

// GetLogstashPipelines returns all logstash pipelines
func (s *Service) GetLogstashPipelines() ([]Pipeline, error) {
	logger.Debug("Getting Logstash pipelines")
	return s.repo.GetLogstashPipelines()
}

// ToggleLogstashPipeline enables or disables a pipeline
func (s *Service) ToggleLogstashPipeline(pipelineID string, enabled bool) error {
	logger.Infof("Toggling pipeline %s to %v", pipelineID, enabled)

	if err := s.repo.ToggleLogstashPipeline(pipelineID, enabled); err != nil {
		return err
	}

	// Restart logstash to apply changes
	return s.repo.RestartService("logstash")
}

// GetLogstashLogs gets logstash logs
func (s *Service) GetLogstashLogs(lines int) ([]LogEntry, error) {
	logger.Debugf("Getting last %d Logstash log entries", lines)
	return s.repo.TailServiceLogs("logstash", lines)
}

// Suricata-specific methods

// CheckSuricataConfig checks suricata configuration
func (s *Service) CheckSuricataConfig() (bool, string, error) {
	logger.Debug("Checking Suricata configuration")
	return s.repo.CheckSuricataConfig(s.config.Services.SuricataPath)
}

// GetSuricataInterface gets the current monitoring interface
func (s *Service) GetSuricataInterface() (string, error) {
	logger.Debug("Getting Suricata interface")
	return s.repo.GetSuricataInterface(s.config.Services.SuricataPath)
}

// UpdateSuricataInterface updates the monitoring interface
func (s *Service) UpdateSuricataInterface(newInterface string) error {
	logger.Infof("Updating Suricata interface to %s", newInterface)

	// Get current interface
	currentInterface, err := s.repo.GetSuricataInterface(s.config.Services.SuricataPath)
	if err != nil {
		return err
	}

	// Update interface
	if err := s.repo.UpdateSuricataInterface(s.config.Services.SuricataPath, currentInterface, newInterface); err != nil {
		return err
	}

	// Restart Suricata
	return s.repo.RestartService("suricata")
}

// GetSuricataHomeNet gets the HOME_NET configuration
func (s *Service) GetSuricataHomeNet() ([]string, error) {
	logger.Debug("Getting Suricata HOME_NET")
	return s.repo.GetSuricataHomeNet(s.config.Services.SuricataPath)
}

// UpdateSuricataHomeNet updates the HOME_NET configuration
func (s *Service) UpdateSuricataHomeNet(networks []string) error {
	logger.Infof("Updating Suricata HOME_NET to %v", networks)

	if err := s.repo.UpdateSuricataHomeNet(s.config.Services.SuricataPath, networks); err != nil {
		return err
	}

	// Restart Suricata
	return s.repo.RestartService("suricata")
}

// UpdateSuricataRules updates Suricata rules
func (s *Service) UpdateSuricataRules() error {
	logger.Info("Updating Suricata rules")
	return s.repo.UpdateSuricataRules()
}

// GetSuricataAlerts gets recent alerts
func (s *Service) GetSuricataAlerts(limit int) ([]Alert, error) {
	logger.Debugf("Getting last %d Suricata alerts", limit)
	return s.repo.GetSuricataAlerts(limit)
}

// GetSuricataLogs gets suricata logs
func (s *Service) GetSuricataLogs(lines int) ([]LogEntry, error) {
	logger.Debugf("Getting last %d Suricata log entries", lines)
	return s.repo.TailServiceLogs("suricata", lines)
}

// RPort-specific methods

// CheckRPortConnection checks RPort connection status
func (s *Service) CheckRPortConnection() (*ConnectionStatus, error) {
	logger.Debug("Checking RPort connection status")
	return s.repo.CheckRPortConnection()
}

// GetRPortLogs gets rport logs
func (s *Service) GetRPortLogs(lines int) ([]LogEntry, error) {
	logger.Debugf("Getting last %d RPort log entries", lines)
	return s.repo.TailServiceLogs("rport", lines)
}

// Rsyslog-specific methods

// CheckRsyslogConfig checks rsyslog configuration status
func (s *Service) CheckRsyslogConfig() (*RsyslogConfigStatus, error) {
	logger.Debug("Checking Rsyslog configuration")
	return s.repo.CheckRsyslogConfig()
}

// ToggleRsyslogFirewallConfig enables or disables firewall log collection
func (s *Service) ToggleRsyslogFirewallConfig(enabled bool) error {
	logger.Infof("Toggling Rsyslog firewall config to %v", enabled)
	return s.repo.ToggleRsyslogConfig(enabled)
}

// GetRsyslogFirewallLogs gets firewall logs
func (s *Service) GetRsyslogFirewallLogs(lines int) ([]LogEntry, error) {
	logger.Debugf("Getting last %d firewall log entries", lines)
	return s.repo.TailRsyslogFirewallLogs(lines)
}

// GetRsyslogLogs gets rsyslog service logs
func (s *Service) GetRsyslogLogs(lines int) ([]LogEntry, error) {
	logger.Debugf("Getting last %d Rsyslog log entries", lines)
	return s.repo.TailServiceLogs("rsyslog", lines)
}

// Wazuh-specific methods

// GetWazuhLogs gets wazuh agent logs
func (s *Service) GetWazuhLogs(lines int) ([]LogEntry, error) {
	logger.Debugf("Getting last %d Wazuh log entries", lines)
	return s.repo.TailWazuhLogs(lines)
}
