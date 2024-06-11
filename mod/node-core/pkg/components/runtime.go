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
	"cosmossdk.io/core/log"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// BeaconState is a type alias for the BeaconState.
type BeaconState = core.BeaconState[
	*types.BeaconBlockHeader, *types.Eth1Data,
	*types.ExecutionPayloadHeader, *types.Fork,
	*types.Validator, *engineprimitives.Withdrawal,
]

// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
type BeaconKitRuntime = runtime.BeaconKitRuntime[
	*dastore.Store[*types.BeaconBlockBody],
	*types.BeaconBlock,
	*types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
]

// RuntimeInput is the input for the runtime provider.
type RuntimeInput struct {
	depinject.In
	ChainSpec       primitives.ChainSpec
	Logger          log.Logger
	ServiceRegistry *service.Registry
	StorageBackend  StorageBackend
	TelemetrySink   *metrics.TelemetrySink
}

// ProvideRuntime is a depinject provider that returns a BeaconKitRuntime.
func ProvideRuntime(
	in RuntimeInput,
) (*BeaconKitRuntime, error) {
	// Build the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		blockchain.StorageBackend[
			*dastore.Store[*types.BeaconBlockBody],
			*types.BeaconBlockBody,
			BeaconState,
			*datypes.BlobSidecars,
			*types.Deposit,
			*depositdb.KVStore[*types.Deposit],
		],
	](
		in.ChainSpec,
		in.Logger,
		in.ServiceRegistry,
		in.StorageBackend,
		in.TelemetrySink,
	)
}
