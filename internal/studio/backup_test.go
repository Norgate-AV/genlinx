package studio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// ui32
// ---------------------------------------------------------------------------

type UI32Suite struct{ suite.Suite }

func TestUI32Suite(t *testing.T) { suite.Run(t, new(UI32Suite)) }

func (s *UI32Suite) TestUint32() {
	m := map[string]any{"k": uint32(42)}
	assert.Equal(s.T(), 42, ui32(m, "k"))
}

func (s *UI32Suite) TestUint64() {
	m := map[string]any{"k": uint64(99)}
	assert.Equal(s.T(), 99, ui32(m, "k"))
}

func (s *UI32Suite) TestMissingKey() {
	assert.Equal(s.T(), 0, ui32(map[string]any{}, "missing"))
}

func (s *UI32Suite) TestWrongType() {
	m := map[string]any{"k": "not-an-int"}
	assert.Equal(s.T(), 0, ui32(m, "k"))
}

// ---------------------------------------------------------------------------
// str
// ---------------------------------------------------------------------------

type StrSuite struct{ suite.Suite }

func TestStrSuite(t *testing.T) { suite.Run(t, new(StrSuite)) }

func (s *StrSuite) TestStringValue() {
	m := map[string]any{"k": "hello"}
	assert.Equal(s.T(), "hello", str(m, "k"))
}

func (s *StrSuite) TestMissingKey() {
	assert.Equal(s.T(), "", str(map[string]any{}, "missing"))
}

func (s *StrSuite) TestWrongType() {
	m := map[string]any{"k": uint32(5)}
	assert.Equal(s.T(), "", str(m, "k"))
}

func (s *StrSuite) TestEmptyString() {
	m := map[string]any{"k": ""}
	assert.Equal(s.T(), "", str(m, "k"))
}

// ---------------------------------------------------------------------------
// buildStyles
// ---------------------------------------------------------------------------

type BuildStylesSuite struct{ suite.Suite }

func TestBuildStylesSuite(t *testing.T) { suite.Run(t, new(BuildStylesSuite)) }

func (s *BuildStylesSuite) TestEmptyMap() {
	assert.Empty(s.T(), buildStyles(map[string]any{}, "Key"))
}

func (s *BuildStylesSuite) TestSingleStyle() {
	m := map[string]any{
		"StyleType-0":       uint32(1),
		"StyleName-0":       "Comment",
		"Bold-0":            uint32(0),
		"Italic-0":          uint32(1),
		"Underline-0":       uint32(0),
		"BackgroundColor-0": uint32(0xFFFFFF),
		"ForegroundColor-0": uint32(0x008000),
		"Internal-0":        uint32(0),
	}
	styles := buildStyles(m, "MyKey")
	s.Require().Len(styles, 1)
	assert.Equal(s.T(), "MyKey", styles[0].StyleKey)
	assert.Equal(s.T(), "Comment", styles[0].StyleName)
	assert.Equal(s.T(), 1, styles[0].StyleType)
	assert.Equal(s.T(), 0, styles[0].Bold)
	assert.Equal(s.T(), 1, styles[0].Italic)
	assert.Equal(s.T(), 0xFFFFFF, styles[0].BackgroundColor)
	assert.Equal(s.T(), 0x008000, styles[0].ForegroundColor)
}

func (s *BuildStylesSuite) TestMultipleStylesAreOrderedByIndex() {
	m := map[string]any{
		"StyleType-2": uint32(3),
		"StyleName-2": "Third",
		"StyleType-0": uint32(1),
		"StyleName-0": "First",
		"StyleType-1": uint32(2),
		"StyleName-1": "Second",
	}
	styles := buildStyles(m, "K")
	s.Require().Len(styles, 3)
	assert.Equal(s.T(), "First", styles[0].StyleName)
	assert.Equal(s.T(), "Second", styles[1].StyleName)
	assert.Equal(s.T(), "Third", styles[2].StyleName)
}

func (s *BuildStylesSuite) TestNonContiguousIndices() {
	m := map[string]any{
		"StyleType-5":  uint32(6),
		"StyleName-5":  "Five",
		"StyleType-10": uint32(11),
		"StyleName-10": "Ten",
	}
	styles := buildStyles(m, "K")
	s.Require().Len(styles, 2)
	assert.Equal(s.T(), "Five", styles[0].StyleName)
	assert.Equal(s.T(), "Ten", styles[1].StyleName)
}

func (s *BuildStylesSuite) TestStyleKeySetOnAllEntries() {
	m := map[string]any{
		"StyleType-0": uint32(1),
		"StyleType-1": uint32(2),
	}
	styles := buildStyles(m, "TheKey")
	for _, st := range styles {
		assert.Equal(s.T(), "TheKey", st.StyleKey)
	}
}

// ---------------------------------------------------------------------------
// buildTCPIPHistory
// ---------------------------------------------------------------------------

type BuildTCPIPHistorySuite struct{ suite.Suite }

func TestBuildTCPIPHistorySuite(t *testing.T) { suite.Run(t, new(BuildTCPIPHistorySuite)) }

func (s *BuildTCPIPHistorySuite) TestOnlyTCPEntriesIncluded() {
	m := map[string]any{
		"Recent Connection History0": "T-192.168.1.1|1319|1|Lab||",
		"Recent Connection History1": "S-COM1|9600|0|Serial||",
		"Recent Connection History2": "U-device|0|0|USB||",
		"Recent Connection History3": "T-10.0.0.1|1319|0|Office||",
	}
	h := buildTCPIPHistory(m)
	s.Require().Len(h.Entries, 2)
	assert.Equal(s.T(), "192.168.1.1", h.Entries[0].Host)
	assert.Equal(s.T(), "10.0.0.1", h.Entries[1].Host)
}

func (s *BuildTCPIPHistorySuite) TestNumericalIndexOrdering() {
	// Entries are stored in a map so iteration order is random; ordering must
	// be by the numeric suffix of the key, not map iteration order.
	m := map[string]any{
		"Recent Connection History10": "T-10.0.0.10|1319|1|Z||",
		"Recent Connection History2":  "T-10.0.0.2|1319|1|A||",
		"Recent Connection History1":  "T-10.0.0.1|1319|1|B||",
	}
	h := buildTCPIPHistory(m)
	s.Require().Len(h.Entries, 3)
	assert.Equal(s.T(), "10.0.0.1", h.Entries[0].Host)
	assert.Equal(s.T(), "10.0.0.2", h.Entries[1].Host)
	assert.Equal(s.T(), "10.0.0.10", h.Entries[2].Host)
}

func (s *BuildTCPIPHistorySuite) TestNonNumericSuffixSkipped() {
	m := map[string]any{
		"Recent Connection History0": "T-192.168.1.1|1319|1|Lab||",
		"Recent Connection HistoryX": "T-should-be-skipped|1319|1|||",
	}
	h := buildTCPIPHistory(m)
	s.Require().Len(h.Entries, 1)
	assert.Equal(s.T(), "192.168.1.1", h.Entries[0].Host)
}

func (s *BuildTCPIPHistorySuite) TestNonStringValueSkipped() {
	m := map[string]any{
		"Recent Connection History0": uint32(42),
		"Recent Connection History1": "T-192.168.1.1|1319|1|Lab||",
	}
	h := buildTCPIPHistory(m)
	s.Require().Len(h.Entries, 1)
	assert.Equal(s.T(), "192.168.1.1", h.Entries[0].Host)
}

func (s *BuildTCPIPHistorySuite) TestUnrelatedKeysSkipped() {
	m := map[string]any{
		"SomethingElse":              "T-should-skip|1319|1|||",
		"Recent Connection History0": "T-192.168.1.1|1319|1|Lab||",
	}
	h := buildTCPIPHistory(m)
	s.Require().Len(h.Entries, 1)
}

func (s *BuildTCPIPHistorySuite) TestEmptyMap() {
	h := buildTCPIPHistory(map[string]any{})
	assert.Empty(s.T(), h.Entries)
}

func (s *BuildTCPIPHistorySuite) TestPrefixStrippedFromValue() {
	m := map[string]any{
		"Recent Connection History0": "T-192.168.1.1|1319|1|Room||",
	}
	h := buildTCPIPHistory(m)
	s.Require().Len(h.Entries, 1)
	// The "T-" prefix must be stripped — the host should just be "192.168.1.1".
	assert.Equal(s.T(), "192.168.1.1", h.Entries[0].Host)
}

// ---------------------------------------------------------------------------
// BuildPreferences
// ---------------------------------------------------------------------------

type BuildPreferencesSuite struct{ suite.Suite }

func TestBuildPreferencesSuite(t *testing.T) { suite.Run(t, new(BuildPreferencesSuite)) }

// makeSettings returns a RegistrySettings with known values for all fields
// so every section of BuildPreferences can be exercised.
func (s *BuildPreferencesSuite) makeSettings() *RegistrySettings {
	return &RegistrySettings{
		EditorPreferences: map[string]any{
			"AutoIdent":                uint32(1),
			"ShowLineNumbers":          uint32(1),
			"CodeFolding":              uint32(0),
			"AutoSuggest":              uint32(1),
			"IndentGuides":             uint32(0),
			"CallTips":                 uint32(1),
			"UTF8FormatEnabled":        uint32(1),
			"IndentWidth":              uint32(4),
			"TabWidth":                 uint32(4),
			"PointSize":                uint32(10),
			"FontName":                 "Courier New",
			"UseOnlyFixedPitchFonts":   uint32(1),
			"PrinterColorMode":         uint32(2),
			"RightTrimLines":           uint32(1),
			"AutoSuggestInComments":    uint32(0),
			"ColumnEdgeMarker":         uint32(1),
			"ColumnEdgeColumnNumber":   uint32(80),
			"SaveBookmarksOnClose":     uint32(1),
			"HightlightMatchingBraces": uint32(1), // registry has this typo
			"UseTabs":                  uint32(0),
		},
		ClipboardHistory: map[string]any{
			"Clipboard History Max Items":         uint32(20),
			"Clipboard History Max Display Width": uint32(60),
		},
		DiagnosticPreferences: map[string]any{
			"DisplayLineNumbersNotification":      uint32(1),
			"DisplayLineNumbersDiagnostics":       uint32(1),
			"DisplayTimeStampDiagnostics":         uint32(0),
			"DisplayTimeStampNotification":        uint32(1),
			"DisplayTimeMillisecondsDiagnostics":  uint32(0),
			"DisplayTimeMillisecondsNotification": uint32(0),
			"DisplayMessagesNotification":         uint32(1),
			"DisplayMessagesDiagnostics":          uint32(1),
			"LinesToRead":                         uint32(500),
			"BufferFileSize":                      uint32(10),
			"DisplayDateStampNotification":        uint32(0),
			"DisplayDateStampDiagnostics":         uint32(0),
		},
		GeneralOptions: map[string]any{
			"Enable Auto Time Stamp Source File":      uint32(0),
			"Restore Workspace":                       uint32(1),
			"Enable AutoSave":                         uint32(1),
			"Enable Save Before Compile":              uint32(1),
			"Display CommConfig In System Identifier": uint32(0),
			"AutoSave Delay":                          uint32(5),
			"Enable Window Tabs":                      uint32(1),
			"Enable Icons In Tabs":                    uint32(1),
			"Enable Close Button In Tabs":             uint32(1),
			"Workspace-ShowIRTab":                     uint32(0),
			"Workspace-ShowZeroConfigTab":             uint32(1),
			"Close Associated Workspace Files":        uint32(1),
			"Close Removed File":                      uint32(0),
			"MaxMRUList":                              uint32(9),
			"Terminal Bckgrnd":                        uint32(0x000000),
			"Terminal Text":                           uint32(0xFFFFFF),
			"TerminalWindow-FontNameSize":             "Courier New",
			"TerminalWindow-PointSize":                uint32(9),
		},
		AXSViewColorPreferences: map[string]any{
			"StyleType-0": uint32(1),
			"StyleName-0": "AXSStyle",
		},
		TextViewColorPreferences: map[string]any{
			"StyleType-0": uint32(2),
			"StyleName-0": "TextStyle",
		},
		CompilerOptions: map[string]any{
			"BuildWithSource":    uint32(1),
			"BuildWithDebugInfo": uint32(0),
			"PasswordProtect":    uint32(0),
			"EnableWC":           uint32(1),
			"Password":           "",
			"ShowDebugWindow":    uint32(0),
			"ShowMainAXS":        uint32(1),
		},
		BatchTransferOptions: map[string]any{
			"TKN Reboot":         uint32(1),
			"XDD Reboot":         uint32(0),
			"JAR Reboot":         uint32(0),
			"TP4 Smart Transfer": uint32(1),
			"TP5 Smart Transfer": uint32(1),
			"Auto Send SRC":      uint32(0),
			"Report Number of Items Loaded From List": uint32(1),
		},
		KITFiles: map[string]any{
			"HttpServerPort":     uint32(80),
			"HttpServerSelected": uint32(1),
		},
		ConnectionHistory: map[string]any{
			"Recent Connection History0": "T-192.168.1.1|1319|1|Lab||",
			"Recent Connection History1": "S-COM1|9600|0|Serial||", // should be excluded
		},
		IncludeDirs: []string{`C:\AMX\NetLinx\include`},
		ModuleDirs:  []string{`C:\AMX\NetLinx\module`},
		LibraryDirs: []string{`C:\AMX\NetLinx\lib`},
	}
}

func (s *BuildPreferencesSuite) TestRegistryTypoMappedToCorrectEPXField() {
	prefs := BuildPreferences(s.makeSettings())
	// The registry key is "HightlightMatchingBraces" (missing 'l') but
	// the EPX element is "HighlightMatchingBraces" (correct spelling).
	assert.Equal(s.T(), 1, prefs.EditorSettings.HighlightMatchingBraces,
		"HightlightMatchingBraces (registry typo) must map to HighlightMatchingBraces")
}

func (s *BuildPreferencesSuite) TestHttpServerFieldsComefromKITFiles() {
	settings := s.makeSettings()
	// Ensure HttpServerPort is NOT in BatchTransferOptions so we confirm it
	// must come from KITFiles.
	delete(settings.BatchTransferOptions, "HttpServerPort")
	delete(settings.BatchTransferOptions, "HttpServerSelected")

	prefs := BuildPreferences(settings)
	assert.Equal(s.T(), 80, prefs.FileTransferSettings.HttpServerPort)
	assert.Equal(s.T(), 1, prefs.FileTransferSettings.HttpServerSelected)
}

func (s *BuildPreferencesSuite) TestClipboardFieldsFromClipboardHistory() {
	prefs := BuildPreferences(s.makeSettings())
	assert.Equal(s.T(), 20, prefs.EditorSettings.ClipboardTextBufferMaxItems)
	assert.Equal(s.T(), 60, prefs.EditorSettings.ClipboardTextBufferMaxWidth)
}

func (s *BuildPreferencesSuite) TestStylesOrderIsTextViewThenAXSView() {
	prefs := BuildPreferences(s.makeSettings())
	s.Require().Len(prefs.EditorSettings.Styles, 2)
	assert.Equal(s.T(), "TextViewColorPreferences", prefs.EditorSettings.Styles[0].StyleKey,
		"TextViewColorPreferences styles must appear before AXSViewColorPreferences styles")
	assert.Equal(s.T(), "AXSViewColorPreferences", prefs.EditorSettings.Styles[1].StyleKey)
}

func (s *BuildPreferencesSuite) TestDirListsPassThrough() {
	settings := s.makeSettings()
	prefs := BuildPreferences(settings)
	assert.Equal(s.T(), settings.IncludeDirs, prefs.NetlinxCompilerSettings.IncludeDirs)
	assert.Equal(s.T(), settings.ModuleDirs, prefs.NetlinxCompilerSettings.ModuleDirs)
	assert.Equal(s.T(), settings.LibraryDirs, prefs.NetlinxCompilerSettings.LibraryDirs)
}

func (s *BuildPreferencesSuite) TestTCPIPHistoryExcludesNonTCP() {
	prefs := BuildPreferences(s.makeSettings())
	// ConnectionHistory has one T- entry and one S- entry; only T- should appear.
	s.Require().Len(prefs.TCPIPHistory.Entries, 1)
	assert.Equal(s.T(), "192.168.1.1", prefs.TCPIPHistory.Entries[0].Host)
}

func (s *BuildPreferencesSuite) TestEditorSettings() {
	prefs := BuildPreferences(s.makeSettings())
	e := prefs.EditorSettings
	assert.Equal(s.T(), 1, e.AutoIndentEnabled)
	assert.Equal(s.T(), 1, e.ShowLineNumbersEnabled)
	assert.Equal(s.T(), 0, e.CodeFoldingEnabled)
	assert.Equal(s.T(), 4, e.IndentationWidth)
	assert.Equal(s.T(), 4, e.TabSpaces)
	assert.Equal(s.T(), 10, e.FontPointSize)
	assert.Equal(s.T(), "Courier New", e.FontName)
	assert.Equal(s.T(), 1, e.OnlyFixedPitchFontsEnabled)
	assert.Equal(s.T(), 80, e.ColumnEdgeMarkerValue)
	assert.Equal(s.T(), 0, e.UseTabEnabled)
}

func (s *BuildPreferencesSuite) TestWorkspaceSettings() {
	prefs := BuildPreferences(s.makeSettings())
	w := prefs.WorkspaceSettings
	assert.Equal(s.T(), 1, w.RestoreWorkspaceEnabled)
	assert.Equal(s.T(), 1, w.AutoSaveEnabled)
	assert.Equal(s.T(), 1, w.SaveBeforeCompileEnabled)
	assert.Equal(s.T(), 5, w.AutoSaveDelay)
	assert.Equal(s.T(), 1, w.ShowZeroConfigTab)
}

func (s *BuildPreferencesSuite) TestTerminalSettings() {
	prefs := BuildPreferences(s.makeSettings())
	t := prefs.TerminalSettings
	assert.Equal(s.T(), 0x000000, t.TerminalActiveBackgroundClr)
	assert.Equal(s.T(), 0xFFFFFF, t.TerminalActiveTextClr)
	assert.Equal(s.T(), "Courier New", t.TerminalWindowFontName)
	assert.Equal(s.T(), 9, t.TerminalWindowPointSize)
}

func (s *BuildPreferencesSuite) TestCompilerSettings() {
	prefs := BuildPreferences(s.makeSettings())
	c := prefs.NetlinxCompilerSettings
	assert.Equal(s.T(), 1, c.BuildWithSource)
	assert.Equal(s.T(), 0, c.BuildWithDebugInfo)
	assert.Equal(s.T(), 1, c.EnableWCPreprocessor)
	assert.Equal(s.T(), 1, c.ShowMainAXSOnSessionStart)
}

func (s *BuildPreferencesSuite) TestFileTransferSettings() {
	prefs := BuildPreferences(s.makeSettings())
	f := prefs.FileTransferSettings
	assert.Equal(s.T(), 1, f.TKNReboot)
	assert.Equal(s.T(), 0, f.XDDReboot)
	assert.Equal(s.T(), 1, f.TP4SmartTransfer)
	assert.Equal(s.T(), 1, f.TP5SmartTransfer)
	assert.Equal(s.T(), 1, f.ReportNumberofItemsLoadedFromList)
}

func (s *BuildPreferencesSuite) TestDiagnosticsSettings() {
	prefs := BuildPreferences(s.makeSettings())
	d := prefs.DiagnosticsSettings
	assert.Equal(s.T(), 1, d.DisplayLineNumbersNotification)
	assert.Equal(s.T(), 1, d.DisplayLineNumbersDiagnostics)
	assert.Equal(s.T(), 1, d.DisplayTimeStampNotification)
	assert.Equal(s.T(), 500, d.LinesToRead)
	assert.Equal(s.T(), 10, d.BufferFileSize)
}
