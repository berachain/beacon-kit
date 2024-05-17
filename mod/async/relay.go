package async

import (
	"context"
	"sync"
)

// Relay handles the subscribers and dispatches notifications.
type Relay[T any] struct {
	mu          sync.RWMutex
	n           uint32
	subscribers map[uint32]*Subscriber[T]
}

// NewRelay creates a new Relay.
func NewRelay[T any]() *Relay[T] {
	return &Relay[T]{
		subscribers: make(map[uint32]*Subscriber[T]),
	}
}

// AddSubscriber creates a new subscriber with a given channel capacity.
// TODO: allow the relay to track the actual number of real subscribers.
func (r *Relay[T]) AddSubscriber(capacity int) *Subscriber[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	subscriber := &Subscriber[T]{
		ch:    make(chan T, capacity),
		id:    r.n,
		relay: r,
	}
	r.subscribers[r.n] = subscriber
	r.n++
	return subscriber
}

// Notify sends a notification to all subscribers.
func (r *Relay[T]) Notify(v T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.subscribers {
		select {
		case client.ch <- v:
		default:
			// Log or handle the case where the channel is full
		}
	}
}

// NotifyCtx sends a notification to all subscribers until the context times out or is canceled.
func (r *Relay[T]) NotifyCtx(ctx context.Context, v T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.subscribers {
		select {
		case client.ch <- v:
		case <-ctx.Done():
			return
		}
	}
}

// Broadcast sends a notification to all subscribers in a non-blocking manner.
func (r *Relay[T]) Broadcast(v T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.subscribers {
		select {
		case client.ch <- v:
		default:
			// Log or handle the case where the channel is full
		}
	}
}

// Close closes the relay and all its subscribers.
func (r *Relay[T]) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, client := range r.subscribers {
		r.closeSubscriber(client)
	}
	r.subscribers = nil
}

// closeSubscriber closes the subscriber.
func (r *Relay[T]) closeSubscriber(l *Subscriber[T]) {
	l.once.Do(func() {
		close(l.ch)
		delete(r.subscribers, l.id)
	})
}
