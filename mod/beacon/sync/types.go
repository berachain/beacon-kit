package sync

// CLSyncUpdateEvent represents an interface for block events.
type CLSyncUpdateEvent interface {
	// Data returns a boolean indicating the event data.
	Data() bool
}

// EventFeed is a generic interface for sending events.
type EventFeed[
	CLSyncUpdateEventT CLSyncUpdateEvent,
	SubscriptionT interface {
		// Unsubscribe terminates the subscription.
		Unsubscribe()
	},
] interface {
	// Subscribe subscribes to the event feed and returns a subscription.
	Subscribe(chan<- CLSyncUpdateEventT) SubscriptionT
}
