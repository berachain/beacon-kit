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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/storage"
	depositstorecommon "github.com/berachain/beacon-kit/storage/deposit/common"
	"github.com/berachain/beacon-kit/storage/encoding"
	dbm "github.com/cosmos/cosmos-db"
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

	// logger is used for logging information and errors.
	logger log.Logger
}

// closure type for closing the store.
type CloseFunc func() error

// NewStore creates a new deposit store.
func NewStore(
	baseDB dbm.DB,
	logger log.Logger,
) *KVStore {
	spdbV1 := depositstorecommon.NewSynced(baseDB)
	kvsp := NewKVStoreProvider(spdbV1)
	closeFunc := spdbV1.Close

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
// Note: we return the hash root of the selected deposits to simplify block
// building pre migration to deposit store V2.
func (kv *KVStore) GetDepositsByIndex(
	ctx context.Context,
	startIndex uint64,
	depRange uint64,
) (ctypes.Deposits, common.Root, error) {
	var (
		deposits = make(ctypes.Deposits, 0, depRange)
		endIdx   = startIndex + depRange
	)

	done := false
	for i := startIndex; i < endIdx && !done; i++ {
		deposit, err := kv.store.Get(ctx, i)
		switch {
		case err == nil:
			deposits = append(deposits, deposit)
		case errors.Is(err, sdkcollections.ErrNotFound):
			done = true // normal happy path, there are less than max allowed deposits
		default:
			return deposits, common.Root{}, errors.Wrapf(
				err, "failed to get deposit %d, start: %d, end: %d", i, startIndex, endIdx,
			)
		}
	}

	kv.logger.Debug("GetDepositsByIndex", "start", startIndex, "end", endIdx)
	rootBytes, err := deposits.HashTreeRoot()
	if err != nil {
		return nil, common.Root{}, err
	}
	return deposits, common.NewRootFromBytes(rootBytes[:]), nil
}

// EnqueueDeposits pushes multiple deposits to the queue.
func (kv *KVStore) EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error {
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

	for i := range end {
		// This only errors if the key passed in cannot be encoded.
		if err := kv.store.Remove(ctx, start+i); err != nil {
			return errors.Wrapf(err, "failed to prune deposit %d", start+i)
		}
	}

	kv.logger.Debug("Pruned deposits", "start", start, "end", end)
	return nil
}
