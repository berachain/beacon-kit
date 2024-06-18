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

package middleware

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
)

// ABCIMiddleware is a middleware between ABCI and the validator logic.
type ABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT BeaconState,
	BlobSidecarsT ssz.Marshallable,
] struct {
	// chainSpec is the chain specification.
	chainSpec primitives.ChainSpec
	// chainService represents the blockchain service.
	chainService BlockchainService[BeaconBlockT, BlobSidecarsT]
	// validatorService is the service responsible for building beacon blocks.
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobSidecarsT,
	]
	// TODO: we will eventually gossip the blobs separately from
	// CometBFT, but for now, these are no-op gossipers.
	blobGossiper p2p.PublisherReceiver[
		BlobSidecarsT,
		[]byte,
		encoding.ABCIRequest,
		BlobSidecarsT,
	]
	// TODO: we will eventually gossip the blocks separately from
	// CometBFT, but for now, these are no-op gossipers.
	beaconBlockGossiper p2p.PublisherReceiver[
		BeaconBlockT,
		[]byte,
		encoding.ABCIRequest,
		BeaconBlockT,
	]
	// resChannel is used to communicate the validator updates to the
	// EndBlock method.
	valUpdatesCh chan transition.ValidatorUpdates
	// errCh is used to communicate errors to the EndBlock method.
	errCh chan error
	// metrics is the metrics emitter.
	metrics *ABCIMiddlewareMetrics
	// logger is the logger for the middleware.
	logger log.Logger[any]
}

// NewABCIMiddleware creates a new instance of the Handler struct.
func NewABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT BeaconState,
	BlobSidecarsT ssz.Marshallable,
](
	chainSpec primitives.ChainSpec,
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobSidecarsT,
	],
	chainService BlockchainService[BeaconBlockT, BlobSidecarsT],
	logger log.Logger[any],
	telemetrySink TelemetrySink,
) *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT,
] {
	return &ABCIMiddleware[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT,
	]{
		chainSpec:        chainSpec,
		validatorService: validatorService,
		chainService:     chainService,
		blobGossiper: rp2p.NewNoopBlobHandler[
			BlobSidecarsT, encoding.ABCIRequest](),
		beaconBlockGossiper: rp2p.
			NewNoopBlockGossipHandler[BeaconBlockT, encoding.ABCIRequest](
			chainSpec,
		),
		logger:       logger,
		valUpdatesCh: make(chan transition.ValidatorUpdates),
		errCh:        make(chan error),
		metrics:      newABCIMiddlewareMetrics(telemetrySink),
	}
}
