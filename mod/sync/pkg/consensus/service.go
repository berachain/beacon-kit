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

package consensus

import (
	"context"
	"sync/atomic"

	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/sync"
)

// defaultSyncStatusUpdateThreshold is the default threshold for updating
// the status of the CL.
const defaultSyncStatusUpdateThreshold = 5

// SyncService is a service that listens to sync events and updates the
// status of the CL.
type SyncService[
	SubscriptionT interface {
		Unsubscribe()
	},
] struct {
	// syncCh is the channel that the sync service listens to.
	syncCh chan *feed.Event[bool]
	// syncFeed is the event feed that the sync service listens to.
	// It notifies the sync service when the CL sync status is updated.
	syncFeed *event.FeedOf[*feed.Event[bool]]
	// syncCount is the number of consecutive sync events that claim
	// the CL is synced to head.
	syncCount atomic.Uint64
	// syncStatusUpdateThreshold is the threshold for updating the status
	// of the CL. The sync service requires this many consecutive sync
	// successful events to claim the CL is synced to head.
	syncStatusUpdateThreshold uint64
	// syncStatus is the current status of the CL.
	syncStatus sync.ClientStatus
	// logger is the logger used by the sync service.
	logger log.Logger[any]
}

// New creates a new sync service.
func NewSyncService[
	SubscriptionT interface {
		Unsubscribe()
	},
](
	syncFeed *event.FeedOf[*feed.Event[bool]],
	logger log.Logger[any],
) *SyncService[SubscriptionT] {
	return &SyncService[SubscriptionT]{
		syncCh:                    make(chan *feed.Event[bool], 1),
		syncFeed:                  syncFeed,
		syncCount:                 atomic.Uint64{},
		syncStatusUpdateThreshold: defaultSyncStatusUpdateThreshold,
		logger:                    logger,
	}
}

// Name returns the name of the service.
func (s *SyncService[SubscriptionT]) Name() string {
	return "cl-sync"
}

// Status returns the status of the service.
func (s *SyncService[SubscriptionT]) Status() error {
	return nil
}

// Start spawns any goroutines required by the service.
func (s *SyncService[SubscriptionT]) Start(
	ctx context.Context,
) error {
	go s.mainLoop(ctx)
	return nil
}

// mainLoop listens to sync events and updates the status of the CL.
func (s *SyncService[SubscriptionT]) mainLoop(ctx context.Context) {
	sub := s.syncFeed.Subscribe(s.syncCh)
	defer sub.Unsubscribe()
	for {
		select {
		case event := <-s.syncCh:
			if event.Is(events.CLSyncUpdate) {
				// Update the status of the CL.
				s.handleCLSyncUpdateEvent(event)
			}
		case <-ctx.Done():
			return
		}
	}
}
