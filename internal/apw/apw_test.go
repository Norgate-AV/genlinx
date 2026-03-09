package apw

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// minimalAPW returns valid APW XML with a single project/system containing
// one Include, one Module, and one MasterSrc file.
func minimalAPW() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE Workspace [
    <!ELEMENT Workspace (Identifier, CreateVersion, Project*)>
    <!ATTLIST Workspace CurrentVersion CDATA #REQUIRED>
    <!ELEMENT Identifier (#PCDATA)>
    <!ELEMENT CreateVersion (#PCDATA)>
    <!ELEMENT Project (Identifier, System*)>
    <!ELEMENT System (Identifier, SysID, File*)>
    <!ATTLIST System IsActive CDATA #REQUIRED Platform CDATA #REQUIRED Transport CDATA #REQUIRED TransportEx CDATA #REQUIRED>
    <!ELEMENT SysID (#PCDATA)>
    <!ELEMENT File (Identifier, FilePathName, Comments?)>
    <!ATTLIST File CompileType CDATA #REQUIRED Type CDATA #REQUIRED>
    <!ELEMENT FilePathName (#PCDATA)>
    <!ELEMENT Comments (#PCDATA)>
]>
<Workspace CurrentVersion="4.0">
    <Identifier>TestWorkspace</Identifier>
    <CreateVersion>4.0</CreateVersion>
    <Project>
        <Identifier>TestProject</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>TestSystem</Identifier>
            <SysID>1</SysID>
            <File CompileType="Netlinx" Type="Include">
                <Identifier>TestInclude</Identifier>
                <FilePathName>Include\TestInclude.axi</FilePathName>
                <Comments></Comments>
            </File>
            <File CompileType="Netlinx" Type="Module">
                <Identifier>TestModule</Identifier>
                <FilePathName>Module\TestModule.axs</FilePathName>
                <Comments></Comments>
            </File>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>TestMain</Identifier>
                <FilePathName>Source\TestMain.axs</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
}

// writeAPW writes APW data to a temp directory and returns (dir, apwPath).
func writeAPW(t *testing.T, name string, data []byte) (string, string) {
	t.Helper()
	dir := t.TempDir()
	apwPath := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(apwPath, data, 0o644))
	return dir, apwPath
}

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

type APWTestSuite struct {
	suite.Suite
}

func TestAPWSuite(t *testing.T) {
	suite.Run(t, new(APWTestSuite))
}

// ---------------------------------------------------------------------------
// Parse – error paths
// ---------------------------------------------------------------------------

func (s *APWTestSuite) TestParse_WrongExtension() {
	_, err := Parse("workspace.txt", minimalAPW())
	s.Error(err)
	s.Contains(err.Error(), "not a NetLinx Workspace file")
}

func (s *APWTestSuite) TestParse_MissingDOCTYPE() {
	data := []byte(`<?xml version="1.0"?><Workspace CurrentVersion="4.0"><Identifier>X</Identifier><CreateVersion>4.0</CreateVersion></Workspace>`)
	_, apwPath := writeAPW(s.T(), "workspace.apw", data)
	_, err := Parse(apwPath, data)
	s.Error(err)
	s.Contains(err.Error(), "not a NetLinx Workspace file")
}

func (s *APWTestSuite) TestParse_MalformedDTD() {
	data := []byte(`<?xml version="1.0"?><!DOCTYPE Workspace [<Workspace/>`)
	_, apwPath := writeAPW(s.T(), "workspace.apw", data)
	_, err := Parse(apwPath, data)
	s.Error(err)
	s.Contains(err.Error(), "malformed APW file")
}

func (s *APWTestSuite) TestParse_MalformedXML() {
	data := []byte(`<?xml version="1.0"?><!DOCTYPE Workspace []><NOT_VALID`)
	_, apwPath := writeAPW(s.T(), "workspace.apw", data)
	_, err := Parse(apwPath, data)
	s.Error(err)
}

// ---------------------------------------------------------------------------
// Parse – happy path
// ---------------------------------------------------------------------------

func (s *APWTestSuite) TestParse_Valid() {
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)
	s.Require().NotNil(a)
}

func (s *APWTestSuite) TestParse_ID() {
	_, apwPath := writeAPW(s.T(), "MyProject.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)
	s.Equal("MyProject", a.ID())
}

func (s *APWTestSuite) TestParse_FilePath_IsAbsolute() {
	_, apwPath := writeAPW(s.T(), "MyProject.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)
	s.True(filepath.IsAbs(a.FilePath()))
}

func (s *APWTestSuite) TestParse_CaseInsensitiveExtension() {
	_, apwPath := writeAPW(s.T(), "workspace.APW", minimalAPW())
	_, err := Parse(apwPath, minimalAPW())
	s.NoError(err)
}

// ---------------------------------------------------------------------------
// AllFiles
// ---------------------------------------------------------------------------

func (s *APWTestSuite) TestAllFiles_IncludesWorkspaceFile() {
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	files := a.AllFiles()

	var found bool
	for _, f := range files {
		if f.Type == FileTypeWorkspace {
			found = true
			break
		}
	}

	s.True(found, "expected workspace file entry in AllFiles()")
}

func (s *APWTestSuite) TestAllFiles_SortedByPath() {
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	files := a.AllFiles()
	for i := 1; i < len(files); i++ {
		s.LessOrEqual(files[i-1].Path, files[i].Path, "AllFiles() should be sorted by path")
	}
}

// ---------------------------------------------------------------------------
// File categorisation methods
// ---------------------------------------------------------------------------

func (s *APWTestSuite) TestModuleFiles() {
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	modules := a.ModuleFiles()
	s.Len(modules, 1)
	s.Equal(FileTypeModule, modules[0].Type)
	s.True(strings.HasSuffix(modules[0].Path, "TestModule.axs"))
}

func (s *APWTestSuite) TestMasterSrcFiles() {
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	masters := a.MasterSrcFiles()
	s.Len(masters, 1)
	s.Equal(FileTypeMasterSrc, masters[0].Type)
	s.True(strings.HasSuffix(masters[0].Path, "TestMain.axs"))
}

func (s *APWTestSuite) TestIncludePath_ReturnsIncludeDir() {
	dir, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	paths := a.IncludePath()
	s.Len(paths, 1)
	s.Equal(filepath.Join(dir, "Include"), paths[0])
}

func (s *APWTestSuite) TestModulePath_ReturnsModuleDir() {
	dir, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	paths := a.ModulePath()
	s.Len(paths, 1)
	s.Equal(filepath.Join(dir, "Module"), paths[0])
}

// ---------------------------------------------------------------------------
// FileIsReadable (table-driven)
// ---------------------------------------------------------------------------

func TestFileIsReadable(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"main.axs", true},
		{"include.axi", true},
		{"main.AXS", true}, // case-insensitive
		{"module.tko", false},
		{"workspace.apw", false},
		{"touch.tp4", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			assert.Equal(t, tt.expected, FileIsReadable(tt.file))
		})
	}
}

// ---------------------------------------------------------------------------
// FileIsOfInterest (table-driven)
// ---------------------------------------------------------------------------

func TestFileIsOfInterest(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"main.axs", true},
		{"include.axi", true},
		{"duet.jar", true},
		{"device.xdd", true},
		{"main.AXS", true},
		{"module.tko", false},
		{"workspace.apw", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			assert.Equal(t, tt.expected, FileIsOfInterest(tt.file))
		})
	}
}

// ---------------------------------------------------------------------------
// GetFileType (table-driven)
// ---------------------------------------------------------------------------

func TestGetFileType(t *testing.T) {
	tests := []struct {
		file     string
		expected FileType
	}{
		{"workspace.apw", FileTypeWorkspace},
		{"module.tko", FileTypeTKO},
		{"unknown.xyz", FileTypeOther},
		{"", FileTypeOther},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetFileType(tt.file))
		})
	}
}

// .axs maps to FileTypeSource, FileTypeMasterSrc, and FileTypeModule (all share
// the .axs extension). Map iteration order is non-deterministic, so any of the
// three is a valid return value — verify it's one of the known AXS types.
func TestGetFileType_AXS_IsKnownAXSType(t *testing.T) {
	result := GetFileType("main.axs")
	known := result == FileTypeSource || result == FileTypeMasterSrc || result == FileTypeModule
	assert.True(t, known,
		"expected FileTypeSource, FileTypeMasterSrc, or FileTypeModule for .axs, got %q", result)
}

// ---------------------------------------------------------------------------
// MasterSrcPath
// ---------------------------------------------------------------------------

func (s *APWTestSuite) TestMasterSrcPath_ReturnsMasterSrcDir() {
	dir, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	paths := a.MasterSrcPath()
	s.Len(paths, 1)
	s.Equal(filepath.Join(dir, "Source"), paths[0])
}

func (s *APWTestSuite) TestMasterSrcPath_EmptyWhenNoMasterSrc() {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE Workspace [
    <!ELEMENT Workspace (Identifier, CreateVersion, Project*)>
    <!ATTLIST Workspace CurrentVersion CDATA #REQUIRED>
    <!ELEMENT Identifier (#PCDATA)>
    <!ELEMENT CreateVersion (#PCDATA)>
    <!ELEMENT Project (Identifier, System*)>
    <!ELEMENT System (Identifier, SysID, File*)>
    <!ATTLIST System IsActive CDATA #REQUIRED Platform CDATA #REQUIRED Transport CDATA #REQUIRED TransportEx CDATA #REQUIRED>
    <!ELEMENT SysID (#PCDATA)>
    <!ELEMENT File (Identifier, FilePathName, Comments?)>
    <!ATTLIST File CompileType CDATA #REQUIRED Type CDATA #REQUIRED>
    <!ELEMENT FilePathName (#PCDATA)>
    <!ELEMENT Comments (#PCDATA)>
]>
<Workspace CurrentVersion="4.0">
    <Identifier>NoMaster</Identifier>
    <CreateVersion>4.0</CreateVersion>
    <Project>
        <Identifier>P</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>S</Identifier>
            <SysID>1</SysID>
            <File CompileType="Netlinx" Type="Include">
                <Identifier>Inc</Identifier>
                <FilePathName>Include\Inc.axi</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
	_, apwPath := writeAPW(s.T(), "NoMaster.apw", data)
	a, err := Parse(apwPath, data)
	s.Require().NoError(err)

	paths := a.MasterSrcPath()
	s.Empty(paths)
}

// ---------------------------------------------------------------------------
// GetExtraFileReferencesFromFile
// ---------------------------------------------------------------------------

func TestGetExtraFileReferencesFromFile_IncludeDirective(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "main.axs")
	require.NoError(t, os.WriteFile(srcPath, []byte("#include 'SomeLibrary'\n"), 0o644))

	// APW with no existing files — nothing is "in workspace" yet.
	_, apwPath := writeAPW(t, "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	require.NoError(t, err)

	refs, err := a.GetExtraFileReferencesFromFile(srcPath)
	require.NoError(t, err)
	assert.Contains(t, refs, "SomeLibrary")
}

func TestGetExtraFileReferencesFromFile_DefineModuleDirective(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "main.axs")
	require.NoError(t, os.WriteFile(srcPath, []byte("define_module 'MyModule' md()\n"), 0o644))

	_, apwPath := writeAPW(t, "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	require.NoError(t, err)

	refs, err := a.GetExtraFileReferencesFromFile(srcPath)
	require.NoError(t, err)
	assert.Contains(t, refs, "MyModule")
}

func TestGetExtraFileReferencesFromFile_ExcludesWorkspaceEntries(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "main.axs")
	// Reference "TestInclude" which is already a file ID in minimalAPW().
	require.NoError(t, os.WriteFile(srcPath, []byte("#include 'TestInclude'\n"), 0o644))

	_, apwPath := writeAPW(t, "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	require.NoError(t, err)

	refs, err := a.GetExtraFileReferencesFromFile(srcPath)
	require.NoError(t, err)
	assert.NotContains(t, refs, "TestInclude", "IDs already in workspace must not be returned as extra refs")
}

func TestGetExtraFileReferencesFromFile_NilForNonReadableFile(t *testing.T) {
	_, apwPath := writeAPW(t, "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	require.NoError(t, err)

	// .tko is not readable per FileIsReadable().
	refs, err := a.GetExtraFileReferencesFromFile("firmware.tko")
	require.NoError(t, err)
	assert.Nil(t, refs)
}

func TestGetExtraFileReferencesFromFile_NoDuplicates(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "main.axs")
	content := "#include 'LibA'\n#include 'LibA'\n"
	require.NoError(t, os.WriteFile(srcPath, []byte(content), 0o644))

	_, apwPath := writeAPW(t, "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	require.NoError(t, err)

	refs, err := a.GetExtraFileReferencesFromFile(srcPath)
	require.NoError(t, err)

	count := 0
	for _, r := range refs {
		if r == "LibA" {
			count++
		}
	}

	assert.Equal(t, 1, count, "duplicate references must be deduplicated")
}

// ---------------------------------------------------------------------------
// GetExtraFileReferences (workspace-wide)
// ---------------------------------------------------------------------------

func (s *APWTestSuite) TestGetExtraFileReferences_EmptyWhenNoExtraRefs() {
	// minimalAPW workspace files don't exist on disk → skipped by GetExtraFileReferences.
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)

	refs, err := a.GetExtraFileReferences()
	s.Require().NoError(err)
	s.Empty(refs)
}
