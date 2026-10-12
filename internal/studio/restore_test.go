package studio

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// ParsePreferences
// ---------------------------------------------------------------------------

type ParsePreferencesSuite struct{ suite.Suite }

func TestParsePreferencesSuite(t *testing.T) { suite.Run(t, new(ParsePreferencesSuite)) }

// marshalPrefs marshals prefs to XML bytes as the backup command would.
func marshalPrefs(t *testing.T, prefs *Preferences) []byte {
	t.Helper()

	b, err := xml.MarshalIndent(prefs, "", "    ")
	require.NoError(t, err)

	content := append([]byte(xml.Header), b...)
	return append(content, '\n')
}

func (s *ParsePreferencesSuite) TestRoundTripEditorSettings() {
	original := BuildPreferences(fullSettings())
	data := marshalPrefs(s.T(), original)

	parsed, err := ParsePreferences(bytes.NewReader(data))
	s.Require().NoError(err)

	e := parsed.EditorSettings
	assert.Equal(s.T(), original.EditorSettings.AutoIndentEnabled, e.AutoIndentEnabled)
	assert.Equal(s.T(), original.EditorSettings.IndentationWidth, e.IndentationWidth)
	assert.Equal(s.T(), original.EditorSettings.FontName, e.FontName)
	assert.Equal(s.T(), original.EditorSettings.HighlightMatchingBraces, e.HighlightMatchingBraces)
	assert.Equal(s.T(), original.EditorSettings.UseTabEnabled, e.UseTabEnabled)
}

func (s *ParsePreferencesSuite) TestRoundTripWorkspaceSettings() {
	original := BuildPreferences(fullSettings())
	data := marshalPrefs(s.T(), original)

	parsed, err := ParsePreferences(bytes.NewReader(data))
	s.Require().NoError(err)

	assert.Equal(s.T(), original.WorkspaceSettings, parsed.WorkspaceSettings)
}

func (s *ParsePreferencesSuite) TestRoundTripTerminalSettings() {
	original := BuildPreferences(fullSettings())
	data := marshalPrefs(s.T(), original)

	parsed, err := ParsePreferences(bytes.NewReader(data))
	s.Require().NoError(err)

	assert.Equal(s.T(), original.TerminalSettings, parsed.TerminalSettings)
}

func (s *ParsePreferencesSuite) TestRoundTripCompilerDirLists() {
	original := BuildPreferences(fullSettings())
	data := marshalPrefs(s.T(), original)

	parsed, err := ParsePreferences(bytes.NewReader(data))
	s.Require().NoError(err)

	c := parsed.NetlinxCompilerSettings
	assert.Equal(s.T(), original.NetlinxCompilerSettings.LibraryDirs, c.LibraryDirs)
	assert.Equal(s.T(), original.NetlinxCompilerSettings.IncludeDirs, c.IncludeDirs)
	assert.Equal(s.T(), original.NetlinxCompilerSettings.ModuleDirs, c.ModuleDirs)
}

func (s *ParsePreferencesSuite) TestRoundTripTCPIPHistory() {
	original := BuildPreferences(fullSettings())
	data := marshalPrefs(s.T(), original)

	parsed, err := ParsePreferences(bytes.NewReader(data))
	s.Require().NoError(err)

	s.Require().Len(parsed.TCPIPHistory.Entries, len(original.TCPIPHistory.Entries))
	assert.Equal(s.T(), original.TCPIPHistory.Entries[0].Host, parsed.TCPIPHistory.Entries[0].Host)
	assert.Equal(s.T(), original.TCPIPHistory.Entries[0].Port, parsed.TCPIPHistory.Entries[0].Port)
}

func (s *ParsePreferencesSuite) TestRoundTripStyles() {
	original := BuildPreferences(fullSettings())
	data := marshalPrefs(s.T(), original)

	parsed, err := ParsePreferences(bytes.NewReader(data))
	s.Require().NoError(err)

	s.Require().Len(parsed.EditorSettings.Styles, len(original.EditorSettings.Styles))
	for i, orig := range original.EditorSettings.Styles {
		got := parsed.EditorSettings.Styles[i]
		assert.Equal(s.T(), orig.StyleKey, got.StyleKey)
		assert.Equal(s.T(), orig.StyleName, got.StyleName)
		assert.Equal(s.T(), orig.ForegroundColor, got.ForegroundColor)
		assert.Equal(s.T(), orig.BackgroundColor, got.BackgroundColor)
	}
}

func (s *ParsePreferencesSuite) TestInvalidXMLReturnsError() {
	_, err := ParsePreferences(strings.NewReader("this is not xml"))
	assert.Error(s.T(), err)
}

func (s *ParsePreferencesSuite) TestEmptyInputReturnsError() {
	_, err := ParsePreferences(strings.NewReader(""))
	assert.Error(s.T(), err)
}

// ---------------------------------------------------------------------------
// PreferencesToRegistryEntries
// ---------------------------------------------------------------------------

type PreferencesToRegistryEntriesSuite struct{ suite.Suite }

func TestPreferencesToRegistryEntriesSuite(t *testing.T) {
	suite.Run(t, new(PreferencesToRegistryEntriesSuite))
}

// findEntry returns the first RegistryEntry matching subKey and valueName.
func findEntry(entries []RegistryEntry, subKey, valueName string) (RegistryEntry, bool) {
	for _, e := range entries {
		if e.SubKey == subKey && e.ValueName == valueName {
			return e, true
		}
	}

	return RegistryEntry{}, false
}

func (s *PreferencesToRegistryEntriesSuite) TestEditorIntFieldMappedAsUint32() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	e, ok := findEntry(entries, "Editor Preferences", "AutoIdent")
	s.Require().True(ok, "AutoIdent entry not found")
	assert.Equal(s.T(), uint32(1), e.Value)
}

func (s *PreferencesToRegistryEntriesSuite) TestEditorStringFieldMappedAsString() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	e, ok := findEntry(entries, "Editor Preferences", "FontName")
	s.Require().True(ok, "FontName entry not found")
	assert.Equal(s.T(), "Courier New", e.Value)
}

func (s *PreferencesToRegistryEntriesSuite) TestRegistryTypoPreservedOnWrite() {
	// The EPX field is HighlightMatchingBraces but the registry value name
	// must be written back as HightlightMatchingBraces (typo preserved).
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	_, wrongOk := findEntry(entries, "Editor Preferences", "HighlightMatchingBraces")
	assert.False(s.T(), wrongOk, "must NOT write 'HighlightMatchingBraces' (correct spelling)")

	e, ok := findEntry(entries, "Editor Preferences", "HightlightMatchingBraces")
	s.Require().True(ok, "HightlightMatchingBraces entry not found")
	assert.Equal(s.T(), uint32(1), e.Value)
}

func (s *PreferencesToRegistryEntriesSuite) TestClipboardHistoryGoesToCorrectSubKey() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	e, ok := findEntry(entries, "Text Clipboard History", "Clipboard History Max Items")
	s.Require().True(ok)
	assert.Equal(s.T(), uint32(20), e.Value)
}

func (s *PreferencesToRegistryEntriesSuite) TestHttpServerGoesToKITFiles() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	e, ok := findEntry(entries, "KIT Files", "HttpServerPort")
	s.Require().True(ok)
	assert.Equal(s.T(), uint32(80), e.Value)

	_, inBatch := findEntry(entries, "Batch Transfer User Options", "HttpServerPort")
	assert.False(s.T(), inBatch, "HttpServerPort must not appear in Batch Transfer User Options")
}

func (s *PreferencesToRegistryEntriesSuite) TestDirListsNotInRegistryEntries() {
	// Directory lists are handled via DiffDirList / replaceDirList, not as
	// individual RegistryEntry values.
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	const libKey = `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\NLXCompiler_Libs`
	_, ok := findEntry(entries, libKey, "Dir000")
	assert.False(s.T(), ok, "NLXCompiler_Libs must not appear in RegistryEntry list")
}

func (s *PreferencesToRegistryEntriesSuite) TestIncludeDirsNotInRegistryEntries() {
	// Include dirs are handled via DiffDirList / replaceDirList.
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	const incKey = `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\NLXCompiler_Includes`
	_, ok := findEntry(entries, incKey, "Dir000")
	assert.False(s.T(), ok, "NLXCompiler_Includes must not appear in RegistryEntry list")
}

func (s *PreferencesToRegistryEntriesSuite) TestTCPIPHistoryNotIncluded() {
	// TCP/IP history is excluded from restore: Studio manages it dynamically.
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	_, ok := findEntry(entries, "RecentConnectionsHistory", "Recent Connection History0")
	assert.False(s.T(), ok, "TCP/IP history must not appear in restore entries")
}

func (s *PreferencesToRegistryEntriesSuite) TestStyleEntriesHaveCorrectSubKey() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	// TextViewColorPreferences style should go to its own sub-key.
	found := false
	for _, e := range entries {
		if e.SubKey == "TextViewColorPreferences" {
			found = true
			break
		}
	}

	assert.True(s.T(), found, "expected entries with SubKey == TextViewColorPreferences")
}

func (s *PreferencesToRegistryEntriesSuite) TestAllValuesAreUint32OrString() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)

	for _, e := range entries {
		switch e.Value.(type) {
		case uint32, string:
			// ok
		default:
			s.Failf("unexpected value type", "entry %s\\%s has type %T", e.SubKey, e.ValueName, e.Value)
		}
	}
}

// ---------------------------------------------------------------------------
// helpers shared across restore_test.go
// ---------------------------------------------------------------------------

// fullSettings returns a RegistrySettings with concrete values for every
// field so that all mapping paths in BuildPreferences are exercised.
func fullSettings() *RegistrySettings {
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
			"HightlightMatchingBraces": uint32(1),
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
			"Terminal Inactive Bckgrnd":               uint32(0x111111),
			"Terminal Inactive Text":                  uint32(0xEEEEEE),
			"Terminal Add LF":                         uint32(0),
			"Terminal Ctrl Chars to Port":             uint32(0),
			"TerminalWindow-FontNameSize":             "Courier New",
			"TerminalWindow-PointSize":                uint32(9),
			"TelnetWindow-FontNameSize":               "Courier New",
			"TelnetWindow-PointSize":                  uint32(9),
			"Default TELNET Program":                  "",
			"OnlineTree-ShowOnlyPortCounts":           uint32(0),
			"OnlineTree-ShowNexusDevices":             uint32(1),
			"OnlineTree-FontNameSize":                 "Arial",
			"OnlineTree-PointSize":                    uint32(9),
			"ZeroConfig-FontNameSize":                 "Arial",
			"ZeroConfig-PointSize":                    uint32(9),
			"StatusWindow-FontNameSize":               "Courier New",
			"StatusWindow-PointSize":                  uint32(9),
			"AsyncWindow-FontNameSize":                "Courier New",
			"AsyncWindow-PointSize":                   uint32(9),
			"DiagsWindow-FontNameSize":                "Courier New",
			"DiagsWindow-PointSize":                   uint32(9),
			"ApplicationLook":                         uint32(0),
			"OutputBar Foreground Color-0":            uint32(0),
			"OutputBar Foreground Color-1":            uint32(1),
			"OutputBar Foreground Color-2":            uint32(2),
			"OutputBar Foreground Color-3":            uint32(3),
			"OutputBar Foreground Color-4":            uint32(4),
			"OutputBar Foreground Color-5":            uint32(5),
			"OutputBar Background Color-0":            uint32(10),
			"OutputBar Background Color-1":            uint32(11),
			"OutputBar Background Color-2":            uint32(12),
			"OutputBar Background Color-3":            uint32(13),
			"OutputBar Background Color-4":            uint32(14),
			"OutputBar Background Color-5":            uint32(15),
			"Workspace Background Color":              uint32(0xFFFFFF),
			"Workspace Foreground Color":              uint32(0x000000),
			"OnLine Tree Background Color":            uint32(0xFFFFFF),
			"OnLine Tree Foreground Color":            uint32(0x000000),
			"Watch Window Background Color":           uint32(0xFFFFFF),
			"Watch Window Foreground Color":           uint32(0x000000),
			"Telnet Window Background Color":          uint32(0xFFFFFF),
			"Telnet Window Foreground Color":          uint32(0x000000),
			"ZeroConfig Window Background Color":      uint32(0xFFFFFF),
			"ZeroConfig Window Foreground Color":      uint32(0x000000),
		},
		AXSViewColorPreferences: map[string]any{
			"StyleType-0":       uint32(1),
			"StyleName-0":       "AXSStyle",
			"Bold-0":            uint32(0),
			"Italic-0":          uint32(0),
			"Underline-0":       uint32(0),
			"BackgroundColor-0": uint32(0xFFFFFF),
			"ForegroundColor-0": uint32(0x000000),
			"Internal-0":        uint32(0),
		},
		TextViewColorPreferences: map[string]any{
			"StyleType-0":       uint32(2),
			"StyleName-0":       "TextStyle",
			"Bold-0":            uint32(1),
			"Italic-0":          uint32(0),
			"Underline-0":       uint32(0),
			"BackgroundColor-0": uint32(0xFFFFFF),
			"ForegroundColor-0": uint32(0x008000),
			"Internal-0":        uint32(0),
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
		},
		IncludeDirs: []string{`C:\AMX\NetLinx\include`},
		ModuleDirs:  []string{`C:\AMX\NetLinx\module`},
		LibraryDirs: []string{`C:\AMX\NetLinx\lib`},
	}
}

// ---------------------------------------------------------------------------
// DiffRegistryEntries
// ---------------------------------------------------------------------------

type DiffRegistryEntriesSuite struct{ suite.Suite }

func TestDiffRegistryEntriesSuite(t *testing.T) { suite.Run(t, new(DiffRegistryEntriesSuite)) }

func (s *DiffRegistryEntriesSuite) TestIdenticalEntriesProduceEmptyDiff() {
	prefs := BuildPreferences(fullSettings())
	entries := PreferencesToRegistryEntries(prefs)
	diff := DiffRegistryEntries(entries, entries)
	assert.Empty(s.T(), diff)
}

func (s *DiffRegistryEntriesSuite) TestChangedValueProducesDiffChangedEntry() {
	current := []RegistryEntry{
		{SubKey: "Editor Preferences", ValueName: "IndentWidth", Value: uint32(2)},
	}
	incoming := []RegistryEntry{
		{SubKey: "Editor Preferences", ValueName: "IndentWidth", Value: uint32(4)},
	}
	diff := DiffRegistryEntries(current, incoming)
	s.Require().Len(diff, 1)
	assert.Equal(s.T(), DiffChanged, diff[0].Status)
	assert.Equal(s.T(), uint32(2), diff[0].OldValue)
	assert.Equal(s.T(), uint32(4), diff[0].Value)
}

func (s *DiffRegistryEntriesSuite) TestNewEntryProducesDiffAddedEntry() {
	incoming := []RegistryEntry{
		{SubKey: "Editor Preferences", ValueName: "IndentWidth", Value: uint32(4)},
	}
	diff := DiffRegistryEntries([]RegistryEntry{}, incoming)
	s.Require().Len(diff, 1)
	assert.Equal(s.T(), DiffAdded, diff[0].Status)
	assert.Nil(s.T(), diff[0].OldValue)
}

func (s *DiffRegistryEntriesSuite) TestUnchangedEntryIsExcluded() {
	entry := RegistryEntry{SubKey: "Editor Preferences", ValueName: "IndentWidth", Value: uint32(4)}
	diff := DiffRegistryEntries([]RegistryEntry{entry}, []RegistryEntry{entry})
	assert.Empty(s.T(), diff)
}

func (s *DiffRegistryEntriesSuite) TestMixedEntriesOnlyReturnsDelta() {
	current := []RegistryEntry{
		{SubKey: "Editor Preferences", ValueName: "IndentWidth", Value: uint32(2)},
		{SubKey: "Editor Preferences", ValueName: "FontName", Value: "Courier New"},
	}
	incoming := []RegistryEntry{
		{SubKey: "Editor Preferences", ValueName: "IndentWidth", Value: uint32(4)},  // changed
		{SubKey: "Editor Preferences", ValueName: "FontName", Value: "Courier New"}, // same
		{SubKey: "Editor Preferences", ValueName: "TabWidth", Value: uint32(4)},     // added
	}
	diff := DiffRegistryEntries(current, incoming)
	s.Require().Len(diff, 2)
}

func (s *DiffRegistryEntriesSuite) TestHiveLMKeyedSeparatelyFromHKCU() {
	current := []RegistryEntry{
		{SubKey: "somekey", ValueName: "Val", Value: uint32(1), HiveLM: false},
		{SubKey: "somekey", ValueName: "Val", Value: uint32(2), HiveLM: true},
	}
	incoming := []RegistryEntry{
		{SubKey: "somekey", ValueName: "Val", Value: uint32(1), HiveLM: false}, // same
		{SubKey: "somekey", ValueName: "Val", Value: uint32(9), HiveLM: true},  // changed
	}
	diff := DiffRegistryEntries(current, incoming)
	s.Require().Len(diff, 1)
	assert.True(s.T(), diff[0].HiveLM)
	assert.Equal(s.T(), DiffChanged, diff[0].Status)
}

func (s *DiffRegistryEntriesSuite) TestEmptyIncomingProducesEmptyDiff() {
	prefs := BuildPreferences(fullSettings())
	current := PreferencesToRegistryEntries(prefs)
	diff := DiffRegistryEntries(current, []RegistryEntry{})
	assert.Empty(s.T(), diff)
}

// ---------------------------------------------------------------------------
// DiffTCPIPHistory
// ---------------------------------------------------------------------------

type DiffTCPIPHistorySuite struct{ suite.Suite }

func TestDiffTCPIPHistorySuite(t *testing.T) { suite.Run(t, new(DiffTCPIPHistorySuite)) }

func tcpEntries(hosts ...string) []TCPIPEntry {
	entries := make([]TCPIPEntry, len(hosts))
	for i, h := range hosts {
		entries[i] = TCPIPEntry{Host: h, Port: 1319, PingTest: true}
	}

	return entries
}

func (s *DiffTCPIPHistorySuite) TestIdenticalReturnsNil() {
	e := tcpEntries("10.0.0.1", "10.0.0.2")
	assert.Nil(s.T(), DiffTCPIPHistory(e, e))
}

func (s *DiffTCPIPHistorySuite) TestBothEmptyReturnsNil() {
	assert.Nil(s.T(), DiffTCPIPHistory(nil, nil))
}

func (s *DiffTCPIPHistorySuite) TestIncomingEntryAdded() {
	current := tcpEntries("10.0.0.1")
	incoming := tcpEntries("10.0.0.1", "10.0.0.2")

	d := DiffTCPIPHistory(current, incoming)
	s.Require().NotNil(d)
	assert.Equal(s.T(), DiffReplaced, d.Status)
	assert.Equal(s.T(), 1, d.OldValue)
	assert.Len(s.T(), d.TCPIPReplace, 2)
}

func (s *DiffTCPIPHistorySuite) TestEntryValueChanged() {
	d := DiffTCPIPHistory(tcpEntries("10.0.0.1"), tcpEntries("10.0.0.99"))
	s.Require().NotNil(d)
	assert.Equal(s.T(), DiffReplaced, d.Status)
}

func (s *DiffTCPIPHistorySuite) TestOrderChangeDetected() {
	current := tcpEntries("10.0.0.1", "10.0.0.2")
	incoming := tcpEntries("10.0.0.2", "10.0.0.1")

	d := DiffTCPIPHistory(current, incoming)
	s.Require().NotNil(d)
	assert.Equal(s.T(), DiffReplaced, d.Status)
}

func (s *DiffTCPIPHistorySuite) TestIncomingEntriesCarried() {
	d := DiffTCPIPHistory(tcpEntries("10.0.0.1", "10.0.0.2"), tcpEntries("10.0.0.3"))
	s.Require().NotNil(d)
	s.Require().Len(d.TCPIPReplace, 1)
	assert.Equal(s.T(), "10.0.0.3", d.TCPIPReplace[0].Host)
}

func (s *DiffTCPIPHistorySuite) TestSubKeyIsRecentConnectionsHistory() {
	d := DiffTCPIPHistory(tcpEntries("a"), tcpEntries("b"))
	s.Require().NotNil(d)
	assert.Equal(s.T(), "RecentConnectionsHistory", d.SubKey)
}

func (s *DiffTCPIPHistorySuite) TestOldValueIsCurrentCount() {
	d := DiffTCPIPHistory(tcpEntries("a", "b", "c"), tcpEntries("d"))
	s.Require().NotNil(d)
	assert.Equal(s.T(), 3, d.OldValue)
}

// ---------------------------------------------------------------------------
// DiffDirList
// ---------------------------------------------------------------------------

type DiffDirListSuite struct{ suite.Suite }

func TestDiffDirListSuite(t *testing.T) { suite.Run(t, new(DiffDirListSuite)) }

const testDirSubKey = `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\NLXCompiler_Libs`

func (s *DiffDirListSuite) TestIdenticalReturnsNil() {
	dirs := []string{`C:\AMX\lib`, `C:\AMX\lib2`}
	assert.Nil(s.T(), DiffDirList(testDirSubKey, dirs, dirs))
}

func (s *DiffDirListSuite) TestBothEmptyReturnsNil() {
	assert.Nil(s.T(), DiffDirList(testDirSubKey, nil, nil))
}

func (s *DiffDirListSuite) TestIncomingEntryAdded() {
	current := []string{`C:\AMX\lib`}
	incoming := []string{`C:\AMX\lib`, `C:\AMX\lib2`}

	d := DiffDirList(testDirSubKey, current, incoming)
	s.Require().NotNil(d)
	assert.Equal(s.T(), DiffReplaced, d.Status)
	assert.Equal(s.T(), 1, d.OldValue)
	assert.Len(s.T(), d.DirReplace, 2)
}

func (s *DiffDirListSuite) TestEntryValueChanged() {
	d := DiffDirList(testDirSubKey, []string{`C:\AMX\lib`}, []string{`C:\OTHER\lib`})
	s.Require().NotNil(d)
	assert.Equal(s.T(), DiffReplaced, d.Status)
}

func (s *DiffDirListSuite) TestOrderChangeDetected() {
	current := []string{`C:\AMX\lib`, `C:\AMX\lib2`}
	incoming := []string{`C:\AMX\lib2`, `C:\AMX\lib`}

	d := DiffDirList(testDirSubKey, current, incoming)
	s.Require().NotNil(d)
	assert.Equal(s.T(), DiffReplaced, d.Status)
}

func (s *DiffDirListSuite) TestIncomingDirsCarried() {
	d := DiffDirList(testDirSubKey, []string{`C:\old\lib`}, []string{`C:\new\lib`, `C:\new\lib2`})
	s.Require().NotNil(d)
	s.Require().Len(d.DirReplace, 2)
	assert.Equal(s.T(), `C:\new\lib`, d.DirReplace[0])
	assert.Equal(s.T(), `C:\new\lib2`, d.DirReplace[1])
}

func (s *DiffDirListSuite) TestHiveLMIsSet() {
	d := DiffDirList(testDirSubKey, []string{`C:\a`}, []string{`C:\b`})
	s.Require().NotNil(d)
	assert.True(s.T(), d.HiveLM, "HKLM dir list entries must have HiveLM=true")
}

func (s *DiffDirListSuite) TestSubKeyPreserved() {
	d := DiffDirList(testDirSubKey, []string{`C:\a`}, []string{`C:\b`})
	s.Require().NotNil(d)
	assert.Equal(s.T(), testDirSubKey, d.SubKey)
}

func (s *DiffDirListSuite) TestOldValueIsCurrentCount() {
	d := DiffDirList(testDirSubKey, []string{`C:\a`, `C:\b`, `C:\c`}, []string{`C:\d`})
	s.Require().NotNil(d)
	assert.Equal(s.T(), 3, d.OldValue)
}
