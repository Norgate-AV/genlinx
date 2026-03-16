package ftl

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/Norgate-AV/genlinx/internal/apw"
)

const (
	// platformNetLinx is the FTL platform code for a NetLinx master controller.
	platformNetLinx = 2

	// ftType is the standard transfer type value used in every FTL item.
	ftType = 1

	// deviceMaster is the device address used for NetLinx master (TKN) items.
	deviceMaster = 0

	// devicePanel is the default device address used for touchpanel items
	// (TP4, TP5, KPB).  AMX Modero-series panels are addressed at 10001 by
	// default.
	devicePanel = 10001
)

// FileTransferList is the root element of a File Transfer List (.ftl) document.
type FileTransferList struct {
	XMLName xml.Name `xml:"Items"`
	Items   []Item   `xml:"Item"`
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

	// Device is the NetLinx device address: 0 for the master (TKN), 10001
	// for a default panel device (TP4 / TP5 / KPB).
	Device int `xml:"Device"`

	// Port is always 1.
	Port int `xml:"Port"`

	// System is always 0 — Studio uses 0 to mean "same system".
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

// fromAPW is the shared implementation.  Empty projectID / systemID act as
// "include all".
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

			for _, fr := range system.Files {
				item := makeItem(
					fr,
					apwDir,
					ws.Identifier,
					project.Identifier,
					system.Identifier,
					a.FilePath(),
					system.TransTCPIPEx,
				)

				if item != nil {
					ftl.Items = append(ftl.Items, *item)
				}
			}
		}
	}

	return ftl
}

// makeItem builds a single FTL Item for the given FileRef.  File types that
// do not produce transferable items (source, include, module, IR, etc.) return
// nil and are silently skipped.
func makeItem(
	fr *apw.FileRef,
	apwDir, workspaceName, projectName, systemName, apwPath, transTCPIPEx string,
) *Item {
	relPath := normalisePath(fr.FilePathName)

	base := &Item{
		Platform:                platformNetLinx,
		FtType:                  ftType,
		WorkspaceName:           workspaceName,
		ProjectName:             projectName,
		SystemName:              systemName,
		WorkspacePathName:       apwPath,
		Port:                    1,
		IsTpdSendBitmapsEnabled: 1,
		IsTpdSendFontsEnabled:   1,
		IsTpdSendIconsEnabled:   1,
		ItemCheckedInList:       1,
		TransportSettings:       "T-" + transTCPIPEx,
	}

	switch fr.Type {
	case apw.FileTypeMasterSrc:
		tknRel := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + apw.FileExtensionTKN
		base.SourceFile = filepath.Join(apwDir, tknRel)
		base.Device = deviceMaster
		base.IsRebootRequired = 1

	case apw.FileTypeTP4:
		base.SourceFile = filepath.Join(apwDir, relPath)
		base.Device = devicePanel
		base.IsTP4SmartTransferEnabled = 1

	case apw.FileTypeTP5:
		base.SourceFile = filepath.Join(apwDir, relPath)
		base.Device = devicePanel
		base.IsTP5SmartTransferEnabled = 1

	case apw.FileTypeKPB:
		base.SourceFile = filepath.Join(apwDir, relPath)
		base.Device = devicePanel
		base.IsKPBSendGlyphEnabled = 1
		base.IsKPBSendFontEnabled = 1

	default:
		return nil
	}

	return base
}

// normalisePath converts an APW file path (which may use either backslash or
// forward-slash) to the OS-native separator so that filepath.Join works
// correctly.
func normalisePath(p string) string {
	return filepath.FromSlash(strings.ReplaceAll(p, `\`, `/`))
}
