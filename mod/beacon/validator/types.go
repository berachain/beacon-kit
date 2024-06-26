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

package validator

import (
	"context"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// BeaconBlock represents a beacon block interface.
type BeaconBlock[
	BeaconBlockT any,
	BeaconBlockBodyT BeaconBlockBody[
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	DepositT,
	Eth1DataT,
	ExecutionPayloadT any,
] interface {
	constraints.SSZMarshallable
	// NewWithVersion creates a new beacon block with the given parameters.
	NewWithVersion(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		forkVersion uint32,
	) (BeaconBlockT, error)
	// GetSlot returns the slot of the beacon block.
	GetSlot() math.Slot
	// GetParentBlockRoot returns the parent block root of the beacon block.
	GetParentBlockRoot() common.Root
	// SetStateRoot sets the state root of the beacon block.
	SetStateRoot(common.Root)
	// GetStateRoot returns the state root of the beacon block.
	GetStateRoot() common.Root
	// GetBody returns the body of the beacon block.
	GetBody() BeaconBlockBodyT
}

// BeaconBlockBody represents a beacon block body interface.
type BeaconBlockBody[
	DepositT, Eth1DataT, ExecutionPayloadT any,
] interface {
	constraints.SSZMarshallable
	constraints.Nillable
	// SetRandaoReveal sets the Randao reveal of the beacon block body.
	SetRandaoReveal(crypto.BLSSignature)
	// SetEth1Data sets the Eth1 data of the beacon block body.
	SetEth1Data(Eth1DataT)
	// SetDeposits sets the deposits of the beacon block body.
	SetDeposits([]DepositT)
	// SetExecutionData sets the execution data of the beacon block body.
	SetExecutionData(ExecutionPayloadT) error
	// SetGraffiti sets the graffiti of the beacon block body.
	SetGraffiti(common.Bytes32)
	// SetBlobKzgCommitments sets the blob KZG commitments of the beacon block
	// body.
	SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
}

// BeaconState represents a beacon state interface.
type BeaconState[ExecutionPayloadHeaderT any] interface {
	// GetBlockRootAtIndex returns the block root at the given index.
	GetBlockRootAtIndex(uint64) (common.Root, error)
	// GetLatestExecutionPayloadHeader returns the latest execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT, error,
	)
	// GetSlot returns the current slot of the beacon state.
	GetSlot() (math.Slot, error)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() ([32]byte, error)
	// ValidatorIndexByPubkey returns the validator index by public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
	// GetEth1DepositIndex returns the latest deposit index from the beacon
	// state.
	GetEth1DepositIndex() (uint64, error)
	// GetGenesisValidatorsRoot returns the genesis validators root.
	GetGenesisValidatorsRoot() (common.Root, error)
}

// BlobFactory represents a blob factory interface.
type BlobFactory[
	BeaconBlockT BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BlobSidecarsT,
	DepositT,
	Eth1DataT,
	ExecutionPayloadT any,
] interface {
	// BuildSidecars builds sidecars for a given block and blobs bundle.
	BuildSidecars(
		blk BeaconBlockT,
		blobs engineprimitives.BlobsBundle,
	) (BlobSidecarsT, error)
}

// DepositStore defines the interface for deposit storage.
type DepositStore[DepositT any] interface {
	// GetDepositsByIndex returns `numView` expected deposits.
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) ([]DepositT, error)
}

// Eth1Data represents the eth1 data interface.
type Eth1Data[T any] interface {
	// New creates a new eth1 data with the given parameters.
	New(
		depositRoot common.Root,
		depositCount math.U64,
		blockHash common.ExecutionHash,
	) T
}

// ExecutionPayloadHeader represents the execution payload header interface.
type ExecutionPayloadHeader interface {
	// GetTimestamp returns the timestamp of the execution payload header.
	GetTimestamp() math.U64
	// GetBlockHash returns the block hash of the execution payload header.
	GetBlockHash() common.ExecutionHash
	// GetParentHash returns the parent hash of the execution payload header.
	GetParentHash() common.ExecutionHash
}

// EventSubscription represents the event subscription interface.
type EventSubscription[T any] chan T

// EventPublisher represents the event publisher interface.
type EventPublisher[T any] interface {
	// PublishEvent publishes an event.
	Publish(context.Context, T) error
}

// ForkData represents the fork data interface.
type ForkData[T any] interface {
	// New creates a new fork data with the given parameters.
	New(
		common.Version,
		common.Root,
	) T
	// ComputeRandaoSigningRoot computes the Randao signing root.
	ComputeRandaoSigningRoot(
		common.DomainType,
		math.Epoch,
	) (common.Root, error)
}

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder[BeaconStateT, ExecutionPayloadT any] interface {
	// RetrievePayload retrieves the payload for the given slot.
	RetrievePayload(
		ctx context.Context,
		slot math.Slot,
		parentBlockRoot common.Root,
	) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error)
	// RequestPayloadSync requests a payload for the given slot and
	// blocks until the payload is delivered.
	RequestPayloadSync(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot common.Root,
		headEth1BlockHash common.ExecutionHash,
		finalEth1BlockHash common.ExecutionHash,
	) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error)
}

// StateProcessor defines the interface for processing the state.
type StateProcessor[
	BeaconBlockT any,
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	ContextT,
	ExecutionPayloadHeaderT any,
] interface {
	// ProcessSlot processes the slot.
	ProcessSlots(
		st BeaconStateT, slot math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition performs the core state transition.
	Transition(
		ctx ContextT,
		st BeaconStateT,
		blk BeaconBlockT,
	) (transition.ValidatorUpdates, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend[
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	ExecutionPayloadHeaderT any,
] interface {
	// DepositStore retrieves the deposit store.
	DepositStore(context.Context) DepositStoreT
	// StateFromContext retrieves the beacon state from the context.
	StateFromContext(context.Context) BeaconStateT
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// IncrementCounter increments a counter metric identified by the provided
	// keys.
	IncrementCounter(key string, args ...string)
	// MeasureSince measures the time since the provided start time,
	// identified by the provided keys.
	MeasureSince(key string, start time.Time, args ...string)
}
