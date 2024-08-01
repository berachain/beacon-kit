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
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	chainservice "github.com/berachain/beacon-kit/mod/node-core/pkg/services/blockchain"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

type ChainServiceInput struct {
	depinject.In
	ChainProcessor    ChainProcessorI
	ChainEventHandler *ChainEventHandler
}

func ProvideChainService(
	in ChainServiceInput,
) *ChainService {
	return service.NewService[
		*ChainEventHandler,
		ChainProcessorI,
	](
		in.ChainEventHandler,
		in.ChainProcessor,
	)
}

type ChainEventHandlerInput struct {
	depinject.In
	BlockBroker           *BlockBroker
	ValidatorUpdateBroker *ValidatorUpdateBroker
	GenesisBroker         *GenesisBroker
	Logger                log.AdvancedLogger[any, sdklog.Logger]
}

func ProvideChainEventHandler(
	in ChainEventHandlerInput,
) *ChainEventHandler {
	return chainservice.NewEventHandler[
		*BeaconBlock,
		*BeaconBlockBody,
		*Deposit,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Genesis,
	](
		in.BlockBroker,
		in.GenesisBroker,
		in.ValidatorUpdateBroker,
		in.Logger,
	)
}

// ChainProcessorInput is the input for the chain service provider.
type ChainProcessorInput struct {
	depinject.In
	BlockBroker           *BlockBroker
	ChainSpec             common.ChainSpec
	Cfg                   *config.Config
	DepositService        *DepositService
	EngineClient          *EngineClient
	ExecutionEngine       *ExecutionEngine
	GenesisBrocker        *GenesisBroker
	LocalBuilder          *LocalBuilder
	Logger                log.AdvancedLogger[any, sdklog.Logger]
	Signer                crypto.BLSSigner
	StateProcessor        *StateProcessor
	StorageBackend        *StorageBackend
	TelemetrySink         *metrics.TelemetrySink
	ValidatorUpdateBroker *ValidatorUpdateBroker
}

// ProvideChainProcessor is a depinject provider for the blockchain service.
func ProvideChainProcessor(
	in ChainProcessorInput,
) *ChainProcessor {
	return blockchain.NewProcessor[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Genesis,
		*PayloadAttributes,
		*Withdrawal,
	](
		in.StorageBackend,
		in.Logger.With("service", "blockchain"),
		in.ChainSpec,
		in.ExecutionEngine,
		in.LocalBuilder,
		in.StateProcessor,
		in.TelemetrySink,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		in.Cfg.Validator.EnableOptimisticPayloadBuilds,
	)
}

// ChainServiceComponents returns a slice of the main components needed to
// depinject the chain service.
func ChainServiceComponents() []any {
	return []any{
		ProvideChainService,
		ProvideChainEventHandler,
		ProvideChainProcessor,
	}
}
