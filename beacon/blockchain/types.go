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

package blockchain

import (
	"context"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, *ctypes.BeaconBlockBody,
	) bool
	// Prune prunes the deposit store of [start, end)
	Prune(start, end uint64) error
}

type ConsensusBlock interface {
	GetBeaconBlock() *ctypes.BeaconBlock

	// GetProposerAddress returns the address of the validator
	// selected by consensus to propose the block
	GetProposerAddress() []byte

	// GetConsensusTime returns the timestamp of current consensus request.
	// It is used to build next payload and to validate currentpayload.
	GetConsensusTime() math.U64
}

type BlobSidecars[T any] interface {
	constraints.SSZMarshallable
	constraints.Empty[T]
	constraints.Nillable
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *ctypes.ForkchoiceUpdateRequest,
	) (*engineprimitives.PayloadID, *common.ExecutionHash, error)
}

// ExecutionPayload is the interface for the execution payload.
type ExecutionPayload interface {
	ExecutionPayloadHeader
	GetNumber() math.U64
}

// ExecutionPayloadHeader is the interface for the execution payload header.
type ExecutionPayloadHeader interface {
	// GetTimestamp returns the timestamp.
	GetTimestamp() math.U64
	// GetBlockHash returns the block hash.
	GetBlockHash() common.ExecutionHash
	// GetParentHash returns the parent hash.
	GetParentHash() common.ExecutionHash
}

// Genesis is the interface for the genesis.
type Genesis interface {
	// GetForkVersion returns the fork version.
	GetForkVersion() common.Version
	// GetDeposits returns the deposits.
	GetDeposits() []*ctypes.Deposit
	// GetExecutionPayloadHeader returns the execution payload header.
	GetExecutionPayloadHeader() *ctypes.ExecutionPayloadHeader
}

// LocalBuilder is the interface for the builder service.
type LocalBuilder interface {
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
}

// ReadOnlyBeaconState defines the interface for accessing various components of
// the beacon state.
type ReadOnlyBeaconState[
	T any,
] interface {
	// Copy creates a copy of the beacon state.
	Copy(context.Context) T
	// GetLatestBlockHeader returns the most recent block header.
	GetLatestBlockHeader() (
		*ctypes.BeaconBlockHeader,
		error,
	)
	// GetLatestExecutionPayloadHeader returns the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
	// GetSlot retrieves the current slot of the beacon state.
	GetSlot() (math.Slot, error)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() common.Root
}

// StateProcessor defines the interface for processing various state transitions
// in the beacon chain.
type StateProcessor[
	ContextT any,
] interface {
	// InitializePreminedBeaconStateFromEth1 initializes the premined beacon
	// state
	// from the eth1 deposits.
	InitializePreminedBeaconStateFromEth1(
		*statedb.StateDB,
		ctypes.Deposits,
		*ctypes.ExecutionPayloadHeader,
		common.Version,
	) (transition.ValidatorUpdates, error)
	// ProcessSlots processes the state transition for a range of slots.
	ProcessSlots(
		*statedb.StateDB, math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition processes the state transition for a given block.
	Transition(
		ContextT,
		*statedb.StateDB,
		*ctypes.BeaconBlock,
	) (transition.ValidatorUpdates, error)
	GetSidecarVerifierFn(*statedb.StateDB) (
		func(
			blkHeader *ctypes.BeaconBlockHeader,
			signature crypto.BLSSignature) error,
		error,
	)
}

// StorageBackend defines an interface for accessing various storage components
// required by the beacon node.
type StorageBackend[
	AvailabilityStoreT any,
	BlockStoreT any,
	DepositStoreT any,
] interface {
	// AvailabilityStore returns the availability store for the given context.
	AvailabilityStore() AvailabilityStoreT
	// StateFromContext retrieves the beacon state from the given context.
	StateFromContext(context.Context) *statedb.StateDB
	// DepositStore retrieves the deposit store.
	DepositStore() DepositStoreT
	// BlockStore retrieves the block store.
	BlockStore() BlockStoreT
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// IncrementCounter increments the counter identified by
	// the provided key.
	IncrementCounter(key string, args ...string)

	// MeasureSince measures the time since the provided start time,
	// identified by the provided keys.
	MeasureSince(key string, start time.Time, args ...string)
}

//nolint:revive // its ok
type BlockchainI interface {
	ProcessGenesisData(
		context.Context, []byte) (transition.ValidatorUpdates, error)
	ProcessProposal(
		sdk.Context,
		*cmtabci.ProcessProposalRequest,
	) (*cmtabci.ProcessProposalResponse, error)
	FinalizeBlock(
		sdk.Context,
		*cmtabci.FinalizeBlockRequest,
	) (transition.ValidatorUpdates, error)
}

type ValidatorUpdates = transition.ValidatorUpdates
