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

package collections_test

import (
	"testing"

	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/lib/store/collections"
	"github.com/stretchr/testify/require"
)

func TestCircularQueuePushPop(t *testing.T) {
	sk, ctx := deps()
	schema := sdkcollections.NewSchemaBuilder(sk)
	queue := collections.NewCircularQueue[uint64](
		schema,
		"testQueue",
		sdkcollections.Uint64Value,
		5,
	)

	// Push elements into the queue
	for i := uint64(0); i < 5; i++ {
		_, err := queue.Push(ctx, i)
		require.NoError(t, err)
	}

	// Push another element, which should cause the first element to be evicted
	// (circular behavior)
	evicted, err := queue.Push(ctx, 5)
	require.NoError(t, err)
	require.Equal(t, uint64(0), evicted)

	// Pop elements and verify order
	var item uint64
	for i := uint64(1); i <= 5; i++ {
		item, err = queue.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, i, item)
		_, err = queue.Push(
			ctx,
			item+5,
		) // Push next element to maintain circularity
		require.NoError(t, err)
	}
}

func TestCircularQueueOverflow(t *testing.T) {
	sk, ctx := deps()
	schema := sdkcollections.NewSchemaBuilder(sk)
	queue := collections.NewCircularQueue[uint64](
		schema,
		"testQueueOverflow",
		sdkcollections.Uint64Value,
		3,
	)

	// Push elements more than the queue size to test overflow
	for i := uint64(0); i < 10; i++ {
		_, err := queue.Push(ctx, i)
		require.NoError(t, err)
	}

	// Verify that the queue only contains the last 3 elements
	expected := []uint64{7, 8, 9}
	for _, e := range expected {
		item, err := queue.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, e, item)
		_, err = queue.Push(
			ctx,
			item,
		) // Push the peeked item back to simulate continuous operation
		require.NoError(t, err)
	}
}
