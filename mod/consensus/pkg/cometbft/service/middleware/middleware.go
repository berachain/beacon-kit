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

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
)

// ABCIMiddleware is a middleware between ABCI and the validator logic.
type ABCIMiddleware[
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BlobSidecarsT BlobSidecars[BlobSidecarsT],
	GenesisT json.Unmarshaler,
	SlotDataT any,
] struct {
	// chainSpec is the chain specification.
	chainSpec  common.ChainSpec
	dispatcher types.Dispatcher
	// metrics is the metrics emitter.
	metrics *ABCIMiddlewareMetrics
	// logger is the logger for the middleware.
	logger log.Logger[any]
	// subscription channels
	subGenDataProcessed      chan async.Event[validatorUpdates]
	subBuiltBeaconBlock      chan async.Event[BeaconBlockT]
	subBuiltSidecars         chan async.Event[BlobSidecarsT]
	subBBVerified            chan async.Event[BeaconBlockT]
	subSCVerified            chan async.Event[BlobSidecarsT]
	subFinalValidatorUpdates chan async.Event[validatorUpdates]
}

// NewABCIMiddleware creates a new instance of the Handler struct.
func NewABCIMiddleware[
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BlobSidecarsT BlobSidecars[BlobSidecarsT],
	GenesisT json.Unmarshaler,
	SlotDataT any,
](
	chainSpec common.ChainSpec,
	logger log.Logger[any],
	telemetrySink TelemetrySink,
	dispatcher types.Dispatcher,
) *ABCIMiddleware[
	BeaconBlockT, BlobSidecarsT, GenesisT, SlotDataT,
] {
	return &ABCIMiddleware[
		BeaconBlockT, BlobSidecarsT, GenesisT, SlotDataT,
	]{
		chainSpec:                chainSpec,
		logger:                   logger,
		metrics:                  newABCIMiddlewareMetrics(telemetrySink),
		dispatcher:               dispatcher,
		subGenDataProcessed:      make(chan async.Event[validatorUpdates]),
		subBuiltBeaconBlock:      make(chan async.Event[BeaconBlockT]),
		subBuiltSidecars:         make(chan async.Event[BlobSidecarsT]),
		subBBVerified:            make(chan async.Event[BeaconBlockT]),
		subSCVerified:            make(chan async.Event[BlobSidecarsT]),
		subFinalValidatorUpdates: make(chan async.Event[validatorUpdates]),
	}
}

// Start subscribes the middleware to the events it needs to listen for.
func (am *ABCIMiddleware[_, _, _, _]) Start(
	_ context.Context,
) error {
	var err error
	if err = am.dispatcher.Subscribe(
		async.GenesisDataProcessed, am.subGenDataProcessed,
	); err != nil {
		return err
	}
	if err = am.dispatcher.Subscribe(
		async.BuiltBeaconBlock, am.subBuiltBeaconBlock,
	); err != nil {
		return err
	}
	if err = am.dispatcher.Subscribe(
		async.BuiltSidecars, am.subBuiltSidecars,
	); err != nil {
		return err
	}
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
func (am *ABCIMiddleware[
	_, _, _, _,
]) Name() string {
	return "abci-middleware"
}
