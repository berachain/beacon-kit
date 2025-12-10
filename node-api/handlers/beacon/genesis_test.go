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

package beacon_test

import (
	"context"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	sdkmetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	beaconlog "github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	kvstorage "github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetGenesis(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	var (
		testGenesisRoot        = common.Root{0x10, 0x20, 0x30}
		testGenesisForkVersion = common.Version{0xff, 0x11, 0x22, 0x33}
		testGenesisTime        = math.U64(123456789)
	)

	testCases := []struct {
		name                string
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "success",
			setMockExpectations: func(b *mocks.Backend) {
				st := initTestGenesisState(t, cs)

				require.NoError(t, st.SetGenesisValidatorsRoot(testGenesisRoot))
				require.NoError(t, st.SetFork(&ctypes.Fork{
					PreviousVersion: testGenesisForkVersion,
					CurrentVersion:  testGenesisForkVersion,
				}))
				require.NoError(t, st.SetLatestExecutionPayloadHeader(&ctypes.ExecutionPayloadHeader{
					Versionable: ctypes.NewVersionable(testGenesisForkVersion),
					Timestamp:   testGenesisTime,
				}))

				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, 0, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenesisResponse{}, res)
				gr, _ := res.(beacontypes.GenesisResponse)

				require.Equal(t, testGenesisRoot, gr.Data.GenesisValidatorsRoot)
				require.Equal(t, testGenesisForkVersion.String(), gr.Data.GenesisForkVersion)
				require.Equal(t, testGenesisTime.Base10(), gr.Data.GenesisTime)
			},
		},
		{
			name: "genesis not ready",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(nil, 0, cometbft.ErrAppNotReady)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.ErrorIs(t, err, types.ErrNotFound)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// setup test
			backend := mocks.NewBackend(t)
			h := beacon.NewHandler(backend, cs, noop.NewLogger[beaconlog.Logger]())
			e := echo.New()
			e.Validator = &middleware.CustomValidator{
				Validator: middleware.ConstructValidator(),
			}

			// set expectations
			tc.setMockExpectations(backend)

			// test
			res, err := h.GetGenesis(nil) // input not relevant for this API

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

type backendKVStoreService struct {
	ctx sdk.Context
}

func (kvs *backendKVStoreService) OpenKVStore(context.Context) corestore.KVStore {
	//nolint:contextcheck // fine with tests
	store := sdk.UnwrapSDKContext(kvs.ctx).KVStore(testStoreKey)
	return kvstorage.NewKVStore(store)
}

var testStoreKey = storetypes.NewKVStoreKey("test-genesis")

func initTestGenesisState(t *testing.T, cs chain.Spec) *state.StateDB {
	t.Helper()

	db, err := db.OpenDB("", dbm.MemDBBackend)
	require.NoError(t, err)

	var (
		nopLog     = log.NewNopLogger()
		nopMetrics = sdkmetrics.NewNoOpMetrics()
	)

	cms := store.NewCommitMultiStore(db, nopLog, nopMetrics)
	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	require.NoError(t, cms.LoadLatestVersion())

	ctx := sdk.NewContext(cms, true, nopLog)
	backendStoreService := &backendKVStoreService{
		ctx: ctx,
	}
	kvStore := beacondb.New(backendStoreService)

	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, nopLog)
	return state.NewBeaconStateFromDB(
		kvStore.WithContext(sdkCtx),
		cs,
		sdkCtx.Logger(),
		metrics.NewNoOpTelemetrySink(),
	)
}
