// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package callback

import (
	"context"
	"errors"
	"reflect" //#nosec:G702 // reflect is required for ABI parsing.

	"github.com/ethereum/go-ethereum/accounts/abi"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// Handler is the interface for all stateful precompiled contracts, which must
// expose their ABI methods and precompile methods for stateful execution.
type Handler interface {
	// ABIEvents() should return a map of Ethereum event names to Go-Ethereum abi `Event`.
	// NOTE: this can be directly loaded from the `Events` field of a Go-Ethereum ABI struct,
	// which can be built for a solidity library, interface, or contract.
	ABIEvents() map[string]abi.Event
}

// wrappedHandler is a container for running the underlying handler functions
// and building.
type wrappedHandler struct {
	// sigsToFns is a map of log signatures to their corresponding handler fn.
	sigsToFns map[logSig]*method
}

// New creates and returns a new `wrappedHandler` with the given method ids
// precompile functions map.
func NewFrom(
	si Handler,
) (LogHandler, error) {
	// We take all of the functions from the handler and build a map that
	// maps an ethereum log to a corresponding function to handle it.
	sigsToFns, err := buildIDsToMethods(si, reflect.ValueOf(si))
	if err != nil {
		return nil, err
	}

	return &wrappedHandler{
		sigsToFns: sigsToFns,
	}, nil
}

// HandleLog calls the function that matches the given log's signature.
func (sc *wrappedHandler) HandleLog(ctx context.Context, log *coretypes.Log) error {
	// Extract the method ID from the input and load the method.
	if log == nil {
		return errors.New("log is nil")
	}
	fn, found := sc.sigsToFns[logSig(log.Topics[0])]
	if !found {
		return ErrHandlerFnNotFound
	}

	return fn.Call(ctx, log)
}
