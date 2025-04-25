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
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

// TestTransitionUpdateValidators shows that when validator is
// updated (increasing amount), corresponding balance is updated.
//
//nolint:paralleltest // uses envars
func TestTransitionUpdateValidators(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance       = math.Gwei(cs.MaxEffectiveBalance())
		increment        = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance       = math.Gwei(cs.EjectionBalance())
		emptyCredentials = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// STEP 0: Setup initial state via genesis
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: emptyCredentials,
				Amount:      minBalance + increment,
				Index:       uint64(0),
			},
			{
				Pubkey:      [48]byte{0x01},
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
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	valDiff, err := sp.InitializeBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.Len(t, valDiff, len(genDeposits))

	// STEP 1: top up a genesis validator balance
	blkDeposit := &types.Deposit{
		Pubkey:      genDeposits[2].Pubkey,
		Credentials: emptyCredentials,
		Amount:      2 * increment, // twice to account for hysteresis
		Index:       uint64(len(genDeposits)),
	}

	depRoot := append(genDeposits, blkDeposit).HashTreeRoot()
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{blkDeposit},
		st.EVMInflationWithdrawal(10),
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blk1.Body.Deposits))

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff) // validators set updates only at epoch turn

	// check validator balances are duly updated, that is:
	// - balance is updated immediately
	// - effective balance is updated only at the epoch turn
	expectedBalance := genDeposits[2].Amount + blkDeposit.Amount
	expectedEffectiveBalance := genDeposits[2].Amount
	idx, err := st.ValidatorIndexByPubkey(genDeposits[2].Pubkey)
	require.NoError(t, err)

	balance, err := st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err := st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)

	// check that validator index is still correct
	latestValIdx, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)+1), latestValIdx)

	// STEP 2: check that effective balance is updated once next epoch arrives
	blk := moveToEndOfEpoch(t, blk1, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, valDiff, 1) // just topped up one validator
	require.Equal(
		t,
		&transition.ValidatorUpdate{
			Pubkey:           blkDeposit.Pubkey,
			EffectiveBalance: expectedBalance,
		},
		valDiff[0],
	)
	expectedEffectiveBalance = expectedBalance

	balance, err = st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err = st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)
}

// TestTransitionCreateValidator shows the lifecycle
// of a validator creation.
//
//nolint:paralleltest // uses envars
func TestTransitionCreateValidator(t *testing.T) {
	// Create state processor to test
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance       = math.Gwei(cs.MaxEffectiveBalance())
		increment        = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance       = math.Gwei(cs.EjectionBalance())
		emptyAddress     = common.ExecutionAddress{}
		emptyCredentials = types.NewCredentialsFromExecutionAddress(emptyAddress)
	)

	// STEP 0: Setup initial state via genesis
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x01},
				Credentials: emptyCredentials,
				Amount:      minBalance + increment,
				Index:       uint64(0),
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	genVals, err := sp.InitializeBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.Len(t, genVals, len(genDeposits))

	// STEP 1: top up a genesis validator balance
	blkDeposit := &types.Deposit{
		Pubkey:      [48]byte{0xff}, // a new key for a new validator
		Credentials: emptyCredentials,
		Amount:      maxBalance,
		Index:       uint64(len(genDeposits)),
	}

	depRoot := append(genDeposits, blkDeposit).HashTreeRoot()
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{blkDeposit},
		st.EVMInflationWithdrawal(10),
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blk1.Body.Deposits))

	// run the test
	valDiff, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff) // validators set updates only at epoch turn

	// check validator balances are duly updated
	var (
		expectedBalance          = blkDeposit.Amount
		expectedEffectiveBalance = expectedBalance
	)
	idx, err := st.ValidatorIndexByPubkey(blkDeposit.Pubkey)
	require.NoError(t, err)

	balance, err := st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err := st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)

	// check that validator index is still correct
	latestValIdx, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)+1), latestValIdx)

	// STEP 2: move the chain to the next epoch and show that
	// the extra validator is eligible for activation
	blk := moveToEndOfEpoch(t, blk1, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff) // new validator is only eligible for activation

	extraValIdx, err := st.ValidatorIndexByPubkey(blkDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEpoch,
	)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	// STEP 3: move the chain to the next epoch and show that
	// the extra validator is activate
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, valDiff, 1)
	require.Equal(
		t,
		&transition.ValidatorUpdate{
			Pubkey:           blkDeposit.Pubkey,
			EffectiveBalance: expectedBalance,
		},
		valDiff[0],
	)

	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(2), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	expectedBalance = blkDeposit.Amount
	expectedEffectiveBalance = expectedBalance

	balance, err = st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err = st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)
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

// TestTransitionHittingValidatorsCap shows that the extra
// validator added when validators set is at cap gets never activated
// and its deposit is returned at after next epoch starts.
func TestTransitionHittingValidatorsCap_ExtraSmall(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance      = math.Gwei(cs.MaxEffectiveBalance())
		ejectionBalance = math.Gwei(cs.EjectionBalance())
		minBalance      = ejectionBalance + math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed         = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with GetValidatorSetCap validators
	// TODO: consider instead setting state artificially
	var (
		genDeposits      = make(types.Deposits, 0, cs.ValidatorSetCap())
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)

	// let genesis define all available validators
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

	// STEP 1: Try and add an extra validator
	extraValKey, rndSeed := generateTestPK(t, rndSeed)
	extraValCreds, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		extraValDeposit = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      minBalance,
			Index:       uint64(len(genDeposits)),
		}
	)

	depRoot := append(genDeposits, extraValDeposit).HashTreeRoot()
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{extraValDeposit},
		st.EVMInflationWithdrawal(10),
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blk1.Body.Deposits))

	// run the test
	valDiff, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	extraValIdx, err := st.ValidatorIndexByPubkey(extraValDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEligibilityEpoch,
	)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEpoch,
	)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	// STEP 2: move the chain to the next epoch and show that
	// the extra validator is eligible for activation
	blk := moveToEndOfEpoch(t, blk1, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	// check extra validator is added with Withdraw epoch duly set
	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEpoch,
	)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	// STEP 3: move the chain to the next epoch and show that the extra
	// validator is activate and immediately marked for exit
	blk = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, constants.GenesisEpoch+1, extraVal.ActivationEligibilityEpoch)
	require.Equal(t, constants.GenesisEpoch+2, extraVal.ActivationEpoch)
	require.Equal(t, constants.GenesisEpoch+2, extraVal.ExitEpoch)
	require.Equal(t, constants.GenesisEpoch+3, extraVal.WithdrawableEpoch)

	// STEP 4: move the chain to the next epoch and show withdrawals
	// for rejected validator are enqueued then
	blk = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch. extra validator deposits
	// will be withdrawn within 3 blocks (#Validator / MaxValidatorsPerWithdrawalsSweep)
	extraValAddr, err := extraValCreds.ToExecutionAddress()
	require.NoError(t, err)
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	withdrawals := []*engineprimitives.Withdrawal{
		st.EVMInflationWithdrawal(blk.GetTimestamp() + 1),
		{
			Index:     0,
			Validator: extraValIdx,
			Address:   extraValAddr,
			Amount:    extraValDeposit.Amount,
		},
	}
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		withdrawals...,
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
}

// TestTransitionHittingValidatorsCap shows that if the extra
// validator added when validators set is at cap improves amount staked
// an existing validator is removed at the beginning of next epoch.
//
//nolint:maintidx // this end‑to‑end staking‑cap scenario is inherently complex
func TestTransitionHittingValidatorsCap_ExtraBig(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance      = math.Gwei(cs.MaxEffectiveBalance())
		ejectionBalance = math.Gwei(cs.EjectionBalance())
		minBalance      = ejectionBalance + math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed         = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with GetValidatorSetCap validators
	// TODO: consider instead setting state artificially
	var (
		genDeposits      = make(types.Deposits, 0, cs.ValidatorSetCap())
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)

	// let genesis define all available validators
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
	// make a deposit small to be ready for eviction
	genDeposits[0].Amount = minBalance

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	genVals, err := sp.InitializeBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.Len(t, genVals, len(genDeposits))

	// STEP 1: Add an extra validator
	extraValKey, rndSeed := generateTestPK(t, rndSeed)
	extraValCreds, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		extraValDeposit = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      maxBalance,
			Index:       uint64(len(genDeposits)),
		}
	)

	depRoot := append(genDeposits, extraValDeposit).HashTreeRoot()
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{extraValDeposit},
		st.EVMInflationWithdrawal(10),
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), blk1.Body.Deposits))

	// run the test
	valDiff, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	extraValIdx, err := st.ValidatorIndexByPubkey(extraValDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEligibilityEpoch,
	)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEpoch,
	)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	smallestValIdx, err := st.ValidatorIndexByPubkey(genDeposits[0].Pubkey)
	require.NoError(t, err)
	smallestVal, err := st.ValidatorByIndex(smallestValIdx)
	require.NoError(t, err)
	require.Equal(t, constants.GenesisEpoch, smallestVal.ActivationEligibilityEpoch)
	require.Equal(t, constants.GenesisEpoch, smallestVal.ActivationEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		smallestVal.ExitEpoch,
	)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		smallestVal.WithdrawableEpoch,
	)

	// STEP 2: move the chain to the next epoch and show that
	// the extra validator is eligible for activation
	blk := moveToEndOfEpoch(t, blk1, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot), blk.GetTimestamp()+1, []*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	// check extra validator is added with Withdraw epoch duly set
	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.ActivationEpoch,
	)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	smallestVal, err = st.ValidatorByIndex(smallestValIdx)
	require.NoError(t, err)
	require.Equal(t, constants.GenesisEpoch, smallestVal.ActivationEligibilityEpoch)
	require.Equal(t, constants.GenesisEpoch, smallestVal.ActivationEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		smallestVal.ExitEpoch,
	)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		smallestVal.WithdrawableEpoch,
	)

	// STEP 3: move the chain to the next epoch and show that the extra
	// validator is activate and genesis validator immediately marked for exit
	blk = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, valDiff, 2)
	require.Equal(
		t,
		&transition.ValidatorUpdate{
			Pubkey:           extraVal.Pubkey,
			EffectiveBalance: extraVal.EffectiveBalance,
		},
		valDiff[0],
	)
	require.Equal(
		t,
		&transition.ValidatorUpdate{
			Pubkey:           smallestVal.Pubkey,
			EffectiveBalance: 0,
		},
		valDiff[1],
	)

	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, constants.GenesisEpoch+1, extraVal.ActivationEligibilityEpoch)
	require.Equal(t, constants.GenesisEpoch+2, extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(
		t,
		math.Epoch(constants.FarFutureEpoch),
		extraVal.WithdrawableEpoch,
	)

	smallestVal, err = st.ValidatorByIndex(smallestValIdx)
	require.NoError(t, err)
	require.Equal(t, constants.GenesisEpoch, smallestVal.ActivationEligibilityEpoch)
	require.Equal(t, constants.GenesisEpoch, smallestVal.ActivationEpoch)
	require.Equal(t, constants.GenesisEpoch+2, smallestVal.ExitEpoch)
	require.Equal(t, constants.GenesisEpoch+3, smallestVal.WithdrawableEpoch)

	// STEP 4: move the chain to the next epoch and show withdrawal
	// for rejected validator is enqueued within 3 blocks
	// (#Validator / MaxValidatorsPerWithdrawalsSweep)
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx, depRoot)

	valToEvict := genDeposits[0]
	valToEvictAddr, err := valToEvict.Credentials.ToExecutionAddress()
	require.NoError(t, err)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)

	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		st.EVMInflationWithdrawal(blk.GetTimestamp()+1),
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)

	withdrawals := []*engineprimitives.Withdrawal{
		st.EVMInflationWithdrawal(blk.GetTimestamp() + 1),
		{
			Index:     0,
			Validator: smallestValIdx,
			Address:   valToEvictAddr,
			Amount:    valToEvict.Amount,
		},
	}
	blk = buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		blk.GetTimestamp()+1,
		[]*types.Deposit{},
		withdrawals...,
	)
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
}

//nolint:paralleltest // uses envars
func TestValidatorNotWithdrawable(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		belowActiveBalance = math.Gwei(cs.EjectionBalance())
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
