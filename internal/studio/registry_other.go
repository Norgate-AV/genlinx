//go:build !windows

package studio

import "fmt"

// ReadRegistry is not supported on non-Windows platforms.
func ReadRegistry() (*RegistrySettings, error) {
	return nil, fmt.Errorf("registry access is only supported on Windows")
}

// IsInstalled always returns false on non-Windows platforms.
func IsInstalled() bool {
	return false
}
