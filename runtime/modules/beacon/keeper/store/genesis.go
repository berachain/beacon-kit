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

// SetGenesisEth1Hash sets the Ethereum 1 genesis hash in the BeaconStore.
func (s *BeaconStore) SetGenesisEth1Hash(eth1GenesisHash common.Hash) {
	if err := s.eth1GenesisHash.Set(s.ctx, eth1GenesisHash); err != nil {
		panic(err)
	}
}

// GenesisEth1Hash retrieves the Ethereum 1 genesis hash from the BeaconStore.
func (s *BeaconStore) GenesisEth1Hash() common.Hash {
	genesisHash, err := s.eth1GenesisHash.Get(s.ctx)
	if err != nil {
		panic("failed to get genesis eth1hash")
	}
	return genesisHash
}
