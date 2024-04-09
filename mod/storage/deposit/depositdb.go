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
	"errors"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	collections "github.com/berachain/beacon-kit/mod/storage/beacondb/collections"
	encoding "github.com/berachain/beacon-kit/mod/storage/beacondb/collections/encoding"
	cdb "github.com/cosmos/cosmos-db"
)

type Store struct {
	cdb.DB
}

const (
	KeyDepositPrefix = "deposit"
)

var _ store.KVStoreService = (*KVStoreProvider)(nil)

type KVStoreProvider struct {
	*PebbleDB
}

func NewKVStoreProvider(name, backend, dir string) (*KVStoreProvider, error) {
	var db *PebbleDB
	var err error
	switch backend {
	case "pebble":
		db, err = NewPebbleDB(name, dir)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported backend")
	}

	return &KVStoreProvider{
		db,
	}, nil
}

// OpenKVStore opens a new KV store.
func (p *KVStoreProvider) OpenKVStore(context.Context) store.KVStore {
	return p.PebbleDB
}

// KVStore is a wrapper around an sdk.Context.
type KVStore struct {
	depositQueue *collections.Queue[*beacontypes.Deposit]
}

// NewStore creates a new deposit store.
func NewStore(kvsp store.KVStoreService) *KVStore {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore{
		depositQueue: collections.NewQueue[*beacontypes.Deposit](
			schemaBuilder,
			KeyDepositPrefix,
			encoding.SSZValueCodec[*beacontypes.Deposit]{},
		),
	}
}

// ExpectedDeposits returns the first numPeek deposits in the queue.
func (kv *KVStore) ExpectedDeposits(
	numView uint64,
) (beacontypes.Deposits, error) {
	return kv.depositQueue.PeekMulti(context.TODO(), numView)
}

// EnqueueDeposits pushes the deposits to the queue.
func (kv *KVStore) EnqueueDeposits(
	deposits beacontypes.Deposits,
) error {
	return kv.depositQueue.PushMulti(context.TODO(), deposits)
}

// DequeueDeposits returns the first numDequeue deposits in the queue.
func (kv *KVStore) DequeueDeposits(
	numDequeue uint64,
) (beacontypes.Deposits, error) {
	return kv.depositQueue.PopMulti(context.TODO(), numDequeue)
}
