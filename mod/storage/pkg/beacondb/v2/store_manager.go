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
	"cosmossdk.io/runtime/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
	indexv2 "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index/v2"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

const ModuleName = "beacon"

// StoreManager is a wrapper around storev2.RootStore
type StoreManager[
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
	store *StateStore
	// Versioning
	// genesisValidatorsRoot is the root of the genesis validators.
	genesisValidatorsRoot collections.ItemKeeper[[]byte]
	// slot is the current slot.
	slot collections.ItemKeeper[uint64]
	// fork is the current fork
	fork collections.ItemKeeper[ForkT]
	// History
	// latestBlockHeader stores the latest beacon block header.
	latestBlockHeader collections.ItemKeeper[BeaconBlockHeaderT]
	// blockRoots stores the block roots for the current epoch.
	blockRoots collections.MapKeeper[uint64, []byte]
	// stateRoots stores the state roots for the current epoch.
	stateRoots collections.MapKeeper[uint64, []byte]
	// Eth1
	// eth1Data stores the latest eth1 data.
	eth1Data collections.ItemKeeper[Eth1DataT]
	// eth1DepositIndex is the index of the latest eth1 deposit.
	eth1DepositIndex collections.ItemKeeper[uint64]
	// latestExecutionPayload stores the latest execution payload version.
	latestExecutionPayloadVersion collections.ItemKeeper[uint32]
	// latestExecutionPayloadCodec is the codec for the latest execution
	// payload, it allows us to update the codec with the latest version.
	latestExecutionPayloadCodec *encoding.
					SSZInterfaceCodec[ExecutionPayloadHeaderT]
	// latestExecutionPayloadHeader stores the latest execution payload header.
	latestExecutionPayloadHeader collections.ItemKeeper[ExecutionPayloadHeaderT]
	// Registry
	// validatorIndex provides the next available index for a new validator.
	validatorIndex collections.Sequence
	// validators stores the list of validators.
	validators *collections.IndexedMapKeeper[
		uint64, ValidatorT, indexv2.ValidatorsIndex[ValidatorT],
	]
	// balances stores the list of balances.
	balances collections.MapKeeper[uint64, uint64]
	// nextWithdrawalIndex stores the next global withdrawal index.
	nextWithdrawalIndex collections.ItemKeeper[uint64]
	// nextWithdrawalValidatorIndex stores the next withdrawal validator index
	// for each validator.
	nextWithdrawalValidatorIndex collections.ItemKeeper[uint64]
	// Randomness
	// randaoMix stores the randao mix for the current epoch.
	randaoMix collections.MapKeeper[uint64, []byte]
	// Slashings
	// slashings stores the slashings for the current epoch.
	slashings collections.MapKeeper[uint64, uint64]
	// totalSlashing stores the total slashing in the vector range.
	totalSlashing collections.ItemKeeper[uint64]
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
) *StoreManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
] {
	storeKey := []byte(ModuleName)
	store := &StoreManager[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	]{
		store: NewStore(),
	}

	store.genesisValidatorsRoot = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.GenesisValidatorsRootPrefix},
		sdkcollections.BytesValue,
		store.stateStore,
	)
	store.slot = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.SlotPrefix},
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.fork = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.ForkPrefix},
		encoding.SSZValueCodec[ForkT]{},
		store.stateStore,
	)
	store.blockRoots = collections.NewMapKeeper(
		storeKey,
		[]byte{keys.BlockRootsPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		store.stateStore,
	)
	store.stateRoots = collections.NewMapKeeper(
		storeKey,
		[]byte{keys.StateRootsPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		store.stateStore,
	)
	store.eth1Data = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.Eth1DataPrefix},
		encoding.SSZValueCodec[Eth1DataT]{},
		store.stateStore,
	)
	store.eth1DepositIndex = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.Eth1DepositIndexPrefix},
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.latestExecutionPayloadVersion = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.LatestExecutionPayloadVersionPrefix},
		sdkcollections.Uint32Value,
		store.stateStore,
	)
	store.latestExecutionPayloadCodec = payloadCodec
	store.latestExecutionPayloadHeader = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.LatestExecutionPayloadHeaderPrefix},
		payloadCodec,
		store.stateStore,
	)
	store.validatorIndex = collections.NewSequence(
		storeKey,
		[]byte{keys.ValidatorIndexPrefix},
		store.stateStore,
	)
	store.validators = collections.NewIndexedMap(
		storeKey,
		[]byte{keys.ValidatorByIndexPrefix},
		sdkcollections.Uint64Key,
		encoding.SSZValueCodec[ValidatorT]{},
		indexv2.NewValidatorsIndex[ValidatorT](store.stateStore),
		store.stateStore,
	)
	store.balances = collections.NewMapKeeper(
		storeKey,
		[]byte{keys.BalancesPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.randaoMix = collections.NewMapKeeper(
		storeKey,
		[]byte{keys.RandaoMixPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.BytesValue,
		store.stateStore,
	)
	store.slashings = collections.NewMapKeeper(
		storeKey,
		[]byte{keys.SlashingsPrefix},
		sdkcollections.Uint64Key,
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.nextWithdrawalIndex = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.NextWithdrawalIndexPrefix},
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.nextWithdrawalValidatorIndex = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.NextWithdrawalValidatorIndexPrefix},
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.totalSlashing = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.TotalSlashingPrefix},
		sdkcollections.Uint64Value,
		store.stateStore,
	)
	store.latestBlockHeader = collections.NewItemKeeper(
		storeKey,
		[]byte{keys.LatestBeaconBlockHeaderPrefix},
		encoding.SSZValueCodec[BeaconBlockHeaderT]{},
		store.stateStore,
	)
	return store
}

// if commit errors should we still reset? maybe just do an
// explicit call instead of defer to prevent that case
// TODO: return store hash
func (s *StoreManager[_, _, _, _, _]) Save() {
	s.store.Save()
}

// TODO: deprecate
func (s *StoreManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) Copy() *StoreManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
] {
	return s
}

func (s *StoreManager[_, _, _, _, _]) SetStateStore(store runtime.Store) {
	s.store.SetStore(store)
}

func (s *StoreManager[_, _, _, _, _]) stateStore() collections.Store {
	return s.store
}
