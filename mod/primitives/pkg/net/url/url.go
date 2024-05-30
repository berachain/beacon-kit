// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package url

import "net/url"

// ConnectionURL is a URL struct that is used to dial the execution client.
type ConnectionURL struct {
	*url.URL
}

// NewDialURL creates a new DialURL.
func NewDialURL(u *url.URL) *ConnectionURL {
	return &ConnectionURL{u}
}

func NewFromRaw(raw string) (*ConnectionURL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	return NewDialURL(u), nil
}

// IsHTTP checks if the DialURL scheme is HTTP.
func (d *ConnectionURL) IsHTTP() bool {
	return d.Scheme == "http"
}

// IsHTTPS checks if the DialURL scheme is HTTPS.
func (d *ConnectionURL) IsHTTPS() bool {
	return d.Scheme == "https"
}

// IsIPC checks if the DialURL scheme is IPC.
func (d *ConnectionURL) IsIPC() bool {
	return d.Scheme == "ipc"
}
