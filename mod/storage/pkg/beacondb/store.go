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
	storev2 "cosmossdk.io/store/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	storectx "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/context"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

const BeaconStoreKey = "beacon"

// Store is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type Store[
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
] struct {
	ctx       *storectx.Context
	rootStore storev2.RootStore
	storeKey  []byte

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
](
	rootStore storev2.RootStore,
	payloadCodec *encoding.SSZInterfaceCodec[ExecutionPayloadHeaderT],
) *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
] {
	store := &Store[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT,
	]{
		storeKey:  []byte(BeaconStoreKey),
		rootStore: rootStore,
		ctx:       nil,
	}
	schemaBuilder := sdkcollections.NewSchemaBuilderFromAccessor(store.accessor)
	store.genesisValidatorsRoot = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.GenesisValidatorsRootPrefix}),
		keys.GenesisValidatorsRootPrefixHumanReadable,
		sdkcollections.BytesValue,
	)
	store.slot = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.SlotPrefix}),
		keys.SlotPrefixHumanReadable,
		sdkcollections.Uint64Value,
	)
	store.fork = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.ForkPrefix}),
		keys.ForkPrefixHumanReadable,
		encoding.SSZValueCodec[ForkT]{},
	)
	store.blockRoots = sdkcollections.NewMap(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.BlockRootsPrefix}),
		keys.BlockRootsPrefixHumanReadable,
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
	)
	store.stateRoots = sdkcollections.NewMap(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.StateRootsPrefix}),
		keys.StateRootsPrefixHumanReadable,
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
	)
	store.eth1Data = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.Eth1DataPrefix}),
		keys.Eth1DataPrefixHumanReadable,
		encoding.SSZValueCodec[Eth1DataT]{},
	)
	store.eth1DepositIndex = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.Eth1DepositIndexPrefix}),
		keys.Eth1DepositIndexPrefixHumanReadable,
		sdkcollections.Uint64Value,
	)
	store.latestExecutionPayloadVersion = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix(
			[]byte{keys.LatestExecutionPayloadVersionPrefix},
		),
		keys.LatestExecutionPayloadVersionPrefixHumanReadable,
		sdkcollections.Uint32Value,
	)
	store.latestExecutionPayloadCodec = payloadCodec
	store.latestExecutionPayloadHeader = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix(
			[]byte{keys.LatestExecutionPayloadHeaderPrefix},
		),
		keys.LatestExecutionPayloadHeaderPrefixHumanReadable,
		payloadCodec,
	)
	store.validatorIndex = sdkcollections.NewSequence(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.ValidatorIndexPrefix}),
		keys.ValidatorIndexPrefixHumanReadable,
	)
	store.validators = sdkcollections.NewIndexedMap(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.ValidatorByIndexPrefix}),
		keys.ValidatorByIndexPrefixHumanReadable,
		sdkcollections.Uint64Key,
		encoding.SSZValueCodec[ValidatorT]{},
		index.NewValidatorsIndex[ValidatorT](schemaBuilder),
	)
	store.balances = sdkcollections.NewMap(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.BalancesPrefix}),
		keys.BalancesPrefixHumanReadable,
		sdkcollections.Uint64Key,
		sdkcollections.Uint64Value,
	)
	store.randaoMix = sdkcollections.NewMap(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.RandaoMixPrefix}),
		keys.RandaoMixPrefixHumanReadable,
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
	)
	store.slashings = sdkcollections.NewMap(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.SlashingsPrefix}),
		keys.SlashingsPrefixHumanReadable,
		sdkcollections.Uint64Key,
		sdkcollections.Uint64Value,
	)
	store.nextWithdrawalIndex = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.NextWithdrawalIndexPrefix}),
		keys.NextWithdrawalIndexPrefixHumanReadable,
		sdkcollections.Uint64Value,
	)
	store.nextWithdrawalValidatorIndex = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix(
			[]byte{keys.NextWithdrawalValidatorIndexPrefix},
		),
		keys.NextWithdrawalValidatorIndexPrefixHumanReadable,
		sdkcollections.Uint64Value,
	)
	store.totalSlashing = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix([]byte{keys.TotalSlashingPrefix}),
		keys.TotalSlashingPrefixHumanReadable,
		sdkcollections.Uint64Value,
	)
	store.latestBlockHeader = sdkcollections.NewItem(
		schemaBuilder,
		sdkcollections.NewPrefix(
			[]byte{keys.LatestBeaconBlockHeaderPrefix},
		),
		keys.LatestBeaconBlockHeaderPrefixHumanReadable,
		encoding.SSZValueCodec[BeaconBlockHeaderT]{},
	)
	return store
}

func (s *Store[_, _, _, _, _, _]) LatestCommitHash() ([]byte, error) {
	commitID, err := s.rootStore.LastCommitID()
	if err != nil {
		return nil, err
	}

	return commitID.Hash, nil
}

func (s *Store[_, _, _, _, _, _]) Commit() ([]byte, error) {
	changes, err := s.ctx.WriterMap.GetStateChanges()
	if err != nil {
		return nil, err
	}
	return s.rootStore.Commit(&store.Changeset{Changes: changes})
}

func (s *Store[_, _, _, _, _, _]) WorkingHash() ([]byte, error) {
	changes, err := s.ctx.WriterMap.GetStateChanges()
	if err != nil {
		return nil, err
	}
	return s.rootStore.WorkingHash(&store.Changeset{Changes: changes})
}

// Copy returns a copy of the Store.
func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) Copy() *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
] {
	// unnnecessary check?
	if s.ctx == nil {
		return nil
	}
	cctx, err := s.ctx.CacheCopy(s.rootStore)
	if err != nil {
		return nil
	}
	return s.WithContext(cctx)
}

// Context returns the context of the Store.
func (s *Store[
	_, _, _, _, _, _,
]) Context() context.Context {
	return s.ctx
}

// WithContext returns a copy of the Store with the given context.
func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) WithContext(
	ctx context.Context,
) *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
] {
	cpy := *s
	cpy.ctx = storectx.Wrap(ctx, s.storeKey)
	if err := cpy.ctx.AttachStore(cpy.rootStore); err != nil {
		panic(err)
	}
	return &cpy
}

func (s *Store[_, _, _, _, _, _]) accessor(rawCtx context.Context) store.KVStore {
	ctx := storectx.Wrap(rawCtx, s.storeKey)
	if err := ctx.AttachStore(s.rootStore); err != nil {
		panic(err)
	}
	writer, err := ctx.WriterMap.GetWriter(s.storeKey)
	if err != nil {
		panic(err)
	}
	return writer
}
