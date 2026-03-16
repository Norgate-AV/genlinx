//go:build !windows

package studio

import "fmt"

// DiffRegistry is not supported on non-Windows platforms.
func DiffRegistry(_ *Preferences) ([]DiffEntry, error) {
	return nil, fmt.Errorf("registry access is only supported on Windows")
}

// ApplyRegistryDiff is not supported on non-Windows platforms.
func ApplyRegistryDiff(_ []DiffEntry) error {
	return fmt.Errorf("registry access is only supported on Windows")
}
