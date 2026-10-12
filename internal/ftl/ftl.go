package ftl

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

const (
	// platformNetLinx is the FTL platform code for a NetLinx master controller.
	platformNetLinx = 2

	// ftType is the standard transfer type value used in every FTL item.
	ftType = 1

	// deviceMaster is the device address for NetLinx master (TKN) items.
	deviceMaster = 0

	// portMaster is the port number for NetLinx master (TKN) items.
	portMaster = 1
)

// FileTransferList is the root element of a File Transfer List (.ftl) document.
type FileTransferList struct {
	XMLName  xml.Name `xml:"Items"`
	Items    []Item   `xml:"Item"`
	Warnings []string `xml:"-"`
}

// Item represents a single transferable file entry within a File Transfer List.
// Field names and ordering match the schema produced by NetLinx Studio.
type Item struct {
	// Platform identifies the target hardware: 2 = NetLinx.
	Platform int `xml:"Platform"`

	// SourceFile is the absolute path to the file to be transferred.
	// For MasterSrc files this is the compiled .tkn output; for UI files it
	// is the panel binary (.TP4, .TP5, .KPB) as-is.
	SourceFile string `xml:"SourceFile"`

	// DestinationFile is unused by Studio's file transfer and left empty.
	DestinationFile string `xml:"DestinationFile"`

	// FtType is always 1 in workspace-derived FTL files.
	FtType int `xml:"FtType"`

	// WorkspaceName, ProjectName, SystemName trace the item back to the
	// hierarchy it was generated from.
	WorkspaceName string `xml:"WorkspaceName"`
	ProjectName   string `xml:"ProjectName"`
	SystemName    string `xml:"SystemName"`

	// WorkspacePathName is the absolute path to the originating .apw file.
	WorkspacePathName string `xml:"WorkspacePathName"`

	// Device is the NetLinx device address: 0 for the master (TKN), or the
	// device number from the panel's DeviceMap DPS address (TP4 / TP5 / KPB).
	Device int `xml:"Device"`

	// Port is 1 for master (TKN) items; for panel items it comes from the
	// port component of the panel's DeviceMap DPS address.
	Port int `xml:"Port"`

	// System carries the APW SysID for master (TKN) items, and the system
	// component of the panel's DeviceMap DPS address for panel items.
	System int `xml:"System"`

	IsDirectToDevice int `xml:"IsDirectToDevice"`

	// UI transfer flags — enabled for all items so Studio can send bitmaps,
	// fonts and icons when the item is a panel file.
	IsTpdSendBitmapsEnabled int `xml:"IsTpdSendBitmapsEnabled"`
	IsTpdSendFontsEnabled   int `xml:"IsTpdSendFontsEnabled"`
	IsTpdSendIconsEnabled   int `xml:"IsTpdSendIconsEnabled"`

	// Smart-transfer flags — set to 1 for the matching panel type only.
	IsTP4SmartTransferEnabled int `xml:"IsTP4SmartTransferEnabled"`
	IsTP5SmartTransferEnabled int `xml:"IsTP5SmartTransferEnabled"`

	// KPB-specific transfer flags.
	IsKPBSendGlyphEnabled int `xml:"IsKPBSendGlyphEnabled"`
	IsKPBSendFontEnabled  int `xml:"IsKPBSendFontEnabled"`

	// IsRebootRequired is 1 for master (TKN) items and 0 for panel items.
	IsRebootRequired int `xml:"IsRebootRequired"`

	IsSendTKNOnly     int    `xml:"IsSendTKNOnly"`
	ItemCheckedInList int    `xml:"ItemCheckedInList"`
	MasterDirectory   string `xml:"MasterDirectory"`

	// TransportSettings encodes the TCP/IP connection to the NetLinx master:
	// "T-<host>|<port>|<type>|<description>||"
	TransportSettings string `xml:"TransportSettings"`
}

// FromAPW converts an entire parsed workspace into a FileTransferList.
// Every system in every project is included.
func FromAPW(a *apw.APW) *FileTransferList {
	return fromAPW(a, "", "")
}

// FromAPWProject converts a single project within the workspace into a
// FileTransferList.  Returns an empty list if projectID is not found.
func FromAPWProject(a *apw.APW, projectID string) *FileTransferList {
	return fromAPW(a, projectID, "")
}

// FromAPWSystem converts a single system within the workspace into a
// FileTransferList.  Both projectID and systemID must be non-empty.
func FromAPWSystem(a *apw.APW, projectID, systemID string) *FileTransferList {
	return fromAPW(a, projectID, systemID)
}

// fromAPW is the shared implementation.  Empty projectFilter / systemFilter
// act as "include all".
func fromAPW(a *apw.APW, projectFilter, systemFilter string) *FileTransferList {
	apwDir := filepath.Dir(a.FilePath())
	ws := a.Workspace()

	ftl := &FileTransferList{}

	for _, project := range ws.Projects {
		if projectFilter != "" && project.Identifier != projectFilter {
			continue
		}

		for _, system := range project.Systems {
			if systemFilter != "" && system.Identifier != systemFilter {
				continue
			}

			sysID := parseSysID(system.SysID)
			reg := buildDeviceRegistry(system, apwDir)

			for _, fr := range system.Files {
				items := makeItems(
					fr,
					apwDir,
					ws.Identifier,
					project.Identifier,
					system.Identifier,
					a.FilePath(),
					system.TransTCPIPEx,
					sysID,
					reg,
					&ftl.Warnings,
				)
				ftl.Items = append(ftl.Items, items...)
			}
		}
	}

	return ftl
}

// parseSysID converts the APW SysID string to an int, defaulting to 0.
func parseSysID(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}

	return n
}

// makeItems builds the FTL Item(s) for a single FileRef.
//
// For MasterSrc files one item is returned if the compiled .tkn exists on
// disk, otherwise nil.  For panel files (TP4, TP5, KPB) one item is produced
// per DeviceMap entry whose DevAddr can be resolved to a concrete DPS address;
// entries with unresolvable symbolic names are silently skipped.  File types
// that are never transferred (Include, Module, Source, etc.) return nil.
func makeItems(
	fr *apw.FileRef,
	apwDir, workspaceName, projectName, systemName, apwPath, transTCPIPEx string,
	sysID int,
	reg deviceRegistry,
	warns *[]string,
) []Item {
	relPath := normalisePath(fr.FilePathName)

	baseItem := Item{
		Platform:                platformNetLinx,
		FtType:                  ftType,
		WorkspaceName:           workspaceName,
		ProjectName:             projectName,
		SystemName:              systemName,
		WorkspacePathName:       apwPath,
		IsTpdSendBitmapsEnabled: 1,
		IsTpdSendFontsEnabled:   1,
		IsTpdSendIconsEnabled:   1,
		ItemCheckedInList:       1,
		TransportSettings:       "T-" + transTCPIPEx,
	}

	switch fr.Type {
	case apw.FileTypeMasterSrc:
		tknRel := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + apw.FileExtensionTKN
		tknPath := filepath.Join(apwDir, tknRel)
		if _, err := os.Stat(tknPath); err != nil {
			*warns = append(*warns, fmt.Sprintf("[%s] TKN not found: %s", systemName, tknRel))
			return nil
		}

		item := baseItem
		item.SourceFile = tknPath
		item.Device = deviceMaster
		item.Port = portMaster
		item.System = sysID
		item.IsRebootRequired = 1

		return []Item{item}

	case apw.FileTypeTP4, apw.FileTypeTP5, apw.FileTypeKPB:
		return makePanelItems(fr, baseItem, apwDir, relPath, reg, systemName, warns)

	default:
		return nil
	}
}

// makePanelItems resolves each DeviceMap on a panel FileRef and returns one
// Item per successfully resolved DPS address.
func makePanelItems(
	fr *apw.FileRef,
	base Item,
	apwDir, relPath string,
	reg deviceRegistry,
	systemName string,
	warns *[]string,
) []Item {
	if len(fr.DeviceMaps) == 0 {
		return nil
	}

	srcPath := filepath.Join(apwDir, relPath)
	var items []Item

	for _, dm := range fr.DeviceMaps {
		dps, ok := resolveDevAddr(dm.DevAddr, reg)
		if !ok {
			*warns = append(*warns, fmt.Sprintf("[%s] unresolved DevAddr %q for %s", systemName, dm.DevAddr, filepath.Base(relPath)))
			continue
		}

		item := base
		item.SourceFile = srcPath
		item.Device = dps.Device
		item.Port = dps.Port
		item.System = dps.System

		switch fr.Type {
		case apw.FileTypeTP4:
			item.IsTP4SmartTransferEnabled = 1
		case apw.FileTypeTP5:
			item.IsTP5SmartTransferEnabled = 1
		case apw.FileTypeKPB:
			item.IsKPBSendGlyphEnabled = 1
			item.IsKPBSendFontEnabled = 1
		}

		items = append(items, item)
	}

	return items
}

// normalisePath converts an APW file path (which may use either backslash or
// forward-slash) to the OS-native separator so that filepath.Join works
// correctly.
func normalisePath(p string) string {
	return filepath.FromSlash(strings.ReplaceAll(p, `\`, `/`))
}
