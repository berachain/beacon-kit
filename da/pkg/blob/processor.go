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

	"github.com/berachain/beacon-kit/da/pkg/kzg"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/pkg/common"
	"github.com/berachain/beacon-kit/primitives/pkg/math"
)

// Processor is the blob processor that handles the processing and verification
// of blob sidecars.
type Processor[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	ConsensusSidecarsT ConsensusSidecars[BlobSidecarsT, BeaconBlockHeaderT],
	BlobSidecarT Sidecar[BeaconBlockHeaderT],
	BlobSidecarsT Sidecars[BlobSidecarT],
] struct {
	// logger is used to log information and errors.
	logger log.Logger
	// chainSpec defines the specifications of the blockchain.
	chainSpec common.ChainSpec
	// verifier is responsible for verifying the blobs.
	verifier *verifier[BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT]
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
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	ConsensusSidecarsT ConsensusSidecars[BlobSidecarsT, BeaconBlockHeaderT],
	BlobSidecarT Sidecar[BeaconBlockHeaderT],
	BlobSidecarsT Sidecars[BlobSidecarT],
](
	logger log.Logger,
	chainSpec common.ChainSpec,
	proofVerifier kzg.BlobProofVerifier,
	blockBodyOffsetFn func(math.Slot, common.ChainSpec) uint64,
	telemetrySink TelemetrySink,
) *Processor[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
	ConsensusSidecarsT, BlobSidecarT, BlobSidecarsT,
] {
	verifier := newVerifier[
		BeaconBlockHeaderT,
		BlobSidecarT,
		BlobSidecarsT,
	](proofVerifier, telemetrySink)
	return &Processor[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
		ConsensusSidecarsT, BlobSidecarT, BlobSidecarsT,
	]{
		logger:            logger,
		chainSpec:         chainSpec,
		verifier:          verifier,
		blockBodyOffsetFn: blockBodyOffsetFn,
		metrics:           newProcessorMetrics(telemetrySink),
	}
}

// VerifySidecars verifies the blobs and ensures they match the local state.
func (sp *Processor[
	AvailabilityStoreT, _, _, ConsensusSidecarsT, _, _,
]) VerifySidecars(
	cs ConsensusSidecarsT,
) error {
	var (
		sidecars  = cs.GetSidecars()
		blkHeader = cs.GetHeader()
	)
	defer sp.metrics.measureVerifySidecarsDuration(
		time.Now(), math.U64(sidecars.Len()),
	)

	// Abort if there are no blobs to store.
	if sidecars.Len() == 0 {
		return nil
	}

	// Verify the blobs and ensure they match the local state.
	return sp.verifier.verifySidecars(
		sidecars,
		sp.blockBodyOffsetFn(
			sidecars.Get(0).GetBeaconBlockHeader().GetSlot(),
			sp.chainSpec,
		),
		blkHeader,
	)
}

// slot :=  processes the blobs and ensures they match the local state.
func (sp *Processor[
	AvailabilityStoreT, _, _, _, _, BlobSidecarsT,
]) ProcessSidecars(
	avs AvailabilityStoreT,
	sidecars BlobSidecarsT,
) error {
	defer sp.metrics.measureProcessSidecarsDuration(
		time.Now(), math.U64(sidecars.Len()),
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
