//go:build test
// +build test

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

package core_test

import (
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

func TestPartialWithdrawalRequestGenesisValidators(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	// make sure Electra is active
	require.True(t, version.EqualsOrIsAfter(cs.GenesisForkVersion(), version.Electra()))

	var (
		maxBalance = cs.MaxEffectiveBalance()
		minBalance = cs.MinActivationBalance()

		addr    = common.ExecutionAddress{0x01}
		creds   = types.NewCredentialsFromExecutionAddress(addr)
		badAddr = common.ExecutionAddress{0x20}
	)

	// Add a single validator to which we will target withdrawal requests
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: creds,
				Amount:      maxBalance,
				Index:       0,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))

	// Send withdrawal requests and see the valid ones being carried out
	vals, err := st.GetValidators()
	require.NoError(t, err)
	require.Len(t, vals, 1)
	genVal := vals[0]
	genValIdx, err := st.ValidatorIndexByPubkey(genVal.GetPubkey())
	require.NoError(t, err)

	genValPubKey := genVal.GetPubkey()
	wrs := []*types.WithdrawalRequest{
		{ // valid request
			SourceAddress:   addr,
			ValidatorPubKey: genValPubKey,
			Amount:          1,
		},
		{ // valid request
			SourceAddress:   addr,
			ValidatorPubKey: genValPubKey,
			Amount:          10,
		},
		{ // invalid request, invalid address
			SourceAddress:   badAddr,
			ValidatorPubKey: genValPubKey,
			Amount:          10,
		},
		{ // invalid request, invalid pub key
			SourceAddress:   addr,
			ValidatorPubKey: crypto.BLSPubkey(append([]byte{0xff}, genValPubKey[1:]...)),
			Amount:          10,
		},
		{ // valid request, largest withdrawable amount
			SourceAddress:   addr,
			ValidatorPubKey: genValPubKey,
			Amount:          maxBalance - 1 - 10 - minBalance, // remaining amount to minBalance
		},
		{ // invalid request (can't go below min activation balance even by 1 bera)
			SourceAddress:   addr,
			ValidatorPubKey: genValPubKey,
			Amount:          1,
		},
		{ // invalid request (full withdraw ignored when partial withdraws are ongoing)
			SourceAddress:   addr,
			ValidatorPubKey: genValPubKey,
			Amount:          0,
		},
	}

	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), 0, uint64(len(genDeposits)))
	require.NoError(t, err)

	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(10),
	)

	// Run the test.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// check withdrawal request is enqueued
	expectedWithdrawalEpoch := 1 + cs.MinValidatorWithdrawabilityDelay()
	pr, err := st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Len(t, pr, 3)
	require.Equal(t,
		[]*types.PendingPartialWithdrawal{
			{
				ValidatorIndex:    0,
				Amount:            wrs[0].Amount,
				WithdrawableEpoch: expectedWithdrawalEpoch,
			},
			{
				ValidatorIndex:    0,
				Amount:            wrs[1].Amount,
				WithdrawableEpoch: expectedWithdrawalEpoch,
			},
			{
				ValidatorIndex:    0,
				Amount:            wrs[4].Amount,
				WithdrawableEpoch: expectedWithdrawalEpoch,
			},
		},
		pr,
	)

	// check that the request is eventually fulfilled
	for range expectedWithdrawalEpoch - 1 {
		_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

		// This is just because we cannot chain moveToEndOfEpoch
		// back to back. TODO: fix
		timestamp := blk.Body.ExecutionPayload.Timestamp + 1
		blk = buildNextBlock(
			t,
			cs,
			st,
			types.NewEth1Data(depRoot),
			10,
			[]*types.Deposit{},
			&types.ExecutionRequests{},
			st.EVMInflationWithdrawal(timestamp),
		)
		_, err = sp.Transition(ctx, st, blk)
		require.NoError(t, err)
	}
	for range cs.SlotsPerEpoch() - 1 {
		timestamp := blk.Body.ExecutionPayload.Timestamp + 1
		blk = buildNextBlock(
			t,
			cs,
			st,
			types.NewEth1Data(depRoot),
			10,
			[]*types.Deposit{},
			&types.ExecutionRequests{},
			st.EVMInflationWithdrawal(timestamp),
		)
		_, err = sp.Transition(ctx, st, blk)
		require.NoError(t, err)
	}

	// finally the withdrawal
	timestamp := blk.Body.ExecutionPayload.Timestamp + 1
	withdrawals := []*engineprimitives.Withdrawal{
		// The first withdrawal is always for EVM inflation.
		st.EVMInflationWithdrawal(10),
		{
			Index:     0,
			Validator: genValIdx,
			Amount:    wrs[0].Amount,
			Address:   wrs[0].SourceAddress,
		},
		{
			Index:     1,
			Validator: genValIdx,
			Amount:    wrs[1].Amount,
			Address:   wrs[1].SourceAddress,
		},
		{
			Index:     2,
			Validator: genValIdx,
			Amount:    wrs[4].Amount,
			Address:   wrs[4].SourceAddress,
		},
	}
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		timestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		withdrawals...,
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// no more pending withdrawals
	pr, err = st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Empty(t, pr)
}

func TestFullWithdrawalRequestGenesisValidators(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	// make sure Electra is active
	require.True(t, version.EqualsOrIsAfter(cs.GenesisForkVersion(), version.Electra()))

	var (
		maxBalance = cs.MaxEffectiveBalance()

		addr1  = common.ExecutionAddress{0x01}
		creds1 = types.NewCredentialsFromExecutionAddress(addr1)
		addr2  = common.ExecutionAddress{0x01}
		creds2 = types.NewCredentialsFromExecutionAddress(addr2)
	)

	// Add a couple of validators and fully withdraw one of them
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: creds1,
				Amount:      maxBalance,
				Index:       0,
			},
			{
				Pubkey:      [48]byte{0x01},
				Credentials: creds2,
				Amount:      maxBalance,
				Index:       1,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))

	vals, err := st.GetValidators()
	require.NoError(t, err)
	require.Len(t, vals, 2)
	valToRm := vals[0]

	wrs := []*types.WithdrawalRequest{
		{
			SourceAddress:   addr1,
			ValidatorPubKey: valToRm.GetPubkey(),
			Amount:          0,
		},
	}

	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), 0, uint64(len(genDeposits)))
	require.NoError(t, err)

	blkTimestamp := math.U64(10)
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Run the test.
	valDiff, err := sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	// check that valToRm has initiated exit
	expectedExitEpoch := math.Epoch(1)
	expectedWithdrawalEpoch := expectedExitEpoch + cs.MinValidatorWithdrawabilityDelay()

	valToRmIdx, err := st.ValidatorIndexByPubkey(valToRm.GetPubkey())
	require.NoError(t, err)
	valToRm, err = st.ValidatorByIndex(valToRmIdx)
	require.NoError(t, err)
	require.Equal(t, expectedExitEpoch, valToRm.ExitEpoch)
	require.Equal(t, expectedWithdrawalEpoch, valToRm.WithdrawableEpoch)

	// no pending withdrawals, full withdrawals are executed right away
	pr, err := st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Empty(t, pr)

	// check the validator duly exits validator set
	blk = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(blkTimestamp),
	)
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Equal(t,
		transition.ValidatorUpdates{
			{
				Pubkey:           valToRm.Pubkey,
				EffectiveBalance: 0,
			},
		},
		valDiff,
	)

	// no more partial withdrawals are possible for an exited validator
	wrs = []*types.WithdrawalRequest{
		{
			SourceAddress:   addr1,
			ValidatorPubKey: valToRm.GetPubkey(),
			Amount:          1,
		},
	}
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Run the test.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	pr, err = st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Empty(t, pr)

	// check that balance is still locked and it will be
	// returned after MinValidatorWithdrawabilityDelay epochs
	// check that the request is eventually fulfilled
	for range expectedWithdrawalEpoch - 2 {
		_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

		// This is just because we cannot chain moveToEndOfEpoch
		// back to back. TODO: fix
		blkTimestamp = blk.GetTimestamp() + 1
		blk = buildNextBlock(
			t,
			cs,
			st,
			types.NewEth1Data(depRoot),
			blkTimestamp,
			[]*types.Deposit{},
			&types.ExecutionRequests{},
			st.EVMInflationWithdrawal(blkTimestamp),
		)
		_, err = sp.Transition(ctx, st, blk)
		require.NoError(t, err)
	}
	for range cs.SlotsPerEpoch() - 1 {
		blkTimestamp = blk.GetTimestamp() + 1
		blk = buildNextBlock(
			t,
			cs,
			st,
			types.NewEth1Data(depRoot),
			blkTimestamp,
			[]*types.Deposit{},
			&types.ExecutionRequests{},
			st.EVMInflationWithdrawal(blkTimestamp),
		)
		_, err = sp.Transition(ctx, st, blk)
		require.NoError(t, err)
	}

	// finally the withdrawal
	timestamp := blk.Body.ExecutionPayload.Timestamp + 1
	withdrawals := []*engineprimitives.Withdrawal{
		// The first withdrawal is always for EVM inflation.
		st.EVMInflationWithdrawal(10),
		{
			Index:     0,
			Validator: valToRmIdx,
			Amount:    valToRm.EffectiveBalance, // wrs request has zero to signal full withdrawal
			Address:   wrs[0].SourceAddress,
		},
	}
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		timestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		withdrawals...,
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// Check that validator balance is 0
	valBalance, err := st.GetBalance(valToRmIdx)
	require.NoError(t, err)
	require.Equal(t, math.U64(0), valBalance)

	// Check that effective balance has not updated yet. It will update next epoch
	// as part of processEffectiveBalanceUpdates
	valToRm, err = st.ValidatorByIndex(valToRmIdx)
	require.NoError(t, err)
	require.Equal(t, maxBalance, valToRm.GetEffectiveBalance())

	// Move forward one more epoch to trigger the effective balance update
	{
		_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)
		// This is just because we cannot chain moveToEndOfEpoch
		// back to back. TODO: fix
		blkTimestamp = blk.GetTimestamp() + 1
		blk = buildNextBlock(
			t,
			cs,
			st,
			types.NewEth1Data(depRoot),
			blkTimestamp,
			[]*types.Deposit{},
			&types.ExecutionRequests{},
			st.EVMInflationWithdrawal(blkTimestamp),
		)
		_, err = sp.Transition(ctx, st, blk)
		require.NoError(t, err)
	}

	// Check that effective balance is now 0
	valToRm, err = st.ValidatorByIndex(valToRmIdx)
	require.NoError(t, err)
	require.Equal(t, math.Gwei(0), valToRm.GetEffectiveBalance())
}

func TestWithdrawalRequestsNonGenesisValidators(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	// make sure Electra is active
	require.True(t, version.EqualsOrIsAfter(cs.GenesisForkVersion(), version.Electra()))

	var (
		maxBalance = cs.MaxEffectiveBalance()

		genAddr  = common.ExecutionAddress{0x01}
		genCreds = types.NewCredentialsFromExecutionAddress(genAddr)
		valAddr  = common.ExecutionAddress{0x01}
		valCreds = types.NewCredentialsFromExecutionAddress(valAddr)
	)

	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: genCreds,
				Amount:      maxBalance,
				Index:       0,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))

	// add a validator and test withdrawals through its lifetime
	blkDeposit := &types.Deposit{
		Pubkey:      [48]byte{0xff},
		Credentials: valCreds,
		Amount:      maxBalance,
		Index:       uint64(len(genDeposits)),
	}
	blkDeposits := []*types.Deposit{blkDeposit}

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blkDeposits))
	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), uint64(len(genDeposits)), cs.MaxDepositsPerBlock())
	require.NoError(t, err)

	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		blkDeposits,
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)

	// run the test
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// assert that validator is not even eligible for activation yet
	idx, err := st.ValidatorIndexByPubkey(blkDeposit.Pubkey)
	require.NoError(t, err)
	val, err := st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, constants.FarFutureEpoch, val.ActivationEligibilityEpoch)

	// validator is not even in activation queue, any withdrawal request is dropped
	wrs := []*types.WithdrawalRequest{
		{
			SourceAddress:   valAddr,
			ValidatorPubKey: val.GetPubkey(),
			Amount:          1,
		},
	}
	blkTimestamp := blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Show that no withdrawal is enqueued
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	pr, err := st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Empty(t, pr)

	// try again with full withdrawal
	wrs = []*types.WithdrawalRequest{
		{
			SourceAddress:   valAddr,
			ValidatorPubKey: val.GetPubkey(),
			Amount:          0,
		},
	}
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Show that validator is not marked for exit
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	val, err = st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, constants.FarFutureEpoch, val.ActivationEligibilityEpoch)
	require.Equal(t, constants.FarFutureEpoch, val.ExitEpoch)

	// make validator eligible for activation and show that withdrawals are not yet allowed
	blk = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(blkTimestamp),
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	val, err = st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), val.ActivationEligibilityEpoch)
	require.Equal(t, constants.FarFutureEpoch, val.ActivationEpoch)

	// validator eligible for activation but not active yet. Requests dropped
	wrs = []*types.WithdrawalRequest{
		{
			SourceAddress:   valAddr,
			ValidatorPubKey: val.GetPubkey(),
			Amount:          1,
		},
	}
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Show that no withdrawal is enqueued
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	pr, err = st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Empty(t, pr)

	// try again with full withdrawal
	wrs = []*types.WithdrawalRequest{
		{
			SourceAddress:   valAddr,
			ValidatorPubKey: val.GetPubkey(),
			Amount:          0,
		},
	}
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Show that validator is not marked for exit
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	val, err = st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, constants.FarFutureEpoch, val.ActivationEpoch)
	require.Equal(t, constants.FarFutureEpoch, val.ExitEpoch)

	// finally when validator is active withdrawals will work
	blk = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	wrs = []*types.WithdrawalRequest{
		{
			SourceAddress:   valAddr,
			ValidatorPubKey: val.GetPubkey(),
			Amount:          1,
		},
	}
	blkTimestamp = blk.GetTimestamp() + 1
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Run the test.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	pr, err = st.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Len(t, pr, 1)
	require.Equal(t,
		[]*types.PendingPartialWithdrawal{
			{
				ValidatorIndex:    1,
				Amount:            wrs[0].Amount,
				WithdrawableEpoch: 3 + cs.MinValidatorWithdrawabilityDelay(),
			},
		},
		pr,
	)
}

// Check that if the withdrawal request comes for a validator about to
// be evicted, this double eviction is duly handled
func TestConcurrentAutomaticAndVoluntaryWithdrawalRequests(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	// make sure Electra is active
	require.True(t, version.EqualsOrIsAfter(cs.GenesisForkVersion(), version.Electra()))

	// Make sure we have as many validators as the cap allows
	var (
		maxBalance = cs.MaxEffectiveBalance()
		rndSeed    = 2024 // seed used to generate unique random value
	)

	var (
		genDeposits      = make(types.Deposits, 0, cs.ValidatorSetCap())
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)

	// Step1: let blockchain have as many validators as cap allows
	for idx := range cs.ValidatorSetCap() {
		var (
			key   bytes.B48
			creds types.WithdrawalCredentials
		)
		key, rndSeed = generateTestPK(t, rndSeed)
		creds, rndSeed = generateTestExecutionAddress(t, rndSeed)

		genDeposits = append(
			genDeposits,
			&types.Deposit{
				Pubkey:      key,
				Credentials: creds,
				Amount:      maxBalance,
				Index:       idx,
			},
		)
	}

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	_, err := sp.InitializeBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	// Step 2: add a deposit which will evict one of the existing validators
	newValKey, rndSeed := generateTestPK(t, rndSeed)
	newValCreds, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		newValDeposit = &types.Deposit{
			Pubkey:      newValKey,
			Credentials: newValCreds,
			Amount:      maxBalance,
			Index:       uint64(len(genDeposits)),
		}
		blkDeposits = []*types.Deposit{newValDeposit}
	)

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blkDeposits))
	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), uint64(len(genDeposits)), cs.MaxDepositsPerBlock())
	require.NoError(t, err)

	blkTimestamp := math.U64(10)
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		blkDeposits,
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// check the deposit has been accepted
	_, err = st.ValidatorIndexByPubkey(newValDeposit.Pubkey)
	require.NoError(t, err)

	// move chain on till the new validator is about to be activated
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	// right at the block where we a validator will be evicted
	// we add a full withdrawal request for it. We expect this request
	// to be simply dropped since the validator is evicted upon ProcessEpoch
	evictedValAddr, err := genDeposits[0].Credentials.ToExecutionAddress()
	require.NoError(t, err)
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: []*types.WithdrawalRequest{
				{
					SourceAddress:   evictedValAddr,
					ValidatorPubKey: genDeposits[0].Pubkey,
					Amount:          0,
				},
			},
		},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)
	valDiff, err := sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, valDiff, 2)
	require.Equal(
		t,
		&transition.ValidatorUpdate{
			Pubkey:           newValDeposit.Pubkey,
			EffectiveBalance: newValDeposit.Amount,
		},
		valDiff[0],
	)
	require.Equal(
		t,
		&transition.ValidatorUpdate{
			Pubkey:           genDeposits[0].Pubkey,
			EffectiveBalance: 0,
		},
		valDiff[1],
	)

	evictedValIdx, err := st.ValidatorIndexByPubkey(genDeposits[0].Pubkey)
	require.NoError(t, err)
	evictedVal, err := st.ValidatorByIndex(evictedValIdx)
	require.NoError(t, err)

	expectedExitEpoch := math.Epoch(2)
	expectedWithdrawalEpoch := math.Epoch(2) + cs.MinValidatorWithdrawabilityDelay()
	require.Equal(t, expectedExitEpoch, evictedVal.ExitEpoch)
	require.Equal(t, expectedWithdrawalEpoch, evictedVal.WithdrawableEpoch)
}

// Check that two full withdrawals requests issued back to back
// are idempotent. The second request simply dropped
func TestDoubleFullWithdrawalRequests(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	// make sure Electra is active
	require.True(t, version.EqualsOrIsAfter(cs.GenesisForkVersion(), version.Electra()))

	var (
		maxBalance = cs.MaxEffectiveBalance()

		addr1  = common.ExecutionAddress{0x01}
		creds1 = types.NewCredentialsFromExecutionAddress(addr1)
		addr2  = common.ExecutionAddress{0x01}
		creds2 = types.NewCredentialsFromExecutionAddress(addr2)
	)

	// Add a couple of validators and fully withdraw one of them
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: creds1,
				Amount:      maxBalance,
				Index:       0,
			},
			{
				Pubkey:      [48]byte{0x01},
				Credentials: creds2,
				Amount:      maxBalance,
				Index:       1,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))

	vals, err := st.GetValidators()
	require.NoError(t, err)
	require.Len(t, vals, 2)
	valToRm := vals[0]

	wrs := []*types.WithdrawalRequest{
		{
			SourceAddress:   addr1,
			ValidatorPubKey: valToRm.GetPubkey(),
			Amount:          0,
		},
	}

	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), 0, uint64(len(genDeposits)))
	require.NoError(t, err)

	blkTimestamp := math.U64(10)
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)

	// Run the test.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// check that valToRm has initiated exit
	expectedExitEpoch := math.Epoch(1)
	expectedWithdrawalEpoch := expectedExitEpoch + cs.MinValidatorWithdrawabilityDelay()

	valToRmIdx, err := st.ValidatorIndexByPubkey(valToRm.GetPubkey())
	require.NoError(t, err)
	valToRm, err = st.ValidatorByIndex(valToRmIdx)
	require.NoError(t, err)
	require.Equal(t, expectedExitEpoch, valToRm.ExitEpoch)
	require.Equal(t, expectedWithdrawalEpoch, valToRm.WithdrawableEpoch)

	// issue another full withdrawal request, which should be
	// processed without errors and simply dropped
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blkTimestamp,
		[]*types.Deposit{},
		&types.ExecutionRequests{
			Withdrawals: wrs,
		},
		st.EVMInflationWithdrawal(blkTimestamp),
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
}

func TestPartialWithdrawalsOfBalanceAboveMaxEffectiveBalance(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance   = cs.MaxEffectiveBalance()
		minBalance   = cs.EffectiveBalanceIncrement()
		credentials0 = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
		address1     = common.ExecutionAddress{0x01}
		credentials1 = types.NewCredentialsFromExecutionAddress(address1)
	)

	// Setup initial state so that validator 1 is partially withdrawable.
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: credentials0,
				Amount:      maxBalance - 3*minBalance,
				Index:       0,
			},
			{
				Pubkey:      [48]byte{0x01},
				Credentials: credentials1,
				Amount:      maxBalance + minBalance,
				Index:       1,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	// Assert validator 1 balance before withdrawal.
	val1Bal, err := st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val1Bal)

	// Create test inputs.
	withdrawals := []*engineprimitives.Withdrawal{
		// The first withdrawal is always for EVM inflation.
		st.EVMInflationWithdrawal(10),
		// Partially withdraw validator 1 by minBalance.
		{
			Index:     0,
			Validator: 1,
			Amount:    minBalance,
			Address:   address1,
		},
	}

	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), 0, uint64(len(genDeposits)))
	require.NoError(t, err)

	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		withdrawals...,
	)

	// Run the test.
	_, err = sp.Transition(ctx, st, blk)

	// Check outputs and ensure withdrawals in payload is consistent with
	// statedb expected withdrawals.
	require.NoError(t, err)

	// Assert validator 1 balance after withdrawal.
	val1BalAfter, err := st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance, val1BalAfter)
}

func TestTransitionMaxWithdrawals(t *testing.T) {
	t.Parallel()
	// Use custom chain spec with max withdrawals set to 2.
	csData := spec.DevnetChainSpecData()
	csData.MaxWithdrawalsPerPayload = 2
	csData.MaxValidatorsPerWithdrawalsSweep = 2
	cs, err := chain.NewSpec(csData)
	require.NoError(t, err)

	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance   = cs.MaxEffectiveBalance()
		minBalance   = cs.EffectiveBalanceIncrement()
		address0     = common.ExecutionAddress{}
		credentials0 = types.NewCredentialsFromExecutionAddress(address0)
		address1     = common.ExecutionAddress{0x01}
		credentials1 = types.NewCredentialsFromExecutionAddress(address1)
	)

	// Setup initial state so that both validators are partially withdrawable.
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: credentials0,
				Amount:      maxBalance + minBalance,
				Index:       0,
			},
			{
				Pubkey:      [48]byte{0x01},
				Credentials: credentials1,
				Amount:      maxBalance + minBalance,
				Index:       1,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	_, err = sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	// Assert validator balances before withdrawal.
	val0Bal, err := st.GetBalance(math.U64(0))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val0Bal)

	val1Bal, err := st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val1Bal)

	// Create test inputs.
	withdrawals := []*engineprimitives.Withdrawal{
		// The first withdrawal is always for EVM inflation.
		st.EVMInflationWithdrawal(10),
		// Partially withdraw validator 0 by minBalance.
		{
			Index:     0,
			Validator: 0,
			Amount:    minBalance,
			Address:   address0,
		},
	}

	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), 0, uint64(len(genDeposits)))
	require.NoError(t, err)

	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		withdrawals...,
	)

	// Run the test.
	_, err = sp.Transition(ctx, st, blk)

	// Check outputs and ensure withdrawals in payload is consistent with
	// statedb expected withdrawals.
	require.NoError(t, err)

	// Assert validator balances after withdrawal, ensuring only validator 0 is
	// withdrawn from.
	val0BalAfter, err := st.GetBalance(math.U64(0))
	require.NoError(t, err)
	require.Equal(t, maxBalance, val0BalAfter)

	val1BalAfter, err := st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val1BalAfter)

	// Process the next block, ensuring that validator 1 is also withdrawn from,
	// also ensuring that the state's next withdrawal (validator) index is
	// appropriately incremented.

	withdrawals = []*engineprimitives.Withdrawal{
		// The first withdrawal is always for EVM inflation.
		st.EVMInflationWithdrawal(blk.GetTimestamp() + 1),
		// Partially withdraw validator 1 by minBalance.
		{
			Index:     1,
			Validator: 1,
			Amount:    minBalance,
			Address:   address1,
		},
	}
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		withdrawals...,
	)
	// Run the test.
	vals, err := sp.Transition(ctx, st, blk)

	// Check outputs and ensure withdrawals in payload is consistent with
	// statedb expected withdrawals.
	require.NoError(t, err)
	require.Zero(t, vals)

	// Validator 1 is now withdrawn from.
	val1BalAfter, err = st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance, val1BalAfter)
}

func TestValidatorNotWithdrawable(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		belowActiveBalance = cs.MinActivationBalance() - cs.EffectiveBalanceIncrement()
		maxBalance         = cs.MaxEffectiveBalance()
		validCredentials   = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup initial state with one validator
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: validCredentials,
				Amount:      maxBalance,
				Index:       0,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	// Create the block deposit with a non-ETH1 withdrawal credentials. This stake should not
	// be lost.
	invalidCredentials := types.WithdrawalCredentials(validCredentials[:])
	invalidCredentials[1] = 0x01
	blkDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: invalidCredentials,
		Amount:      belowActiveBalance,
		Index:       1,
	}
	blkDeposits := []*types.Deposit{blkDeposit}

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blkDeposits))
	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), uint64(len(genDeposits)), cs.MaxDepositsPerBlock())
	require.NoError(t, err)

	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		blkDeposits,
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)

	// Run transition.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// Check that validator 0x01 is part of beacon state with below active balance.
	validator, err := st.ValidatorByIndex(1)
	require.NoError(t, err)
	require.Equal(t, belowActiveBalance, validator.EffectiveBalance)
}
