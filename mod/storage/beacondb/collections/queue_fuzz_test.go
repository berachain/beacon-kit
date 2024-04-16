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

	sdk "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/mod/storage/beacondb/collections"
	"github.com/stretchr/testify/require"
)

func FuzzQueueSimple(f *testing.F) {
	f.Add(int64(1), int64(2), int64(3), int64(4))
	sk, ctx := deps()
	sb := sdk.NewSchemaBuilder(sk)
	q := collections.NewQueue[int64](sb, "queue", sdk.Int64Value)
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4 int64) {
		if n1 < 0 || n1 > n2 || n4 < 0 || n3 < n4 {
			t.Skip()
		}

		trackedItems := make([]int64, 0)
		for i := range n1 {
			require.NoError(t, q.Push(ctx, i))
			trackedItems = append(trackedItems, i)
		}

		l, err := q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, n1, int64(l))

		// n2 >= n1
		for i := range n2 {
			var item int64
			item, err = q.Pop(ctx)
			if i < n1 {
				require.NoError(t, err)
				require.Equal(t, trackedItems[0], item)
				trackedItems = trackedItems[1:]
			} else {
				require.Equal(t, sdk.ErrNotFound, err)
			}
		}

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(0), l)

		for i := range n3 {
			require.NoError(t, q.Push(ctx, n1+i))
			trackedItems = append(trackedItems, n1+i)
		}

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(n3), l)

		// n3 >= n4
		require.GreaterOrEqual(t, n3, n4)
		for range n4 {
			var item int64
			item, err = q.Pop(ctx)
			require.NoError(t, err)
			require.Equal(t, trackedItems[0], item)
			trackedItems = trackedItems[1:]
		}

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(n3-n4), l)

		for i := n4; i < n3; i++ {
			var item int64
			item, err = q.Pop(ctx)
			require.NoError(t, err)
			require.Equal(t, trackedItems[0], item)
			trackedItems = trackedItems[1:]
		}

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(0), l)
	})
}

func FuzzQueueMulti(f *testing.F) {
	f.Add(int64(1), int64(2), int64(3), int64(4))
	sk, ctx := deps()
	sb := sdk.NewSchemaBuilder(sk)
	q := collections.NewQueue[int64](sb, "queue", sdk.Int64Value)
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4 int64) {
		if n1 < 0 || n1 > n2 || n4 < 0 || n3 < n4 {
			t.Skip()
		}

		trackedItems := make([]int64, 0)
		for i := range n1 {
			trackedItems = append(trackedItems, i)
		}
		require.Len(t, trackedItems, int(n1))

		require.NoError(t, q.PushMulti(ctx, trackedItems))

		l, err := q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(n1), l)

		ml, err := mapLen[int64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(
			t, l, uint64(ml), "Queue length should match container length")

		// n2 >= n1
		poppedItems, err := q.PopMulti(ctx, uint64(n2))
		require.NoError(t, err)
		require.Len(t, poppedItems, int(n1))
		for i := range trackedItems {
			require.Equal(t, trackedItems[i], poppedItems[i])
		}

		trackedItems = trackedItems[:0]
		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(0), l)

		ml, err = mapLen[int64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(
			t, l, uint64(ml), "Queue length should match container length")

		for i := range n3 {
			trackedItems = append(trackedItems, n1+i)
		}
		require.Len(t, trackedItems, int(n3))

		require.NoError(t, q.PushMulti(ctx, trackedItems))

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Len(t, trackedItems, int(l))
		require.Equal(t, uint64(n3), l)

		ml, err = mapLen[int64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(
			t, l, uint64(ml), "Queue length should match container length")

		// n3 >= n4
		poppedItems, err = q.PopMulti(ctx, uint64(n4))
		require.NoError(t, err)
		require.Len(t, poppedItems, int(n4))
		for i := range n4 {
			require.Equal(t, trackedItems[i], poppedItems[i])
		}

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(n3-n4), l)

		ml, err = mapLen[int64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(
			t, l, uint64(ml), "Queue length should match container length")

		poppedItems, err = q.PopMulti(ctx, uint64(n3))
		require.NoError(t, err)
		require.Len(t, poppedItems, int(n3-n4))
		for i := range n3 - n4 {
			require.Equal(t, trackedItems[n4+i], poppedItems[i])
		}

		l, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), l)

		ml, err = mapLen[int64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(
			t, l, uint64(ml), "Queue length should match container length")
	})
}
