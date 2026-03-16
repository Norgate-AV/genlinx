//go:build windows

package studio

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// DiffRegistry reads the current registry state and returns the diff between
// it and the entries derived from prefs. The result is the set of entries that
// would actually be written — ADDED entries that don't yet exist and CHANGED
// entries whose value differs from what is currently stored.
func DiffRegistry(prefs *Preferences) ([]DiffEntry, error) {
	currentSettings, err := ReadRegistry()
	if err != nil {
		return nil, fmt.Errorf("read current registry for diff: %w", err)
	}

	current := PreferencesToRegistryEntries(BuildPreferences(currentSettings))
	incoming := PreferencesToRegistryEntries(prefs)

	return DiffRegistryEntries(current, incoming), nil
}

// ApplyRegistryDiff writes a pre-computed delta to the registry. HKLM entries
// (directory lists) require administrator privileges.
func ApplyRegistryDiff(delta []DiffEntry) error {
	if len(delta) == 0 {
		return nil
	}

	hkcuBase, err := registry.OpenKey(
		registry.CURRENT_USER,
		nlxKeyBase,
		registry.WRITE,
	)
	if err != nil {
		return fmt.Errorf("open HKCU\\%s for write: %w", nlxKeyBase, err)
	}

	defer func() { _ = hkcuBase.Close() }()

	for _, d := range delta {
		entry := d.RegistryEntry

		if entry.HiveLM {
			if err := writeHKLMEntry(entry); err != nil {
				return err
			}

			continue
		}

		k, _, err := registry.CreateKey(hkcuBase, entry.SubKey, registry.SET_VALUE)
		if err != nil {
			return fmt.Errorf("create/open HKCU\\%s\\%s: %w", nlxKeyBase, entry.SubKey, err)
		}

		if err := setRegistryValue(k, entry.ValueName, entry.Value); err != nil {
			_ = k.Close()
			return fmt.Errorf("set %s\\%s: %w", entry.SubKey, entry.ValueName, err)
		}

		_ = k.Close()
	}

	return nil
}

func writeHKLMEntry(entry RegistryEntry) error {
	k, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE,
		entry.SubKey,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("create/open HKLM\\%s (requires admin): %w", entry.SubKey, err)
	}

	defer func() { _ = k.Close() }()

	return setRegistryValue(k, entry.ValueName, entry.Value)
}

func setRegistryValue(k registry.Key, name string, value any) error {
	switch v := value.(type) {
	case uint32:
		return k.SetDWordValue(name, v)
	case string:
		return k.SetStringValue(name, v)
	default:
		return fmt.Errorf("unsupported registry value type %T for %q", value, name)
	}
}
