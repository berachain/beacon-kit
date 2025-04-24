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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // uses envars
func TestWithdrawalRequestLifecycle(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	// make sure Electra is active
	require.Equal(t, version.Electra(), cs.GenesisForkVersion())

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		minBalance = math.Gwei(cs.MinActivationBalance())

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

	// Send a small withdrawal request and see it being carried out
	vals, err := st.GetValidators()
	require.NoError(t, err)
	require.Len(t, vals, 1)
	genVal := vals[0]
	genValIdx, err := st.ValidatorIndexByPubkey(genVal.GetPubkey())
	require.NoError(t, err)

	wrs := []*types.WithdrawalRequest{
		{ // valid request
			SourceAddress:   addr,
			ValidatorPubKey: genVal.GetPubkey(),
			Amount:          1,
		},
		{ // valid request
			SourceAddress:   addr,
			ValidatorPubKey: genVal.GetPubkey(),
			Amount:          10,
		},
		{ // invalid request, invalid address
			SourceAddress:   badAddr,
			ValidatorPubKey: genVal.GetPubkey(),
			Amount:          10,
		},
		{ // valid request, largest withdrawable amount
			SourceAddress:   addr,
			ValidatorPubKey: genVal.GetPubkey(),
			Amount:          maxBalance - 1 - 10 - minBalance,
		},
		{ // invalid request (can't go below min activation balance)
			SourceAddress:   addr,
			ValidatorPubKey: genVal.GetPubkey(),
			Amount:          minBalance,
		},
		{ // invalid request (full withdraw ignored when partial withdraws are ongoing)
			SourceAddress:   addr,
			ValidatorPubKey: genVal.GetPubkey(),
			Amount:          0,
		},
	}

	depRoot := genDeposits.HashTreeRoot()
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
	expectedWithdrawalEpoch := math.Epoch(33)
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
				Amount:            wrs[3].Amount,
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
			Amount:    wrs[3].Amount,
			Address:   wrs[3].SourceAddress,
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

//nolint:paralleltest // uses envars
func TestTransitionWithdrawals(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		minBalance   = math.Gwei(cs.EffectiveBalanceIncrement())
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
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(genDeposits.HashTreeRoot()),
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
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		minBalance   = math.Gwei(cs.EffectiveBalanceIncrement())
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
	depRoot := genDeposits.HashTreeRoot()
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

//nolint:paralleltest // uses envars
func TestValidatorNotWithdrawable(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		belowActiveBalance = math.Gwei(cs.MinActivationBalance() - cs.EffectiveBalanceIncrement())
		maxBalance         = math.Gwei(cs.MaxEffectiveBalance())
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
	blockDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: invalidCredentials,
			Amount:      belowActiveBalance,
			Index:       1,
		},
	}

	depRoot := append(genDeposits, blockDeposits...).HashTreeRoot()
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		blockDeposits,
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blockDeposits))

	// Run transition.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	// Check that validator 0x01 is part of beacon state with below active balance.
	validator, err := st.ValidatorByIndex(1)
	require.NoError(t, err)
	require.Equal(t, belowActiveBalance, validator.EffectiveBalance)
}
