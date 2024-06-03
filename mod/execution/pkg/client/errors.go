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

package client

import (
	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/http"
	jsonrpc "github.com/berachain/beacon-kit/mod/primitives/pkg/net/json-rpc"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
)

// ErrUnauthenticatedConnection indicates that the connection is not
// authenticated.
//
//nolint:lll
const (
	UnauthenticatedConnectionErrorStr = `could not verify execution chain ID as your 
	connection is not authenticated. If connecting to your execution client via HTTP, you 
	will need to set up JWT authentication...`

	AuthErrMsg = "HTTP authentication to your execution client " +
		"is not working. Please ensure you are setting a correct " +
		"value for the JWT secret path" +
		"is set correctly, or use an IPC " +
		"connection if on the same machine."
)

var (
	// ErrNotStarted indicates that the execution client is not started.
	ErrNotStarted = errors.New("engine client is not started")
)

// Handles errors received from the RPC server according to the specification.
func (s *EngineClient[ExecutionPayloadDenebT]) handleRPCError(err error) error {
	// Exit early if there is no error.
	if err == nil {
		return nil
	}

	// Check for timeout errors.
	if http.IsTimeoutError(err) {
		s.metrics.incrementHTTPTimeoutCounter()
		return http.ErrTimeout
	}

	// Check for connection errors.
	//
	//nolint:errorlint // from prysm.
	e, ok := err.(jsonrpc.Error)
	if !ok {
		if jsonrpc.IsUnauthorizedError(e) {
			return http.ErrUnauthorized
		}
		return errors.Wrapf(
			err,
			"got an unexpected server error in JSON-RPC response "+
				"failed to convert from jsonrpc.Error",
		)
	}

	// Otherwise check for our engine errors.
	switch e.ErrorCode() {
	case -32700:
		s.metrics.incrementParseErrorCounter()
		return jsonrpc.ErrParse
	case -32600:
		s.metrics.incrementInvalidRequestCounter()
		return jsonrpc.ErrInvalidRequest
	case -32601:
		s.metrics.incrementMethodNotFoundCounter()
		return jsonrpc.ErrMethodNotFound
	case -32602:
		s.metrics.incrementInvalidParamsCounter()
		return jsonrpc.ErrInvalidParams
	case -32603:
		s.metrics.incrementInternalErrorCounter()
		return jsonrpc.ErrInternal
	case -38001:
		s.metrics.incrementUnknownPayloadErrorCounter()
		return engineerrors.ErrUnknownPayload
	case -38002:
		s.metrics.incrementInvalidForkchoiceStateCounter()
		return engineerrors.ErrInvalidForkchoiceState
	case -38003:
		s.metrics.incrementInvalidPayloadAttributesCounter()
		return engineerrors.ErrInvalidPayloadAttributes
	case -38004:
		s.metrics.incrementRequestTooLargeCounter()
		return engineerrors.ErrRequestTooLarge
	case -32000:
		s.metrics.incrementInternalServerErrorCounter()
		// Only -32000 status codes are data errors in the RPC specification.
		var errWithData gethRPC.DataError
		errWithData, ok = err.(gethRPC.DataError) //nolint:errorlint // from prysm.
		if !ok {
			return errors.Wrapf(
				errors.Join(jsonrpc.ErrServer, err),
				"got an unexpected data error in JSON-RPC response",
			)
		}
		return errors.Wrapf(jsonrpc.ErrServer, "%v", errWithData.Error())
	default:
		return err
	}
}
