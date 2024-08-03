// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package beacondb

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KVStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type KVStore[
	BeaconBlockHeaderT constraints.SSZMarshallable,
	Eth1DataT constraints.SSZMarshallable,
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
	ForkT constraints.SSZMarshallable,
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
	// latestExecutionPayloadVersion stores the latest execution payload
	// version.
	latestExecutionPayloadVersion sdkcollections.Item[uint32]
	// latestExecutionPayloadCodec is the codec for the latest execution
	// payload, it allows us to update the codec with the latest version.
	latestExecutionPayloadCodec *encoding.
					SSZInterfaceCodec[ExecutionPayloadHeaderT]
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

// New creates a new instance of Store.
//
//nolint:funlen // its not overly complex.
func New[
	BeaconBlockHeaderT constraints.SSZMarshallable,
	Eth1DataT constraints.SSZMarshallable,
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
	ForkT constraints.SSZMarshallable,
	ValidatorT Validator,
](
	kss store.KVStoreService,
	payloadCodec *encoding.SSZInterfaceCodec[ExecutionPayloadHeaderT],
) *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kss)
	return &KVStore[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	]{
		ctx: nil,
		genesisValidatorsRoot: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.GenesisValidatorsRootPrefix}),
			keys.GenesisValidatorsRootPrefixHumanReadable,
			sdkcollections.BytesValue,
		),
		slot: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.SlotPrefix}),
			keys.SlotPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		fork: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.ForkPrefix}),
			keys.ForkPrefixHumanReadable,
			encoding.SSZValueCodec[ForkT]{},
		),
		blockRoots: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.BlockRootsPrefix}),
			keys.BlockRootsPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		stateRoots: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.StateRootsPrefix}),
			keys.StateRootsPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		eth1Data: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.Eth1DataPrefix}),
			keys.Eth1DataPrefixHumanReadable,
			encoding.SSZValueCodec[Eth1DataT]{},
		),
		eth1DepositIndex: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.Eth1DepositIndexPrefix}),
			keys.Eth1DepositIndexPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		latestExecutionPayloadVersion: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.LatestExecutionPayloadVersionPrefix},
			),
			keys.LatestExecutionPayloadVersionPrefixHumanReadable,
			sdkcollections.Uint32Value,
		),
		latestExecutionPayloadCodec: payloadCodec,
		latestExecutionPayloadHeader: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.LatestExecutionPayloadHeaderPrefix},
			),
			keys.LatestExecutionPayloadHeaderPrefixHumanReadable,
			payloadCodec,
		),
		validatorIndex: sdkcollections.NewSequence(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.ValidatorIndexPrefix}),
			keys.ValidatorIndexPrefixHumanReadable,
		),
		validators: sdkcollections.NewIndexedMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.ValidatorByIndexPrefix}),
			keys.ValidatorByIndexPrefixHumanReadable,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[ValidatorT]{},
			index.NewValidatorsIndex[ValidatorT](schemaBuilder),
		),
		balances: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.BalancesPrefix}),
			keys.BalancesPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		randaoMix: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.RandaoMixPrefix}),
			keys.RandaoMixPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		slashings: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.SlashingsPrefix}),
			keys.SlashingsPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalIndex: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.NextWithdrawalIndexPrefix}),
			keys.NextWithdrawalIndexPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalValidatorIndex: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.NextWithdrawalValidatorIndexPrefix},
			),
			keys.NextWithdrawalValidatorIndexPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		totalSlashing: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.TotalSlashingPrefix}),
			keys.TotalSlashingPrefixHumanReadable,
			sdkcollections.Uint64Value,
		),
		latestBlockHeader: sdkcollections.NewItem(
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
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) Copy() *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
] {
	// TODO: Decouple the KVStore type from the Cosmos-SDK.
	cctx, write := sdk.UnwrapSDKContext(kv.ctx).CacheContext()
	ss := kv.WithContext(cctx)
	ss.write = write
	return ss
}

// Context returns the context of the Store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) Context() context.Context {
	return kv.ctx
}

// WithContext returns a copy of the Store with the given context.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) WithContext(
	ctx context.Context,
) *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
] {
	cpy := *kv
	cpy.ctx = ctx
	return &cpy
}

// Save saves the Store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) Save() {
	if kv.write != nil {
		kv.write()
	}
}
