// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package deposit

import (
	"context"
	"errors"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
)

const KeyDepositPrefix = "deposit"

// KVStore is a simple KV store based implementation that assumes
// the deposit indexes are tracked outside of the kv store.
type KVStore[DepositT Deposit[DepositT]] struct {
	store sdkcollections.Map[uint64, DepositT]
	mu    sync.RWMutex
}

// NewStore creates a new deposit store.
func NewStore[DepositT Deposit[DepositT]](
	kvsp store.KVStoreService,
) *KVStore[DepositT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore[DepositT]{
		store: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte(KeyDepositPrefix)),
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
		switch {
		case err == nil:
			deposits = append(deposits, deposit)
		case errors.Is(err, sdkcollections.ErrNotFound):
			// not more deposits from index i on.
			return deposits, nil
		default:
			return nil, err
		}
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
	return kv.store.Set(context.TODO(), deposit.GetIndex().Unwrap(), deposit)
}

// Prune removes the [start, end) deposits from the store.
func (kv *KVStore[DepositT]) Prune(start, end uint64) error {
	if start > end {
		return pruner.ErrInvalidRange
	}

	var ctx = context.TODO()
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for i := range end {
		if err := kv.store.Remove(ctx, start+i); err != nil {
			return err
		}
	}
	return nil
}
