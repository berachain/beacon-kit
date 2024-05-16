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

package errors

import (
	"github.com/berachain/beacon-kit/mod/errors"
)

var (
	// ErrUnknownPayload indicates an unavailable or non-existent payload
	// (JSON-RPC code -38001).
	ErrUnknownPayload = errors.New(
		"payload does not exist or is not available")

	// ErrInvalidForkchoiceState indicates an invalid fork choice state
	// (JSON-RPC code -38002).
	ErrInvalidForkchoiceState = errors.New(
		"invalid forkchoice state")

	// ErrInvalidPayloadAttributes indicates invalid or inconsistent payload
	// attributes
	// (JSON-RPC code -38003).
	ErrInvalidPayloadAttributes = errors.New(
		"payload attributes are invalid / inconsistent")

	// ErrRequestTooLarge indicates that the request is too large
	// (JSON-RPC code -38004).
	ErrRequestTooLarge = errors.New(
		"request is too large",
	)

	// ErrUnknownPayloadStatus indicates an unknown payload status.
	ErrUnknownPayloadStatus = errors.New(
		"unknown payload status")

	// ErrAcceptedPayloadStatus indicates a payload status of ACCEPTED.
	ErrAcceptedPayloadStatus = errors.New(
		"payload status is ACCEPTED")

	// ErrSyncingPayloadStatus indicates a payload status of SYNCING.
	ErrSyncingPayloadStatus = errors.New(
		"payload status is SYNCING",
	)

	// ErrInvalidPayloadStatus indicates an invalid payload status.
	ErrInvalidPayloadStatus = errors.New(
		"payload status is INVALID")

	// ErrInvalidBlockHashPayloadStatus indicates a failure in validating the
	// block hash for the payload.
	ErrInvalidBlockHashPayloadStatus = errors.New(
		"payload status is INVALID_BLOCK_HASH")

	// ErrNilForkchoiceResponse indicates a nil forkchoice response.
	ErrNilForkchoiceResponse = errors.New(
		"nil forkchoice response",
	)
	/// ErrNilBlobsBundle is returned when nil blobs bundle is received.
	ErrNilBlobsBundle = errors.New(
		"nil blobs bundle received from execution client")

	// ErrInvalidPayloadAttributeVersion indicates an invalid version of payload
	// attributes was provided.
	ErrInvalidPayloadAttributeVersion = errors.New(
		"invalid payload attribute version")
	// ErrInvalidPayloadType indicates an invalid payload type
	// was provided for an RPC call.
	ErrInvalidPayloadType = errors.New("invalid payload type for RPC call")

	// ErrInvalidGetPayloadVersion indicates that an unknown fork version was
	// provided for getting a payload.
	ErrInvalidGetPayloadVersion = errors.New("unknown fork for get payload")

	// ErrUnsupportedVersion indicates a request for a block type with an
	// unknown ExecutionPayload schema.
	ErrUnsupportedVersion = errors.New(
		"unknown ExecutionPayload schema for block version")

	// ErrNilJWTSecret indicates that the JWT secret is nil.
	ErrNilJWTSecret = errors.New("nil JWT secret")

	// ErrNilAttributesPassedToClient is returned when nil attributes are
	// passed to the client.
	ErrNilAttributesPassedToClient = errors.New(
		"nil attributes passed to client",
	)

	// ErrNilExecutionPayloadEnvelope is returned when nil execution payload
	// envelope is received.
	ErrNilExecutionPayloadEnvelope = errors.New(
		"nil execution payload envelope received from execution client")

	// ErrNilExecutionPayload is returned when nil execution payload
	// envelope is received.
	ErrNilExecutionPayload = errors.New(
		"nil execution payload received from execution client")

	// ErrNilPayloadStatus is returned when nil payload status is received.
	ErrNilPayloadStatus = errors.New(
		"nil payload status received from execution client",
	)
)
