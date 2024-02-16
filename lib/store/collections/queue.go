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

	sdk "cosmossdk.io/collections"
	sdkcodec "cosmossdk.io/collections/codec"
)

// Queue is a simple queue implementation that uses a map and two sequences.
type Queue[V any] struct {
	// container is a map that holds the queue elements.
	container sdk.Map[uint64, V]
	// headSeq is a sequence that points to the head of the queue.
	headSeq sdk.Sequence // inclusive
	// tailSeq is a sequence that points to the tail of the queue.
	tailSeq sdk.Sequence // exclusive
	// mu is a mutex that protects the queue.
	mu sync.RWMutex
}

// NewQueue creates a new queue with the provided prefix and name.
func NewQueue[V any](
	schema *sdk.SchemaBuilder, name string,
	valueCodec sdkcodec.ValueCodec[V],
) Queue[V] {
	var (
		queueName   = name + "_queue"
		headSeqName = name + "_head"
		tailSeqName = name + "_tail"
	)
	return Queue[V]{
		container: sdk.NewMap[uint64, V](
			schema,
			sdk.NewPrefix(queueName),
			queueName, sdk.Uint64Key, valueCodec,
		),
		headSeq: sdk.NewSequence(schema, sdk.NewPrefix(headSeqName), headSeqName),
		tailSeq: sdk.NewSequence(schema, sdk.NewPrefix(tailSeqName), tailSeqName),
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
		return v, sdk.ErrNotFound
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

// Push adds a new element to the queue.
func (q *Queue[V]) Push(ctx context.Context, value V) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	var (
		tailIdx uint64
		err     error
	)

	// If the queue is empty, set the head sequence to 0.
	if tailIdx, err = q.tailSeq.Peek(ctx); err != nil {
		return err
	} else if err = q.container.Set(ctx, tailIdx, value); err != nil {
		return err
	}

	// If the pop is successful, increment the tail sequence.
	_, err = q.tailSeq.Next(ctx)
	return err
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
