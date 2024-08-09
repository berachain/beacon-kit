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
	sdklog "cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
)

// DispatcherInput is the input for the Dispatcher.
type DispatcherInput struct {
	depinject.In
	EventServer                   *EventServer
	MessageServer                 *MessageServer
	Logger                        log.AdvancedLogger[any, sdklog.Logger]
	BeaconBlockFinalizedPublisher *BeaconBlockFinalizedPublisher
	Routes                        []asynctypes.MessageRoute
}

// ProvideDispatcher provides a new Dispatcher.
func ProvideDispatcher(
	in DispatcherInput,
) (*Dispatcher, error) {
	d := dispatcher.NewDispatcher(
		in.EventServer,
		in.MessageServer,
		in.Logger.With("service", "dispatcher"),
	)
	d.RegisterPublisher(
		in.BeaconBlockFinalizedPublisher.EventID(),
		in.BeaconBlockFinalizedPublisher,
	)
	for _, route := range in.Routes {
		if err := d.RegisterRoute(route.MessageID(), route); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func DefaultDispatcherComponents() []any {
	return []any{
		ProvideDispatcher,
		ProvideMessageRoutes,
		ProvideBeaconBlockFinalizedPublisher,
		ProvideMessageServer,
		ProvideEventServer,
	}
}
