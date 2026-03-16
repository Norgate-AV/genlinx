//go:build windows

package studio

import (
	"fmt"
	"sort"

	"golang.org/x/sys/windows/registry"
)

const nlxKeyBase = `Software\AMX Corp.\NetLinx Studio`

// IsInstalled reports whether NetLinx Studio is installed by checking for
// the existence of its HKCU registry key.
func IsInstalled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, nlxKeyBase, registry.READ)
	if err != nil {
		return false
	}
	_ = k.Close()
	return true
}

// ReadRegistry reads all known NetLinx Studio settings from HKCU / HKLM.
func ReadRegistry() (*RegistrySettings, error) {
	s := &RegistrySettings{
		EditorPreferences:        make(map[string]any),
		ClipboardHistory:         make(map[string]any),
		DiagnosticPreferences:    make(map[string]any),
		GeneralOptions:           make(map[string]any),
		AXSViewColorPreferences:  make(map[string]any),
		TextViewColorPreferences: make(map[string]any),
		CompilerOptions:          make(map[string]any),
		BatchTransferOptions:     make(map[string]any),
		KITFiles:                 make(map[string]any),
		ConnectionHistory:        make(map[string]any),
	}

	hkcu, err := registry.OpenKey(registry.CURRENT_USER, nlxKeyBase, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("open HKCU\\%s: %w", nlxKeyBase, err)
	}

	defer func() { _ = hkcu.Close() }()

	for _, spec := range []struct {
		subKey string
		dest   map[string]any
	}{
		{"Editor Preferences", s.EditorPreferences},
		{"Text Clipboard History", s.ClipboardHistory},
		{"Diagnostic Preferences", s.DiagnosticPreferences},
		{"General Options", s.GeneralOptions},
		{"AXSViewColorPreferences", s.AXSViewColorPreferences},
		{"TextViewColorPreferences", s.TextViewColorPreferences},
		{"NLXCompiler_Options", s.CompilerOptions},
		{"Batch Transfer User Options", s.BatchTransferOptions},
		{"KIT Files", s.KITFiles},
		{"RecentConnectionsHistory", s.ConnectionHistory},
	} {
		if err := readSubKeyValues(hkcu, spec.subKey, spec.dest); err != nil {
			return nil, err
		}
	}

	s.IncludeDirs, err = readDirList(registry.LOCAL_MACHINE, "NLXCompiler_Includes")
	if err != nil {
		return nil, err
	}

	s.ModuleDirs, err = readDirList(registry.LOCAL_MACHINE, "NLXCompiler_Modules")
	if err != nil {
		return nil, err
	}

	s.LibraryDirs, err = readDirList(registry.LOCAL_MACHINE, "NLXCompiler_Libs")
	if err != nil {
		return nil, err
	}

	return s, nil
}

// readSubKeyValues opens a child key relative to parent and reads all DWORD
// and string values into dest. Missing keys are silently skipped.
func readSubKeyValues(parent registry.Key, subKey string, dest map[string]any) error {
	k, err := registry.OpenKey(parent, subKey, registry.READ)
	if err != nil {
		return nil
	}

	defer func() { _ = k.Close() }()

	names, err := k.ReadValueNames(0)
	if err != nil {
		return fmt.Errorf("read value names from %q: %w", subKey, err)
	}

	for _, name := range names {
		_, valType, err := k.GetValue(name, nil)
		if err != nil {
			continue
		}

		switch valType {
		case registry.DWORD, registry.QWORD:
			v, _, err := k.GetIntegerValue(name)
			if err == nil {
				dest[name] = uint32(v)
			}
		case registry.SZ, registry.EXPAND_SZ:
			v, _, err := k.GetStringValue(name)
			if err == nil {
				dest[name] = v
			}
		}
	}

	return nil
}

// readDirList reads all string values from an HKLM WOW6432Node sub-key and
// returns them sorted by value name (preserving Dir000, Dir001, ... order).
func readDirList(root registry.Key, subKeyName string) ([]string, error) {
	path := `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\` + subKeyName
	k, err := registry.OpenKey(root, path, registry.READ)
	if err != nil {
		return nil, nil
	}

	defer func() { _ = k.Close() }()

	names, err := k.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("read value names from %q: %w", subKeyName, err)
	}

	sort.Strings(names)

	var dirs []string
	for _, name := range names {
		v, _, err := k.GetStringValue(name)
		if err == nil && v != "" {
			dirs = append(dirs, v)
		}
	}

	return dirs, nil
}
