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

package state_test

import (
	"context"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var testStoreKey = storetypes.NewKVStoreKey("state-transition-tests")

type testKVStoreService struct{}

func (kvs *testKVStoreService) OpenKVStore(ctx context.Context) corestore.KVStore {
	store := sdk.UnwrapSDKContext(ctx).KVStore(testStoreKey)
	return storage.NewKVStore(store)
}

func TestStateCopy(t *testing.T) {
	t.Parallel()

	// define specs
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	// create state backing a KV Store
	db, err := db.OpenDB("", dbm.MemDBBackend)
	require.NoError(t, err)

	nopLog := log.NewNopLogger()
	cms := store.NewCommitMultiStore(db, nopLog, metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	require.NoError(t, cms.LoadLatestVersion())

	ctx := sdk.NewContext(cms.CacheMultiStore(), true, nopLog)
	testStoreService := &testKVStoreService{}
	kv := beacondb.New(testStoreService)

	// create beacon state
	st := state.NewBeaconStateFromDB(kv.WithContext(ctx), cs)
	// basically what StateFromContext does. WithContext is key since
	// is sets the context containing the persistence layer

	// Store some data in st
	v1 := math.Slot(2025)
	require.NoError(t, st.SetSlot(v1))

	// Make another state to pass to the same kvStore and show that data are not there
	cpyCtx := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())
	cpy := state.NewBeaconStateFromDB(kv.WithContext(cpyCtx), cs)

	_, err = cpy.GetSlot()
	require.Error(t, err)
}
