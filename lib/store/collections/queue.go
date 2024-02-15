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

	sdk "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

type Queue[V any] struct {
	container sdk.Map[uint64, V]
	headSeq   sdk.Sequence // inclusive
	tailSeq   sdk.Sequence // exclusive
}

const (
	_ = iota
	headSeqOffset
	tailSeqOffset
)

// NewQueue creates a new queue with the provided prefix and name.
func NewQueue[V any](
	schema *sdk.SchemaBuilder,
	startPrefixID int, name string,
	valueCodec codec.ValueCodec[V]) Queue[V] {
	return Queue[V]{
		container: sdk.NewMap[uint64, V](
			schema, sdk.NewPrefix(startPrefixID),
			name+"_queue", sdk.Uint64Key, valueCodec),
		headSeq: sdk.NewSequence(schema, sdk.NewPrefix(startPrefixID+headSeqOffset), name+"_head"),
		tailSeq: sdk.NewSequence(schema, sdk.NewPrefix(startPrefixID+tailSeqOffset), name+"_tail"),
	}
}

// Peek returns the top element of the queue, or ErrNotFound if the queue is empty.
func (q *Queue[V]) Peek(ctx context.Context) (V, error) {
	var v V
	headIdx, err := q.headSeq.Peek(ctx)
	if err != nil {
		return v, err
	}
	tailIdx, err := q.tailSeq.Peek(ctx)
	if err != nil {
		return v, err
	}
	if headIdx >= tailIdx {
		return v, sdk.ErrNotFound
	}
	v, err = q.container.Get(ctx, headIdx)
	if err != nil {
		return v, err
	}
	return v, nil
}

// Next returns the top element of the queue and removes it from the queue.
func (q *Queue[V]) Next(ctx context.Context) (V, error) {
	v, err := q.Peek(ctx)
	if err != nil {
		return v, err
	}
	headIdx, err := q.headSeq.Next(ctx)
	if err != nil {
		return v, err
	}
	err = q.container.Remove(ctx, headIdx)
	if err != nil {
		return v, err
	}
	return v, nil
}

// Push adds a new element to the queue.
func (q *Queue[V]) Push(ctx context.Context, value V) error {
	tailIdx, err := q.tailSeq.Peek(ctx)
	if err != nil {
		return err
	}
	err = q.container.Set(ctx, tailIdx, value)
	if err != nil {
		return err
	}
	_, err = q.tailSeq.Next(ctx)
	return err
}
