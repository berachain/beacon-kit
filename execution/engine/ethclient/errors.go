// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package ethclient

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/telemetry"
	gethRPC "github.com/ethereum/go-ethereum/rpc"

	"github.com/pkg/errors"
)

// ErrUnauthenticatedConnection indicates that the connection is not authenticated.
const UnauthenticatedConnectionErrorStr = `could not verify execution chain ID as your 
	connection is not authenticated. If connecting to your execution client via HTTP, you 
	will need to set up JWT authentication...`

var (
	// ErrInvalidJWTSecretLength signifies the JWT secret length is not correct.
	ErrInvalidJWTSecretLength = errors.New("invalid JWT secret length")
	// ErrHTTPTimeout indicates a timeout error from http.Client.
	ErrHTTPTimeout = errors.New("timeout from http.Client")
	// ErrEmptyBlockHash denotes an empty block hash scenario.
	ErrEmptyBlockHash = errors.New("block hash is empty 0x0000...000")
	// ErrNilResponse signifies that the execution client returned a nil response.
	ErrNilResponse = errors.New("nil response from execution client")
	// ErrParse is associated with JSON-RPC code -32700, indicating a parsing error.
	ErrParse = errors.New("invalid JSON was received by the server")
	// ErrInvalidRequest is associated with JSON-RPC code -32600,
	// indicating an invalid request object.
	ErrInvalidRequest = errors.New("JSON sent is not valid request object")
	// ErrMethodNotFound is associated with JSON-RPC code -32601,
	// indicating a non-existent method was called.
	ErrMethodNotFound = errors.New("method not found")
	// ErrInvalidParams is associated with JSON-RPC code -32602,
	// indicating method parameters were invalid.
	ErrInvalidParams = errors.New("invalid method parameter(s)")
	// ErrInternal is associated with JSON-RPC code -32603,
	// indicating an internal JSON-RPC error.
	ErrInternal = errors.New("internal JSON-RPC error")
	// ErrServer is associated with JSON-RPC code -32000,
	// indicating a client-side error during request processing.
	ErrServer = errors.New("client error while processing request")
	// ErrUnknownPayload is associated with JSON-RPC code -38001,
	// indicating an unavailable or non-existent payload.
	ErrUnknownPayload = errors.New("payload does not exist or is not available")
	// ErrInvalidForkchoiceState is associated with JSON-RPC code -38002,
	// indicating an invalid fork choice state.
	ErrInvalidForkchoiceState = errors.New("invalid forkchoice state")
	// ErrInvalidPayloadAttributes is associated with JSON-RPC code -38003,
	// indicating invalid or inconsistent payload attributes.
	ErrInvalidPayloadAttributes = errors.New("payload attributes are invalid / inconsistent")
	// ErrUnknownPayloadStatus indicates an unknown status for the payload.
	ErrUnknownPayloadStatus = errors.New("unknown payload status")
	// ErrAcceptedSyncingPayloadStatus indicates a payload status of either SYNCING or ACCEPTED.
	ErrAcceptedSyncingPayloadStatus = errors.New("payload status is SYNCING or ACCEPTED")
	// ErrInvalidPayloadStatus indicates an invalid status for the payload.
	ErrInvalidPayloadStatus = errors.New("payload status is INVALID")
	// ErrInvalidBlockHashPayloadStatus indicates a failure in validating the
	// block hash for the payload.
	ErrInvalidBlockHashPayloadStatus = errors.New("payload status is INVALID_BLOCK_HASH")
	// ErrRequestTooLarge indicates that the request size exceeded the limit.
	ErrRequestTooLarge = errors.New("request too large")
	// ErrUnsupportedVersion indicates a request for a block type with an unknown
	// ExecutionPayload schema.
	ErrUnsupportedVersion = errors.New("unknown ExecutionPayload schema for block version")
	// ErrNilJWTSecret indicates that the JWT secret is nil.
	ErrNilJWTSecret = errors.New("nil JWT secret")
)

// Handles errors received from the RPC server according to the specification.
func (s *Eth1Client) handleRPCError(err error) error {
	if err == nil {
		return nil
	}
	if isTimeout(err) {
		return ErrHTTPTimeout
	}
	e, ok := err.(gethRPC.Error) //nolint:errorlint // from prysm.
	if !ok {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			s.logger.Error("HTTP authentication to your execution client is not working. " +
				"Please ensure you are setting a correct value for the --jwt-secret flag in " +
				"Prysm, or use an IPC connection if on the same machine.")
			return fmt.Errorf("could not authenticate connection to execution client: %w", err)
		}
		return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
	}
	switch e.ErrorCode() {
	case -32700:
		telemetry.IncrCounter(1, MetricKeyParseErrorCount)
		return ErrParse
	case -32600:
		telemetry.IncrCounter(1, MetricKeyInvalidRequestCount)
		return ErrInvalidRequest
	case -32601:
		telemetry.IncrCounter(1, MetricKeyMethodNotFoundCount)
		return ErrMethodNotFound
	case -32602:
		telemetry.IncrCounter(1, MetricKeyInvalidParamsCount)
		return ErrInvalidParams
	case -32603:
		telemetry.IncrCounter(1, MetricKeyInternalErrorCount)
		return ErrInternal
	case -38001:
		telemetry.IncrCounter(1, MetricKeyUnknownPayloadErrorCount)
		return ErrUnknownPayload
	case -38002:
		telemetry.IncrCounter(1, MetricKeyInvalidForkchoiceStateCount)
		return ErrInvalidForkchoiceState
	case -38003:
		telemetry.IncrCounter(1, MetricKeyInvalidPayloadAttributesCount)
		return ErrInvalidPayloadAttributes
	case -38004:
		telemetry.IncrCounter(1, MetricKeyRequestTooLargeCount)
		return ErrRequestTooLarge
	case -32000:
		telemetry.IncrCounter(1, MetricKeyInternalServerErrorCount)
		// Only -32000 status codes are data errors in the RPC specification.
		var errWithData gethRPC.DataError
		errWithData, ok = err.(gethRPC.DataError) //nolint:errorlint // from prysm.
		if !ok {
			return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
		}
		return errors.Wrapf(ErrServer, "%v", errWithData.Error())
	default:
		return err
	}
}

type httpTimeoutError interface {
	Error() string
	Timeout() bool
}

func isTimeout(e error) bool {
	t, ok := e.(httpTimeoutError) //nolint:errorlint // from prysm.
	return ok && t.Timeout()
}
