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

// ForkchoiceKVStore is the interface for the forkchoice KV store.
type ForkchoiceKVStore interface {
	WithContext(ctx context.Context) ForkchoiceKVStore

	// TODO: eventually the forkchoicer shouldn't have
	// anything to do with beacon, since this forkchoicer
	// is only needed for finalizing the eth1 blocks.
	SetLastSeenBeaconBlock(blockHash [32]byte)
	GetLastSeenBeaconBlock() [32]byte

	// SetSafeEth1BlockHash sets the safe block hash in the store.
	SetSafeEth1BlockHash(blockHash common.Hash)
	// GetSafeEth1BlockHash retrieves the safe block hash from the store.
	GetSafeEth1BlockHash() common.Hash

	// SetFinalizedEth1BlockHash sets the finalized block hash in the store.
	SetFinalizedEth1BlockHash(blockHash common.Hash)
	// GetFinalizedEth1BlockHash retrieves the finalized block hash from the
	// store.
	GetFinalizedEth1BlockHash() common.Hash

	GenesisEth1Hash() common.Hash
	SetGenesisEth1Hash(common.Hash)
}

// ForkChoice represents an extremely simple single slot finality
// forkchoice ruler for the execution chain.
type ForkChoice struct {
	kv ForkchoiceKVStore
}

func New(kv ForkchoiceKVStore) *ForkChoice {
	return &ForkChoice{
		kv: kv,
	}
}

func (f *ForkChoice) WithContext(ctx context.Context) *ForkChoice {
	f.kv = f.kv.WithContext(ctx)
	return f
}

func (f *ForkChoice) HeadBeaconBlock() [32]byte {
	return f.kv.GetLastSeenBeaconBlock()
}

func (f *ForkChoice) UpdateHeadBeaconBlock(
	blockHash [32]byte,
) {
	f.kv.SetLastSeenBeaconBlock(blockHash)
}

func (f *ForkChoice) UpdateFinalizedCheckpoint(
	hash common.Hash,
) error {
	f.kv.SetFinalizedEth1BlockHash(hash)
	return nil
}

// UpdateFinalizedCheckpoint updates the finalized
// checkpoint in the forkchoice.
func (f *ForkChoice) UpdateJustifiedCheckpoint(
	hash common.Hash,
) error {
	f.kv.SetSafeEth1BlockHash(hash)
	return nil
}

func (f *ForkChoice) JustifiedCheckpoint() common.Hash {
	return f.kv.GetSafeEth1BlockHash()
}

func (f *ForkChoice) FinalizedCheckpoint() common.Hash {
	return f.kv.GetFinalizedEth1BlockHash()
}
