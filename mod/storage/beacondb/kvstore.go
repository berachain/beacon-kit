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
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections/encoding"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/index"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SSZMarshallable is an interface that combines the ssz.Marshaler and
// ssz.Unmarshaler interfaces.
type SSZMarshallable interface {
	// MarshalSSZTo marshals the object into the provided byte slice and returns
	// it along with any error.
	MarshalSSZTo([]byte) ([]byte, error)
	// MarshalSSZ marshals the object into a new byte slice and returns it along
	// with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
	// SizeSSZ returns the size in bytes that the object would take when
	// marshaled.
	SizeSSZ() int
}

type Validator interface {
	SSZMarshallable
	GetPubkey() primitives.BLSPubkey
	GetEffectiveBalance() primitives.Gwei
	IsActive(primitives.Epoch) bool
}

// KVStore is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type KVStore[
	DepositT SSZMarshallable,
	ForkT SSZMarshallable,
	BeaconBlockHeaderT SSZMarshallable,
	ExecutionPayloadT SSZMarshallable,
	Eth1DataT SSZMarshallable,
	ValidatorT Validator,
] struct {
	ctx   context.Context
	write func()

	// Versioning
	// genesisValidatorsRoot is the root of the genesis validators.
	genesisValidatorsRoot sdkcollections.Item[[32]byte]
	// slot is the current slot.
	slot sdkcollections.Item[uint64]
	// fork is the current fork
	fork sdkcollections.Item[ForkT]

	// History
	// latestBlockHeader stores the latest beacon block header.
	latestBlockHeader sdkcollections.Item[BeaconBlockHeaderT]
	// blockRoots stores the block roots for the current epoch.
	blockRoots sdkcollections.Map[uint64, [32]byte]
	// stateRoots stores the state roots for the current epoch.
	stateRoots sdkcollections.Map[uint64, [32]byte]

	// Eth1
	// latestExecutionPayload stores the latest execution payload.

	latestExecutionPayload sdkcollections.Item[ExecutionPayloadT]

	// eth1Data stores the latest eth1 data.
	eth1Data sdkcollections.Item[Eth1DataT]
	// eth1DepositIndex is the index of the latest eth1 deposit.
	eth1DepositIndex sdkcollections.Item[uint64]

	// Registry
	// validatorIndex provides the next available index for a new validator.
	validatorIndex sdkcollections.Sequence
	// validators stores the list of validators.
	validators *sdkcollections.IndexedMap[
		uint64, ValidatorT, index.ValidatorsIndex[ValidatorT],
	]
	// balances stores the list of balances.
	balances sdkcollections.Map[uint64, uint64]

	// depositQueue is a list of deposits that are queued to be processed.
	depositQueue *collections.Queue[DepositT]

	// nextWithdrawalIndex stores the next global withdrawal index.
	nextWithdrawalIndex sdkcollections.Item[uint64]

	// nextWithdrawalValidatorIndex stores the next withdrawal validator index
	// for each validator.
	nextWithdrawalValidatorIndex sdkcollections.Item[uint64]

	// Randomness
	// randaoMix stores the randao mix for the current epoch.
	randaoMix sdkcollections.Map[uint64, [32]byte]

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
	DepositT SSZMarshallable,
	ForkT SSZMarshallable,
	BeaconBlockHeaderT SSZMarshallable,
	ExecutionPayloadT SSZMarshallable,
	Eth1DataT SSZMarshallable,
	ValidatorT Validator,
](
	kss store.KVStoreService,
	executionPayloadFactory func() ExecutionPayloadT,
) *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kss)
	return &KVStore[
		DepositT, ForkT, BeaconBlockHeaderT,
		ExecutionPayloadT, Eth1DataT, ValidatorT,
	]{
		ctx: nil,
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
		fork: sdkcollections.NewItem[ForkT](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.ForkPrefix),
			keys.ForkPrefix,
			encoding.SSZValueCodec[ForkT]{},
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

		latestExecutionPayload: sdkcollections.NewItem[ExecutionPayloadT](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.LatestExecutionPayloadPrefix),
			keys.LatestExecutionPayloadPrefix,
			encoding.SSZInterfaceCodec[ExecutionPayloadT]{
				Factory: executionPayloadFactory,
			},
		),
		eth1Data: sdkcollections.NewItem[Eth1DataT](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.Eth1DataPrefix),
			keys.Eth1DataPrefix,
			encoding.SSZValueCodec[Eth1DataT]{},
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
			uint64, ValidatorT,
		](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.ValidatorByIndexPrefix),
			keys.ValidatorByIndexPrefix,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[ValidatorT]{},
			index.NewValidatorsIndex[ValidatorT](schemaBuilder),
		),
		balances: sdkcollections.NewMap[uint64, uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.BalancesPrefix),
			keys.BalancesPrefix,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Value,
		),
		depositQueue: collections.NewQueue[DepositT](
			schemaBuilder,
			keys.DepositQueuePrefix,
			encoding.SSZValueCodec[DepositT]{},
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
		nextWithdrawalIndex: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.NextWithdrawalIndexPrefix),
			keys.NextWithdrawalIndexPrefix,
			sdkcollections.Uint64Value,
		),
		nextWithdrawalValidatorIndex: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.NextWithdrawalValidatorIndexPrefix),
			keys.NextWithdrawalValidatorIndexPrefix,
			sdkcollections.Uint64Value,
		),

		totalSlashing: sdkcollections.NewItem[uint64](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.TotalSlashingPrefix),
			keys.TotalSlashingPrefix,
			sdkcollections.Uint64Value,
		),

		latestBlockHeader: sdkcollections.NewItem[BeaconBlockHeaderT](
			schemaBuilder,
			sdkcollections.NewPrefix(keys.LatestBeaconBlockHeaderPrefix),
			keys.LatestBeaconBlockHeaderPrefix,
			encoding.SSZValueCodec[BeaconBlockHeaderT]{},
		),
	}
}

// Copy returns a copy of the Store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) Copy() *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
] {
	// TODO: Decouple the KVStore type from the Cosmos-SDK.
	cctx, write := sdk.UnwrapSDKContext(kv.ctx).CacheContext()
	ss := kv.WithContext(cctx)
	ss.write = write
	return ss
}

// Context returns the context of the Store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) Context() context.Context {
	return kv.ctx
}

// WithContext returns a copy of the Store with the given context.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) WithContext(
	ctx context.Context,
) *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
] {
	cpy := *kv
	cpy.ctx = ctx
	return &cpy
}

// Save saves the Store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) Save() {
	if kv.write != nil {
		kv.write()
	}
}
