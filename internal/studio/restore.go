package studio

import "fmt"

// DiffStatus describes whether a registry entry would be added or changed.
type DiffStatus string

const (
	DiffAdded   DiffStatus = "ADDED"
	DiffChanged DiffStatus = "CHANGED"
)

// DiffEntry pairs a registry entry with its diff status and (for CHANGED
// entries) the value currently in the registry.
type DiffEntry struct {
	RegistryEntry
	Status   DiffStatus
	OldValue any // populated for DiffChanged; nil for DiffAdded
}

type entryKey struct {
	hiveLM    bool
	subKey    string
	valueName string
}

// DiffRegistryEntries returns the subset of incoming entries that differ from
// the corresponding entries in current. An entry is ADDED when its key does
// not exist in current; it is CHANGED when the key exists but the value
// differs. Identical entries are omitted.
//
// This function is platform-neutral and fully testable without registry access.
func DiffRegistryEntries(current, incoming []RegistryEntry) []DiffEntry {
	index := make(map[entryKey]any, len(current))
	for _, e := range current {
		index[entryKey{e.HiveLM, e.SubKey, e.ValueName}] = e.Value
	}

	var diff []DiffEntry

	for _, e := range incoming {
		k := entryKey{e.HiveLM, e.SubKey, e.ValueName}
		oldVal, exists := index[k]

		if !exists {
			diff = append(diff, DiffEntry{RegistryEntry: e, Status: DiffAdded})
		} else if oldVal != e.Value {
			diff = append(diff, DiffEntry{RegistryEntry: e, Status: DiffChanged, OldValue: oldVal})
		}
	}

	return diff
}

// RegistryEntry describes a single registry value to be written: the HKCU
// sub-key path relative to the NetLinx Studio base key, the value name, and
// the value itself (always uint32 or string).
type RegistryEntry struct {
	// SubKey is the path relative to HKCU\Software\AMX Corp.\NetLinx Studio,
	// or the full HKLM WOW6432Node path when HiveLM is true.
	SubKey    string
	ValueName string
	Value     any // uint32 or string
	// HiveLM indicates the entry belongs to HKLM rather than HKCU.
	HiveLM bool
}

// PreferencesToRegistryEntries converts a Preferences struct into a flat list
// of registry entries ready to be written by WriteRegistry. All logic for
// mapping EPX fields back to registry key paths and value names lives here,
// keeping WriteRegistry itself a trivial write loop.
//
// This function is platform-neutral and fully testable without registry access.
func PreferencesToRegistryEntries(prefs *Preferences) []RegistryEntry {
	var entries []RegistryEntry

	add := func(subKey, name string, value any) {
		entries = append(entries, RegistryEntry{SubKey: subKey, ValueName: name, Value: value})
	}

	addHKLM := func(subKey, name string, value any) {
		entries = append(entries, RegistryEntry{SubKey: subKey, ValueName: name, Value: value, HiveLM: true})
	}

	// -----------------------------------------------------------------------
	// Editor Preferences  (also covers Clipboard History sub-key)
	// -----------------------------------------------------------------------
	e := prefs.EditorSettings
	add("Editor Preferences", "AutoIdent", uint32(e.AutoIndentEnabled))
	add("Editor Preferences", "ShowLineNumbers", uint32(e.ShowLineNumbersEnabled))
	add("Editor Preferences", "CodeFolding", uint32(e.CodeFoldingEnabled))
	add("Editor Preferences", "AutoSuggest", uint32(e.AutoSuggestEnabled))
	add("Editor Preferences", "IndentGuides", uint32(e.IndentGuidesEnabled))
	add("Editor Preferences", "CallTips", uint32(e.CallTipsEnabled))
	add("Editor Preferences", "UTF8FormatEnabled", uint32(e.UTF8FormatEnabled))
	add("Editor Preferences", "IndentWidth", uint32(e.IndentationWidth))
	add("Editor Preferences", "TabWidth", uint32(e.TabSpaces))
	add("Editor Preferences", "PointSize", uint32(e.FontPointSize))
	add("Editor Preferences", "FontName", e.FontName)
	add("Editor Preferences", "UseOnlyFixedPitchFonts", uint32(e.OnlyFixedPitchFontsEnabled))
	add("Editor Preferences", "PrinterColorMode", uint32(e.PrinterColorMode))
	add("Editor Preferences", "RightTrimLines", uint32(e.TrimLinesRight))
	add("Editor Preferences", "AutoSuggestInComments", uint32(e.AutoSuggestInComments))
	add("Editor Preferences", "ColumnEdgeMarker", uint32(e.ColumnEdgeMarker))
	add("Editor Preferences", "ColumnEdgeColumnNumber", uint32(e.ColumnEdgeMarkerValue))
	add("Editor Preferences", "SaveBookmarksOnClose", uint32(e.RetainBookmarksOnFileClose))
	// Preserve the registry typo on write so Studio can read it back.
	add("Editor Preferences", "HightlightMatchingBraces", uint32(e.HighlightMatchingBraces))
	add("Editor Preferences", "UseTabs", uint32(e.UseTabEnabled))

	// Clipboard History sub-key
	add("Text Clipboard History", "Clipboard History Max Items", uint32(e.ClipboardTextBufferMaxItems))
	add("Text Clipboard History", "Clipboard History Max Display Width", uint32(e.ClipboardTextBufferMaxWidth))

	// Styles → TextViewColorPreferences / AXSViewColorPreferences sub-keys
	for i, s := range e.Styles {
		sfx := fmt.Sprintf("-%d", i)
		add(s.StyleKey, "StyleName"+sfx, s.StyleName)
		add(s.StyleKey, "StyleType"+sfx, uint32(s.StyleType))
		add(s.StyleKey, "Bold"+sfx, uint32(s.Bold))
		add(s.StyleKey, "Italic"+sfx, uint32(s.Italic))
		add(s.StyleKey, "Underline"+sfx, uint32(s.Underline))
		add(s.StyleKey, "BackgroundColor"+sfx, uint32(s.BackgroundColor))
		add(s.StyleKey, "ForegroundColor"+sfx, uint32(s.ForegroundColor))
		add(s.StyleKey, "Internal"+sfx, uint32(s.Internal))
	}

	// -----------------------------------------------------------------------
	// General Options  (covers WorkspaceSettings, TerminalSettings,
	//                    GeneralSettings)
	// -----------------------------------------------------------------------
	w := prefs.WorkspaceSettings
	add("General Options", "Enable Auto Time Stamp Source File", uint32(w.AutoDateStampFileEnabled))
	add("General Options", "Restore Workspace", uint32(w.RestoreWorkspaceEnabled))
	add("General Options", "Enable AutoSave", uint32(w.AutoSaveEnabled))
	add("General Options", "Enable Save Before Compile", uint32(w.SaveBeforeCompileEnabled))
	add("General Options", "Display CommConfig In System Identifier", uint32(w.EnableCommConfigInSystemID))
	add("General Options", "AutoSave Delay", uint32(w.AutoSaveDelay))
	add("General Options", "Enable Window Tabs", uint32(w.EnableWindowTabs))
	add("General Options", "Enable Icons In Tabs", uint32(w.EnableIconsInTabs))
	add("General Options", "Enable Close Button In Tabs", uint32(w.EnableCloseWindowInTabs))
	add("General Options", "Workspace-ShowIRTab", uint32(w.ShowIRTab))
	add("General Options", "Workspace-ShowZeroConfigTab", uint32(w.ShowZeroConfigTab))
	add("General Options", "Close Associated Workspace Files", uint32(w.PromptToCloseAssociatedFiles))
	add("General Options", "Close Removed File", uint32(w.PromptToCloseRemovedFile))

	t := prefs.TerminalSettings
	add("General Options", "Terminal Bckgrnd", uint32(t.TerminalActiveBackgroundClr))
	add("General Options", "Terminal Text", uint32(t.TerminalActiveTextClr))
	add("General Options", "Terminal Inactive Bckgrnd", uint32(t.TerminalInactiveBackgroundClr))
	add("General Options", "Terminal Inactive Text", uint32(t.TerminalInactiveTextClr))
	add("General Options", "Terminal Add LF", uint32(t.AddLineFeed))
	add("General Options", "Terminal Ctrl Chars to Port", uint32(t.CtrlCharsToPort))
	add("General Options", "TerminalWindow-FontNameSize", t.TerminalWindowFontName)
	add("General Options", "TerminalWindow-PointSize", uint32(t.TerminalWindowPointSize))
	add("General Options", "TelnetWindow-FontNameSize", t.TelnetWindowFontName)
	add("General Options", "TelnetWindow-PointSize", uint32(t.TelnetWindowPointSize))
	add("General Options", "Default TELNET Program", t.TelnetDefaultPGM)

	g := prefs.GeneralSettings
	add("General Options", "MaxMRUList", uint32(g.MRUListSizxeTesting))
	add("General Options", "OnlineTree-ShowOnlyPortCounts", uint32(g.OLTShowPortCounts))
	add("General Options", "OnlineTree-ShowNexusDevices", uint32(g.OLTShowNexusDevices))
	add("General Options", "OnlineTree-FontNameSize", g.OnlineTreeFontName)
	add("General Options", "OnlineTree-PointSize", uint32(g.OnlineTreePointSize))
	add("General Options", "ZeroConfig-FontNameSize", g.ZeroConfigFontName)
	add("General Options", "ZeroConfig-PointSize", uint32(g.ZeroConfigPointSize))
	add("General Options", "StatusWindow-FontNameSize", g.StatusWindowFontName)
	add("General Options", "StatusWindow-PointSize", uint32(g.StatusWindowPointSize))
	add("General Options", "AsyncWindow-FontNameSize", g.AsyncWindowFontName)
	add("General Options", "AsyncWindow-PointSize", uint32(g.AsyncWindowPointSize))
	add("General Options", "DiagsWindow-FontNameSize", g.DiagsWindowFontName)
	add("General Options", "DiagsWindow-PointSize", uint32(g.DiagsWindowPointSize))
	add("General Options", "ApplicationLook", uint32(g.ApplicationLook))
	add("General Options", "OutputBar Foreground Color-0", uint32(g.GCOStatusFGC))
	add("General Options", "OutputBar Foreground Color-1", uint32(g.GCOFindInFilesFGC))
	add("General Options", "OutputBar Foreground Color-2", uint32(g.GCOFindIRLFGC))
	add("General Options", "OutputBar Foreground Color-3", uint32(g.GCOFileTransferStatusFGC))
	add("General Options", "OutputBar Foreground Color-4", uint32(g.GCOAsyncNotificationsFGC))
	add("General Options", "OutputBar Foreground Color-5", uint32(g.GCODiagnosticsFGC))
	add("General Options", "OutputBar Background Color-0", uint32(g.GCOStatusBGC))
	add("General Options", "OutputBar Background Color-1", uint32(g.GCOFindInFilesBGC))
	add("General Options", "OutputBar Background Color-2", uint32(g.GCOFindIRLBGC))
	add("General Options", "OutputBar Background Color-3", uint32(g.GCOFileTransferStatusBGC))
	add("General Options", "OutputBar Background Color-4", uint32(g.GCOAsyncNotificationsBGC))
	add("General Options", "OutputBar Background Color-5", uint32(g.GCODiagnosticsBGC))
	add("General Options", "Workspace Background Color", uint32(g.GCOWorkspaceProjectTreeBGC))
	add("General Options", "Workspace Foreground Color", uint32(g.GCOWorkspaceProjectTreeFGC))
	add("General Options", "OnLine Tree Background Color", uint32(g.GCOOnlineTreeBGC))
	add("General Options", "OnLine Tree Foreground Color", uint32(g.GCOOnlineTreeFGC))
	add("General Options", "Watch Window Background Color", uint32(g.GCOWatchWindowBGC))
	add("General Options", "Watch Window Foreground Color", uint32(g.GCOWatchWindowFGC))
	add("General Options", "Telnet Window Background Color", uint32(g.GCOTelnetWindowBGC))
	add("General Options", "Telnet Window Foreground Color", uint32(g.GCOTelnetWindowFGC))
	add("General Options", "ZeroConfig Window Background Color", uint32(g.GCOZeroConfigWindowBGC))
	add("General Options", "ZeroConfig Window Foreground Color", uint32(g.GCOZeroConfigWindowFGC))

	// -----------------------------------------------------------------------
	// Compiler Options  (HKCU)
	// -----------------------------------------------------------------------
	c := prefs.NetlinxCompilerSettings
	add("NLXCompiler_Options", "BuildWithSource", uint32(c.BuildWithSource))
	add("NLXCompiler_Options", "BuildWithDebugInfo", uint32(c.BuildWithDebugInfo))
	add("NLXCompiler_Options", "PasswordProtect", uint32(c.PasswordProtect))
	add("NLXCompiler_Options", "EnableWC", uint32(c.EnableWCPreprocessor))
	add("NLXCompiler_Options", "Password", c.Password)
	add("NLXCompiler_Options", "ShowDebugWindow", uint32(c.ShowDebugWindowOnSessionClose))
	add("NLXCompiler_Options", "ShowMainAXS", uint32(c.ShowMainAXSOnSessionStart))

	// Directory lists live in HKLM WOW6432Node.
	const dirBase = `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\`
	for i, dir := range c.LibraryDirs {
		addHKLM(dirBase+"NLXCompiler_Libs", fmt.Sprintf("Dir%03d", i), dir)
	}

	for i, dir := range c.IncludeDirs {
		addHKLM(dirBase+"NLXCompiler_Includes", fmt.Sprintf("Dir%03d", i), dir)
	}

	for i, dir := range c.ModuleDirs {
		addHKLM(dirBase+"NLXCompiler_Modules", fmt.Sprintf("Dir%03d", i), dir)
	}

	// -----------------------------------------------------------------------
	// Batch Transfer Options
	// -----------------------------------------------------------------------
	ft := prefs.FileTransferSettings
	add("Batch Transfer User Options", "TKN Reboot", uint32(ft.TKNReboot))
	add("Batch Transfer User Options", "XDD Reboot", uint32(ft.XDDReboot))
	add("Batch Transfer User Options", "JAR Reboot", uint32(ft.JARReboot))
	add("Batch Transfer User Options", "TP4 Smart Transfer", uint32(ft.TP4SmartTransfer))
	add("Batch Transfer User Options", "TP5 Smart Transfer", uint32(ft.TP5SmartTransfer))
	add("Batch Transfer User Options", "Auto Send SRC", uint32(ft.AutoSendSRC))
	add("Batch Transfer User Options", "Report Number of Items Loaded From List", uint32(ft.ReportNumberofItemsLoadedFromList))

	// HttpServer fields go to the KIT Files sub-key.
	add("KIT Files", "HttpServerPort", uint32(ft.HttpServerPort))
	add("KIT Files", "HttpServerSelected", uint32(ft.HttpServerSelected))

	// -----------------------------------------------------------------------
	// Diagnostic Preferences
	// -----------------------------------------------------------------------
	d := prefs.DiagnosticsSettings
	add("Diagnostic Preferences", "DisplayLineNumbersNotification", uint32(d.DisplayLineNumbersNotification))
	add("Diagnostic Preferences", "DisplayLineNumbersDiagnostics", uint32(d.DisplayLineNumbersDiagnostics))
	add("Diagnostic Preferences", "DisplayTimeStampDiagnostics", uint32(d.DisplayTimeStampDiagnostics))
	add("Diagnostic Preferences", "DisplayTimeStampNotification", uint32(d.DisplayTimeStampNotification))
	add("Diagnostic Preferences", "DisplayTimeMillisecondsDiagnostics", uint32(d.DisplayTimeMillisecondsDiagnostics))
	add("Diagnostic Preferences", "DisplayTimeMillisecondsNotification", uint32(d.DisplayTimeMillisecondsNotification))
	add("Diagnostic Preferences", "DisplayMessagesNotification", uint32(d.DisplayMessagesNotification))
	add("Diagnostic Preferences", "DisplayMessagesDiagnostics", uint32(d.DisplayMessagesDiagnostics))
	add("Diagnostic Preferences", "LinesToRead", uint32(d.LinesToRead))
	add("Diagnostic Preferences", "BufferFileSize", uint32(d.BufferFileSize))
	add("Diagnostic Preferences", "DisplayDateStampNotification", uint32(d.DisplayDateStampNotification))
	add("Diagnostic Preferences", "DisplayDateStampDiagnostics", uint32(d.DisplayDateStampDiagnostics))

	// -----------------------------------------------------------------------
	// TCP/IP Connection History
	// -----------------------------------------------------------------------
	for i, entry := range prefs.TCPIPHistory.Entries {
		add("RecentConnectionsHistory",
			fmt.Sprintf("Recent Connection History%d", i),
			"T-"+entry.encode())
	}

	return entries
}
