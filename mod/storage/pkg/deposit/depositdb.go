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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	encoding "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/collections/encoding"
)

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
	depositQueue *Queue
}

// NewStore creates a new deposit store.
func NewStore(kvsp store.KVStoreService) *KVStore {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore{
		depositQueue: NewQueue(
			schemaBuilder,
			KeyDepositPrefix,
			encoding.SSZValueCodec[*consensus.Deposit]{},
		),
	}
}

// ExpectedDeposits returns the first numPeek deposits in the queue.
func (kv *KVStore) ExpectedDeposits(
	numView uint64,
) ([]*consensus.Deposit, error) {
	return kv.depositQueue.PeekMulti(context.TODO(), numView)
}

// EnqueueDeposit pushes the deposit to the queue.
func (kv *KVStore) EnqueueDeposit(deposit *consensus.Deposit) error {
	return kv.depositQueue.Push(context.TODO(), deposit)
}

// EnqueueDeposits pushes multiple deposits to the queue.
func (kv *KVStore) EnqueueDeposits(deposits []*consensus.Deposit) error {
	return kv.depositQueue.PushMulti(context.TODO(), deposits)
}

// DequeueDeposits returns the first numDequeue deposits in the queue.
func (kv *KVStore) DequeueDeposits(
	numDequeue uint64,
) ([]*consensus.Deposit, error) {
	return kv.depositQueue.PopMulti(context.TODO(), numDequeue)
}

// PruneToIndex removes all deposits up to the given index.
func (kv *KVStore) PruneToIndex(
	index uint64,
) error {
	length, err := kv.depositQueue.Len(context.TODO())
	if err != nil {
		return err
	} else if length == 0 {
		return nil
	}

	head, err := kv.depositQueue.Peek(context.TODO())
	if err != nil {
		return err
	}

	numPop := min(index-head.Index+1, length)
	_, err = kv.depositQueue.PopMulti(context.TODO(), numPop)
	return err
}
