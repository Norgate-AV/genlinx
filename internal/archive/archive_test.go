package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx/internal/apw"
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
            <File CompileType="Netlinx" Type="Module">
                <Identifier>TestModule</Identifier>
                <FilePathName>Module\TestModule.axs</FilePathName>
                <Comments></Comments>
            </File>
            <File CompileType="Netlinx" Type="Include">
                <Identifier>TestInclude</Identifier>
                <FilePathName>Include\TestInclude.axi</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
}

// setupWorkspace writes a minimal APW and all its referenced files (Source,
// Module, Include) to a temp directory, changes the working directory to that
// temp dir (restored after the test), and returns the parsed APW.
func setupWorkspace(t *testing.T, id string) *apw.APW {
	t.Helper()
	dir := t.TempDir()

	data := minimalAPW(id)
	apwPath := filepath.Join(dir, id+".apw")
	require.NoError(t, os.WriteFile(apwPath, data, 0o644))

	// Create all referenced files so they are flagged as Exists=true in the APW.
	for _, sub := range []struct{ dir, name string }{
		{"Source", "TestMain.axs"},
		{"Module", "TestModule.axs"},
		{"Include", "TestInclude.axi"},
	} {
		subDir := filepath.Join(dir, sub.dir)
		require.NoError(t, os.MkdirAll(subDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, sub.name), []byte(`PROGRAM_NAME='test'`), 0o644))
	}

	a, err := apw.Parse(apwPath, data)
	require.NoError(t, err)

	// Build() writes the zip to the current working directory.
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(dir))

	return a
}

// minimalAPWFull is like minimalAPW but with configurable project and system identifiers.
func minimalAPWFull(wsID, projectID, systemID string) []byte {
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
    <Identifier>` + wsID + `</Identifier>
    <CreateVersion>4.0</CreateVersion>
    <Project>
        <Identifier>` + projectID + `</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>` + systemID + `</Identifier>
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

// setupWorkspaceWithIDs is like setupWorkspace but with configurable project
// and system identifiers, allowing tests to use identifiers that contain spaces.
func setupWorkspaceWithIDs(t *testing.T, wsID, projectID, systemID string) *apw.APW {
	t.Helper()
	dir := t.TempDir()

	data := minimalAPWFull(wsID, projectID, systemID)
	apwPath := filepath.Join(dir, wsID+".apw")
	require.NoError(t, os.WriteFile(apwPath, data, 0o644))

	subDir := filepath.Join(dir, "Source")
	require.NoError(t, os.MkdirAll(subDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "TestMain.axs"), []byte(`PROGRAM_NAME='test'`), 0o644))

	a, err := apw.Parse(apwPath, data)
	require.NoError(t, err)

	oldWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(dir))

	return a
}

func defaultOpts() *Options {
	return &Options{
		OutputFileSuffix:           "zip",
		IncludeCompiledSourceFiles: false,
		IncludeCompiledModuleFiles: false,
		IncludeFilesNotInWorkspace: false,
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
	defer func() { _ = zr.Close() }()

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
	defer func() { _ = zr.Close() }()

	var found bool
	for _, f := range zr.File {
		if f.Name == "TestMain.axs" {
			found = true
			break
		}
	}

	s.True(found, "zip should contain TestMain.axs at the archive root")
}

func (s *ArchiveTestSuite) TestBuild_ZipEntry_ModuleUsesRelativePath() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	var found bool
	for _, f := range zr.File {
		if f.Name == "TestModule.axs" {
			found = true
			break
		}
	}

	s.True(found, "zip should contain TestModule.axs at the archive root")
}

func (s *ArchiveTestSuite) TestBuild_ZipEntry_IncludeUsesRelativePath() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	var found bool
	for _, f := range zr.File {
		if f.Name == "TestInclude.axi" {
			found = true
			break
		}
	}

	s.True(found, "zip should contain TestInclude.axi at the archive root")
}

func (s *ArchiveTestSuite) TestBuild_ZipEntry_NoAbsolutePaths() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	for _, f := range zr.File {
		// A zip entry must never start with a drive letter (C:/) or a leading slash.
		s.False(strings.HasPrefix(f.Name, "/"), "entry %q must not start with /", f.Name)
		if len(f.Name) >= 2 {
			s.NotEqual(':', rune(f.Name[1]), "entry %q must not be an absolute Windows path", f.Name)
		}
	}
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
// Build – compiled file inclusion
// ---------------------------------------------------------------------------

func (s *ArchiveTestSuite) TestBuild_CompiledModuleFile_AddedWithRelativePath() {
	a := setupWorkspace(s.T(), "TestWorkspace")

	// Create the pre-compiled .tko file next to the module source.
	wd, _ := os.Getwd()
	tkoPath := filepath.Join(wd, "Module", "TestModule.tko")
	s.Require().NoError(os.WriteFile(tkoPath, []byte("compiled-module"), 0o644))

	opts := defaultOpts()
	opts.IncludeCompiledModuleFiles = true
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	var found bool
	for _, f := range zr.File {
		if f.Name == "TestModule.tko" {
			found = true
			break
		}
	}

	s.True(found, "zip should contain TestModule.tko at the archive root")
}

func (s *ArchiveTestSuite) TestBuild_CompiledModuleFile_SkippedWhenOptionFalse() {
	a := setupWorkspace(s.T(), "TestWorkspace")

	// Even if the .tko exists on disk it should NOT appear in the zip when
	// IncludeCompiledModuleFiles is false (the default).
	wd, _ := os.Getwd()
	s.Require().NoError(os.WriteFile(
		filepath.Join(wd, "Module", "TestModule.tko"),
		[]byte("compiled-module"), 0o644,
	))

	opts := defaultOpts() // IncludeCompiledModuleFiles = false
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	for _, f := range zr.File {
		s.NotEqual("TestModule.tko", f.Name,
			"compiled .tko should not be in zip when IncludeCompiledModuleFiles=false")
	}
}

func (s *ArchiveTestSuite) TestBuild_CompiledSourceFile_AddedWithRelativePath() {
	a := setupWorkspace(s.T(), "TestWorkspace")

	// TestMain is type MasterSrc — compiled to .tkn.
	wd, _ := os.Getwd()
	tknPath := filepath.Join(wd, "Source", "TestMain.tkn")
	s.Require().NoError(os.WriteFile(tknPath, []byte("compiled-source"), 0o644))

	opts := defaultOpts()
	opts.IncludeCompiledSourceFiles = true
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	var found bool
	for _, f := range zr.File {
		if f.Name == "TestMain.tkn" {
			found = true
			break
		}
	}

	s.True(found, "zip should contain TestMain.tkn at the archive root")
}

func (s *ArchiveTestSuite) TestBuild_ExtraModuleFile_TKOAddedToArchiveRoot() {
	// Verify that an extra .axs file treated as a Module adds both itself and
	// its compiled .tko counterpart at the archive root (flat layout).
	a := setupWorkspace(s.T(), "TestWorkspace")
	wd, _ := os.Getwd()

	// Simulate an extra module file that was located on disk.
	// We exercise addFileToArchive directly by creating the apw.File manually
	// and calling Build with a builder that has a pre-populated locatedExtraRefs.
	// Instead, we test via the public Build() path using a fabricated IsExtra file
	// injected through addModuleItem, which we call directly.

	// Create a fake extra module file on disk (with its .tko counterpart).
	extraDir := filepath.Join(wd, "ExtraModules")
	s.Require().NoError(os.MkdirAll(extraDir, 0o755))
	extraAxsPath := filepath.Join(extraDir, "ExtraLib.axs")
	extraTkoPath := filepath.Join(extraDir, "ExtraLib.tko")
	s.Require().NoError(os.WriteFile(extraAxsPath, []byte("PROGRAM_NAME='extra'"), 0o644))
	s.Require().NoError(os.WriteFile(extraTkoPath, []byte("compiled-extra-module"), 0o644))

	opts := defaultOpts()
	opts.IncludeCompiledModuleFiles = true

	b := NewBuilder(a, opts)

	// Open a real zip writer so we can call addModuleItem directly.
	zipPath := filepath.Join(wd, "extra_test.zip")
	f, err := os.Create(zipPath)
	s.Require().NoError(err)
	b.zipWriter = zip.NewWriter(f)

	extraFile := apw.File{
		Type:    apw.FileTypeModule,
		Path:    extraAxsPath,
		Exists:  true,
		IsExtra: true,
	}

	s.Require().NoError(b.addModuleItem(extraFile))
	s.Require().NoError(b.zipWriter.Close())
	s.Require().NoError(f.Close())

	zr, err := zip.OpenReader(zipPath)
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	var axsFound, tkoFound bool
	for _, zf := range zr.File {
		if zf.Name == "ExtraLib.axs" {
			axsFound = true
		}
		if zf.Name == "ExtraLib.tko" {
			tkoFound = true
		}
	}

	s.True(axsFound, "extra .axs should be added at the archive root")
	s.True(tkoFound, "compiled .tko should be added at the archive root alongside extra .axs")
}

// ---------------------------------------------------------------------------
// Build – extra files (IncludeFilesNotInWorkspace)
// ---------------------------------------------------------------------------

// writeDefineModule overwrites the master source file with content that
// references ExtraLib via a define_module directive, triggering the extra-file
// discovery pipeline when IncludeFilesNotInWorkspace is true.
func writeDefineModule(t *testing.T, wd string) {
	t.Helper()
	src := filepath.Join(wd, "Source", "TestMain.axs")
	require.NoError(t, os.WriteFile(src, []byte("PROGRAM_NAME='test'\ndefine_module 'ExtraLib'\n"), 0o644))
}

// createExtraLibDir creates a search directory containing ExtraLib.axs and
// returns the directory path.
func createExtraLibDir(t *testing.T, wd string) string {
	t.Helper()
	dir := filepath.Join(wd, "ExtraLibs")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ExtraLib.axs"), []byte("PROGRAM_NAME='extra'"), 0o644))
	return dir
}

func (s *ArchiveTestSuite) TestBuild_CompiledSourceFile_SkippedWhenOptionFalse() {
	a := setupWorkspace(s.T(), "TestWorkspace")

	// Even if the .tkn exists on disk it should NOT appear in the zip when
	// IncludeCompiledSourceFiles is false (the default).
	wd, _ := os.Getwd()
	s.Require().NoError(os.WriteFile(
		filepath.Join(wd, "Source", "TestMain.tkn"),
		[]byte("compiled-source"), 0o644,
	))

	opts := defaultOpts() // IncludeCompiledSourceFiles = false
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	for _, f := range zr.File {
		s.NotEqual("TestMain.tkn", f.Name,
			"compiled .tkn should not be in zip when IncludeCompiledSourceFiles=false")
	}
}

func (s *ArchiveTestSuite) TestBuild_IncludeFilesNotInWorkspace_False_NoExtraFiles() {
	// Even when a define_module reference and a matching file on disk exist,
	// no extra files should be added when IncludeFilesNotInWorkspace is false.
	a := setupWorkspace(s.T(), "TestWorkspace")
	wd, _ := os.Getwd()

	writeDefineModule(s.T(), wd)
	extraDir := createExtraLibDir(s.T(), wd)

	opts := defaultOpts() // IncludeFilesNotInWorkspace = false
	opts.ExtraFileSearchLocations = []string{extraDir}
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	for _, f := range zr.File {
		s.NotEqual("ExtraLib.axs", f.Name,
			"no extra files should appear when IncludeFilesNotInWorkspace=false; found %q", f.Name)
	}
}

func (s *ArchiveTestSuite) TestBuild_IncludeFilesNotInWorkspace_True_AddsExtraFiles() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	wd, _ := os.Getwd()

	writeDefineModule(s.T(), wd)
	extraDir := createExtraLibDir(s.T(), wd)

	opts := defaultOpts()
	opts.IncludeFilesNotInWorkspace = true
	opts.ExtraFileSearchLocations = []string{extraDir}
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	var found bool
	for _, f := range zr.File {
		if f.Name == "ExtraLib.axs" {
			found = true
			break
		}
	}

	s.True(found, "extra file should be added at the archive root when IncludeFilesNotInWorkspace=true")
}

func (s *ArchiveTestSuite) TestBuild_IgnoredFiles_ExcludesFromExtraSearch() {
	a := setupWorkspace(s.T(), "TestWorkspace")
	wd, _ := os.Getwd()

	writeDefineModule(s.T(), wd)
	extraDir := createExtraLibDir(s.T(), wd)

	opts := defaultOpts()
	opts.IncludeFilesNotInWorkspace = true
	opts.ExtraFileSearchLocations = []string{extraDir}
	opts.IgnoredFiles = []string{"ExtraLib.axs"}
	err := NewBuilder(a, opts).Build()
	s.Require().NoError(err)

	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	for _, f := range zr.File {
		s.NotEqual("ExtraLib.axs", f.Name,
			"ignored file should not appear in zip even when IncludeFilesNotInWorkspace=true")
	}
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

// ---------------------------------------------------------------------------
// Build – verbose output (covers displayZippedFiles + logVerbose)
// ---------------------------------------------------------------------------

// captureStdout redirects os.Stdout for the duration of fn and returns the
// captured output. It is defined locally because the archive package is an
// internal package; the cmd package helper is not accessible here.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w
	fn()
	require.NoError(t, w.Close())
	os.Stdout = orig

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)

	return buf.String()
}

func (s *ArchiveTestSuite) TestBuild_Verbose_PrintsEntries() {
	a := setupWorkspace(s.T(), "TestWorkspace")

	opts := defaultOpts()
	opts.Verbose = true

	out := captureStdout(s.T(), func() {
		err := NewBuilder(a, opts).Build()
		s.Require().NoError(err)
	})

	// displayZippedFiles prints "--> <entry>" for every zip entry.
	s.Contains(out, "-->", "verbose output should list entries with --> prefix")
}

func (s *ArchiveTestSuite) TestBuild_NonVerbose_NoEntryList() {
	a := setupWorkspace(s.T(), "TestWorkspace")

	opts := defaultOpts()
	opts.Verbose = false

	out := captureStdout(s.T(), func() {
		err := NewBuilder(a, opts).Build()
		s.Require().NoError(err)
	})

	// Without verbose, displayZippedFiles must not run.
	s.NotContains(out, "-->", "non-verbose build should not list entries")
}

func (s *ArchiveTestSuite) TestBuild_Verbose_LogsExtraFileSearch() {
	// Exercises getFileReferencesFromFiles' logVerbose("Searching %s ...") branch.
	a := setupWorkspace(s.T(), "TestWorkspace")
	wd, _ := os.Getwd()
	writeDefineModule(s.T(), wd)
	extraDir := createExtraLibDir(s.T(), wd)

	opts := defaultOpts()
	opts.Verbose = true
	opts.IncludeFilesNotInWorkspace = true
	opts.ExtraFileSearchLocations = []string{extraDir}

	out := captureStdout(s.T(), func() {
		err := NewBuilder(a, opts).Build()
		s.Require().NoError(err)
	})

	// getFileReferencesFromFiles logs "Searching <file> for references..."
	s.Contains(out, "Searching", "verbose output should mention file search activity")
}

// ---------------------------------------------------------------------------
// Build – flat APW in zip
// ---------------------------------------------------------------------------

func (s *ArchiveTestSuite) TestBuild_WorkspaceAPW_HasFlatPaths() {
	// After Build() the .apw entry in the zip must have FilePathName values that
	// are bare filenames (no subdirectory prefix), because the archive is flat.
	a := setupWorkspace(s.T(), "TestWorkspace")
	err := NewBuilder(a, defaultOpts()).Build()
	s.Require().NoError(err)

	wd, _ := os.Getwd()
	zr, err := zip.OpenReader(filepath.Join(wd, "TestWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	// Find and read the .apw entry.
	var apwData []byte
	for _, f := range zr.File {
		if f.Name == "TestWorkspace.apw" {
			rc, err := f.Open()
			s.Require().NoError(err)
			apwData, err = io.ReadAll(rc)
			_ = rc.Close()
			s.Require().NoError(err)
			break
		}
	}

	s.Require().NotEmpty(apwData, "zip must contain TestWorkspace.apw")

	// Parse the archived APW and verify all FileRef paths are bare filenames.
	parsed, err := apw.Parse("TestWorkspace.apw", apwData)
	s.Require().NoError(err)

	for _, proj := range parsed.Workspace().Projects {
		for _, sys := range proj.Systems {
			for _, fr := range sys.Files {
				base := filepath.Base(fr.FilePathName)
				s.Equal(base, fr.FilePathName,
					"FileRef path in archived .apw should be a bare filename, got %q", fr.FilePathName)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Build – deduplication
// ---------------------------------------------------------------------------

// sharedIncludeAPW returns a workspace with two systems that both reference the
// same shared include file, simulating the real-world pattern that triggers
// duplicate zip entries without deduplication.
func sharedIncludeAPW(id string) []byte {
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
            <Identifier>System1</Identifier>
            <SysID>1</SysID>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>Main1</Identifier>
                <FilePathName>Source\Main1.axs</FilePathName>
                <Comments></Comments>
            </File>
            <File CompileType="Netlinx" Type="Include">
                <Identifier>Shared</Identifier>
                <FilePathName>Include\Shared.axi</FilePathName>
                <Comments></Comments>
            </File>
        </System>
        <System IsActive="false" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>System2</Identifier>
            <SysID>2</SysID>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>Main2</Identifier>
                <FilePathName>Source\Main2.axs</FilePathName>
                <Comments></Comments>
            </File>
            <File CompileType="Netlinx" Type="Include">
                <Identifier>Shared</Identifier>
                <FilePathName>Include\Shared.axi</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
}

func (s *ArchiveTestSuite) TestBuild_SharedFile_AddedOnlyOnce() {
	// A file referenced by multiple systems must appear exactly once in the zip.
	dir := s.T().TempDir()

	data := sharedIncludeAPW("SharedWorkspace")
	apwPath := filepath.Join(dir, "SharedWorkspace.apw")
	s.Require().NoError(os.WriteFile(apwPath, data, 0o644))

	for _, sub := range []struct{ d, n string }{
		{"Source", "Main1.axs"},
		{"Source", "Main2.axs"},
		{"Include", "Shared.axi"},
	} {
		s.Require().NoError(os.MkdirAll(filepath.Join(dir, sub.d), 0o755))
		s.Require().NoError(os.WriteFile(filepath.Join(dir, sub.d, sub.n), []byte(`PROGRAM_NAME='test'`), 0o644))
	}

	a, err := apw.Parse(apwPath, data)
	s.Require().NoError(err)

	oldWd, _ := os.Getwd()
	s.T().Cleanup(func() { _ = os.Chdir(oldWd) })
	s.Require().NoError(os.Chdir(dir))

	s.Require().NoError(NewBuilder(a, defaultOpts()).Build())

	zr, err := zip.OpenReader(filepath.Join(dir, "SharedWorkspace.zip"))
	s.Require().NoError(err)
	defer func() { _ = zr.Close() }()

	count := 0
	for _, f := range zr.File {
		if f.Name == "Shared.axi" {
			count++
		}
	}

	s.Equal(1, count, "Shared.axi should appear exactly once in the zip, got %d entries", count)
}

// ---------------------------------------------------------------------------
// sanitizeSegment
// ---------------------------------------------------------------------------

func TestSanitizeSegment_NoSpaces(t *testing.T) {
	require.Equal(t, "KingstonUniversity", sanitizeSegment("KingstonUniversity"))
}

func TestSanitizeSegment_SingleSpace(t *testing.T) {
	require.Equal(t, "Kingston-University", sanitizeSegment("Kingston University"))
}

func TestSanitizeSegment_MultipleSpaces(t *testing.T) {
	require.Equal(t, "Large-Classroom-Type-C", sanitizeSegment("Large Classroom Type C"))
}

func TestSanitizeSegment_AlreadyHyphenated(t *testing.T) {
	require.Equal(t, "KU-Large-Classroom", sanitizeSegment("KU-Large-Classroom"))
}

func TestSanitizeSegment_Empty(t *testing.T) {
	require.Equal(t, "", sanitizeSegment(""))
}

// ---------------------------------------------------------------------------
// Build – output filename sanitization and scoped naming
// ---------------------------------------------------------------------------

func (s *ArchiveTestSuite) TestBuild_OutputFile_SpacesInWorkspaceID_ReplacedWithHyphens() {
	a := setupWorkspace(s.T(), "Kingston University")
	opts := defaultOpts()
	opts.OutputFileSuffix = "archive.zip"
	s.Require().NoError(NewBuilder(a, opts).Build())

	wd, _ := os.Getwd()
	_, statErr := os.Stat(filepath.Join(wd, "Kingston-University.archive.zip"))
	s.NoError(statErr, "spaces in workspace ID should be replaced with hyphens in archive filename")
}

func (s *ArchiveTestSuite) TestBuild_OutputFile_SpacesInProjectID_ReplacedWithHyphens() {
	a := setupWorkspaceWithIDs(s.T(), "TestWorkspace", "My Project", "TestSystem")
	opts := defaultOpts()
	opts.OutputFileSuffix = "archive.zip"
	opts.ProjectID = "My Project"
	s.Require().NoError(NewBuilder(a, opts).Build())

	wd, _ := os.Getwd()
	_, statErr := os.Stat(filepath.Join(wd, "TestWorkspace-My-Project.archive.zip"))
	s.NoError(statErr, "spaces in project ID should be replaced with hyphens in archive filename")
}

func (s *ArchiveTestSuite) TestBuild_OutputFile_SpacesInSystemID_ReplacedWithHyphens() {
	a := setupWorkspaceWithIDs(s.T(), "TestWorkspace", "TestProject", "My System")
	opts := defaultOpts()
	opts.OutputFileSuffix = "archive.zip"
	opts.ProjectID = "TestProject"
	opts.SystemID = "My System"
	s.Require().NoError(NewBuilder(a, opts).Build())

	wd, _ := os.Getwd()
	_, statErr := os.Stat(filepath.Join(wd, "TestWorkspace-TestProject-My-System.archive.zip"))
	s.NoError(statErr, "spaces in system ID should be replaced with hyphens in archive filename")
}

func (s *ArchiveTestSuite) TestBuild_OutputFile_AllSegmentsWithSpaces_AllSanitized() {
	a := setupWorkspaceWithIDs(s.T(), "Kingston University", "My Project", "My System")
	opts := defaultOpts()
	opts.OutputFileSuffix = "archive.zip"
	opts.ProjectID = "My Project"
	opts.SystemID = "My System"
	s.Require().NoError(NewBuilder(a, opts).Build())

	wd, _ := os.Getwd()
	_, statErr := os.Stat(filepath.Join(wd, "Kingston-University-My-Project-My-System.archive.zip"))
	s.NoError(statErr, "all three segments with spaces should produce a fully hyphenated filename")
}

func (s *ArchiveTestSuite) TestBuild_OutputFile_ProjectAndSystemIDAppended() {
	// Regression: when both ProjectID and SystemID are set, both must appear in
	// the filename (guards against the auto-resolve bug where SystemID was dropped).
	a := setupWorkspace(s.T(), "TestWorkspace")
	opts := defaultOpts()
	opts.OutputFileSuffix = "archive.zip"
	opts.ProjectID = "TestProject"
	opts.SystemID = "TestSystem"
	s.Require().NoError(NewBuilder(a, opts).Build())

	wd, _ := os.Getwd()
	_, statErr := os.Stat(filepath.Join(wd, "TestWorkspace-TestProject-TestSystem.archive.zip"))
	s.NoError(statErr, "archive filename should include both project and system ID suffixes")
}
