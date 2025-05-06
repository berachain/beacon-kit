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

package validator

import (
	"context"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
)

// BlobFactory represents a blob factory interface.
type BlobFactory interface {
	// BuildSidecars builds sidecars for a given block and blobs bundle.
	BuildSidecars(
		signedBlk *ctypes.SignedBeaconBlock,
		blobs engineprimitives.BlobsBundle,
	) (datypes.BlobSidecars, error)
}

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder interface {
	// Enabled may be enabled (e.g. for validators)
	// or disabled (e.g. full nodes)
	Enabled() bool
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

// StateProcessor defines the interface for processing the state.
type StateProcessor interface {
	// ProcessFork prepares the state for the fork version at the given timestamp.
	ProcessFork(
		st *statedb.StateDB, timestamp math.U64, logUpgrade bool,
	) error
	// ProcessSlots processes the slot.
	ProcessSlots(
		st *statedb.StateDB, slot math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition performs the core state transition.
	Transition(
		ctx core.ReadOnlyContext,
		st *statedb.StateDB,
		blk *ctypes.BeaconBlock,
	) (transition.ValidatorUpdates, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend interface {
	// DepositStore retrieves the deposit store.
	DepositStore() *depositdb.KVStore
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
		*types.SlotData,
	) ([]byte, []byte, error)
}

// ChainSpec defines an interface for accessing chain-specific parameters.
type ChainSpec interface {
	SlotsPerHistoricalRoot() uint64
	DomainTypeRandao() common.DomainType
	MaxDepositsPerBlock() uint64
	ActiveForkVersionForTimestamp(timestamp math.U64) common.Version
	SlotToEpoch(slot math.Slot) math.Epoch
	ctypes.ProposerDomain
}
