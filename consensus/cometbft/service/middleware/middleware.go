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

	"github.com/berachain/beacon-kit/async/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/async"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
)

// ABCIMiddleware is a middleware between ABCI and the validator logic.
type ABCIMiddleware[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockHeaderT],
	BeaconBlockHeaderT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT],
	GenesisT json.Unmarshaler,
	SlotDataT any,
] struct {
	// chainSpec is the chain specification.
	chainSpec common.ChainSpec
	// dispatcher is the central dispatcher to
	dispatcher types.EventDispatcher
	// metrics is the metrics emitter.
	metrics *ABCIMiddlewareMetrics
	// logger is the logger for the middleware.
	logger log.Logger
	// subBBVerified is the channel to hold BeaconBlockVerified events.
	subBBVerified chan async.Event[BeaconBlockT]
	// subSCVerified is the channel to hold SidecarsVerified events.
	subSCVerified chan async.Event[BlobSidecarsT]
	// subFinalValidatorUpdates is the channel to hold
	// FinalValidatorUpdatesProcessed events.
	subFinalValidatorUpdates chan async.Event[validatorUpdates]
}

// NewABCIMiddleware creates a new instance of the Handler struct.
func NewABCIMiddleware[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockHeaderT],
	BeaconBlockHeaderT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT],
	GenesisT json.Unmarshaler,
	SlotDataT any,
](
	chainSpec common.ChainSpec,
	dispatcher types.EventDispatcher,
	logger log.Logger,
	telemetrySink TelemetrySink,
) *ABCIMiddleware[
	BeaconBlockT, BeaconBlockHeaderT, BlobSidecarsT, GenesisT, SlotDataT,
] {
	return &ABCIMiddleware[
		BeaconBlockT, BeaconBlockHeaderT, BlobSidecarsT, GenesisT, SlotDataT,
	]{
		chainSpec:                chainSpec,
		dispatcher:               dispatcher,
		logger:                   logger,
		metrics:                  newABCIMiddlewareMetrics(telemetrySink),
		subBBVerified:            make(chan async.Event[BeaconBlockT]),
		subSCVerified:            make(chan async.Event[BlobSidecarsT]),
		subFinalValidatorUpdates: make(chan async.Event[validatorUpdates]),
	}
}

// Start subscribes the middleware to the events it needs to listen for.
func (am *ABCIMiddleware[_, _, _, _, _]) Start(
	_ context.Context,
) error {
	var err error
	if err = am.dispatcher.Subscribe(
		async.BeaconBlockVerified, am.subBBVerified,
	); err != nil {
		return err
	}
	if err = am.dispatcher.Subscribe(
		async.SidecarsVerified, am.subSCVerified,
	); err != nil {
		return err
	}
	if err = am.dispatcher.Subscribe(
		async.FinalValidatorUpdatesProcessed, am.subFinalValidatorUpdates,
	); err != nil {
		return err
	}
	return nil
}

// Name returns the name of the middleware.
func (am *ABCIMiddleware[_, _, _, _, _]) Name() string {
	return "abci-middleware"
}
