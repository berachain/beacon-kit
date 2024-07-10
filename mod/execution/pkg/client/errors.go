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

package client

import (
	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/rpc"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/http"
	jsonrpc "github.com/berachain/beacon-kit/mod/primitives/pkg/net/json-rpc"
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

	// ErrFailedToRefreshJWT indicates that the JWT could not be refreshed.
	ErrFailedToRefreshJWT = errors.New("failed to refresh auth token")

	// ErrMismatchedEth1ChainID is returned when the chainID does not
	// match the expected chain ID.
	ErrMismatchedEth1ChainID = errors.New("mismatched chain ID")
)

// Handles errors received from the RPC server according to the specification.
func (s *EngineClient[
	_, _, _, _, _,
]) handleRPCError(
	err error,
) error {
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
		var errWithData rpc.DataError
		errWithData, ok = err.(rpc.DataError) //nolint:errorlint // from prysm.
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
