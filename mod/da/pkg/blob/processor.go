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
	"time"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Processor is the blob processor that handles the processing and verification
// of blob sidecars.
type Processor[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader,
	BlobSidecarT Sidecar[BeaconBlockHeaderT],
	BlobSidecarsT Sidecars[BlobSidecarT],
] struct {
	// logger is used to log information and errors.
	logger log.Logger[any]
	// chainSpec defines the specifications of the blockchain.
	chainSpec common.ChainSpec
	// verifier is responsible for verifying the blobs.
	verifier *Verifier[
		BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT,
	]
	// blockBodyOffsetFn is a function that calculates the block body offset
	// based on the slot and chain specifications.
	blockBodyOffsetFn func(math.Slot, common.ChainSpec) uint64
	// metrics is used to collect and report processor metrics.
	metrics *processorMetrics
}

// NewProcessor creates a new blob processor.
func NewProcessor[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader,
	BlobSidecarT Sidecar[BeaconBlockHeaderT],
	BlobSidecarsT Sidecars[BlobSidecarT],
](
	logger log.Logger[any],
	chainSpec common.ChainSpec,
	verifier *Verifier[
		BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT,
	],
	blockBodyOffsetFn func(math.Slot, common.ChainSpec) uint64,
	telemetrySink TelemetrySink,
) *Processor[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BlobSidecarT, BlobSidecarsT,
] {
	return &Processor[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BlobSidecarT, BlobSidecarsT,
	]{
		logger:            logger,
		chainSpec:         chainSpec,
		verifier:          verifier,
		blockBodyOffsetFn: blockBodyOffsetFn,
		metrics:           newProcessorMetrics(telemetrySink),
	}
}

// VerifySidecars verifies the blobs and ensures they match the local state.
func (sp *Processor[AvailabilityStoreT, _, _, _, BlobSidecarsT]) VerifySidecars(
	sidecars BlobSidecarsT,
) error {
	startTime := time.Now()
	defer sp.metrics.measureVerifySidecarsDuration(
		startTime, math.U64(sidecars.Len()),
	)

	// Abort if there are no blobs to store.
	if sidecars.Len() == 0 {
		return nil
	}

	// Verify the blobs and ensure they match the local state.
	return sp.verifier.VerifySidecars(
		sidecars,
		sp.blockBodyOffsetFn(
			sidecars.Get(0).GetBeaconBlockHeader().GetSlot(),
			sp.chainSpec,
		),
	)
}

// slot :=  processes the blobs and ensures they match the local state.
func (sp *Processor[
	AvailabilityStoreT, _, _, _, BlobSidecarsT,
]) ProcessSidecars(
	avs AvailabilityStoreT,
	sidecars BlobSidecarsT,
) error {
	startTime := time.Now()
	defer sp.metrics.measureProcessSidecarsDuration(
		startTime, math.U64(sidecars.Len()),
	)

	// Abort if there are no blobs to store.
	if sidecars.Len() == 0 {
		return nil
	}

	// If we have reached this point, we can safely assume that the blobs are
	// valid and can be persisted, as well as that index 0 is filled.
	return avs.Persist(
		sidecars.Get(0).GetBeaconBlockHeader().GetSlot(),
		sidecars,
	)
}
