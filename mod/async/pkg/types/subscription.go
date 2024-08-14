package types

import "context"

type Subscription[T any] chan T

func NewSubscription[T any]() Subscription[T] {
	return make(chan T)
}

func (s Subscription[T]) Await(ctx context.Context) (T, error) {
	select {
	case event := <-s:
		return event, nil
	case <-ctx.Done():
		return *new(T), ctx.Err()
	}
}
