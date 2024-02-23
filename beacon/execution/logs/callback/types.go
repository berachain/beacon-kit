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

package callback

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// LogHandler represents a struct that has the ability to ingest
// an ethereum log and handle it.
type LogHandler interface {
	HandleLog(ctx context.Context, log *coretypes.Log) error
}

// ContractHandler is the interface for all contracts in the execution layer,
// which must expose their ABI methods and events so that callback
// can select corresponding methods to process them.
type ContractHandler interface {
	// ABIEvents() should return a map of Ethereum event names to Go-Ethereum
	// abi `Event`. NOTE: this can be directly loaded from the `Events` field of
	// a Go-Ethereum ABI struct,
	// which can be built for a solidity library, interface, or contract.
	ABIEvents() map[string]abi.Event
}
