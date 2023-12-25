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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

// CallbackHandler is the interface for all stateful precompiled contracts, which must
// expose their ABI methods and precompile methods for stateful execution.
type CallbackHandler interface {
	// Registrable

	// ABIEvents() should return a map of Ethereum event names to Go-Ethereum abi `Event`.
	// NOTE: this can be directly loaded from the `Events` field of a Go-Ethereum ABI struct,
	// which can be built for a solidity library, interface, or contract.
	ABIEvents() map[string]abi.Event
}

// wrappedCallbackHandler is a container for running wrappedCallbackHandler and precompiled contracts.
type wrappedCallbackHandler struct {
	// CallbackHandler is the base precompile implementation.
	CallbackHandler
	idsToMethods map[methodID]*method
}

// New creates and returns a new `wrappedCallbackHandler` with the given method ids
// precompile functions map.
func New(
	si CallbackHandler, idsToMethods map[methodID]*method,
) (*wrappedCallbackHandler, error) {
	return &wrappedCallbackHandler{
		CallbackHandler: si,
		idsToMethods:    idsToMethods,
	}, nil
}

func (sc *wrappedCallbackHandler) HandleLog(ctx context.Context, log types.Log) error {
	// Extract the method ID from the input and load the method.
	method, found := sc.idsToMethods[methodID(log.Topics[0].Bytes())]
	if !found {
		return errors.New("method not found")
	}

	return method.Call(ctx, log)
}
