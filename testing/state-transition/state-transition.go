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
//go:build test
// +build test

package statetransition

import (
	"context"
	"fmt"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log/noop"
	nodemetrics "github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/bytes"
	cryptomocks "github.com/berachain/beacon-kit/primitives/crypto/mocks"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/state-transition/core/mocks"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	depositstore "github.com/berachain/beacon-kit/storage/deposit"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	TestBeaconStateMarshallableT = types.BeaconState
	TestBeaconStateT             = statedb.StateDB
	TestStateProcessorT          = core.StateProcessor
)

type testKVStoreService struct{}

func (kvs *testKVStoreService) OpenKVStore(ctx context.Context) corestore.KVStore {
	store := sdk.UnwrapSDKContext(ctx).KVStore(testStoreKey)
	return storage.NewKVStore(store)
}

//nolint:gochecknoglobals // unexported and use only in tests
var testStoreKey = storetypes.NewKVStoreKey("state-transition-tests")

func BuildTestStores() (
	storetypes.CommitMultiStore,
	*beacondb.KVStore,
	*depositstore.KVStore,
	error) {
	db, err := db.OpenDB("", dbm.MemDBBackend)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed opening mem db: %w", err)
	}
	var (
		nopLog        = log.NewNopLogger()
		noopCloseFunc = func() error { return nil }
		nopMetrics    = metrics.NewNoOpMetrics()
	)

	cms := store.NewCommitMultiStore(
		db,
		nopLog,
		nopMetrics,
	)

	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	if err = cms.LoadLatestVersion(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load latest version: %w", err)
	}

	testStoreService := &testKVStoreService{}
	return cms,
		beacondb.New(testStoreService),
		depositstore.NewStore(testStoreService, noopCloseFunc, nopLog),
		nil
}

func SetupTestState(t *testing.T, cs chain.Spec) (
	*TestStateProcessorT,
	*TestBeaconStateT,
	*depositstore.KVStore,
	core.ReadOnlyContext,
) {
	t.Helper()

	execEngine := mocks.NewExecutionEngine(t)

	mocksSigner := &cryptomocks.BLSSigner{}
	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	dummyProposerAddr := []byte{0xff}

	cms, kvStore, depositStore, err := BuildTestStores()
	require.NoError(t, err)

	ms := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())
	beaconState := statedb.NewBeaconStateFromDB(kvStore.WithContext(ms), cs)

	sp := core.NewStateProcessor(
		noop.NewLogger[any](),
		cs,
		execEngine,
		depositStore,
		mocksSigner,
		func(bytes.B48) ([]byte, error) {
			return dummyProposerAddr, nil
		},
		nodemetrics.NewNoOpTelemetrySink(),
	)

	ctx := transition.NewTransitionCtx(
		ms,
		0, // time
		dummyProposerAddr,
	).
		WithVerifyPayload(false).
		WithVerifyRandao(false).
		WithVerifyResult(false).
		WithMeterGas(false).
		WithOptimisticEngine(true)

	return sp, beaconState, depositStore, ctx
}
