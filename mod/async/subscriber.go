package async

import "sync"

// Subscriber represents a Relay subscriber.
type Subscriber[T any] struct {
	ch    chan T
	id    uint32
	relay *Relay[T]
	once  sync.Once
}

// Channel returns the Subscriber's channel.
func (s *Subscriber[T]) Channel() <-chan T {
	return s.ch
}

// Close closes the subscriber.
func (s *Subscriber[T]) Close() {
	s.once.Do(func() {
		s.relay.closeSubscriber(s)
	})
}
