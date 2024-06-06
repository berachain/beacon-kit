// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/events"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/ethereum/go-ethereum/event"
)

type BeaconState = core.BeaconState[
	*types.BeaconBlockHeader, *types.ExecutionPayloadHeader, *types.Fork,
	*types.Validator, *engineprimitives.Withdrawal,
]

// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
type BeaconKitRuntime = runtime.BeaconKitRuntime[
	*dastore.Store[types.BeaconBlockBody],
	*types.BeaconBlock,
	types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
]

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // bullish.
func ProvideRuntime(
	cfg *config.Config,
	blobProofVerifier kzg.BlobProofVerifier,
	chainSpec primitives.ChainSpec,
	signer crypto.BLSSigner,
	engineClient *engineclient.EngineClient[*types.ExecutionPayload],
	executionEngine *execution.Engine[*types.ExecutionPayload],
	beaconDepositContract *deposit.WrappedBeaconDepositContract[
		*types.Deposit, types.WithdrawalCredentials,
	],
	localBuilder *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	],
	blockFeed *event.FeedOf[events.Block[*types.BeaconBlock]],
	storageBackend blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
	stateProcessor blockchain.StateProcessor[
		*types.BeaconBlock,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
	],
	telemetrySink *metrics.TelemetrySink,
	logger log.Logger,
	ServiceRegistry *service.Registry,
) (*BeaconKitRuntime, error) {

	// Pass all the services and options into the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime[
		*dastore.Store[types.BeaconBlockBody],
		*types.BeaconBlock,
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		blockchain.StorageBackend[
			*dastore.Store[types.BeaconBlockBody],
			types.BeaconBlockBody,
			BeaconState,
			*datypes.BlobSidecars,
			*types.Deposit,
			*depositdb.KVStore[*types.Deposit],
		],
	](
		chainSpec,
		logger,
		ServiceRegistry,
		storageBackend,
		telemetrySink,
	)
}
