package sync

import (
	"context"
	"sync/atomic"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
)

type Service[
	SubscriptionT interface {
		Unsubscribe()
	},
] struct {
	// CL
	clSyncFeed                  EventFeed[*feed.Event[bool], SubscriptionT]
	clSyncCount                 atomic.Uint64
	clSyncStatusUpdateThreshold uint64
	clSyncStatus                uint8

	logger log.Logger[any]
}

func New[SubscriptionT interface {
	Unsubscribe()
}](
	clSyncFeed EventFeed[*feed.Event[bool], SubscriptionT],
	logger log.Logger[any],
) *Service[SubscriptionT] {
	return &Service[SubscriptionT]{
		clSyncFeed:  clSyncFeed,
		clSyncCount: atomic.Uint64{},
		//nolint:mnd // todo configurable.
		clSyncStatusUpdateThreshold: 10,
		logger:                      logger,
	}
}

// Name returns the name of the service.
func (s *Service[SubscriptionT]) Name() string {
	return "sync"
}

func (s *Service[SubscriptionT]) Status() error {
	return nil
}

func (s *Service[SubscriptionT]) Start(
	ctx context.Context,
) error {

	ch := make(chan *feed.Event[bool])
	sub := s.clSyncFeed.Subscribe(ch)
	defer sub.Unsubscribe()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-ch:
				if event.Is(events.CLSyncUpdate) {
					s.handleClSyncUpdateEvent(event)
				} else {
					s.logger.Warn("unexpected event", "event", event)
				}
			}
		}
	}()
	return nil
}

// handleClSyncUpdateEvent processes a CL sync update event.
func (s *Service[SubscriptionT]) handleClSyncUpdateEvent(event *feed.Event[bool]) {

	// 1. If we are not sync'd and the event is true, increment the count.
	// 2. If the event is false, reset the count and sync status.
	// 3. Otherwise ignore everything.
	if s.clSyncStatus == 0 && event.Data() {
		s.clSyncCount.Add(1)
	} else if !event.Data() {
		s.clSyncCount.Store(0)
		s.clSyncStatus = 0
	} else {
		return
	}

	// Otherwise update the CL status.
	if s.clSyncCount.Load() >= s.clSyncStatusUpdateThreshold {
		s.clSyncStatus = 1
	}
}
