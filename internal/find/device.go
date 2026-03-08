package find

// Date holds numeric and text representations of a date.
type Date struct {
	Numeric string `json:"numeric"`
	Text    string `json:"text"`
}

// Device represents a discovered NetLinx device.
type Device struct {
	IP       string `json:"ip"`
	System   int    `json:"system"`
	Date     Date   `json:"date"`
	Time     string `json:"time"`
	Day      string `json:"day"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
}
