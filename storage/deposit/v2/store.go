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
	"fmt"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/encoding"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint:gochecknoglobals // storeKey is a singleton.
var (
	DepositStoreKey    = storetypes.NewKVStoreKey("deposits")
	depositStorePrefix = sdkcollections.NewPrefix(0)
	depositStoreName   = "deposits"
)

type KVStoreService struct{}

func (k KVStoreService) OpenKVStore(ctx context.Context) corestore.KVStore {
	return storage.NewKVStore(sdk.UnwrapSDKContext(ctx).KVStore(DepositStoreKey))
}

// closure type for closing the store.
// TODO: consider integrating this store into consensus service one (separate store)
type CloseFunc func() error

type KVStore struct {
	store sdkcollections.Map[uint64, *ctypes.Deposit]

	// TODO ABENEGIA: consolidate within consensus service multistore
	cms storetypes.CommitMultiStore
	mu  sync.RWMutex // mu protects store for concurrent access

	closeFunc CloseFunc
	once      sync.Once

	depositsRoot common.Root
}

func NewStore(
	baseDB dbm.DB,
	metrics metrics.StoreMetrics,
) *KVStore {
	db := NewSynced(baseDB)
	closeFn := db.Close

	// TODO ABENEGIA: fix logging
	cms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics)
	cms.MountStoreWithDB(DepositStoreKey, storetypes.StoreTypeIAVL, nil)
	if err := cms.LoadLatestVersion(); err != nil {
		panic(fmt.Errorf("deposit store v2: failed loading latest version: %w", err))
	}

	schemaBuilder := sdkcollections.NewSchemaBuilder(KVStoreService{})
	store := sdkcollections.NewMap(
		schemaBuilder,
		depositStorePrefix,
		depositStoreName,
		sdkcollections.Uint64Key,
		encoding.SSZValueCodec[*ctypes.Deposit]{
			NewEmptyF: ctypes.NewEmptyDeposit,
		},
	)
	if _, err := schemaBuilder.Build(); err != nil {
		panic(fmt.Errorf("failed building deposits store schema: %w", err))
	}
	root, err := bytesToRoot(cms.WorkingHash())
	if err != nil {
		panic(err)
	}

	return &KVStore{
		store:        store,
		cms:          cms,
		closeFunc:    closeFn,
		depositsRoot: root,
	}
}

// Close closes the store by calling the closeFunc. It ensures that the
// closeFunc is called at most once.
func (kv *KVStore) Close() error {
	var err error
	kv.once.Do(func() { err = kv.closeFunc() })
	return err
}

func (kv *KVStore) EnqueueDeposits( /*ctx context.Context,*/ deposits []*ctypes.Deposit) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// create context out of commit multistore, simiarly to what we do in consensus service
	// ms := kv.cms.CacheMultiStore()
	sdkCtx := sdk.NewContext(kv.cms /*ms*/, false, log.NewNopLogger()) // .WithContext(ctx)

	for _, deposit := range deposits {
		idx := deposit.GetIndex().Unwrap()
		if err := kv.store.Set(sdkCtx, idx, deposit); err != nil {
			return errors.Wrapf(err, "failed to enqueue deposit %d", idx)
		}
	}
	// ms.Write()
	kv.cms.Commit()

	root, err := bytesToRoot(kv.cms.WorkingHash())
	if err != nil {
		panic(err)
	}
	kv.depositsRoot = root

	// TODO ABENEGIA: re-add logging
	// if len(deposits) > 0 {
	// 	kv.logger.Debug(
	// 		"EnqueueDeposit", "enqueued", len(deposits),
	// 		"start", deposits[0].GetIndex(), "end", deposits[len(deposits)-1].GetIndex(),
	// 	)
	// }
	return nil
}

func (kv *KVStore) GetDepositsByIndex(
	// ctx context.Context, // we use the internal context here
	startIndex uint64,
	depRange uint64,
) (
	ctypes.Deposits,
	common.Root, // deposits common root
	error,
) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	var (
		deposits = make(ctypes.Deposits, 0, depRange)
		endIdx   = startIndex + depRange
		ctx      = sdk.NewContext(kv.cms, false, log.NewNopLogger())
	)

	for i := startIndex; i < endIdx; i++ {
		deposit, err := kv.store.Get(ctx, i)
		switch {
		case err == nil:
			deposits = append(deposits, deposit)
		case errors.Is(err, sdkcollections.ErrNotFound):
			return deposits, kv.depositsRoot, nil
		default:
			return deposits, kv.depositsRoot, errors.Wrapf(
				err, "failed to get deposit %d, start: %d, end: %d", i, startIndex, endIdx,
			)
		}
	}

	// TODO ABENEGIA: re-add logging
	// kv.logger.Debug("GetDepositsByIndex", "start", startIndex, "end", endIdx)
	return deposits, kv.depositsRoot, nil
}

func bytesToRoot(b []byte) (common.Root, error) {
	if len(b) != len(common.Root{}) {
		return common.Root{}, fmt.Errorf(
			"working has length %d not compatible with common Root",
			len(b),
		)
	}
	return common.Root(b), nil
}
