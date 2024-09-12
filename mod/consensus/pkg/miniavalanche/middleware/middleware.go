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
	"time"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	mava "github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

// AwaitTimeout is the timeout for awaiting events.
const AwaitTimeout = 2 * time.Second

type VMMiddleware struct {
	// dispatcher is the central dispatcher to
	dispatcher types.EventDispatcher
	// logger is the logger for the middleware.
	logger log.Logger
	// subGenDataProcessed is the channel to hold GenesisDataProcessed events.
	subGenDataProcessed chan async.Event[mava.ValidatorUpdates]
	// subBuiltBeaconBlock is the channel to hold BuiltBeaconBlock events.
	subBuiltBeaconBlock chan async.Event[mava.BeaconBlockT]
	// subBuiltSidecars is the channel to hold BuiltSidecars events.
	subBuiltSidecars chan async.Event[mava.BlobSidecarsT]
	// subBBVerified is the channel to hold BeaconBlockVerified events.
	subBBVerified chan async.Event[mava.BeaconBlockT]
	// subSCVerified is the channel to hold SidecarsVerified events.
	subSCVerified chan async.Event[mava.BlobSidecarsT]
	// subValidatorUpdates is the channel to hold
	// FinalValidatorUpdatesProcessed events.
	subValidatorUpdates chan async.Event[mava.ValidatorUpdates]
}

func NewABCIMiddleware(
	dispatcher types.EventDispatcher,
	logger log.Logger,
) *VMMiddleware {
	return &VMMiddleware{
		dispatcher:          dispatcher,
		logger:              logger,
		subGenDataProcessed: make(chan async.Event[mava.ValidatorUpdates]),
		subBuiltBeaconBlock: make(chan async.Event[mava.BeaconBlockT]),
		subBuiltSidecars:    make(chan async.Event[mava.BlobSidecarsT]),
		subBBVerified:       make(chan async.Event[mava.BeaconBlockT]),
		subSCVerified:       make(chan async.Event[mava.BlobSidecarsT]),
		subValidatorUpdates: make(chan async.Event[mava.ValidatorUpdates]),
	}
}

// Should this be called upon VM.Initialize or VM.SetState(normalOp) ??
func (vm *VMMiddleware) Start(_ context.Context) error {
	var err error
	if err = vm.dispatcher.Subscribe(
		async.GenesisDataProcessed, vm.subGenDataProcessed,
	); err != nil {
		return err
	}
	if err = vm.dispatcher.Subscribe(
		async.BuiltBeaconBlock, vm.subBuiltBeaconBlock,
	); err != nil {
		return err
	}
	if err = vm.dispatcher.Subscribe(
		async.BuiltSidecars, vm.subBuiltSidecars,
	); err != nil {
		return err
	}
	if err = vm.dispatcher.Subscribe(
		async.BeaconBlockVerified, vm.subBBVerified,
	); err != nil {
		return err
	}
	if err = vm.dispatcher.Subscribe(
		async.SidecarsVerified, vm.subSCVerified,
	); err != nil {
		return err
	}
	if err = vm.dispatcher.Subscribe(
		async.FinalValidatorUpdatesProcessed, vm.subValidatorUpdates,
	); err != nil {
		return err
	}
	return nil
}

// Name returns the name of the middleware.
func (vm *VMMiddleware) Name() string {
	return "abci-middleware"
}
