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
	"github.com/stretchr/testify/require"

	"github.com/itsdevbear/bolaris/lib/store/collections"
)

func FuzzQueueSimple(f *testing.F) {
	f.Add(int64(1), int64(2), int64(3), int64(4))
	sk, ctx := deps()
	sb := sdk.NewSchemaBuilder(sk)
	q := collections.NewQueue[int64](sb, "queue", sdk.Int64Value)
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4 int64) {
		if n1 < 0 || n1 > n2 || n2 < n3 || n3 > n4 {
			t.Skip()
		}

		trackedItems := make([]int64, 0)
		for i := int64(0); i < n1; i++ {
			i := i // must capture loop var
			require.NoError(t, q.Push(ctx, i))
			trackedItems = append(trackedItems, i)
		}

		l, err := q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(len(trackedItems)), int64(l))
		require.Equal(t, n1, int64(l))

		for i := int64(0); i < n2 && len(trackedItems) > 0; i++ {
			var item int64
			item, err = q.Pop(ctx)
			require.NoError(t, err)
			require.Equal(t, trackedItems[0], item)
			trackedItems = trackedItems[1:]
		}

		for i := int64(0); i < n3; i++ {
			i := i // must capture loop var
			require.NoError(t, q.Push(ctx, n1+i))
			trackedItems = append(trackedItems, n1+i)
		}

		for i := int64(0); i < n4 && len(trackedItems) > 0; i++ {
			var item int64
			item, err = q.Pop(ctx)
			require.NoError(t, err)
			require.Equal(t, trackedItems[0], item)
			trackedItems = trackedItems[1:]
		}
	})
}
