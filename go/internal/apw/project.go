package apw

import "encoding/xml"

type Project struct {
	XMLName       xml.Name  `xml:"Project"`
	Identifier    string    `xml:"Identifier"`
	Designer      string    `xml:"Designer,omitempty"`
	DealerID      string    `xml:"DealerID,omitempty"`
	SalesOrder    string    `xml:"SalesOrder,omitempty"`
	PurchaseOrder string    `xml:"PurchaseOrder,omitempty"`
	Comments      string    `xml:"Comments,omitempty"`
	Systems       []*System `xml:"System"`
}
