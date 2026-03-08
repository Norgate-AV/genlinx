package apw

import "encoding/xml"

type Workspace struct {
	XMLName        xml.Name   `xml:"Workspace"`
	Identifier     string     `xml:"Identifier"`
	CreateVersion  string     `xml:"CreateVersion"`
	PJSFile        string     `xml:"PJS_File,omitempty"`
	PJSConvertDate string     `xml:"PJS_ConvertDate,omitempty"`
	PJSCreateDate  string     `xml:"PJS_CreateDate,omitempty"`
	Comments       string     `xml:"Comments,omitempty"`
	Projects       []*Project `xml:"Project"`
	CurrentVersion string     `xml:"CurrentVersion,attr"`
}

func NewWorkspace(id string) *Workspace {
	return &Workspace{
		Identifier: id,
	}
}
