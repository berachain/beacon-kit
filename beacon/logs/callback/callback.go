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
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	evmv1 "github.com/itsdevbear/bolaris/types/evm/v1"
)

type LogHandler interface {
	HandleLog(ctx context.Context, log *evmv1.Log) error
}

// Handler is the interface for all stateful precompiled contracts, which must
// expose their ABI methods and precompile methods for stateful execution.
type Handler interface {
	// Registrable

	// ABIEvents() should return a map of Ethereum event names to Go-Ethereum abi `Event`.
	// NOTE: this can be directly loaded from the `Events` field of a Go-Ethereum ABI struct,
	// which can be built for a solidity library, interface, or contract.
	ABIEvents() map[string]abi.Event
}

// wrappedHandler is a container for running wrappedHandler and precompiled contracts.
type wrappedHandler struct {
	// Handler is the base precompile implementation.
	Handler
	idsToMethods map[logSig]*method
}

// New creates and returns a new `wrappedHandler` with the given method ids
// precompile functions map.
func NewFrom(
	si Handler,
) (LogHandler, error) {
	idsToMethods, err := buildIdsToMethods(si, reflect.ValueOf(si))
	if err != nil {
		return nil, err
	}

	return &wrappedHandler{
		Handler:      si,
		idsToMethods: idsToMethods,
	}, nil
}

// HandleLog calls the function that matches the given log's signature.
func (sc *wrappedHandler) HandleLog(ctx context.Context, log *evmv1.Log) error {
	// Extract the method ID from the input and load the method.
	method, found := sc.idsToMethods[logSig(log.Topics[0])]
	if !found {
		return errors.New("method not found")
	}

	return method.Call(ctx, log)
}
