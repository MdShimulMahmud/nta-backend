package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/sharetrip/nta-backend/internal/logger"
)

// Repository handles system service operations
type Repository struct{}

// NewRepository creates a new services repository
func NewRepository() *Repository {
	return &Repository{}
}

// GetServiceStatus gets the status of a systemd service
func (r *Repository) GetServiceStatus(serviceName string) (*ServiceStatus, error) {
	cmd := exec.Command("systemctl", "status", serviceName)
	output, _ := cmd.CombinedOutput()

	status := &ServiceStatus{
		Name:    serviceName,
		Details: make(map[string]interface{}),
	}
	outputStr := string(output)

	// Check if active
	status.Active = strings.Contains(outputStr, "Active: active")
	status.Running = strings.Contains(outputStr, "active (running)")
	status.Enabled = r.isServiceEnabled(serviceName)

	// Parse status
	if status.Active {
		status.Status = "active"
	} else if strings.Contains(outputStr, "Active: inactive") {
		status.Status = "inactive"
	} else if strings.Contains(outputStr, "Active: failed") {
		status.Status = "failed"
	} else {
		status.Status = "unknown"
	}

	return status, nil
}

// isServiceEnabled checks if a service is enabled
func (r *Repository) isServiceEnabled(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-enabled", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "enabled"
}

// StartService starts a systemd service
func (r *Repository) StartService(serviceName string) error {
	cmd := exec.Command("systemctl", "start", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start service %s: %s", serviceName, output)
	}
	logger.Infof("Service %s started successfully", serviceName)
	return nil
}

// StopService stops a systemd service
func (r *Repository) StopService(serviceName string) error {
	cmd := exec.Command("systemctl", "stop", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop service %s: %s", serviceName, output)
	}
	logger.Infof("Service %s stopped successfully", serviceName)
	return nil
}

// RestartService restarts a systemd service
func (r *Repository) RestartService(serviceName string) error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart service %s: %s", serviceName, output)
	}
	logger.Infof("Service %s restarted successfully", serviceName)
	return nil
}

// ReloadService reloads a systemd service
func (r *Repository) ReloadService(serviceName string) error {
	cmd := exec.Command("systemctl", "reload", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload service %s: %s", serviceName, output)
	}
	logger.Infof("Service %s reloaded successfully", serviceName)
	return nil
}

// CheckLogstashConfig tests logstash configuration
func (r *Repository) CheckLogstashConfig(configPath string) (bool, string, error) {
	cmd := exec.Command("/usr/share/logstash/bin/logstash", "--config.test_and_exit", "-f", configPath)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return false, outputStr, fmt.Errorf("configuration test failed: %s", outputStr)
	}

	if strings.Contains(outputStr, "Configuration OK") {
		return true, "Configuration is valid", nil
	}
	return false, outputStr, nil
}

// GetLogstashPipelines gets running logstash pipelines
func (r *Repository) GetLogstashPipelines() ([]Pipeline, error) {
	// Try API first
	resp, err := http.Get("http://localhost:9600/_node/pipelines")
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
				if pipelines, ok := result["pipelines"].(map[string]interface{}); ok {
					var pipelineList []Pipeline
					for id := range pipelines {
						pipelineList = append(pipelineList, Pipeline{
							ID:      id,
							Enabled: true,
							Status:  "running",
						})
					}
					return pipelineList, nil
				}
			}
		}
	}

	// Fallback to config file parsing
	return r.parseLogstashPipelinesConfig()
}

// parseLogstashPipelinesConfig parses pipelines.yml
func (r *Repository) parseLogstashPipelinesConfig() ([]Pipeline, error) {
	data, err := os.ReadFile("/etc/logstash/pipelines.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to read pipelines.yml: %w", err)
	}

	var pipelines []Pipeline
	content := string(data)

	// Simple parsing for pipeline IDs
	re := regexp.MustCompile(`pipeline\.id:\s*([^\s#]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			pipelineID := strings.TrimSpace(match[1])
			enabled := !strings.Contains(content, fmt.Sprintf("#pipeline.id: %s", pipelineID))

			pipelines = append(pipelines, Pipeline{
				ID:      pipelineID,
				Enabled: enabled,
			})
		}
	}

	return pipelines, nil
}

// ToggleLogstashPipeline enables or disables a pipeline
func (r *Repository) ToggleLogstashPipeline(pipelineID string, enabled bool) error {
	pipelinesFile := "/etc/logstash/pipelines.yml"
	data, err := os.ReadFile(pipelinesFile)
	if err != nil {
		return fmt.Errorf("failed to read pipelines.yml: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var newLines []string
	inPipeline := false
	pipelineFound := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if we're entering the target pipeline section
		if strings.Contains(trimmed, fmt.Sprintf("pipeline.id: %s", pipelineID)) ||
			strings.Contains(trimmed, fmt.Sprintf("#pipeline.id: %s", pipelineID)) {
			inPipeline = true
			pipelineFound = true
			if enabled {
				// Remove comment if present
				line = strings.Replace(line, "#", "", 1)
			} else {
				// Add comment if not present
				if !strings.HasPrefix(trimmed, "#") {
					line = "# " + line
				}
			}
		} else if inPipeline && strings.Contains(trimmed, "pipeline.id:") {
			// We've reached another pipeline
			inPipeline = false
		}

		// If we're in the pipeline section, comment/uncomment lines
		if inPipeline && trimmed != "" && !strings.Contains(trimmed, "pipeline.id:") {
			if enabled {
				line = strings.Replace(line, "# ", "", 1)
				line = strings.Replace(line, "#", "", 1)
			} else {
				if !strings.HasPrefix(trimmed, "#") {
					// Preserve indentation
					spaces := len(line) - len(strings.TrimLeft(line, " \t"))
					indent := line[:spaces]
					line = indent + "# " + strings.TrimLeft(line, " \t")
				}
			}
		}

		newLines = append(newLines, line)
	}

	if !pipelineFound {
		return fmt.Errorf("pipeline %s not found in configuration", pipelineID)
	}

	// Write back
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(pipelinesFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write pipelines.yml: %w", err)
	}

	logger.Infof("Pipeline %s %s", pipelineID, map[bool]string{true: "enabled", false: "disabled"}[enabled])
	return nil
}

// CheckSuricataConfig tests suricata configuration
func (r *Repository) CheckSuricataConfig(configPath string) (bool, string, error) {
	cmd := exec.Command("suricata", "-T", "-c", configPath)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return false, outputStr, fmt.Errorf("configuration test failed: %s", outputStr)
	}

	return true, "Configuration is valid", nil
}

// GetSuricataInterface gets the current Suricata interface
func (r *Repository) GetSuricataInterface(configPath string) (string, error) {
	cmd := exec.Command("grep", "-A2", "af-packet:", configPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Suricata interface: %w", err)
	}

	// Parse interface from output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "interface:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("interface not found in configuration")
}

// UpdateSuricataInterface updates the Suricata monitoring interface
func (r *Repository) UpdateSuricataInterface(configPath, oldInterface, newInterface string) error {
	cmd := exec.Command("sed", "-i",
		fmt.Sprintf("s/interface: %s/interface: %s/", oldInterface, newInterface),
		configPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update Suricata interface: %s", output)
	}

	logger.Infof("Suricata interface updated from %s to %s", oldInterface, newInterface)
	return nil
}

// GetSuricataHomeNet gets the HOME_NET configuration
func (r *Repository) GetSuricataHomeNet(configPath string) ([]string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Suricata config: %w", err)
	}

	content := string(data)
	re := regexp.MustCompile(`HOME_NET:\s*"?\[([^\]]+)\]"?`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		var result []string
		networks := strings.Split(matches[1], ",")
		for _, net := range networks {
			result = append(result, strings.TrimSpace(net))
		}
		return result, nil
	}

	return []string{}, nil
}

// UpdateSuricataHomeNet updates the HOME_NET configuration
func (r *Repository) UpdateSuricataHomeNet(configPath string, networks []string) error {
	networksStr := strings.Join(networks, ",")
	cmd := exec.Command("sed", "-i",
		fmt.Sprintf("s/HOME_NET:.*/HOME_NET: \"[%s]\"/", networksStr),
		configPath)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update HOME_NET: %s", output)
	}

	logger.Infof("Suricata HOME_NET updated to: [%s]", networksStr)
	return nil
}

// UpdateSuricataRules updates Suricata rules
func (r *Repository) UpdateSuricataRules() error {
	cmd := exec.Command("suricata-update")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update Suricata rules: %s", output)
	}

	logger.Info("Suricata rules updated successfully")
	return nil
}

// GetSuricataAlerts gets recent Suricata alerts
func (r *Repository) GetSuricataAlerts(limit int) ([]Alert, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", limit), "/var/log/suricata/eve.json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read alerts: %w", err)
	}

	var alerts []Alert
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if eventType, ok := entry["event_type"].(string); ok && eventType == "alert" {
			alert := Alert{
				Timestamp: getString(entry, "timestamp"),
				EventType: eventType,
				SrcIP:     getString(entry, "src_ip"),
				DestIP:    getString(entry, "dest_ip"),
				Proto:     getString(entry, "proto"),
			}

			if srcPort, ok := entry["src_port"].(float64); ok {
				alert.SrcPort = int(srcPort)
			}

			if destPort, ok := entry["dest_port"].(float64); ok {
				alert.DestPort = int(destPort)
			}

			if alertData, ok := entry["alert"].(map[string]interface{}); ok {
				alert.Alert = alertData
			}

			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// CheckRPortConnection checks RPort connection status
func (r *Repository) CheckRPortConnection() (*ConnectionStatus, error) {
	cmd := exec.Command("journalctl", "-u", "rport", "-n", "50", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get RPort logs: %w", err)
	}

	status := &ConnectionStatus{
		Service:   "rport",
		Connected: false,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Connected to") {
		status.Connected = true
		status.Message = "RPort is connected to server"

		// Extract connection message
		lines := strings.Split(outputStr, "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.Contains(lines[i], "Connected to") {
				status.Message = lines[i]
				break
			}
		}
	} else {
		status.Message = "RPort connection status unknown"
	}

	return status, nil
}

// CheckRsyslogConfig checks if rsyslog firewall config is enabled
func (r *Repository) CheckRsyslogConfig() (*RsyslogConfigStatus, error) {
	configPath := "/etc/rsyslog.d/01-firewall.conf"
	status := &RsyslogConfigStatus{
		ListeningPort: 514,
	}

	// Check if config is enabled (not commented out)
	cmd := exec.Command("grep", "-v", "^#", configPath)
	output, _ := cmd.Output()
	outputStr := string(output)
	status.Enabled = strings.Contains(outputStr, "input(type=\"imudp\" port=\"514\"")

	// Check if port is listening
	cmd = exec.Command("ss", "-tulpn")
	output, _ = cmd.Output()
	status.ConfigValid = strings.Contains(string(output), ":514")

	if status.Enabled {
		status.Message = "Firewall log collection is enabled"
	} else {
		status.Message = "Firewall log collection is disabled"
	}

	return status, nil
}

// ToggleRsyslogConfig enables or disables rsyslog firewall configuration
func (r *Repository) ToggleRsyslogConfig(enabled bool) error {
	configPath := "/etc/rsyslog.d/01-firewall.conf"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read rsyslog config: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var newLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if enabled {
			// Remove leading comment
			if strings.HasPrefix(trimmed, "#") {
				line = strings.TrimPrefix(line, "#")
			}
		} else {
			// Add comment if not already commented
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				line = "#" + line
			}
		}

		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write rsyslog config: %w", err)
	}

	// Restart rsyslog
	if err := r.RestartService("rsyslog"); err != nil {
		return fmt.Errorf("failed to restart rsyslog: %w", err)
	}

	logger.Infof("Rsyslog firewall config %s", map[bool]string{true: "enabled", false: "disabled"}[enabled])
	return nil
}

// TailRsyslogFirewallLogs tails firewall logs
func (r *Repository) TailRsyslogFirewallLogs(lines int) ([]LogEntry, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), "/var/log/firewall.log")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read firewall logs: %w", err)
	}

	var logs []LogEntry
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		logs = append(logs, LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Message:   line,
		})
	}

	return logs, nil
}

// TailWazuhLogs reads Wazuh agent logs
func (r *Repository) TailWazuhLogs(lines int) ([]LogEntry, error) {
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), "/var/ossec/logs/ossec.log")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read Wazuh logs: %w", err)
	}

	var logs []LogEntry
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		logs = append(logs, LogEntry{
			Message: line,
		})
	}

	return logs, nil
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// TailServiceLogs tails logs for a service using journalctl
func (r *Repository) TailServiceLogs(serviceName string, lines int) ([]LogEntry, error) {
	cmd := exec.Command("journalctl", "-u", serviceName, "-n", fmt.Sprintf("%d", lines), "--no-pager", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get service logs: %w", err)
	}

	var logs []LogEntry
	decoder := json.NewDecoder(bytes.NewReader(output))

	for {
		var entry map[string]interface{}
		if err := decoder.Decode(&entry); err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		log := LogEntry{
			Message: getString(entry, "MESSAGE"),
		}

		if timestamp, ok := entry["__REALTIME_TIMESTAMP"].(string); ok {
			log.Timestamp = timestamp
		}

		if priority, ok := entry["PRIORITY"].(string); ok {
			log.Level = priority
		}

		logs = append(logs, log)
	}

	return logs, nil
}
