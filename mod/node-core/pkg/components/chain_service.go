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
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/ethereum/go-ethereum/event"
)

// ChainServiceInput is the input for the chain service provider.
type ChainServiceInput struct {
	depinject.In
	BlobProcessor *dablob.Processor[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
	]
	BlockFeed      *event.FeedOf[*feed.Event[*types.BeaconBlock]]
	ChainSpec      primitives.ChainSpec
	Cfg            *config.Config
	DepositService *deposit.Service[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*feed.Event[*types.BeaconBlock],
		*types.Deposit,
		*types.ExecutionPayload,
		event.Subscription,
		types.WithdrawalCredentials,
	]
	EngineClient    *engineclient.EngineClient[*types.ExecutionPayload]
	ExecutionEngine *execution.Engine[*types.ExecutionPayload]
	LocalBuilder    *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	]
	Logger         log.Logger
	Signer         crypto.BLSSigner
	StateProcessor blockchain.StateProcessor[
		*types.BeaconBlock,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
	]
	StorageBackend blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	TelemetrySink *metrics.TelemetrySink
}

// ProvideChainService is a depinject provider for the blockchain service.
func ProvideChainService(
	in ChainServiceInput,
) *blockchain.Service[
	*dastore.Store[*types.BeaconBlockBody],
	*types.BeaconBlock,
	*types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*types.Deposit,
	*depositdb.KVStore[*types.Deposit],
] {
	return blockchain.NewService[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
	](
		in.StorageBackend,
		in.Logger.With("service", "blockchain"),
		in.ChainSpec,
		in.ExecutionEngine,
		in.LocalBuilder,
		in.BlobProcessor,
		in.StateProcessor,
		in.TelemetrySink,
		in.BlockFeed,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		in.Cfg.Validator.EnableOptimisticPayloadBuilds,
	)
}
