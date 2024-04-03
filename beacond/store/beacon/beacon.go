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

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/beacond/store/beacon/collections"
	"github.com/berachain/beacon-kit/beacond/store/beacon/collections/encoding"
	"github.com/berachain/beacon-kit/beacond/store/beacon/index"
	"github.com/berachain/beacon-kit/beacond/store/beacon/keys"
	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/state"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Store is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type Store struct {
	ctx   context.Context
	write func()
	cfg   *params.BeaconChainConfig

	// genesisValidatorsRoot is the root of the genesis validators.
	genesisValidatorsRoot sdkcollections.Item[[32]byte]

	// slot is the current slot.
	slot sdkcollections.Item[uint64]

	// latestBeaconBlockHeader stores the latest beacon block header.
	latestBeaconBlockHeader sdkcollections.Item[*beacontypes.BeaconBlockHeader]

	// blockRoots stores the block roots for the current epoch.
	blockRoots sdkcollections.Map[uint64, [32]byte]

	// stateRoots stores the state roots for the current epoch.
	stateRoots sdkcollections.Map[uint64, [32]byte]

	// eth1BlockHash stores the block hash of the latest eth1 block.
	eth1BlockHash sdkcollections.Item[[32]byte]

	// eth1DepositIndex is the index of the latest eth1 deposit.
	eth1DepositIndex sdkcollections.Item[uint64]

	// validatorIndex is a sequence that provides the next
	// available index for a new validator.
	validatorIndex sdkcollections.Sequence

	// validators stores the list of validators.
	validators *sdkcollections.IndexedMap[
		uint64, *beacontypes.Validator, index.ValidatorsIndex,
	]

	// balances stores the list of balances.
	balances sdkcollections.Map[uint64, uint64]

	// depositQueue is a list of deposits that are queued to be processed.
	depositQueue *collections.Queue[*beacontypes.Deposit]

	// withdrawalQueue is a list of withdrawals that are queued to be processed.
	withdrawalQueue *collections.Queue[*primitives.Withdrawal]

	// randaoMix stores the randao mix for the current epoch.
	randaoMix sdkcollections.Map[uint64, [32]byte]

	// slashings stores the slashings for the current epoch.
	slashings sdkcollections.Map[uint64, uint64]

	// totalSlashing stores the total slashing in the vector range.
	totalSlashing sdkcollections.Item[uint64]
}

// Store creates a new instance of Store.
//
//nolint:funlen // its not overly complex.
func NewStore(
	env appmodule.Environment,
	cfg *params.BeaconChainConfig,
) *Store {
	schemaBuilder := sdkcollections.NewSchemaBuilder(env.KVStoreService)
	return &Store{
		ctx: nil,
		cfg: cfg,
		genesisValidatorsRoot: sdkcollections.NewItem[[32]byte](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.GenesisValidatorsRootPrefix),
			keys.GenesisValidatorsRootPrefix,
			encoding.Bytes32ValueCodec{},
		),
		slot: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.SlotPrefix),
			keys.SlotPrefix,
			sdkcollections.Uint64Value,
		),
		blockRoots: sdkcollections.NewMap[uint64, [32]byte](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.BlockRootsPrefix),
			keys.BlockRootsPrefix,
			sdkcollections.Uint64Key,
			encoding.Bytes32ValueCodec{},
		),
		stateRoots: sdkcollections.NewMap[uint64, [32]byte](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.StateRootsPrefix),
			keys.StateRootsPrefix,
			sdkcollections.Uint64Key,
			encoding.Bytes32ValueCodec{},
		),
		eth1BlockHash: sdkcollections.NewItem[[32]byte](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.Eth1BlockHashPrefix),
			keys.Eth1BlockHashPrefix,
			encoding.Bytes32ValueCodec{},
		),
		eth1DepositIndex: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.Eth1DepositIndexPrefix),
			keys.Eth1DepositIndexPrefix,
			sdkcollections.Uint64Value,
		),
		validatorIndex: sdkcollections.NewSequence(
			schemaBuilder,
			sdkcollections.NewPrefix(keys.ValidatorIndexPrefix),
			keys.ValidatorIndexPrefix,
		),
		validators: sdkcollections.NewIndexedMap[
			uint64, *beacontypes.Validator,
		](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.ValidatorByIndexPrefix),
			keys.ValidatorByIndexPrefix,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[*beacontypes.Validator]{},
			index.NewValidatorsIndex(schemaBuilder),
		),
		balances: sdkcollections.NewMap[uint64, uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.BalancesPrefix),
			keys.BalancesPrefix,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		depositQueue: collections.NewQueue[*beacontypes.Deposit](
			schemaBuilder,
			keys.DepositQueuePrefix,
			encoding.SSZValueCodec[*beacontypes.Deposit]{},
		),
		withdrawalQueue: collections.NewQueue[*primitives.Withdrawal](
			schemaBuilder,
			keys.WithdrawalQueuePrefix,
			encoding.SSZValueCodec[*primitives.Withdrawal]{},
		),
		randaoMix: sdkcollections.NewMap[uint64, [32]byte](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.RandaoMixPrefix),
			keys.RandaoMixPrefix,
			sdkcollections.Uint64Key,
			encoding.Bytes32ValueCodec{},
		),
		slashings: sdkcollections.NewMap[uint64, uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.SlashingsPrefix),
			keys.SlashingsPrefix,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		totalSlashing: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.TotalSlashingPrefix),
			keys.TotalSlashingPrefix,
			sdkcollections.Uint64Value,
		),
		//nolint:lll
		latestBeaconBlockHeader: sdkcollections.NewItem[*beacontypes.BeaconBlockHeader](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.LatestBeaconBlockHeaderPrefix),
			keys.LatestBeaconBlockHeaderPrefix,
			encoding.SSZValueCodec[*beacontypes.BeaconBlockHeader]{},
		),
	}
}

// Copy returns a copy of the Store.
func (s *Store) Copy() state.BeaconState {
	cctx, write := sdk.UnwrapSDKContext(s.ctx).CacheContext()
	ss := s.WithContext(cctx)
	ss.write = write
	return ss
}

// Context returns the context of the Store.
func (s *Store) Context() context.Context {
	return s.ctx
}

// WithContext returns a copy of the Store with the given context.
func (s *Store) WithContext(ctx context.Context) *Store {
	cpy := *s
	cpy.ctx = ctx
	return &cpy
}

// Save saves the Store.
func (s *Store) Save() {
	if s.write != nil {
		s.write()
	}
}
