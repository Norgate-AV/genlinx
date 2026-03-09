package apw

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

type MarshalTestSuite struct {
	suite.Suite
}

func TestMarshalSuite(t *testing.T) {
	suite.Run(t, new(MarshalTestSuite))
}

// ---------------------------------------------------------------------------
// Marshal — document structure
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestMarshal_ContainsXMLHeader() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	s.True(bytes.HasPrefix(data, []byte("<?xml version")),
		"output must begin with XML declaration")
}

func (s *MarshalTestSuite) TestMarshal_ContainsDOCTYPE() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	s.Contains(string(data), "<!DOCTYPE Workspace [",
		"output must contain DOCTYPE so Parse can accept it")
}

func (s *MarshalTestSuite) TestMarshal_ContainsWorkspaceElement() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	s.Contains(string(data), "<Workspace ")
	s.Contains(string(data), "</Workspace>")
}

func (s *MarshalTestSuite) TestMarshal_ContainsCurrentVersionAttr() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	s.Contains(string(data), `CurrentVersion="4.0"`)
}

func (s *MarshalTestSuite) TestMarshal_PreservesIdentifier() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	s.Contains(string(data), "<Identifier>TestWorkspace</Identifier>")
}

func (s *MarshalTestSuite) TestMarshal_PreservesProjectIdentifier() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	s.Contains(string(data), "<Identifier>TestProject</Identifier>")
}

func (s *MarshalTestSuite) TestMarshal_SystemAttributesPresent() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, `IsActive="true"`)
	s.Contains(out, `Platform="Netlinx"`)
	s.Contains(out, `Transport="Serial"`)
	s.Contains(out, `TransportEx="TCPIP"`)
}

func (s *MarshalTestSuite) TestMarshal_FileRefAttributesPresent() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, `Type="Include"`)
	s.Contains(out, `Type="Module"`)
	s.Contains(out, `Type="MasterSrc"`)
	s.Contains(out, `CompileType="NetLinx"`)
}

func (s *MarshalTestSuite) TestMarshal_FilePathNamesPresent() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, "TestInclude.axi")
	s.Contains(out, "TestModule.axs")
	s.Contains(out, "TestMain.axs")
}

func (s *MarshalTestSuite) TestMarshal_NilWorkspace_ReturnsError() {
	_, err := Marshal(nil)
	s.Error(err)
	s.Contains(err.Error(), "nil")
}

// ---------------------------------------------------------------------------
// Marshal — optional metadata fields round-trip
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestMarshal_OptionalFields_Preserved() {
	ws := minimalWorkspace()
	ws.Comments = "project-level comment"
	ws.Projects[0].Designer = "Jane Smith"
	ws.Projects[0].DealerID = "D-001"
	ws.Projects[0].SalesOrder = "SO-999"
	ws.Projects[0].Systems[0].Comments = "system comment"

	data, err := Marshal(ws)
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, "project-level comment")
	s.Contains(out, "Jane Smith")
	s.Contains(out, "D-001")
	s.Contains(out, "SO-999")
	s.Contains(out, "system comment")
}

func (s *MarshalTestSuite) TestMarshal_EmptyOptionalFields_Omitted() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	out := string(data)
	// Designer and DealerID are genuinely optional — omitted when not set.
	s.NotContains(out, "<Designer>")
	s.NotContains(out, "<DealerID>")
}

func (s *MarshalTestSuite) TestMarshal_WorkspaceFields_AlwaysPresent() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)
	out := string(data)
	// NetLinx Studio always writes these elements even when empty.
	s.Contains(out, "<PJS_File></PJS_File>")
	s.Contains(out, "<PJS_ConvertDate></PJS_ConvertDate>")
	s.Contains(out, "<PJS_CreateDate></PJS_CreateDate>")
}

// ---------------------------------------------------------------------------
// Marshal — round-trip via Parse
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestMarshal_RoundTrip_ParseSucceeds() {
	data, err := Marshal(minimalWorkspace())
	s.Require().NoError(err)

	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "TestWorkspace.apw")
	s.Require().NoError(os.WriteFile(apwPath, data, 0o644))

	_, err = Parse(apwPath, data)
	s.NoError(err, "marshaled output must be accepted by Parse")
}

func (s *MarshalTestSuite) TestMarshal_RoundTrip_PreservesFileIdentifiers() {
	ws := minimalWorkspace()
	data, err := Marshal(ws)
	s.Require().NoError(err)

	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "TestWorkspace.apw")
	s.Require().NoError(os.WriteFile(apwPath, data, 0o644))

	a, err := Parse(apwPath, data)
	s.Require().NoError(err)

	// Collect IDs from re-parsed workspace files (excluding workspace-type entry).
	ids := make(map[string]bool)
	for _, f := range a.AllFiles() {
		if f.Type != FileTypeWorkspace {
			ids[f.ID] = true
		}
	}

	s.True(ids["TestInclude"], "expected TestInclude after round-trip")
	s.True(ids["TestModule"], "expected TestModule after round-trip")
	s.True(ids["TestMain"], "expected TestMain after round-trip")
}

func (s *MarshalTestSuite) TestMarshal_RoundTrip_PreservesWorkspaceIdentifier() {
	ws := minimalWorkspace()
	ws.Identifier = "My Custom Workspace"
	data, err := Marshal(ws)
	s.Require().NoError(err)

	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "CustomName.apw")
	s.Require().NoError(os.WriteFile(apwPath, data, 0o644))

	a, err := Parse(apwPath, data)
	s.Require().NoError(err)
	// The APW id comes from the filename, not ws.Identifier, so verify
	// the workspace struct identifier is preserved in the XML content.
	s.Contains(string(data), "My Custom Workspace")
	// a.ID() comes from the filename stem
	s.Equal("CustomName", a.ID())
}

// ---------------------------------------------------------------------------
// APW.Bytes
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestBytes_ProducesSameAsMarshal() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "TestWorkspace.apw")

	ws := minimalWorkspace()
	expected, err := Marshal(ws)
	s.Require().NoError(err)

	a := NewAPW("TestWorkspace", apwPath)
	a.SetWorkspace(ws)

	actual, err := a.Bytes()
	s.Require().NoError(err)
	s.Equal(expected, actual)
}

func (s *MarshalTestSuite) TestBytes_ReturnsErrorForInvalidWorkspace() {
	a := NewAPW("X", "X.apw")
	a.ws = nil // force nil to trigger error path
	_, err := a.Bytes()
	s.Error(err)
}

// ---------------------------------------------------------------------------
// APW.Write
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestWrite_CreatesFileOnDisk() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "TestWorkspace.apw")

	a := NewAPW("TestWorkspace", apwPath)
	s.Require().NoError(a.Write(apwPath))

	_, err := os.Stat(apwPath)
	s.NoError(err, "Write must create the file")
}

func (s *MarshalTestSuite) TestWrite_FileCanBeParsedBack() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "TestWorkspace.apw")

	a := NewAPW("TestWorkspace", apwPath)
	s.Require().NoError(a.Write(apwPath))

	data, err := os.ReadFile(apwPath)
	s.Require().NoError(err)

	_, err = Parse(apwPath, data)
	s.NoError(err, "file written by Write must be accepted by Parse")
}

func (s *MarshalTestSuite) TestWrite_Error_BadPath() {
	a := NewAPW("X", "/tmp/X.apw")
	err := a.Write("/nonexistent-dir/sub/X.apw")
	s.Error(err)
}

// ---------------------------------------------------------------------------
// NewAPW
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestNewAPW_ID() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.Equal("MyProject", a.ID())
}

func (s *MarshalTestSuite) TestNewAPW_PathIsAbsolute() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.True(filepath.IsAbs(a.FilePath()), "FilePath must be absolute")
}

func (s *MarshalTestSuite) TestNewAPW_WorkspaceIdentifier() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.Equal("MyProject", a.ws.Identifier)
}

func (s *MarshalTestSuite) TestNewAPW_DefaultVersions() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.Equal("4.0", a.ws.CreateVersion)
	s.Equal("4.0", a.ws.CurrentVersion)
}

func (s *MarshalTestSuite) TestNewAPW_WorkspaceIsEmpty() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.Empty(a.ws.Projects, "NewAPW workspace must have no projects — use Workspace().AddProject()")
}

func (s *MarshalTestSuite) TestNewAPW_WorkspaceAccessor() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.NotNil(a.Workspace(), "Workspace() must return non-nil")
	s.Same(a.ws, a.Workspace(), "Workspace() must return the inner workspace pointer")
}

func (s *MarshalTestSuite) TestNewAPW_FilesMapInitialised() {
	a := NewAPW("MyProject", "MyProject.apw")
	s.NotNil(a.files, "files map must be non-nil")
	s.Empty(a.files, "files map must be empty on creation")
}

func (s *MarshalTestSuite) TestNewAPW_CanBeMarshaledAndParsed() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "NewProject.apw")

	a := NewAPW("NewProject", apwPath)
	data, err := a.Bytes()
	s.Require().NoError(err)

	s.Require().NoError(os.WriteFile(apwPath, data, 0o644))

	_, err = Parse(apwPath, data)
	s.NoError(err, "NewAPW output must round-trip through Parse")
}

// ---------------------------------------------------------------------------
// Workspace.AddProject
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestWorkspace_AddProject_AppearsInProjects() {
	a := NewAPW("P", "P.apw")
	proj := NewProject("MyProject")
	returned := a.Workspace().AddProject(proj)
	s.Same(proj, returned, "AddProject must return the same project pointer")
	s.Len(a.ws.Projects, 1)
	s.Equal("MyProject", a.ws.Projects[0].Identifier)
}

func (s *MarshalTestSuite) TestWorkspace_AddProject_MultipleProjects() {
	a := NewAPW("P", "P.apw")
	a.Workspace().AddProject(NewProject("Alpha"))
	a.Workspace().AddProject(NewProject("Beta"))
	s.Len(a.ws.Projects, 2)
}

// ---------------------------------------------------------------------------
// Workspace.RemoveProject
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestWorkspace_RemoveProject_RemovesMatchingProject() {
	a := NewAPW("P", "P.apw")
	a.Workspace().AddProject(NewProject("Alpha"))
	a.Workspace().AddProject(NewProject("Beta"))
	removed := a.Workspace().RemoveProject("Alpha")
	s.True(removed)
	s.Len(a.ws.Projects, 1)
	s.Equal("Beta", a.ws.Projects[0].Identifier)
}

func (s *MarshalTestSuite) TestWorkspace_RemoveProject_ReturnsFalseWhenNotFound() {
	a := NewAPW("P", "P.apw")
	a.Workspace().AddProject(NewProject("Alpha"))
	removed := a.Workspace().RemoveProject("NoSuchProject")
	s.False(removed)
	s.Len(a.ws.Projects, 1)
}

func (s *MarshalTestSuite) TestProject_AddSystem_AppearsInSystems() {
	proj := NewProject("P")
	sys := NewSystem("MySystem")
	returned := proj.AddSystem(sys)
	s.Same(sys, returned, "AddSystem must return the same system pointer")
	s.Len(proj.Systems, 1)
	s.Equal("MySystem", proj.Systems[0].Identifier)
}

func (s *MarshalTestSuite) TestProject_AddSystem_MultipleSystems() {
	proj := NewProject("P")
	proj.AddSystem(NewSystem("S1"))
	proj.AddSystem(NewSystem("S2"))
	s.Len(proj.Systems, 2)
}

// ---------------------------------------------------------------------------
// Project.RemoveSystem
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestProject_RemoveSystem_RemovesMatchingSystem() {
	proj := NewProject("P")
	proj.AddSystem(NewSystem("S1"))
	proj.AddSystem(NewSystem("S2"))
	removed := proj.RemoveSystem("S1")
	s.True(removed)
	s.Len(proj.Systems, 1)
	s.Equal("S2", proj.Systems[0].Identifier)
}

func (s *MarshalTestSuite) TestProject_RemoveSystem_ReturnsFalseWhenNotFound() {
	proj := NewProject("P")
	proj.AddSystem(NewSystem("S1"))
	removed := proj.RemoveSystem("NoSuchSystem")
	s.False(removed)
	s.Len(proj.Systems, 1)
}

func (s *MarshalTestSuite) TestSystem_AddFile_AppearsInFiles() {
	sys := newDefaultSystem()
	fr := NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude)
	returned := sys.AddFile(fr)
	s.Same(fr, returned, "AddFile must return the same FileRef pointer")
	s.Len(sys.Files, 1)
	s.Equal("Inc1", sys.Files[0].Identifier)
}

func (s *MarshalTestSuite) TestSystem_AddFile_InfersCompileType_NetLinx() {
	sys := newDefaultSystem()
	fr := NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude)
	s.Empty(fr.CompileType)
	sys.AddFile(fr)
	s.Equal(FileCompileTypeNetLinx, fr.CompileType)
}

func (s *MarshalTestSuite) TestSystem_AddFile_InfersCompileType_None() {
	sys := newDefaultSystem()
	fr := NewFileRef("Panel1", "Panel1.tp4", FileTypeTP4)
	sys.AddFile(fr)
	s.Equal(FileCompileTypeNone, fr.CompileType)
}

func (s *MarshalTestSuite) TestSystem_AddFile_PreservesExplicitCompileType() {
	sys := newDefaultSystem()
	fr := NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude)
	fr.CompileType = FileCompileTypeAxcess
	sys.AddFile(fr)
	s.Equal(FileCompileTypeAxcess, fr.CompileType)
}

func (s *MarshalTestSuite) TestSystem_AddFile_MultipleFiles() {
	sys := newDefaultSystem()
	sys.AddFile(NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude))
	sys.AddFile(NewFileRef("Mod1", "Module/Mod1.axs", FileTypeModule))
	sys.AddFile(NewFileRef("Main", "Source/Main.axs", FileTypeMasterSrc))
	s.Len(sys.Files, 3)
}

// ---------------------------------------------------------------------------
// System.RemoveFile
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestSystem_RemoveFile_RemovesMatchingFile() {
	sys := newDefaultSystem()
	sys.AddFile(NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude))
	sys.AddFile(NewFileRef("Main", "Source/Main.axs", FileTypeMasterSrc))
	removed := sys.RemoveFile("Inc1")
	s.True(removed)
	s.Len(sys.Files, 1)
	s.Equal("Main", sys.Files[0].Identifier)
}

func (s *MarshalTestSuite) TestSystem_RemoveFile_ReturnsFalseWhenNotFound() {
	sys := newDefaultSystem()
	sys.AddFile(NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude))
	removed := sys.RemoveFile("NoSuchFile")
	s.False(removed)
	s.Len(sys.Files, 1)
}

func (s *MarshalTestSuite) TestFileRef_AddDeviceMap_AppearsInDeviceMaps() {
	fr := NewFileRef("Panel1", "Panel1.tp4", FileTypeTP4)
	dm := NewDeviceMap("10001:1:0", "dvTP")
	returned := fr.AddDeviceMap(dm)
	s.Same(dm, returned, "AddDeviceMap must return the same DeviceMap pointer")
	s.Len(fr.DeviceMaps, 1)
	s.Equal("10001:1:0", fr.DeviceMaps[0].DevAddr)
	s.Equal("dvTP", fr.DeviceMaps[0].DevName)
}

func (s *MarshalTestSuite) TestFileRef_AddDeviceMap_MultipleDeviceMaps() {
	fr := NewFileRef("Panel1", "Panel1.tp4", FileTypeTP4)
	fr.AddDeviceMap(NewDeviceMap("10001:1:0", "dvTP1"))
	fr.AddDeviceMap(NewDeviceMap("10002:1:0", "dvTP2"))
	s.Len(fr.DeviceMaps, 2)
}

// ---------------------------------------------------------------------------
// FileRef.RemoveDeviceMap
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestFileRef_RemoveDeviceMap_RemovesMatchingDeviceMap() {
	fr := NewFileRef("Panel1", "Panel1.tp4", FileTypeTP4)
	fr.AddDeviceMap(NewDeviceMap("10001:1:0", "dvTP1"))
	fr.AddDeviceMap(NewDeviceMap("10002:1:0", "dvTP2"))
	removed := fr.RemoveDeviceMap("10001:1:0")
	s.True(removed)
	s.Len(fr.DeviceMaps, 1)
	s.Equal("10002:1:0", fr.DeviceMaps[0].DevAddr)
}

func (s *MarshalTestSuite) TestFileRef_RemoveDeviceMap_ReturnsFalseWhenNotFound() {
	fr := NewFileRef("Panel1", "Panel1.tp4", FileTypeTP4)
	fr.AddDeviceMap(NewDeviceMap("10001:1:0", "dvTP1"))
	removed := fr.RemoveDeviceMap("99999:1:0")
	s.False(removed)
	s.Len(fr.DeviceMaps, 1)
}

func (s *MarshalTestSuite) TestFileRef_AddDeviceMap_AppearsInMarshaledOutput() {
	dir := s.T().TempDir()
	a := NewAPW("P", filepath.Join(dir, "P.apw"))
	proj := a.Workspace().AddProject(NewProject("P"))
	sys := proj.AddSystem(newDefaultSystem())
	fr := NewFileRef("Panel1", "User Interface/Panel1.tp4", FileTypeTP4)
	fr.AddDeviceMap(NewDeviceMap("10001:1:0", "dvTP_Main"))
	sys.AddFile(fr)

	data, err := a.Bytes()
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, "DeviceMap")
	s.Contains(out, "10001:1:0")
	s.Contains(out, "dvTP_Main")
}

// ---------------------------------------------------------------------------
// FileRef.AddIRDB
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestFileRef_AddIRDB_AppearsInIRDBs() {
	fr := NewFileRef("IR1", "IR/Remote.irl", FileTypeIR)
	irdb := NewIRDB("Samsung_TV", "Samsung_TV.irdb")
	returned := fr.AddIRDB(irdb)
	s.Same(irdb, returned, "AddIRDB must return the same IRDB pointer")
	s.Len(fr.IRDBs, 1)
	s.Equal("Samsung_TV", fr.IRDBs[0].Property)
}

func (s *MarshalTestSuite) TestFileRef_AddIRDB_MultipleIRDBs() {
	fr := NewFileRef("IR1", "IR/Remote.irl", FileTypeIR)
	fr.AddIRDB(NewIRDB("Samsung_TV", "Samsung_TV.irdb"))
	fr.AddIRDB(NewIRDB("Sony_DVD", "Sony_DVD.irdb"))
	s.Len(fr.IRDBs, 2)
}

// ---------------------------------------------------------------------------
// FileRef.RemoveIRDB
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestFileRef_RemoveIRDB_RemovesMatchingIRDB() {
	fr := NewFileRef("IR1", "IR/Remote.irl", FileTypeIR)
	fr.AddIRDB(NewIRDB("Samsung_TV", "Samsung_TV.irdb"))
	fr.AddIRDB(NewIRDB("Sony_DVD", "Sony_DVD.irdb"))
	removed := fr.RemoveIRDB("Samsung_TV")
	s.True(removed)
	s.Len(fr.IRDBs, 1)
	s.Equal("Sony_DVD", fr.IRDBs[0].Property)
}

func (s *MarshalTestSuite) TestFileRef_RemoveIRDB_ReturnsFalseWhenNotFound() {
	fr := NewFileRef("IR1", "IR/Remote.irl", FileTypeIR)
	fr.AddIRDB(NewIRDB("Samsung_TV", "Samsung_TV.irdb"))
	removed := fr.RemoveIRDB("NoSuchDB")
	s.False(removed)
	s.Len(fr.IRDBs, 1)
}

func (s *MarshalTestSuite) TestFileRef_AddIRDB_AppearsInMarshaledOutput() {
	dir := s.T().TempDir()
	a := NewAPW("P", filepath.Join(dir, "P.apw"))
	proj := a.Workspace().AddProject(NewProject("P"))
	sys := proj.AddSystem(newDefaultSystem())
	fr := NewFileRef("IR1", "IR/Remote.irl", FileTypeIR)
	fr.AddIRDB(NewIRDB("Samsung_TV", "Samsung/Samsung_TV.irdb"))
	sys.AddFile(fr)

	data, err := a.Bytes()
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, "Samsung_TV")
	s.Contains(out, "Samsung_TV.irdb")
}

// ---------------------------------------------------------------------------
// APW.SetWorkspace / APW.Rebuild
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestAPW_SetWorkspace_ReplacesWorkspace() {
	a := NewAPW("P", "P.apw")
	ws := minimalWorkspace()
	a.SetWorkspace(ws)
	s.Same(ws, a.ws, "SetWorkspace must replace the inner workspace")
}

func (s *MarshalTestSuite) TestAPW_SetWorkspace_RebuildsFilesMap() {
	_, apwPath := writeAPW(s.T(), "TestWorkspace.apw", minimalAPW())
	a, err := Parse(apwPath, minimalAPW())
	s.Require().NoError(err)
	prevLen := len(a.files)
	s.Greater(prevLen, 0)

	// Replace with an empty workspace — files map must be cleared.
	a.SetWorkspace(NewWorkspace("Empty"))
	s.Empty(a.files)
}

func (s *MarshalTestSuite) TestAPW_Rebuild_SyncsFilesMap() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "P.apw")
	a := NewAPW("P", apwPath)

	proj := a.Workspace().AddProject(NewProject("P"))
	sys := proj.AddSystem(newDefaultSystem())
	sys.AddFile(NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude))
	sys.AddFile(NewFileRef("Main", "Source/Main.axs", FileTypeMasterSrc))

	// Before Rebuild, files map is empty (NewAPW starts with an empty workspace).
	s.Empty(a.files)

	a.Rebuild()

	s.Len(a.files, 2, "Rebuild must populate the files map from the workspace tree")
	ids := collectFileIDs(a.AllFiles())
	s.Contains(ids, "Inc1")
	s.Contains(ids, "Main")
}

func (s *MarshalTestSuite) TestAPW_Rebuild_UpdatesFilesByType() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "P.apw")
	a := NewAPW("P", apwPath)

	proj := a.Workspace().AddProject(NewProject("P"))
	sys := proj.AddSystem(newDefaultSystem())
	sys.AddFile(NewFileRef("Mod1", "Module/Mod1.axs", FileTypeModule))

	a.Rebuild()

	modules := a.ModuleFiles()
	s.Len(modules, 1)
	s.Equal("Mod1", modules[0].ID)
}

// ---------------------------------------------------------------------------
// Marshal output — files added via hierarchy
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestHierarchicalBuild_AppearsInMarshaledOutput() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "P.apw")
	a := NewAPW("P", apwPath)

	proj := a.Workspace().AddProject(NewProject("P"))
	sys := proj.AddSystem(newDefaultSystem())
	sys.AddFile(NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude))
	sys.AddFile(NewFileRef("Main", "Source/Main.axs", FileTypeMasterSrc))

	data, err := a.Bytes()
	s.Require().NoError(err)
	out := string(data)
	s.Contains(out, "Inc1")
	s.Contains(out, "Inc1.axi")
	s.Contains(out, "Main")
	s.Contains(out, "Main.axs")
	s.Contains(out, `Type="Include"`)
	s.Contains(out, `Type="MasterSrc"`)
	s.Contains(out, `CompileType="NetLinx"`)
}

func (s *MarshalTestSuite) TestHierarchicalBuild_BackslashPath_NormalizedInFilesMap() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "P.apw")
	a := NewAPW("P", apwPath)

	proj := a.Workspace().AddProject(NewProject("P"))
	sys := proj.AddSystem(newDefaultSystem())
	sys.AddFile(NewFileRef("Inc1", `Include\Inc1.axi`, FileTypeInclude))

	a.Rebuild()

	s.Len(a.files, 1)
	for k := range a.files {
		s.True(strings.HasSuffix(k, "Inc1.axi"))
	}
}

// ---------------------------------------------------------------------------
// compileTypeFor (table-driven)
// ---------------------------------------------------------------------------

func TestCompileTypeFor(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected FileCompileType
	}{
		{FileTypeSource, FileCompileTypeNetLinx},
		{FileTypeMasterSrc, FileCompileTypeNetLinx},
		{FileTypeInclude, FileCompileTypeNetLinx},
		{FileTypeModule, FileCompileTypeNetLinx},
		{FileTypeTP4, FileCompileTypeNone},
		{FileTypeTP5, FileCompileTypeNone},
		{FileTypeTPD, FileCompileTypeNone},
		{FileTypeIR, FileCompileTypeNone},
		{FileTypeDuet, FileCompileTypeNone},
		{FileTypeXDD, FileCompileTypeNone},
		{FileTypeAXB, FileCompileTypeNone},
		{FileTypeTKO, FileCompileTypeNone},
		{FileTypeIRDB, FileCompileTypeNone},
		{FileTypeOther, FileCompileTypeNone},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			assert.Equal(t, tt.expected, compileTypeFor(tt.fileType))
		})
	}
}

// ---------------------------------------------------------------------------
// Full build-and-write integration
// ---------------------------------------------------------------------------

func (s *MarshalTestSuite) TestBuildAndWrite_RoundTrip() {
	dir := s.T().TempDir()
	apwPath := filepath.Join(dir, "Built.apw")

	a := NewAPW("Built", apwPath)
	proj := a.Workspace().AddProject(NewProject("Built"))
	sys := proj.AddSystem(newDefaultSystem())
	sys.AddFile(NewFileRef("Inc1", "Include/Inc1.axi", FileTypeInclude))
	sys.AddFile(NewFileRef("Mod1", "Module/Mod1.axs", FileTypeModule))
	sys.AddFile(NewFileRef("Main", "Source/Main.axs", FileTypeMasterSrc))
	s.Require().NoError(a.Write(apwPath))

	data, err := os.ReadFile(apwPath)
	s.Require().NoError(err)

	parsed, err := Parse(apwPath, data)
	s.Require().NoError(err)

	ids := collectFileIDs(parsed.AllFiles())
	s.Contains(ids, "Inc1")
	s.Contains(ids, "Mod1")
	s.Contains(ids, "Main")

	s.Len(parsed.ModuleFiles(), 1)
	s.Len(parsed.IncludePath(), 1)
	s.Len(parsed.MasterSrcPath(), 1)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// minimalWorkspace returns a *Workspace that mirrors the structure of
// minimalAPW() used throughout the APW test suite.
func minimalWorkspace() *Workspace {
	ws := NewWorkspace("TestWorkspace")
	ws.CreateVersion = "4.0"
	ws.CurrentVersion = "4.0"

	proj := NewProject("TestProject")

	sys := NewSystem("TestSystem")
	sys.SysID = "1"
	sys.IsActive = "true"
	sys.Platform = "Netlinx"
	sys.Transport = "Serial"
	sys.TransportEx = "TCPIP"
	sys.Files = []*FileRef{
		{
			Identifier:   "TestInclude",
			FilePathName: "Include/TestInclude.axi",
			Type:         FileTypeInclude,
			CompileType:  FileCompileTypeNetLinx,
		},
		{
			Identifier:   "TestModule",
			FilePathName: "Module/TestModule.axs",
			Type:         FileTypeModule,
			CompileType:  FileCompileTypeNetLinx,
		},
		{
			Identifier:   "TestMain",
			FilePathName: "Source/TestMain.axs",
			Type:         FileTypeMasterSrc,
			CompileType:  FileCompileTypeNetLinx,
		},
	}

	proj.Systems = append(proj.Systems, sys)
	ws.Projects = append(ws.Projects, proj)

	return ws
}

// newDefaultSystem returns a *System with the required XML attributes set to
// their NetLinx Studio defaults, used in tests that exercise the hierarchy
// builder methods (AddFile, AddDeviceMap, AddIRDB).
func newDefaultSystem() *System {
	sys := NewSystem("TestSystem")
	sys.SysID = "1"
	sys.IsActive = "true"
	sys.Platform = "Netlinx"
	sys.Transport = "Serial"
	sys.TransportEx = "TCPIP"
	return sys
}

// collectFileIDs returns the set of file IDs from a slice of Files, excluding
// Workspace-type entries whose ID equals the workspace stem.
func collectFileIDs(files []File) map[string]bool {
	ids := make(map[string]bool)
	for _, f := range files {
		if f.Type != FileTypeWorkspace {
			ids[f.ID] = true
		}
	}
	return ids
}
