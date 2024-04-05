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

package beacon

// import (
// 	"context"

// 	sdkcollections "cosmossdk.io/collections"
// 	"github.com/berachain/beacon-kit/beacon/core/randao/types"
// 	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
// 	"github.com/berachain/beacon-kit/lib/store/collections/encoding"
// 	"github.com/berachain/beacon-kit/light/provider"
// 	"github.com/berachain/beacon-kit/light/store"
// )

// // Store is a wrapper around an sdk provides access to all beacon related data.
// type Store struct {
// 	ctx      context.Context
// 	provider *provider.Provider

// 	// genesisValidatorsRootCodec is the codec for the genesis validators root.
// 	genesisValidatorsRootCodec store.Codec[store.None, [32]byte]

// 	// slotCodec is the codec for the slot.
// 	slotCodec store.Codec[store.None, uint64]

// 	// latestBeaconBlockHeaderCodec is the codec for the latest beacon block header.
// 	latestBeaconBlockHeaderCodec store.Codec[store.None, *beacontypes.BeaconBlockHeader]

// 	// blockRootsCodec is the codec for the block roots.
// 	blockRootsCodec store.Codec[uint64, [32]byte]

// 	// stateRootsCodec is the codec for the state roots.
// 	stateRootsCodec store.Codec[uint64, [32]byte]

// 	// eth1DepositIndexCodec is the codec for the eth1 deposit index.
// 	eth1DepositIndexCodec store.Codec[store.None, uint64]

// 	// validatorIndexCodec is the codec for the validator index.
// 	validatorIndexCodec store.Codec[store.None, uint64]

// 	// validatorsCodec is the codec for the validators.

// 	// balancesCodec is the codec for the balances.
// 	balancesCodec store.Codec[uint64, uint64]

// 	// depositQueueCodec is the codec for the deposit queue.
// 	// TODO: Implement queue codecs

// 	// withdrawalQueueCodec is the codec for the withdrawal queue.

// 	// randaoMixCodec is the codec for the randao mix.
// 	randaoMixCodec store.Codec[uint64, [types.MixLength]byte]
// }

// // Store creates a new instance of Store.
// func NewStore(
// 	provider *provider.Provider,
// ) *Store {
// 	return &Store{
// 		provider: provider,
// 		genesisValidatorsRootCodec: store.Codec[store.None, [32]byte]{
// 			Value: encoding.Bytes32ValueCodec{},
// 		},
// 		slotCodec: store.Codec[store.None, uint64]{
// 			Value: sdkcollections.Uint64Value,
// 		},
// 		blockRootsCodec: store.Codec[uint64, [32]byte]{
// 			Key:   sdkcollections.Uint64Key,
// 			Value: encoding.Bytes32ValueCodec{},
// 		},
// 		stateRootsCodec: store.Codec[uint64, [32]byte]{
// 			Key:   sdkcollections.Uint64Key,
// 			Value: encoding.Bytes32ValueCodec{},
// 		},
// 		eth1DepositIndexCodec: store.Codec[store.None, uint64]{
// 			Value: sdkcollections.Uint64Value,
// 		},
// 		validatorIndexCodec: store.Codec[store.None, uint64]{
// 			Value: sdkcollections.Uint64Value,
// 		},
// 		balancesCodec: store.Codec[uint64, uint64]{
// 			Key:   sdkcollections.Uint64Key,
// 			Value: sdkcollections.Uint64Value,
// 		},
// 		randaoMixCodec: store.Codec[uint64, [types.MixLength]byte]{
// 			Key:   sdkcollections.Uint64Key,
// 			Value: encoding.Bytes32ValueCodec{},
// 		},
// 		//nolint:lll
// 		latestBeaconBlockHeaderCodec: store.Codec[store.None, *beacontypes.BeaconBlockHeader]{
// 			Value: encoding.SSZValueCodec[*beacontypes.BeaconBlockHeader]{},
// 		},
// 	}
// }

// // Context returns the context of the Store.
// func (s *Store) Context() context.Context {
// 	return s.ctx
// }

// // WithContext returns a copy of the Store with the given context.
// func (s *Store) WithContext(ctx context.Context) *Store {
// 	cpy := *s
// 	cpy.ctx = ctx
// 	return &cpy
// }
