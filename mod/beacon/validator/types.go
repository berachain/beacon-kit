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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// BeaconBlock represents a beacon block interface.
type BeaconBlock[BeaconBlockT any, BeaconBlockBodyT BeaconBlockBody[
	*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
]] interface {
	ssz.Marshallable
	// NewWithVersion creates a new beacon block with the given parameters.
	NewWithVersion(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		forkVersion uint32,
	) (BeaconBlockT, error)

	// IsNil checks if the beacon block is nil.
	IsNil() bool
	// Version returns the version of the beacon block.
	Version() uint32
	// GetSlot returns the slot of the beacon block.
	GetSlot() math.Slot
	// GetProposerIndex returns the proposer index of the beacon block.
	GetProposerIndex() math.ValidatorIndex
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
	ssz.Marshallable
	// IsNil checks if the beacon block body is nil.
	IsNil() bool
	// SetRandaoReveal sets the Randao reveal of the beacon block body.
	SetRandaoReveal(crypto.BLSSignature)
	// SetEth1Data sets the Eth1 data of the beacon block body.
	SetEth1Data(Eth1DataT)
	// GetDeposits returns the deposits of the beacon block body.
	GetDeposits() []DepositT
	// SetDeposits sets the deposits of the beacon block body.
	SetDeposits([]DepositT)
	// SetExecutionData sets the execution data of the beacon block body.
	SetExecutionData(ExecutionPayloadT) error
	// GetBlobKzgCommitments returns the blob KZG commitments of the beacon
	// block body.
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
	// SetBlobKzgCommitments sets the blob KZG commitments of the beacon block
	// body.
	SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
	// GetExecutionPayload returns the execution payload of the beacon block
	// body.
	GetExecutionPayload() ExecutionPayloadT
}

// BeaconState represents a beacon state interface.
type BeaconState[
	BeaconBlockHeader interface{ HashTreeRoot() ([32]byte, error) },
	BeaconStateT, ExecutionPayloadHeaderT any,
] interface {
	// Copy creates a copy of the beacon state.
	Copy() BeaconStateT
	// GetBlockRootAtIndex returns the block root at the given index.
	GetBlockRootAtIndex(uint64) (primitives.Root, error)
	// GetLatestExecutionPayloadHeader returns the latest execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT, error,
	)
	// GetLatestBlockHeader returns the latest block header.
	GetLatestBlockHeader() (
		BeaconBlockHeader,
		error,
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
	GetGenesisValidatorsRoot() (primitives.Root, error)
}

// BlobFactory represents a blob factory interface.
type BlobFactory[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
	],
	BlobSidecarsT BlobSidecars,
] interface {
	// BuildSidecars builds sidecars for a given block and blobs bundle.
	BuildSidecars(
		blk BeaconBlockT,
		blobs engineprimitives.BlobsBundle,
	) (BlobSidecarsT, error)
}

// BlobProcessor represents a blob processor interface.
type BlobProcessor[
	BlobSidecarsT BlobSidecars,
] interface {
	// VerifyBlobs verifies the blobs and ensures they match the local state.
	VerifyBlobs(
		slot math.Slot,
		sidecars BlobSidecarsT,
	) error
}

// BlobSidecars represents a blob sidecars interface.
type BlobSidecars interface {
	// BlobSidecars must be ssz.Marshallable.
	ssz.Marshallable
	// IsNil checks if the blob sidecars is nil.
	IsNil() bool
	// Len returns the length of the blob sidecars.
	Len() int
}

// DepositStore defines the interface for deposit storage.
type DepositStore[DepositT any] interface {
	// GetDepositsByIndex returns `numView` expected deposits.
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) ([]DepositT, error)
}

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder[BeaconStateT, ExecutionPayloadT any] interface {
	// Enabled returns true if the payload builder is enabled.
	Enabled() bool
	// RetrievePayload retrieves the payload for the given slot.
	RetrievePayload(
		ctx context.Context,
		slot math.Slot,
		parentBlockRoot primitives.Root,
	) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error)
	// RequestPayloadAsync requests a payload for the given slot and returns
	// immediately.
	RequestPayloadAsync(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot primitives.Root,
		headEth1BlockHash common.ExecutionHash,
		finalEth1BlockHash common.ExecutionHash,
	) (*engineprimitives.PayloadID, error)
	// RequestPayloadSync requests a payload for the given slot and
	// blocks until the payload is delivered.
	RequestPayloadSync(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot primitives.Root,
		headEth1BlockHash common.ExecutionHash,
		finalEth1BlockHash common.ExecutionHash,
	) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error)
	// SendForceHeadFCU sends a force head FCU to the execution client.
	SendForceHeadFCU(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
	) error
}

// StateProcessor defines the interface for processing the state.
type StateProcessor[
	BeaconBlockT any,
	BeaconStateT BeaconState[
		*types.BeaconBlockHeader,
		BeaconStateT,
		*types.ExecutionPayloadHeader,
	],
	ContextT any,
] interface {
	// ProcessSlot processes the slot.
	ProcessSlots(
		st BeaconStateT, slot math.Slot,
	) ([]*transition.ValidatorUpdate, error)

	// Transition performs the core state transition.
	Transition(
		ctx ContextT,
		st BeaconStateT,
		blk BeaconBlockT,
	) ([]*transition.ValidatorUpdate, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend[
	BeaconStateT BeaconState[
		*types.BeaconBlockHeader,
		BeaconStateT,
		*types.ExecutionPayloadHeader,
	],
	DepositT any,
	DepositStoreT DepositStore[DepositT],
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
