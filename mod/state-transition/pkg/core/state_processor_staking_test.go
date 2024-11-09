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
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	cryptomocks "github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestTransitionUpdateValidators shows that when validator is
// updated (increasing amount), corresponding balance is updated.
func TestTransitionUpdateValidators(t *testing.T) {
	// Create state processor to test
	cs := spec.BetnetChainSpec()
	execEngine := mocks.NewExecutionEngine[
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		engineprimitives.Withdrawals,
	](t)
	mocksSigner := &cryptomocks.BLSSigner{}
	dummyProposerAddr := []byte{0xff}

	sp := createStateProcessor(
		cs,
		execEngine,
		mocksSigner,
		func(bytes.B48) ([]byte, error) {
			return dummyProposerAddr, nil
		},
	)

	kvStore, err := initStore()
	require.NoError(t, err)
	beaconState := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	var (
		maxBalance       = math.Gwei(cs.MaxEffectiveBalance())
		increment        = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance       = math.Gwei(cs.EjectionBalance())
		emptyAddress     = common.ExecutionAddress{}
		emptyCredentials = types.NewCredentialsFromExecutionAddress(
			emptyAddress,
		)
	)

	// STEP 0: Setup initial state via genesis
	var (
		genDeposits = []*types.Deposit{
			{
				Pubkey:      [48]byte{0x01},
				Credentials: emptyCredentials,
				Amount:      minBalance + increment,
				Index:       uint64(0),
			},
			{
				Pubkey:      [48]byte{0x02},
				Credentials: emptyCredentials,
				Amount:      maxBalance - 6*increment,
				Index:       uint64(1),
			},
			{
				Pubkey:      [48]byte{0x03},
				Credentials: emptyCredentials,
				Amount:      maxBalance - 3*increment,
				Index:       uint64(2),
			},
		}
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	genVals, err := sp.InitializePreminedBeaconStateFromEth1(
		beaconState,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)
	require.Len(t, genVals, len(genDeposits))

	// STEP 1: top up a genesis validator balance
	var (
		ctx = &transition.Context{
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			ProposerAddress:         dummyProposerAddr,
		}
		blkDeposit = &types.Deposit{
			Pubkey:      genDeposits[2].Pubkey,
			Credentials: emptyCredentials,
			Amount:      2 * increment, // twice to account for hysteresis
			Index:       uint64(len(genDeposits)),
		}
	)

	blk1 := buildNextBlock(
		t,
		beaconState,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:     10,
				ExtraData:     []byte("testing"),
				Transactions:  [][]byte{},
				Withdrawals:   []*engineprimitives.Withdrawal{}, // no withdrawals
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{blkDeposit},
		},
	)

	// run the test
	updatedVals, err := sp.Transition(ctx, beaconState, blk1)
	require.NoError(t, err)
	require.Empty(t, updatedVals) // validators set updates only at epoch turn

	// check validator balances are duly updated, that is:
	// - balance is updated immediately
	// - effective balance is updated only at the epoch turn
	expectedBalance := genDeposits[2].Amount + blkDeposit.Amount
	expectedEffectiveBalance := genDeposits[2].Amount
	idx, err := beaconState.ValidatorIndexByPubkey(genDeposits[2].Pubkey)
	require.NoError(t, err)

	balance, err := beaconState.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err := beaconState.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)

	// check that validator index is still correct
	latestValIdx, err := beaconState.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)), latestValIdx)

	// STEP 2: check that effective balance is updated once next epoch arrives
	var blk = blk1
	for i := 1; i < int(cs.SlotsPerEpoch())-1; i++ {
		blk = buildNextBlock(
			t,
			beaconState,
			&types.BeaconBlockBody{
				ExecutionPayload: &types.ExecutionPayload{
					Timestamp:     blk.Body.ExecutionPayload.Timestamp + 1,
					ExtraData:     []byte("testing"),
					Transactions:  [][]byte{},
					Withdrawals:   []*engineprimitives.Withdrawal{},
					BaseFeePerGas: math.NewU256(0),
				},
				Eth1Data: &types.Eth1Data{},
				Deposits: []*types.Deposit{},
			},
		)

		updatedVals, err = sp.Transition(ctx, beaconState, blk)
		require.NoError(t, err)
		require.Empty(t, updatedVals) // validators set updates only at epoch
	}

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		beaconState,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:     blk.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:     []byte("testing"),
				Transactions:  [][]byte{},
				Withdrawals:   []*engineprimitives.Withdrawal{},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	newEpochVals, err := sp.Transition(ctx, beaconState, blk)
	require.NoError(t, err)
	require.Len(t, newEpochVals, len(genDeposits)) // just topped up one validator

	// Assuming genesis order is preserved here which is not necessary
	// TODO: remove this assumption

	// all genesis validators other than the last are unchanged
	for i := range len(genDeposits) - 1 {
		require.Equal(t, genVals[i], newEpochVals[i], fmt.Sprintf("idx: %d", i))
	}

	expectedBalance = genDeposits[2].Amount + blkDeposit.Amount
	expectedEffectiveBalance = expectedBalance

	balance, err = beaconState.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err = beaconState.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)
}
