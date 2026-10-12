package cmd

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx/internal/studio"
)

// ---------------------------------------------------------------------------
// studioBackupOutput
// ---------------------------------------------------------------------------

type StudioBackupOutputSuite struct{ suite.Suite }

func TestStudioBackupOutputSuite(t *testing.T) { suite.Run(t, new(StudioBackupOutputSuite)) }

// minimalPrefs returns a Preferences value built from an empty RegistrySettings
// so we have a realistic structure without needing Windows registry access.
func (s *StudioBackupOutputSuite) minimalPrefs() *studio.Preferences {
	return studio.BuildPreferences(&studio.RegistrySettings{})
}

func (s *StudioBackupOutputSuite) marshal() []byte {
	b, err := xml.MarshalIndent(s.minimalPrefs(), "", "    ")
	s.Require().NoError(err)
	return b
}

func (s *StudioBackupOutputSuite) fullContent() []byte {
	b := s.marshal()
	content := append([]byte(xml.Header), b...)
	return append(content, '\n')
}

func (s *StudioBackupOutputSuite) TestTopLevelChildElementsIndentedFourSpaces() {
	lines := strings.Split(string(s.marshal()), "\n")

	// The first line is <Preferences>. Find the first indented child element.
	var childLine string
	for _, l := range lines[1:] {
		if strings.HasPrefix(l, " ") {
			childLine = l
			break
		}
	}

	s.Require().NotEmpty(childLine, "expected at least one indented child element in output")

	leadingSpaces := len(childLine) - len(strings.TrimLeft(childLine, " "))
	assert.Equal(s.T(), 4, leadingSpaces,
		"top-level child elements must be indented by exactly 4 spaces")
}

func (s *StudioBackupOutputSuite) TestSecondLevelChildElementsIndentedEightSpaces() {
	lines := strings.Split(string(s.marshal()), "\n")

	// Second-level children (e.g. fields inside <EditorSettings>) must have
	// exactly 8 leading spaces.
	var found bool
	for _, l := range lines {
		trimmed := strings.TrimLeft(l, " ")
		if len(l)-len(trimmed) == 8 && strings.HasPrefix(trimmed, "<") {
			found = true
			break
		}
	}

	assert.True(s.T(), found,
		"expected second-level child elements to be indented by exactly 8 spaces")
}

func (s *StudioBackupOutputSuite) TestNoTabIndentation() {
	assert.NotContains(s.T(), string(s.marshal()), "\t",
		"output must use spaces, not tabs, for indentation")
}

func (s *StudioBackupOutputSuite) TestFullContentStartsWithXMLHeader() {
	assert.True(s.T(), strings.HasPrefix(string(s.fullContent()), xml.Header),
		"backup file must begin with the XML declaration header")
}

func (s *StudioBackupOutputSuite) TestFullContentEndsWithNewline() {
	assert.True(s.T(), strings.HasSuffix(string(s.fullContent()), "\n"),
		"backup file must end with a trailing newline")
}

// ---------------------------------------------------------------------------
// defaultBackupFilename
// ---------------------------------------------------------------------------

type DefaultBackupFilenameSuite struct{ suite.Suite }

func TestDefaultBackupFilenameSuite(t *testing.T) {
	suite.Run(t, new(DefaultBackupFilenameSuite))
}

func (s *DefaultBackupFilenameSuite) TestHasEPXExtension() {
	name := time.Now().Format("netlinx-studio-backup-2006-01-02-150405") + ".epx"
	assert.True(s.T(), strings.HasSuffix(name, ".epx"),
		"default filename must end with .epx")
}

func (s *DefaultBackupFilenameSuite) TestHasPrefix() {
	name := time.Now().Format("netlinx-studio-backup-2006-01-02-150405") + ".epx"
	assert.True(s.T(), strings.HasPrefix(name, "netlinx-studio-backup-"),
		"default filename must start with 'netlinx-studio-backup-'")
}

func (s *DefaultBackupFilenameSuite) TestContainsDateAndTime() {
	now := time.Now()
	name := now.Format("netlinx-studio-backup-2006-01-02-150405") + ".epx"
	// Date portion e.g. "2026-03-16"
	assert.Contains(s.T(), name, now.Format("2006-01-02"),
		"default filename must contain the current date")
	// Time portion e.g. "150405" (HHmmss)
	assert.Contains(s.T(), name, now.Format("150405"),
		"default filename must contain the current time to second precision")
}
