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
	"testing"

	cosmoslog "cosmossdk.io/log"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx // multiple test cases
func TestFilterValidators(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	// Create some input validators and store them to a readonly state
	stateValidators := []*beacontypes.ValidatorData{
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   0,
				Balance: cs.MaxEffectiveBalance().Unwrap(),
			},
			Status: constants.ValidatorStatusPendingInitialized,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x01},
					WithdrawalCredentials:      [32]byte{0x02},
					EffectiveBalance:           cs.MaxEffectiveBalance(),
					Slashed:                    false,
					ActivationEligibilityEpoch: constants.FarFutureEpoch,
					ActivationEpoch:            constants.FarFutureEpoch,
					ExitEpoch:                  constants.FarFutureEpoch,
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   1,
				Balance: cs.MaxEffectiveBalance().Unwrap() * 3 / 4,
			},
			Status: constants.ValidatorStatusPendingQueued,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x03},
					WithdrawalCredentials:      [32]byte{0x04},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 2,
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            constants.FarFutureEpoch,
					ExitEpoch:                  constants.FarFutureEpoch,
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   2,
				Balance: cs.MaxEffectiveBalance().Unwrap() / 4,
			},
			Status: constants.ValidatorStatusActiveOngoing,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x05},
					WithdrawalCredentials:      [32]byte{0x06},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 3,
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  constants.FarFutureEpoch,
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   3,
				Balance: cs.MaxEffectiveBalance().Unwrap() / 4,
			},
			Status: constants.ValidatorStatusActiveSlashed,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x15},
					WithdrawalCredentials:      [32]byte{0x16},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 3,
					Slashed:                    true,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(5),
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   4,
				Balance: cs.MaxEffectiveBalance().Unwrap() / 4,
			},
			Status: constants.ValidatorStatusActiveExiting,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x17},
					WithdrawalCredentials:      [32]byte{0x18},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 3,
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(5),
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   5,
				Balance: cs.MaxEffectiveBalance().Unwrap() / 2,
			},
			Status: constants.ValidatorStatusExitedUnslashed,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x07},
					WithdrawalCredentials:      [32]byte{0x08},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 4,
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   6,
				Balance: cs.MaxEffectiveBalance().Unwrap() / 2,
			},
			Status: constants.ValidatorStatusExitedSlashed,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x27},
					WithdrawalCredentials:      [32]byte{0x28},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 4,
					Slashed:                    true,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          constants.FarFutureEpoch,
				},
			),
		},
		{
			ValidatorBalanceData: beacontypes.ValidatorBalanceData{
				Index:   7,
				Balance: cs.MinActivationBalance().Unwrap() - cs.EffectiveBalanceIncrement().Unwrap(),
			},
			Status: constants.ValidatorStatusWithdrawalPossible,
			Validator: beacontypes.ValidatorFromConsensus(
				&ctypes.Validator{
					Pubkey:                     [48]byte{0x09},
					WithdrawalCredentials:      [32]byte{0x10},
					EffectiveBalance:           cs.MaxEffectiveBalance() / 5,
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
					EffectiveBalance:           cs.MaxEffectiveBalance() / 5,
					Slashed:                    false,
					ActivationEligibilityEpoch: math.Epoch(0),
					ActivationEpoch:            math.Epoch(0),
					ExitEpoch:                  math.Epoch(0),
					WithdrawableEpoch:          math.Epoch(0),
				},
			),
		},
	}

	testCases := []struct {
		name                string
		inputs              func() ([]string, []string)
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res []*beacontypes.ValidatorData, err error)
	}{
		{
			name: "all validators",
			inputs: func() ([]string, []string) {
				return nil, nil
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				require.Len(t, res, len(stateValidators))
				for i := range res {
					require.Equal(t, stateValidators[i], res[i], "index %d", i)
				}
			},
		},
		{
			name: "some validators by indexes",
			inputs: func() ([]string, []string) {
				return []string{"1", "3"}, nil
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				expectedRes := []*beacontypes.ValidatorData{
					stateValidators[1],
					stateValidators[3],
				}
				require.Len(t, res, len(expectedRes))
				for i := range res {
					require.Equal(t, expectedRes[i], res[i], "index %d", i)
				}
			},
		},
		{
			name: "some validators by pub keys",
			inputs: func() ([]string, []string) {
				return []string{
					stateValidators[2].Validator.PublicKey,
					stateValidators[4].Validator.PublicKey,
				}, nil
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				expectedRes := []*beacontypes.ValidatorData{
					stateValidators[2],
					stateValidators[4],
				}
				require.Len(t, res, len(expectedRes))
				for i := range res {
					require.Equal(t, expectedRes[i], res[i], "index %d", i)
				}
			},
		},
		{
			name: "some validators by status",
			inputs: func() ([]string, []string) {
				return nil, []string{
					constants.ValidatorStatusActiveOngoing,
					constants.ValidatorStatusActiveSlashed,
					constants.ValidatorStatusActiveExiting,
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				expectedRes := []*beacontypes.ValidatorData{
					stateValidators[2],
					stateValidators[3],
					stateValidators[4],
				}
				require.Len(t, res, len(expectedRes))
				for i := range res {
					require.Equal(t, expectedRes[i], res[i], "index %d", i)
				}
			},
		},
		{
			name: "chain not ready",
			inputs: func() ([]string, []string) {
				return nil, nil
			},
			setMockExpectations: func(b *mocks.Backend) {
				// cometbft.ErrAppNotReady is the error flag returned when
				// genesis has not yet been processed and chain is not ready.
				b.EXPECT().StateAtSlot(mock.Anything).Return(nil, math.Slot(0), cometbft.ErrAppNotReady)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorData, err error) {
				t.Helper()

				// handlertypes.ErrNotFound is the error flag used to return 404 error code
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
				require.Nil(t, res)
			},
		},
		{
			name: "height requested too high",
			inputs: func() ([]string, []string) {
				return nil, nil
			},
			setMockExpectations: func(b *mocks.Backend) {
				// sdkerrors.ErrInvalidHeight is the error flag returned when
				// requested height is not in the state.
				b.EXPECT().StateAtSlot(mock.Anything).Return(nil, math.Slot(0), sdkerrors.ErrInvalidHeight)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorData, err error) {
				t.Helper()

				// handlertypes.ErrNotFound is the error flag used to return 404 error code
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// setup test
			backend := mocks.NewBackend(t)
			h := beacon.NewHandler(backend, cs, noop.NewLogger[log.Logger]())

			// set expectations
			tc.setMockExpectations(backend)

			// test
			ids, statuses := tc.inputs()
			res, err := h.FilterValidators(math.Slot(0), ids, statuses)

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

func addTestValidators(t *testing.T, stateValidators []*beacontypes.ValidatorData, st *statedb.StateDB) {
	t.Helper()

	for _, in := range stateValidators {
		val, errVal := beacontypes.ValidatorToConsensus(in.Validator)
		require.NoError(t, errVal)
		require.NoError(t, st.AddValidator(val))

		require.NoError(t, st.SetBalance(math.ValidatorIndex(in.Index), math.Gwei(in.Balance)))
	}
}

func makeTestState(t *testing.T, cs chain.Spec) *statedb.StateDB {
	t.Helper()

	cms, kvStore, _, errSt := statetransition.BuildTestStores()
	require.NoError(t, errSt)
	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, cosmoslog.NewNopLogger())
	st := statedb.NewBeaconStateFromDB(
		kvStore.WithContext(sdkCtx), cs, sdkCtx.Logger(), metrics.NewNoOpTelemetrySink(),
	)
	return st
}
