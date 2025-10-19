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
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	sdkmetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestStateProtect(t *testing.T) {
	t.Parallel()

	db, err := db.OpenDB("", dbm.MemDBBackend)
	require.NoError(t, err)

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	var (
		nopLog     = log.NewNopLogger()
		nopMetrics = sdkmetrics.NewNoOpMetrics()
	)

	cms := store.NewCommitMultiStore(db, nopLog, nopMetrics)
	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	require.NoError(t, cms.LoadLatestVersion())

	backendStoreService := &storage.KVStoreService{Key: testStoreKey}
	kvStore := beacondb.New(backendStoreService)

	ms := cms.CacheMultiStore()
	sdkCtx := sdk.NewContext(ms, true, nopLog)
	originalState := state.NewBeaconStateFromDB(
		kvStore.WithContext(sdkCtx),
		cs,
		sdkCtx.Logger(),
		metrics.NewNoOpTelemetrySink(),
	)

	protectingState := originalState.Protect(sdkCtx)

	// 1- set an attribute in the original state and show
	// that value is carried over the protecting state
	wantSlot := math.Slot(1234)
	require.NoError(t, originalState.SetSlot(wantSlot))

	gotSlot, err := protectingState.GetSlot()
	require.NoError(t, err)
	require.Equal(t, wantSlot, gotSlot)

	// 2- Show that modifying the protecting state
	// does not affect the original state
	wantFork := &ctypes.Fork{
		PreviousVersion: common.Version{0x11, 0x22, 0x33, 0x44},
		CurrentVersion:  common.Version{0xff, 0xff, 0xff, 0xff},
		Epoch:           math.Epoch(1234),
	}
	require.NoError(t, protectingState.SetFork(wantFork))

	_, err = originalState.GetFork()
	require.ErrorIs(t, err, collections.ErrNotFound)

	// 3- Show that changes made to original state POST COPY
	// are carried over the protecting state
	wantEthIdx := uint64(1987)
	require.NoError(t, originalState.SetEth1DepositIndex(wantEthIdx))

	gotEthIdx, err := protectingState.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, wantEthIdx, gotEthIdx)
}

var testStoreKey = storetypes.NewKVStoreKey("test-stateDB")
