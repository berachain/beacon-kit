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

package kv

import (
	"cosmossdk.io/store"
	"github.com/ethereum/go-ethereum/common"
)

// ForkChoiceStore represents the fork choice rule in the blockchain.
type ForkChoiceStore struct {
	store store.KVStore

	lastValidHead common.Hash
}

// NewForkchoice creates a new instance of ForkChoiceStore.
func NewForkchoice(store store.KVStore) *ForkChoiceStore {
	return &ForkChoiceStore{
		store: store,
	}
}

// SetSafeEth1BlockHash sets the safe block hash in the store.
func (f *ForkChoiceStore) SetSafeEth1BlockHash(safeBlockHash common.Hash) {
	f.store.Set([]byte("forkchoice_safe"), safeBlockHash[:])
}

// GetSafeEth1BlockHash retrieves the safe block hash from the store.
func (f *ForkChoiceStore) GetSafeEth1BlockHash() common.Hash {
	bz := f.store.Get([]byte("forkchoice_safe"))
	if bz == nil {
		return common.Hash{}
	}
	var safeBlockHash common.Hash
	copy(safeBlockHash[:], bz)
	return safeBlockHash
}

// SetFinalizedEth1BlockHash sets the finalized block hash in the store.
func (f *ForkChoiceStore) SetFinalizedEth1BlockHash(finalizedBlockHash common.Hash) {
	f.store.Set([]byte("forkchoice_finalized"), finalizedBlockHash[:])
}

// GetFinalizedEth1BlockHash retrieves the finalized block hash from the store.
func (f *ForkChoiceStore) GetFinalizedEth1BlockHash() common.Hash {
	bz := f.store.Get([]byte("forkchoice_finalized"))
	if bz == nil {
		return common.Hash{}
	}
	var finalizedBlockHash common.Hash
	copy(finalizedBlockHash[:], bz)
	return finalizedBlockHash
}

// SetLastValidHead sets the last valid head in the store.
func (f *ForkChoiceStore) SetLastValidHead(lastValidHead common.Hash) {
	f.lastValidHead = lastValidHead
}

// GetLastValidHead retrieves the last valid head from the store.
func (f *ForkChoiceStore) GetLastValidHead() common.Hash {
	return f.lastValidHead
}
