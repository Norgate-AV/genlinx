package studio

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// ParsePreferences reads and decodes a NetLinx Studio .epx preferences file
// from r. The returned Preferences can be passed directly to WriteRegistry.
func ParsePreferences(r io.Reader) (*Preferences, error) {
	var prefs Preferences
	if err := xml.NewDecoder(r).Decode(&prefs); err != nil {
		return nil, fmt.Errorf("decode preferences: %w", err)
	}

	return &prefs, nil
}

// ParsePreferencesFile is a convenience wrapper that opens path and calls
// ParsePreferences.
func ParsePreferencesFile(path string) (*Preferences, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	defer func() { _ = f.Close() }()

	return ParsePreferences(f)
}
