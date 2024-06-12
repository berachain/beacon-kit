package cl

import (
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
			s.syncStatus = sync.CLStatusSynced
		}

	// 2. If we see an event that tells us we are not synced to head
	// immediately reset the counter and mark the CL as `NOT_SYNCED`.
	case !event.Data():
		s.syncCount.Store(0)
		s.syncStatus = sync.CLStatusNotSynced
	}
}
