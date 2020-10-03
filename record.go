package main

type SIPRecord struct {
	addressOfRecord string
	tenantId        string
	uri             string
	contact         string
	path            []string
	source          string
	target          string
	userAgent       string
	rawUserAgent    string
	created         string
	lineId          string
}
