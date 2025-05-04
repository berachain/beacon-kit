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
	"fmt"
	"sync"
	"time"

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
	// errChan is used to collect errors from background goroutines
	errChan chan error
	// goroutineCtx is the context used for all background goroutines
	goroutineCtx context.Context
	// goroutineCancel is the cancel function for goroutineCtx
	goroutineCancel context.CancelFunc
}

// NewService creates a new validator service.
func NewService(
	storageBackend StorageBackend,
	blobProcessor BlobProcessor,
	depositContract deposit.Contract,
	eth1FollowDistance math.U64,
	logger log.Logger,
	chainSpec ServiceChainSpec,
	executionEngine ExecutionEngine,
	localBuilder LocalBuilder,
	stateProcessor StateProcessor,
	telemetrySink TelemetrySink,
	optimisticPayloadBuilds bool,
) *Service {
	// Using a modest buffer size of 32 for the error channel to handle potential
	// error bursts without blocking goroutines, while keeping memory usage reasonable
	const errorChannelBufferSize = 32

	return &Service{
		storageBackend:          storageBackend,
		blobProcessor:           blobProcessor,
		depositContract:         depositContract,
		eth1FollowDistance:      eth1FollowDistance,
		failedBlocks:            make(map[math.U64]struct{}),
		logger:                  logger,
		chainSpec:               chainSpec,
		executionEngine:         executionEngine,
		localBuilder:            localBuilder,
		stateProcessor:          stateProcessor,
		metrics:                 newChainMetrics(telemetrySink),
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
		errChan:                 make(chan error, errorChannelBufferSize),
	}
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return "blockchain"
}

// Start starts the blockchain service.
func (s *Service) Start(ctx context.Context) error {
	s.goroutineCtx, s.goroutineCancel = context.WithCancel(ctx)

	// Start monitoring goroutine errors
	go s.monitorGoroutineErrors()

	// Start the deposit catchup fetcher to retry processing deposits from blocks that failed previously.
	// This is a critical component for ensuring all deposits are properly captured and processed.
	go func() {
		// Use a separate function to handle any panics
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("depositCatchupFetcher panicked", "panic", r)
				s.errChan <- fmt.Errorf("depositCatchupFetcher panic: %v", r)
			}
		}()

		s.depositCatchupFetcher(s.goroutineCtx)
	}()

	s.logger.Info("Blockchain service started")
	return nil
}

// Stop stops the blockchain service and closes the deposit store.
func (s *Service) Stop() error {
	s.logger.Info("Stopping blockchain service")

	// First cancel context to notify all goroutines they should exit
	if s.goroutineCancel != nil {
		s.goroutineCancel()
	}
	
	// Add a wait time to ensure goroutines have a chance to observe the context cancellation
	// and stop sending to the error channel before we close it
	safeShutdownTimeout := 2 * time.Second
	timer := time.NewTimer(safeShutdownTimeout)
	defer timer.Stop()
	
	// Only proceed with cleanup after timeout or when the monitoring goroutine exits
	<-timer.C
	
	// Now it's safer to close the error channel as goroutines should have exited
	if s.errChan != nil {
		close(s.errChan)
	}

	// Close the deposit store
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

// monitorGoroutineErrors monitors the error channel and logs any errors.
// If a critical error occurs, it can trigger a service shutdown.
func (s *Service) monitorGoroutineErrors() {
	for {
		select {
		case <-s.goroutineCtx.Done():
			return
		case err := <-s.errChan:
			if err != nil {
				s.logger.Error("Background goroutine error", "error", err)
				// Optionally, implement logic to determine if this is a critical error
				// that should trigger a service shutdown
			}
		}
	}
}
