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
	"bytes"
	"context"
	"errors"
	"fmt"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/index"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
	sdk "github.com/cosmos/cosmos-sdk/types"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/assert"
)

// KVStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type KVStore[
	BeaconBlockHeaderT interface {
		constraints.Empty[BeaconBlockHeaderT]
		constraints.SSZMarshallable
		sszdb.Treeable
	},
	Eth1DataT interface {
		constraints.Empty[Eth1DataT]
		constraints.SSZMarshallable
		sszdb.Treeable
	},
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		constraints.Empty[ExecutionPayloadHeaderT]
		sszdb.Treeable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
	ForkT interface {
		constraints.Empty[ForkT]
		constraints.SSZMarshallable
		sszdb.Treeable
	},
	ValidatorT Validator[ValidatorT],
	ValidatorsT ~[]ValidatorT,
] struct {
	ctx   context.Context
	sszDB *sszdb.SchemaDB
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
		sszdb.Treeable
	},
	Eth1DataT interface {
		constraints.Empty[Eth1DataT]
		constraints.SSZMarshallable
		sszdb.Treeable
	},
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		constraints.Empty[ExecutionPayloadHeaderT]
		sszdb.Treeable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
	ForkT interface {
		constraints.Empty[ForkT]
		constraints.SSZMarshallable
		sszdb.Treeable
	},
	ValidatorT Validator[ValidatorT],
	ValidatorsT ~[]ValidatorT,
](
	kss store.KVStoreService,
	payloadCodec *encoding.SSZInterfaceCodec[ExecutionPayloadHeaderT],
	sszDB *sszdb.SchemaDB,
) *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kss)
	return &KVStore[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT,
	]{
		ctx:   nil,
		sszDB: sszDB,
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
	ForkT, ValidatorT, ValidatorsT,
]) Copy() *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
] {
	// TODO: Decouple the KVStore type from the Cosmos-SDK.
	cctx, _ := sdk.UnwrapSDKContext(kv.ctx).CacheContext()
	ss := kv.WithContext(cctx)
	return ss
}

// Context returns the context of the Store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) Context() context.Context {
	return kv.ctx
}

// WithContext returns a copy of the Store with the given context.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) WithContext(
	ctx context.Context,
) *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
] {
	cpy := *kv
	cpy.ctx = ctx
	return &cpy
}

func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) debugFieldAssertions() error {
	// genesisValidatorsRoot
	genesisValidatorsRoot, err := kv.GetGenesisValidatorsRoot()
	if err != nil {
		return err
	}
	genesisValidatorsRoot2, err := kv.sszDB.GetPath(
		kv.ctx,
		"genesis_validators_root",
	)
	if err != nil {
		return err
	}
	if !bytes.Equal(genesisValidatorsRoot[:], genesisValidatorsRoot2) {
		return errors.New("genesisValidatorsRoot not equal")
	}

	// slot
	slot, err := kv.GetSlot()
	if err != nil {
		return err
	}
	slot2Bz, err := kv.sszDB.GetPath(kv.ctx, "slot")
	if err != nil {
		return err
	}
	slot2 := fastssz.UnmarshallUint64(slot2Bz)
	if slot2 != slot.Unwrap() {
		return errors.New("slot not equal")
	}

	// fork
	fork, err := kv.GetFork()
	if err != nil {
		return err
	}
	var fork2 ForkT
	fork2 = fork2.Empty()
	err = kv.sszDB.GetObject(kv.ctx, "fork", fork2)
	if err != nil {
		return err
	}
	if !assert.ObjectsAreEqual(fork, fork2) {
		return errors.New("fork not equal")
	}

	// latestBlockHeader
	latestBlockHeader, err := kv.GetLatestBlockHeader()
	if err != nil {
		return err
	}
	var latestBlockHeader2 BeaconBlockHeaderT
	latestBlockHeader2 = latestBlockHeader2.Empty()
	err = kv.sszDB.GetObject(
		kv.ctx,
		"latest_block_header",
		latestBlockHeader2,
	)
	if err != nil {
		return err
	}
	if !assert.ObjectsAreEqual(latestBlockHeader, latestBlockHeader2) {
		return errors.New("latestBlockHeader not equal")
	}

	// blockRoots
	const slotsPerHistoricalRoot = 8
	blockRoots := make([]common.Root, slotsPerHistoricalRoot)
	blockRoots2 := make([]common.Root, slotsPerHistoricalRoot)
	for i := range blockRoots {
		blockRoots[i], err = kv.GetBlockRootAtIndex(uint64(i))
		if err != nil {
			return err
		}
		op := fmt.Sprintf("block_roots/%d", i)
		var bz []byte
		bz, err = kv.sszDB.GetPath(kv.ctx, sszdb.ObjectPath(op))
		if err != nil {
			return err
		}
		blockRoots2[i] = common.Root(bz)
	}
	if !assert.ObjectsAreEqual(blockRoots, blockRoots2) {
		return errors.New("blockRoots not equal")
	}

	// state roots
	stateRoots := make([]common.Root, slotsPerHistoricalRoot)
	stateRoots2 := make([]common.Root, slotsPerHistoricalRoot)
	for i := range stateRoots {
		stateRoots[i], err = kv.StateRootAtIndex(uint64(i))
		if err != nil {
			return err
		}
		op := fmt.Sprintf("state_roots/%d", i)
		var bz []byte
		bz, err = kv.sszDB.GetPath(kv.ctx, sszdb.ObjectPath(op))
		if err != nil {
			return err
		}
		stateRoots2[i] = common.Root(bz)
	}
	if !assert.ObjectsAreEqual(stateRoots, stateRoots2) {
		return errors.New("stateRoots not equal")
	}

	// eth1Data
	eth1Data, err := kv.GetEth1Data()
	if err != nil {
		return err
	}
	var eth1Data2 Eth1DataT
	eth1Data2 = eth1Data2.Empty()
	err = kv.sszDB.GetObject(kv.ctx, "eth1_data", eth1Data2)
	if err != nil {
		return err
	}
	if !assert.ObjectsAreEqual(eth1Data, eth1Data2) {
		return errors.New("eth1Data not equal")
	}

	// eth1DepositIndex
	eth1DepositIndex, err := kv.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	eth1DepositIndex2Bz, err := kv.sszDB.GetPath(kv.ctx, "eth1_deposit_index")
	if err != nil {
		return err
	}
	eth1DepositIndex2 := fastssz.UnmarshallUint64(eth1DepositIndex2Bz)
	if eth1DepositIndex2 != eth1DepositIndex {
		return errors.New("eth1DepositIndex not equal")
	}

	// latestExecutionPayloadHeader
	latestExecutionPayloadHeader, err := kv.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}
	var latestExecutionPayloadHeader2 ExecutionPayloadHeaderT
	latestExecutionPayloadHeader2 = latestExecutionPayloadHeader2.Empty()
	err = kv.sszDB.GetObject(
		kv.ctx,
		"latest_execution_payload_header",
		latestExecutionPayloadHeader2,
	)
	if err != nil {
		return err
	}
	if !assert.ObjectsAreEqual(
		latestExecutionPayloadHeader,
		latestExecutionPayloadHeader2,
	) {
		return errors.New("latestExecutionPayloadHeader not equal")
	}

	// validators
	validators, err := kv.GetValidators()
	if err != nil {
		return err
	}
	numValidators, err := kv.sszDB.GetListLength(kv.ctx, "validators")
	if err != nil {
		return err
	}
	if numValidators != uint64(len(validators)) {
		return errors.New("validators length mismatch")
	}
	validators2 := make([]ValidatorT, numValidators)
	for i := range validators2 {
		validators2[i] = validators2[i].Empty()
		err = kv.sszDB.GetObject(
			kv.ctx,
			sszdb.ObjectPath(fmt.Sprintf("validators/%d", i)),
			validators2[i],
		)
		if err != nil {
			return err
		}
	}
	for i, v := range validators {
		if !assert.ObjectsAreEqual(v, validators2[i]) {
			return errors.New("validators not equal")
		}
	}

	// balances
	balances, err := kv.GetBalances()
	if err != nil {
		return err
	}
	numBalances, err := kv.sszDB.GetListLength(kv.ctx, "balances")
	if err != nil {
		return err
	}
	if numBalances != uint64(len(balances)) {
		return errors.New("balances length mismatch")
	}
	balances2 := make([]uint64, numBalances)
	for i := range balances2 {
		var bz []byte
		bz, err = kv.sszDB.GetPath(
			kv.ctx,
			sszdb.ObjectPath(fmt.Sprintf("balances/%d", i)),
		)
		if err != nil {
			return err
		}
		balances2[i] = fastssz.UnmarshallUint64(bz)
	}
	for i, b := range balances {
		if b != balances2[i] {
			return errors.New("balances not equal")
		}
	}

	// randaoMixes
	const epochsPerHistoricalVector = 8
	randaoMixes := make([]common.Bytes32, epochsPerHistoricalVector)
	randaoMixes2 := make([]common.Bytes32, epochsPerHistoricalVector)
	for i := range randaoMixes {
		randaoMixes[i], err = kv.GetRandaoMixAtIndex(uint64(i))
		if err != nil {
			return err
		}
		op := fmt.Sprintf("randao_mixes/%d", i)
		var bz []byte
		bz, err = kv.sszDB.GetPath(kv.ctx, sszdb.ObjectPath(op))
		if err != nil {
			return err
		}
		randaoMixes2[i] = common.Bytes32(bz)
	}
	if !assert.ObjectsAreEqual(randaoMixes, randaoMixes2) {
		return errors.New("randaoMixes not equal")
	}

	// nextWithdrawalIndex
	nextWithdrawalIndex, err := kv.GetNextWithdrawalIndex()
	if err != nil {
		return err
	}
	nextWithdrawalIndex2Bz, err := kv.sszDB.GetPath(
		kv.ctx,
		"next_withdrawal_index",
	)
	if err != nil {
		return err
	}
	nextWithdrawalIndex2 := fastssz.UnmarshallUint64(nextWithdrawalIndex2Bz)
	if nextWithdrawalIndex2 != nextWithdrawalIndex {
		return fmt.Errorf(
			"nextWithdrawalIndex not equal, expected %d, got %d",
			nextWithdrawalIndex,
			nextWithdrawalIndex2,
		)
	}

	// nextWithdrawalValidatorIndex
	nextWithdrawalValidatorIndex, err := kv.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return err
	}
	nextWithdrawalValidatorIndex2Bz, err := kv.sszDB.GetPath(
		kv.ctx,
		"next_withdrawal_validator_index",
	)
	if err != nil {
		return err
	}
	nextWithdrawalValidatorIndex2 := fastssz.UnmarshallUint64(
		nextWithdrawalValidatorIndex2Bz,
	)
	if nextWithdrawalValidatorIndex2 != nextWithdrawalValidatorIndex.Unwrap() {
		return errors.New("nextWithdrawalValidatorIndex not equal")
	}

	// slashings
	slashings, err := kv.GetSlashings()
	if err != nil {
		return err
	}
	numSlashings, err := kv.sszDB.GetListLength(kv.ctx, "slashings")
	if err != nil {
		return err
	}
	if numSlashings != uint64(len(slashings)) {
		return errors.New("slashings length mismatch")
	}
	slashings2 := make([]uint64, numSlashings)
	for i := range slashings2 {
		var bz []byte
		bz, err = kv.sszDB.GetPath(
			kv.ctx,
			sszdb.ObjectPath(fmt.Sprintf("slashings/%d", i)),
		)
		if err != nil {
			return err
		}
		slashings2[i] = fastssz.UnmarshallUint64(bz)
	}
	for i, s := range slashings {
		if s != slashings2[i] {
			return errors.New("slashings not equal")
		}
	}

	// totalSlashing
	totalSlashing, err := kv.GetTotalSlashing()
	if err != nil {
		return err
	}
	totalSlashing2Bz, err := kv.sszDB.GetPath(kv.ctx, "total_slashing")
	if err != nil {
		return err
	}
	totalSlashing2 := fastssz.UnmarshallUint64(totalSlashing2Bz)
	if totalSlashing2 != totalSlashing.Unwrap() {
		return errors.New("totalSlashing not equal")
	}
	return nil
}

func (kv *KVStore[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
	ValidatorsT,
]) RootHash() ([]byte, error) {
	// debug: preform equivalence checks
	// kv.debugFieldAssertions()

	return kv.sszDB.Hash(kv.ctx)
}
