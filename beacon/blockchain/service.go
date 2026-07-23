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
	"sync/atomic"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
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
	// blobRequester distributes and retrieves sidecars over the blob reactor
	// p2p channel once blob consensus is enabled.
	blobRequester BlobRequester
	// blobReconstructor rebuilds sidecars from the local execution client's
	// blob pool.
	blobReconstructor BlobReconstructor
	// blobFetcher fetches sidecars of finalized blocks in the background.
	blobFetcher BlobFetcher
	// depositContract is the contract interface for interacting with the
	// deposit contract.
	depositContract deposit.Contract
	// eth1FollowDistance is the follow distance for Ethereum 1.0 blocks.
	eth1FollowDistance math.U64
	// logger is used for logging messages in the service.
	logger log.Logger
	// chainSpec holds the chain specifications.
	chainSpec ServiceChainSpec
	// executionEngine is the execution engine responsible for processing
	// execution payloads.
	executionEngine ExecutionEngine
	// localBuilder is a local builder for constructing new beacon states.
	localBuilder LocalBuilder
	// stateProcessor is the state processor for beacon blocks and states.
	stateProcessor StateProcessor
	// metrics is the metrics for the service.
	metrics *chainMetrics
	// forceStartupSyncOnce is used to force a sync of the startup head.
	forceStartupSyncOnce *sync.Once
	// latestFcuReq holds a copy of the latest FCU sent to the execution layer.
	// It helps avoid resending the same FCU data (and spares a network call)
	// in case optimistic block building is active
	latestFcuReq atomic.Pointer[engineprimitives.ForkchoiceStateV1]
}

// NewService creates a new validator service.
func NewService(
	storageBackend StorageBackend,
	blobProcessor BlobProcessor,
	blobRequester BlobRequester,
	blobReconstructor BlobReconstructor,
	blobFetcher BlobFetcher,
	depositContract deposit.Contract,
	logger log.Logger,
	chainSpec ServiceChainSpec,
	executionEngine ExecutionEngine,
	localBuilder LocalBuilder,
	stateProcessor StateProcessor,
	telemetrySink TelemetrySink,
) *Service {
	return &Service{
		storageBackend:       storageBackend,
		blobProcessor:        blobProcessor,
		blobRequester:        blobRequester,
		blobReconstructor:    blobReconstructor,
		blobFetcher:          blobFetcher,
		depositContract:      depositContract,
		eth1FollowDistance:   math.U64(chainSpec.Eth1FollowDistance()),
		logger:               logger,
		chainSpec:            chainSpec,
		executionEngine:      executionEngine,
		localBuilder:         localBuilder,
		stateProcessor:       stateProcessor,
		metrics:              newChainMetrics(telemetrySink),
		forceStartupSyncOnce: new(sync.Once),
	}
}

// PendingBlobRequests reports queued background blob fetches; a node with
// pending in-window fetches must not report itself as synced.
func (s *Service) PendingBlobRequests() int {
	return s.blobFetcher.PendingRequests()
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return "blockchain"
}

// Start starts the blockchain service.
func (s *Service) Start(ctx context.Context) error {
	// Start the background blob fetcher.
	s.blobFetcher.Start(ctx)
	return nil
}

// Stop stops the blockchain service and closes the deposit store.
func (s *Service) Stop() error {
	s.logger.Info("Stopping blockchain service")

	s.blobFetcher.Stop()

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

// PruneOrphanedBlobs removes any orphaned blob sidecars that may exist from incomplete block finalization.
func (s *Service) PruneOrphanedBlobs(lastBlockHeight int64) error {
	orphanedSlot := math.Slot(lastBlockHeight + 1) // #nosec G115

	// Once blob consensus is enabled, replayed blocks no longer carry sidecars, so these may be the only local
	// copy and replay depends on them. Keep them. FinalizeSidecars binds them to the replayed block's header
	// before reuse, so a set from a different block at this height is replaced through the normal fetch path.
	if s.chainSpec.IsBlobConsensusEnabled(lastBlockHeight + 1) {
		return nil
	}

	// Check if any blob sidecars exist at the potentially orphaned slot
	sidecars, err := s.storageBackend.AvailabilityStore().GetBlobSidecars(orphanedSlot)
	if err != nil {
		return fmt.Errorf("failed to read blob sidecars at slot %d: %w", orphanedSlot, err)
	}

	// If no sidecars exist at this slot, nothing to clean up
	if len(sidecars) == 0 {
		return nil
	}

	// Sidecars exist at this slot - they are orphaned, so delete them
	s.logger.Warn("Found orphaned blob sidecars from incomplete block finalization, removing",
		"slot", orphanedSlot.Base10(),
		"num_sidecars", len(sidecars),
	)

	err = s.storageBackend.AvailabilityStore().DeleteBlobSidecars(orphanedSlot)
	if err != nil {
		return fmt.Errorf("failed to delete orphaned sidecars at slot %d: %w", orphanedSlot, err)
	}

	s.logger.Info("Successfully removed orphaned blob sidecars", "slot", orphanedSlot.Base10())

	return nil
}
