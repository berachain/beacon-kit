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

// UpdateEth1BlockHash sets the Eth1 hash in the BeaconStore.
func (s *StateDB) UpdateEth1BlockHash(
	hash primitives.ExecutionHash,
) error {
	return s.eth1BlockHash.Set(s.ctx, hash)
}

// GetEth1Hash retrieves the Eth1 hash from the BeaconStore.
func (s *StateDB) GetEth1BlockHash() (primitives.ExecutionHash, error) {
	return s.eth1BlockHash.Get(s.ctx)
}

// GetEth1DepositIndex retrieves the eth1 deposit index from the beacon state.
func (s *StateDB) GetEth1DepositIndex() (uint64, error) {
	return s.eth1DepositIndex.Get(s.ctx)
}

// SetEth1DepositIndex sets the eth1 deposit index in the beacon state.
func (s *StateDB) SetEth1DepositIndex(index uint64) error {
	return s.eth1DepositIndex.Set(s.ctx, index)
}
