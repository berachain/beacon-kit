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
	"fmt"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/examples/berad/pkg/storage/keys"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KVStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type KVStore[
	BeaconBlockHeaderT interface {
		constraints.Empty[BeaconBlockHeaderT]
		constraints.SSZMarshallable
	},
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
	ForkT interface {
		constraints.Empty[ForkT]
		constraints.SSZMarshallable
	},
	ValidatorT Validator[ValidatorT],
	ValidatorsT ~[]ValidatorT,
	WithdrawalT any,
	WithdrawalsT ~[]WithdrawalT,
] struct {
	ctx context.Context
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
	// Randomness
	// randaoMix stores the randao mix for the current epoch.
	randaoMix sdkcollections.Map[uint64, []byte]
	// Staking
	// withdrawals stores a list of initiated withdrawals.
	withdrawals sdkcollections.Item[WithdrawalsT]
}

// New creates a new instance of Store.
//
//nolint:funlen // its not overly complex.
func New[
	BeaconBlockHeaderT interface {
		constraints.Empty[BeaconBlockHeaderT]
		constraints.SSZMarshallable
	},
	Eth1DataT interface {
		constraints.Empty[Eth1DataT]
		constraints.SSZMarshallable
	},
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
	ForkT interface {
		constraints.Empty[ForkT]
		constraints.SSZMarshallable
	},
	ValidatorT Validator[ValidatorT],
	ValidatorsT ~[]ValidatorT,
	WithdrawalT interface {
		constraints.Empty[WithdrawalT]
		constraints.SSZMarshallable
	},
	WithdrawalsT interface {
		~[]WithdrawalT
		constraints.Empty[WithdrawalsT]
		constraints.SSZMarshallable
	},
](
	kss store.KVStoreService,
	payloadCodec *encoding.SSZInterfaceCodec[ExecutionPayloadHeaderT],
) *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kss)
	res := &KVStore[
		BeaconBlockHeaderT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
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
		randaoMix: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.RandaoMixPrefix}),
			keys.RandaoMixPrefixHumanReadable,
			sdkcollections.Uint64Key,
			sdkcollections.BytesValue,
		),
		latestBlockHeader: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix(
				[]byte{keys.LatestBeaconBlockHeaderPrefix},
			),
			keys.LatestBeaconBlockHeaderPrefixHumanReadable,
			encoding.SSZValueCodec[BeaconBlockHeaderT]{},
		),
		withdrawals: sdkcollections.NewItem(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{keys.WithdrawalsPrefix}),
			keys.WithdrawalsPrefixHumanReadable,
			encoding.SSZValueCodec[WithdrawalsT]{},
		),
	}
	if _, err := schemaBuilder.Build(); err != nil {
		panic(fmt.Errorf("failed building KVStore schema: %w", err))
	}
	return res
}

// Copy returns a copy of the Store.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) Copy() *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
] {
	// TODO: Decouple the KVStore type from the Cosmos-SDK.
	cctx, _ := sdk.UnwrapSDKContext(kv.ctx).CacheContext()
	ss := kv.WithContext(cctx)
	return ss
}

// Context returns the context of the Store.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) Context() context.Context {
	return kv.ctx
}

// WithContext returns a copy of the Store with the given context.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) WithContext(
	ctx context.Context,
) *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
] {
	cpy := *kv
	cpy.ctx = ctx
	return &cpy
}
