package apw

import "encoding/xml"

type Project struct {
	XMLName       xml.Name  `xml:"Project"`
	Identifier    string    `xml:"Identifier"`
	Designer      string    `xml:"Designer,omitempty"`
	DealerID      string    `xml:"DealerID,omitempty"`
	SalesOrder    string    `xml:"SalesOrder"`
	PurchaseOrder string    `xml:"PurchaseOrder"`
	Comments      string    `xml:"Comments"`
	Systems       []*System `xml:"System"`
}

func NewProject(id string) *Project {
	return &Project{
		Identifier: id,
	}
}

// AddSystem appends s to the project and returns it so calls can be chained.
func (p *Project) AddSystem(s *System) *System {
	p.Systems = append(p.Systems, s)
	return s
}
