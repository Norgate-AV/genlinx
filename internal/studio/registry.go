package studio

import (
	"fmt"
	"sort"
)

const nlxKeyBase = `Software\AMX Corp.\NetLinx Studio`

// RegistrySettings holds the NetLinx Studio settings read from the Windows registry.
type RegistrySettings struct {
	EditorPreferences        map[string]any
	ClipboardHistory         map[string]any
	DiagnosticPreferences    map[string]any
	GeneralOptions           map[string]any
	AXSViewColorPreferences  map[string]any
	TextViewColorPreferences map[string]any
	CompilerOptions          map[string]any
	BatchTransferOptions     map[string]any
	KITFiles                 map[string]any
	ConnectionHistory        map[string]any
	IncludeDirs              []string
	ModuleDirs               []string
	LibraryDirs              []string
}

// Print prints all registry settings to stdout in a human-readable format.
func (s *RegistrySettings) Print() {
	printSection("Editor Preferences", s.EditorPreferences)
	printSection("Clipboard History", s.ClipboardHistory)
	printSection("Diagnostic Preferences", s.DiagnosticPreferences)
	printSection("General Options", s.GeneralOptions)
	printSection("AXS View Color Preferences", s.AXSViewColorPreferences)
	printSection("Text View Color Preferences", s.TextViewColorPreferences)
	printSection("Compiler Options", s.CompilerOptions)
	printSection("Batch Transfer Options", s.BatchTransferOptions)
	printSection("KIT Files", s.KITFiles)
	printSection("Connection History", s.ConnectionHistory)
	printDirList("Include Directories", s.IncludeDirs)
	printDirList("Module Directories", s.ModuleDirs)
	printDirList("Library Directories", s.LibraryDirs)
}

func printSection(title string, values map[string]any) {
	fmt.Printf("\n=== %s ===\n", title)

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("  %-45s %v\n", k, values[k])
	}
}

func printDirList(title string, dirs []string) {
	fmt.Printf("\n=== %s ===\n", title)
	for i, d := range dirs {
		fmt.Printf("  [%d] %s\n", i, d)
	}
}
