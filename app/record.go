package main

// SIPRegistration represents a single SIP registration
type SIPRegistration struct {
	AddressOfRecord string   `json:"addressOfRecord"`
	TenantID        string   `json:"tenantId"`
	URI             string   `json:"uri"`
	Contact         string   `json:"contact"`
	Path            []string `json:"path"`
	Source          string   `json:"source"`
	Target          string   `json:"target"`
	UserAgent       string   `json:"userAgent"`
	RawUserAgent    string   `json:"rawUserAgent"`
	Created         string   `json:"created"`
	LineID          string   `json:"lineId"`
}
