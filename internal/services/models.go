package services

// ServiceStatus represents the status of a system service
type ServiceStatus struct {
	Name      string                 `json:"name"`
	Active    bool                   `json:"active"`
	Status    string                 `json:"status"`
	Running   bool                   `json:"running"`
	Enabled   bool                   `json:"enabled"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

// ServiceAction represents an action to perform on a service
type ServiceAction string

const (
	ActionStart   ServiceAction = "start"
	ActionStop    ServiceAction = "stop"
	ActionRestart ServiceAction = "restart"
	ActionReload  ServiceAction = "reload"
	ActionStatus  ServiceAction = "status"
	ActionEnable  ServiceAction = "enable"
	ActionDisable ServiceAction = "disable"
)

// ServiceActionRequest represents a request to perform an action on a service
type ServiceActionRequest struct {
	Action ServiceAction `json:"action" binding:"required,oneof=start stop restart reload status"`
}

// Pipeline represents a logstash pipeline
type Pipeline struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Status  string `json:"status,omitempty"`
}

// PipelineToggleRequest represents request to enable/disable pipeline
type PipelineToggleRequest struct {
	PipelineID string `json:"pipeline_id" binding:"required"`
	Enabled    bool   `json:"enabled"`
}

// SuricataConfig represents Suricata-specific configuration
type SuricataConfig struct {
	Interface string   `json:"interface"`
	HomeNet   []string `json:"home_net"`
}

// UpdateSuricataInterfaceRequest represents request to update Suricata interface
type UpdateSuricataInterfaceRequest struct {
	Interface string `json:"interface" binding:"required"`
}

// UpdateS	uricataHomeNetRequest represents request to update Suricata home network
type UpdateSuricataHomeNetRequest struct {
	Networks []string `json:"networks" binding:"required"`
}

// Alert represents a Suricata alert
type Alert struct {
	Timestamp   string                 `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	SrcIP       string                 `json:"src_ip"`
	SrcPort     int                    `json:"src_port"`
	DestIP      string                 `json:"dest_ip"`
	DestPort    int                    `json:"dest_port"`
	Proto       string                 `json:"proto"`
	Alert       map[string]interface{} `json:"alert,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// LogEntry represents a generic log entry
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Level     string `json:"level,omitempty"`
	Details   string `json:"details,omitempty"`
}

// ConnectionStatus represents service connection status
type ConnectionStatus struct {
	Service   string `json:"service"`
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp,omitempty"`
}

// RsyslogConfigStatus represents rsyslog configuration status
type RsyslogConfigStatus struct {
	Enabled       bool   `json:"enabled"`
	ListeningPort int    `json:"listening_port"`
	ConfigValid   bool   `json:"config_valid"`
	Message       string `json:"message,omitempty"`
}
