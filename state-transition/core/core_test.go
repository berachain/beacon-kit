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
	"strconv"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-core/components"
	nodemetrics "github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	cryptomocks "github.com/berachain/beacon-kit/primitives/crypto/mocks"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/state-transition/core/mocks"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	depositstore "github.com/berachain/beacon-kit/storage/deposit"
	"github.com/berachain/beacon-kit/storage/encoding"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	TestBeaconStateMarshallableT = types.BeaconState

	TestKVStoreT = beacondb.KVStore

	TestBeaconStateT = statedb.StateDB

	TestStateProcessorT = core.StateProcessor[
		*transition.Context,
		*TestKVStoreT,
	]
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
	testStoreKey = storetypes.NewKVStoreKey("state-transition-tests")
	testCodec    = &encoding.SSZInterfaceCodec[*types.ExecutionPayloadHeader]{}
)

func initTestStores() (*beacondb.KVStore, *depositstore.KVStore, error) {
	db, err := db.OpenDB("", dbm.MemDBBackend)
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening mem db: %w", err)
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

	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	if err = cms.LoadLatestVersion(); err != nil {
		return nil, nil, fmt.Errorf("failed to load latest version: %w", err)
	}

	ctx := sdk.NewContext(cms, true, nopLog)
	testStoreService := &testKVStoreService{ctx: ctx}
	return beacondb.New(
			testStoreService,
			testCodec,
		),
		depositstore.NewStore(testStoreService, nopLog),
		nil
}

func setupChain(t *testing.T, chainSpecType string) chain.Spec[
	bytes.B4, math.U64, math.U64, any,
] {
	t.Helper()

	t.Setenv(components.ChainSpecTypeEnvVar, chainSpecType)
	cs, err := components.ProvideChainSpec()
	require.NoError(t, err)

	return cs
}

func setupState(
	t *testing.T, cs chain.Spec[
		bytes.B4, math.U64, math.U64, any,
	],
) (
	*TestStateProcessorT,
	*TestBeaconStateT,
	*depositstore.KVStore,
	*transition.Context,
) {
	t.Helper()

	execEngine := mocks.NewExecutionEngine(t)

	mocksSigner := &cryptomocks.BLSSigner{}
	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	dummyProposerAddr := []byte{0xff}

	kvStore, depositStore, err := initTestStores()
	require.NoError(t, err)
	beaconState := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	sp := core.NewStateProcessor[
		*transition.Context,
		*TestKVStoreT,
	](
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

	ctx := &transition.Context{
		SkipPayloadVerification: true,
		SkipValidateResult:      true,
		ProposerAddress:         dummyProposerAddr,
	}

	return sp, beaconState, depositStore, ctx
}

func progressStateToSlot(
	t *testing.T,
	beaconState *TestBeaconStateT,
	slot math.U64,
) {
	t.Helper()

	if slot == math.U64(0) {
		t.Fatal("for genesis slot, use InitializePreminedBeaconStateFromEth1")
	}

	err := beaconState.SetSlot(slot)
	require.NoError(t, err)
	err = beaconState.SetLatestBlockHeader(types.NewBeaconBlockHeader(
		slot,
		math.U64(0),
		common.Root{},
		common.Root{},
		common.Root{},
	))
	require.NoError(t, err)
}

func buildNextBlock(
	t *testing.T,
	beaconState *TestBeaconStateT,
	nextBlkBody *types.BeaconBlockBody,
) *types.BeaconBlock {
	t.Helper()

	// first update state root, similarly to what we do in processSlot
	parentBlkHeader, err := beaconState.GetLatestBlockHeader()
	require.NoError(t, err)
	root := beaconState.HashTreeRoot()
	parentBlkHeader.SetStateRoot(root)

	// finally build the block
	return &types.BeaconBlock{
		Slot:          parentBlkHeader.GetSlot() + 1,
		ProposerIndex: parentBlkHeader.GetProposerIndex(),
		ParentRoot:    parentBlkHeader.HashTreeRoot(),
		StateRoot:     common.Root{},
		Body:          nextBlkBody,
	}
}

func generateTestExecutionAddress(
	t *testing.T,
	rndSeed int,
) (types.WithdrawalCredentials, int) {
	t.Helper()

	addrStr := strconv.Itoa(rndSeed)
	addrBytes := bytes.ExtendToSize([]byte(addrStr), bytes.B20Size)
	execAddr, err := bytes.ToBytes20(addrBytes)
	require.NoError(t, err)
	rndSeed++
	return types.NewCredentialsFromExecutionAddress(
		common.ExecutionAddress(execAddr),
	), rndSeed
}

func generateTestPK(t *testing.T, rndSeed int) (bytes.B48, int) {
	t.Helper()
	keyStr := strconv.Itoa(rndSeed)
	keyBytes := bytes.ExtendToSize([]byte(keyStr), bytes.B48Size)
	key, err := bytes.ToBytes48(keyBytes)
	require.NoError(t, err)
	rndSeed++
	return key, rndSeed
}

func moveToEndOfEpoch(
	t *testing.T,
	tip *types.BeaconBlock,
	cs chain.Spec[bytes.B4, math.U64, math.U64, any],
	sp *TestStateProcessorT,
	st *TestBeaconStateT,
	ctx *transition.Context,
	depRoot common.Root,
) *types.BeaconBlock {
	t.Helper()
	blk := tip
	currEpoch := cs.SlotToEpoch(blk.GetSlot())
	for currEpoch == cs.SlotToEpoch(blk.GetSlot()+1) {
		blk = buildNextBlock(
			t,
			st,
			&types.BeaconBlockBody{
				ExecutionPayload: &types.ExecutionPayload{
					Timestamp:    blk.Body.ExecutionPayload.Timestamp + 1,
					ExtraData:    []byte("testing"),
					Transactions: [][]byte{},
					Withdrawals: []*engineprimitives.Withdrawal{
						st.EVMInflationWithdrawal(),
					},
					BaseFeePerGas: math.NewU256(0),
				},
				Eth1Data: &types.Eth1Data{DepositRoot: depRoot},
				Deposits: []*types.Deposit{},
			},
		)

		vals, err := sp.Transition(ctx, st, blk)
		require.NoError(t, err)
		require.Empty(t, vals) // no vals changes expected before next epoch
	}
	return blk
}
