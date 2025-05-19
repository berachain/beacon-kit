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
	"math"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	sdklog "cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/storage"
	depositstorecommon "github.com/berachain/beacon-kit/storage/deposit/common"
	"github.com/berachain/beacon-kit/storage/encoding"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint:gochecknoglobals // storeKey is a singleton.
var (
	DepositStoreKey    = storetypes.NewKVStoreKey("deposits")
	depositStorePrefix = sdkcollections.NewPrefix(0)
	depositStoreName   = "deposits"

	// migrationFlagDeposit is used to understand whether migration
	// from storeV1 is ongoing. The flag deposit has an index high enough
	// so that it could only be a placeholder
	migrationFlagDeposit = &ctypes.Deposit{
		Index: math.MaxUint64,
	}
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

	closeFunc CloseFunc
	once      sync.Once

	logger log.Logger

	depositsRoot common.Root
}

func NewStore(
	baseDB dbm.DB,
	logger log.Logger,
) *KVStore {
	db := depositstorecommon.NewSynced(baseDB)
	closeFn := db.Close

	// TODO ABENEGIA: fix logging
	cms := store.NewCommitMultiStore(db, sdklog.NewNopLogger(), storemetrics.NewNoOpMetrics())
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
		logger:       logger,
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

func (kv *KVStore) EnqueueDeposits(_ context.Context, deposits []*ctypes.Deposit) error {
	// create context out of commit multistore, simiarly to what we do in consensus service
	ms := kv.cms.CacheMultiStore()
	sdkCtx := sdk.NewContext(ms, false, sdklog.NewNopLogger()) // .WithContext(ctx)

	for _, deposit := range deposits {
		idx := deposit.GetIndex().Unwrap()
		//nolint:contextcheck // TODO ABENEGIA: to be fixed
		if err := kv.store.Set(sdkCtx, idx, deposit); err != nil {
			return errors.Wrapf(err, "failed to enqueue deposit %d", idx)
		}
	}
	ms.Write()
	commit := kv.cms.Commit()

	root, err := bytesToRoot(commit.Hash)
	if err != nil {
		panic(err)
	}
	kv.depositsRoot = root

	if len(deposits) > 0 {
		kv.logger.Debug(
			"EnqueueDeposit", "enqueued", len(deposits),
			"start", deposits[0].GetIndex(), "end", deposits[len(deposits)-1].GetIndex(),
		)
	}
	return nil
}

func (kv *KVStore) GetDepositsByIndex(
	_ context.Context, // we use the internal context here
	startIndex uint64,
	depRange uint64,
) (
	ctypes.Deposits,
	common.Root, // deposits common root
	error,
) {
	var (
		deposits = make(ctypes.Deposits, 0, depRange)
		endIdx   = startIndex + depRange
		sdkCtx   = sdk.NewContext(kv.cms, false, sdklog.NewNopLogger())
	)

	for i := startIndex; i < endIdx; i++ {
		//nolint:contextcheck // TODO ABENEGIA: to be fixed
		deposit, err := kv.store.Get(sdkCtx, i)
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

	kv.logger.Debug("GetDepositsByIndex", "start", startIndex, "end", endIdx)
	return deposits, kv.depositsRoot, nil
}

func (kv *KVStore) Prune(_ context.Context, start, end uint64) error {
	kv.logger.Debug("GetDepositsByIndex", "start", start, "end", end)
	return nil
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

func (kv *KVStore) MarkMigrationFromV1Started(ctx context.Context) error {
	founds, _, err := kv.GetDepositsByIndex(ctx, migrationFlagDeposit.Index, 1)
	if err != nil {
		return fmt.Errorf("failed checking whether migration has started: %w", err)
	}

	if len(founds) != 0 {
		kv.logger.Warn("Deposit DB  migration already started")
		return nil
	}

	toEnqueue := []*ctypes.Deposit{migrationFlagDeposit}
	if err = kv.EnqueueDeposits(ctx, toEnqueue); err != nil {
		return fmt.Errorf("failed marking migration has started: %w", err)
	}
	return nil
}

func (kv *KVStore) MarkMigrationFromV1Done(_ context.Context) error {
	sdkCtx := sdk.NewContext(kv.cms, false, sdklog.NewNopLogger())
	//nolint:contextcheck // TODO ABENEGIA: to be fixed
	if err := kv.store.Remove(sdkCtx, migrationFlagDeposit.Index); err != nil {
		return fmt.Errorf("failed marking migration has completed: %w", err)
	}
	return nil
}

func (kv *KVStore) HasMigrationFromV1Completed(ctx context.Context) (bool, error) {
	// Check whether migration is ongoing
	founds, _, err := kv.GetDepositsByIndex(ctx, migrationFlagDeposit.Index, 1)
	if err != nil {
		return false, fmt.Errorf("failed checking whether migration has started: %w", err)
	}

	if len(founds) != 0 {
		kv.logger.Warn("Deposit DB  migration already started but not completed")
		return false, nil
	}

	// Migration has happened if there is at least one deposit in the store
	founds, _, err = kv.GetDepositsByIndex(ctx, constants.FirstDepositIndex, 1)
	if err != nil {
		return false, fmt.Errorf("failed checking whether first deposit is present: %w", err)
	}

	return len(founds) != 0, nil
}
