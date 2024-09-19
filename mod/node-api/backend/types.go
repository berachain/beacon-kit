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

package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// BeaconBlockHeader is the interface for a beacon block header.
type BeaconBlockHeader[BeaconBlockHeaderT any] interface {
	constraints.SSZMarshallableRootable
	New(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		stateRoot common.Root,
		bodyRoot common.Root,
	) BeaconBlockHeaderT
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	SetStateRoot(common.Root)
	GetBodyRoot() common.Root
}

// BeaconState is the interface for the beacon state.
type BeaconState[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT any,
] interface {
	// SetSlot sets the slot on the beacon state.
	SetSlot(math.Slot) error

	core.ReadOnlyBeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT, WithdrawalT,
	]
}

// BlockStore is the interface for block storage.
type BlockStore[BeaconBlockT any] interface {
	// GetSlotByBlockRoot retrieves the slot by a given block root.
	GetSlotByBlockRoot(root common.Root) (math.Slot, error)
	// GetSlotByStateRoot retrieves the slot by a given state root.
	GetSlotByStateRoot(root common.Root) (math.Slot, error)
	// GetSlotByExecutionNumber retrieves the slot by a given execution number.
	GetSlotByExecutionNumber(executionNumber math.U64) (math.Slot, error)
}

// DepositStore defines the interface for deposit storage.
type DepositStore[DepositT any] interface {
	// GetDepositsByIndex returns `numView` expected deposits.
	GetDepositsByIndex(startIndex uint64, numView uint64) ([]DepositT, error)
	// Prune prunes the deposit store of [start, end)
	Prune(start, end uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []DepositT) error
}

// Node is the interface for a node.
type Node[ContextT any] interface {
	// CreateQueryContext creates a query context for a given height and proof
	// flag.
	CreateQueryContext(height int64, prove bool) (ContextT, error)
}

type StateProcessor[BeaconStateT any] interface {
	ProcessSlots(BeaconStateT, math.Slot) (transition.ValidatorUpdates, error)
}

// StorageBackend is the interface for the storage backend.
type StorageBackend[
	AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT any,
] interface {
	AvailabilityStore() AvailabilityStoreT
	BlockStore() BlockStoreT
	DepositStore() DepositStoreT
	StateFromContext(context.Context) BeaconStateT
}

// Validator represents an interface for a validator.
type Validator interface {
	GetWithdrawalCredentials() types.WithdrawalCredentials
	GetPubkey() crypto.BLSPubkey
	GetEffectiveBalance() math.Gwei
	IsSlashed() bool
	GetActivationEligibilityEpoch() math.Epoch
	GetActivationEpoch() math.Epoch
	GetExitEpoch() math.Epoch
	GetWithdrawableEpoch() math.Epoch
}

// Withdrawal represents an interface for a withdrawal.
type Withdrawal[T any] interface {
	New(
		index math.U64,
		validator math.ValidatorIndex,
		address common.ExecutionAddress,
		amount math.Gwei,
	) T
}

type Fork[ForkT any] interface {
	GetPreviousVersion() common.Version
}
