package utils

import (
	"os"
	"os/exec"
	"strings"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsRoot checks if the current user is root
func IsRoot() bool {
	return os.Geteuid() == 0
}

// IsServiceInstalled checks if a systemd service is installed
func IsServiceInstalled(serviceName string) bool {
	cmd := exec.Command("systemctl", "list-unit-files", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), serviceName)
}

// EnsureDirectory ensures a directory exists
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// ReadFile reads a file and returns its content
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file
func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
