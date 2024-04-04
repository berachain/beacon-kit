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

package statedb

import (
	"github.com/berachain/beacon-kit/mod/primitives"
)

// UpdateBlockRootAtIndex sets a block root in the BeaconStore.
func (s *StateDB) UpdateBlockRootAtIndex(
	index uint64,
	root primitives.Root,
) error {
	return s.blockRoots.Set(s.ctx, index, root)
}

// GetBlockRoot retrieves the block root from the BeaconStore.
func (s *StateDB) GetBlockRootAtIndex(
	index uint64,
) (primitives.Root, error) {
	return s.blockRoots.Get(s.ctx, index)
}

// SetLatestBlockHeader sets the latest block header in the BeaconStore.
func (s *StateDB) SetLatestBlockHeader(
	header *primitives.BeaconBlockHeader,
) error {
	return s.latestBlockHeader.Set(s.ctx, header)
}

// GetLatestBlockHeader retrieves the latest block header from the BeaconStore.
func (s *StateDB) GetLatestBlockHeader() (
	*primitives.BeaconBlockHeader, error,
) {
	return s.latestBlockHeader.Get(s.ctx)
}

// UpdateStateRootAtIndex updates the state root at the given slot.
func (s *StateDB) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot primitives.Root,
) error {
	return s.stateRoots.Set(s.ctx, idx, stateRoot)
}

// StateRootAtIndex returns the state root at the given slot.
func (s *StateDB) StateRootAtIndex(idx uint64) (primitives.Root, error) {
	return s.stateRoots.Get(s.ctx, idx)
}
