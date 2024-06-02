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
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
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

// KVStore is a simple KV store based implementation that assumes
// the deposit indexes are tracked outside of the kv store.
type KVStore[DepositT Deposit] struct {
	store sdkcollections.Map[uint64, DepositT]
	mu    sync.RWMutex
}

// NewStore creates a new deposit store.
func NewStore[DepositT Deposit](kvsp store.KVStoreService) *KVStore[DepositT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore[DepositT]{
		store: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{uint8(0)}),
			KeyDepositPrefix,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[DepositT]{},
		),
	}
}

// GetDepositsByIndex returns the first N deposits starting from the given
// index. If N is greater than the number of deposits, it returns up to the
// last deposit.
func (kv *KVStore[DepositT]) GetDepositsByIndex(
	startIndex uint64,
	numView uint64,
) ([]DepositT, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	deposits := []DepositT{}
	for i := range numView {
		deposit, err := kv.store.Get(context.TODO(), startIndex+i)
		if errors.Is(err, sdkcollections.ErrNotFound) {
			return deposits, nil
		}
		if err != nil {
			return deposits, err
		}
		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

// EnqueueDeposit pushes the deposit to the queue.
func (kv *KVStore[DepositT]) EnqueueDeposit(deposit DepositT) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	return kv.setDeposit(deposit)
}

// EnqueueDeposits pushes multiple deposits to the queue.
func (kv *KVStore[DepositT]) EnqueueDeposits(deposits []DepositT) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for _, deposit := range deposits {
		if err := kv.setDeposit(deposit); err != nil {
			return err
		}
	}
	return nil
}

// setDeposit sets the deposit in the store.
func (kv *KVStore[DepositT]) setDeposit(deposit DepositT) error {
	return kv.store.Set(context.TODO(), deposit.GetIndex(), deposit)
}

// PruneFromInclusive removes up to N deposits from the given starting index.
func (kv *KVStore[DepositT]) PruneFromInclusive(
	index uint64,
	numPrune uint64,
) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for i := range numPrune {
		// This only errors if the key passed in cannot be encoded.
		err := kv.store.Remove(context.TODO(), index+i)
		if err != nil {
			return err
		}
	}
	return nil
}
