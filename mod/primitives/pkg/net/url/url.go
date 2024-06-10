// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
