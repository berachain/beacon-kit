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
