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
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/runtime/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
	indexv2 "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index/v2"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

const ModuleName = "beacon"

// Store is a wrapper around storev2.RootStore
type Store[
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
	runtime.Store
	changeSet *store.Changeset
	// Versioning
	// genesisValidatorsRoot is the root of the genesis validators.
	genesisValidatorsRoot collections.Item[[]byte]
	// slot is the current slot.
	slot collections.Item[uint64]
	// fork is the current fork
	fork collections.Item[ForkT]
	// History
	// latestBlockHeader stores the latest beacon block header.
	latestBlockHeader collections.Item[BeaconBlockHeaderT]
	// blockRoots stores the block roots for the current epoch.
	blockRoots collections.Map[uint64, []byte]
	// stateRoots stores the state roots for the current epoch.
	stateRoots collections.Map[uint64, []byte]
	// Eth1
	// eth1Data stores the latest eth1 data.
	eth1Data collections.Item[Eth1DataT]
	// eth1DepositIndex is the index of the latest eth1 deposit.
	eth1DepositIndex collections.Item[uint64]
	// latestExecutionPayload stores the latest execution payload version.
	latestExecutionPayloadVersion collections.Item[uint32]
	// latestExecutionPayloadCodec is the codec for the latest execution
	// payload, it allows us to update the codec with the latest version.
	latestExecutionPayloadCodec *encoding.
					SSZInterfaceCodec[ExecutionPayloadHeaderT]
	// latestExecutionPayloadHeader stores the latest execution payload header.
	latestExecutionPayloadHeader collections.Item[ExecutionPayloadHeaderT]
	// Registry
	// validatorIndex provides the next available index for a new validator.
	validatorIndex collections.Sequence
	// validators stores the list of validators.
	validators *collections.IndexedMap[
		uint64, ValidatorT, indexv2.ValidatorsIndex[ValidatorT],
	]
	// balances stores the list of balances.
	balances collections.Map[uint64, uint64]
	// nextWithdrawalIndex stores the next global withdrawal index.
	nextWithdrawalIndex collections.Item[uint64]
	// nextWithdrawalValidatorIndex stores the next withdrawal validator index
	// for each validator.
	nextWithdrawalValidatorIndex collections.Item[uint64]
	// Randomness
	// randaoMix stores the randao mix for the current epoch.
	randaoMix collections.Map[uint64, []byte]
	// Slashings
	// slashings stores the slashings for the current epoch.
	slashings collections.Map[uint64, uint64]
	// totalSlashing stores the total slashing in the vector range.
	totalSlashing collections.Item[uint64]
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
	payloadCodec *encoding.SSZInterfaceCodec[ExecutionPayloadHeaderT],
) *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
] {
	storeKey := []byte(ModuleName)
	store := &Store[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	]{
		changeSet: store.NewChangeset(),
	}

	store.genesisValidatorsRoot = collections.NewItem(
		storeKey,
		[]byte{keys.GenesisValidatorsRootPrefix},
		sdkcollections.BytesValue,
		store.accessor,
	)
	store.slot = collections.NewItem(
		storeKey,
		[]byte{keys.SlotPrefix},
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.fork = collections.NewItem(
		storeKey,
		[]byte{keys.ForkPrefix},
		encoding.SSZValueCodec[ForkT]{},
		store.accessor,
	)
	store.blockRoots = collections.NewMap(
		storeKey,
		[]byte{keys.BlockRootsPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		store.accessor,
	)
	store.stateRoots = collections.NewMap(
		storeKey,
		[]byte{keys.StateRootsPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		store.accessor,
	)
	store.eth1Data = collections.NewItem(
		storeKey,
		[]byte{keys.Eth1DataPrefix},
		encoding.SSZValueCodec[Eth1DataT]{},
		store.accessor,
	)
	store.eth1DepositIndex = collections.NewItem(
		storeKey,
		[]byte{keys.Eth1DepositIndexPrefix},
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.latestExecutionPayloadVersion = collections.NewItem(
		storeKey,
		[]byte{keys.LatestExecutionPayloadVersionPrefix},
		sdkcollections.Uint32Value,
		store.accessor,
	)
	store.latestExecutionPayloadCodec = payloadCodec
	store.latestExecutionPayloadHeader = collections.NewItem(
		storeKey,
		[]byte{keys.LatestExecutionPayloadHeaderPrefix},
		payloadCodec,
		store.accessor,
	)
	store.validatorIndex = collections.NewSequence(
		storeKey,
		[]byte{keys.ValidatorIndexPrefix},
		store.accessor,
	)
	store.validators = collections.NewIndexedMap(
		storeKey,
		[]byte{keys.ValidatorByIndexPrefix},
		sdkcollections.Uint64Key,
		encoding.SSZValueCodec[ValidatorT]{},
		indexv2.NewValidatorsIndex[ValidatorT](store.accessor),
		store.accessor,
	)
	store.balances = collections.NewMap(
		storeKey,
		[]byte{keys.BalancesPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.randaoMix = collections.NewMap(
		storeKey,
		[]byte{keys.RandaoMixPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		store.accessor,
	)
	store.slashings = collections.NewMap(
		storeKey,
		[]byte{keys.SlashingsPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.nextWithdrawalIndex = collections.NewItem(
		storeKey,
		[]byte{keys.NextWithdrawalIndexPrefix},
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.nextWithdrawalValidatorIndex = collections.NewItem(
		storeKey,
		[]byte{keys.NextWithdrawalValidatorIndexPrefix},
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.totalSlashing = collections.NewItem(
		storeKey,
		[]byte{keys.TotalSlashingPrefix},
		sdkcollections.Uint64Value,
		store.accessor,
	)
	store.latestBlockHeader = collections.NewItem(
		storeKey,
		[]byte{keys.LatestBeaconBlockHeaderPrefix},
		encoding.SSZValueCodec[BeaconBlockHeaderT]{},
		store.accessor,
	)
	return store
}

// if commit errors should we still reset? maybe just do an
// explicit call instead of defer to prevent that case
// TODO: return store hash
func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) Save() {
	// reset the changeset following the commit
	defer func() {
		s.changeSet = store.NewChangeset()
	}()
	if s.changeSet.Size() == 0 {
		return
	}
	s.Store.Commit(s.changeSet)
}

// TODO: deprecate
func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) Copy() *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
] {
	return s
}

// Note: this function does not enforce the invariant that
// the changeset must not be nil. more performant ish but less safe
func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) AddChange(storeKey []byte, key []byte, value []byte) {
	s.changeSet.Add(storeKey, key, value, false)
}

func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) SetBackendStore(store runtime.Store) {
	s.Store = store
}

func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) accessor() collections.Store {
	return s
}
