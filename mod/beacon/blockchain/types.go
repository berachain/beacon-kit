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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT any, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// ReadOnlyBeaconState defines the interface for accessing various components of
// the beacon state.
type ReadOnlyBeaconState[T any] interface {
	// GetSlot retrieves the current slot of the beacon state.
	GetSlot() (math.Slot, error)
	// GetLatestExecutionPayloadHeader returns the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		*types.ExecutionPayloadHeader,
		error,
	)
	// GetEth1DepositIndex returns the index of the most recent eth1 deposit.
	GetEth1DepositIndex() (uint64, error)
	// GetLatestBlockHeader returns the most recent block header.
	GetLatestBlockHeader() (
		*types.BeaconBlockHeader,
		error,
	)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() ([32]byte, error)
	// Copy creates a copy of the beacon state.
	Copy() T
	// ValidatorIndexByPubkey finds the index of a validator based on their
	// public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)
}

// BlobProcessor is the interface for the blobs processor.
type BlobProcessor[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BlobSidecarsT any,
] interface {
	// ProcessBlobs processes the blobs and ensures they match the local state.
	ProcessBlobs(
		slot math.Slot,
		avs AvailabilityStoreT,
		sidecars BlobSidecarsT,
	) error

	// VerifyBlobs verifies the blobs and ensures they match the local state.
	VerifyBlobs(
		slot math.Slot,
		sidecars BlobSidecarsT,
	) error
}

// BlobSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	ssz.Marshallable
	IsNil() bool
	Len() int
}

// DepositStore defines the interface for managing deposit operations.
type DepositStore[DepositT any] interface {
	// Prune prunes the deposit store of [start, end)
	Prune(start, end uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []DepositT) error
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *engineprimitives.GetPayloadRequest,
	) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error)
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *engineprimitives.ForkchoiceUpdateRequest,
	) (*engineprimitives.PayloadID, *common.ExecutionHash, error)
	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution client.
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest[
			*types.ExecutionPayload, *engineprimitives.Withdrawal],
	) error
}

// EventFeed is a generic interface for sending events.
type EventFeed[EventT any] interface {
	// Send sends an event and returns the number of
	// subscribers that received it.
	Send(event EventT) int
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
		parentBlockRoot primitives.Root,
		headEth1BlockHash common.ExecutionHash,
		finalEth1BlockHash common.ExecutionHash,
	) (*engineprimitives.PayloadID, error)
	SendForceHeadFCU(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
	) error
}

// StateProcessor defines the interface for processing various state transitions
// in the beacon chain.
type StateProcessor[
	BeaconBlockT,
	BeaconStateT,
	BlobSidecarsT,
	ContextT,
	DepositT any,
] interface {
	// InitializePreminedBeaconStateFromEth1 initializes the premined beacon
	// state
	// from the eth1 deposits.
	InitializePreminedBeaconStateFromEth1(
		BeaconStateT,
		[]DepositT,
		*types.ExecutionPayloadHeader,
		primitives.Version,
	) ([]*transition.ValidatorUpdate, error)
	// ProcessSlots processes the state transition for a range of slots.
	ProcessSlots(
		BeaconStateT, math.Slot,
	) ([]*transition.ValidatorUpdate, error)
	// Transition processes the state transition for a given block.
	Transition(
		ContextT,
		BeaconStateT,
		BeaconBlockT,
	) ([]*transition.ValidatorUpdate, error)
}

// StorageBackend defines an interface for accessing various storage components
// required by the beacon node.
type StorageBackend[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositT any,
	DepositStoreT DepositStore[DepositT],

] interface {
	// AvailabilityStore returns the availability store for the given context.
	AvailabilityStore(context.Context) AvailabilityStoreT
	// StateFromContext retrieves the beacon state from the given context.
	StateFromContext(context.Context) BeaconStateT
	// DepositStore returns the deposit store for the given context.
	DepositStore(context.Context) DepositStoreT
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
