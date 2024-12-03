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

package jsonrpc

import (
	"strings"

	"github.com/berachain/beacon-kit/errors"
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

var (
	// ErrParse indicates that invalid JSON was received by the server.
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

	// ErrServer is reserved for implementation-defined server errors.
	// (code: -32000 to -32099).
	ErrServer = errors.New(
		"received implementation-defined server-error (code: -32000 to -32099)",
	)

	// ErrServerParse indicates an error occurred on the server while
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
