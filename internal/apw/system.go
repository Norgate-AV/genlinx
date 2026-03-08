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
	UserName                 string     `xml:"UserName,omitempty"`
	Password                 string     `xml:"Password,omitempty"`
	Comments                 string     `xml:"Comments,omitempty"`
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
