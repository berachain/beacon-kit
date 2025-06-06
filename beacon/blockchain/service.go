// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package blockchain

import (
	"context"
	"sync"

	"github.com/berachain/beacon-kit/execution/deposit"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Service is the blockchain service.
type Service struct {
	// storageBackend represents the backend storage for not state-enforced data.
	storageBackend StorageBackend
	// blobProcessor is used for processing sidecars.
	blobProcessor BlobProcessor
	// depositContract is the contract interface for interacting with the
	// deposit contract.
	depositContract deposit.Contract
	// eth1FollowDistance is the follow distance for Ethereum 1.0 blocks.
	eth1FollowDistance math.U64
	// failedBlocksMu protects failedBlocks for concurrent access.
	failedBlocksMu sync.RWMutex
	// failedBlocks is a map of blocks that failed to be processed
	// and should be retried.
	failedBlocks map[math.U64]struct{}
	// logger is used for logging messages in the service.
	logger log.Logger
	// chainSpec holds the chain specifications.
	chainSpec ServiceChainSpec
	// executionEngine is the execution engine responsible for processing
	//
	// execution payloads.
	executionEngine ExecutionEngine
	// localBuilder is a local builder for constructing new beacon states.
	localBuilder LocalBuilder
	// stateProcessor is the state processor for beacon blocks and states.
	stateProcessor StateProcessor
	// metrics is the metrics for the service.
	metrics *chainMetrics
	// optimisticPayloadBuilds is a flag used when the optimistic payload
	// builder is enabled.
	optimisticPayloadBuilds bool
	// forceStartupSyncOnce is used to force a sync of the startup head.
	forceStartupSyncOnce *sync.Once
}

// NewService creates a new validator service.
func NewService(
	storageBackend StorageBackend,
	blobProcessor BlobProcessor,
	depositContract deposit.Contract,
	logger log.Logger,
	chainSpec ServiceChainSpec,
	executionEngine ExecutionEngine,
	localBuilder LocalBuilder,
	stateProcessor StateProcessor,
	telemetrySink TelemetrySink,
	optimisticPayloadBuilds bool,
) *Service {
	return &Service{
		storageBackend:          storageBackend,
		blobProcessor:           blobProcessor,
		depositContract:         depositContract,
		eth1FollowDistance:      math.U64(chainSpec.Eth1FollowDistance()),
		failedBlocks:            make(map[math.Slot]struct{}),
		logger:                  logger,
		chainSpec:               chainSpec,
		executionEngine:         executionEngine,
		localBuilder:            localBuilder,
		stateProcessor:          stateProcessor,
		metrics:                 newChainMetrics(telemetrySink),
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
	}
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return "blockchain"
}

// Start starts the blockchain service.
func (s *Service) Start(ctx context.Context) error {
	// Catchup deposits for failed blocks. TODO: remove.
	go s.depositCatchupFetcher(ctx)

	return nil
}

// Stop stops the blockchain service and closes the deposit store.
func (s *Service) Stop() error {
	s.logger.Info("Stopping blockchain service")

	err := s.storageBackend.DepositStore().Close()
	if err != nil {
		s.logger.Error("failed to close deposit store", "err", err)
	}

	return nil
}

// StorageBackend returns the storage backend.
func (s *Service) StorageBackend() StorageBackend {
	return s.storageBackend
}
