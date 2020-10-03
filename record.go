package main

type SIPRecord struct {
	AddressOfRecord string   `json:addressOfRecord`
	TenantId        string   `json:tenantId`
	Uri             string   `json:uri`
	Contact         string   `json:contact`
	Path            []string `json:path`
	Source          string   `json:source`
	Target          string   `json:target`
	UserAgent       string   `json:userAgent`
	RawUserAgent    string   `json:rawUserAgent`
	Created         string   `json:created`
	LineId          string   `json:lineId`
}
