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

package beacondb

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/light/mod/provider"
	"github.com/berachain/beacon-kit/light/mod/storage/codec"
	"github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections/encoding"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KVStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type KVStore struct {
	ctx      context.Context
	write    func()
	provider *provider.Provider

	// Versioning
	// genesisValidatorsRoot is the root of the genesis validators.
	genesisValidatorsRoot codec.Item[[32]byte]
	// slot is the current slot.
	slot codec.Item[uint64]
	// fork is the current fork
	fork codec.Item[*primitives.Fork]

	// History
	// latestBlockHeader stores the latest beacon block header.
	latestBlockHeader codec.Item[*primitives.BeaconBlockHeader]
	// blockRoots stores the block roots for the current epoch.
	blockRoots codec.Map[uint64, [32]byte]
	// stateRoots stores the state roots for the current epoch.
	stateRoots codec.Map[uint64, [32]byte]

	// Eth1
	// latestExecutionPayload stores the latest execution payload.
	latestExecutionPayload codec.Item[types.ExecutionPayload]

	// eth1Data stores the latest eth1 data.
	eth1Data codec.Item[*primitives.Eth1Data]
	// eth1DepositIndex is the index of the latest eth1 deposit.
	eth1DepositIndex codec.Item[uint64]

	// Registry
	// // validatorIndex provides the next available index for a new validator.
	// validatorIndex sdkcollections.Sequence
	// // validators stores the list of validators.
	// validators *sdkcollections.IndexedMap[
	// 	uint64, *beacontypes.Validator, index.ValidatorsIndex,
	// ]
	// balances stores the list of balances.
	balances codec.Map[uint64, uint64]

	// depositQueue is a list of deposits that are queued to be processed.
	depositQueue *collections.Queue[*primitives.Deposit]

	// // withdrawalQueue is a list of withdrawals that are queued to be processed.
	// withdrawalQueue *collections.Queue[*primitives.Withdrawal]

	// nextWithdrawalIndex stores the next global withdrawal index.
	nextWithdrawalIndex codec.Item[uint64]

	// nextWithdrawalValidatorIndex stores the next withdrawal validator index
	// for each validator.
	nextWithdrawalValidatorIndex codec.Item[uint64]

	// Randomness
	// randaoMix stores the randao mix for the current epoch.
	randaoMix codec.Map[uint64, [32]byte]

	// Slashings
	// slashings stores the slashings for the current epoch.
	slashings codec.Map[uint64, uint64]
	// totalSlashing stores the total slashing in the vector range.
	totalSlashing codec.Item[uint64]
}

// Store creates a new instance of Store.
//
//nolint:funlen // its not overly complex.
func New(
	provider *provider.Provider,
) *KVStore {
	return &KVStore{
		ctx:      nil,
		provider: provider,
		genesisValidatorsRoot: codec.NewItem[[32]byte](
			keys.GenesisValidatorsRootPrefix,
			encoding.Bytes32ValueCodec{},
		),
		slot: codec.NewItem[uint64](
			keys.SlotPrefix,
			sdkcollections.Uint64Value,
		),
		fork: codec.NewItem[*primitives.Fork](
			keys.ForkPrefix,
			encoding.SSZValueCodec[*primitives.Fork]{},
		),
		blockRoots: codec.NewMap[uint64, [32]byte](
			keys.BlockRootsPrefix,
			sdkcollections.Uint64Key,
			encoding.Bytes32ValueCodec{},
		),
		stateRoots: codec.NewMap[uint64, [32]byte](
			keys.StateRootsPrefix,
			sdkcollections.Uint64Key,
			encoding.Bytes32ValueCodec{},
		),
		latestExecutionPayload: codec.NewItem[types.ExecutionPayload](
			keys.LatestExecutionPayloadPrefix,
			encoding.SSZInterfaceCodec[types.ExecutionPayload]{
				Factory: func() types.ExecutionPayload {
					return &types.ExecutableDataDeneb{}
				},
			},
		),
		eth1Data: codec.NewItem[*primitives.Eth1Data](
			keys.Eth1DataPrefix,
			encoding.SSZValueCodec[*primitives.Eth1Data]{},
		),
		eth1DepositIndex: codec.NewItem[uint64](
			keys.Eth1DepositIndexPrefix,
			sdkcollections.Uint64Value,
		),
		// validatorIndex: sdkcollections.NewSequence(
		// 	schemaBuilder,
		// 	sdkcollections.NewPrefix(keys.ValidatorIndexPrefix),
		// 	keys.ValidatorIndexPrefix,
		// ),
		// validators: sdkcollections.NewIndexedMap[
		// 	uint64, *beacontypes.Validator,
		// ](
		// 	schemaBuilder,
		// 	sdkcollections.NewPrefix(keys.ValidatorByIndexPrefix),
		// 	keys.ValidatorByIndexPrefix,
		// 	sdkcollections.Uint64Key,
		// 	encoding.SSZValueCodec[*beacontypes.Validator]{},
		// 	index.NewValidatorsIndex(schemaBuilder),
		// ),
		balances: codec.NewMap[uint64, uint64](
			keys.BalancesPrefix,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		// depositQueue: collections.NewQueue[*primitives.Deposit](
		// 	schemaBuilder,
		// 	keys.DepositQueuePrefix,
		// 	encoding.SSZValueCodec[*primitives.Deposit]{},
		// ),
		// withdrawalQueue: collections.NewQueue[*primitives.Withdrawal](
		// 	schemaBuilder,
		// 	keys.WithdrawalQueuePrefix,
		// 	encoding.SSZValueCodec[*primitives.Withdrawal]{},
		// ),
		randaoMix: codec.NewMap[uint64, [32]byte](
			keys.RandaoMixPrefix,
			sdkcollections.Uint64Key,
			encoding.Bytes32ValueCodec{},
		),
		slashings: codec.NewMap[uint64, uint64](
			keys.SlashingsPrefix,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalIndex: codec.NewItem[uint64](
			keys.NextWithdrawalIndexPrefix,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalValidatorIndex: codec.NewItem[uint64](
			keys.NextWithdrawalValidatorIndexPrefix,
			sdkcollections.Uint64Value,
		),
		totalSlashing: codec.NewItem[uint64](
			keys.TotalSlashingPrefix,
			sdkcollections.Uint64Value,
		),

		latestBlockHeader: codec.NewItem[*primitives.BeaconBlockHeader](
			keys.LatestBeaconBlockHeaderPrefix,
			encoding.SSZValueCodec[*primitives.BeaconBlockHeader]{},
		),
	}
}

// Copy returns a copy of the Store.
func (kv *KVStore) Copy() *KVStore {
	cctx, write := sdk.UnwrapSDKContext(kv.ctx).CacheContext()
	ss := kv.WithContext(cctx)
	ss.write = write
	return ss
}

// Context returns the context of the Store.
func (kv *KVStore) Context() context.Context {
	return kv.ctx
}

// WithContext returns a copy of the Store with the given context.
func (kv *KVStore) WithContext(ctx context.Context) *KVStore {
	cpy := *kv
	cpy.ctx = ctx
	return &cpy
}

// Save saves the Store.
func (kv *KVStore) Save() {
	if kv.write != nil {
		kv.write()
	}
}
