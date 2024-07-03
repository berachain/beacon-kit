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
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/middleware"
)

// ABCIMiddlewareInput is the input for the validator middleware provider.
type ABCIMiddlewareInput struct {
	depinject.In
	BeaconBlockFeed       *BlockBroker
	ChainSpec             common.ChainSpec
	GenesisBroker         *GenesisBroker
	Logger                log.Logger[any]
	SidecarsFeed          *SidecarsBroker
	SlotBroker            *SlotBroker
	TelemetrySink         *metrics.TelemetrySink
	ValidatorUpdateBroker *ValidatorUpdateBroker
}

// ProvideABCIMiddleware is a depinject provider for the validator
// middleware.
func ProvideABCIMiddleware(
	in ABCIMiddlewareInput,
) (*ABCIMiddleware, error) {
	validatorUpdatesSub, err := in.ValidatorUpdateBroker.Subscribe()
	if err != nil {
		return nil, err
	}
	return middleware.NewABCIMiddleware[
		*AvailabilityStore, *BeaconBlock, *BlobSidecars,
		*Deposit, *ExecutionPayload, *Genesis,
	](
		in.ChainSpec,
		in.Logger,
		in.TelemetrySink,
		in.GenesisBroker,
		in.BeaconBlockFeed,
		in.SidecarsFeed,
		in.SlotBroker,
		validatorUpdatesSub,
	), nil
}
