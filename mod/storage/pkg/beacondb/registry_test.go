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

package beacondb_test

import (
	"context"
	"fmt"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/db"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type testKVStoreService struct {
	ctx sdk.Context
}

func (kvs *testKVStoreService) OpenKVStore(context.Context) corestore.KVStore {
	//nolint:contextcheck // fine with tests
	return components.NewKVStore(
		sdk.UnwrapSDKContext(kvs.ctx).KVStore(testStoreKey),
	)
}

var (
	testStoreKey = storetypes.NewKVStoreKey("storage-tests")
	testCodec    = &encoding.SSZInterfaceCodec[*types.ExecutionPayloadHeader]{}
)

func TestGetBalances(t *testing.T) {
	ctx, err := initKvStore()
	require.NoError(t, err)
	testStoreService := &testKVStoreService{
		ctx: ctx,
	}

	kvStore := beacondb.New[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		[]*types.Validator,
	](
		testStoreService,
		testCodec,
	)

	// no balance to start
	res, err := kvStore.GetBalances()
	require.NoError(t, err)
	require.Zero(t, res)

	// add balances
	var (
		idx1, idx2     = math.U64(1_987), math.U64(1_989)
		inBal1, inBal2 = math.U64(8_992), math.U64(10_000)
	)
	require.NoError(t, kvStore.SetBalance(idx1, inBal1))
	require.NoError(t, kvStore.SetBalance(idx2, inBal2))

	// check we can query added balances
	balRes, err := kvStore.GetBalance(idx1)
	require.NoError(t, err)
	require.Equal(t, balRes, inBal1)

	balRes, err = kvStore.GetBalance(idx2)
	require.NoError(t, err)
	require.Equal(t, balRes, inBal2)

	res, err = kvStore.GetBalances()
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, res[0], inBal1.Unwrap())
	require.Equal(t, res[1], inBal2.Unwrap())

	// update existing balances
	newInBal1, newInBal2 := math.U64(0), inBal2*2
	require.NoError(t, kvStore.SetBalance(idx1, newInBal1))
	require.NoError(t, kvStore.SetBalance(idx2, newInBal2))

	// check we can query updated balances
	balRes, err = kvStore.GetBalance(idx1)
	require.NoError(t, err)
	require.Equal(t, balRes, newInBal1)

	balRes, err = kvStore.GetBalance(idx2)
	require.NoError(t, err)
	require.Equal(t, balRes, newInBal2)

	res, err = kvStore.GetBalances()
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, res[0], newInBal1.Unwrap())
	require.Equal(t, res[1], newInBal2.Unwrap())
}

func initKvStore() (sdk.Context, error) {
	db, err := db.OpenDB("", dbm.MemDBBackend)
	if err != nil {
		return sdk.Context{}, fmt.Errorf("failed opening mem db: %w", err)
	}
	var (
		nopLog     = log.NewNopLogger()
		nopMetrics = metrics.NewNoOpMetrics()
	)

	cms := store.NewCommitMultiStore(
		db,
		nopLog,
		nopMetrics,
	)

	ctx := sdk.NewContext(cms, true, nopLog)
	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	if err = cms.LoadLatestVersion(); err != nil {
		return sdk.Context{}, fmt.Errorf("failed to load latest version: %w", err)
	}
	return ctx, nil
}
