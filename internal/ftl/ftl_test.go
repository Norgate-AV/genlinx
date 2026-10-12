package ftl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

// ---------------------------------------------------------------------------
// Test helper: sysHelper
//
// Builds a single-system APW inside a t.TempDir() so tests that need real
// files on disk (TKN existence checks, DEFINE_DEVICE scanning) work correctly.
// ---------------------------------------------------------------------------

type sysHelper struct {
	t          *testing.T
	dir        string
	ws         *apw.Workspace
	proj       *apw.Project
	sys        *apw.System
	sourceName string // base name of the MasterSrc file, without extension
}

func newSysHelper(t *testing.T, sysID, transTCPIPEx string) *sysHelper {
	t.Helper()
	dir := t.TempDir()
	ws := apw.NewWorkspace("TestWorkspace")
	p := apw.NewProject("TestProject")
	s := apw.NewSystem("TestSystem")
	s.SysID = sysID
	s.TransTCPIPEx = transTCPIPEx
	p.AddSystem(s)
	ws.AddProject(p)

	return &sysHelper{t: t, dir: dir, ws: ws, proj: p, sys: s}
}

// withMasterSrc adds a MasterSrc FileRef.  Call withTKN() afterwards to also
// create the compiled .tkn file on disk.
func (h *sysHelper) withMasterSrc(name string) *sysHelper {
	h.t.Helper()
	h.sourceName = name
	h.sys.AddFile(apw.NewFileRef(name, `Source\`+name+`.axs`, apw.FileTypeMasterSrc))
	return h
}

// withTKN creates the compiled .tkn file so makeItems does not skip the entry.
func (h *sysHelper) withTKN() *sysHelper {
	h.t.Helper()
	dir := filepath.Join(h.dir, "Source")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		h.t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, h.sourceName+".tkn"), []byte{}, 0o644); err != nil {
		h.t.Fatal(err)
	}

	return h
}

// withDefineDevice writes the MasterSrc .axs file with the supplied body so
// that buildDeviceRegistry can scan it for device definitions.
func (h *sysHelper) withDefineDevice(axsBody string) *sysHelper {
	h.t.Helper()
	dir := filepath.Join(h.dir, "Source")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		h.t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, h.sourceName+".axs"), []byte(axsBody), 0o644); err != nil {
		h.t.Fatal(err)
	}

	return h
}

// withPanel adds a panel FileRef with the given type, relative path under the
// APW directory, and one DeviceMap per supplied devAddr string.
func (h *sysHelper) withPanel(fileType apw.FileType, relPath string, devAddrs ...string) *sysHelper {
	h.t.Helper()
	fr := apw.NewFileRef(filepath.Base(relPath), relPath, fileType)

	for _, addr := range devAddrs {
		fr.AddDeviceMap(apw.NewDeviceMap(addr, addr))
	}

	h.sys.AddFile(fr)
	return h
}

// withInclude adds a plain Include FileRef (non-transferable).
func (h *sysHelper) withInclude(name string) *sysHelper {
	h.sys.AddFile(apw.NewFileRef(name, `Include\`+name+`.axi`, apw.FileTypeInclude))
	return h
}

// build finalises the APW.
func (h *sysHelper) build() *apw.APW {
	h.t.Helper()
	apwPath := filepath.Join(h.dir, "TestWorkspace.apw")
	a := apw.NewAPW(h.ws.Identifier, apwPath)
	a.SetWorkspace(h.ws)
	return a
}

// ---------------------------------------------------------------------------
// DevResolverSuite — unit tests for devresolver.go
// ---------------------------------------------------------------------------

type DevResolverSuite struct{ suite.Suite }

func TestDevResolverSuite(t *testing.T) { suite.Run(t, new(DevResolverSuite)) }

func (s *DevResolverSuite) TestCustomAddr_Parsed() {
	dps, ok := resolveDevAddr("Custom [10001:1:2]", nil)
	s.Require().True(ok)
	assert.Equal(s.T(), 10001, dps.Device)
	assert.Equal(s.T(), 1, dps.Port)
	assert.Equal(s.T(), 2, dps.System)
}

func (s *DevResolverSuite) TestCustomAddr_SystemZero() {
	dps, ok := resolveDevAddr("Custom [10001:1:0]", nil)
	s.Require().True(ok)
	assert.Equal(s.T(), 0, dps.System)
}

func (s *DevResolverSuite) TestCustomAddr_WhitespaceVariants() {
	_, ok1 := resolveDevAddr("Custom [10001:1:0]", nil)
	_, ok2 := resolveDevAddr("Custom  [10001:1:0]", nil) // extra space
	s.Require().True(ok1)
	s.Require().True(ok2)
}

func (s *DevResolverSuite) TestSymbolicAddr_FoundInRegistry() {
	reg := deviceRegistry{"dvTP_Main": {Device: 10001, Port: 1, System: 1}}
	dps, ok := resolveDevAddr("dvTP_Main", reg)
	s.Require().True(ok)
	assert.Equal(s.T(), DPS{10001, 1, 1}, dps)
}

func (s *DevResolverSuite) TestSymbolicAddr_NotFoundInRegistry() {
	reg := deviceRegistry{}
	_, ok := resolveDevAddr("dvTP_Missing", reg)
	assert.False(s.T(), ok)
}

func (s *DevResolverSuite) TestSymbolicAddr_NilRegistry() {
	_, ok := resolveDevAddr("dvTP_Main", nil)
	assert.False(s.T(), ok)
}

func (s *DevResolverSuite) TestBuildRegistry_ParsesDefineDevice() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice(`
PROGRAM_NAME='Test'

DEFINE_DEVICE

dvTP_Main        = 10001:1:0
dvTP_Lectern     = 10001:1:1
dvTP_Control     = 10002:1:2

DEFINE_CONSTANT
`)

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)

	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvTP_Main"])
	assert.Equal(s.T(), DPS{10001, 1, 1}, reg["dvTP_Lectern"])
	assert.Equal(s.T(), DPS{10002, 1, 2}, reg["dvTP_Control"])
}

func (s *DevResolverSuite) TestBuildRegistry_IgnoresNonDeviceLines() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice(`
DEFINE_DEVICE

dvTP_Main = 10001:1:0

DEFINE_CONSTANT

NOT_A_DEV = 99:9:9  // should be ignored — outside DEFINE_DEVICE
`)

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)

	assert.Len(s.T(), reg, 1)
	_, hasNotADev := reg["NOT_A_DEV"]
	assert.False(s.T(), hasNotADev)
}

func (s *DevResolverSuite) TestBuildRegistry_LineCommentsStripped() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice(`
DEFINE_DEVICE

dvTP_Main = 10001:1:0    // main touchpanel

DEFINE_CONSTANT
`)

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvTP_Main"])
}

func (s *DevResolverSuite) TestBuildRegistry_BlockCommentsStripped() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice(`
DEFINE_DEVICE

(* this device is commented out
dvTP_Hidden = 99:9:9
*)
dvTP_Main = 10001:1:0

DEFINE_CONSTANT
`)

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Len(s.T(), reg, 1)
	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvTP_Main"])
}

func (s *DevResolverSuite) TestBuildRegistry_ScansIncludeFile() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice(`
PROGRAM_NAME='Test'
#INCLUDE 'Devices.axi'
DEFINE_CONSTANT
`)

	// Write the included file with DEFINE_DEVICE.
	inclDir := filepath.Join(h.dir, "Include")
	s.Require().NoError(os.MkdirAll(inclDir, 0o755))
	s.Require().NoError(os.WriteFile(
		filepath.Join(inclDir, "Devices.axi"),
		[]byte("DEFINE_DEVICE\ndvTP_Main = 10001:1:0\nDEFINE_CONSTANT\n"),
		0o644,
	))

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvTP_Main"])
}

func (s *DevResolverSuite) TestBuildRegistry_NoMasterSrc_EmptyRegistry() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withInclude("SomeHeader")
	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Empty(s.T(), reg)
}

func (s *DevResolverSuite) TestStripBlockComments_ParensAsterisk() {
	src := "before (* removed *) after"
	assert.Equal(s.T(), "before  after", stripNetLinxBlockComments(src))
}

func (s *DevResolverSuite) TestStripBlockComments_CStyle() {
	src := "before /* removed */ after"
	assert.Equal(s.T(), "before  after", stripNetLinxBlockComments(src))
}

func (s *DevResolverSuite) TestStripBlockComments_PreservesNewlines() {
	src := "a\n(* line1\nline2 *)\nb"
	result := stripNetLinxBlockComments(src)
	assert.Equal(s.T(), 3, strings.Count(result, "\n"), "newlines inside block comment must be preserved")
}

func (s *DevResolverSuite) TestParseSysID_Valid() {
	assert.Equal(s.T(), 2, parseSysID("2"))
	assert.Equal(s.T(), 11, parseSysID("11"))
}

func (s *DevResolverSuite) TestParseSysID_Zero() {
	assert.Equal(s.T(), 0, parseSysID("0"))
}

func (s *DevResolverSuite) TestParseSysID_EmptyDefaultsToZero() {
	assert.Equal(s.T(), 0, parseSysID(""))
}

func (s *DevResolverSuite) TestParseSysID_InvalidDefaultsToZero() {
	assert.Equal(s.T(), 0, parseSysID("N/A"))
}

func (s *DevResolverSuite) TestBuildRegistry_ScansSourceSubdirInclude() {
	// resolveIncludePath falls back to Source/ when the file isn't found in
	// the APW root or Include/.
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice("PROGRAM_NAME='Test'\n#INCLUDE 'Devices.axi'\nDEFINE_CONSTANT\n")

	// Place the included file in Source/ (not in Include/).
	s.Require().NoError(os.WriteFile(
		filepath.Join(h.dir, "Source", "Devices.axi"),
		[]byte("DEFINE_DEVICE\ndvTP_Main = 10001:1:0\nDEFINE_CONSTANT\n"),
		0o644,
	))

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvTP_Main"])
}

func (s *DevResolverSuite) TestBuildRegistry_FirstDefinitionWins() {
	// If the same device name appears in two separate DEFINE_DEVICE blocks
	// (or across the main file and an include), the first occurrence is kept.
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice("DEFINE_DEVICE\ndvTP_Main = 10001:1:0\nDEFINE_CONSTANT\nDEFINE_DEVICE\ndvTP_Main = 99:9:9\nDEFINE_CONSTANT\n")

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvTP_Main"], "first definition must be kept")
}

func (s *DevResolverSuite) TestBuildRegistry_CircularIncludeNoHang() {
	// A #INCLUDEs B; B #INCLUDEs A.  The visited-map guard must prevent
	// infinite recursion and the scanner must still find both devices.
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice("PROGRAM_NAME='Test'\n#INCLUDE 'A.axi'\nDEFINE_CONSTANT\n")

	inclDir := filepath.Join(h.dir, "Include")
	s.Require().NoError(os.MkdirAll(inclDir, 0o755))
	s.Require().NoError(os.WriteFile(
		filepath.Join(inclDir, "A.axi"),
		[]byte("#INCLUDE 'B.axi'\nDEFINE_DEVICE\ndvA = 10001:1:0\nDEFINE_CONSTANT\n"),
		0o644,
	))
	s.Require().NoError(os.WriteFile(
		filepath.Join(inclDir, "B.axi"),
		[]byte("#INCLUDE 'A.axi'\nDEFINE_DEVICE\ndvB = 10002:1:0\nDEFINE_CONSTANT\n"),
		0o644,
	))

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Equal(s.T(), DPS{10001, 1, 0}, reg["dvA"])
	assert.Equal(s.T(), DPS{10002, 1, 0}, reg["dvB"])
}

func (s *DevResolverSuite) TestBuildRegistry_MaxIncludeDepthNotExceeded() {
	// Build a chain: MasterSrc(depth 0) → inc0(1) → … → inc9(10) → inc10(11).
	// inc9 (depth 10, allowed) defines dvAtDepth10.
	// inc10 (depth 11, rejected) defines dvTooDeep — must not appear.
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Test||").
		withMasterSrc("Main").
		withDefineDevice("#INCLUDE 'inc0.axi'\n")

	inclDir := filepath.Join(h.dir, "Include")
	s.Require().NoError(os.MkdirAll(inclDir, 0o755))

	for i := 0; i < 10; i++ {
		var content string
		if i == 9 {
			// inc9 is scanned at depth 10 — device must be in the registry.
			content = fmt.Sprintf("DEFINE_DEVICE\ndvAtDepth10 = 9999:1:0\nDEFINE_CONSTANT\n#INCLUDE 'inc%d.axi'\n", i+1)
		} else {
			content = fmt.Sprintf("#INCLUDE 'inc%d.axi'\n", i+1)
		}
		s.Require().NoError(os.WriteFile(
			filepath.Join(inclDir, fmt.Sprintf("inc%d.axi", i)),
			[]byte(content),
			0o644,
		))
	}

	// inc10 is called at depth 11 — must be skipped entirely.
	s.Require().NoError(os.WriteFile(
		filepath.Join(inclDir, "inc10.axi"),
		[]byte("DEFINE_DEVICE\ndvTooDeep = 10001:1:0\nDEFINE_CONSTANT\n"),
		0o644,
	))

	a := h.build()
	reg := buildDeviceRegistry(a.Workspace().Projects[0].Systems[0], h.dir)
	assert.Equal(s.T(), DPS{9999, 1, 0}, reg["dvAtDepth10"], "device at depth 10 must be found")
	_, tooDeep := reg["dvTooDeep"]
	assert.False(s.T(), tooDeep, "device at depth 11 must be skipped")
}

// ---------------------------------------------------------------------------
// FromAPWSuite
// ---------------------------------------------------------------------------

type FromAPWSuite struct{ suite.Suite }

func TestFromAPWSuite(t *testing.T) { suite.Run(t, new(FromAPWSuite)) }

func (s *FromAPWSuite) TestEmptyWorkspaceReturnsEmptyList() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||")
	list := FromAPW(h.build())
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWSuite) TestMasterSrcWithTKNOnDisk_ProducesItem() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withMasterSrc("Main").withTKN()
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".tkn", strings.ToLower(filepath.Ext(list.Items[0].SourceFile)))
}

func (s *FromAPWSuite) TestMasterSrcWithoutTKNOnDisk_Skipped() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withMasterSrc("Main") // no withTKN()
	list := FromAPW(h.build())
	assert.Empty(s.T(), list.Items)
	if s.Len(list.Warnings, 1) {
		assert.Contains(s.T(), list.Warnings[0], "TKN not found")
		assert.Contains(s.T(), list.Warnings[0], "[TestSystem]")
	}
}

func (s *FromAPWSuite) TestMasterSrcWithoutTKNOnDisk_WarnContainsTKNPath() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withMasterSrc("Main") // no withTKN()
	list := FromAPW(h.build())
	s.Require().Len(list.Warnings, 1)
	assert.Contains(s.T(), list.Warnings[0], "Main.tkn")
}

func (s *FromAPWSuite) TestTP4_WithCustomDevAddr_ProducesItem() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withPanel(apw.FileTypeTP4, `User Interface\Panel.TP4`, "Custom [10001:1:0]")
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".TP4", filepath.Ext(list.Items[0].SourceFile))
}

func (s *FromAPWSuite) TestTP4_WithSymbolicDevAddr_Resolved() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withMasterSrc("Main").
		withDefineDevice("DEFINE_DEVICE\ndvTP_Main = 10001:1:0\nDEFINE_CONSTANT\n").
		withPanel(apw.FileTypeTP4, `User Interface\Panel.TP4`, "dvTP_Main")
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), 10001, list.Items[0].Device)
}

func (s *FromAPWSuite) TestTP4_WithUnresolvableSymbolicDevAddr_Skipped() {
	// No MasterSrc so the registry stays empty — dvTP_Missing cannot resolve.
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withPanel(apw.FileTypeTP4, `User Interface\Panel.TP4`, "dvTP_Missing")
	list := FromAPW(h.build())
	assert.Empty(s.T(), list.Items)
	if s.Len(list.Warnings, 1) {
		assert.Contains(s.T(), list.Warnings[0], "unresolved DevAddr")
		assert.Contains(s.T(), list.Warnings[0], "dvTP_Missing")
		assert.Contains(s.T(), list.Warnings[0], "Panel.TP4")
		assert.Contains(s.T(), list.Warnings[0], "[TestSystem]")
	}
}

func (s *FromAPWSuite) TestTP4_WithNoDeviceMap_Skipped() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||")
	fr := apw.NewFileRef("Panel", `User Interface\Panel.TP4`, apw.FileTypeTP4)
	h.sys.AddFile(fr)
	list := FromAPW(h.build())
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWSuite) TestTP4_TwoDeviceMaps_TwoItems() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withPanel(apw.FileTypeTP4, `User Interface\Panel.TP4`,
			"Custom [10001:1:0]",
			"Custom [10002:1:0]",
		)
	list := FromAPW(h.build())
	assert.Len(s.T(), list.Items, 2)
}

func (s *FromAPWSuite) TestTP5_WithCustomDevAddr_ProducesItem() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withPanel(apw.FileTypeTP5, `User Interface\Panel.TP5`, "Custom [10001:1:0]")
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), ".TP5", filepath.Ext(list.Items[0].SourceFile))
}

func (s *FromAPWSuite) TestKPB_WithCustomDevAddr_ProducesItem() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").
		withPanel(apw.FileTypeKPB, `User Interface\Panel.KPB`, "Custom [10001:1:0]")
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
}

func (s *FromAPWSuite) TestIncludeModuleSourceAlwaysSkipped() {
	h := newSysHelper(s.T(), "0", "10.0.0.1|1319|1|Desc||").withInclude("Header")
	h.sys.AddFile(apw.NewFileRef("Mod", `Module\Mod.axs`, apw.FileTypeModule))
	h.sys.AddFile(apw.NewFileRef("Src", `Source\Other.axs`, apw.FileTypeSource))
	list := FromAPW(h.build())
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWSuite) TestAllProjectsIncludedByDefault() {
	dir := s.T().TempDir()
	ws := apw.NewWorkspace("Multi")

	for i, pid := range []string{"ProjectA", "ProjectB"} {
		p := apw.NewProject(pid)
		sys := apw.NewSystem("Sys")
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		name := pid + "Main"
		sys.AddFile(apw.NewFileRef(name, `Source\`+name+`.axs`, apw.FileTypeMasterSrc))
		p.AddSystem(sys)
		ws.AddProject(p)
		// Create the TKN file.
		tknDir := filepath.Join(dir, "Source")
		_ = os.MkdirAll(tknDir, 0o755)
		_ = os.WriteFile(filepath.Join(tknDir, name+".tkn"), []byte{}, 0o644)
		_ = i
	}

	a := apw.NewAPW(ws.Identifier, filepath.Join(dir, "Multi.apw"))
	a.SetWorkspace(ws)
	list := FromAPW(a)
	assert.Len(s.T(), list.Items, 2)
}

// ---------------------------------------------------------------------------
// FromAPWProjectSuite
// ---------------------------------------------------------------------------

type FromAPWProjectSuite struct{ suite.Suite }

func TestFromAPWProjectSuite(t *testing.T) { suite.Run(t, new(FromAPWProjectSuite)) }

// filteredAPW builds a workspace with two projects each containing a panel
// with a Custom DevAddr so no disk access is needed.
func filteredAPW(t *testing.T) *apw.APW {
	t.Helper()
	dir := t.TempDir()
	ws := apw.NewWorkspace("WS")

	for i, pid := range []string{"Alpha", "Beta"} {
		p := apw.NewProject(pid)
		sys := apw.NewSystem("Sys")
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		fr := apw.NewFileRef("Panel", `UI\Panel.TP4`, apw.FileTypeTP4)
		fr.AddDeviceMap(apw.NewDeviceMap("Custom [10001:1:0]", "Custom [10001:1:0]"))
		sys.AddFile(fr)
		p.AddSystem(sys)
		ws.AddProject(p)
		_ = i
	}

	a := apw.NewAPW(ws.Identifier, filepath.Join(dir, "WS.apw"))
	a.SetWorkspace(ws)
	return a
}

func (s *FromAPWProjectSuite) TestKnownProjectReturnsOnlyItsItems() {
	a := filteredAPW(s.T())
	list := FromAPWProject(a, "Alpha")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "Alpha", list.Items[0].ProjectName)
}

func (s *FromAPWProjectSuite) TestUnknownProjectReturnsEmpty() {
	list := FromAPWProject(filteredAPW(s.T()), "NoSuch")
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWProjectSuite) TestOtherProjectsExcluded() {
	// Alpha and Beta both exist; requesting Beta should only return Beta.
	list := FromAPWProject(filteredAPW(s.T()), "Beta")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "Beta", list.Items[0].ProjectName)
}

// ---------------------------------------------------------------------------
// FromAPWSystemSuite
// ---------------------------------------------------------------------------

type FromAPWSystemSuite struct{ suite.Suite }

func TestFromAPWSystemSuite(t *testing.T) { suite.Run(t, new(FromAPWSystemSuite)) }

func multiSystemAPW(t *testing.T, sysIDs ...string) *apw.APW {
	t.Helper()
	dir := t.TempDir()
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")

	for _, sid := range sysIDs {
		sys := apw.NewSystem(sid)
		sys.TransTCPIPEx = "10.0.0.1|1319|1|Desc||"
		fr := apw.NewFileRef("Panel", `UI\Panel.TP4`, apw.FileTypeTP4)
		fr.AddDeviceMap(apw.NewDeviceMap("Custom [10001:1:0]", "Custom [10001:1:0]"))
		sys.AddFile(fr)
		p.AddSystem(sys)
	}

	ws.AddProject(p)
	a := apw.NewAPW(ws.Identifier, filepath.Join(dir, "WS.apw"))
	a.SetWorkspace(ws)
	return a
}

func (s *FromAPWSystemSuite) TestKnownSystemReturnsItsItems() {
	a := multiSystemAPW(s.T(), "Room1", "Room2")
	list := FromAPWSystem(a, "Proj", "Room1")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "Room1", list.Items[0].SystemName)
}

func (s *FromAPWSystemSuite) TestUnknownSystemReturnsEmpty() {
	list := FromAPWSystem(multiSystemAPW(s.T(), "Room1"), "Proj", "NoSuch")
	assert.Empty(s.T(), list.Items)
}

func (s *FromAPWSystemSuite) TestOtherSystemsExcluded() {
	a := multiSystemAPW(s.T(), "A", "B", "C")
	list := FromAPWSystem(a, "Proj", "B")
	s.Require().Len(list.Items, 1)
	assert.Equal(s.T(), "B", list.Items[0].SystemName)
}

// ---------------------------------------------------------------------------
// ItemFieldsSuite — validate field values on produced items
// ---------------------------------------------------------------------------

type ItemFieldsSuite struct{ suite.Suite }

func TestItemFieldsSuite(t *testing.T) { suite.Run(t, new(ItemFieldsSuite)) }

func (s *ItemFieldsSuite) tknItem(sysID string) Item {
	s.T().Helper()
	h := newSysHelper(s.T(), sysID, "10.0.0.5|1319|1|Lab||").
		withMasterSrc("Main").withTKN()
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
	return list.Items[0]
}

func (s *ItemFieldsSuite) panelItem(ft apw.FileType, devAddr string) Item {
	s.T().Helper()
	ext := apw.AmxExtensions[ft]
	h := newSysHelper(s.T(), "0", "10.0.0.5|1319|1|Lab||").
		withPanel(ft, `User Interface\Panel`+ext, devAddr)
	list := FromAPW(h.build())
	s.Require().Len(list.Items, 1)
	return list.Items[0]
}

// --- Common fields ---

func (s *ItemFieldsSuite) TestPlatformIsNetLinx() {
	assert.Equal(s.T(), platformNetLinx, s.tknItem("0").Platform)
}

func (s *ItemFieldsSuite) TestFtTypeIsOne() {
	assert.Equal(s.T(), ftType, s.tknItem("0").FtType)
}

func (s *ItemFieldsSuite) TestWorkspaceNameCarried() {
	assert.Equal(s.T(), "TestWorkspace", s.tknItem("0").WorkspaceName)
}

func (s *ItemFieldsSuite) TestProjectNameCarried() {
	assert.Equal(s.T(), "TestProject", s.tknItem("0").ProjectName)
}

func (s *ItemFieldsSuite) TestSystemNameCarried() {
	assert.Equal(s.T(), "TestSystem", s.tknItem("0").SystemName)
}

func (s *ItemFieldsSuite) TestWorkspacePathNameCarried() {
	path := s.tknItem("0").WorkspacePathName
	assert.True(s.T(),
		strings.HasSuffix(filepath.ToSlash(path), "TestWorkspace.apw"),
		"unexpected WorkspacePathName: %s", path,
	)
}

func (s *ItemFieldsSuite) TestTpdFlagsEnabled() {
	item := s.tknItem("0")
	assert.Equal(s.T(), 1, item.IsTpdSendBitmapsEnabled)
	assert.Equal(s.T(), 1, item.IsTpdSendFontsEnabled)
	assert.Equal(s.T(), 1, item.IsTpdSendIconsEnabled)
}

func (s *ItemFieldsSuite) TestItemCheckedInListIsOne() {
	assert.Equal(s.T(), 1, s.tknItem("0").ItemCheckedInList)
}

func (s *ItemFieldsSuite) TestTransportSettingsPrefixIsT() {
	assert.True(s.T(),
		strings.HasPrefix(s.tknItem("0").TransportSettings, "T-"),
	)
}

func (s *ItemFieldsSuite) TestTransportSettingsContainsHost() {
	assert.Contains(s.T(), s.tknItem("0").TransportSettings, "10.0.0.5")
}

// --- MasterSrc / TKN fields ---

func (s *ItemFieldsSuite) TestMasterSrcDeviceIsZero() {
	assert.Equal(s.T(), deviceMaster, s.tknItem("0").Device)
}

func (s *ItemFieldsSuite) TestMasterSrcPortIsOne() {
	assert.Equal(s.T(), portMaster, s.tknItem("0").Port)
}

func (s *ItemFieldsSuite) TestMasterSrcSystemFromSysID() {
	assert.Equal(s.T(), 2, s.tknItem("2").System)
}

func (s *ItemFieldsSuite) TestMasterSrcSystemDefaultsToZero() {
	assert.Equal(s.T(), 0, s.tknItem("0").System)
}

func (s *ItemFieldsSuite) TestMasterSrcIsRebootRequired() {
	assert.Equal(s.T(), 1, s.tknItem("0").IsRebootRequired)
}

func (s *ItemFieldsSuite) TestMasterSrcSourceFileIsTKN() {
	assert.Equal(s.T(), apw.FileExtensionTKN,
		strings.ToLower(filepath.Ext(s.tknItem("0").SourceFile)),
	)
}

func (s *ItemFieldsSuite) TestMasterSrcNoSmartTransferFlags() {
	item := s.tknItem("0")
	assert.Equal(s.T(), 0, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 0, item.IsTP5SmartTransferEnabled)
}

func (s *ItemFieldsSuite) TestMasterSrcSourceFileIsAbsolute() {
	assert.True(s.T(), filepath.IsAbs(s.tknItem("0").SourceFile))
}

// --- Panel DPS fields (Custom [D:P:S] DevAddr) ---

func (s *ItemFieldsSuite) TestPanelDeviceFromDPS() {
	item := s.panelItem(apw.FileTypeTP4, "Custom [10001:1:0]")
	assert.Equal(s.T(), 10001, item.Device)
}

func (s *ItemFieldsSuite) TestPanelPortFromDPS() {
	item := s.panelItem(apw.FileTypeTP4, "Custom [10001:1:0]")
	assert.Equal(s.T(), 1, item.Port)
}

func (s *ItemFieldsSuite) TestPanelSystemFromDPS() {
	item := s.panelItem(apw.FileTypeTP4, "Custom [10001:1:2]")
	assert.Equal(s.T(), 2, item.System)
}

func (s *ItemFieldsSuite) TestPanelNoReboot() {
	assert.Equal(s.T(), 0, s.panelItem(apw.FileTypeTP4, "Custom [10001:1:0]").IsRebootRequired)
}

// --- TP4-specific ---

func (s *ItemFieldsSuite) TestTP4SmartTransferEnabled() {
	item := s.panelItem(apw.FileTypeTP4, "Custom [10001:1:0]")
	assert.Equal(s.T(), 1, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 0, item.IsTP5SmartTransferEnabled)
}

// --- TP5-specific ---

func (s *ItemFieldsSuite) TestTP5SmartTransferEnabled() {
	item := s.panelItem(apw.FileTypeTP5, "Custom [10001:1:0]")
	assert.Equal(s.T(), 0, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 1, item.IsTP5SmartTransferEnabled)
}

// --- KPB-specific ---

func (s *ItemFieldsSuite) TestKPBGlyphAndFontFlagsEnabled() {
	item := s.panelItem(apw.FileTypeKPB, "Custom [10001:1:0]")
	assert.Equal(s.T(), 1, item.IsKPBSendGlyphEnabled)
	assert.Equal(s.T(), 1, item.IsKPBSendFontEnabled)
}

func (s *ItemFieldsSuite) TestKPBSmartTransferFlagsZero() {
	item := s.panelItem(apw.FileTypeKPB, "Custom [10001:1:0]")
	assert.Equal(s.T(), 0, item.IsTP4SmartTransferEnabled)
	assert.Equal(s.T(), 0, item.IsTP5SmartTransferEnabled)
}

// ---------------------------------------------------------------------------
// TransportSettingsSuite
// ---------------------------------------------------------------------------

type TransportSettingsSuite struct{ suite.Suite }

func TestTransportSettingsSuite(t *testing.T) { suite.Run(t, new(TransportSettingsSuite)) }

func (s *TransportSettingsSuite) TestEachSystemUsesItsOwnTransport() {
	dir := s.T().TempDir()
	ws := apw.NewWorkspace("WS")
	p := apw.NewProject("Proj")

	for i, host := range []string{"192.168.1.1", "192.168.1.2"} {
		sys := apw.NewSystem("Sys" + strings.Repeat("I", i+1))
		sys.TransTCPIPEx = host + "|1319|1|Room||"
		name := "Main" + strings.Repeat("X", i)
		sys.AddFile(apw.NewFileRef(name, `Source\`+name+`.axs`, apw.FileTypeMasterSrc))
		p.AddSystem(sys)
		tknDir := filepath.Join(dir, "Source")
		_ = os.MkdirAll(tknDir, 0o755)
		_ = os.WriteFile(filepath.Join(tknDir, name+".tkn"), []byte{}, 0o644)
	}

	ws.AddProject(p)
	a := apw.NewAPW(ws.Identifier, filepath.Join(dir, "WS.apw"))
	a.SetWorkspace(ws)

	list := FromAPW(a)
	s.Require().Len(list.Items, 2)
	assert.NotEqual(s.T(), list.Items[0].TransportSettings, list.Items[1].TransportSettings)
	assert.Contains(s.T(), list.Items[0].TransportSettings, "192.168.1.1")
	assert.Contains(s.T(), list.Items[1].TransportSettings, "192.168.1.2")
}

func (s *TransportSettingsSuite) TestMultipleFilesInSystemShareTransport() {
	h := newSysHelper(s.T(), "0", "10.0.0.99|1319|1|Lab||").
		withMasterSrc("Main").withTKN().
		withPanel(apw.FileTypeTP4, `UI\Panel.TP4`, "Custom [10001:1:0]")

	list := FromAPW(h.build())
	s.Require().Len(list.Items, 2)
	assert.Equal(s.T(), list.Items[0].TransportSettings, list.Items[1].TransportSettings)
	assert.Contains(s.T(), list.Items[0].TransportSettings, "10.0.0.99")
}
