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

// import (
// 	"context"

// 	sdkcollections "cosmossdk.io/collections"
// 	"cosmossdk.io/core/store"
// 	"github.com/berachain/beacon-kit/beacon/forkchoice/ssf"
// 	"github.com/berachain/beacon-kit/lib/store/collections/encoding"
// 	"github.com/berachain/beacon-kit/primitives"
// )

// // TODO: Decouple from the Specific SingleSlotFinalityStore Impl.
// var _ ssf.SingleSlotFinalityStore = &Store{}

// type Store struct {
// 	ctx context.Context

// 	// fcSafeEth1BlockHash is the safe block hash.
// 	fcSafeEth1BlockHash sdkcollections.Item[[32]byte]

// 	// fcFinalizedEth1BlockHash is the finalized block hash.
// 	fcFinalizedEth1BlockHash sdkcollections.Item[[32]byte]

// 	// eth1GenesisHash is the Eth1 genesis hash.
// 	eth1GenesisHash sdkcollections.Item[[32]byte]
// }

// func NewStore(
// 	kvs store.KVStoreService,
// ) *Store {
// 	kvSchemaBuilder := sdkcollections.NewSchemaBuilder(kvs)

// 	fcSafeEth1BlockHash := sdkcollections.NewItem[[32]byte](
// 		kvSchemaBuilder,
// 		sdkcollections.NewPrefix(fcSafeEth1BlockHashPrefix),
// 		fcSafeEth1BlockHashPrefix,
// 		encoding.Bytes32ValueCodec{},
// 	)
// 	fcFinalizedEth1BlockHash := sdkcollections.NewItem[[32]byte](
// 		kvSchemaBuilder,
// 		sdkcollections.NewPrefix(fcFinalizedEth1BlockHashPrefix),
// 		fcFinalizedEth1BlockHashPrefix,
// 		encoding.Bytes32ValueCodec{},
// 	)
// 	eth1GenesisHash := sdkcollections.NewItem[[32]byte](
// 		kvSchemaBuilder,
// 		sdkcollections.NewPrefix(eth1GenesisHashPrefix),
// 		eth1GenesisHashPrefix,
// 		encoding.Bytes32ValueCodec{},
// 	)

// 	return &Store{
// 		fcSafeEth1BlockHash:      fcSafeEth1BlockHash,
// 		fcFinalizedEth1BlockHash: fcFinalizedEth1BlockHash,
// 		eth1GenesisHash:          eth1GenesisHash,
// 	}
// }

// // SetSafeEth1BlockHash sets the safe block hash in the store.
// func (s *Store) SetSafeEth1BlockHash(blockHash primitives.ExecutionHash) {}

// // GetSafeEth1BlockHash retrieves the sprimitives.ExecutionHashash from the
// // store.
// func (s *Store) GetSafeEth1BlockHash() primitives.ExecutionHash {
// 	panic("not implemented")
// }

// // SetFinalizedEth1BlockHash sets the finalized block hash in the store.
// func (s *Store) SetFinalizedEth1BlockHash(blockHash primitives.ExecutionHash) {}

// // GetFinalizedEth1BlockHash retrieves the finalized block hash from the store.
// func (s *Store) GetFinalizedEth1BlockHash() primitives.ExecutionHash {
// 	panic("not implemented")
// }

// // SetGenesisEth1Hash sets the Ethereum 1 genesis hash in the BeaconStore.
// func (s *Store) SetGenesisEth1Hash(eth1GenesisHash primitives.ExecutionHash) {}

// // GenesisEth1Hash retrieves the Ethereum 1 genesis hash from the BeaconStore.
// func (s *Store) GenesisEth1Hash() primitives.ExecutionHash {
// 	panic("not implemented")
// }

// // WithContext returns the Store with the given context.
// func (s *Store) WithContext(ctx context.Context) ssf.SingleSlotFinalityStore {
// 	panic("not implemented")
// }
