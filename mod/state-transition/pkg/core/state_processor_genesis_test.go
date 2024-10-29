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

package core_test

import (
	"context"
	"fmt"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	statedb "github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/db"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	TestBeaconStateMarshallableT = types.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.BeaconBlockHeader,
		types.Eth1Data,
		types.ExecutionPayloadHeader,
		types.Fork,
		types.Validator,
	]

	TestKVStoreT = beacondb.KVStore[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.Validators,
	]

	TestBeaconStateT = statedb.StateDB[
		*types.BeaconBlockHeader,
		*TestBeaconStateMarshallableT,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*TestKVStoreT,
		*types.Validator,
		types.Validators,
		*engineprimitives.Withdrawal,
		types.WithdrawalCredentials,
	]
)

func TestInitialize(t *testing.T) {
	cs := spec.TestnetChainSpec()
	// in.ExecutionEngine,
	mocksSigner := &mocks.BLSSigner{}

	sp := core.NewStateProcessor[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		*TestBeaconStateT,
		*transition.Context,
		*types.Deposit,
		*types.Eth1Data,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*TestKVStoreT,
		*types.Validator,
		types.Validators,
		*engineprimitives.Withdrawal,
		engineprimitives.Withdrawals,
		types.WithdrawalCredentials,
	](
		cs,
		nil,
		mocksSigner,
	)

	kvStore, err := initTestStore()
	require.NoError(t, err)

	var (
		beaconState = new(TestBeaconStateT).NewFromDB(kvStore, cs)
		deposits    = []*types.Deposit{
			{
				Pubkey: [48]byte{0x01},
				Amount: math.Gwei(1_000),
				Index:  uint64(0),
			},
			{
				Pubkey: [48]byte{0x02},
				Amount: math.Gwei(2_000),
				Index:  uint64(1),
			},
		}
		executionPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genesisVersion         = version.FromUint32[common.Version](version.Deneb)
	)

	mocksSigner.On(
		"VerifySignature",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	vals, err := sp.InitializePreminedBeaconStateFromEth1(
		beaconState,
		deposits,
		executionPayloadHeader,
		genesisVersion,
	)
	require.NoError(t, err)
	require.Len(t, vals, len(deposits))
}

// Unit tests helpers

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
	testStoreKey = storetypes.NewKVStoreKey("state-transition-tests")
	testCodec    = &encoding.SSZInterfaceCodec[*types.ExecutionPayloadHeader]{}
)

func initTestStore() (
	*beacondb.KVStore[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.Validators,
	], error) {
	db, err := db.OpenDB("", dbm.MemDBBackend)
	if err != nil {
		return nil, fmt.Errorf("failed opening mem db: %w", err)
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
		return nil, fmt.Errorf("failed to load latest version: %w", err)
	}
	testStoreService := &testKVStoreService{ctx: ctx}

	return beacondb.New[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.Validators,
	](
		testStoreService,
		testCodec,
	), nil
}
