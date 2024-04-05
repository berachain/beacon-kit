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
	"context"

	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/light/mod/state/codec"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/statedb/collections/encoding"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Copy returns a copy of the Store.
func (s *StateDB) Copy() *StateDB {
	cctx, write := sdk.UnwrapSDKContext(s.ctx).CacheContext()
	ss := s.WithContext(cctx)
	ss.write = write
	return ss
}

// Context returns the context of the Store.
func (s *StateDB) Context() context.Context {
	return s.ctx
}

// WithContext returns a copy of the Store with the given context.
func (s *StateDB) WithContext(ctx context.Context) *StateDB {
	cpy := *s
	cpy.ctx = ctx
	return &cpy
}

// Save saves the Store.
func (s *StateDB) Save() {
	if s.write != nil {
		s.write()
	}
}

// Store is a wrapper around an sdk provides access to all beacon related data.
type StateDB struct {
	ctx   context.Context
	write func()
	// provider *provider.Provider

	// genesisValidatorsRootCodec is the codec for the genesis validators root.
	genesisValidatorsRootCodec codec.Codec[codec.None, [32]byte]

	// slotCodec is the codec for the slot.
	slotCodec codec.Codec[codec.None, uint64]

	// latestBeaconBlockHeaderCodec is the codec for the latest beacon block header.
	latestBeaconBlockHeaderCodec codec.Codec[codec.None, *primitives.BeaconBlockHeader]

	// blockRootsCodec is the codec for the block roots.
	blockRootsCodec codec.Codec[uint64, [32]byte]

	// stateRootsCodec is the codec for the state roots.
	stateRootsCodec codec.Codec[uint64, [32]byte]

	// eth1BlockHashCodec is the codec for the eth1 block hash.
	eth1BlockHashCodec codec.Codec[codec.None, [32]byte]

	// eth1DepositIndexCodec is the codec for the eth1 deposit index.
	eth1DepositIndexCodec codec.Codec[codec.None, uint64]

	// validatorIndexCodec is the codec for the validator index.
	validatorIndexCodec codec.Codec[codec.None, uint64]

	// validatorsCodec is the codec for the validators.

	// balancesCodec is the codec for the balances.
	balancesCodec codec.Codec[uint64, uint64]

	// depositQueueCodec is the codec for the deposit queue.
	// TODO: Implement queue codecs

	// withdrawalQueueCodec is the codec for the withdrawal queue.

	// randaoMixCodec is the codec for the randao mix.
	randaoMixCodec codec.Codec[uint64, [32]byte]

	// slashingsCodec is the codec for the slashings.
	slashingsCodec codec.Codec[uint64, uint64]

	// totalSlashingCodec is the codec for the total slashing.
	totalSlashingCodec codec.Codec[codec.None, uint64]
}

// Store creates a new instance of Store.
func NewStateDB(
// provider *provider.Provider,
) *StateDB {
	return &StateDB{
		// provider: provider,
		genesisValidatorsRootCodec: codec.Codec[codec.None, [32]byte]{
			Value: encoding.Bytes32ValueCodec{},
		},
		slotCodec: codec.Codec[codec.None, uint64]{
			Value: sdkcollections.Uint64Value,
		},
		blockRootsCodec: codec.Codec[uint64, [32]byte]{
			Key:   sdkcollections.Uint64Key,
			Value: encoding.Bytes32ValueCodec{},
		},
		stateRootsCodec: codec.Codec[uint64, [32]byte]{
			Key:   sdkcollections.Uint64Key,
			Value: encoding.Bytes32ValueCodec{},
		},
		eth1BlockHashCodec: codec.Codec[codec.None, [32]byte]{
			Value: encoding.Bytes32ValueCodec{},
		},
		eth1DepositIndexCodec: codec.Codec[codec.None, uint64]{
			Value: sdkcollections.Uint64Value,
		},
		validatorIndexCodec: codec.Codec[codec.None, uint64]{
			Value: sdkcollections.Uint64Value,
		},
		balancesCodec: codec.Codec[uint64, uint64]{
			Key:   sdkcollections.Uint64Key,
			Value: sdkcollections.Uint64Value,
		},
		randaoMixCodec: codec.Codec[uint64, [32]byte]{
			Key:   sdkcollections.Uint64Key,
			Value: encoding.Bytes32ValueCodec{},
		},
		slashingsCodec: codec.Codec[uint64, uint64]{
			Key:   sdkcollections.Uint64Key,
			Value: sdkcollections.Uint64Value,
		},
		totalSlashingCodec: codec.Codec[codec.None, uint64]{
			Value: sdkcollections.Uint64Value,
		},
		//nolint:lll
		latestBeaconBlockHeaderCodec: codec.Codec[codec.None, *primitives.BeaconBlockHeader]{
			Value: encoding.SSZValueCodec[*primitives.BeaconBlockHeader]{},
		},
	}
}
