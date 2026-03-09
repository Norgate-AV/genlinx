package apw

import "encoding/xml"

type System struct {
	XMLName                  xml.Name   `xml:"System"`
	Identifier               string     `xml:"Identifier"`
	SysID                    string     `xml:"SysID"`
	TransTCPIP               string     `xml:"TransTCPIP,omitempty"`
	TransSerial              string     `xml:"TransSerial,omitempty"`
	TransTCPIPEx             string     `xml:"TransTCPIPEx,omitempty"`
	TransSerialEx            string     `xml:"TransSerialEx,omitempty"`
	TransUSBEx               string     `xml:"TransUSBEx,omitempty"`
	TransVNMEx               string     `xml:"TransVNMEx,omitempty"`
	VirtualNetLinxMasterFlag string     `xml:"VirtualNetLinxMasterFlag,omitempty"`
	VNMSystemID              string     `xml:"VNMSystemID,omitempty"`
	VNMIPAddress             string     `xml:"VNMIPAddress,omitempty"`
	VNMMaskAddress           string     `xml:"VNMMaskAddress,omitempty"`
	UserName                 string     `xml:"UserName"`
	Password                 string     `xml:"Password"`
	Comments                 string     `xml:"Comments"`
	Files                    []*FileRef `xml:"File"`
	IsActive                 string     `xml:"IsActive,attr"`
	Platform                 string     `xml:"Platform,attr"`
	Transport                string     `xml:"Transport,attr"`
	TransportEx              string     `xml:"TransportEx,attr"`
}

func NewSystem(id string) *System {
	return &System{
		Identifier: id,
	}
}

// AddFile appends fr to the system's file list and returns it so calls can be
// chained. If fr.CompileType is empty it is inferred from fr.Type.
func (s *System) AddFile(fr *FileRef) *FileRef {
	if fr.CompileType == "" {
		fr.CompileType = compileTypeFor(fr.Type)
	}

	s.Files = append(s.Files, fr)

	return fr
}

// RemoveFile removes the first FileRef whose Identifier matches id.
// Returns true if a file was removed, false if no match was found.
func (s *System) RemoveFile(id string) bool {
	for i, f := range s.Files {
		if f.Identifier == id {
			s.Files = append(s.Files[:i], s.Files[i+1:]...)
			return true
		}
	}

	return false
}
