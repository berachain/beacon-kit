package client

import "net/url"

// DialURL is a URL struct that is used to dial the execution client.
type DialURL struct {
	*url.URL
}

// NewDialURL creates a new DialURL.
func NewDialURL(u *url.URL) *DialURL {
	return &DialURL{u}
}

// IsHTTP checks if the DialURL scheme is HTTP.
func (d *DialURL) IsHTTP() bool {
	return d.Scheme == "http"
}

// IsHTTPS checks if the DialURL scheme is HTTPS.
func (d *DialURL) IsHTTPS() bool {
	return d.Scheme == "https"
}

// IsIPC checks if the DialURL scheme is IPC.
func (d *DialURL) IsIPC() bool {
	return d.Scheme == "ipc"
}
