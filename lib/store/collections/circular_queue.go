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
	"errors"

	"cosmossdk.io/collections"
	sdkcodec "cosmossdk.io/collections/codec"
)

// CircularQueue is a simple queue implementation that uses a map and two
// sequences.
type CircularQueue[V any] struct {
	// container is a map that holds the queue elements.
	container collections.Map[uint64, V]

	// size is the size of the queue.
	size uint64
}

// NewCircularQueue creates a new queue with the provided prefix and name.
func NewCircularQueue[V any](
	schema *collections.SchemaBuilder,
	name string,
	valueCodec sdkcodec.ValueCodec[V],
	size uint64,
) *CircularQueue[V] {
	return &CircularQueue[V]{
		container: collections.NewMap(
			schema,
			collections.NewPrefix(name),
			name,
			collections.Uint64Key,
			valueCodec,
		),
		size: size,
	}
}

// Push pushes a new element to the queue.
func (q *CircularQueue[V]) Push(
	ctx context.Context,
	index uint64,
	item V,
) error {
	return q.container.Set(ctx, index%q.size, item)
}

// Peek returns the element at the given index.
func (q *CircularQueue[V]) Peek(ctx context.Context, index uint64) (V, error) {
	v, err := q.container.Get(ctx, index%q.size)
	if errors.Is(collections.ErrNotFound, err) {
		return v, nil
	} else if err != nil {
		return v, err
	}
	return v, nil
}
