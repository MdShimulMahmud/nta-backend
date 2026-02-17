package interfaces

import "time"

// Interface represents a network interface
type Interface struct {
	Name        string               `json:"name"`
	State       string               `json:"state"`
	Promiscuous bool                 `json:"promiscuous"`
	MAC         string               `json:"mac"`
	IPv4        string               `json:"ipv4,omitempty"`
	IPv6        string               `json:"ipv6,omitempty"`
	Speed       int                  `json:"speed,omitempty"`
	Duplex      string               `json:"duplex,omitempty"`
	MTU         int                  `json:"mtu"`
	Driver      string               `json:"driver,omitempty"`
	Type        string               `json:"type"`
	Carrier     bool                 `json:"carrier"`
	Statistics  *InterfaceStatistics `json:"statistics,omitempty"`
}

// InterfaceStatistics represents network interface statistics
type InterfaceStatistics struct {
	RxPackets uint64 `json:"rx_packets"`
	TxPackets uint64 `json:"tx_packets"`
	RxBytes   uint64 `json:"rx_bytes"`
	TxBytes   uint64 `json:"tx_bytes"`
	RxErrors  uint64 `json:"rx_errors"`
	TxErrors  uint64 `json:"tx_errors"`
}

// InterfaceStatus represents the status of an interface operation
type InterfaceStatus struct {
	Interface string `json:"interface"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// NetplanConfig represents the netplan configuration
type NetplanConfig struct {
	Network NetplanNetwork `json:"network" yaml:"network"`
}

// NetplanNetwork represents the network section in netplan
type NetplanNetwork struct {
	Version   int                           `json:"version" yaml:"version"`
	Ethernets map[string]NetplanEthernetCfg `json:"ethernets" yaml:"ethernets"`
}

// NetplanEthernetCfg represents an ethernet interface configuration
type NetplanEthernetCfg struct {
	Addresses   []string              `json:"addresses,omitempty" yaml:"addresses,omitempty"`
	DHCP4       bool                  `json:"dhcp4,omitempty" yaml:"dhcp4,omitempty"`
	DHCP6       bool                  `json:"dhcp6,omitempty" yaml:"dhcp6,omitempty"`
	Gateway4    string                `json:"gateway4,omitempty" yaml:"gateway4,omitempty"`
	Nameservers *NetplanNameservers   `json:"nameservers,omitempty" yaml:"nameservers,omitempty"`
	Routes      []NetplanRoute        `json:"routes,omitempty" yaml:"routes,omitempty"`
	Optional    bool                  `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// NetplanNameservers represents nameserver configuration
type NetplanNameservers struct {
	Addresses []string `json:"addresses,omitempty" yaml:"addresses,omitempty"`
	Search    []string `json:"search,omitempty" yaml:"search,omitempty"`
}

// NetplanRoute represents a route configuration
type NetplanRoute struct {
	To  string `json:"to" yaml:"to"`
	Via string `json:"via" yaml:"via"`
}

// TCPDumpOutput represents tcpdump output
type TCPDumpOutput struct {
	Interface   string    `json:"interface"`
	IsRunning   bool      `json:"is_running"`
	PacketCount int       `json:"packet_count"`
	Output      []string  `json:"output"`
	StartedAt   time.Time `json:"started_at,omitempty"`
}

// UpdateInterfaceStateRequest represents request to update interface state
type UpdateInterfaceStateRequest struct {
	State string `json:"state" binding:"required,oneof=up down"`
}

// UpdatePromiscModeRequest represents request to update promiscuous mode
type UpdatePromiscModeRequest struct {
	Enabled bool `json:"enabled" binding:"required"`
}

// UpdateNetplanRequest represents request to update netplan config
type UpdateNetplanRequest struct {
	Interface   string              `json:"interface" binding:"required"`
	Addresses   []string            `json:"addresses,omitempty"`
	DHCP4       *bool               `json:"dhcp4,omitempty"`
	DHCP6       *bool               `json:"dhcp6,omitempty"`
	Gateway4    string              `json:"gateway4,omitempty"`
	Nameservers *NetplanNameservers `json:"nameservers,omitempty"`
	Routes      []NetplanRoute      `json:"routes,omitempty"`
	Optional    *bool               `json:"optional,omitempty"`
}
