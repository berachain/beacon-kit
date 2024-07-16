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

package block

import (
	"context"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
)

// Service is a service that listens for blocks and stores them in a KVStore.
type Service[BeaconBlockT BeaconBlock] struct {
	// logger is used for logging information and errors.
	logger    log.Logger[any]
	blkBroker EventFeed[*asynctypes.Event[BeaconBlockT]]
	store     *block.KVStore[BeaconBlockT]
}

// NewService creates a new block service.
func NewService[BeaconBlockT BeaconBlock](
	logger log.Logger[any],
	blkBroker EventFeed[*asynctypes.Event[BeaconBlockT]],
	store *block.KVStore[BeaconBlockT],
) *Service[BeaconBlockT] {
	return &Service[BeaconBlockT]{
		logger:    logger,
		blkBroker: blkBroker,
		store:     store,
	}
}

// Name returns the name of the service.
func (s *Service[BeaconBlockT]) Name() string {
	return "block-service"
}

// Start starts the block service.
func (s *Service[BeaconBlockT]) Start(ctx context.Context) error {
	subBlkCh, err := s.blkBroker.Subscribe()
	if err != nil {
		s.logger.Error("failed to subscribe to block events", "error", err)
		return err
	}
	go s.listenAndStore(ctx, subBlkCh)
	return nil
}

// listenAndStore listens for blocks and stores them in the KVStore.
func (s *Service[BeaconBlockT]) listenAndStore(
	ctx context.Context,
	subBlkCh <-chan *asynctypes.Event[BeaconBlockT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-subBlkCh:
			if msg.Is(events.BeaconBlockFinalized) {
				slot := msg.Data().GetSlot()
				if err := s.store.Set(slot.Unwrap(), msg.Data()); err != nil {
					s.logger.Error("failed to store block", "error", err)
				}
			}
		}
	}
}
