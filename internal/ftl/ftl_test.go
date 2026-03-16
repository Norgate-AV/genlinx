package ftl

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// buildAPW constructs a minimal *apw.APW programmatically so tests have no
// dependency on disk files.
func buildAPW(t *testing.T, apwPath string, ws *apw.Workspace) *apw.APW {
	t.Helper()
	a := apw.NewAPW(ws.Identifier, apwPath)
	a.SetWorkspace(ws)
	return a
}

// singleSystemAPW returns an APW with one project, one system, and the given
// file types as files within that system.
func singleSystemAPW(t *testing.T, apwPath, transTCPIPEx string, types ...apw.FileType) *apw.APW {
	t.Helper()

	ws := apw.NewWorkspace("TestWorkspace")
	p := apw.NewProject("TestProject")
	s := apw.NewSystem("TestSystem")
	s.TransTCPIPEx = transTCPIPEx

	for _, ft := range types {
		ext, ok := apw.AmxExtensions[ft]
		if !ok || ext == "" {
			ext = ".bin"
		}

		s.AddFile(apw.NewFileRef("TestFile-"+string(ft), "Source\\TestFile"+ext, ft))
	}

	p.AddSystem(s)
	ws.AddProject(p)

	return buildAPW(t, apwPath, ws)
}

// ---------------------------------------------------------------------------
// FromAPW
// ---------------------------------------------------------------------------

type FromAPWSuite struct{ suite.Suite }

func TestFromAPWSuite(t *testing.T) { suite.Run(t, new(FromAPWSuite)) }

func (s *FromAPWSuite) TestEmptyWorkspaceReturnsEmptyList() {
	ws := apw.NewWorkspace("Empty")
	a := buildAPW(s.T(), "/workspace/Empty.apw", ws)

	list := FromAPW(a)
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWSuite) TestMasterSrcProducesOneTKNItem() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeMasterSrc,
	)

	list := FromAPW(a)
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".tkn", strings.ToLower(filepath.Ext(list.Items[0].SourceFile)))
}

func (s *FromAPWSuite) TestTP4ProducesOneItem() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeTP4,
	)

	list := FromAPW(a)
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".tp4", strings.ToLower(filepath.Ext(list.Items[0].SourceFile)))
}

func (s *FromAPWSuite) TestTP5ProducesOneItem() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeTP5,
	)

	list := FromAPW(a)
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".tp5", strings.ToLower(filepath.Ext(list.Items[0].SourceFile)))
}

func (s *FromAPWSuite) TestKPBProducesOneItem() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeKPB,
	)

	list := FromAPW(a)
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".kpb", strings.ToLower(filepath.Ext(list.Items[0].SourceFile)))
}

func (s *FromAPWSuite) TestIncludeModuleSourceSkipped() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeInclude,
		apw.FileTypeModule,
		apw.FileTypeSource,
	)

	list := FromAPW(a)
	assert.Empty(s.T(), list.Items, "Include/Module/Source files must not appear in FTL")
}

func (s *FromAPWSuite) TestMixedFilesOnlyTransferable() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeInclude,
		apw.FileTypeMasterSrc,
		apw.FileTypeModule,
		apw.FileTypeTP4,
	)

	list := FromAPW(a)
	s.Require().Len(list.Items, 2, "only MasterSrc and TP4 should appear")
}

func (s *FromAPWSuite) TestAllProjectsIncludedByDefault() {
	ws := apw.NewWorkspace("Multi")

	for _, pid := range []string{"ProjectA", "ProjectB"} {
		p := apw.NewProject(pid)
		sys := apw.NewSystem("Sys")
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		sys.AddFile(apw.NewFileRef("Main", "Source\\Main.axs", apw.FileTypeMasterSrc))
		p.AddSystem(sys)
		ws.AddProject(p)
	}

	a := buildAPW(s.T(), "/ws/Multi.apw", ws)
	list := FromAPW(a)
	assert.Len(s.T(), list.Items, 2, "one TKN item per project")
}

// ---------------------------------------------------------------------------
// FromAPWProject
// ---------------------------------------------------------------------------

type FromAPWProjectSuite struct{ suite.Suite }

func TestFromAPWProjectSuite(t *testing.T) { suite.Run(t, new(FromAPWProjectSuite)) }

func (s *FromAPWProjectSuite) TestKnownProjectReturnsItems() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("ProjectA")
	sys := apw.NewSystem("Sys")
	sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
	sys.AddFile(apw.NewFileRef("Main", "Source\\Main.axs", apw.FileTypeMasterSrc))
	p.AddSystem(sys)
	ws.AddProject(p)

	p2 := apw.NewProject("ProjectB")
	sys2 := apw.NewSystem("SysB")
	sys2.TransTCPIPEx = "10.0.0.2|1319|1|Desc||"
	sys2.AddFile(apw.NewFileRef("MainB", "Source\\MainB.axs", apw.FileTypeMasterSrc))
	p2.AddSystem(sys2)
	ws.AddProject(p2)

	a := buildAPW(s.T(), "/ws/WS.apw", ws)

	list := FromAPWProject(a, "ProjectA")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "ProjectA", list.Items[0].ProjectName)
}

func (s *FromAPWProjectSuite) TestUnknownProjectReturnsEmpty() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeMasterSrc,
	)

	list := FromAPWProject(a, "NoSuchProject")
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWProjectSuite) TestOtherProjectsExcluded() {
	ws := apw.NewWorkspace("WS")

	for _, pid := range []string{"Alpha", "Beta", "Gamma"} {
		p := apw.NewProject(pid)
		sys := apw.NewSystem("Sys")
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		sys.AddFile(apw.NewFileRef("Main", "Source\\Main.axs", apw.FileTypeMasterSrc))
		p.AddSystem(sys)
		ws.AddProject(p)
	}

	a := buildAPW(s.T(), "/ws/WS.apw", ws)

	list := FromAPWProject(a, "Beta")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "Beta", list.Items[0].ProjectName)
}

// ---------------------------------------------------------------------------
// FromAPWSystem
// ---------------------------------------------------------------------------

type FromAPWSystemSuite struct{ suite.Suite }

func TestFromAPWSystemSuite(t *testing.T) { suite.Run(t, new(FromAPWSystemSuite)) }

func (s *FromAPWSystemSuite) TestKnownSystemReturnsItems() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")

	for _, sid := range []string{"Room1", "Room2"} {
		sys := apw.NewSystem(sid)
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		sys.AddFile(apw.NewFileRef("Main", "Source\\Main.axs", apw.FileTypeMasterSrc))
		p.AddSystem(sys)
	}

	ws.AddProject(p)
	a := buildAPW(s.T(), "/ws/WS.apw", ws)

	list := FromAPWSystem(a, "Proj", "Room1")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "Room1", list.Items[0].SystemName)
}

func (s *FromAPWSystemSuite) TestUnknownSystemReturnsEmpty() {
	a := singleSystemAPW(s.T(), "/ws/Test.apw", "10.0.0.1|1319|1|Desc||",
		apw.FileTypeMasterSrc,
	)

	list := FromAPWSystem(a, "TestProject", "NoSuchSystem")
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWSystemSuite) TestOtherSystemsExcluded() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")

	for _, sid := range []string{"A", "B", "C"} {
		sys := apw.NewSystem(sid)
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		sys.AddFile(apw.NewFileRef("Main", "Source\\Main.axs", apw.FileTypeMasterSrc))
		p.AddSystem(sys)
	}

	ws.AddProject(p)
	a := buildAPW(s.T(), "/ws/WS.apw", ws)

	list := FromAPWSystem(a, "Proj", "B")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "B", list.Items[0].SystemName)
}

// ---------------------------------------------------------------------------
// Item field values
// ---------------------------------------------------------------------------

type ItemFieldsSuite struct{ suite.Suite }

func TestItemFieldsSuite(t *testing.T) { suite.Run(t, new(ItemFieldsSuite)) }

func (s *ItemFieldsSuite) item(ft apw.FileType) Item {
	s.T().Helper()
	a := singleSystemAPW(s.T(), "/ws/MyWorkspace.apw", "10.0.0.5|1319|1|Lab||", ft)
	list := FromAPW(a)
	s.Require().Len(list.Items, 1)
	return list.Items[0]
}

// --- Common fields shared by all item types ---

func (s *ItemFieldsSuite) TestPlatformIsNetLinx() {
	assert.Equal(s.T(), platformNetLinx, s.item(apw.FileTypeMasterSrc).Platform)
}

func (s *ItemFieldsSuite) TestFtTypeIsOne() {
	assert.Equal(s.T(), ftType, s.item(apw.FileTypeMasterSrc).FtType)
}

func (s *ItemFieldsSuite) TestWorkspaceNameCarried() {
	assert.Equal(s.T(), "TestWorkspace", s.item(apw.FileTypeMasterSrc).WorkspaceName)
}

func (s *ItemFieldsSuite) TestProjectNameCarried() {
	assert.Equal(s.T(), "TestProject", s.item(apw.FileTypeMasterSrc).ProjectName)
}

func (s *ItemFieldsSuite) TestSystemNameCarried() {
	assert.Equal(s.T(), "TestSystem", s.item(apw.FileTypeMasterSrc).SystemName)
}

func (s *ItemFieldsSuite) TestWorkspacePathNameCarried() {
	path := s.item(apw.FileTypeMasterSrc).WorkspacePathName
	assert.True(s.T(),
		strings.HasSuffix(filepath.ToSlash(path), "ws/MyWorkspace.apw"),
		"unexpected WorkspacePathName: %s", path,
	)
}

func (s *ItemFieldsSuite) TestPortIsOne() {
	assert.Equal(s.T(), 1, s.item(apw.FileTypeMasterSrc).Port)
}

func (s *ItemFieldsSuite) TestSystemIsZero() {
	assert.Equal(s.T(), 0, s.item(apw.FileTypeMasterSrc).System)
}

func (s *ItemFieldsSuite) TestIsDirectToDeviceIsZero() {
	assert.Equal(s.T(), 0, s.item(apw.FileTypeMasterSrc).IsDirectToDevice)
}

func (s *ItemFieldsSuite) TestTpdFlagsEnabled() {
	item := s.item(apw.FileTypeMasterSrc)
	assert.Equal(s.T(), 1, item.IsTpdSendBitmapsEnabled)
	assert.Equal(s.T(), 1, item.IsTpdSendFontsEnabled)
	assert.Equal(s.T(), 1, item.IsTpdSendIconsEnabled)
}

func (s *ItemFieldsSuite) TestItemCheckedInListIsOne() {
	assert.Equal(s.T(), 1, s.item(apw.FileTypeMasterSrc).ItemCheckedInList)
}

func (s *ItemFieldsSuite) TestTransportSettingsPrefixIsT() {
	assert.True(s.T(),
		strings.HasPrefix(s.item(apw.FileTypeMasterSrc).TransportSettings, "T-"),
		"TransportSettings must start with 'T-'",
	)
}

func (s *ItemFieldsSuite) TestTransportSettingsContainsHost() {
	assert.Contains(s.T(), s.item(apw.FileTypeMasterSrc).TransportSettings, "10.0.0.5")
}

// --- MasterSrc-specific ---

func (s *ItemFieldsSuite) TestMasterSrcDeviceIsZero() {
	assert.Equal(s.T(), deviceMaster, s.item(apw.FileTypeMasterSrc).Device)
}

func (s *ItemFieldsSuite) TestMasterSrcIsRebootRequired() {
	assert.Equal(s.T(), 1, s.item(apw.FileTypeMasterSrc).IsRebootRequired)
}

func (s *ItemFieldsSuite) TestMasterSrcSourceFileIsTKN() {
	assert.Equal(s.T(), apw.FileExtensionTKN,
		strings.ToLower(filepath.Ext(s.item(apw.FileTypeMasterSrc).SourceFile)),
	)
}

func (s *ItemFieldsSuite) TestMasterSrcNoSmartTransferFlags() {
	item := s.item(apw.FileTypeMasterSrc)
	assert.Equal(s.T(), 0, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 0, item.IsTP5SmartTransferEnabled)
}

// --- TP4-specific ---

func (s *ItemFieldsSuite) TestTP4DeviceIsPanel() {
	assert.Equal(s.T(), devicePanel, s.item(apw.FileTypeTP4).Device)
}

func (s *ItemFieldsSuite) TestTP4IsRebootNotRequired() {
	assert.Equal(s.T(), 0, s.item(apw.FileTypeTP4).IsRebootRequired)
}

func (s *ItemFieldsSuite) TestTP4SmartTransferEnabled() {
	item := s.item(apw.FileTypeTP4)
	assert.Equal(s.T(), 1, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 0, item.IsTP5SmartTransferEnabled)
}

// --- TP5-specific ---

func (s *ItemFieldsSuite) TestTP5DeviceIsPanel() {
	assert.Equal(s.T(), devicePanel, s.item(apw.FileTypeTP5).Device)
}

func (s *ItemFieldsSuite) TestTP5IsRebootNotRequired() {
	assert.Equal(s.T(), 0, s.item(apw.FileTypeTP5).IsRebootRequired)
}

func (s *ItemFieldsSuite) TestTP5SmartTransferEnabled() {
	item := s.item(apw.FileTypeTP5)
	assert.Equal(s.T(), 0, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 1, item.IsTP5SmartTransferEnabled)
}

// --- KPB-specific ---

func (s *ItemFieldsSuite) TestKPBDeviceIsPanel() {
	assert.Equal(s.T(), devicePanel, s.item(apw.FileTypeKPB).Device)
}

func (s *ItemFieldsSuite) TestKPBIsRebootNotRequired() {
	assert.Equal(s.T(), 0, s.item(apw.FileTypeKPB).IsRebootRequired)
}

func (s *ItemFieldsSuite) TestKPBGlyphAndFontFlagsEnabled() {
	item := s.item(apw.FileTypeKPB)
	assert.Equal(s.T(), 1, item.IsKPBSendGlyphEnabled)
	assert.Equal(s.T(), 1, item.IsKPBSendFontEnabled)
}

func (s *ItemFieldsSuite) TestKPBSmartTransferFlagsZero() {
	item := s.item(apw.FileTypeKPB)
	assert.Equal(s.T(), 0, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 0, item.IsTP5SmartTransferEnabled)
}

// ---------------------------------------------------------------------------
// Source file path derivation
// ---------------------------------------------------------------------------

type SourcePathSuite struct{ suite.Suite }

func TestSourcePathSuite(t *testing.T) { suite.Run(t, new(SourcePathSuite)) }

func (s *SourcePathSuite) TestMasterSrcTKNReplacesAXSExtension() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")
	sys := apw.NewSystem("Sys")
	sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
	sys.AddFile(apw.NewFileRef("Main", `Source\MyController.axs`, apw.FileTypeMasterSrc))
	p.AddSystem(sys)
	ws.AddProject(p)

	a := buildAPW(s.T(), `/work/WS.apw`, ws)
	list := FromAPW(a)
	s.Require().Len(list.Items, 1)

	src := list.Items[0].SourceFile
	assert.Equal(s.T(), ".tkn", strings.ToLower(filepath.Ext(src)))
	assert.True(s.T(), strings.HasSuffix(src, "MyController.tkn"),
		"expected filename MyController.tkn, got %s", src)
}

func (s *SourcePathSuite) TestTP4SourceFilePathUnchanged() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")
	sys := apw.NewSystem("Sys")
	sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
	sys.AddFile(apw.NewFileRef("Panel", `User Interface\MyPanel.TP4`, apw.FileTypeTP4))
	p.AddSystem(sys)
	ws.AddProject(p)

	a := buildAPW(s.T(), `/work/WS.apw`, ws)
	list := FromAPW(a)
	s.Require().Len(list.Items, 1)

	src := list.Items[0].SourceFile
	assert.Equal(s.T(), ".TP4", filepath.Ext(src))
	assert.True(s.T(), strings.HasSuffix(src, "MyPanel.TP4"),
		"expected filename MyPanel.TP4, got %s", src)
}

func (s *SourcePathSuite) TestSourceFileIsAbsolutePath() {
	a := singleSystemAPW(s.T(), `/work/WS.apw`, "10.0.0.1|1319|1|Desc||",
		apw.FileTypeMasterSrc,
	)

	list := FromAPW(a)
	s.Require().Len(list.Items, 1)

	assert.True(s.T(), filepath.IsAbs(list.Items[0].SourceFile),
		"SourceFile must be an absolute path")
}

// ---------------------------------------------------------------------------
// Multi-system workspace: per-system TransportSettings
// ---------------------------------------------------------------------------

type TransportSettingsSuite struct{ suite.Suite }

func TestTransportSettingsSuite(t *testing.T) { suite.Run(t, new(TransportSettingsSuite)) }

func (s *TransportSettingsSuite) TestEachSystemUsesItsOwnTransport() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")

	sysA := apw.NewSystem("SysA")
	sysA.TransTCPIPEx = "192.168.1.1|1319|1|Room A||"
	sysA.AddFile(apw.NewFileRef("MainA", `Source\MainA.axs`, apw.FileTypeMasterSrc))
	p.AddSystem(sysA)

	sysB := apw.NewSystem("SysB")
	sysB.TransTCPIPEx = "192.168.1.2|1319|1|Room B||"
	sysB.AddFile(apw.NewFileRef("MainB", `Source\MainB.axs`, apw.FileTypeMasterSrc))
	p.AddSystem(sysB)

	ws.AddProject(p)
	a := buildAPW(s.T(), "/ws/WS.apw", ws)

	list := FromAPW(a)
	s.Require().Len(list.Items, 2)

	tsA := list.Items[0].TransportSettings
	tsB := list.Items[1].TransportSettings
	assert.NotEqual(s.T(), tsA, tsB, "each system must have its own transport settings")
	assert.Contains(s.T(), tsA, "192.168.1.1")
	assert.Contains(s.T(), tsB, "192.168.1.2")
}

func (s *TransportSettingsSuite) TestSystemWithMultipleFilesShareTransport() {
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")
	sys := apw.NewSystem("Sys")
	sys.TransTCPIPEx = "10.0.0.99|1319|1|Lab||"
	sys.AddFile(apw.NewFileRef("Main", `Source\Main.axs`, apw.FileTypeMasterSrc))
	sys.AddFile(apw.NewFileRef("Panel", `UI\Panel.TP4`, apw.FileTypeTP4))
	p.AddSystem(sys)
	ws.AddProject(p)

	a := buildAPW(s.T(), "/ws/WS.apw", ws)
	list := FromAPW(a)
	s.Require().Len(list.Items, 2)

	assert.Equal(s.T(), list.Items[0].TransportSettings, list.Items[1].TransportSettings)
	assert.Contains(s.T(), list.Items[0].TransportSettings, "10.0.0.99")
}
