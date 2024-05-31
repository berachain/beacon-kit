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

package deposit

import (
	"context"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	sdkcodec "cosmossdk.io/collections/codec"
	"github.com/berachain/beacon-kit/mod/errors"
)

// Queue is a simple queue implementation that uses a map and two sequences.
// TODO: Check atomicity of write operations.
type Queue[DepositT Deposit] struct {
	// container is a map that holds the queue elements.
	container sdkcollections.Map[uint64, DepositT]
	// headSeq is a sequence that points to the head of the queue.
	headSeq sdkcollections.Sequence // inclusive
	// length is an item that holds the length of the queue.
	length sdkcollections.Item[uint64]
	// mu is a mutex that protects the queue.
	mu sync.RWMutex
}

// NewQueue creates a new queue with the provided prefix and name.
func NewQueue[DepositT Deposit](
	schema *sdkcollections.SchemaBuilder, name string,
	valueCodec sdkcodec.ValueCodec[DepositT],
) *Queue[DepositT] {
	var (
		queueName   = name + "_queue"
		headSeqName = name + "_head"
		lengthName  = name + "_length"
	)
	return &Queue[DepositT]{
		container: sdkcollections.NewMap(
			schema,
			sdkcollections.NewPrefix(queueName),
			queueName, sdkcollections.Uint64Key, valueCodec,
		),
		headSeq: sdkcollections.NewSequence(
			schema, sdkcollections.NewPrefix(headSeqName), headSeqName),
		length: sdkcollections.NewItem(
			schema,
			sdkcollections.NewPrefix(lengthName),
			lengthName,
			sdkcollections.Uint64Value,
		),
	}
}

func (q *Queue[DepositT]) Init(ctx context.Context) error {
	if err := q.headSeq.Set(ctx, 0); err != nil {
		return err
	}

	// The length starts at 0.
	return q.length.Set(ctx, 0)
}

// Peek wraps the peek method with a read lock.
func (q *Queue[DepositT]) Peek(ctx context.Context) (DepositT, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.UnsafePeek(ctx)
}

// UnsafePeek returns the top element of the queue without removing it.
// It is unsafe to call this method without acquiring the read lock.
func (q *Queue[DepositT]) UnsafePeek(
	ctx context.Context,
) (DepositT, error) {
	var (
		v       DepositT
		headIdx uint64
		length  uint64
		err     error
	)
	if headIdx, err = q.headSeq.Peek(ctx); err != nil {
		return v, err
	} else if length, err = q.len(ctx); err != nil || length == 0 {
		return v, err
	}
	return q.container.Get(ctx, headIdx)
}

// Pop returns the top element of the queue and removes it from the queue.
func (q *Queue[DepositT]) Pop(ctx context.Context) (DepositT, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var (
		v       DepositT
		headIdx uint64
		err     error
	)

	if v, err = q.UnsafePeek(ctx); err != nil {
		return v, err
	} else if headIdx, err = q.headSeq.Next(ctx); err != nil {
		return v, err
	}
	err = q.container.Remove(ctx, headIdx)
	if err != nil {
		return v, err
	}
	length, err := q.len(ctx)
	if err != nil {
		return v, err
	}
	err = q.length.Set(ctx, length-1)
	if err != nil {
		return v, err
	}
	return v, err
}

// PeekMulti returns the top n elements of the queue.
func (q *Queue[DepositT]) PeekMulti(
	ctx context.Context,
	n uint64,
) ([]DepositT, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	var err error
	headIdx, err := q.headSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	length, err := q.len(ctx)
	if err != nil {
		return nil, err
	}
	endIdx := min(headIdx+length, headIdx+n)
	ranger := new(sdkcollections.Range[uint64]).
		StartInclusive(headIdx).EndExclusive(endIdx)
	iter, err := q.container.Iterate(ctx, ranger)
	if err != nil {
		return nil, err
	}
	// iter.Values already closes the iterator.
	values, err := iter.Values()
	if err != nil {
		return nil, err
	}

	return values, nil
}

// PopMulti returns the top n elements of the queue and removes them from the
// queue.
func (q *Queue[DepositT]) PopMulti(
	ctx context.Context,
	n uint64,
) ([]DepositT, error) {
	if n == 0 {
		return nil, nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	var err error

	headIdx, err := q.headSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	length, err := q.len(ctx)
	if err != nil {
		return nil, err
	}
	endIdx := min(headIdx+length, headIdx+n)
	ranger := new(sdkcollections.Range[uint64]).
		StartInclusive(headIdx).EndExclusive(endIdx)
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
	err = q.length.Set(ctx, length-n)
	if err != nil {
		return nil, err
	}
	return values, nil
}

// Push adds a new element to the queue.
func (q *Queue[DepositT]) Push(
	ctx context.Context,
	value DepositT,
) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if err := q.container.Set(ctx, value.GetIndex(), value); err != nil {
		return err
	}
	length, err := q.len(ctx)
	if err != nil {
		return err
	}
	return q.length.Set(ctx, length+1)
}

// PushMulti adds multiple new elements to the queue.
func (q *Queue[DepositT]) PushMulti(
	ctx context.Context,
	values []DepositT,
) error {
	if len(values) == 0 {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, value := range values {
		if err := q.container.Set(ctx, value.GetIndex(), value); err != nil {
			return err
		}
	}
	length, err := q.len(ctx)
	if err != nil {
		return err
	}
	return q.length.Set(ctx, length+uint64(len(values)))
}

// Len returns the length of the queue. len assumes that the lock is already
// held.
func (q *Queue[DepositT]) len(ctx context.Context) (uint64, error) {
	length, err := q.length.Get(ctx)
	if errors.Is(err, sdkcollections.ErrNotFound) {
		return 0, nil
	}
	return length, err
}

// Len returns the length of the queue.
func (q *Queue[DepositT]) Len(ctx context.Context) (uint64, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	length, err := q.length.Get(ctx)
	if errors.Is(err, sdkcollections.ErrNotFound) {
		return 0, nil
	}
	return length, err
}

// Container returns the underlying map container of the queue.
func (q *Queue[DepositT]) Container() sdkcollections.Map[
	uint64, DepositT,
] {
	return q.container
}
