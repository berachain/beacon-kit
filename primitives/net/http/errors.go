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

package http

import "errors"

var (
	// ErrTimeout indicates a timeout error from http.Client.
	ErrTimeout = errors.New("timeout from HTTP client")

	// ErrUnauthorized indicates an unauthorized error from http.Client.
	ErrUnauthorized = errors.New("401 unauthorized request")
)

// TimeoutError defines an interface for timeout errors.
// It includes methods for error message retrieval and timeout status checking.
type TimeoutError interface {
	// Error returns the error message.
	Error() string
	// Timeout indicates whether the error is a timeout error.
	Timeout() bool
}

// isTimeout checks if the given error is a timeout error.
// It asserts the error to the httpTimeoutError interface and checks its Timeout
// status.
// Returns true if the error is a timeout error, false otherwise.
func IsTimeoutError(e error) bool {
	if e == nil {
		return false
	}
	//nolint:errorlint // by design.
	t, ok := e.(TimeoutError)
	return ok && t.Timeout()
}
