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

package backend_test

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/backend"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	types "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	nodemetrics "github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	cryptomocks "github.com/berachain/beacon-kit/primitives/crypto/mocks"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/state-transition/core/mocks"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/beacondb"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestFilteredValidators(t *testing.T) {
	t.Parallel()

	// Build backend to test
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	cms, kvStore, depositStore, err := statetransition.BuildTestStores()
	require.NoError(t, err)
	sb := storage.NewBackend(cs, nil, kvStore, depositStore, nil)

	sp := core.NewStateProcessor(
		noop.NewLogger[any](),
		cs,
		mocks.NewExecutionEngine(t),
		depositStore,
		&cryptomocks.BLSSigner{},
		func(bytes.B48) ([]byte, error) { return nil, nil },
		nodemetrics.NewNoOpTelemetrySink(),
	)

	b := backend.New(sb, cs, sp)
	tcs := &testConsensusService{
		cms:     cms,
		kvStore: kvStore,
		cs:      cs,
	}
	b.AttachQueryBackend(tcs)

	// refSlot to allow validators in multiple states
	// from initializing to withdrawned
	refSlot := math.Slot(cs.SlotsPerEpoch() * 3)

	// Set relevant quantities in initial status and
	// write them to make them available to caches built
	// on top of cms
	stateValidators := []*types.ValidatorData{
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   0,
				Balance: cs.MaxEffectiveBalance(),
			},
			Status: constants.ValidatorStatusPendingInitialized,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x01},
					WithdrawalCredentials:      [32]byte{0x02},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance()),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(constants.FarFutureEpoch),
					ActivationEpoch:            math.Epoch(constants.FarFutureEpoch),
					ExitEpoch:                  math.Epoch(constants.FarFutureEpoch),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   1,
				Balance: cs.MaxEffectiveBalance() * 3 / 4,
			},
			Status: constants.ValidatorStatusPendingQueued,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x03},
					WithdrawalCredentials:      [32]byte{0x04},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 2),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(constants.FarFutureEpoch),
					ExitEpoch:                  math.Epoch(constants.FarFutureEpoch),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   2,
				Balance: cs.MaxEffectiveBalance() / 4,
			},
			Status: constants.ValidatorStatusActiveOngoing,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x05},
					WithdrawalCredentials:      [32]byte{0x06},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 3),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(constants.FarFutureEpoch),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   3,
				Balance: cs.MaxEffectiveBalance() / 4,
			},
			Status: constants.ValidatorStatusActiveSlashed,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x15},
					WithdrawalCredentials:      [32]byte{0x16},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 3),
					Slashed:                    true,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(5),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   4,
				Balance: cs.MaxEffectiveBalance() / 4,
			},
			Status: constants.ValidatorStatusActiveExiting,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x17},
					WithdrawalCredentials:      [32]byte{0x18},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 3),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(5),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   5,
				Balance: cs.MaxEffectiveBalance() / 2,
			},
			Status: constants.ValidatorStatusExitedUnslashed,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x07},
					WithdrawalCredentials:      [32]byte{0x08},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 4),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   6,
				Balance: cs.MaxEffectiveBalance() / 2,
			},
			Status: constants.ValidatorStatusExitedSlashed,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x27},
					WithdrawalCredentials:      [32]byte{0x28},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 4),
					Slashed:                    true,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   7,
				Balance: cs.EjectionBalance(),
			},
			Status: constants.ValidatorStatusWithdrawalPossible,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x09},
					WithdrawalCredentials:      [32]byte{0x10},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 5),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          math.Epoch(0),
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   8,
				Balance: 0,
			},
			Status: constants.ValidatorStatusWithdrawalPossible,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x39},
					WithdrawalCredentials:      [32]byte{0x40},
					EffectiveBalance:           math.Gwei(cs.MaxEffectiveBalance() / 5),
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          math.Epoch(0),
				},
			),
		},
	}
	setupTestFilteredValidatorsState(
		t,
		cms, kvStore, cs,
		stateValidators,
	)

	// test cases
	tests := []struct {
		name        string
		inputsF     func() ([]string /*ids*/, []string /*statuses*/)
		expectedErr error
		checkF      func(t *testing.T, res []*types.ValidatorData)
	}{
		{
			name: "all validators",
			inputsF: func() ([]string, []string) {
				return nil, nil
			},
			expectedErr: nil,
			checkF: func(t *testing.T, res []*types.ValidatorData) {
				require.Len(t, res, len(stateValidators))
				for i := range len(res) {
					require.Equal(t, stateValidators[i], res[i], fmt.Sprintf("index %d", i))
				}
			},
		},
		{
			name: "some validators by indexes",
			inputsF: func() ([]string, []string) {
				return []string{"1", "3"}, nil
			},
			expectedErr: nil,
			checkF: func(t *testing.T, res []*types.ValidatorData) {
				expectedRes := []*types.ValidatorData{
					stateValidators[1],
					stateValidators[3],
				}
				require.Len(t, res, len(expectedRes))
				for i := range len(res) {
					require.Equal(t, expectedRes[i], res[i], fmt.Sprintf("index %d", i))
				}
			},
		},
		{
			name: "some validators by status",
			inputsF: func() ([]string, []string) {
				return nil, []string{
					constants.ValidatorStatusActiveOngoing,
					constants.ValidatorStatusActiveSlashed,
					constants.ValidatorStatusActiveExiting,
				}
			},
			expectedErr: nil,
			checkF: func(t *testing.T, res []*types.ValidatorData) {
				expectedRes := []*types.ValidatorData{
					stateValidators[2],
					stateValidators[3],
					stateValidators[4],
				}
				require.Len(t, res, len(expectedRes))
				for i := range len(res) {
					require.Equal(t, expectedRes[i], res[i], fmt.Sprintf("index %d", i))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, statuses := tt.inputsF()
			res, err := b.FilteredValidators(refSlot, ids, statuses)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.expectedErr)
			}
			tt.checkF(t, res)
		})
	}
}

func setupTestFilteredValidatorsState(
	t *testing.T,
	cms storetypes.CommitMultiStore,
	kvStore *beacondb.KVStore,
	cs chain.Spec,
	stateValidators []*types.ValidatorData,
) {
	t.Helper()

	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())
	st := statedb.NewBeaconStateFromDB(kvStore.WithContext(sdkCtx), cs)

	for _, in := range stateValidators {
		val, err := beacontypes.ValidatorToConsensus(in.Validator)
		require.NoError(t, err)
		require.NoError(t, st.AddValidator(val))

		require.NoError(t, st.SetBalance(math.ValidatorIndex(in.Index), math.Gwei(in.Balance)))
	}

	setupStateDummyParts(t, cs, st)

	// finally write it all to underlying cms
	sdkCtx.MultiStore().(storetypes.CacheMultiStore).Write()
}

func setupStateDummyParts(t *testing.T, cs chain.Spec, st *statedb.StateDB) {
	dummySlot := math.Slot(2025)
	require.NoError(t, st.SetSlot(dummySlot))

	fork := ctypes.NewFork(version.Genesis(), version.Genesis(), constants.GenesisEpoch)
	require.NoError(t, st.SetFork(fork))

	require.NoError(t, st.SetGenesisValidatorsRoot(common.Root{}))

	blkHeader := &ctypes.BeaconBlockHeader{
		Slot:            constants.GenesisSlot,
		ProposerIndex:   0,
		ParentBlockRoot: common.Root{},
		StateRoot:       common.Root{},
		BodyRoot:        common.Root{},
	}
	require.NoError(t, st.SetLatestBlockHeader(blkHeader))

	for i := range cs.HistoricalRootsLimit() {
		require.NoError(t, st.UpdateBlockRootAtIndex(i, common.Root{}))
		require.NoError(t, st.UpdateStateRootAtIndex(i, common.Root{}))
	}

	payload, err := ctypes.DefaultGenesisExecutionPayloadHeader()
	require.NoError(t, err)
	require.NoError(t, st.SetLatestExecutionPayloadHeader(payload))

	eth1Data := &ctypes.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 0,
		BlockHash:    payload.GetBlockHash(),
	}
	require.NoError(t, st.SetEth1Data(eth1Data))
	require.NoError(t, st.SetEth1DepositIndex(constants.FirstDepositIndex))

	for i := range cs.EpochsPerHistoricalVector() {
		require.NoError(t, st.UpdateRandaoMixAtIndex(
			i,
			common.Bytes32(payload.GetBlockHash()),
		))
	}

	require.NoError(t, st.SetNextWithdrawalIndex(0))
	require.NoError(t, st.SetNextWithdrawalValidatorIndex(0))
	require.NoError(t, st.SetTotalSlashing(0))
}

var errTestMemberNotImplemented = errors.New("not implemented")

// testConsensusService stubs consensus service
type testConsensusService struct {
	cms     storetypes.CommitMultiStore
	kvStore *beacondb.KVStore
	cs      chain.Spec
}

func (t *testConsensusService) CreateQueryContext(height int64, prove bool) (sdk.Context, error) {
	sdkCtx := sdk.NewContext(t.cms.CacheMultiStore(), true, log.NewNopLogger())

	// there validations mimics consensus service, not sure if they are necessary
	tmpState := statedb.NewBeaconStateFromDB(t.kvStore.WithContext(sdkCtx), t.cs)
	slot, err := tmpState.GetSlot()
	if err != nil {
		return sdk.Context{}, sdkerrors.ErrInvalidHeight
	}
	if height > int64(slot.Unwrap()) {
		return sdk.Context{}, sdkerrors.ErrInvalidHeight
	}
	// end of possibly unnecessary validations

	return sdkCtx, nil
}

func (t *testConsensusService) Start(ctx context.Context) error {
	return errTestMemberNotImplemented
}

func (t *testConsensusService) Stop() error {
	return errTestMemberNotImplemented
}

func (t *testConsensusService) Name() string {
	panic(errTestMemberNotImplemented)
}

func (t *testConsensusService) LastBlockHeight() int64 {
	panic(errTestMemberNotImplemented)
}
