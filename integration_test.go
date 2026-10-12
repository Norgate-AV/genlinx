package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx/internal/apw"
	archivepkg "github.com/Norgate-AV/genlinx/internal/archive"
	cfgpkg "github.com/Norgate-AV/genlinx/internal/cfg"
	"github.com/Norgate-AV/genlinx/internal/options"
)

// ---------------------------------------------------------------------------
// APW XML templates
// ---------------------------------------------------------------------------

// singleProjectAPW returns XML bytes for a minimal single-project workspace
// with one MasterSrc file (Source\Main.axs) and one Module file (Module\Module1.axs).
func singleProjectAPW(wsID string) []byte {
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
        <Identifier>MainProject</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>MainSystem</Identifier>
            <SysID>1</SysID>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>Main</Identifier>
                <FilePathName>Source\Main.axs</FilePathName>
                <Comments></Comments>
            </File>
            <File CompileType="Netlinx" Type="Module">
                <Identifier>Module1</Identifier>
                <FilePathName>Module\Module1.axs</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
}

// multiProjectAPW returns XML bytes for a workspace with two projects,
// "Alpha" (Source\MainA.axs) and "Beta" (Source\MainB.axs), for testing
// scoped archive builds.
func multiProjectAPW(wsID string) []byte {
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
        <Identifier>Alpha</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>SystemA</Identifier>
            <SysID>1</SysID>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>MainA</Identifier>
                <FilePathName>Source\MainA.axs</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
    <Project>
        <Identifier>Beta</Identifier>
        <System IsActive="true" Platform="Netlinx" Transport="Serial" TransportEx="TCPIP">
            <Identifier>SystemB</Identifier>
            <SysID>2</SysID>
            <File CompileType="Netlinx" Type="MasterSrc">
                <Identifier>MainB</Identifier>
                <FilePathName>Source\MainB.axs</FilePathName>
                <Comments></Comments>
            </File>
        </System>
    </Project>
</Workspace>`)
}

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

// IntegrationTestSuite exercises genlinx subsystems end-to-end.
type IntegrationTestSuite struct {
	suite.Suite
	tempDir string
	origWd  string
}

// SetupTest runs before each test method. It creates a fresh temporary
// directory (with any symlinks resolved for stable path comparisons) and
// records the original working directory so TearDownTest can restore it.
func (suite *IntegrationTestSuite) SetupTest() {
	t := suite.T()

	dir := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}

	suite.tempDir = dir

	wd, err := os.Getwd()
	require.NoError(t, err)
	suite.origWd = wd
}

// TearDownTest restores the working directory changed by any test method.
func (suite *IntegrationTestSuite) TearDownTest() {
	_ = os.Chdir(suite.origWd)
}

// setupWorkspace writes the APW and the listed source files under suite.tempDir
// and returns the parsed *apw.APW. Pass a nil or empty slice when source files
// are not needed on disk (e.g. CFG-only tests that never stat the files).
func (suite *IntegrationTestSuite) setupWorkspace(wsID string, data []byte, files []struct{ dir, name string }) *apw.APW {
	t := suite.T()
	t.Helper()

	apwPath := filepath.Join(suite.tempDir, wsID+".apw")
	require.NoError(t, os.WriteFile(apwPath, data, 0o644))

	for _, f := range files {
		subDir := filepath.Join(suite.tempDir, f.dir)
		require.NoError(t, os.MkdirAll(subDir, 0o755))
		require.NoError(t, os.WriteFile(
			filepath.Join(subDir, f.name),
			[]byte("PROGRAM_NAME='integration_test'"),
			0o644,
		))
	}

	ws, err := apw.Parse(apwPath, data)
	require.NoError(t, err)

	return ws
}

// chdirToTemp changes the working directory to suite.tempDir.
// TearDownTest restores the original directory after every test.
func (suite *IntegrationTestSuite) chdirToTemp() {
	require.NoError(suite.T(), os.Chdir(suite.tempDir))
}

// containsEntry reports whether name appears in the slice.
func containsEntry(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}

	return false
}

// ---------------------------------------------------------------------------
// Test 1 – APW parsing
// ---------------------------------------------------------------------------

// TestAPWParsing_ValidWorkspaceFile verifies that a well-formed APW file is
// parsed correctly and that the workspace identifier and file listings are
// populated.
func (suite *IntegrationTestSuite) TestAPWParsing_ValidWorkspaceFile() {
	t := suite.T()
	wsID := "ParseTest"
	data := singleProjectAPW(wsID)
	ws := suite.setupWorkspace(wsID, data, []struct{ dir, name string }{
		{"Source", "Main.axs"},
		{"Module", "Module1.axs"},
	})

	assert.Equal(t, wsID, ws.ID())
	assert.NotEmpty(t, ws.MasterSrcFiles(), "should have at least one MasterSrc file")
	assert.NotEmpty(t, ws.ModuleFiles(), "should have at least one Module file")
}

// ---------------------------------------------------------------------------
// Test 2 – CFG pipeline: expected content
// ---------------------------------------------------------------------------

// TestCfgBuild_ContainsExpectedContent verifies that Build returns a string
// containing all the required CFG sections and key=value entries.
func (suite *IntegrationTestSuite) TestCfgBuild_ContainsExpectedContent() {
	t := suite.T()
	wsID := "CfgTest"
	ws := suite.setupWorkspace(wsID, singleProjectAPW(wsID), nil)

	opts := &cfgpkg.Options{
		OutputLogFileSuffix:    "build.log",
		OutputLogFileOption:    "N",
		OutputLogConsoleOption: true,
	}
	output := cfgpkg.NewBuilder(ws, opts).Build()

	assert.Contains(t, output, "AXSFILE=")
	assert.Contains(t, output, "Main.axs")
	assert.Contains(t, output, "Module1.axs")
	assert.Contains(t, output, "MainAXSRootDirectory=")
	assert.Contains(t, output, "OutputLogFile=CfgTest.build.log")
	assert.Contains(t, output, "OutputLogFileOption=N")
	assert.Contains(t, output, "BuildWithDebugInformation=N")
}

// ---------------------------------------------------------------------------
// Test 3 – CFG pipeline: CRLF regression
// ---------------------------------------------------------------------------

// TestCfgBuild_NoCRLFInOutput is a regression test for the bug where git's
// autocrlf on Windows caused the embedded LICENSE to contain \r\n, which
// propagated into the generated CFG output via convertToComment.
func (suite *IntegrationTestSuite) TestCfgBuild_NoCRLFInOutput() {
	t := suite.T()
	wsID := "CRLFTest"
	ws := suite.setupWorkspace(wsID, singleProjectAPW(wsID), nil)

	output := cfgpkg.NewBuilder(ws, &cfgpkg.Options{}).Build()

	assert.False(t, strings.Contains(output, "\r"), "CFG output must not contain carriage returns")
}

// ---------------------------------------------------------------------------
// Test 4 – Local config find-up
// ---------------------------------------------------------------------------

// TestLocalConfig_FindUp verifies that LoadLocalConfig walks up the directory
// tree and discovers a config file placed in an ancestor directory.
func (suite *IntegrationTestSuite) TestLocalConfig_FindUp() {
	t := suite.T()

	configPath := filepath.Join(suite.tempDir, ".genlinxrc.json")
	require.NoError(t, os.WriteFile(configPath, []byte(`{}`), 0o644))

	subDir := filepath.Join(suite.tempDir, "sub")
	require.NoError(t, os.MkdirAll(subDir, 0o755))
	require.NoError(t, os.Chdir(subDir))

	result, err := options.LoadLocalConfig()
	require.NoError(t, err)

	assert.True(t, result.Found, "config should be discovered via find-up")
	assert.Equal(t, configPath, result.Path)
}

// ---------------------------------------------------------------------------
// Test 5 – All six config filename formats are discovered
// ---------------------------------------------------------------------------

// TestLocalConfig_AllFormatsDiscovered verifies that every supported config
// filename variant is recognised by LoadLocalConfig.
func (suite *IntegrationTestSuite) TestLocalConfig_AllFormatsDiscovered() {
	t := suite.T()

	formats := []struct {
		name    string
		content string
	}{
		{".genlinxrc.json", `{}`},
		{".genlinxrc.yaml", `{}`},
		{".genlinxrc.yml", `{}`},
		{".genlinx.json", `{}`},
		{".genlinx.yaml", `{}`},
		{".genlinx.yml", `{}`},
	}

	for _, f := range formats {
		f := f
		t.Run(f.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, f.name), []byte(f.content), 0o644))

			oldWd, err := os.Getwd()
			require.NoError(t, err)
			t.Cleanup(func() { _ = os.Chdir(oldWd) })
			require.NoError(t, os.Chdir(dir))

			result, err := options.LoadLocalConfig()
			require.NoError(t, err)
			assert.True(t, result.Found, "config file %s should be discovered", f.name)
		})
	}
}

// ---------------------------------------------------------------------------
// Test 6 – LoadCfgOptions merges local config paths
// ---------------------------------------------------------------------------

// TestCfgOptions_LocalConfigMergesIncludePaths verifies that includePath and
// modulePath values from a local .genlinxrc.json are resolved to absolute
// paths and exposed through LoadCfgOptions.
func (suite *IntegrationTestSuite) TestCfgOptions_LocalConfigMergesIncludePaths() {
	t := suite.T()

	// Isolate the global config lookup so the real machine's config cannot
	// interfere with the expected path values.
	globalIsolateDir := t.TempDir()
	t.Setenv("GENLINX_CONFIG_DIR", globalIsolateDir)

	const configJSON = `{
		"compiler": {
			"includePath": ["custom/include"],
			"modulePath":  ["custom/modules"]
		}
	}`
	require.NoError(t, os.WriteFile(
		filepath.Join(suite.tempDir, ".genlinxrc.json"),
		[]byte(configJSON),
		0o644,
	))

	suite.chdirToTemp()

	opts, info, err := options.LoadCfgOptions(nil)
	require.NoError(t, err)
	require.True(t, info.LocalResult.Found, "local config should be found")

	// resolveConfigPaths converts relative paths to absolute using the CWD.
	wantInclude := filepath.Join(suite.tempDir, "custom", "include")
	wantModule := filepath.Join(suite.tempDir, "custom", "modules")

	assert.Contains(t, opts.IncludePath, wantInclude)
	assert.Contains(t, opts.ModulePath, wantModule)
}

// ---------------------------------------------------------------------------
// Test 7 – Archive pipeline: valid ZIP with workspace files
// ---------------------------------------------------------------------------

// TestArchiveBuild_ProducesValidZip verifies that Build creates a zip archive
// in the current directory containing the workspace file and all source files.
func (suite *IntegrationTestSuite) TestArchiveBuild_ProducesValidZip() {
	t := suite.T()
	wsID := "ArchiveTest"
	ws := suite.setupWorkspace(wsID, singleProjectAPW(wsID), []struct{ dir, name string }{
		{"Source", "Main.axs"},
		{"Module", "Module1.axs"},
	})

	suite.chdirToTemp()

	require.NoError(t, archivepkg.NewBuilder(ws, &archivepkg.Options{OutputFileSuffix: "zip"}).Build())

	r, err := zip.OpenReader(filepath.Join(suite.tempDir, wsID+".zip"))
	require.NoError(t, err)
	defer func() { require.NoError(t, r.Close()) }()

	var names []string
	for _, f := range r.File {
		names = append(names, f.Name)
	}

	assert.True(t, containsEntry(names, wsID+".apw"), "zip should contain the workspace file")
	assert.True(t, containsEntry(names, "Main.axs"), "zip should contain the MasterSrc file")
	assert.True(t, containsEntry(names, "Module1.axs"), "zip should contain the Module file")
}

// ---------------------------------------------------------------------------
// Test 8 – Archive pipeline: scoped to a single project
// ---------------------------------------------------------------------------

// TestArchiveBuild_ScopedToProject verifies that when ProjectID is set, Build
// archives only the files belonging to that project and excludes all others.
func (suite *IntegrationTestSuite) TestArchiveBuild_ScopedToProject() {
	t := suite.T()
	wsID := "ScopedWS"
	ws := suite.setupWorkspace(wsID, multiProjectAPW(wsID), []struct{ dir, name string }{
		{"Source", "MainA.axs"},
		{"Source", "MainB.axs"},
	})

	suite.chdirToTemp()

	opts := &archivepkg.Options{
		OutputFileSuffix: "zip",
		ProjectID:        "Alpha",
	}
	require.NoError(t, archivepkg.NewBuilder(ws, opts).Build())

	r, err := zip.OpenReader(filepath.Join(suite.tempDir, wsID+"-Alpha.zip"))
	require.NoError(t, err)
	defer func() { require.NoError(t, r.Close()) }()

	var names []string
	for _, f := range r.File {
		names = append(names, f.Name)
	}

	assert.True(t, containsEntry(names, "MainA.axs"), "scoped zip should contain the Alpha project's source file")
	assert.False(t, containsEntry(names, "MainB.axs"), "scoped zip must not contain the Beta project's source file")
}

// ---------------------------------------------------------------------------
// Runner
// ---------------------------------------------------------------------------

// TestIntegrationTestSuite runs the integration test suite.
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
