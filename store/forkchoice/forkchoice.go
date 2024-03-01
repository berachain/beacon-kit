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

package forkchoice

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/lib/store/collections/encoding"
)

type Store struct {
	ctx context.Context

	// fcHeadEth1BlockHash is the head block hash.
	fcHeadEth1BlockHash *common.Hash

	// fcSafeEth1BlockHash is the safe block hash.
	fcSafeEth1BlockHash sdkcollections.Item[[32]byte]

	// fcFinalizedEth1BlockHash is the finalized block hash.
	fcFinalizedEth1BlockHash sdkcollections.Item[[32]byte]

	// eth1GenesisHash is the Eth1 genesis hash.
	eth1GenesisHash sdkcollections.Item[[32]byte]
}

func NewStore(
	kvs store.KVStoreService,
) *Store {
	kvSchemaBuilder := sdkcollections.NewSchemaBuilder(kvs)

	fcSafeEth1BlockHash := sdkcollections.NewItem[[32]byte](
		kvSchemaBuilder,
		sdkcollections.NewPrefix(fcSafeEth1BlockHashPrefix),
		fcSafeEth1BlockHashPrefix,
		encoding.Bytes32ValueCodec{},
	)
	fcFinalizedEth1BlockHash := sdkcollections.NewItem[[32]byte](
		kvSchemaBuilder,
		sdkcollections.NewPrefix(fcFinalizedEth1BlockHashPrefix),
		fcFinalizedEth1BlockHashPrefix,
		encoding.Bytes32ValueCodec{},
	)
	eth1GenesisHash := sdkcollections.NewItem[[32]byte](
		kvSchemaBuilder,
		sdkcollections.NewPrefix(eth1GenesisHashPrefix),
		eth1GenesisHashPrefix,
		encoding.Bytes32ValueCodec{},
	)

	return &Store{
		fcHeadEth1BlockHash:      nil,
		fcSafeEth1BlockHash:      fcSafeEth1BlockHash,
		fcFinalizedEth1BlockHash: fcFinalizedEth1BlockHash,
		eth1GenesisHash:          eth1GenesisHash,
	}
}

// SetLastValidHead sets the last valid head in the store.
// TODO: Make this in-mem thing more robust.
func (s *Store) SetLastValidHead(blockHash common.Hash) {
	s.fcHeadEth1BlockHash = &blockHash
}

// GetHeadHash retrieves the last valid head from the store.
// TODO: Make this in-mem thing more robust.
func (s *Store) GetLastValidHead() common.Hash {
	if s.fcHeadEth1BlockHash == nil {
		return common.Hash{}
	}
	return *s.fcHeadEth1BlockHash
}

// SetSafeEth1BlockHash sets the safe block hash in the store.
func (s *Store) SetSafeEth1BlockHash(blockHash common.Hash) {
	if err := s.fcSafeEth1BlockHash.Set(s.ctx, blockHash); err != nil {
		panic(err)
	}
}

// GetSafeEth1BlockHash retrieves the safe block hash from the store.
func (s *Store) GetSafeEth1BlockHash() common.Hash {
	safeHash, err := s.fcSafeEth1BlockHash.Get(s.ctx)
	if err != nil {
		safeHash = common.Hash{}
	}
	return safeHash
}

// SetFinalizedEth1BlockHash sets the finalized block hash in the store.
func (s *Store) SetFinalizedEth1BlockHash(blockHash common.Hash) {
	if err := s.fcFinalizedEth1BlockHash.Set(s.ctx, blockHash); err != nil {
		panic(err)
	}
}

// GetFinalizedEth1BlockHash retrieves the finalized block hash from the store.
func (s *Store) GetFinalizedEth1BlockHash() common.Hash {
	finalHash, err := s.fcFinalizedEth1BlockHash.Get(s.ctx)
	if err != nil {
		finalHash = common.Hash{}
	}
	return finalHash
}

// SetGenesisEth1Hash sets the Ethereum 1 genesis hash in the BeaconStore.
func (s *Store) SetGenesisEth1Hash(eth1GenesisHash common.Hash) {
	if err := s.eth1GenesisHash.Set(s.ctx, eth1GenesisHash); err != nil {
		panic(err)
	}
}

// GenesisEth1Hash retrieves the Ethereum 1 genesis hash from the BeaconStore.
func (s *Store) GenesisEth1Hash() common.Hash {
	genesisHash, err := s.eth1GenesisHash.Get(s.ctx)
	if err != nil {
		panic("failed to get genesis eth1hash")
	}
	return genesisHash
}

// WithContext( returns the Store with the given context.
func (s *Store) WithContext(ctx context.Context) *Store {
	s.ctx = ctx
	return s
}
