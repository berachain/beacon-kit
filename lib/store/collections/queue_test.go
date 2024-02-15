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
	"github.com/stretchr/testify/require"

	"github.com/itsdevbear/bolaris/lib/store/collections"
)

func Test_Queue(t *testing.T) {
	t.Run("should return correct items", func(t *testing.T) {
		sk, ctx := deps()
		sb := sdk.NewSchemaBuilder(sk)
		q := collections.NewQueue[uint64](sb, "queue", sdk.Uint64Value)

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

type Store struct {
	dba.Store
}

func (s Store) Get(key []byte) ([]byte, error) {
	return s.Store.Get(key), nil
}

func (s Store) Has(key []byte) (bool, error) {
	return s.Store.Has(key), nil
}

func (s Store) Iterator(start, end []byte) (db.Iterator, error) {
	return s.Store.Iterator(start, end), nil
}

func (s Store) ReverseIterator(start, end []byte) (db.Iterator, error) {
	return s.Store.ReverseIterator(start, end), nil
}

func (s Store) Set(key, value []byte) error {
	s.Store.Set(key, value)
	return nil
}

func (s Store) Delete(key []byte) error {
	s.Store.Delete(key)
	return nil
}

func (s Store) OpenKVStore(ctx context.Context) store.KVStore {
	return s
}

func deps() (store.KVStoreService, context.Context) {
	db := db.NewMemDB()
	return &Store{Store: dba.Store{DB: db}}, context.Background()
}
