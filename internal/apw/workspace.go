package apw

import "encoding/xml"

type Workspace struct {
	XMLName        xml.Name   `xml:"Workspace"`
	Identifier     string     `xml:"Identifier"`
	CreateVersion  string     `xml:"CreateVersion"`
	PJSFile        string     `xml:"PJS_File"`
	PJSConvertDate string     `xml:"PJS_ConvertDate"`
	PJSCreateDate  string     `xml:"PJS_CreateDate"`
	Comments       string     `xml:"Comments"`
	Projects       []*Project `xml:"Project"`
	CurrentVersion string     `xml:"CurrentVersion,attr"`
}

func NewWorkspace(id string) *Workspace {
	return &Workspace{
		Identifier: id,
	}
}

// AddProject appends p to the workspace and returns it so calls can be chained.
func (w *Workspace) AddProject(p *Project) *Project {
	w.Projects = append(w.Projects, p)
	return p
}
