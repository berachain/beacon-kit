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
	"encoding/json"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
)

// ABCIMiddleware is a middleware between ABCI and the validator logic.
type ABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BeaconBlockBundleT BeaconBlockBundle[BeaconBlockT, BlobSidecarsT],
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
	dispatcher asynctypes.Dispatcher
	// metrics is the metrics emitter.
	metrics *ABCIMiddlewareMetrics
	// logger is the logger for the middleware.
	logger log.Logger[any]

	subGenDataProcessed      asynctypes.Subscription[asynctypes.Event[transition.ValidatorUpdates]]
	subBuiltBeaconBlock      asynctypes.Subscription[asynctypes.Event[BeaconBlockT]]
	subBuiltSidecars         asynctypes.Subscription[asynctypes.Event[BlobSidecarsT]]
	subBBVerified            asynctypes.Subscription[asynctypes.Event[BeaconBlockT]]
	subSCVerified            asynctypes.Subscription[asynctypes.Event[BlobSidecarsT]]
	subFinalValidatorUpdates asynctypes.Subscription[asynctypes.Event[transition.ValidatorUpdates]]
}

// NewABCIMiddleware creates a new instance of the Handler struct.
func NewABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BeaconBlockBundleT BeaconBlockBundle[BeaconBlockT, BlobSidecarsT],
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
	dispatcher asynctypes.Dispatcher,
) *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBundleT, BlobSidecarsT, DepositT,
	ExecutionPayloadT, GenesisT, SlotDataT,
] {
	return &ABCIMiddleware[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBundleT, BlobSidecarsT, DepositT,
		ExecutionPayloadT, GenesisT, SlotDataT,
	]{
		chainSpec: chainSpec,
		blobGossiper: rp2p.NewNoopBlobHandler[
			BlobSidecarsT, encoding.ABCIRequest,
		](),
		beaconBlockGossiper: rp2p.
			NewNoopBlockGossipHandler[
			BeaconBlockT, encoding.ABCIRequest,
		](
			chainSpec,
		),
		logger:                   logger,
		metrics:                  newABCIMiddlewareMetrics(telemetrySink),
		dispatcher:               dispatcher,
		subGenDataProcessed:      asynctypes.NewSubscription[asynctypes.Event[transition.ValidatorUpdates]](),
		subBuiltBeaconBlock:      asynctypes.NewSubscription[asynctypes.Event[BeaconBlockT]](),
		subBuiltSidecars:         asynctypes.NewSubscription[asynctypes.Event[BlobSidecarsT]](),
		subBBVerified:            asynctypes.NewSubscription[asynctypes.Event[BeaconBlockT]](),
		subSCVerified:            asynctypes.NewSubscription[asynctypes.Event[BlobSidecarsT]](),
		subFinalValidatorUpdates: asynctypes.NewSubscription[asynctypes.Event[transition.ValidatorUpdates]](),
	}
}

func (am *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBundleT, BlobSidecarsT, DepositT,
	ExecutionPayloadT, GenesisT, SlotDataT,
]) Start() error {
	if err := am.dispatcher.Subscribe(
		messages.GenesisDataProcessed, am.subGenDataProcessed,
	); err != nil {
		return err
	}
	if err := am.dispatcher.Subscribe(
		messages.BuiltBeaconBlock, am.subBuiltBeaconBlock,
	); err != nil {
		return err
	}
	if err := am.dispatcher.Subscribe(
		messages.BuiltSidecars, am.subBuiltSidecars,
	); err != nil {
		return err
	}
	if err := am.dispatcher.Subscribe(
		messages.BeaconBlockVerified, am.subBBVerified,
	); err != nil {
		return err
	}
	if err := am.dispatcher.Subscribe(
		messages.SidecarsVerified, am.subSCVerified,
	); err != nil {
		return err
	}
	if err := am.dispatcher.Subscribe(
		messages.FinalValidatorUpdatesProcessed, am.subFinalValidatorUpdates,
	); err != nil {
		return err
	}
	return nil
}

// Name returns the name of the middleware.
func (am *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBundleT, BlobSidecarsT, DepositT,
	ExecutionPayloadT, GenesisT, SlotDataT,
]) Name() string {
	return "abci-middleware"
}
