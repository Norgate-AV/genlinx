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

func NewIRDB(id, path string) *IRDB {
	return &IRDB{
		Property:       id,
		DOSName:        path,
		UserDBPathName: path,
		DBKey:          id,
	}
}
