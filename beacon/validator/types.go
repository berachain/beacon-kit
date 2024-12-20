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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// BeaconBlock represents a beacon block interface.
type BeaconBlock[
	T any,
] interface {
	constraints.SSZMarshallable
	// NewWithVersion creates a new beacon block with the given parameters.
	NewWithVersion(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		forkVersion uint32,
	) (T, error)
	// GetSlot returns the slot of the beacon block.
	GetSlot() math.Slot
	// GetParentBlockRoot returns the parent block root of the beacon block.
	GetParentBlockRoot() common.Root
	// SetStateRoot sets the state root of the beacon block.
	SetStateRoot(common.Root)
	// GetStateRoot returns the state root of the beacon block.
	GetStateRoot() common.Root
	// GetBody returns the body of the beacon block.
	GetBody() *ctypes.BeaconBlockBody
}

// BeaconBlockBody represents a beacon block body interface.
type BeaconBlockBody interface {
	constraints.SSZMarshallable
	constraints.Nillable
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
	// SetBlobKzgCommitments sets the blob KZG commitments of the beacon block
	// body.
	SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
}

// BeaconState represents a beacon state interface.
type BeaconState interface {
	// GetBlockRootAtIndex returns the block root at the given index.
	GetBlockRootAtIndex(uint64) (common.Root, error)
	// GetLatestExecutionPayloadHeader returns the latest execution payload
	// header.
	GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
	// GetSlot returns the current slot of the beacon state.
	GetSlot() (math.Slot, error)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() common.Root
	// ValidatorIndexByPubkey returns the validator index by public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
	// GetEth1DepositIndex returns the latest deposit index from the beacon
	// state.
	GetEth1DepositIndex() (uint64, error)
	// GetGenesisValidatorsRoot returns the genesis validators root.
	GetGenesisValidatorsRoot() (common.Root, error)
}

// BlobFactory represents a blob factory interface.
type BlobFactory interface {
	// BuildSidecars builds sidecars for a given block and blobs bundle.
	BuildSidecars(
		blk *ctypes.BeaconBlock,
		blobs ctypes.BlobsBundle,
		signer crypto.BLSSigner,
		forkData *ctypes.ForkData,
	) (datypes.BlobSidecars, error)
}

// DepositStore defines the interface for deposit storage.
type DepositStore interface {
	// GetDepositsByIndex returns `numView` expected deposits.
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) (ctypes.Deposits, error)
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
	) common.Root
	// ComputeDomain computes the fork data domain for a given domain type.
	ComputeDomain(common.DomainType) common.Domain
}

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder interface {
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

// SlotData represents the slot data interface.
type SlotData interface {
	// GetSlot returns the slot of the incoming slot.
	GetSlot() math.Slot
	// GetAttestationData returns the attestation data of the incoming slot.
	GetAttestationData() []*ctypes.AttestationData
	// GetSlashingInfo returns the slashing info of the incoming slot.
	GetSlashingInfo() []*ctypes.SlashingInfo
	// GetProposerAddress returns the address of the validator
	// selected by consensus to propose the block
	GetProposerAddress() []byte
	// GetConsensusTime returns the timestamp of current consensus request.
	// It is used to build next payload and to validate currentpayload.
	GetConsensusTime() math.U64
}

// StateProcessor defines the interface for processing the state.
type StateProcessor[
	ContextT any,
] interface {
	// ProcessSlot processes the slot.
	ProcessSlots(
		st *statedb.StateDB, slot math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition performs the core state transition.
	Transition(
		ctx ContextT,
		st *statedb.StateDB,
		blk *ctypes.BeaconBlock,
	) (transition.ValidatorUpdates, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend[
	DepositStoreT any,
] interface {
	// DepositStore retrieves the deposit store.
	DepositStore() DepositStoreT
	// StateFromContext retrieves the beacon state from the context.
	StateFromContext(context.Context) *statedb.StateDB
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

type BlockBuilderI interface {
	BuildBlockAndSidecars(
		context.Context,
		types.SlotData,
	) ([]byte, []byte, error)
}
