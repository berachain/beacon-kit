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

package eth

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/telemetry"
	gethRPC "github.com/ethereum/go-ethereum/rpc"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
)

// ErrUnauthenticatedConnection indicates that the connection is not authenticated.
const UnauthenticatedConnectionErrorStr = `could not verify execution chain ID as your 
	connection is not authenticated. If connecting to your execution client via HTTP, you 
	will need to set up JWT authentication...`

var (

	// ErrInvalidJWTSecretLength indicates that the JWT secret length is invalid.
	ErrInvalidJWTSecretLength = errors.New("invalid JWT secret length")

	// ErrHTTPTimeout returns true if the error is a http.Client timeout error.
	ErrHTTPTimeout = errors.New("timeout from http.Client")

	// ErrEmptyBlockHash indicates that the block hash is empty.
	ErrEmptyBlockHash = errors.New("block hash is empty 0x0000...000")
)

// Handles errors received from the RPC server according to the specification.
func (s *Eth1Client) handleRPCError(err error) error {
	if err == nil {
		return nil
	}
	if isTimeout(err) {
		return execution.ErrHTTPTimeout
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
		return execution.ErrParse
	case -32600:
		telemetry.IncrCounter(1, MetricKeyInvalidRequestCount)
		return execution.ErrInvalidRequest
	case -32601:
		telemetry.IncrCounter(1, MetricKeyMethodNotFoundCount)
		return execution.ErrMethodNotFound
	case -32602:
		telemetry.IncrCounter(1, MetricKeyInvalidParamsCount)
		return execution.ErrInvalidParams
	case -32603:
		telemetry.IncrCounter(1, MetricKeyInternalErrorCount)
		return execution.ErrInternal
	case -38001:
		telemetry.IncrCounter(1, MetricKeyUnknownPayloadErrorCount)
		return execution.ErrUnknownPayload
	case -38002:
		telemetry.IncrCounter(1, MetricKeyInvalidForkchoiceStateCount)
		return execution.ErrInvalidForkchoiceState
	case -38003:
		telemetry.IncrCounter(1, MetricKeyInvalidPayloadAttributesCount)
		return execution.ErrInvalidPayloadAttributes
	case -38004:
		telemetry.IncrCounter(1, MetricKeyRequestTooLargeCount)
		return execution.ErrRequestTooLarge
	case -32000:
		telemetry.IncrCounter(1, MetricKeyInternalServerErrorCount)
		// Only -32000 status codes are data errors in the RPC specification.
		var errWithData gethRPC.DataError
		errWithData, ok = err.(gethRPC.DataError) //nolint:errorlint // from prysm.
		if !ok {
			return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
		}
		return errors.Wrapf(execution.ErrServer, "%v", errWithData.Error())
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
