package types

import "context"

type Subscription[T any] chan Event[T]

func (s Subscription[T]) Await(ctx context.Context) (Event[T], error) {
	select {
	case event := <-s:
		return event, nil
	case <-ctx.Done():
		return event[T]{}, ctx.Err()
	}
}
