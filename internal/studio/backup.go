package studio

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// BuildPreferences constructs a Preferences struct from registry settings
// read by ReadRegistry. The result can be marshalled to a NetLinx Studio
// compatible .epx file using encoding/xml.
func BuildPreferences(s *RegistrySettings) *Preferences {
	// TextViewColorPreferences styles appear first in the EPX <Styles> block,
	// followed by AXSViewColorPreferences styles.
	styles := buildStyles(s.TextViewColorPreferences, "TextViewColorPreferences")
	styles = append(styles, buildStyles(s.AXSViewColorPreferences, "AXSViewColorPreferences")...)

	return &Preferences{
		EditorSettings: EditorSettings{
			ClipboardTextBufferMaxItems: ui32(s.ClipboardHistory, "Clipboard History Max Items"),
			ClipboardTextBufferMaxWidth: ui32(s.ClipboardHistory, "Clipboard History Max Display Width"),
			AutoIndentEnabled:           ui32(s.EditorPreferences, "AutoIdent"),
			ShowLineNumbersEnabled:      ui32(s.EditorPreferences, "ShowLineNumbers"),
			CodeFoldingEnabled:          ui32(s.EditorPreferences, "CodeFolding"),
			AutoSuggestEnabled:          ui32(s.EditorPreferences, "AutoSuggest"),
			IndentGuidesEnabled:         ui32(s.EditorPreferences, "IndentGuides"),
			CallTipsEnabled:             ui32(s.EditorPreferences, "CallTips"),
			UTF8FormatEnabled:           ui32(s.EditorPreferences, "UTF8FormatEnabled"),
			IndentationWidth:            ui32(s.EditorPreferences, "IndentWidth"),
			TabSpaces:                   ui32(s.EditorPreferences, "TabWidth"),
			FontPointSize:               ui32(s.EditorPreferences, "PointSize"),
			FontName:                    str(s.EditorPreferences, "FontName"),
			OnlyFixedPitchFontsEnabled:  ui32(s.EditorPreferences, "UseOnlyFixedPitchFonts"),
			PrinterColorMode:            ui32(s.EditorPreferences, "PrinterColorMode"),
			TrimLinesRight:              ui32(s.EditorPreferences, "RightTrimLines"),
			AutoSuggestInComments:       ui32(s.EditorPreferences, "AutoSuggestInComments"),
			ColumnEdgeMarker:            ui32(s.EditorPreferences, "ColumnEdgeMarker"),
			ColumnEdgeMarkerValue:       ui32(s.EditorPreferences, "ColumnEdgeColumnNumber"),
			RetainBookmarksOnFileClose:  ui32(s.EditorPreferences, "SaveBookmarksOnClose"),
			// Registry has a typo ("Hight" instead of "High") — map it correctly.
			HighlightMatchingBraces: ui32(s.EditorPreferences, "HightlightMatchingBraces"),
			UseTabEnabled:           ui32(s.EditorPreferences, "UseTabs"),
			Styles:                  styles,
		},
		WorkspaceSettings: WorkspaceSettings{
			AutoDateStampFileEnabled:     ui32(s.GeneralOptions, "Enable Auto Time Stamp Source File"),
			RestoreWorkspaceEnabled:      ui32(s.GeneralOptions, "Restore Workspace"),
			AutoSaveEnabled:              ui32(s.GeneralOptions, "Enable AutoSave"),
			SaveBeforeCompileEnabled:     ui32(s.GeneralOptions, "Enable Save Before Compile"),
			EnableCommConfigInSystemID:   ui32(s.GeneralOptions, "Display CommConfig In System Identifier"),
			AutoSaveDelay:                ui32(s.GeneralOptions, "AutoSave Delay"),
			EnableWindowTabs:             ui32(s.GeneralOptions, "Enable Window Tabs"),
			EnableIconsInTabs:            ui32(s.GeneralOptions, "Enable Icons In Tabs"),
			EnableCloseWindowInTabs:      ui32(s.GeneralOptions, "Enable Close Button In Tabs"),
			ShowIRTab:                    ui32(s.GeneralOptions, "Workspace-ShowIRTab"),
			ShowZeroConfigTab:            ui32(s.GeneralOptions, "Workspace-ShowZeroConfigTab"),
			PromptToCloseAssociatedFiles: ui32(s.GeneralOptions, "Close Associated Workspace Files"),
			PromptToCloseRemovedFile:     ui32(s.GeneralOptions, "Close Removed File"),
		},
		TerminalSettings: TerminalSettings{
			TerminalActiveBackgroundClr:   ui32(s.GeneralOptions, "Terminal Bckgrnd"),
			TerminalActiveTextClr:         ui32(s.GeneralOptions, "Terminal Text"),
			TerminalInactiveBackgroundClr: ui32(s.GeneralOptions, "Terminal Inactive Bckgrnd"),
			TerminalInactiveTextClr:       ui32(s.GeneralOptions, "Terminal Inactive Text"),
			AddLineFeed:                   ui32(s.GeneralOptions, "Terminal Add LF"),
			CtrlCharsToPort:               ui32(s.GeneralOptions, "Terminal Ctrl Chars to Port"),
			TerminalWindowFontName:        str(s.GeneralOptions, "TerminalWindow-FontNameSize"),
			TerminalWindowPointSize:       ui32(s.GeneralOptions, "TerminalWindow-PointSize"),
			TelnetWindowFontName:          str(s.GeneralOptions, "TelnetWindow-FontNameSize"),
			TelnetWindowPointSize:         ui32(s.GeneralOptions, "TelnetWindow-PointSize"),
			TelnetDefaultPGM:              str(s.GeneralOptions, "Default TELNET Program"),
		},
		GeneralSettings: GeneralSettings{
			MRUListSizxeTesting:        ui32(s.GeneralOptions, "MaxMRUList"),
			OLTShowPortCounts:          ui32(s.GeneralOptions, "OnlineTree-ShowOnlyPortCounts"),
			OLTShowNexusDevices:        ui32(s.GeneralOptions, "OnlineTree-ShowNexusDevices"),
			OnlineTreeFontName:         str(s.GeneralOptions, "OnlineTree-FontNameSize"),
			OnlineTreePointSize:        ui32(s.GeneralOptions, "OnlineTree-PointSize"),
			ZeroConfigFontName:         str(s.GeneralOptions, "ZeroConfig-FontNameSize"),
			ZeroConfigPointSize:        ui32(s.GeneralOptions, "ZeroConfig-PointSize"),
			StatusWindowFontName:       str(s.GeneralOptions, "StatusWindow-FontNameSize"),
			StatusWindowPointSize:      ui32(s.GeneralOptions, "StatusWindow-PointSize"),
			AsyncWindowFontName:        str(s.GeneralOptions, "AsyncWindow-FontNameSize"),
			AsyncWindowPointSize:       ui32(s.GeneralOptions, "AsyncWindow-PointSize"),
			DiagsWindowFontName:        str(s.GeneralOptions, "DiagsWindow-FontNameSize"),
			DiagsWindowPointSize:       ui32(s.GeneralOptions, "DiagsWindow-PointSize"),
			ApplicationLook:            ui32(s.GeneralOptions, "ApplicationLook"),
			GCOStatusFGC:               ui32(s.GeneralOptions, "OutputBar Foreground Color-0"),
			GCOFindInFilesFGC:          ui32(s.GeneralOptions, "OutputBar Foreground Color-1"),
			GCOFindIRLFGC:              ui32(s.GeneralOptions, "OutputBar Foreground Color-2"),
			GCOFileTransferStatusFGC:   ui32(s.GeneralOptions, "OutputBar Foreground Color-3"),
			GCOAsyncNotificationsFGC:   ui32(s.GeneralOptions, "OutputBar Foreground Color-4"),
			GCODiagnosticsFGC:          ui32(s.GeneralOptions, "OutputBar Foreground Color-5"),
			GCOStatusBGC:               ui32(s.GeneralOptions, "OutputBar Background Color-0"),
			GCOFindInFilesBGC:          ui32(s.GeneralOptions, "OutputBar Background Color-1"),
			GCOFindIRLBGC:              ui32(s.GeneralOptions, "OutputBar Background Color-2"),
			GCOFileTransferStatusBGC:   ui32(s.GeneralOptions, "OutputBar Background Color-3"),
			GCOAsyncNotificationsBGC:   ui32(s.GeneralOptions, "OutputBar Background Color-4"),
			GCODiagnosticsBGC:          ui32(s.GeneralOptions, "OutputBar Background Color-5"),
			GCOWorkspaceProjectTreeBGC: ui32(s.GeneralOptions, "Workspace Background Color"),
			GCOWorkspaceProjectTreeFGC: ui32(s.GeneralOptions, "Workspace Foreground Color"),
			GCOOnlineTreeBGC:           ui32(s.GeneralOptions, "OnLine Tree Background Color"),
			GCOOnlineTreeFGC:           ui32(s.GeneralOptions, "OnLine Tree Foreground Color"),
			GCOWatchWindowBGC:          ui32(s.GeneralOptions, "Watch Window Background Color"),
			GCOWatchWindowFGC:          ui32(s.GeneralOptions, "Watch Window Foreground Color"),
			GCOTelnetWindowBGC:         ui32(s.GeneralOptions, "Telnet Window Background Color"),
			GCOTelnetWindowFGC:         ui32(s.GeneralOptions, "Telnet Window Foreground Color"),
			GCOZeroConfigWindowBGC:     ui32(s.GeneralOptions, "ZeroConfig Window Background Color"),
			GCOZeroConfigWindowFGC:     ui32(s.GeneralOptions, "ZeroConfig Window Foreground Color"),
		},
		NetlinxCompilerSettings: NetlinxCompilerSettings{
			BuildWithSource:               ui32(s.CompilerOptions, "BuildWithSource"),
			BuildWithDebugInfo:            ui32(s.CompilerOptions, "BuildWithDebugInfo"),
			PasswordProtect:               ui32(s.CompilerOptions, "PasswordProtect"),
			EnableWCPreprocessor:          ui32(s.CompilerOptions, "EnableWC"),
			Password:                      str(s.CompilerOptions, "Password"),
			ShowDebugWindowOnSessionClose: ui32(s.CompilerOptions, "ShowDebugWindow"),
			ShowMainAXSOnSessionStart:     ui32(s.CompilerOptions, "ShowMainAXS"),
			LibraryDirs:                   s.LibraryDirs,
			IncludeDirs:                   s.IncludeDirs,
			ModuleDirs:                    s.ModuleDirs,
		},
		FileTransferSettings: FileTransferSettings{
			TKNReboot:                         ui32(s.BatchTransferOptions, "TKN Reboot"),
			XDDReboot:                         ui32(s.BatchTransferOptions, "XDD Reboot"),
			JARReboot:                         ui32(s.BatchTransferOptions, "JAR Reboot"),
			TP4SmartTransfer:                  ui32(s.BatchTransferOptions, "TP4 Smart Transfer"),
			TP5SmartTransfer:                  ui32(s.BatchTransferOptions, "TP5 Smart Transfer"),
			AutoSendSRC:                       ui32(s.BatchTransferOptions, "Auto Send SRC"),
			ReportNumberofItemsLoadedFromList: ui32(s.BatchTransferOptions, "Report Number of Items Loaded From List"),
			// HttpServer fields live in the KIT Files key, not Batch Transfer.
			HttpServerPort:     ui32(s.KITFiles, "HttpServerPort"),
			HttpServerSelected: ui32(s.KITFiles, "HttpServerSelected"),
		},
		DiagnosticsSettings: DiagnosticsSettings{
			DisplayLineNumbersNotification:      ui32(s.DiagnosticPreferences, "DisplayLineNumbersNotification"),
			DisplayLineNumbersDiagnostics:       ui32(s.DiagnosticPreferences, "DisplayLineNumbersDiagnostics"),
			DisplayTimeStampDiagnostics:         ui32(s.DiagnosticPreferences, "DisplayTimeStampDiagnostics"),
			DisplayTimeStampNotification:        ui32(s.DiagnosticPreferences, "DisplayTimeStampNotification"),
			DisplayTimeMillisecondsDiagnostics:  ui32(s.DiagnosticPreferences, "DisplayTimeMillisecondsDiagnostics"),
			DisplayTimeMillisecondsNotification: ui32(s.DiagnosticPreferences, "DisplayTimeMillisecondsNotification"),
			DisplayMessagesNotification:         ui32(s.DiagnosticPreferences, "DisplayMessagesNotification"),
			DisplayMessagesDiagnostics:          ui32(s.DiagnosticPreferences, "DisplayMessagesDiagnostics"),
			LinesToRead:                         ui32(s.DiagnosticPreferences, "LinesToRead"),
			BufferFileSize:                      ui32(s.DiagnosticPreferences, "BufferFileSize"),
			DisplayDateStampNotification:        ui32(s.DiagnosticPreferences, "DisplayDateStampNotification"),
			DisplayDateStampDiagnostics:         ui32(s.DiagnosticPreferences, "DisplayDateStampDiagnostics"),
		},
		TCPIPHistory: buildTCPIPHistory(s.ConnectionHistory),
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// ui32 extracts a uint32/uint64 value from a registry map and returns it as int.
func ui32(m map[string]any, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case uint32:
			return int(n)
		case uint64:
			return int(n)
		}
	}

	return 0
}

// str extracts a string value from a registry map.
func str(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

// buildStyles reads colour/font style entries from a flat registry map (keyed
// as "Bold-N", "ForegroundColor-N", etc.) and returns ordered Style values.
func buildStyles(prefs map[string]any, styleKey string) []Style {
	seen := map[int]bool{}

	for k := range prefs {
		if after, ok := strings.CutPrefix(k, "StyleType-"); ok {
			if n, err := strconv.Atoi(after); err == nil {
				seen[n] = true
			}
		}
	}

	indices := make([]int, 0, len(seen))
	for n := range seen {
		indices = append(indices, n)
	}

	sort.Ints(indices)

	styles := make([]Style, 0, len(indices))
	for _, n := range indices {
		sfx := fmt.Sprintf("-%d", n)
		styles = append(styles, Style{
			StyleKey:        styleKey,
			StyleName:       str(prefs, "StyleName"+sfx),
			StyleType:       ui32(prefs, "StyleType"+sfx),
			Bold:            ui32(prefs, "Bold"+sfx),
			Italic:          ui32(prefs, "Italic"+sfx),
			Underline:       ui32(prefs, "Underline"+sfx),
			BackgroundColor: ui32(prefs, "BackgroundColor"+sfx),
			ForegroundColor: ui32(prefs, "ForegroundColor"+sfx),
			Internal:        ui32(prefs, "Internal"+sfx),
		})
	}

	return styles
}

// buildTCPIPHistory extracts TCP/IP connection history entries from the
// RecentConnectionsHistory registry map. Serial (S-) and USB (U-) entries
// are silently skipped to match the EPX file format which is TCP-only.
func buildTCPIPHistory(history map[string]any) TCPIPHistory {
	type indexed struct {
		n   int
		raw string
	}

	const prefix = "Recent Connection History"
	var entries []indexed

	for k, v := range history {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		n, err := strconv.Atoi(strings.TrimPrefix(k, prefix))
		if err != nil {
			continue
		}

		s, ok := v.(string)
		if !ok || !strings.HasPrefix(s, "T-") {
			continue
		}

		entries = append(entries, indexed{n: n, raw: strings.TrimPrefix(s, "T-")})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].n < entries[j].n })

	h := TCPIPHistory{}
	for _, e := range entries {
		h.Entries = append(h.Entries, parseTCPIPEntry(e.raw))
	}

	return h
}
