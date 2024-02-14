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
	db "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	"github.com/itsdevbear/bolaris/collections"
)

func Test_Queue(t *testing.T) {
	t.Run("should return correct items", func(t *testing.T) {
		sk, ctx := deps()
		sb := sdk.NewSchemaBuilder(sk)
		q := collections.NewQueue[uint64](sb, sdk.NewPrefix(0), "queue", sdk.Uint64Value)

		_, err := q.Peek(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		_, err = q.Next(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		err = q.Push(ctx, 1)
		require.NoError(t, err)
		err = q.Push(ctx, 2)
		require.NoError(t, err)

		v, err := q.Next(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), v)
		v, err = q.Next(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), v)
		_, err = q.Next(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		err = q.Push(ctx, 3)
		require.NoError(t, err)

		v, err = q.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), v)

		err = q.Push(ctx, 4)
		require.NoError(t, err)

		v, err = q.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), v)

		v, err = q.Next(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), v)

		v, err = q.Next(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(4), v)

		_, err = q.Peek(ctx)
		require.Equal(t, sdk.ErrNotFound, err)

		_, err = q.Next(ctx)
		require.Equal(t, sdk.ErrNotFound, err)
	})
}

type testStore struct {
	db db.DB
}

func (t testStore) OpenKVStore(ctx context.Context) store.KVStore {
	return t
}

func (t testStore) Get(key []byte) ([]byte, error) {
	return t.db.Get(key)
}

func (t testStore) Has(key []byte) (bool, error) {
	return t.db.Has(key)
}

func (t testStore) Set(key, value []byte) error {
	return t.db.Set(key, value)
}

func (t testStore) Delete(key []byte) error {
	return t.db.Delete(key)
}

func (t testStore) Iterator(start, end []byte) (store.Iterator, error) {
	return t.db.Iterator(start, end)
}

func (t testStore) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return t.db.ReverseIterator(start, end)
}

var _ store.KVStore = testStore{}

func deps() (store.KVStoreService, context.Context) {
	kv := db.NewMemDB()
	return &testStore{kv}, context.Background()
}
