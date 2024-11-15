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

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
}

type ConsensusBlock[BeaconBlockT any] interface {
	GetBeaconBlock() BeaconBlockT

	// GetProposerAddress returns the address of the validator
	// selected by consensus to propose the block
	GetProposerAddress() []byte

	// GetNextPayloadTimestamp returns the timestamp proposed by
	// consensus for the next payload to be proposed. It is also
	// used to bound current payload upon validation
	GetNextPayloadTimestamp() math.U64

	// GetConsensusSyncing signals whether consensus is working normally
	// or is still syncing. In the former case we can skip execution payload
	// verification on finalized blocks.
	GetConsensusSyncing() bool
}

// BeaconBlock represents a beacon block interface.
type BeaconBlock[BeaconBlockBodyT any] interface {
	constraints.SSZMarshallableRootable
	constraints.Nillable
	// GetSlot returns the slot of the beacon block.
	GetSlot() math.Slot
	// GetStateRoot returns the state root of the beacon block.
	GetStateRoot() common.Root
	// GetBody returns the body of the beacon block.
	GetBody() BeaconBlockBodyT
}

// BeaconBlockBody represents the interface for the beacon block body.
type BeaconBlockBody[ExecutionPayloadT any] interface {
	constraints.SSZMarshallableRootable
	constraints.Nillable
	// GetExecutionPayload returns the execution payload of the beacon block
	// body.
	GetExecutionPayload() ExecutionPayloadT
}

// BeaconBlockHeader represents the interface for the beacon block header.
type BeaconBlockHeader interface {
	constraints.SSZMarshallableRootable
	// SetStateRoot sets the state root of the beacon block header.
	SetStateRoot(common.Root)
	// GetStateRoot returns the state root of the beacon block header.
	GetStateRoot() common.Root
}

// BlobSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	constraints.SSZMarshallable
	constraints.Nillable
	// Len returns the length of the blobs sidecars.
	Len() int
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine[PayloadAttributesT any] interface {
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *engineprimitives.ForkchoiceUpdateRequest[PayloadAttributesT],
	) (*engineprimitives.PayloadID, *common.ExecutionHash, error)
}

// ExecutionPayload is the interface for the execution payload.
type ExecutionPayload interface {
	ExecutionPayloadHeader
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
type Genesis[DepositT any, ExecutionPayloadHeaderT any] interface {
	// GetForkVersion returns the fork version.
	GetForkVersion() common.Version
	// GetDeposits returns the deposits.
	GetDeposits() []DepositT
	// GetExecutionPayloadHeader returns the execution payload header.
	GetExecutionPayloadHeader() ExecutionPayloadHeaderT
}

// LocalBuilder is the interface for the builder service.
type LocalBuilder[BeaconStateT any] interface {
	// Enabled returns true if the local builder is enabled.
	Enabled() bool
	// RequestPayloadAsync requests a new payload for the given slot.
	RequestPayloadAsync(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot common.Root,
		headEth1BlockHash common.ExecutionHash,
		finalEth1BlockHash common.ExecutionHash,
	) (*engineprimitives.PayloadID, error)
	// SendForceHeadFCU sends a force head FCU request.
	SendForceHeadFCU(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
	) error
}

type PayloadAttributes interface {
	IsNil() bool
	Version() uint32
	GetSuggestedFeeRecipient() common.ExecutionAddress
}

// ReadOnlyBeaconState defines the interface for accessing various components of
// the beacon state.
type ReadOnlyBeaconState[
	T any,
	BeaconBlockHeaderT any,
	ExecutionPayloadHeaderT any,
] interface {
	// Copy creates a copy of the beacon state.
	Copy() T
	// GetLatestBlockHeader returns the most recent block header.
	GetLatestBlockHeader() (
		BeaconBlockHeaderT,
		error,
	)
	// GetLatestExecutionPayloadHeader returns the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT,
		error,
	)
	// GetSlot retrieves the current slot of the beacon state.
	GetSlot() (math.Slot, error)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() common.Root
}

// StateProcessor defines the interface for processing various state transitions
// in the beacon chain.
type StateProcessor[
	BeaconBlockT,
	BeaconStateT,
	ContextT,
	DepositT,
	ExecutionPayloadHeaderT any,
] interface {
	// InitializePreminedBeaconStateFromEth1 initializes the premined beacon
	// state
	// from the eth1 deposits.
	InitializePreminedBeaconStateFromEth1(
		BeaconStateT,
		[]DepositT,
		ExecutionPayloadHeaderT,
		common.Version,
	) (transition.ValidatorUpdates, error)
	// ProcessSlots processes the state transition for a range of slots.
	ProcessSlots(
		BeaconStateT, math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition processes the state transition for a given block.
	Transition(
		ContextT,
		BeaconStateT,
		BeaconBlockT,
	) (transition.ValidatorUpdates, error)
}

// StorageBackend defines an interface for accessing various storage components
// required by the beacon node.
type StorageBackend[
	AvailabilityStoreT any,
	BeaconStateT any,
] interface {
	// AvailabilityStore returns the availability store for the given context.
	AvailabilityStore() AvailabilityStoreT
	// StateFromContext retrieves the beacon state from the given context.
	StateFromContext(context.Context) BeaconStateT
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

type ValidatorUpdates = transition.ValidatorUpdates
