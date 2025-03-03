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
	"encoding/json"

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
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/block"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
	v1 "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	fastssz "github.com/ferranbt/fastssz"
)

type (
	// AttributesFactory is the interface for the attributes factory.
	AttributesFactory interface {
		BuildPayloadAttributes(
			st *statedb.StateDB,
			slot math.Slot,
			timestamp uint64,
			prevHeadRoot [32]byte,
		) (*engineprimitives.PayloadAttributes, error)
	}

	ConsensusBlock interface {
		GetBeaconBlock() *ctypes.BeaconBlock

		// GetProposerAddress returns the address of the validator
		// selected by consensus to propose the block
		GetProposerAddress() []byte

		// GetConsensusTime returns the timestamp of current consensus request.
		// It is used to build next payload and to validate currentpayload.
		GetConsensusTime() math.U64
	}

	// BeaconBlock represents a generic interface for a beacon block.
	BeaconBlock[
		T any,
	] interface {
		constraints.Nillable
		constraints.Empty[T]
		constraints.Versionable
		constraints.SSZMarshallableRootable

		NewFromSSZ([]byte, uint32) (T, error)
		// NewWithVersion creates a new beacon block with the given parameters.
		NewWithVersion(
			slot math.Slot,
			proposerIndex math.ValidatorIndex,
			parentBlockRoot common.Root,
			forkVersion uint32,
		) (T, error)
		// SetStateRoot sets the state root of the beacon block.
		SetStateRoot(common.Root)
		// GetProposerIndex returns the index of the proposer.
		GetProposerIndex() math.ValidatorIndex
		// GetSlot returns the slot number of the block.
		GetSlot() math.Slot
		// GetBody returns the body of the block.
		GetBody() *ctypes.BeaconBlockBody
		// GetHeader returns the header of the block.
		GetHeader() *ctypes.BeaconBlockHeader
		// GetParentBlockRoot returns the root of the parent block.
		GetParentBlockRoot() common.Root
		// GetStateRoot returns the state root of the block.
		GetStateRoot() common.Root
		// GetTimestamp returns the timestamp of the block from the execution
		// payload.
		GetTimestamp() math.U64
	}

	// BeaconBlockBody represents a generic interface for the body of a beacon
	// block.
	BeaconBlockBody[
		T any,
	] interface {
		constraints.Nillable
		constraints.EmptyWithVersion[T]
		constraints.SSZMarshallableRootable
		Length() uint64
		GetTopLevelRoots() []common.Root
		// GetRandaoReveal returns the RANDAO reveal signature.
		GetRandaoReveal() crypto.BLSSignature
		// GetExecutionPayload returns the execution payload.
		GetExecutionPayload() *ctypes.ExecutionPayload
		// GetDeposits returns the list of deposits.
		GetDeposits() []*ctypes.Deposit
		// GetBlobKzgCommitments returns the KZG commitments for the blobs.
		GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
		// SetRandaoReveal sets the Randao reveal of the beacon block body.
		SetRandaoReveal(crypto.BLSSignature)
		// SetEth1Data sets the Eth1 data of the beacon block body.
		SetEth1Data(*ctypes.Eth1Data)
		// SetDeposits sets the deposits of the beacon block body.
		SetDeposits([]*ctypes.Deposit)
		// SetExecutionPayload sets the execution data of the beacon block body.
		SetExecutionPayload(*ctypes.ExecutionPayload)
		// SetGraffiti sets the graffiti of the beacon block body.
		SetGraffiti(common.Bytes32)
		// SetAttestations sets the attestations of the beacon block body.
		SetAttestations([]*ctypes.AttestationData)
		// SetSlashingInfo sets the slashing info of the beacon block body.
		SetSlashingInfo([]*ctypes.SlashingInfo)
		// SetBlobKzgCommitments sets the blob KZG commitments of the beacon
		// block body.
		SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
	}

	// BeaconStateMarshallable represents an interface for a beacon state
	// with generic types.
	BeaconStateMarshallable[
		T any,
	] interface {
		constraints.SSZMarshallableRootable
		GetTree() (*fastssz.Node, error)
		// New returns a new instance of the BeaconStateMarshallable.
		New(
			forkVersion uint32,
			genesisValidatorsRoot common.Root,
			slot math.U64,
			fork *ctypes.Fork,
			latestBlockHeader *ctypes.BeaconBlockHeader,
			blockRoots []common.Root,
			stateRoots []common.Root,
			eth1Data *ctypes.Eth1Data,
			eth1DepositIndex uint64,
			latestExecutionPayloadHeader *ctypes.ExecutionPayloadHeader,
			validators []*ctypes.Validator,
			balances []uint64,
			randaoMixes []common.Bytes32,
			nextWithdrawalIndex uint64,
			nextWithdrawalValidatorIndex math.U64,
			slashings []math.U64, totalSlashing math.U64,
		) (T, error)
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

	ConsensusSidecars interface {
		GetSidecars() datypes.BlobSidecars
		GetHeader() *ctypes.BeaconBlockHeader
	}

	ConsensusEngine interface {
		PrepareProposal(
			ctx sdk.Context, req *v1.PrepareProposalRequest,
		) (*v1.PrepareProposalResponse, error)
		ProcessProposal(
			ctx sdk.Context, req *v1.ProcessProposalRequest,
		) (*v1.ProcessProposalResponse, error)
	}

	// Deposit is the interface for a deposit.
	Deposit[
		T any,
	] interface {
		constraints.Empty[T]
		constraints.SSZMarshallableRootable
		// New creates a new deposit.
		New(
			crypto.BLSPubkey,
			ctypes.WithdrawalCredentials,
			math.U64,
			crypto.BLSSignature,
			uint64,
		) T
		// Equals returns true if the Deposit is equal to the other.
		Equals(T) bool
		// GetIndex returns the index of the deposit.
		GetIndex() math.U64
		// GetAmount returns the amount of the deposit.
		GetAmount() math.Gwei
		// GetPubkey returns the public key of the validator.
		GetPubkey() crypto.BLSPubkey
		// GetWithdrawalCredentials returns the withdrawal credentials.
		GetWithdrawalCredentials() ctypes.WithdrawalCredentials
		// HasEth1WithdrawalCredentials returns true if the deposit has eth1
		// withdrawal credentials.
		HasEth1WithdrawalCredentials() bool
		// VerifySignature verifies the deposit and creates a validator.
		VerifySignature(
			forkData *ctypes.ForkData,
			domainType common.DomainType,
			signatureVerificationFn func(
				pubkey crypto.BLSPubkey,
				message []byte, signature crypto.BLSSignature,
			) error,
		) error
	}

	// Genesis is the interface for the genesis.
	Genesis interface {
		json.Unmarshaler
		// GetForkVersion returns the fork version.
		GetForkVersion() common.Version
		// GetDeposits returns the deposits.
		GetDeposits() []*ctypes.Deposit
		// GetExecutionPayloadHeader returns the execution payload header.
		GetExecutionPayloadHeader() *ctypes.ExecutionPayloadHeader
	}

	// IndexDB is the interface for the range DB.
	IndexDB interface {
		Has(index uint64, key []byte) (bool, error)
		Get(index uint64, key []byte) ([]byte, error)
		Set(index uint64, key []byte, value []byte) error
		Prune(start uint64, end uint64) error
		GetByIndex(index uint64) ([][]byte, error)
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
			timestamp uint64,
			parentBlockRoot common.Root,
			headEth1BlockHash common.ExecutionHash,
			finalEth1BlockHash common.ExecutionHash,
		) (*engineprimitives.PayloadID, error)
		// SendForceHeadFCU sends a force head FCU request.
		SendForceHeadFCU(
			ctx context.Context,
			st *statedb.StateDB,
			slot math.Slot,
		) error
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
			timestamp uint64,
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
		// InitializePreminedBeaconStateFromEth1 initializes the premined beacon
		// state
		// from the eth1 deposits.
		InitializePreminedBeaconStateFromEth1(
			*statedb.StateDB,
			ctypes.Deposits,
			*ctypes.ExecutionPayloadHeader,
			common.Version,
		) (transition.ValidatorUpdates, error)
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
			blobs ctypes.BlobsBundle,
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
		// GetTotalActiveBalances retrieves the total active balances.
		GetTotalActiveBalances(uint64) (math.Gwei, error)
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
		// GetValidatorsByEffectiveBalance retrieves validators by effective
		// balance.
		GetValidatorsByEffectiveBalance() ([]*ctypes.Validator, error)
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
		GetTotalActiveBalances(uint64) (math.Gwei, error)
		GetValidators() (ctypes.Validators, error)
		GetSlashingAtIndex(uint64) (math.Gwei, error)
		GetTotalSlashing() (math.Gwei, error)
		GetNextWithdrawalIndex() (uint64, error)
		GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
		GetTotalValidators() (uint64, error)
		GetValidatorsByEffectiveBalance() ([]*ctypes.Validator, error)
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
		UpdateSlashingAtIndex(uint64, math.Gwei) error
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
		ChainSpec() chain.Spec
		GetSlotByBlockRoot(root common.Root) (math.Slot, error)
		GetSlotByStateRoot(root common.Root) (math.Slot, error)
		GetParentSlotByTimestamp(timestamp math.U64) (math.Slot, error)

		NodeAPIBeaconBackend
		NodeAPIProofBackend
	}

	// NodeAPIBackend is the interface for backend of the beacon API.
	NodeAPIBeaconBackend interface {
		GenesisBackend
		BlobBackend
		BlockBackend
		RandaoBackend
		StateBackend
		ValidatorBackend
		HistoricalBackend
		// GetSlotByBlockRoot retrieves the slot by a given root from the store.
		GetSlotByBlockRoot(root common.Root) (math.Slot, error)
		// GetSlotByStateRoot retrieves the slot by a given root from the store.
		GetSlotByStateRoot(root common.Root) (math.Slot, error)
	}

	// NodeAPIProofBackend is the interface for backend of the proof API.
	NodeAPIProofBackend interface {
		BlockBackend
		StateBackend
		GetParentSlotByTimestamp(timestamp math.U64) (math.Slot, error)
	}

	GenesisBackend interface {
		GenesisValidatorsRoot(slot math.Slot) (common.Root, error)
	}

	HistoricalBackend interface {
		StateRootAtSlot(slot math.Slot) (common.Root, error)
		StateForkAtSlot(slot math.Slot) (*ctypes.Fork, error)
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
		StateRootAtSlot(slot math.Slot) (common.Root, error)
		StateForkAtSlot(slot math.Slot) (*ctypes.Fork, error)
		StateFromSlotForProof(slot math.Slot) (*statedb.StateDB, math.Slot, error)
		StateAtSlot(slot math.Slot) (*statedb.StateDB, error)
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
