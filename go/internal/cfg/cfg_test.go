package cfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx-go/internal/apw"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// minimalAPW returns valid APW XML with one Include, one Module and one MasterSrc.
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

func parseAPW(t *testing.T, id string) *apw.APW {
	t.Helper()
	dir := t.TempDir()
	data := minimalAPW(id)
	apwPath := filepath.Join(dir, id+".apw")
	require.NoError(t, os.WriteFile(apwPath, data, 0o644))
	a, err := apw.Parse(apwPath, data)
	require.NoError(t, err)
	return a
}

func defaultOpts() *Options {
	return &Options{
		OutputFileSuffix:          "build.cfg",
		OutputLogFileSuffix:       "build.log",
		OutputLogFileOption:       "N",
		OutputLogConsoleOption:    true,
		BuildWithDebugInformation: false,
		BuildWithSource:           false,
	}
}

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

type CfgTestSuite struct {
	suite.Suite
}

func TestCfgSuite(t *testing.T) {
	suite.Run(t, new(CfgTestSuite))
}

// ---------------------------------------------------------------------------
// deduplicate
// ---------------------------------------------------------------------------

func TestDeduplicate_UniqueItems(t *testing.T) {
	input := []string{`C:\Foo`, `C:\Bar`, `C:\Baz`}
	result := deduplicate(input)
	assert.Len(t, result, 3)
}

func TestDeduplicate_ExactDuplicates(t *testing.T) {
	input := []string{`C:\Foo`, `C:\Foo`, `C:\Foo`}
	result := deduplicate(input)
	assert.Len(t, result, 1)
}

func TestDeduplicate_SlashVariants(t *testing.T) {
	// Forward-slash and back-slash forms of the same path must deduplicate.
	input := []string{
		`C:\Program Files (x86)\Common Files\AMXShare\AXIs`,
		`C:/Program Files (x86)/Common Files/AMXShare/AXIs`,
	}
	result := deduplicate(input)
	assert.Len(t, result, 1)
	// The retained entry must have OS-native separators.
	assert.Equal(t, filepath.FromSlash(input[0]), result[0])
}

func TestDeduplicate_EmptyInput(t *testing.T) {
	result := deduplicate(nil)
	assert.Empty(t, result)
}

func TestDeduplicate_PreservesOrder(t *testing.T) {
	input := []string{`C:\C`, `C:\A`, `C:\B`}
	result := deduplicate(input)
	assert.Equal(t, []string{`C:\C`, `C:\A`, `C:\B`}, result)
}

// ---------------------------------------------------------------------------
// Builder.Build – structural keys
// ---------------------------------------------------------------------------

func (s *CfgTestSuite) TestBuild_ContainsMainAXSRootDirectory() {
	a := parseAPW(s.T(), "TestWorkspace")
	output := NewBuilder(a, defaultOpts()).Build()
	s.Contains(output, "MainAXSRootDirectory=")
}

func (s *CfgTestSuite) TestBuild_ContainsOutputLogFile() {
	a := parseAPW(s.T(), "TestWorkspace")
	output := NewBuilder(a, defaultOpts()).Build()
	s.Contains(output, "OutputLogFile=TestWorkspace.build.log")
}

func (s *CfgTestSuite) TestBuild_ContainsOutputLogFileOption() {
	a := parseAPW(s.T(), "TestWorkspace")
	output := NewBuilder(a, defaultOpts()).Build()
	s.Contains(output, "OutputLogFileOption=N")
}

func (s *CfgTestSuite) TestBuild_OutputLogConsoleOption_Y() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	opts.OutputLogConsoleOption = true
	output := NewBuilder(a, opts).Build()
	s.Contains(output, "OutputLogConsoleOption=Y")
}

func (s *CfgTestSuite) TestBuild_OutputLogConsoleOption_N() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	opts.OutputLogConsoleOption = false
	output := NewBuilder(a, opts).Build()
	s.Contains(output, "OutputLogConsoleOption=N")
}

func (s *CfgTestSuite) TestBuild_BuildWithDebugInformation_Y() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	opts.BuildWithDebugInformation = true
	output := NewBuilder(a, opts).Build()
	s.Contains(output, "BuildWithDebugInformation=Y")
}

func (s *CfgTestSuite) TestBuild_BuildWithDebugInformation_N() {
	a := parseAPW(s.T(), "TestWorkspace")
	output := NewBuilder(a, defaultOpts()).Build()
	s.Contains(output, "BuildWithDebugInformation=N")
}

func (s *CfgTestSuite) TestBuild_BuildWithSource_Y() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	opts.BuildWithSource = true
	output := NewBuilder(a, opts).Build()
	s.Contains(output, "BuildWithSource=Y")
}

func (s *CfgTestSuite) TestBuild_ContainsAXSFILE_ForMasterSrc() {
	a := parseAPW(s.T(), "TestWorkspace")
	output := NewBuilder(a, defaultOpts()).Build()
	s.Contains(output, "AXSFILE=")
	s.Contains(output, "TestMain.axs")
}

func (s *CfgTestSuite) TestBuild_ContainsAXSFILE_ForModule() {
	a := parseAPW(s.T(), "TestWorkspace")
	output := NewBuilder(a, defaultOpts()).Build()
	s.Contains(output, "TestModule.axs")
}

// ---------------------------------------------------------------------------
// Additional paths – deduplication
// ---------------------------------------------------------------------------

func (s *CfgTestSuite) TestBuild_IncludePaths_Deduplicated() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	// Same path, both slash styles — should appear only once in output.
	opts.IncludePath = []string{
		`C:\Program Files (x86)\Common Files\AMXShare\AXIs`,
		`C:/Program Files (x86)/Common Files/AMXShare/AXIs`,
	}

	output := NewBuilder(a, opts).Build()
	count := strings.Count(output, `AdditionalIncludePath=`+filepath.FromSlash(`C:\Program Files (x86)\Common Files\AMXShare\AXIs`))
	s.Equal(1, count, "duplicate include path should appear exactly once")
}

func (s *CfgTestSuite) TestBuild_ModulePaths_Deduplicated() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	opts.ModulePath = []string{
		`C:\Program Files (x86)\Common Files\AMXShare\Duet\bundle`,
		`C:/Program Files (x86)/Common Files/AMXShare/Duet/bundle`,
	}

	output := NewBuilder(a, opts).Build()
	count := strings.Count(output, `AdditionalModulePath=`+filepath.FromSlash(`C:\Program Files (x86)\Common Files\AMXShare\Duet\bundle`))
	s.Equal(1, count, "duplicate module path should appear exactly once")
}

func (s *CfgTestSuite) TestBuild_LibraryPaths_Deduplicated() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	// This specific case was the bug that was fixed.
	opts.LibraryPath = []string{
		`C:\Program Files (x86)\Common Files\AMXShare\SYCs`,
		`C:/Program Files (x86)/Common Files/AMXShare/SYCs`,
	}

	output := NewBuilder(a, opts).Build()
	count := strings.Count(output, `AdditionalLibraryPath=`+filepath.FromSlash(`C:\Program Files (x86)\Common Files\AMXShare\SYCs`))
	s.Equal(1, count, "duplicate library path should appear exactly once")
}

func (s *CfgTestSuite) TestBuild_IncludePaths_MergedFromWorkspace() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	// The workspace has an Include file, so its directory should appear in the output.
	output := NewBuilder(a, opts).Build()
	s.Contains(output, "AdditionalIncludePath=")
	s.Contains(output, "Include") // workspace Include dir
}

func (s *CfgTestSuite) TestBuild_ModulePaths_MergedFromWorkspace() {
	a := parseAPW(s.T(), "TestWorkspace")
	opts := defaultOpts()
	output := NewBuilder(a, opts).Build()
	s.Contains(output, "AdditionalModulePath=")
	s.Contains(output, "Module") // workspace Module dir
}
