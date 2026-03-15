package studio

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// parseTCPIPEntry
// ---------------------------------------------------------------------------

type ParseTCPIPEntrySuite struct{ suite.Suite }

func TestParseTCPIPEntrySuite(t *testing.T) { suite.Run(t, new(ParseTCPIPEntrySuite)) }

func (s *ParseTCPIPEntrySuite) TestFullSixFieldEntry() {
	e := parseTCPIPEntry("192.168.1.1|1319|1|MyRoom|user|pass")
	assert.Equal(s.T(), "192.168.1.1", e.Host)
	assert.Equal(s.T(), 1319, e.Port)
	assert.True(s.T(), e.PingTest)
	assert.Equal(s.T(), "MyRoom", e.Name)
	assert.Equal(s.T(), "user", e.Username)
	assert.Equal(s.T(), "pass", e.Password)
}

func (s *ParseTCPIPEntrySuite) TestPingFalse() {
	e := parseTCPIPEntry("10.0.0.1|1319|0|Lab||")
	assert.False(s.T(), e.PingTest)
}

func (s *ParseTCPIPEntrySuite) TestPingNonZeroIsTrue() {
	e := parseTCPIPEntry("10.0.0.1|1319|2|Lab||")
	assert.True(s.T(), e.PingTest)
}

func (s *ParseTCPIPEntrySuite) TestShortInputDoesNotPanic() {
	// Fewer than 6 pipe-separated fields — missing fields should default to "".
	e := parseTCPIPEntry("10.0.0.1|1319|1")
	assert.Equal(s.T(), "10.0.0.1", e.Host)
	assert.Equal(s.T(), 1319, e.Port)
	assert.True(s.T(), e.PingTest)
	assert.Equal(s.T(), "", e.Name)
	assert.Equal(s.T(), "", e.Username)
	assert.Equal(s.T(), "", e.Password)
}

func (s *ParseTCPIPEntrySuite) TestEmptyStringDoesNotPanic() {
	e := parseTCPIPEntry("")
	assert.Equal(s.T(), "", e.Host)
	assert.Equal(s.T(), 0, e.Port)
	assert.False(s.T(), e.PingTest)
}

func (s *ParseTCPIPEntrySuite) TestInvalidPortDefaultsToZero() {
	e := parseTCPIPEntry("10.0.0.1|notaport|1|Room||")
	assert.Equal(s.T(), 0, e.Port)
}

func (s *ParseTCPIPEntrySuite) TestEmptyCredentials() {
	e := parseTCPIPEntry("10.0.0.1|1319|1|Room||")
	assert.Equal(s.T(), "", e.Username)
	assert.Equal(s.T(), "", e.Password)
}

// ---------------------------------------------------------------------------
// TCPIPEntry.encode
// ---------------------------------------------------------------------------

type TCPIPEntryEncodeSuite struct{ suite.Suite }

func TestTCPIPEntryEncodeSuite(t *testing.T) { suite.Run(t, new(TCPIPEntryEncodeSuite)) }

func (s *TCPIPEntryEncodeSuite) TestRoundTrip() {
	raw := "192.168.1.1|1319|1|Lab|user|pass"
	assert.Equal(s.T(), raw, parseTCPIPEntry(raw).encode())
}

func (s *TCPIPEntryEncodeSuite) TestRoundTripEmptyCredentials() {
	raw := "10.0.0.1|1319|0|Room||"
	assert.Equal(s.T(), raw, parseTCPIPEntry(raw).encode())
}

func (s *TCPIPEntryEncodeSuite) TestPingFalseEncodesAsZero() {
	e := TCPIPEntry{Host: "10.0.0.1", Port: 1319, PingTest: false}
	assert.Equal(s.T(), "10.0.0.1|1319|0|||", e.encode())
}

func (s *TCPIPEntryEncodeSuite) TestPingTrueEncodesAsOne() {
	e := TCPIPEntry{Host: "10.0.0.1", Port: 1319, PingTest: true, Name: "Room"}
	assert.Equal(s.T(), "10.0.0.1|1319|1|Room||", e.encode())
}

// ---------------------------------------------------------------------------
// TCPIPHistory — MarshalXML / UnmarshalXML
// ---------------------------------------------------------------------------

type TCPIPHistoryXMLSuite struct{ suite.Suite }

func TestTCPIPHistoryXMLSuite(t *testing.T) { suite.Run(t, new(TCPIPHistoryXMLSuite)) }

func (s *TCPIPHistoryXMLSuite) TestMarshalUsesIndexedElementNames() {
	h := TCPIPHistory{
		Entries: []TCPIPEntry{
			{Host: "192.168.1.1", Port: 1319, PingTest: true, Name: "A"},
			{Host: "10.0.0.1", Port: 1319, PingTest: false, Name: "B"},
		},
	}
	data, err := xml.Marshal(h)
	s.Require().NoError(err)
	xmlStr := string(data)
	assert.Contains(s.T(), xmlStr, "TCPIPHistoryEX0")
	assert.Contains(s.T(), xmlStr, "TCPIPHistoryEX1")
	assert.Contains(s.T(), xmlStr, "192.168.1.1|1319|1|A||")
	assert.Contains(s.T(), xmlStr, "10.0.0.1|1319|0|B||")
}

func (s *TCPIPHistoryXMLSuite) TestMarshalEmptyHistory() {
	data, err := xml.Marshal(TCPIPHistory{})
	s.Require().NoError(err)
	assert.NotContains(s.T(), string(data), "TCPIPHistoryEX")
}

func (s *TCPIPHistoryXMLSuite) TestUnmarshalRoundTrip() {
	original := TCPIPHistory{
		Entries: []TCPIPEntry{
			{Host: "192.168.1.1", Port: 1319, PingTest: true, Name: "A"},
			{Host: "10.0.0.1", Port: 1319, PingTest: false, Name: "B"},
		},
	}
	data, err := xml.Marshal(original)
	s.Require().NoError(err)

	var decoded TCPIPHistory
	s.Require().NoError(xml.Unmarshal(data, &decoded))
	s.Require().Len(decoded.Entries, 2)
	assert.Equal(s.T(), original.Entries[0].Host, decoded.Entries[0].Host)
	assert.Equal(s.T(), original.Entries[0].PingTest, decoded.Entries[0].PingTest)
	assert.Equal(s.T(), original.Entries[1].Host, decoded.Entries[1].Host)
	assert.Equal(s.T(), original.Entries[1].PingTest, decoded.Entries[1].PingTest)
}

func (s *TCPIPHistoryXMLSuite) TestUnmarshalPreservesInsertionOrder() {
	xmlData := `<TCPIPHistory>` +
		`<TCPIPHistoryEX0>192.168.1.1|1319|1|First||</TCPIPHistoryEX0>` +
		`<TCPIPHistoryEX1>10.0.0.1|1319|0|Second||</TCPIPHistoryEX1>` +
		`<TCPIPHistoryEX2>172.16.0.1|1319|1|Third||</TCPIPHistoryEX2>` +
		`</TCPIPHistory>`

	var h TCPIPHistory
	s.Require().NoError(xml.Unmarshal([]byte(xmlData), &h))
	s.Require().Len(h.Entries, 3)
	assert.Equal(s.T(), "First", h.Entries[0].Name)
	assert.Equal(s.T(), "Second", h.Entries[1].Name)
	assert.Equal(s.T(), "Third", h.Entries[2].Name)
}

func (s *TCPIPHistoryXMLSuite) TestUnmarshalEmptyBlock() {
	var h TCPIPHistory
	s.Require().NoError(xml.Unmarshal([]byte(`<TCPIPHistory></TCPIPHistory>`), &h))
	assert.Empty(s.T(), h.Entries)
}

// ---------------------------------------------------------------------------
// NetlinxCompilerSettings — MarshalXML / UnmarshalXML
// ---------------------------------------------------------------------------

type NetlinxCompilerSettingsXMLSuite struct{ suite.Suite }

func TestNetlinxCompilerSettingsXMLSuite(t *testing.T) {
	suite.Run(t, new(NetlinxCompilerSettingsXMLSuite))
}

func (s *NetlinxCompilerSettingsXMLSuite) TestMarshalProducesZeroPaddedDirElements() {
	cs := NetlinxCompilerSettings{
		BuildWithSource: 1,
		LibraryDirs:     []string{`C:\lib\a`, `C:\lib\b`},
		IncludeDirs:     []string{`C:\inc\a`},
		ModuleDirs:      []string{},
	}
	data, err := xml.Marshal(cs)
	s.Require().NoError(err)
	xmlStr := string(data)
	assert.Contains(s.T(), xmlStr, "NLXLibraryDir000")
	assert.Contains(s.T(), xmlStr, "NLXLibraryDir001")
	assert.Contains(s.T(), xmlStr, "NLXIncludeDir000")
	assert.NotContains(s.T(), xmlStr, "NLXModuleDir000")
}

func (s *NetlinxCompilerSettingsXMLSuite) TestMarshalFixedFieldsPresent() {
	cs := NetlinxCompilerSettings{
		BuildWithSource:               1,
		BuildWithDebugInfo:            0,
		PasswordProtect:               0,
		EnableWCPreprocessor:          1,
		Password:                      "secret",
		ShowDebugWindowOnSessionClose: 0,
		ShowMainAXSOnSessionStart:     1,
	}
	data, err := xml.Marshal(cs)
	s.Require().NoError(err)
	xmlStr := string(data)
	assert.Contains(s.T(), xmlStr, "BuildWithSource")
	assert.Contains(s.T(), xmlStr, "BuildWithDebugInfo")
	assert.Contains(s.T(), xmlStr, "EnableWCPreprocessor")
	assert.Contains(s.T(), xmlStr, "Password")
	assert.Contains(s.T(), xmlStr, "secret")
	assert.Contains(s.T(), xmlStr, "ShowMainAXSOnSessionStart")
}

func (s *NetlinxCompilerSettingsXMLSuite) TestRoundTrip() {
	original := NetlinxCompilerSettings{
		BuildWithSource:               1,
		BuildWithDebugInfo:            0,
		PasswordProtect:               0,
		EnableWCPreprocessor:          1,
		Password:                      "s3cr3t",
		ShowDebugWindowOnSessionClose: 0,
		ShowMainAXSOnSessionStart:     1,
		LibraryDirs:                   []string{`C:\AMX\lib`},
		IncludeDirs:                   []string{`C:\AMX\inc`},
		ModuleDirs:                    []string{`C:\AMX\mod`},
	}
	data, err := xml.Marshal(original)
	s.Require().NoError(err)

	var decoded NetlinxCompilerSettings
	s.Require().NoError(xml.Unmarshal(data, &decoded))

	assert.Equal(s.T(), original.BuildWithSource, decoded.BuildWithSource)
	assert.Equal(s.T(), original.BuildWithDebugInfo, decoded.BuildWithDebugInfo)
	assert.Equal(s.T(), original.PasswordProtect, decoded.PasswordProtect)
	assert.Equal(s.T(), original.EnableWCPreprocessor, decoded.EnableWCPreprocessor)
	assert.Equal(s.T(), original.Password, decoded.Password)
	assert.Equal(s.T(), original.ShowDebugWindowOnSessionClose, decoded.ShowDebugWindowOnSessionClose)
	assert.Equal(s.T(), original.ShowMainAXSOnSessionStart, decoded.ShowMainAXSOnSessionStart)
	assert.Equal(s.T(), original.LibraryDirs, decoded.LibraryDirs)
	assert.Equal(s.T(), original.IncludeDirs, decoded.IncludeDirs)
	assert.Equal(s.T(), original.ModuleDirs, decoded.ModuleDirs)
}

func (s *NetlinxCompilerSettingsXMLSuite) TestRoundTripNoDirs() {
	original := NetlinxCompilerSettings{BuildWithSource: 1}
	data, err := xml.Marshal(original)
	s.Require().NoError(err)

	var decoded NetlinxCompilerSettings
	s.Require().NoError(xml.Unmarshal(data, &decoded))
	assert.Equal(s.T(), 1, decoded.BuildWithSource)
	assert.Empty(s.T(), decoded.LibraryDirs)
	assert.Empty(s.T(), decoded.IncludeDirs)
	assert.Empty(s.T(), decoded.ModuleDirs)
}

func (s *NetlinxCompilerSettingsXMLSuite) TestDirIndexesAreSequential() {
	cs := NetlinxCompilerSettings{
		LibraryDirs: []string{`C:\a`, `C:\b`, `C:\c`},
	}
	data, err := xml.Marshal(cs)
	s.Require().NoError(err)
	xmlStr := string(data)
	assert.Contains(s.T(), xmlStr, "NLXLibraryDir000")
	assert.Contains(s.T(), xmlStr, "NLXLibraryDir001")
	assert.Contains(s.T(), xmlStr, "NLXLibraryDir002")
	assert.NotContains(s.T(), xmlStr, "NLXLibraryDir003")
}
