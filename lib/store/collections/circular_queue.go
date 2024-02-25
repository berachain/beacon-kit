package collections

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	sdkcodec "cosmossdk.io/collections/codec"
)

// CircularQueue is a simple queue implementation that uses a map and two sequences.
// TODO: Check atomicity of write operations.
type CircularQueue[V any] struct {
	queue *Queue[V]
	size  uint64
}

// NewCircularQueue creates a new queue with the provided prefix and name.
func NewCircularQueue[V any](
	schema *sdkcollections.SchemaBuilder, name string,
	valueCodec sdkcodec.ValueCodec[V],
	size uint64,
) *CircularQueue[V] {
	return &CircularQueue[V]{
		queue: NewQueue[V](
			schema, name, valueCodec,
		),
		size: size,
	}
}

// Peek wraps the peek method with a read lock.
func (q *CircularQueue[V]) Peek(ctx context.Context) (V, error) {
	return q.queue.Peek(ctx)
}

// Push pushes a new element to the queue and returns the element that was evicted
// by the circular property of the queue.
func (q *CircularQueue[V]) Push(ctx context.Context, item V) (V, error) {
	var v V
	if err := q.queue.Push(ctx, item); err != nil {
		return v, err
	}

	len, err := q.queue.Len(ctx)
	if err != nil {
		return v, err
	}

	if len > q.size {
		v, err = q.queue.Pop(ctx)
	}

	return v, err
}
