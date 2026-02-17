package interfaces

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sharetrip/nta-backend/internal/logger"
	"gopkg.in/yaml.v3"
)

// Repository handles network interface operations
type Repository struct{}

// NewRepository creates a new interface repository
func NewRepository() *Repository {
	return &Repository{}
}

// ListInterfaces returns all network interfaces
func (r *Repository) ListInterfaces() ([]Interface, error) {
	netDir := "/sys/class/net"
	entries, err := os.ReadDir(netDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read network interfaces: %w", err)
	}

	var interfaces []Interface
	       for _, entry := range entries {
		       // Skip loopback
		       if entry.Name() == "lo" {
			       continue
		       }

		       logger.Debugf("Processing interface: %s", entry.Name())
		       iface, err := r.GetInterface(entry.Name())
		       if err != nil {
			       logger.Warnf("Failed to get interface %s: %v", entry.Name(), err)
			       continue
		       }
		       interfaces = append(interfaces, *iface)
	       }

	return interfaces, nil
}

// GetInterface returns details for a specific interface
func (r *Repository) GetInterface(name string) (*Interface, error) {
	       logger.Debugf("Gathering details for interface: %s", name)
	basePath := filepath.Join("/sys/class/net", name)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		       logger.Warnf("Base path missing for interface %s: %v", name, err)
		return nil, fmt.Errorf("interface %s not found", name)
	}

	iface := &Interface{
		Name: name,
	}

	// Get state
	if state, err := r.readFile(filepath.Join(basePath, "operstate")); err == nil {
		iface.State = strings.TrimSpace(state)
	       } else {
		       logger.Warnf("Missing operstate for %s: %v", name, err)
	}

	// Get MAC address
	if mac, err := r.readFile(filepath.Join(basePath, "address")); err == nil {
		iface.MAC = strings.TrimSpace(mac)
	       } else {
		       logger.Warnf("Missing address for %s: %v", name, err)
	}

	// Get MTU
	if mtuStr, err := r.readFile(filepath.Join(basePath, "mtu")); err == nil {
		if mtu, err := strconv.Atoi(strings.TrimSpace(mtuStr)); err == nil {
			iface.MTU = mtu
		}
	       } else {
		       logger.Warnf("Missing mtu for %s: %v", name, err)
	}

	// Check promiscuous mode
	iface.Promiscuous = r.isPromiscuous(name)

	// Get type
	iface.Type = r.getInterfaceType(basePath)

	// Get carrier status
	if carrierStr, err := r.readFile(filepath.Join(basePath, "carrier")); err == nil {
		iface.Carrier = strings.TrimSpace(carrierStr) == "1"
	       } else {
		       logger.Warnf("Missing carrier for %s: %v", name, err)
	}

	// Get speed and duplex if interface is up
	if iface.State == "up" {
		if speedStr, err := r.readFile(filepath.Join(basePath, "speed")); err == nil {
			if speed, err := strconv.Atoi(strings.TrimSpace(speedStr)); err == nil {
				iface.Speed = speed
			}
		       } else {
			       logger.Warnf("Missing speed for %s: %v", name, err)
		}

		if duplex, err := r.readFile(filepath.Join(basePath, "duplex")); err == nil {
			iface.Duplex = strings.TrimSpace(duplex)
		       } else {
			       logger.Warnf("Missing duplex for %s: %v", name, err)
		}
	}

	// Get driver
	driverPath := filepath.Join(basePath, "device", "driver")
	if link, err := os.Readlink(driverPath); err == nil {
		iface.Driver = filepath.Base(link)
	       } else {
		       logger.Warnf("Missing driver for %s: %v", name, err)
	}

	// Get IP addresses
	r.getIPAddresses(name, iface)

	// Get statistics
	iface.Statistics = r.getStatistics(basePath)

	return iface, nil
}

// SetInterfaceState sets the interface state (up/down)
func (r *Repository) SetInterfaceState(name, state string) error {
	cmd := exec.Command("ip", "link", "set", name, state)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set interface %s to %s: %s", name, state, output)
	}

	logger.Infof("Interface %s set to %s", name, state)
	return nil
}

// SetPromiscuousMode enables or disables promiscuous mode
func (r *Repository) SetPromiscuousMode(name string, enabled bool) error {
	mode := "off"
	if enabled {
		mode = "on"
	}

	cmd := exec.Command("ip", "link", "set", name, "promisc", mode)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set promiscuous mode on %s: %s", name, output)
	}

	logger.Infof("Promiscuous mode on %s set to %s", name, mode)
	return nil
}

// GetNetplanConfig reads the netplan configuration
func (r *Repository) GetNetplanConfig(path string) (*NetplanConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read netplan config: %w", err)
	}

	var config NetplanConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse netplan config: %w", err)
	}

	return &config, nil
}

// UpdateNetplanConfig updates the netplan configuration
func (r *Repository) UpdateNetplanConfig(path string, config *NetplanConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal netplan config: %w", err)
	}

	// Create backup
	backupPath := path + ".backup"
	if err := r.copyFile(path, backupPath); err != nil {
		logger.Warnf("Failed to create backup: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write netplan config: %w", err)
	}

	logger.Infof("Netplan config updated at %s", path)
	return nil
}

// ApplyNetplanConfig applies the netplan configuration
func (r *Repository) ApplyNetplanConfig() error {
	cmd := exec.Command("netplan", "apply")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply netplan config: %s", output)
	}

	logger.Info("Netplan config applied successfully")
	return nil
}

// CheckTCPDump checks if tcpdump is capturing traffic
func (r *Repository) CheckTCPDump(iface string, count int) (*TCPDumpOutput, error) {
	cmd := exec.Command("timeout", "5", "tcpdump", "-i", iface, "-c", strconv.Itoa(count), "-n")
	output, err := cmd.CombinedOutput()

	result := &TCPDumpOutput{
		Interface: iface,
		IsRunning: err == nil,
		Output:    []string{},
	}

	if len(output) > 0 {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				result.Output = append(result.Output, line)
			}
		}

		result.PacketCount = len(result.Output)
	}

	return result, nil
}

// Helper methods
func (r *Repository) readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *Repository) isPromiscuous(name string) bool {
	cmd := exec.Command("ip", "link", "show", name)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "PROMISC")
}

func (r *Repository) getInterfaceType(basePath string) string {
	typeStr, err := r.readFile(filepath.Join(basePath, "type"))
	if err != nil {
		return "Unknown"
	}

	typeNum := strings.TrimSpace(typeStr)
	switch typeNum {
	case "1":
		return "Ethernet"
	case "2":
		return "Loopback"
	case "3":
		return "Bridge"
	default:
		return fmt.Sprintf("Unknown (%s)", typeNum)
	}
}

func (r *Repository) getIPAddresses(name string, iface *Interface) {
	// Get IPv4
	cmd := exec.Command("ip", "-4", "addr", "show", name)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "inet ") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "inet" && i+1 < len(fields) {
						iface.IPv4 = fields[i+1]
						break
					}
				}
			}
		}
	}

	// Get IPv6
	cmd = exec.Command("ip", "-6", "addr", "show", name)
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "inet6 ") && !strings.Contains(line, "fe80") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "inet6" && i+1 < len(fields) {
						iface.IPv6 = fields[i+1]
						break
					}
				}
			}
		}
	}
}

func (r *Repository) getStatistics(basePath string) *InterfaceStatistics {
	statsPath := filepath.Join(basePath, "statistics")
	stats := &InterfaceStatistics{}

	if val, err := r.readFile(filepath.Join(statsPath, "rx_packets")); err == nil {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64); err == nil {
			stats.RxPackets = parsed
		}
	}

	if val, err := r.readFile(filepath.Join(statsPath, "tx_packets")); err == nil {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64); err == nil {
			stats.TxPackets = parsed
		}
	}

	if val, err := r.readFile(filepath.Join(statsPath, "rx_bytes")); err == nil {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64); err == nil {
			stats.RxBytes = parsed
		}
	}

	if val, err := r.readFile(filepath.Join(statsPath, "tx_bytes")); err == nil {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64); err == nil {
			stats.TxBytes = parsed
		}
	}

	if val, err := r.readFile(filepath.Join(statsPath, "rx_errors")); err == nil {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64); err == nil {
			stats.RxErrors = parsed
		}
	}

	if val, err := r.readFile(filepath.Join(statsPath, "tx_errors")); err == nil {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64); err == nil {
			stats.TxErrors = parsed
		}
	}

	return stats
}

func (r *Repository) copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
