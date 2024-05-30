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
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KVStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type KVStore[
	ForkT SSZMarshallable,
	BeaconBlockHeaderT SSZMarshallable,
	ExecutionPayloadHeaderT SSZMarshallable,
	Eth1DataT SSZMarshallable,
	ValidatorT Validator,
] struct {
	ctx   context.Context
	write func()
	// Versioning
	// genesisValidatorsRoot is the root of the genesis validators.
	genesisValidatorsRoot sdkcollections.Item[[]byte]
	// slot is the current slot.
	slot sdkcollections.Item[uint64]
	// fork is the current fork
	fork sdkcollections.Item[ForkT]
	// History
	// latestBlockHeader stores the latest beacon block header.
	latestBlockHeader sdkcollections.Item[BeaconBlockHeaderT]
	// blockRoots stores the block roots for the current epoch.
	blockRoots sdkcollections.Map[uint64, []byte]
	// stateRoots stores the state roots for the current epoch.
	stateRoots sdkcollections.Map[uint64, []byte]
	// Eth1
	// eth1Data stores the latest eth1 data.
	eth1Data sdkcollections.Item[Eth1DataT]
	// eth1DepositIndex is the index of the latest eth1 deposit.
	eth1DepositIndex sdkcollections.Item[uint64]
	// latestExecutionPayloadHeader stores the latest execution payload header.
	latestExecutionPayloadHeader sdkcollections.Item[ExecutionPayloadHeaderT]
	// Registry
	// validatorIndex provides the next available index for a new validator.
	validatorIndex sdkcollections.Sequence
	// validators stores the list of validators.
	validators *sdkcollections.IndexedMap[
		uint64, ValidatorT, index.ValidatorsIndex[ValidatorT],
	]
	// balances stores the list of balances.
	balances sdkcollections.Map[uint64, uint64]
	// nextWithdrawalIndex stores the next global withdrawal index.
	nextWithdrawalIndex sdkcollections.Item[uint64]
	// nextWithdrawalValidatorIndex stores the next withdrawal validator index
	// for each validator.
	nextWithdrawalValidatorIndex sdkcollections.Item[uint64]
	// Randomness
	// randaoMix stores the randao mix for the current epoch.
	randaoMix sdkcollections.Map[uint64, []byte]
	// Slashings
	// slashings stores the slashings for the current epoch.
	slashings sdkcollections.Map[uint64, uint64]
	// totalSlashing stores the total slashing in the vector range.
	totalSlashing sdkcollections.Item[uint64]
}

// Store creates a new instance of Store.
//
//nolint:funlen // its not overly complex.
func New[
	ForkT SSZMarshallable,
	BeaconBlockHeaderT SSZMarshallable,
	ExecutionPayloadHeaderT SSZMarshallable,
	Eth1DataT SSZMarshallable,
	ValidatorT Validator,
](
	kss store.KVStoreService,
	executionPayloadHeaderFactory func() ExecutionPayloadHeaderT,
) *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadHeaderT, Eth1DataT, ValidatorT,
] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kss)
	return &KVStore[
		ForkT, BeaconBlockHeaderT,
		ExecutionPayloadHeaderT, Eth1DataT, ValidatorT,
	]{
		ctx: nil,
		genesisValidatorsRoot: sdkcollections.NewItem[[]byte](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.GenesisValidatorsRootPrefix}),
			keys.GenesisValidatorsRootPrefixHumanReadable,
			sdkcollections.BytesValue,
		),
		slot: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.SlotPrefix}),
			keys.SlotPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		fork: sdkcollections.NewItem[ForkT](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.ForkPrefix}),
			keys.ForkPrefixHumanReadable,
			encoding.SSZValueCodec[ForkT]{},
		),
		blockRoots: sdkcollections.NewMap[uint64, []byte](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.BlockRootsPrefix}),
			keys.BlockRootsPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		stateRoots: sdkcollections.NewMap[uint64, []byte](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.StateRootsPrefix}),
			keys.StateRootsPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		eth1Data: sdkcollections.NewItem[Eth1DataT](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.Eth1DataPrefix}),
			keys.Eth1DataPrefixHumanReadable,
			encoding.SSZValueCodec[Eth1DataT]{},
		),
		eth1DepositIndex: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.Eth1DepositIndexPrefix}),
			keys.Eth1DepositIndexPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		latestExecutionPayloadHeader: sdkcollections.NewItem[ExecutionPayloadHeaderT](
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.LatestExecutionPayloadHeaderPrefix},
			),
			keys.LatestExecutionPayloadHeaderPrefixHumanReadable,
			encoding.SSZInterfaceCodec[ExecutionPayloadHeaderT]{
				Factory: executionPayloadHeaderFactory,
			},
		),
		validatorIndex: sdkcollections.NewSequence(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.ValidatorIndexPrefix}),
			keys.ValidatorIndexPrefixHumanReadable,
		),
		validators: sdkcollections.NewIndexedMap[
			uint64, ValidatorT,
		](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.ValidatorByIndexPrefix}),
			keys.ValidatorByIndexPrefixHumanReadable,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[ValidatorT]{},
			index.NewValidatorsIndex[ValidatorT](schemaBuilder),
		),
		balances: sdkcollections.NewMap[uint64, uint64](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.BalancesPrefix}),
			keys.BalancesPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		randaoMix: sdkcollections.NewMap[uint64, []byte](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.RandaoMixPrefix}),
			keys.RandaoMixPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		slashings: sdkcollections.NewMap[uint64, uint64](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.SlashingsPrefix}),
			keys.SlashingsPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalIndex: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.NextWithdrawalIndexPrefix}),
			keys.NextWithdrawalIndexPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalValidatorIndex: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.NextWithdrawalValidatorIndexPrefix},
			),
			keys.NextWithdrawalValidatorIndexPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		totalSlashing: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.TotalSlashingPrefix}),
			keys.TotalSlashingPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		latestBlockHeader: sdkcollections.NewItem[BeaconBlockHeaderT](
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.LatestBeaconBlockHeaderPrefix},
			),
			keys.LatestBeaconBlockHeaderPrefixHumanReadable,
			encoding.SSZValueCodec[BeaconBlockHeaderT]{},
		),
	}
}

// Copy returns a copy of the Store.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) Copy() *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
] {
	// TODO: Decouple the KVStore type from the Cosmos-SDK.
	cctx, write := sdk.UnwrapSDKContext(kv.ctx).CacheContext()
	ss := kv.WithContext(cctx)
	ss.write = write
	return ss
}

// Context returns the context of the Store.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) Context() context.Context {
	return kv.ctx
}

// WithContext returns a copy of the Store with the given context.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) WithContext(
	ctx context.Context,
) *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
] {
	cpy := *kv
	cpy.ctx = ctx
	return &cpy
}

// Save saves the Store.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) Save() {
	if kv.write != nil {
		kv.write()
	}
}
