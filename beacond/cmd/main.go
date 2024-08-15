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

package main

import (
	"log/slog"
	"os"

	clibuilder "github.com/berachain/beacon-kit/mod/cli/pkg/builder"
	clicomponents "github.com/berachain/beacon-kit/mod/cli/pkg/components"
	nodebuilder "github.com/berachain/beacon-kit/mod/node-core/pkg/builder"
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"go.uber.org/automaxprocs/maxprocs"
)

type node = types.Node

// run runs the beacon node.
func run() error {
	// Set the uber max procs
	if _, err := maxprocs.Set(); err != nil {
		return err
	}

	// Build the node using the node-core.
	nb := nodebuilder.New(
		// Set the Runtime Components to the Default.
		nodebuilder.WithComponents[node](
			nodecomponents.DefaultComponents[
				*AttestationData,
				*AttributesFactory,
				*AvailabilityStore,
				*BeaconBlock,
				*BeaconBlockBody,
				*BeaconBlockHeader,
				*BeaconState,
				*BeaconStateMarshallable,
				*BlobProcessor,
				*BlobSidecar,
				*BlobSidecars,
				*BlobFactory,
				*BlockStore,
				*Context,
				*Deposit,
				*DepositStore,
				*Eth1Data,
				*EthClient,
				*EngineClient,
				*ExecutionEngine,
				*ExecutionPayload,
				*ExecutionPayloadHeader,
				*Fork,
				*ForkData,
				*Genesis,
				*IndexDB,
				*KVStore,
				KZGBlobProofVerifier,
				LegacyKey,
				*LocalBuilder,
				*Logger,
				*ABCIMiddleware,
				node,
				*PayloadAttributes,
				PayloadID,
				*SlashingInfo,
				*SlotData,
				*StateProcessor,
				*StorageBackend,
				*Validator,
				Validators,
				*ValidatorUpdate,
				ValidatorUpdates,
				*Withdrawal,
				Withdrawals,
				WithdrawalCredentials,
				*BlockBroker,
				*BlockStoreService,
				*ChainService,
				*DAService,
				*DBManager,
				*DepositService,
				*GenesisBroker,
				*NodeAPIServer,
				*ReportingService,
				*SidecarsBroker,
				*SlotBroker,
				*ValidatorService,
				*ValidatorUpdateBroker,
				*BlockEvent,
				*StatusEvent,
				DAPruner,
				BlockPruner,
				DepositPruner,
				*BeaconAPIHandler,
				*BuilderAPIHandler,
				*ConfigAPIHandler,
				*DebugAPIHandler,
				*EventsAPIHandler,
				*NodeAPIHandler,
				*ProofAPIHandler,
				NodeAPIContext,
				*NodeAPIEngine,
				*NodeAPIBackend,
				*NodeAPIBackend,
			](),
		),
	)

	// Build the root command using the builder
	cb := clibuilder.New(
		// Set the Name to the Default.
		clibuilder.WithName[
			node, *ExecutionPayload, LegacyKey, *Logger,
		](
			"BeaconKit",
		),
		// Set the Description to the Default.
		clibuilder.WithDescription[
			node, *ExecutionPayload, LegacyKey, *Logger,
		](
			"A basic beacon node, usable most standard networks.",
		),
		// Set the Runtime Components to the Default.
		clibuilder.WithComponents[
			node, *ExecutionPayload, LegacyKey, *Logger,
		](
			append(
				clicomponents.DefaultClientComponents(),
				// TODO: remove these, and eventually pull cfg and chainspec
				// from built node
				nodecomponents.ProvideConfig,
				nodecomponents.ProvideChainSpec,
			),
		),
		clibuilder.SupplyModuleDeps[
			node, *ExecutionPayload, LegacyKey, *Logger,
		](
			[]any{
				&ABCIMiddleware{},
				&runtime.App{},
				&StorageBackend{},
			},
		),
		// Set the Run Handler to the Default.
		clibuilder.WithRunHandler[
			node, *ExecutionPayload, LegacyKey, *Logger,
		](
			server.InterceptConfigsPreRunHandler,
		),
		// Set the NodeBuilderFunc to the NodeBuilder Build.
		clibuilder.WithNodeBuilderFunc[
			node, *ExecutionPayload, LegacyKey, *Logger,
		](nb.Build),
	)

	cmd, err := cb.Build()
	if err != nil {
		return err
	}

	// eventually we want to decouple from cosmos cli, and just pass in a built
	// Node and Cmd to a runner

	// for now, running the cmd will start the node
	return cmd.Run(clicomponents.DefaultNodeHome)
}

// main is the entry point.
func main() {
	if err := run(); err != nil {
		//nolint:sloglint // todo fix.
		slog.Error("startup failure", "error", err)
		os.Exit(1)
	}
}
