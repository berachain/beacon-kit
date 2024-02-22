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
	"context"
	"testing"

	sdk "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	dba "cosmossdk.io/store/dbadapter"
	db "github.com/cosmos/cosmos-db"
	"github.com/itsdevbear/bolaris/lib/store/collections"
	"github.com/stretchr/testify/require"
)

func Test_Queue(t *testing.T) {
	t.Run("should return correct items and lengths", func(t *testing.T) {
		sk, ctx := deps()
		sb := sdk.NewSchemaBuilder(sk)
		q := collections.NewQueue[uint64](sb, "queue", sdk.Uint64Value)

		// Test initial length of the queue
		qlen, err := q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), qlen, "Queue should be empty initially")

		_, err = q.Peek(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		_, err = q.Pop(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		err = q.Push(ctx, 1)
		require.NoError(t, err)

		// Test length after first push
		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), qlen, "Queue should have 1 item after first push")

		err = q.Push(ctx, 2)
		require.NoError(t, err)

		// Test length after second push
		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), qlen, "Queue should have 2 items after second push")

		v, err := q.Pop(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), v)

		// Test length after first pop
		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), qlen, "Queue should have 1 item after first pop")

		v, err = q.Pop(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), v)

		// Test length after clearing the queue
		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), qlen, "Queue should be empty after clearing")

		// Attempt to peek at the top value of the queue, expecting an error
		// since the queue should now be empty
		_, err = q.Peek(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		// Attempt to pop an item from the queue, expecting an error since the
		// queue is empty
		_, err = q.Pop(ctx)
		require.Equal(t, sdk.ErrNotFound, err)
	})
}

func Test_PopMulti(t *testing.T) {
	t.Run("should return correct items and lengths", func(t *testing.T) {
		sk, ctx := deps()
		sb := sdk.NewSchemaBuilder(sk)
		q := collections.NewQueue[uint64](sb, "queue", sdk.Uint64Value)

		err := q.Push(ctx, 1)
		require.NoError(t, err)

		err = q.Push(ctx, 2)
		require.NoError(t, err)

		err = q.Push(ctx, 3)
		require.NoError(t, err)

		// Test length after pushes
		qlen, err := q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), qlen,
			"Queue should have 3 items after 3 pushes")

		mlen, err := mapLen[uint64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(t, qlen, uint64(mlen),
			"Queue length should match container length")

		// Pop 2 items from the queue
		items, err := q.PopMulti(ctx, 2)
		require.NoError(t, err)
		require.Equal(t, []uint64{1, 2}, items)

		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), qlen,
			"Queue should have 1 item after popping 2 items")

		mlen, err = mapLen[uint64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(t, qlen, uint64(mlen),
			"Queue length should match container length")

		items, err = q.PopMulti(ctx, 3)
		require.NoError(t, err)
		require.Equal(t, []uint64{3}, items)

		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), qlen,
			"Queue should be empty after popping all items")

		mlen, err = mapLen[uint64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(t, qlen, uint64(mlen),
			"Queue length should match container length")

		items, err = q.PopMulti(ctx, 3)
		require.NoError(t, err)
		require.Empty(t, items)
	})
}

func Test_PushMulti(t *testing.T) {
	t.Run("should return correct items and lengths", func(t *testing.T) {
		sk, ctx := deps()
		sb := sdk.NewSchemaBuilder(sk)
		q := collections.NewQueue[uint64](sb, "queue", sdk.Uint64Value)

		err := q.PushMulti(ctx, []uint64{1, 2, 3})
		require.NoError(t, err)

		// Test length after pushes
		qlen, err := q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), qlen,
			"Queue should have 3 items after 3 pushes")

		mlen, err := mapLen[uint64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(t, qlen, uint64(mlen),
			"Queue length should match container length")

		items, err := q.PopMulti(ctx, 4)
		require.NoError(t, err)
		require.Equal(t, []uint64{1, 2, 3}, items)

		qlen, err = q.Len(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), qlen,
			"Queue should be empty after popping all items")

		mlen, err = mapLen[uint64](ctx, q.Container())
		require.NoError(t, err)
		require.Equal(t, qlen, uint64(mlen),
			"Queue length should match container length")
	})
}

func mapLen[V any](ctx context.Context, m sdk.Map[uint64, V]) (int, error) {
	iter, err := m.IterateRaw(ctx, nil, nil, sdk.OrderAscending)
	if err != nil {
		return 0, err
	}
	keys, err := iter.Keys()
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

// MockStore wraps the dba.Store to implement additional functionalities.
type MockStore struct {
	dba.Store
}

// OpenKVStore returns the underlying KVStore from the Store.
func (s MockStore) OpenKVStore(context.Context) store.KVStore {
	return s
}

func (s MockStore) Delete(key []byte) error {
	s.Store.Delete(key)
	return nil
}

func (s MockStore) Set(key, value []byte) error {
	s.Store.Set(key, value)
	return nil
}

func (s MockStore) Get(key []byte) ([]byte, error) {
	return s.Store.Get(key), nil
}

func (s MockStore) Has(key []byte) (bool, error) {
	return s.Store.Has(key), nil
}

// Iterator wraps the underlying DB's Iterator method panicing on error.
func (s MockStore) Iterator(start, end []byte) (store.Iterator, error) {
	iter, err := s.DB.Iterator(start, end)
	if err != nil {
		return nil, err
	}

	return iter, nil
}

// ReverseIterator wraps the underlying DB's ReverseIterator method panicing on error.
func (s MockStore) ReverseIterator(start, end []byte) (store.Iterator, error) {
	iter, err := s.DB.ReverseIterator(start, end)
	if err != nil {
		return nil, err
	}

	return iter, nil
}

// deps initializes dependencies for testing, returning a KVStoreService
// and a context.
func deps() (store.KVStoreService, context.Context) {
	return &MockStore{
		Store: dba.Store{DB: db.NewMemDB()},
	}, context.Background()
}
