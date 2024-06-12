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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/sync"
)

// handleCLSyncUpdateEvent processes a CL sync update event.
func (s *SyncService[SubscriptionT]) handleCLSyncUpdateEvent(
	event *feed.Event[bool],
) {
	switch {
	// 1. If we are not synced to head, and we have
	// a synced event, increment the sync count.
	case s.syncStatus == sync.CLStatusNotSynced && event.Data():
		s.syncCount.Add(1)
		// If the sync count is greater than or equal to the
		// threshold, mark the CL as `SYNCED`.
		if s.syncCount.Load() >= s.syncStatusUpdateThreshold {
			s.logger.Info("marking consensus client as synced 🎉")
			s.syncStatus = sync.CLStatusSynced
			s.syncFeed.Send(
				feed.NewEvent(context.Background(), events.CLSyncStatus, true),
			)
		}

	// 2. If we see an event that tells us we are not synced to head
	// immediately reset the counter and mark the CL as `NOT_SYNCED`.
	case !event.Data():
		s.syncCount.Store(0)
		s.syncStatus = sync.CLStatusNotSynced
		s.syncFeed.Send(
			feed.NewEvent(context.Background(), events.CLSyncStatus, false),
		)
	}
}
