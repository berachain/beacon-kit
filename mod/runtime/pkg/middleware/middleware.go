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
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
)

// ABCIMiddleware is a middleware between ABCI and the validator logic.
type ABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BlobSidecarsT interface {
		constraints.SSZMarshallable
		Empty() BlobSidecarsT
	},
	DepositT,
	ExecutionPayloadT any,
	GenesisT json.Unmarshaler,
	SlotDataT any,
] struct {
	// chainSpec is the chain specification.
	chainSpec common.ChainSpec
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
	slotBroker *broker.Broker[*asynctypes.Event[SlotDataT]]

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
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BlobSidecarsT interface {
		constraints.SSZMarshallable
		Empty() BlobSidecarsT
	},
	DepositT,
	ExecutionPayloadT any,
	GenesisT json.Unmarshaler,
	SlotDataT any,
](
	chainSpec common.ChainSpec,
	logger log.Logger[any],
	telemetrySink TelemetrySink,
	genesisBroker *broker.Broker[*asynctypes.Event[GenesisT]],
	blkBroker *broker.Broker[*asynctypes.Event[BeaconBlockT]],
	sidecarsBroker *broker.Broker[*asynctypes.Event[BlobSidecarsT]],
	slotBroker *broker.Broker[*asynctypes.Event[SlotDataT]],
	valUpdateSub chan *asynctypes.Event[transition.ValidatorUpdates],
) *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BlobSidecarsT, DepositT,
	ExecutionPayloadT, GenesisT, SlotDataT,
] {
	return &ABCIMiddleware[
		AvailabilityStoreT, BeaconBlockT, BlobSidecarsT, DepositT,
		ExecutionPayloadT, GenesisT, SlotDataT,
	]{
		chainSpec: chainSpec,
		blobGossiper: rp2p.NewNoopBlobHandler[
			BlobSidecarsT, encoding.ABCIRequest,
		](),
		beaconBlockGossiper: rp2p.NewNoopBlockGossipHandler[
			BeaconBlockT, encoding.ABCIRequest,
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
	AvailabilityStoreT, BeaconBlockT, BlobSidecarsT, DepositT,
	ExecutionPayloadT, GenesisT, SlotDataT,
]) Name() string {
	return "abci-middleware"
}

// Start the middleware.
func (am *ABCIMiddleware[
	_, _, _, _, _, _, _,
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
	_, BeaconBlockT, BlobSidecarsT, _, _, _, _,
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
