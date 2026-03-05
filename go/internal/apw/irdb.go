package apw

import "encoding/xml"

type IRDB struct {
	XMLName        xml.Name `xml:"IRDB"`
	Property       string   `xml:"Property"`
	DOSName        string   `xml:"DOSName"`
	UserDBPathName string   `xml:"UserDBPathName"`
	Notes          string   `xml:"Notes"`
	DBKey          string   `xml:"DBKey,attr"`
}
