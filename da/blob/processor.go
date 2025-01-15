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

	"github.com/berachain/beacon-kit/chain-spec/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/kzg"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Processor is the blob processor that handles the processing and verification
// of blob sidecars.
type Processor[
	AvailabilityStoreT AvailabilityStore,
	ConsensusSidecarsT ConsensusSidecars,
] struct {
	// logger is used to log information and errors.
	logger log.Logger
	// chainSpec defines the specifications of the blockchain.
	chainSpec chain.ChainSpec
	// verifier is responsible for verifying the blobs.
	verifier *verifier
	// metrics is used to collect and report processor metrics.
	metrics *processorMetrics
}

// NewProcessor creates a new blob processor.
func NewProcessor[
	AvailabilityStoreT AvailabilityStore,
	ConsensusSidecarsT ConsensusSidecars,
](
	logger log.Logger,
	chainSpec chain.ChainSpec,
	proofVerifier kzg.BlobProofVerifier,
	telemetrySink TelemetrySink,
) *Processor[
	AvailabilityStoreT,
	ConsensusSidecarsT,
] {
	verifier := newVerifier(proofVerifier, telemetrySink, chainSpec)

	return &Processor[
		AvailabilityStoreT,
		ConsensusSidecarsT,
	]{
		logger:    logger,
		chainSpec: chainSpec,
		verifier:  verifier,
		metrics:   newProcessorMetrics(telemetrySink),
	}
}

// VerifySidecars verifies the blobs and ensures they match the local state.
func (sp *Processor[
	AvailabilityStoreT, ConsensusSidecarsT,
]) VerifySidecars(
	cs ConsensusSidecarsT,
	verifierFn func(
		blkHeader *ctypes.BeaconBlockHeader,
		signature crypto.BLSSignature,
	) error,
) error {
	var (
		sidecars  = cs.GetSidecars()
		blkHeader = cs.GetHeader()
	)
	defer sp.metrics.measureVerifySidecarsDuration(
		time.Now(), math.U64(len(sidecars)),
	)

	// Abort if there are no blobs to store.
	if len(sidecars) == 0 {
		return nil
	}

	// Verify the blobs and ensure they match the local state.
	return sp.verifier.verifySidecars(
		sidecars,
		blkHeader,
		verifierFn,
	)
}

// ProcessSidecars processes the blobs and ensures they match the local state.
func (sp *Processor[
	AvailabilityStoreT, _,
]) ProcessSidecars(
	avs AvailabilityStoreT,
	sidecars datypes.BlobSidecars,
) error {
	defer sp.metrics.measureProcessSidecarsDuration(
		time.Now(), math.U64(len(sidecars)),
	)

	// Abort if there are no blobs to store.
	if len(sidecars) == 0 {
		return nil
	}

	// If we have reached this point, we can safely assume that the blobs are
	// valid and can be persisted, as well as that index 0 is filled.
	return avs.Persist(
		sidecars[0].GetSignedBeaconBlockHeader().GetHeader().GetSlot(),
		sidecars,
	)
}
