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

package ssf

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// SingleSlotFinalityStore is the interface for the required storage
// backend for the SingleSlotFinality forkchoice store.
type SingleSlotFinalityStore interface {
	// WithContext sets the context for the store.
	WithContext(ctx context.Context) SingleSlotFinalityStore

	// SetSafeEth1BlockHash sets the safe block hash in the store.
	SetSafeEth1BlockHash(blockHash common.Hash)

	// GetSafeEth1BlockHash retrieves the safe block hash from the store.
	GetSafeEth1BlockHash() common.Hash

	// SetFinalizedEth1BlockHash sets the finalized block hash in the store.
	SetFinalizedEth1BlockHash(blockHash common.Hash)
	// GetFinalizedEth1BlockHash retrieves the finalized block hash from the
	// store.
	GetFinalizedEth1BlockHash() common.Hash

	// GenesisEth1Hash retrieves the Ethereum 1 genesis hash from the
	// store.
	GenesisEth1Hash() common.Hash

	// SetGenesisEth1Hash sets the Ethereum 1 genesis hash in the store.
	SetGenesisEth1Hash(common.Hash)
}
