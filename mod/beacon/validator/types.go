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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	ssz "github.com/ferranbt/fastssz"
)

// BeaconBlock is the interface for a beacon block.
type BeaconBlock[BeaconBlockT any, BeaconBlockBodyT BeaconBlockBody[
	*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
]] interface {
	NewWithVersion(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		forkVersion uint32,
	) (BeaconBlockT, error)
	SetStateRoot(common.Root)
	GetStateRoot() common.Root
	ReadOnlyBeaconBlock[BeaconBlockBodyT]
}

// ReadOnlyBeaconBlock is the interface for a read-only beacon block.
type ReadOnlyBeaconBlock[
	BodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
	]] interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	Version() uint32
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	GetBody() BodyT
}

type BeaconBlockBody[
	DepositT, Eth1DataT, ExecutionPayloadT any,
] interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	SetRandaoReveal(crypto.BLSSignature)
	SetEth1Data(Eth1DataT)
	GetDeposits() []DepositT
	SetDeposits([]DepositT)
	SetExecutionData(ExecutionPayloadT) error
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
	SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
	GetExecutionPayload() ExecutionPayloadT
}

// BeaconState defines the interface for accessing various components of the
// beacon state.
type BeaconState[BeaconStateT any] interface {
	Copy() BeaconStateT
	// GetBlockRootAtIndex fetches the block root at a specified index.
	GetBlockRootAtIndex(uint64) (primitives.Root, error)
	// GetLatestExecutionPayloadHeader returns the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		*types.ExecutionPayloadHeader, error,
	)
	// GetLatestBlockHeader
	GetLatestBlockHeader() (
		*types.BeaconBlockHeader,
		error,
	)
	// GetSlot retrieves the current slot of the beacon state.
	GetSlot() (math.Slot, error)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() ([32]byte, error)
	// ValidatorIndexByPubkey finds the index of a validator based on their
	// public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
	// GetEth1DepositIndex retrieves the latest deposit index from the
	// beacon state.
	GetEth1DepositIndex() (uint64, error)
	// GetGenesisValidatorsRoot retrieves the genesis validators root.
	GetGenesisValidatorsRoot() (primitives.Root, error)
}

// BlobFactory is the interface for building blobs.
type BlobFactory[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
	],
	BlobSidecarsT BlobSidecars,
] interface {
	// BuildSidecars generates sidecars for a given block and blobs bundle.
	BuildSidecars(
		blk BeaconBlockT,
		blobs engineprimitives.BlobsBundle,
	) (BlobSidecarsT, error)
}

// BlobSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	ssz.Marshaler
	ssz.Unmarshaler
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
type PayloadBuilder[BeaconStateT BeaconState[BeaconStateT]] interface {
	// RetrievePayload retrieves the payload for the given slot.
	RetrievePayload(
		ctx context.Context,
		slot math.Slot,
		parentBlockRoot primitives.Root,
	) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error)
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
	BeaconStateT BeaconState[BeaconStateT],
	ContextT any,
] interface {
	// ProcessSlot processes the slot.
	ProcessSlot(
		st BeaconStateT,
	) ([]*transition.ValidatorUpdate, error)

	// Transition performs the core state transition.
	Transition(
		ctx ContextT,
		st BeaconStateT,
		blk BeaconBlockT,
	) ([]*transition.ValidatorUpdate, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend[BeaconStateT BeaconState[BeaconStateT]] interface {
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
