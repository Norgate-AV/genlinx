package apw

import "encoding/xml"

type DeviceMap struct {
	XMLName xml.Name `xml:"DeviceMap"`
	DevName string   `xml:"DevName,omitempty"`
	DevAddr string   `xml:"DevAddr,attr"`
}

func NewDeviceMap(addr, name string) *DeviceMap {
	return &DeviceMap{
		DevName: name,
		DevAddr: addr,
	}
}
