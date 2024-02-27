// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package collections

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	sdkcodec "cosmossdk.io/collections/codec"
)

// CircularQueue is a simple queue implementation that uses a map and two
// sequences.
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

// Peek wraps the peek method.
func (q *CircularQueue[V]) Peek(ctx context.Context) (V, error) {
	return q.queue.Peek(ctx)
}

// Push pushes a new element to the queue and returns the element that was
// evicted
// by the circular property of the queue.
func (q *CircularQueue[V]) Push(ctx context.Context, item V) (V, error) {
	var v V
	if err := q.queue.Push(ctx, item); err != nil {
		return v, err
	}

	if length, err := q.queue.Len(ctx); err != nil {
		return v, err
	} else if length > q.size {
		v, err = q.queue.Pop(ctx)
		if err != nil {
			return v, err
		}
	}

	return v, nil
}
