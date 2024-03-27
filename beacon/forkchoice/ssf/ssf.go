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

	"github.com/berachain/beacon-kit/primitives"
)

// ForkChoice represents the single-slot finality forkchoice algoritmn.
type ForkChoice struct {
	kv SingleSlotFinalityStore
}

// New creates a new forkchoice instance.
func New(kv SingleSlotFinalityStore) *ForkChoice {
	return &ForkChoice{
		kv: kv,
	}
}

// SetContext sets the context for the forkchoice.
func (f *ForkChoice) SetContext(ctx context.Context) {
	if f == nil || f.kv == nil {
		panic("ssf: uninitialized forkchoice instance")
	} else if ctx == nil {
		panic("ssf: nil context")
	}
	f.kv.WithContext(ctx)
}

// InsertNode inserts a new node into the forkchoice.
func (f ForkChoice) InsertNode(
	hash primitives.ExecutionHash,
) error {
	// Since this is single slot finality, we can just set the safe and
	// finalized block hash to the same value immediately.
	if f.kv == nil {
		panic("ssf: uninitialized forkchoice instance")
	}
	f.kv.SetFinalizedEth1BlockHash(hash)
	f.kv.SetSafeEth1BlockHash(hash)
	return nil
}

// JustifiedPayloadBlockHash returns the justified checkpoint.
func (f *ForkChoice) JustifiedPayloadBlockHash() primitives.ExecutionHash {
	return f.kv.GetSafeEth1BlockHash()
}

// FinalizedPayloadBlockHash returns the finalized checkpoint.
func (f *ForkChoice) FinalizedPayloadBlockHash() primitives.ExecutionHash {
	return f.kv.GetFinalizedEth1BlockHash()
}
