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

package store

import (
	"github.com/ethereum/go-ethereum/common"
)

// SetLastValidHead sets the last valid head in the store.
// TODO: Make this in-mem thing more robust.
func (s *BeaconStore) SetLastValidHead(lastValidHead common.Hash) {
	if s.lastValidHash == nil {
		s.lastValidHash = new(common.Hash)
	}
	*s.lastValidHash = lastValidHead
}

// GetLastValidHead retrieves the last valid head from the store.
// TODO: Make this in-mem thing more robust.
func (s *BeaconStore) GetLastValidHead() common.Hash {
	if s.lastValidHash == nil {
		return s.GetSafeEth1BlockHash()
	}
	return *s.lastValidHash
}

// SetSafeEth1BlockHash sets the safe block hash in the store.
func (s *BeaconStore) SetSafeEth1BlockHash(blockHash common.Hash) {
	if err := s.fcSafeEth1BlockHash.Set(s.ctx, blockHash); err != nil {
		panic(err)
	}
}

// GetSafeEth1BlockHash retrieves the safe block hash from the store.
func (s *BeaconStore) GetSafeEth1BlockHash() common.Hash {
	safeHash, err := s.fcSafeEth1BlockHash.Get(s.ctx)
	if err != nil {
		safeHash = common.Hash{}
	}
	return safeHash
}

// SetFinalizedEth1BlockHash sets the finalized block hash in the store.
func (s *BeaconStore) SetFinalizedEth1BlockHash(blockHash common.Hash) {
	if err := s.fcFinalizedEth1BlockHash.Set(s.ctx, blockHash); err != nil {
		panic(err)
	}
}

// GetFinalizedEth1BlockHash retrieves the finalized block hash from the store.
func (s *BeaconStore) GetFinalizedEth1BlockHash() common.Hash {
	finalHash, err := s.fcFinalizedEth1BlockHash.Get(s.ctx)
	if err != nil {
		finalHash = common.Hash{}
	}
	return finalHash
}
