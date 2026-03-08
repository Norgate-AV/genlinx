package archive

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx-go/internal/apw"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func minimalAPW(id string) []byte {
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
    <Identifier>` + id + `</Identifier>
    <CreateVersion>4.0</CreateVersion>
    <Project>
        <Identifier>TestProject</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>TestSystem</Identifier>
            <SysID>1</SysID>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>TestMain</Identifier>
                <FilePathName>Source\TestMain.axs</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
}

// setupWorkspace writes a minimal APW and its referenced source file to a
// temp directory, changes the working directory to that temp dir (restored
// after the test), and returns the parsed APW.
func setupWorkspace(t *testing.T, id string) *apw.APW {
	t.Helper()
	dir := t.TempDir()

	data := minimalAPW(id)
	apwPath := filepath.Join(dir, id+".apw")
	require.NoError(t, os.WriteFile(apwPath, data, 0o644))

	// Create the source file so it's flagged as Exists=true in the APW.
	sourceDir := filepath.Join(dir, "Source")
	require.NoError(t, os.MkdirAll(sourceDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "TestMain.axs"), []byte(`PROGRAM_NAME='test'`), 0o644))

	a, err := apw.Parse(apwPath, data)
	require.NoError(t, err)

	// Build() writes the zip to the current working directory.
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(dir))

	return a
}

func defaultOpts() *Options {
	return &Options{
		OutputFileSuffix:           "zip",
		IncludeCompiledSourceFiles: false,
		IncludeCompiledModuleFiles: false,
		IncludeFilesNotInWorkspace: false,
		ExtraFileArchiveLocation:   "Extra",
	}
}

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

type ArchiveTestSuite struct {
	suite.Suite
}

func TestArchiveSuite(t *testing.T) {
	suite.Run(t, new(ArchiveTestSuite))
}

// ---------------------------------------------------------------------------
// NewBuilder
// ---------------------------------------------------------------------------

func (s *ArchiveTestSuite) TestNewBuilder_NotNil() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	b := NewBuilder(a, defaultOpts())
	s.Require().NotNil(b)
}

// ---------------------------------------------------------------------------
// Build – file creation
// ---------------------------------------------------------------------------

func (s *ArchiveTestSuite) TestBuild_CreatesZipFile() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zipPath := filepath.Join(wd, "TestWorkspace.zip")
	_, statErr := os.Stat(zipPath)
	s.NoError(statErr, "zip file should exist after Build()")
}

func (s *ArchiveTestSuite) TestBuild_OutputFileNamed_IDdotSuffix() {
	a := setupWorkspace(s.T(), "MyProject")
	opts := defaultOpts()
	opts.OutputFileSuffix = "archive.zip"
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	_, statErr := os.Stat(filepath.Join(wd, "MyProject.archive.zip"))
	s.NoError(statErr, "output file should follow <id>.<suffix> naming")
}

// ---------------------------------------------------------------------------
// Build – zip contents
// ---------------------------------------------------------------------------

func (s *ArchiveTestSuite) TestBuild_ZipContains_WorkspaceFile() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zipPath := filepath.Join(wd, "TestWorkspace.zip")
	zr, err := zip.OpenReader(zipPath)
	s.Require().NoError(err)
	defer zr.Close()

	var found bool
	for _, f := range zr.File {
		if f.Name == "TestWorkspace.apw" {
			found = true
			break
		}
	}
	s.True(found, "zip should contain the workspace .apw file")
}

func (s *ArchiveTestSuite) TestBuild_ZipContains_SourceFile() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zipPath := filepath.Join(wd, "TestWorkspace.zip")
	zr, err := zip.OpenReader(zipPath)
	s.Require().NoError(err)
	defer zr.Close()

	// archiveDir returns the full absolute directory of the file, so the zip
	// entry is the full forward-slash path — just check the filename suffix.
	var found bool
	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, "TestMain.axs") {
			found = true
			break
		}
	}
	s.True(found, "zip should contain an entry for TestMain.axs")
}

func (s *ArchiveTestSuite) TestBuild_ValidZip() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	data, err := os.ReadFile(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	s.True(bytes.HasPrefix(data, []byte("PK")), "output should be a valid zip file (PK magic bytes)")
}

// ---------------------------------------------------------------------------
// zipEntryPath
// ---------------------------------------------------------------------------

func TestZipEntryPath_ForwardSlashes(t *testing.T) {
	result := zipEntryPath(`Source\SubDir\file.axs`)
	if result != "Source/SubDir/file.axs" {
		t.Errorf("expected forward slashes, got %q", result)
	}
}

func TestZipEntryPath_AlreadyForwardSlashes(t *testing.T) {
	result := zipEntryPath("Source/SubDir/file.axs")
	if result != "Source/SubDir/file.axs" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestZipEntryPath_CleansDots(t *testing.T) {
	result := zipEntryPath("Source/./SubDir/../file.axs")
	if result != "Source/file.axs" {
		t.Errorf("expected cleaned path, got %q", result)
	}
}
