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
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/middleware"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// ABCIMiddlewareInput is the input for the validator middleware provider.
type ABCIMiddlewareInput[
	BeaconBlockT any,
	BlobSidecarsT any,
	LoggerT log.Logger[any],
] struct {
	depinject.In
	ChainSpec     common.ChainSpec
	Logger        LoggerT
	TelemetrySink *metrics.TelemetrySink
	Dispatcher    *Dispatcher
}

// ProvideABCIMiddleware is a depinject provider for the validator
// middleware.
func ProvideABCIMiddleware[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT any,
	BeaconBlockHeaderT any,
	BlobSidecarT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT, BlobSidecarT],
	DepositT any,
	GenesisT Genesis[DepositT, *ExecutionPayloadHeader],
	LoggerT log.Logger[any],
](
	in ABCIMiddlewareInput[BeaconBlockT, BlobSidecarsT, LoggerT],
) (*middleware.ABCIMiddleware[
	BeaconBlockT, BlobSidecarsT, GenesisT, *SlotData,
], error) {
	return middleware.NewABCIMiddleware[
		BeaconBlockT,
		BlobSidecarsT,
		GenesisT,
		*SlotData,
	](
		in.ChainSpec,
		in.Logger,
		in.TelemetrySink,
		in.Dispatcher,
	), nil
}
