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

// RemoveProject removes the first project whose Identifier matches id.
// Returns true if a project was removed, false if no match was found.
func (w *Workspace) RemoveProject(id string) bool {
	for i, p := range w.Projects {
		if p.Identifier == id {
			w.Projects = append(w.Projects[:i], w.Projects[i+1:]...)
			return true
		}
	}

	return false
}
