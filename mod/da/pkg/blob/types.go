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

package blob

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[
	BeaconBlockBodyT any,
	BlobSidecarsT any,
	BeaconBlockHeaderT any,
] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(context.Context, math.Slot, BeaconBlockBodyT) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
	// GetBlobSideCars returns the sidecars for the given slot.
	GetBlobSideCars(math.Slot) (*[]backend.BlobSideCar[BeaconBlockHeaderT], error)
}

type BeaconBlock[
	BeaconBlockBodyT any,
	BeaconBlockHeaderT any,
] interface {
	GetBody() BeaconBlockBodyT
	GetHeader() BeaconBlockHeaderT
}

type BeaconBlockBody interface {
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
	GetTopLevelRoots() []common.Root
	Length() uint64
}

type BeaconBlockHeader interface {
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetStateRoot() common.Root
	GetParentBlockRoot() common.Root
	GetBodyRoot() common.Root
}

//nolint:revive // name conflict
type BlobVerifier[BlobSidecarsT any] interface {
	VerifyInclusionProofs(scs BlobSidecarsT, kzgOffset uint64) error
	VerifyKZGProofs(scs BlobSidecarsT) error
	VerifySidecars(sidecars BlobSidecarsT, kzgOffset uint64) error
}

type Sidecar[BeaconBlockHeaderT any] interface {
	GetBeaconBlockHeader() BeaconBlockHeaderT
	GetBlob() eip4844.Blob
	GetKzgProof() eip4844.KZGProof
	GetKzgCommitment() eip4844.KZGCommitment
	GetIndex() uint64
	GetInclusionProof() []common.Root
}

type Sidecars[SidecarT any] interface {
	Len() int
	Get(index int) SidecarT
	GetSidecars() []SidecarT
	ValidateBlockRoots() error
	VerifyInclusionProofs(kzgOffset uint64) error
}

// ChainSpec represents a chain spec.
type ChainSpec interface {
	MaxBlobCommitmentsPerBlock() uint64
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// MeasureSince measures the time since the provided start time,
	// identified by the provided keys.
	MeasureSince(key string, start time.Time, args ...string)
}
