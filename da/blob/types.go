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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(context.Context, math.Slot, *ctypes.BeaconBlockBody) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, datypes.BlobSidecars) error
}

type ConsensusSidecars interface {
	GetSidecars() datypes.BlobSidecars
	GetHeader() *ctypes.BeaconBlockHeader
}

type Sidecar interface {
	GetIndex() uint64
	GetBlob() eip4844.Blob
	GetKzgProof() eip4844.KZGProof
	GetKzgCommitment() eip4844.KZGCommitment
	GetSignedBeaconBlockHeader() *ctypes.SignedBeaconBlockHeader
}

type Sidecars[SidecarT any] interface {
	Len() int
	Get(index int) SidecarT
	GetSidecars() []SidecarT
	ValidateBlockRoots() error
	VerifyInclusionProofs(
		kzgOffset uint64,
		inclusionProofDepth uint8,
	) error
}

// ChainSpec represents a chain spec.
type ChainSpec interface {
	MaxBlobCommitmentsPerBlock() uint64
	DomainTypeProposer() common.DomainType
	ActiveForkVersionForSlot(slot math.Slot) uint32
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// MeasureSince measures the time since the provided start time,
	// identified by the provided keys.
	MeasureSince(key string, start time.Time, args ...string)
}
