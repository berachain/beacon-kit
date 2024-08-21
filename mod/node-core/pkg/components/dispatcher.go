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

package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
	asynchelpers "github.com/berachain/beacon-kit/mod/async/pkg/helpers"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

// DispatcherInput is the input for the Dispatcher.
type DispatcherInput[
	LoggerT any,
] struct {
	depinject.In
	Logger LoggerT
}

// ProvideDispatcher provides a new Dispatcher.
func ProvideDispatcher[
	BeaconBlockT any,
	BlobSidecarsT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	in DispatcherInput[LoggerT],
) (*Dispatcher, error) {
	var err error
	d := dispatcher.New(
		in.Logger.With("service", "dispatcher"),
	)
	// Register the GenesisDataReceived event.
	if err = asynchelpers.RegisterEvent[GenesisEvent](
		d, async.GenesisDataReceived,
	); err != nil {
		return nil, err
	}
	// Register the GenesisDataProcessed event.
	if err = asynchelpers.RegisterEvent[ValidatorUpdateEvent](
		d, async.GenesisDataProcessed,
	); err != nil {
		return nil, err
	}
	// Register the NewSlot event.
	if err = asynchelpers.RegisterEvent[SlotEvent](
		d, async.NewSlot,
	); err != nil {
		return nil, err
	}
	// Register the BuiltBeaconBlock event.
	if err = asynchelpers.RegisterEvent[async.Event[BeaconBlockT]](
		d, async.BuiltBeaconBlock,
	); err != nil {
		return nil, err
	}
	// Register the BuiltSidecars event.
	if err = asynchelpers.RegisterEvent[async.Event[BlobSidecarsT]](
		d, async.BuiltSidecars,
	); err != nil {
		return nil, err
	}
	// Register the BeaconBlockReceived event.
	if err = asynchelpers.RegisterEvent[async.Event[BeaconBlockT]](
		d, async.BeaconBlockReceived,
	); err != nil {
		return nil, err
	}
	// Register the SidecarsReceived event.
	if err = asynchelpers.RegisterEvent[async.Event[BlobSidecarsT]](
		d, async.SidecarsReceived,
	); err != nil {
		return nil, err
	}
	// Register the BeaconBlockVerified event.
	if err = asynchelpers.RegisterEvent[async.Event[BeaconBlockT]](
		d, async.BeaconBlockVerified,
	); err != nil {
		return nil, err
	}
	// Register the SidecarsVerified event.
	if err = asynchelpers.RegisterEvent[async.Event[BlobSidecarsT]](
		d, async.SidecarsVerified,
	); err != nil {
		return nil, err
	}
	// Register the FinalBeaconBlockReceived event.
	if err = asynchelpers.RegisterEvent[async.Event[BeaconBlockT]](
		d, async.FinalBeaconBlockReceived,
	); err != nil {
		return nil, err
	}
	// Register the FinalSidecarsReceived event.
	if err = asynchelpers.RegisterEvent[async.Event[BlobSidecarsT]](
		d, async.FinalSidecarsReceived,
	); err != nil {
		return nil, err
	}
	// Register the FinalValidatorUpdatesProcessed event.
	if err = asynchelpers.RegisterEvent[ValidatorUpdateEvent](
		d, async.FinalValidatorUpdatesProcessed,
	); err != nil {
		return nil, err
	}
	// Register the BeaconBlockFinalized event.
	if err = asynchelpers.RegisterEvent[async.Event[BeaconBlockT]](
		d, async.BeaconBlockFinalizedEvent,
	); err != nil {
		return nil, err
	}
	return d, nil
}
