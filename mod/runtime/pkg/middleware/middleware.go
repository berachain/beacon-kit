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
	"context"

	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/beacon"
	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/telemetry"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

// ABCIMiddleware is a middleware between ABCI and the validator logic.
type ABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT types.BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT any,
	BlobSidecarsT constraints.SSZMarshallable,
	DepositT,
	Eth1DataT any,
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT types.Genesis[DepositT, ExecutionPayloadHeaderT],
] struct {
	// chainSpec is the chain specification.
	chainSpec common.ChainSpec
	// chainService represents the blockchain service.
	chainService beacon.BlockchainService[
		BeaconBlockT, BlobSidecarsT, DepositT,
		ExecutionPayloadHeaderT, GenesisT,
	]
	// TODO: we will eventually gossip the blobs separately from
	// CometBFT, but for now, these are no-op gossipers.
	blobGossiper p2p.PublisherReceiver[
		BlobSidecarsT,
		[]byte,
		runtime.ABCIRequest,
		BlobSidecarsT,
	]
	// TODO: we will eventually gossip the blocks separately from
	// CometBFT, but for now, these are no-op gossipers.
	beaconBlockGossiper p2p.PublisherReceiver[
		BeaconBlockT,
		[]byte,
		runtime.ABCIRequest,
		BeaconBlockT,
	]
	// metrics is the metrics emitter.
	metrics *ABCIMiddlewareMetrics
	// logger is the logger for the middleware.
	logger log.Logger[any]

	// Feeds
	//
	// genesisBroker is a feed for genesis data.
	genesisBroker *broker.Broker[*asynctypes.Event[GenesisT]]
	// blkBroker is a feed for blocks.
	blkBroker *broker.Broker[*asynctypes.Event[BeaconBlockT]]
	// sidecarsBroker is a feed for sidecars.
	sidecarsBroker *broker.Broker[*asynctypes.Event[BlobSidecarsT]]
	// slotBroker is a feed for slots.
	slotBroker *broker.Broker[*asynctypes.Event[math.Slot]]

	// TODO: this is a temporary hack.
	req *cmtabci.FinalizeBlockRequest

	// Channels
	// blkCh is used to communicate the beacon block to the EndBlock method.
	blkCh chan *asynctypes.Event[BeaconBlockT]
	// sidecarsCh is used to communicate the sidecars to the EndBlock method.
	sidecarsCh chan *asynctypes.Event[BlobSidecarsT]
	// valUpdateSub is the channel for listening for incoming validator set
	// updates.
	valUpdateSub chan *asynctypes.Event[transition.ValidatorUpdates]
}

// NewABCIMiddleware creates a new instance of the Handler struct.
func NewABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT types.BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT any,
	BlobSidecarsT constraints.SSZMarshallable,
	DepositT,
	Eth1DataT any,
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT types.Genesis[DepositT, ExecutionPayloadHeaderT],
](
	chainSpec common.ChainSpec,
	chainService beacon.BlockchainService[
		BeaconBlockT, BlobSidecarsT, DepositT,
		ExecutionPayloadHeaderT, GenesisT,
	],
	logger log.Logger[any],
	telemetrySink telemetry.Sink,
	genesisBroker *broker.Broker[*asynctypes.Event[GenesisT]],
	blkBroker *broker.Broker[*asynctypes.Event[BeaconBlockT]],
	sidecarsBroker *broker.Broker[*asynctypes.Event[BlobSidecarsT]],
	slotBroker *broker.Broker[*asynctypes.Event[math.Slot]],
	valUpdateSub chan *asynctypes.Event[transition.ValidatorUpdates],
) *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconBlockHeaderT, BlobSidecarsT, DepositT, Eth1DataT,
	ExecutionPayloadT, ExecutionPayloadHeaderT, GenesisT,
] {
	return &ABCIMiddleware[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconBlockHeaderT, BlobSidecarsT, DepositT, Eth1DataT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, GenesisT,
	]{
		chainSpec:    chainSpec,
		chainService: chainService,
		blobGossiper: rp2p.NewNoopBlobHandler[
			BlobSidecarsT, runtime.ABCIRequest,
		](),
		beaconBlockGossiper: rp2p.
			NewNoopBlockGossipHandler[
			BeaconBlockT, runtime.ABCIRequest,
		](
			chainSpec,
		),
		logger:         logger,
		metrics:        newABCIMiddlewareMetrics(telemetrySink),
		genesisBroker:  genesisBroker,
		blkBroker:      blkBroker,
		sidecarsBroker: sidecarsBroker,
		slotBroker:     slotBroker,
		blkCh: make(
			chan *asynctypes.Event[BeaconBlockT],
			1,
		),
		sidecarsCh: make(
			chan *asynctypes.Event[BlobSidecarsT],
			1,
		),
		valUpdateSub: valUpdateSub,
	}
}

// Name returns the name of the middleware.
func (am *ABCIMiddleware[
	_, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "abci-middleware"
}

// Start the middleware.
func (am *ABCIMiddleware[
	_, _, _, _, _, _, _, _, _, _,
]) Start(ctx context.Context) error {
	subBlkCh, err := am.blkBroker.Subscribe()
	if err != nil {
		return err
	}

	subSidecarsCh, err := am.sidecarsBroker.Subscribe()
	if err != nil {
		return err
	}

	go am.start(ctx, subBlkCh, subSidecarsCh)
	return nil
}

// start starts the middleware.
func (am *ABCIMiddleware[
	_, BeaconBlockT, _, _, BlobSidecarsT, _, _, _, _, _,
]) start(
	ctx context.Context,
	blkCh chan *asynctypes.Event[BeaconBlockT],
	sidecarsCh chan *asynctypes.Event[BlobSidecarsT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-blkCh:
			switch msg.Type() {
			case events.BeaconBlockBuilt:
				fallthrough
			case events.BeaconBlockVerified:
				am.blkCh <- msg
			}
		case msg := <-sidecarsCh:
			switch msg.Type() {
			case events.BlobSidecarsBuilt:
				fallthrough
			case events.BlobSidecarsProcessed:
				am.sidecarsCh <- msg
			}
		}
	}
}
