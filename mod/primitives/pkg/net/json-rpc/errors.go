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

package jsonrpc

import (
	"strings"

	"github.com/berachain/beacon-kit/mod/errors"
)

// Error wraps RPC errors, which contain an error code in addition to the
// message.
type Error interface {
	// Error returns the error message.
	Error() string

	// ErrorCode returns the JSON-RPC error code.
	ErrorCode() int
}

// IsPreDefinedError returns true if the given
// error is a predefined JSON-RPC errors.
func IsPreDefinedError(err error) bool {
	return errors.IsAny(
		err,
		ErrParse,
		ErrInvalidRequest,
		ErrMethodNotFound,
		ErrInvalidParams,
		ErrInternal,
		ErrServer,
		ErrServerParse,
	)
}

// GetPredefinedError returns a predefined error for the given error code.
func GetPredefinedError(e Error) error {
	switch e.ErrorCode() {
	case -32700:
		// telemetry.IncrCounter(1, MetricKeyParseErrorCount)
		return ErrParse
	case -32600:
		// telemetry.IncrCounter(1, MetricKeyInvalidRequestCount)
		return ErrInvalidRequest
	case -32601:
		// telemetry.IncrCounter(1, MetricKeyMethodNotFoundCount)
		return ErrMethodNotFound
	case -32602:
		// telemetry.IncrCounter(1, MetricKeyInvalidParamsCount)
		return ErrInvalidParams
	case -32603:
		// telemetry.IncrCounter(1, MetricKeyInternalErrorCount)
		return ErrInternal
	default:
		return nil
	}
}

var (
	// ErrParseError indicates that invalid JSON was received by the server.
	// (code: -32700).
	ErrParse = errors.New(
		"invalid JSON was received by the server (code: -32700)",
	)

	// ErrInvalidRequest indicates that the JSON sent is not a valid Request
	// object. (code: -32600).
	ErrInvalidRequest = errors.New(
		"the JSON sent is not a valid Request object (code: -32600)",
	)

	// ErrMethodNotFound indicates that the method does not exist or is not
	// available. (code: -32601).
	ErrMethodNotFound = errors.New(
		"the method does not exist / is not available (code: -32601)",
	)

	// ErrInvalidParams indicates invalid method parameters. (code: -32602).
	ErrInvalidParams = errors.New("invalid method parameter(s) (code: -32602)")

	// ErrInternal indicates an internal JSON-RPC error. (code: -32603).
	ErrInternal = errors.New("internal JSON-RPC error (code: -32603)")

	// ErrServerError is reserved for implementation-defined server errors.
	// (code: -32000 to -32099).
	ErrServer = errors.New(
		"reserved for implementation-defined server-errors (code: -32000 to -32099)",
	)

	// ErrServerParseError indicates an error occurred on the server while
	// parsing the JSON text. (code: -32700).
	ErrServerParse = errors.New(
		"an error occurred on the server while parsing the JSON text (code: -32700)",
	)
)

// IsUnauthorizedError defines an interface for unauthorized errors.
func IsUnauthorizedError(err error) bool {
	if err == nil {
		return false
	}

	//nolint:errorlint // its'ok.
	e, ok := err.(Error)
	if !ok {
		if strings.Contains(
			e.Error(), "401 Unauthorized",
		) {
			return true
		}
	}
	return false
}
