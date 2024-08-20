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

package blockstore

import (
	"context"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

// Service is a Service that listens for blocks and stores them in a KVStore.
type Service[
	BeaconBlockT BeaconBlock,
	BlockStoreT BlockStore[BeaconBlockT],
] struct {
	// config is the configuration for the block service.
	config Config
	// logger is used for logging information and errors.
	logger log.Logger[any]
	// dispatcher is the dispatcher for the service.
	dispatcher asynctypes.EventDispatcher
	// store is the block store for the service.
	store BlockStoreT
	// subFinalizedBlkEvents is a channel for receiving finalized block events.
	subFinalizedBlkEvents chan async.Event[BeaconBlockT]
}

// NewService creates a new block service.
func NewService[
	BeaconBlockT BeaconBlock,
	BlockStoreT BlockStore[BeaconBlockT],
](
	config Config,
	logger log.Logger[any],
	dispatcher asynctypes.EventDispatcher,
	store BlockStoreT,
) *Service[BeaconBlockT, BlockStoreT] {
	return &Service[BeaconBlockT, BlockStoreT]{
		config:                config,
		logger:                logger,
		dispatcher:            dispatcher,
		store:                 store,
		subFinalizedBlkEvents: make(chan async.Event[BeaconBlockT]),
	}
}

// Name returns the name of the service.
func (s *Service[_, _]) Name() string {
	return "block-service"
}

// Start starts the block service.
func (s *Service[BeaconBlockT, _]) Start(ctx context.Context) error {
	if !s.config.Enabled {
		s.logger.Warn("block service is disabled, skipping storing blocks")
		return nil
	}

	// subscribe a channel to the finalized block events.
	if err := s.dispatcher.Subscribe(
		async.BeaconBlockFinalizedEvent, s.subFinalizedBlkEvents,
	); err != nil {
		s.logger.Error("failed to subscribe to block events", "error", err)
		return err
	}

	go s.eventLoop(ctx)
	return nil
}

func (s *Service[BeaconBlockT, BlockStoreT]) eventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.subFinalizedBlkEvents:
			s.onFinalizeBlock(event)
		}
	}
}

// onFinalizeBlock is triggered when a finalized block event is received.
// It stores the block in the KVStore.
func (s *Service[BeaconBlockT, _]) onFinalizeBlock(
	event async.Event[BeaconBlockT],
) {
	slot := event.Data().GetSlot()
	if err := s.store.Set(event.Data()); err != nil {
		s.logger.Error(
			"failed to store block", "slot", slot, "error", err,
		)
	}
}
