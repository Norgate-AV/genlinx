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
	diff := DiffRegistryEntries(current, incoming)

	// TCP/IP history and HKLM dir lists are diffed as ordered sets and replaced
	// wholesale, matching the behaviour of the native NetLinx Studio settings import.
	const dirBase = `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\`

	for _, spec := range []struct {
		subKey   string
		current  []string
		incoming []string
	}{
		{dirBase + "NLXCompiler_Libs", currentSettings.LibraryDirs, prefs.NetlinxCompilerSettings.LibraryDirs},
		{dirBase + "NLXCompiler_Includes", currentSettings.IncludeDirs, prefs.NetlinxCompilerSettings.IncludeDirs},
		{dirBase + "NLXCompiler_Modules", currentSettings.ModuleDirs, prefs.NetlinxCompilerSettings.ModuleDirs},
	} {
		if d := DiffDirList(spec.subKey, spec.current, spec.incoming); d != nil {
			diff = append(diff, *d)
		}
	}

	currentHistory := buildTCPIPHistory(currentSettings.ConnectionHistory)
	if d := DiffTCPIPHistory(currentHistory.Entries, prefs.TCPIPHistory.Entries); d != nil {
		diff = append(diff, *d)
	}

	return diff, nil
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

		if d.Status == DiffReplaced {
			if d.SubKey == "RecentConnectionsHistory" {
				if err := replaceTCPIPHistory(hkcuBase, d.TCPIPReplace); err != nil {
					return err
				}
			} else {
				if err := replaceDirList(d.SubKey, d.DirReplace); err != nil {
					return err
				}
			}

			continue
		}

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

// replaceDirList deletes all existing Dir* values under the given HKLM subkey
// and writes dirs verbatim from index 0, matching the behaviour of the native
// NetLinx Studio settings import.
func replaceDirList(subKey string, dirs []string) error {
	k, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE,
		subKey,
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		return fmt.Errorf("open HKLM\\%s for write (requires admin): %w", subKey, err)
	}

	defer func() { _ = k.Close() }()

	names, err := k.ReadValueNames(0)
	if err != nil {
		return fmt.Errorf("read value names from %s: %w", subKey, err)
	}

	for _, name := range names {
		if err := k.DeleteValue(name); err != nil {
			return fmt.Errorf("delete %s\\%s: %w", subKey, name, err)
		}
	}

	for i, dir := range dirs {
		name := fmt.Sprintf("Dir%03d", i)
		if err := k.SetStringValue(name, dir); err != nil {
			return fmt.Errorf("write %s\\%s: %w", subKey, name, err)
		}
	}

	return nil
}

// replaceTCPIPHistory deletes all existing Recent Connection History* values
// under RecentConnectionsHistory and writes incoming verbatim from index 0,
// matching the behaviour of the native NetLinx Studio settings import.
func replaceTCPIPHistory(hkcuBase registry.Key, entries []TCPIPEntry) error {
	const subKey = "RecentConnectionsHistory"

	k, _, err := registry.CreateKey(hkcuBase, subKey, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return fmt.Errorf("open HKCU\\...\\%s for write: %w", subKey, err)
	}

	defer func() { _ = k.Close() }()

	// Delete all existing values.
	names, err := k.ReadValueNames(0)
	if err != nil {
		return fmt.Errorf("read value names from %s: %w", subKey, err)
	}

	for _, name := range names {
		if err := k.DeleteValue(name); err != nil {
			return fmt.Errorf("delete %s\\%s: %w", subKey, name, err)
		}
	}

	// Write incoming entries sequentially from index 0.
	for i, entry := range entries {
		name := fmt.Sprintf("Recent Connection History%d", i)
		if err := k.SetStringValue(name, "T-"+entry.encode()); err != nil {
			return fmt.Errorf("write %s\\%s: %w", subKey, name, err)
		}
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
