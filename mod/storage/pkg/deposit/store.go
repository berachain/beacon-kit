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

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	encoding "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
)

const (
	KeyDepositPrefix = "deposit"
)

type KVStoreProvider struct {
	store.KVStoreWithBatch
}

// OpenKVStore opens a new KV store.
func (p *KVStoreProvider) OpenKVStore(context.Context) store.KVStore {
	return p.KVStoreWithBatch
}

// KVStore is a wrapper around an sdk.Context.
type KVStore[DepositT Deposit] struct {
	depositQueue *Queue[DepositT]
}

// NewStore creates a new deposit store.
func NewStore[DepositT Deposit](kvsp store.KVStoreService) *KVStore[DepositT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore[DepositT]{
		depositQueue: NewQueue(
			schemaBuilder,
			KeyDepositPrefix,
			encoding.SSZValueCodec[DepositT]{},
		),
	}
}

// ExpectedDeposits returns the first numPeek deposits in the queue.
func (kv *KVStore[DepositT]) ExpectedDeposits(
	numView uint64,
) ([]DepositT, error) {
	return kv.depositQueue.PeekMulti(context.TODO(), numView)
}

// EnqueueDeposit pushes the deposit to the queue.
func (kv *KVStore[DepositT]) EnqueueDeposit(deposit DepositT) error {
	return kv.depositQueue.Push(context.TODO(), deposit)
}

// EnqueueDeposits pushes multiple deposits to the queue.
func (kv *KVStore[DepositT]) EnqueueDeposits(deposits []DepositT) error {
	return kv.depositQueue.PushMulti(context.TODO(), deposits)
}

// DequeueDeposits returns the first numDequeue deposits in the queue.
func (kv *KVStore[DepositT]) DequeueDeposits(
	numDequeue uint64,
) ([]DepositT, error) {
	return kv.depositQueue.PopMulti(context.TODO(), numDequeue)
}

// Prune removes all deposits up to the given index.
func (kv *KVStore[DepositT]) Prune(
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

	numPop := min(index-head.GetIndex()+1, length)
	_, err = kv.depositQueue.PopMulti(context.TODO(), numPop)
	return err
}
