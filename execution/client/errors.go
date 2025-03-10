// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/net/http"
	jsonrpc "github.com/berachain/beacon-kit/primitives/net/json-rpc"
)

const (
	UnauthenticatedConnectionErrorStr = `could not verify execution chain ID as your
	connection is not authenticated. If connecting to your execution client via HTTP, you
	will need to set up JWT authentication...`
)

var (

	// ErrMismatchedEth1ChainID is returned when the chainID does not
	// match the expected chain ID.
	ErrMismatchedEth1ChainID = errors.New("mismatched chain ID")

	// ErrBadConnection indicates that the http.Client was unable to
	// establish a connection.
	ErrBadConnection = errors.New("connection error")
)

// Handles errors received from the RPC server according to the specification.
func (s *EngineClient) handleRPCError(
	err error,
) error {
	if err == nil {
		//nolint:nilerr // appease nilaway
		return err
	}

	// Check for timeout errors.
	if errors.Is(err, engineerrors.ErrEngineAPITimeout) {
		s.metrics.incrementEngineAPITimeout()
		return err
	}
	if http.IsTimeoutError(err) {
		s.metrics.incrementHTTPTimeoutCounter()
		return http.ErrTimeout
	}
	// Check for authorization errors
	if errors.Is(err, http.ErrUnauthorized) {
		return err
	}
	// Check for connection errors.
	var e jsonrpc.Error
	ok := errors.As(err, &e)
	if !ok || e == nil {
		return errors.Join(ErrBadConnection, err)
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
		// // Only -32000 status codes are data errors in the RPC specification.
		// var errWithData rpc.DataError
		// errWithData, ok = err.(rpc.DataError) //nolint:errorlint // from
		// prysm.
		// if !ok {
		// 	return errors.Wrapf(
		// 		errors.Join(jsonrpc.ErrServer, err),
		// 		"got an unexpected data error in JSON-RPC response",
		// 	)
		// }
		return errors.Wrapf(jsonrpc.ErrServer, "%v", e.Error())
	default:
		return err
	}
}

// IsFatalError defines errors that indicate a bad request or an otherwise
// unusable client.
func IsFatalError(err error) bool {
	return jsonrpc.IsPreDefinedError(err) || errors.IsAny(
		err,
		ErrBadConnection,
	)
}

// IsNonFatalError defines errors that should be ephemeral and can be
// recovered from simply by retrying.
func IsNonFatalError(err error) bool {
	return errors.IsAny(
		err,
		engineerrors.ErrEngineAPITimeout,
		http.ErrTimeout,
	)
}
