package interfaces

import (
	"fmt"

	"github.com/sharetrip/nta-backend/internal/config"
	"github.com/sharetrip/nta-backend/internal/logger"
)

// Service handles business logic for network interfaces
type Service struct {
	repo   *Repository
	config *config.Config
}

// NewService creates a new interface service
func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
	}
}

// GetAllInterfaces returns all network interfaces
func (s *Service) GetAllInterfaces() ([]Interface, error) {
	logger.Debug("Fetching all network interfaces")
	interfaces, err := s.repo.ListInterfaces()
	if err != nil {
		logger.Errorf("Failed to get interfaces: %v", err)
		return nil, err
	}

	logger.Debugf("Found %d interfaces", len(interfaces))
	return interfaces, nil
}

// GetInterface returns a specific interface
func (s *Service) GetInterface(name string) (*Interface, error) {
	logger.Debugf("Fetching interface: %s", name)
	iface, err := s.repo.GetInterface(name)
	if err != nil {
		logger.Errorf("Failed to get interface %s: %v", name, err)
		return nil, err
	}

	return iface, nil
}

// SetInterfaceUp brings an interface up
func (s *Service) SetInterfaceUp(name string) error {
	logger.Infof("Bringing interface %s up", name)

	if err := s.repo.SetInterfaceState(name, "up"); err != nil {
		logger.Errorf("Failed to bring interface %s up: %v", name, err)
		return err
	}

	return nil
}

// SetInterfaceDown brings an interface down
func (s *Service) SetInterfaceDown(name string) error {
	logger.Infof("Bringing interface %s down", name)

	if err := s.repo.SetInterfaceState(name, "down"); err != nil {
		logger.Errorf("Failed to bring interface %s down: %v", name, err)
		return err
	}

	return nil
}

// UpdateInterfaceState updates interface state
func (s *Service) UpdateInterfaceState(name, state string) error {
	logger.Infof("Updating interface %s state to %s", name, state)

	if state != "up" && state != "down" {
		return fmt.Errorf("invalid state: %s (must be 'up' or 'down')", state)
	}

	if err := s.repo.SetInterfaceState(name, state); err != nil {
		logger.Errorf("Failed to update interface %s state: %v", name, err)
		return err
	}

	return nil
}

// EnablePromiscuousMode enables promiscuous mode on an interface
func (s *Service) EnablePromiscuousMode(name string) error {
	logger.Infof("Enabling promiscuous mode on %s", name)

	if err := s.repo.SetPromiscuousMode(name, true); err != nil {
		logger.Errorf("Failed to enable promiscuous mode on %s: %v", name, err)
		return err
	}

	return nil
}

// DisablePromiscuousMode disables promiscuous mode on an interface
func (s *Service) DisablePromiscuousMode(name string) error {
	logger.Infof("Disabling promiscuous mode on %s", name)

	if err := s.repo.SetPromiscuousMode(name, false); err != nil {
		logger.Errorf("Failed to disable promiscuous mode on %s: %v", name, err)
		return err
	}

	return nil
}

// UpdatePromiscuousMode updates promiscuous mode
func (s *Service) UpdatePromiscuousMode(name string, enabled bool) error {
	logger.Infof("Updating promiscuous mode on %s to %v", name, enabled)

	if err := s.repo.SetPromiscuousMode(name, enabled); err != nil {
		logger.Errorf("Failed to update promiscuous mode on %s: %v", name, err)
		return err
	}

	return nil
}

// GetNetplanConfig retrieves the netplan configuration
func (s *Service) GetNetplanConfig() (*NetplanConfig, error) {
	logger.Debug("Fetching netplan configuration")

	config, err := s.repo.GetNetplanConfig(s.config.Network.NetplanPath)
	if err != nil {
		logger.Errorf("Failed to get netplan config: %v", err)
		return nil, err
	}

	return config, nil
}

// UpdateNetplanInterface updates a specific interface in netplan
func (s *Service) UpdateNetplanInterface(req *UpdateNetplanRequest) error {
	logger.Infof("Updating netplan configuration for interface %s", req.Interface)

	// Get current config
	config, err := s.repo.GetNetplanConfig(s.config.Network.NetplanPath)
	if err != nil {
		return err
	}

	// Ensure ethernets map exists
	if config.Network.Ethernets == nil {
		config.Network.Ethernets = make(map[string]NetplanEthernetCfg)
	}

	// Update or create interface config
	ifaceCfg := NetplanEthernetCfg{}

	if existing, ok := config.Network.Ethernets[req.Interface]; ok {
		ifaceCfg = existing
	}

	// Update fields
	if req.Addresses != nil {
		ifaceCfg.Addresses = req.Addresses
	}

	if req.DHCP4 != nil {
		ifaceCfg.DHCP4 = *req.DHCP4
	}

	if req.DHCP6 != nil {
		ifaceCfg.DHCP6 = *req.DHCP6
	}

	if req.Gateway4 != "" {
		ifaceCfg.Gateway4 = req.Gateway4
	}

	if req.Nameservers != nil {
		ifaceCfg.Nameservers = req.Nameservers
	}

	if req.Routes != nil {
		ifaceCfg.Routes = req.Routes
	}

	if req.Optional != nil {
		ifaceCfg.Optional = *req.Optional
	}

	config.Network.Ethernets[req.Interface] = ifaceCfg

	// Save config
	if err := s.repo.UpdateNetplanConfig(s.config.Network.NetplanPath, config); err != nil {
		logger.Errorf("Failed to update netplan config: %v", err)
		return err
	}

	// Apply config
	if err := s.repo.ApplyNetplanConfig(); err != nil {
		logger.Errorf("Failed to apply netplan config: %v", err)
		return err
	}

	logger.Infof("Netplan configuration updated successfully for %s", req.Interface)
	return nil
}

// CheckTCPDumpNetwork checks network logs with tcpdump
func (s *Service) CheckTCPDumpNetwork(interfaceName string) (*TCPDumpOutput, error) {
	logger.Infof("Checking network logs on %s with tcpdump", interfaceName)

	result, err := s.repo.CheckTCPDump(interfaceName, 10)
	if err != nil {
		logger.Errorf("Failed to check tcpdump on %s: %v", interfaceName, err)
		return nil, err
	}

	return result, nil
}

// CheckTCPDumpFirewall checks firewall logs with tcpdump
func (s *Service) CheckTCPDumpFirewall(interfaceName string) (*TCPDumpOutput, error) {
	logger.Infof("Checking firewall logs on %s with tcpdump", interfaceName)

	// Check for syslog port 514
	result, err := s.repo.CheckTCPDump(interfaceName, 10)
	if err != nil {
		logger.Errorf("Failed to check firewall tcpdump on %s: %v", interfaceName, err)
		return nil, err
	}

	return result, nil
}
