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
	"sync"

	sdkcollections "cosmossdk.io/collections"
	sdkcodec "cosmossdk.io/collections/codec"
)

// Queue is a simple queue implementation that uses a map and two sequences.
// TODO: Check atomicity of write operations.
type Queue[V any] struct {
	// container is a map that holds the queue elements.
	container sdkcollections.Map[uint64, V]
	// headSeq is a sequence that points to the head of the queue.
	headSeq sdkcollections.Sequence // inclusive
	// tailSeq is a sequence that points to the tail of the queue.
	tailSeq sdkcollections.Sequence // exclusive
	// mu is a mutex that protects the queue.
	mu sync.RWMutex
}

// NewQueue creates a new queue with the provided prefix and name.
func NewQueue[V any](
	schema *sdkcollections.SchemaBuilder, name string,
	valueCodec sdkcodec.ValueCodec[V],
) *Queue[V] {
	var (
		queueName   = name + "_queue"
		headSeqName = name + "_head"
		tailSeqName = name + "_tail"
	)
	return &Queue[V]{
		container: sdkcollections.NewMap[uint64, V](
			schema,
			sdkcollections.NewPrefix(queueName),
			queueName, sdkcollections.Uint64Key, valueCodec,
		),
		headSeq: sdkcollections.NewSequence(schema, sdkcollections.NewPrefix(headSeqName), headSeqName),
		tailSeq: sdkcollections.NewSequence(schema, sdkcollections.NewPrefix(tailSeqName), tailSeqName),
	}
}

// Peek wraps the peek method with a read lock.
func (q *Queue[V]) Peek(ctx context.Context) (V, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.UnsafePeek(ctx)
}

// UnsafePeek returns the top element of the queue without removing it.
// It is unsafe to call this method without acquiring the read lock.
func (q *Queue[V]) UnsafePeek(ctx context.Context) (V, error) {
	var (
		v                V
		headIdx, tailIdx uint64
		err              error
	)
	if headIdx, err = q.headSeq.Peek(ctx); err != nil {
		return v, err
	} else if tailIdx, err = q.tailSeq.Peek(ctx); err != nil {
		return v, err
	} else if headIdx >= tailIdx {
		return v, sdkcollections.ErrNotFound
	}
	return q.container.Get(ctx, headIdx)
}

// Pop returns the top element of the queue and removes it from the queue.
func (q *Queue[V]) Pop(ctx context.Context) (V, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var (
		v       V
		headIdx uint64
		err     error
	)

	if v, err = q.UnsafePeek(ctx); err != nil {
		return v, err
	} else if headIdx, err = q.headSeq.Next(ctx); err != nil {
		return v, err
	}
	err = q.container.Remove(ctx, headIdx)
	return v, err
}

// PopMulti returns the top n elements of the queue and removes them from the queue.
func (q *Queue[V]) PopMulti(ctx context.Context, n uint64) ([]V, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var err error

	headIdx, err := q.headSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	tailIdx, err := q.tailSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	endIdx := min(tailIdx, headIdx+n)
	ranger := new(sdkcollections.Range[uint64]).StartInclusive(headIdx).EndExclusive(endIdx)
	iter, err := q.container.Iterate(ctx, ranger)
	if err != nil {
		return nil, err
	}
	// iter.Values already closes the iterator.
	values, err := iter.Values()
	if err != nil {
		return nil, err
	}

	// Clear the range (in batch) after the iteration is done.
	err = q.container.Clear(ctx, ranger)
	if err != nil {
		return nil, err
	}
	err = q.headSeq.Set(ctx, endIdx)
	if err != nil {
		return nil, err
	}
	return values, nil
}

// Push adds a new element to the queue.
func (q *Queue[V]) Push(ctx context.Context, value V) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	var (
		tailIdx uint64
		err     error
	)

	// Get the current tail index.
	if tailIdx, err = q.tailSeq.Peek(ctx); err != nil {
		return err
	} else if err = q.container.Set(ctx, tailIdx, value); err != nil {
		return err
	}

	// If the push is successful, increment the tail sequence.
	_, err = q.tailSeq.Next(ctx)
	return err
}

// PushMulti adds multiple new elements to the queue.
func (q *Queue[V]) PushMulti(ctx context.Context, values []V) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	var (
		tailIdx uint64
		err     error
	)

	// Get the current tail index.
	if tailIdx, err = q.tailSeq.Peek(ctx); err != nil {
		return err
	}
	for _, value := range values {
		if err = q.container.Set(ctx, tailIdx, value); err != nil {
			return err
		}
		tailIdx++
	}

	// If the push is successful, set the tail sequence to the new index.
	return q.tailSeq.Set(ctx, tailIdx)
}

// Len returns the length of the queue.
func (q *Queue[V]) Len(ctx context.Context) (uint64, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	var (
		headIdx, tailIdx uint64
		err              error
	)

	if headIdx, err = q.headSeq.Peek(ctx); err != nil {
		return 0, err
	} else if tailIdx, err = q.tailSeq.Peek(ctx); err != nil {
		return 0, err
	}
	return tailIdx - headIdx, nil
}

// Container returns the underlying map container of the queue.
func (q *Queue[V]) Container() sdkcollections.Map[uint64, V] {
	return q.container
}
