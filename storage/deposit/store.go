// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/encoding"
)

const KeyDepositPrefix = "deposit"

// KVStore is a simple KV store based implementation that assumes
// the deposit indexes are tracked outside of the kv store.
type KVStore struct {
	store sdkcollections.Map[uint64, *ctypes.Deposit]

	// closeFunc is a closure that closes the underlying database
	// used by store to ensure that all writes are flushed to disk.
	// We guarantee that closeFunc is called at maximum only once.
	closeFunc CloseFunc
	once      sync.Once

	// mu protects store for concurrent access
	mu sync.RWMutex

	// logger is used for logging information and errors.
	logger log.Logger
}

// closure type for closing the store.
type CloseFunc func() error

// NewStore creates a new deposit store.
func NewStore(
	kvsp store.KVStoreService,
	closeFunc CloseFunc,
	logger log.Logger,
) *KVStore {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	res := &KVStore{
		store: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte(KeyDepositPrefix)),
			KeyDepositPrefix,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[*ctypes.Deposit]{
				NewEmptyF: ctypes.NewEmptyDeposit,
			},
		),
		closeFunc: closeFunc,
		logger:    logger,
	}
	if _, err := schemaBuilder.Build(); err != nil {
		panic(errors.Wrap(err, "failed building KVStore schema"))
	}
	return res
}

// Close closes the store by calling the closeFunc. It ensures that the
// closeFunc is called at most once.
func (kv *KVStore) Close() error {
	var err error
	kv.once.Do(func() { err = kv.closeFunc() })
	return err
}

// GetDepositsByIndex returns the first N deposits starting from the given
// index. If N is greater than the number of deposits, it returns up to the
// last deposit.
func (kv *KVStore) GetDepositsByIndex(
	ctx context.Context,
	startIndex uint64,
	depRange uint64,
) (ctypes.Deposits, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	var (
		deposits = make(ctypes.Deposits, 0, depRange)
		endIdx   = startIndex + depRange
	)

	for i := startIndex; i < endIdx; i++ {
		deposit, err := kv.store.Get(ctx, i)
		switch {
		case err == nil:
			deposits = append(deposits, deposit)
		case errors.Is(err, sdkcollections.ErrNotFound):
			return deposits, nil
		default:
			return deposits, errors.Wrapf(
				err, "failed to get deposit %d, start: %d, end: %d", i, startIndex, endIdx,
			)
		}
	}

	kv.logger.Debug("GetDepositsByIndex", "start", startIndex, "end", endIdx)
	return deposits, nil
}

// EnqueueDeposits pushes multiple deposits to the queue.
func (kv *KVStore) EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for _, deposit := range deposits {
		idx := deposit.GetIndex().Unwrap()
		if err := kv.store.Set(ctx, idx, deposit); err != nil {
			return errors.Wrapf(err, "failed to enqueue deposit %d", idx)
		}
	}

	if len(deposits) > 0 {
		kv.logger.Debug(
			"EnqueueDeposit", "enqueued", len(deposits),
			"start", deposits[0].GetIndex(), "end", deposits[len(deposits)-1].GetIndex(),
		)
	}
	return nil
}

// Prune removes the [start, end) deposits from the store.
func (kv *KVStore) Prune(ctx context.Context, start, end uint64) error {
	if start > end {
		return errors.Wrapf(
			storage.ErrInvalidRange, "DepositKVStore Prune start: %d, end: %d", start, end)
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()
	for i := range end {
		// This only errors if the key passed in cannot be encoded.
		if err := kv.store.Remove(ctx, start+i); err != nil {
			return errors.Wrapf(err, "failed to prune deposit %d", start+i)
		}
	}

	kv.logger.Debug("Pruned deposits", "start", start, "end", end)
	return nil
}
