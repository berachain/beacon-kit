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
	"github.com/itsdevbear/bolaris/beacon/forkchoice"
)

// ForkChoice represents the single-slot finality forkchoice algoritmn.
type ForkChoice struct {
	kv SingleSlotFinalityStore

	// latestBeaconBlock is the latest beacon block.
	latestBeaconBlock [32]byte
}

// New creates a new forkchoice instance.
func New(kv SingleSlotFinalityStore) *ForkChoice {
	return &ForkChoice{
		kv: kv,
	}
}

// WithContext sets the context for the forkchoice.
func (f *ForkChoice) WithContext(ctx context.Context) forkchoice.ForkChoicer {
	f.kv = f.kv.WithContext(ctx)
	return f
}

// InsertNode inserts a new node into the forkchoice.
func (f *ForkChoice) InsertNode(
	hash common.Hash,
) error {
	// Since this is single slot finality, we can just set the safe and
	// finalized block hash to the same value immediately.
	f.kv.SetFinalizedEth1BlockHash(hash)
	f.kv.SetSafeEth1BlockHash(hash)
	return nil
}

// TODO: Should this live in Forkchoicer?
func (f *ForkChoice) HeadBeaconBlock() [32]byte {
	return f.latestBeaconBlock
}

// TODO: Should this live in Forkchoicer?
func (f *ForkChoice) UpdateHeadBeaconBlock(
	blockHash [32]byte,
) {
	f.latestBeaconBlock = blockHash
}

// JustifiedPayloadBlockHash returns the justified checkpoint.
func (f *ForkChoice) JustifiedPayloadBlockHash() common.Hash {
	return f.kv.GetSafeEth1BlockHash()
}

// FinalizedPayloadBlockHash returns the finalized checkpoint.
func (f *ForkChoice) FinalizedPayloadBlockHash() common.Hash {
	return f.kv.GetFinalizedEth1BlockHash()
}
