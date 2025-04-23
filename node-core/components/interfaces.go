// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package components

import (
	"context"

	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	nodecoretypes "github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/block"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
)

type (
	// AttributesFactory is the interface for the attributes factory.
	AttributesFactory interface {
		BuildPayloadAttributes(
			st *statedb.StateDB,
			slot math.Slot,
			timestamp math.U64,
			prevHeadRoot [32]byte,
		) (*engineprimitives.PayloadAttributes, error)
	}

	// BlobProcessor is the interface for the blobs processor.
	BlobProcessor interface {
		// ProcessSidecars processes the blobs and ensures they match the local
		// state.
		ProcessSidecars(
			avs *dastore.Store,
			sidecars datypes.BlobSidecars,
		) error
		// VerifySidecars verifies the blobs and ensures they match the local
		// state.
		VerifySidecars(
			ctx context.Context,
			sidecars datypes.BlobSidecars,
			blkHeader *ctypes.BeaconBlockHeader,
			kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
		) error
	}

	// LocalBuilder is the interface for the builder service.
	LocalBuilder interface {
		// Enabled returns true if the local builder is enabled.
		Enabled() bool
		// RequestPayloadAsync requests a new payload for the given slot.
		RequestPayloadAsync(
			ctx context.Context,
			st *statedb.StateDB,
			slot math.Slot,
			timestamp math.U64,
			parentBlockRoot common.Root,
			headEth1BlockHash common.ExecutionHash,
			finalEth1BlockHash common.ExecutionHash,
		) (*engineprimitives.PayloadID, common.Version, error)
		// RetrievePayload retrieves the payload for the given slot.
		RetrievePayload(
			ctx context.Context,
			slot math.Slot,
			parentBlockRoot common.Root,
		) (ctypes.BuiltExecutionPayloadEnv, error)
		// RequestPayloadSync requests a payload for the given slot and
		// blocks until the payload is delivered.
		RequestPayloadSync(
			ctx context.Context,
			st *statedb.StateDB,
			slot math.Slot,
			timestamp math.U64,
			parentBlockRoot common.Root,
			headEth1BlockHash common.ExecutionHash,
			finalEth1BlockHash common.ExecutionHash,
		) (ctypes.BuiltExecutionPayloadEnv, error)
	}

	// 	// PayloadAttributes is the interface for the payload attributes.
	// PayloadAttributes[T any, WithdrawalT any] interface {
	// 	engineprimitives.PayloadAttributer
	// 	// New creates a new payload attributes instance.
	// 	New(
	// 		uint32,
	// 		uint64,
	// 		common.Bytes32,
	// 		common.ExecutionAddress,
	// 		[]WithdrawalT,
	// 		common.Root,
	// 	) (T, error)
	// }.

	// StateProcessor defines the interface for processing the state.
	StateProcessor interface {
		// InitializeBeaconStateFromEth1 initializes the premined beacon
		// state
		// from the eth1 deposits.
		InitializeBeaconStateFromEth1(
			*statedb.StateDB,
			ctypes.Deposits,
			*ctypes.ExecutionPayloadHeader,
			common.Version,
		) (transition.ValidatorUpdates, error)
		// ProcessFork prepares the state for the fork version at the given timestamp.
		ProcessFork(
			st *statedb.StateDB, timestamp math.U64, logUpgrade bool,
		) error
		// ProcessSlot processes the slot.
		ProcessSlots(
			st *statedb.StateDB, slot math.Slot,
		) (transition.ValidatorUpdates, error)
		// Transition performs the core state transition.
		Transition(
			ctx core.ReadOnlyContext,
			st *statedb.StateDB,
			blk *ctypes.BeaconBlock,
		) (transition.ValidatorUpdates, error)
		GetSignatureVerifierFn(st *statedb.StateDB) (
			func(blk *ctypes.BeaconBlock, signature crypto.BLSSignature) error,
			error,
		)
	}

	SidecarFactory interface {
		// BuildSidecars builds sidecars for a given block and blobs bundle.
		BuildSidecars(
			signedBlk *ctypes.SignedBeaconBlock,
			blobs engineprimitives.BlobsBundle,
		) (datypes.BlobSidecars, error)
	}

	// StorageBackend defines an interface for accessing various storage
	// components required by the beacon node.
	StorageBackend interface {
		AvailabilityStore() *dastore.Store
		BlockStore() *block.KVStore[*ctypes.BeaconBlock]
		DepositStore() *depositdb.KVStore
		// StateFromContext retrieves the beacon state from the given context.
		StateFromContext(context.Context) *statedb.StateDB
	}

	// 	// TelemetrySink is an interface for sending metrics to a telemetry
	// backend.
	// 	TelemetrySink interface {
	// 		// MeasureSince measures the time since the given time.
	// 		MeasureSince(key string, start time.Time, args ...string)
	// 	}

	// 	// Validator represents an interface for a validator with generic type
	// 	// ValidatorT.
	// 	Validator[
	// 		ValidatorT any,
	// 		WithdrawalCredentialsT any,
	// 	] interface {
	// 		constraints.Empty[ValidatorT]
	// 		constraints.SSZMarshallableRootable
	// 		SizeSSZ() uint32
	// 		// New creates a new validator with the given parameters.
	// 		New(
	// 			pubkey crypto.BLSPubkey,
	// 			withdrawalCredentials WithdrawalCredentialsT,
	// 			amount math.Gwei,
	// 			effectiveBalanceIncrement math.Gwei,
	// 			maxEffectiveBalance math.Gwei,
	// 		) ValidatorT
	// 		// IsSlashed returns true if the validator is slashed.
	// 		IsSlashed() bool
	// 		// IsActive checks if the validator is active at the given epoch.
	// 		IsActive(epoch math.Epoch) bool
	// 		// GetPubkey returns the public key of the validator.
	// 		GetPubkey() crypto.BLSPubkey
	// 		// GetEffectiveBalance returns the effective balance of the validator
	// in
	// 		// Gwei.
	// 		GetEffectiveBalance() math.Gwei
	// 		// SetEffectiveBalance sets the effective balance of the validator in
	// 		// Gwei.
	// 		SetEffectiveBalance(math.Gwei)
	// 		// GetWithdrawableEpoch returns the epoch when the validator can
	// 		// withdraw.
	// 		GetWithdrawableEpoch() math.Epoch
	// 		// GetWithdrawalCredentials returns the withdrawal credentials of the
	// 		// validator.
	// 		GetWithdrawalCredentials() WithdrawalCredentialsT
	// 		// IsFullyWithdrawable checks if the validator is fully withdrawable
	// 		// given a
	// 		// certain Gwei amount and epoch.
	// 		IsFullyWithdrawable(amount math.Gwei, epoch math.Epoch) bool
	// 		// IsPartiallyWithdrawable checks if the validator is partially
	// 		// withdrawable
	// 		// given two Gwei amounts.
	// 		IsPartiallyWithdrawable(amount1 math.Gwei, amount2 math.Gwei) bool
	// 	}

	// 	Validators[ValidatorT any] interface {
	// 		~[]ValidatorT
	// 		HashTreeRoot() common.Root
	// 	}

	// Withdrawal is the interface for a withdrawal.
	Withdrawal[T any] interface {
		New(
			index math.U64,
			validatorIndex math.ValidatorIndex,
			address common.ExecutionAddress,
			amount math.Gwei,
		) T
		// Equals returns true if the withdrawal is equal to the other.
		Equals(T) bool
		// GetAmount returns the amount of the withdrawal.
		GetAmount() math.Gwei
		// GetIndex returns the public key of the validator.
		GetIndex() math.U64
		// GetValidatorIndex returns the index of the validator.
		GetValidatorIndex() math.ValidatorIndex
		// GetAddress returns the address of the withdrawal.
		GetAddress() common.ExecutionAddress
	}

	// // WithdrawalCredentials represents an interface for withdrawal
	// credentials.
	//
	//	WithdrawalCredentials interface {
	//		~[32]byte
	//		// ToExecutionAddress converts the withdrawal credentials to an
	//		// execution
	//		// address.
	//		ToExecutionAddress() (common.ExecutionAddress, error)
	//	}
)

/* -------------------------------------------------------------------------- */
/*                                BeaconState                                 */
/* -------------------------------------------------------------------------- */

type (
	// BeaconStore is the interface for the beacon store.
	BeaconStore[
		T any,
	] interface {
		// Context returns the context of the key-value store.
		Context() context.Context
		// WithContext returns a new key-value store with the given context.
		WithContext(
			ctx context.Context,
		) T
		// Copy returns a copy of the key-value store.
		Copy(context.Context) T
		// GetLatestExecutionPayloadHeader retrieves the latest execution
		// payload
		// header.
		GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
		// SetLatestExecutionPayloadHeader sets the latest execution payload
		// header.
		SetLatestExecutionPayloadHeader(payloadHeader *ctypes.ExecutionPayloadHeader) error
		// GetEth1DepositIndex retrieves the eth1 deposit index.
		GetEth1DepositIndex() (uint64, error)
		// SetEth1DepositIndex sets the eth1 deposit index.
		SetEth1DepositIndex(
			index uint64,
		) error
		// GetBalance retrieves the balance of a validator.
		GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
		// SetBalance sets the balance of a validator.
		SetBalance(idx math.ValidatorIndex, balance math.Gwei) error
		// GetSlot retrieves the current slot.
		GetSlot() (math.Slot, error)
		// SetSlot sets the current slot.
		SetSlot(slot math.Slot) error
		// GetFork retrieves the fork.
		GetFork() (*ctypes.Fork, error)
		// SetFork sets the fork.
		SetFork(fork *ctypes.Fork) error
		// GetGenesisValidatorsRoot retrieves the genesis validators root.
		GetGenesisValidatorsRoot() (common.Root, error)
		// SetGenesisValidatorsRoot sets the genesis validators root.
		SetGenesisValidatorsRoot(root common.Root) error
		// GetLatestBlockHeader retrieves the latest block header.
		GetLatestBlockHeader() (*ctypes.BeaconBlockHeader, error)
		// SetLatestBlockHeader sets the latest block header.
		SetLatestBlockHeader(header *ctypes.BeaconBlockHeader) error
		// GetBlockRootAtIndex retrieves the block root at the given index.
		GetBlockRootAtIndex(index uint64) (common.Root, error)
		// StateRootAtIndex retrieves the state root at the given index.
		StateRootAtIndex(index uint64) (common.Root, error)
		// GetEth1Data retrieves the eth1 data.
		GetEth1Data() (*ctypes.Eth1Data, error)
		// SetEth1Data sets the eth1 data.
		SetEth1Data(data *ctypes.Eth1Data) error
		// GetValidators retrieves all validators.
		GetValidators() (ctypes.Validators, error)
		// GetBalances retrieves all balances.
		GetBalances() ([]uint64, error)
		// GetNextWithdrawalIndex retrieves the next withdrawal index.
		GetNextWithdrawalIndex() (uint64, error)
		// SetNextWithdrawalIndex sets the next withdrawal index.
		SetNextWithdrawalIndex(index uint64) error
		// GetNextWithdrawalValidatorIndex retrieves the next withdrawal
		// validator
		// index.
		GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
		// SetNextWithdrawalValidatorIndex sets the next withdrawal validator
		// index.
		SetNextWithdrawalValidatorIndex(index math.ValidatorIndex) error
		// GetTotalSlashing retrieves the total slashing.
		GetTotalSlashing() (math.Gwei, error)
		// SetTotalSlashing sets the total slashing.
		SetTotalSlashing(total math.Gwei) error
		// GetRandaoMixAtIndex retrieves the randao mix at the given index.
		GetRandaoMixAtIndex(index uint64) (common.Bytes32, error)
		// GetSlashings retrieves all slashings.
		GetSlashings() ([]math.Gwei, error)
		// SetSlashingAtIndex sets the slashing at the given index.
		SetSlashingAtIndex(index uint64, amount math.Gwei) error
		// GetSlashingAtIndex retrieves the slashing at the given index.
		GetSlashingAtIndex(index uint64) (math.Gwei, error)
		// GetTotalValidators retrieves the total validators.
		GetTotalValidators() (uint64, error)
		// ValidatorByIndex retrieves the validator at the given index.
		ValidatorByIndex(index math.ValidatorIndex) (*ctypes.Validator, error)
		// UpdateBlockRootAtIndex updates the block root at the given index.
		UpdateBlockRootAtIndex(index uint64, root common.Root) error
		// UpdateStateRootAtIndex updates the state root at the given index.
		UpdateStateRootAtIndex(index uint64, root common.Root) error
		// UpdateRandaoMixAtIndex updates the randao mix at the given index.
		UpdateRandaoMixAtIndex(index uint64, mix common.Bytes32) error
		// UpdateValidatorAtIndex updates the validator at the given index.
		UpdateValidatorAtIndex(
			index math.ValidatorIndex,
			validator *ctypes.Validator,
		) error
		// ValidatorIndexByPubkey retrieves the validator index by the given
		// pubkey.
		ValidatorIndexByPubkey(
			pubkey crypto.BLSPubkey,
		) (math.ValidatorIndex, error)
		// AddValidator adds a validator.
		AddValidator(val *ctypes.Validator) error
		// ValidatorIndexByCometBFTAddress retrieves the validator index by the
		// given comet BFT address.
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
	}

	// ReadOnlyBeaconState is the interface for a read-only beacon state.
	ReadOnlyBeaconState interface {
		ReadOnlyEth1Data
		ReadOnlyRandaoMixes
		ReadOnlyStateRoots
		ReadOnlyValidators
		ReadOnlyWithdrawals

		// GetBalances retrieves all balances.
		GetBalances() ([]uint64, error)
		GetBalance(math.ValidatorIndex) (math.Gwei, error)
		GetSlot() (math.Slot, error)
		GetFork() (*ctypes.Fork, error)
		GetGenesisValidatorsRoot() (common.Root, error)
		GetBlockRootAtIndex(uint64) (common.Root, error)
		GetLatestBlockHeader() (*ctypes.BeaconBlockHeader, error)
		GetValidators() (ctypes.Validators, error)
		GetSlashingAtIndex(uint64) (math.Gwei, error)
		GetTotalSlashing() (math.Gwei, error)
		GetNextWithdrawalIndex() (uint64, error)
		GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
		GetTotalValidators() (uint64, error)
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
	}

	// WriteOnlyBeaconState is the interface for a write-only beacon state.
	WriteOnlyBeaconState interface {
		WriteOnlyEth1Data
		WriteOnlyRandaoMixes
		WriteOnlyStateRoots
		WriteOnlyValidators

		SetGenesisValidatorsRoot(root common.Root) error
		SetFork(*ctypes.Fork) error
		SetSlot(math.Slot) error
		UpdateBlockRootAtIndex(uint64, common.Root) error
		SetLatestBlockHeader(*ctypes.BeaconBlockHeader) error
		IncreaseBalance(math.ValidatorIndex, math.Gwei) error
		DecreaseBalance(math.ValidatorIndex, math.Gwei) error
		SetNextWithdrawalIndex(uint64) error
		SetNextWithdrawalValidatorIndex(math.ValidatorIndex) error
		SetTotalSlashing(math.Gwei) error
	}

	// WriteOnlyStateRoots defines a struct which only has write access to state
	// roots methods.
	WriteOnlyStateRoots interface {
		UpdateStateRootAtIndex(uint64, common.Root) error
	}

	// ReadOnlyStateRoots defines a struct which only has read access to state
	// roots
	// methods.
	ReadOnlyStateRoots interface {
		StateRootAtIndex(uint64) (common.Root, error)
	}

	// WriteOnlyRandaoMixes defines a struct which only has write access to
	// randao
	// mixes methods.
	WriteOnlyRandaoMixes interface {
		UpdateRandaoMixAtIndex(uint64, common.Bytes32) error
	}

	// ReadOnlyRandaoMixes defines a struct which only has read access to randao
	// mixes methods.
	ReadOnlyRandaoMixes interface {
		GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
	}

	// WriteOnlyValidators has write access to validator methods.
	WriteOnlyValidators interface {
		UpdateValidatorAtIndex(
			math.ValidatorIndex,
			*ctypes.Validator,
		) error

		AddValidator(*ctypes.Validator) error
	}

	// ReadOnlyValidators has read access to validator methods.
	ReadOnlyValidators interface {
		ValidatorIndexByPubkey(
			crypto.BLSPubkey,
		) (math.ValidatorIndex, error)

		ValidatorByIndex(
			math.ValidatorIndex,
		) (*ctypes.Validator, error)
	}

	// WriteOnlyEth1Data has write access to eth1 data.
	WriteOnlyEth1Data interface {
		SetEth1Data(*ctypes.Eth1Data) error
		SetEth1DepositIndex(uint64) error
		SetLatestExecutionPayloadHeader(*ctypes.ExecutionPayloadHeader) error
	}

	// ReadOnlyEth1Data has read access to eth1 data.
	ReadOnlyEth1Data interface {
		GetEth1Data() (*ctypes.Eth1Data, error)
		GetEth1DepositIndex() (uint64, error)
		GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
	}

	// ReadOnlyWithdrawals only has read access to withdrawal methods.
	ReadOnlyWithdrawals interface {
		EVMInflationWithdrawal(math.Slot) *engineprimitives.Withdrawal
		ExpectedWithdrawals() (engineprimitives.Withdrawals, error)
	}
)

// /* --------------------------------------------------------------------------
// */ /*                                  NodeAPI
//    */ /*
// -------------------------------------------------------------------------- */

type (
	NodeAPIContext interface {
		Bind(any) error
		Validate(any) error
	}

	// Engine is a generic interface for an API engine.
	NodeAPIEngine interface {
		Run(addr string) error
		RegisterRoutes(*handlers.RouteSet, log.Logger)
	}

	NodeAPIBackend interface {
		AttachQueryBackend(node nodecoretypes.ConsensusService)
		GetSlotByBlockRoot(root common.Root) (math.Slot, error)
		GetSlotByStateRoot(root common.Root) (math.Slot, error)
		GetParentSlotByTimestamp(timestamp math.U64) (math.Slot, error)

		NodeAPIBeaconBackend
		NodeAPIProofBackend
		NodeAPIConfigBackend
	}

	// NodeAPIBackend is the interface for backend of the beacon API.
	NodeAPIBeaconBackend interface {
		GenesisBackend
		BlobBackend
		BlockBackend
		RandaoBackend
		StateBackend
		ValidatorBackend
		// GetSlotByBlockRoot retrieves the slot by a given root from the store.
		GetSlotByBlockRoot(root common.Root) (math.Slot, error)
		// GetSlotByStateRoot retrieves the slot by a given root from the store.
		GetSlotByStateRoot(root common.Root) (math.Slot, error)
	}

	// NodeAPIConfigBackend is the interface for backend of the config API.
	NodeAPIConfigBackend interface {
		Spec() (chain.Spec, error)
	}

	// NodeAPIProofBackend is the interface for backend of the proof API.
	NodeAPIProofBackend interface {
		BlockBackend
		StateBackend
		GetParentSlotByTimestamp(timestamp math.U64) (math.Slot, error)
	}

	GenesisBackend interface {
		GenesisValidatorsRoot() (common.Root, error)
		GenesisForkVersion() (common.Version, error)
		GenesisTime() (math.U64, error)
	}

	RandaoBackend interface {
		RandaoAtEpoch(slot math.Slot, epoch math.Epoch) (common.Bytes32, error)
	}

	BlobBackend interface {
		BlobSidecarsByIndices(slot math.Slot, indices []uint64) ([]*types.Sidecar, error)
	}

	BlockBackend interface {
		BlockRootAtSlot(slot math.Slot) (common.Root, error)
		BlockRewardsAtSlot(slot math.Slot) (*types.BlockRewardsData, error)
		BlockHeaderAtSlot(slot math.Slot) (*ctypes.BeaconBlockHeader, error)
	}

	StateBackend interface {
		StateAtSlot(slot math.Slot) (*statedb.StateDB, math.Slot, error)
	}

	ValidatorBackend interface {
		ValidatorByID(
			slot math.Slot, id string,
		) (*types.ValidatorData, error)
		FilteredValidators(
			slot math.Slot,
			ids []string,
			statuses []string,
		) ([]*types.ValidatorData, error)
		ValidatorBalancesByIDs(
			slot math.Slot,
			ids []string,
		) ([]*types.ValidatorBalanceData, error)
	}
)
