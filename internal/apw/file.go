package apw

import "encoding/xml"

type FileType string

const (
	FileTypeSource    FileType = "Source"
	FileTypeMasterSrc FileType = "MasterSrc"
	FileTypeInclude   FileType = "Include"
	FileTypeModule    FileType = "Module"
	FileTypeAXB       FileType = "AXB"
	FileTypeIR        FileType = "IR"
	FileTypeTPD       FileType = "TPD"
	FileTypeTP4       FileType = "TP4"
	FileTypeTP5       FileType = "TP5"
	FileTypeKPD       FileType = "KPD"
	FileTypeTKO       FileType = "TKO"
	FileTypeIRDB      FileType = "IRDB"
	FileTypeIRNDB     FileType = "IRNDB"
	FileTypeOther     FileType = "Other"
	FileTypeDuet      FileType = "Duet"
	FileTypeTOK       FileType = "TOK"
	FileTypeTKN       FileType = "TKN"
	FileTypeKPB       FileType = "KPB"
	FileTypeXDD       FileType = "XDD"
	FileTypeWorkspace FileType = "Workspace"
)

type FileCompileType string

const (
	FileCompileTypeNone    FileCompileType = "None"
	FileCompileTypeNetLinx FileCompileType = "NetLinx"
	FileCompileTypeAxcess  FileCompileType = "Axcess"
)

const (
	FileExtensionAPW   = ".apw"
	FileExtensionAXS   = ".axs"
	FileExtensionAXI   = ".axi"
	FileExtensionIRL   = ".irl"
	FileExtensionTP4   = ".tp4"
	FileExtensionTP5   = ".tp5"
	FileExtensionTPD   = ".tpd"
	FileExtensionJAR   = ".jar"
	FileExtensionXDD   = ".xdd"
	FileExtensionKPD   = ".kpd"
	FileExtensionAXB   = ".axb"
	FileExtensionTKO   = ".tko"
	FileExtensionIRDB  = ".irdb"
	FileExtensionIRNDB = ".irndb"
	FileExtensionTOK   = ".tok"
	FileExtensionTKN   = ".tkn"
	FileExtensionKPB   = ".kpb"
)

// AmxExtensions maps file type constants to their file extensions.
var AmxExtensions = map[FileType]string{
	FileTypeWorkspace: FileExtensionAPW,
	FileTypeModule:    FileExtensionAXS,
	FileTypeMasterSrc: FileExtensionAXS,
	FileTypeSource:    FileExtensionAXS,
	FileTypeInclude:   FileExtensionAXI,
	FileTypeIR:        FileExtensionIRL,
	FileTypeTP4:       FileExtensionTP4,
	FileTypeTP5:       FileExtensionTP5,
	FileTypeTPD:       FileExtensionTPD,
	FileTypeDuet:      FileExtensionJAR,
	FileTypeXDD:       FileExtensionXDD,
	FileTypeKPD:       FileExtensionKPD,
	FileTypeAXB:       FileExtensionAXB,
	FileTypeTKO:       FileExtensionTKO,
	FileTypeIRDB:      FileExtensionIRDB,
	FileTypeIRNDB:     FileExtensionIRNDB,
	FileTypeTOK:       FileExtensionTOK,
	FileTypeTKN:       FileExtensionTKN,
	FileTypeKPB:       FileExtensionKPB,
	FileTypeOther:     "",
}

// AmxCompiledExtensions maps file type constants to their compiled output extensions.
var AmxCompiledExtensions = map[FileType]string{
	FileTypeWorkspace: FileExtensionAPW,
	FileTypeModule:    FileExtensionTKO,
	FileTypeMasterSrc: FileExtensionTKN,
	FileTypeSource:    FileExtensionTKN,
	FileTypeInclude:   FileExtensionTKN,
	FileTypeIR:        FileExtensionIRL,
	FileTypeTP4:       FileExtensionTP4,
	FileTypeTP5:       FileExtensionTP5,
	FileTypeTPD:       FileExtensionTPD,
	FileTypeDuet:      FileExtensionJAR,
	FileTypeXDD:       FileExtensionXDD,
	FileTypeKPD:       FileExtensionKPD,
	FileTypeAXB:       FileExtensionAXB,
	FileTypeTKO:       FileExtensionTKO,
	FileTypeIRDB:      FileExtensionIRDB,
	FileTypeIRNDB:     FileExtensionIRNDB,
	FileTypeTOK:       FileExtensionTOK,
	FileTypeTKN:       FileExtensionTKN,
	FileTypeKPB:       FileExtensionKPB,
	FileTypeOther:     "",
}

// File represents a file reference parsed from an APW workspace.
type File struct {
	ID      string
	Type    FileType
	Path    string
	Exists  bool
	IsExtra bool
	Content string // optional in-memory content (e.g. for .env files)
}

type FileRef struct {
	XMLName         xml.Name        `xml:"File"`
	Identifier      string          `xml:"Identifier"`
	FilePathName    string          `xml:"FilePathName"`
	Comments        string          `xml:"Comments,omitempty"`
	MasterDirectory string          `xml:"MasterDirectory,omitempty"`
	DeviceMaps      []*DeviceMap    `xml:"DeviceMap"`
	IRDBs           []*IRDB         `xml:"IRDB"`
	Type            FileType        `xml:"Type,attr"`
	CompileType     FileCompileType `xml:"CompileType,attr"`
}

func NewFileRef(id, path string, t FileType) *FileRef {
	return &FileRef{
		Identifier:   id,
		FilePathName: path,
		Type:         t,
	}
}
