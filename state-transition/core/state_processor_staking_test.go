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
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestTransitionUpdateValidators shows that when validator is
// updated (increasing amount), corresponding balance is updated.
func TestTransitionUpdateValidators(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance       = math.Gwei(cs.MaxEffectiveBalance())
		increment        = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance       = math.Gwei(cs.EjectionBalance())
		emptyCredentials = types.NewCredentialsFromExecutionAddress(
			common.ExecutionAddress{},
		)
	)

	// STEP 0: Setup initial state via genesis
	var (
		genDeposits = []*types.Deposit{
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
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)
	valDiff, err := sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
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

	blk1 := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    10,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{blkDeposit},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

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
	require.Equal(t, uint64(len(genDeposits)), latestValIdx)

	// STEP 2: check that effective balance is updated once next epoch arrives
	blk := moveToEndOfEpoch(t, blk1, cs, sp, st, ctx)

	// finally the block turning epoch
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
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
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
//nolint:lll // let it be
func TestTransitionCreateValidator(t *testing.T) {
	// Create state processor to test
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

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
		}
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	genVals, err := sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
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

	blk1 := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    10,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{blkDeposit},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

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
	require.Equal(t, uint64(len(genDeposits)), latestValIdx)

	// STEP 2: move the chain to the next epoch and show that
	// the extra validator is eligible for activation
	blk := moveToEndOfEpoch(t, blk1, cs, sp, st, ctx)

	// finally the block turning epoch
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
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff) // new validator is only eligible for activation

	extraValIdx, err := st.ValidatorIndexByPubkey(blkDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	// STEP 3: move the chain to the next epoch and show that
	// the extra validator is activate
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
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
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	expectedBalance = blkDeposit.Amount
	expectedEffectiveBalance = expectedBalance

	balance, err = st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err = st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)
}

func TestTransitionWithdrawals(t *testing.T) {
	cs := setupChain(t, components.BoonetChainSpecType)
	sp, st, _, ctx := setupState(t, cs)

	var (
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		minBalance   = math.Gwei(cs.EffectiveBalanceIncrement())
		credentials0 = types.NewCredentialsFromExecutionAddress(
			common.ExecutionAddress{},
		)
		address1     = common.ExecutionAddress{0x01}
		credentials1 = types.NewCredentialsFromExecutionAddress(address1)
	)

	// Setup initial state so that validator 1 is partially withdrawable.
	var (
		genDeposits = []*types.Deposit{
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
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)
	_, err := sp.InitializePreminedBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, genVersion,
	)
	require.NoError(t, err)

	// Progress state to fork 2.
	progressStateToSlot(t, st, math.U64(spec.BoonetFork2Height))

	// Assert validator 1 balance before withdrawal.
	val1Bal, err := st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val1Bal)

	// Create test inputs.
	blk := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    10,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					// The first withdrawal is always for EVM inflation.
					st.EVMInflationWithdrawal(),
					// Partially withdraw validator 1 by minBalance.
					{
						Index:     0,
						Validator: 1,
						Amount:    minBalance,
						Address:   address1,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
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
	// Use custom chain spec with max withdrawals set to 2.
	csData := spec.BaseSpec()
	csData.DepositEth1ChainID = spec.BoonetEth1ChainID
	csData.MaxWithdrawalsPerPayload = 2
	csData.MaxValidatorsPerWithdrawalsSweepPostUpgrade = 2
	cs, err := chain.NewChainSpec(csData)
	require.NoError(t, err)

	sp, st, _, ctx := setupState(t, cs)

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
		genDeposits = []*types.Deposit{
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
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)
	_, err = sp.InitializePreminedBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, genVersion,
	)
	require.NoError(t, err)

	// Progress state to fork 2.
	progressStateToSlot(t, st, math.U64(spec.BoonetFork2Height))

	// Assert validator balances before withdrawal.
	val0Bal, err := st.GetBalance(math.U64(0))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val0Bal)

	val1Bal, err := st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance+minBalance, val1Bal)

	// Create test inputs.
	blk := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    10,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					// The first withdrawal is always for EVM inflation.
					st.EVMInflationWithdrawal(),
					// Partially withdraw validator 0 by minBalance.
					{
						Index:     0,
						Validator: 0,
						Amount:    minBalance,
						Address:   address0,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
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
	blk = buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    11,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					// The first withdrawal is always for EVM inflation.
					st.EVMInflationWithdrawal(),
					// Partially withdraw validator 1 by minBalance.
					{
						Index:     1,
						Validator: 1,
						Amount:    minBalance,
						Address:   address1,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
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
//
//nolint:lll // let it be
func TestTransitionHittingValidatorsCap_ExtraSmall(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance      = math.Gwei(cs.MaxEffectiveBalance())
		ejectionBalance = math.Gwei(cs.EjectionBalance())
		minBalance      = ejectionBalance + math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed         = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with GetValidatorSetCap validators
	// TODO: consider instead setting state artificially
	var (
		genDeposits      = make([]*types.Deposit, 0, cs.ValidatorSetCap())
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	// let genesis define all available validators
	for idx := range cs.ValidatorSetCap() {
		var (
			key   bytes.B48
			creds types.WithdrawalCredentials
		)
		key, rndSeed = generateTestPK(t, rndSeed)
		creds, rndSeed = generateTestExecutionAddress(t, rndSeed)

		genDeposits = append(genDeposits,
			&types.Deposit{
				Pubkey:      key,
				Credentials: creds,
				Amount:      maxBalance,
				Index:       idx,
			},
		)
	}

	_, err := sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
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

	blk1 := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    10,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{extraValDeposit},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

	// run the test
	valDiff, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	extraValIdx, err := st.ValidatorIndexByPubkey(extraValDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	// STEP 2: move the chain to the next epoch and show that
	// the extra validator is eligible for activation
	_ = moveToEndOfEpoch(t, blk1, cs, sp, st, ctx)

	// finally the block turning epoch
	blk := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	// check extra validator is added with Withdraw epoch duly set
	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	// STEP 3: move the chain to the next epoch and show that the extra validator
	// is activate and immediately marked for exit
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(2), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(2), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(3), extraVal.WithdrawableEpoch)

	// STEP 4: move the chain to the next epoch and show withdrawals
	// for rejected validator are enqueued then
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx)

	// finally the block turning epoch
	extraValAddr, err := extraValCreds.ToExecutionAddress()
	require.NoError(t, err)
	blk = buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
					{
						Index:     0,
						Validator: extraValIdx,
						Address:   extraValAddr,
						Amount:    extraValDeposit.Amount,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	// run the test
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
}

// TestTransitionHittingValidatorsCap shows that if the extra
// validator added when validators set is at cap improves amount staked
// an existing validator is removed at the beginning of next epoch.
//
//nolint:lll // let it be
func TestTransitionHittingValidatorsCap_ExtraBig(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance      = math.Gwei(cs.MaxEffectiveBalance())
		ejectionBalance = math.Gwei(cs.EjectionBalance())
		minBalance      = ejectionBalance + math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed         = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with GetValidatorSetCap validators
	// TODO: consider instead setting state artificially
	var (
		genDeposits      = make([]*types.Deposit, 0, cs.ValidatorSetCap())
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	// let genesis define all available validators
	for idx := range cs.ValidatorSetCap() {
		var (
			key   bytes.B48
			creds types.WithdrawalCredentials
		)
		key, rndSeed = generateTestPK(t, rndSeed)
		creds, rndSeed = generateTestExecutionAddress(t, rndSeed)

		genDeposits = append(genDeposits,
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

	var genVals transition.ValidatorUpdates
	genVals, err := sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
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

	blk1 := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    10,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{extraValDeposit},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

	// run the test
	valDiff, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	extraValIdx, err := st.ValidatorIndexByPubkey(extraValDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	smallestValIdx, err := st.ValidatorIndexByPubkey(genDeposits[0].Pubkey)
	require.NoError(t, err)
	smallestVal, err := st.ValidatorByIndex(smallestValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(0), smallestVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(0), smallestVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), smallestVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), smallestVal.WithdrawableEpoch)

	// STEP 2: move the chain to the next epoch and show that
	// the extra validator is eligible for activation
	_ = moveToEndOfEpoch(t, blk1, cs, sp, st, ctx)

	// finally the block turning epoch
	blk := buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	// run the test
	valDiff, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	// check extra validator is added with Withdraw epoch duly set
	extraVal, err = st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	smallestVal, err = st.ValidatorByIndex(smallestValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(0), smallestVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(0), smallestVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), smallestVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), smallestVal.WithdrawableEpoch)

	// STEP 3: move the chain to the next epoch and show that the extra validator
	// is activate and genesis validator immediately marked for exit
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
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
	require.Equal(t, math.Epoch(1), extraVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(2), extraVal.ActivationEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.ExitEpoch)
	require.Equal(t, math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch)

	smallestVal, err = st.ValidatorByIndex(smallestValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Epoch(0), smallestVal.ActivationEligibilityEpoch)
	require.Equal(t, math.Epoch(0), smallestVal.ActivationEpoch)
	require.Equal(t, math.Epoch(2), smallestVal.ExitEpoch)
	require.Equal(t, math.Epoch(3), smallestVal.WithdrawableEpoch)

	// STEP 4: move the chain to the next epoch and show withdrawal
	// for rejected validator is enqueued
	_ = moveToEndOfEpoch(t, blk, cs, sp, st, ctx)

	valToEvict := genDeposits[0]
	valToEvictAddr, err := valToEvict.Credentials.ToExecutionAddress()
	require.NoError(t, err)

	// finally the block turning epoch
	blk = buildNextBlock(
		t,
		st,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					st.EVMInflationWithdrawal(),
					{
						Index:     0,
						Validator: smallestValIdx,
						Address:   valToEvictAddr,
						Amount:    valToEvict.Amount,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	// run the test
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
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
	cs chain.Spec[bytes.B4, math.U64, common.ExecutionAddress, math.U64, any],
	sp *TestStateProcessorT,
	st *TestBeaconStateT,
	ctx *transition.Context,
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
				Eth1Data: &types.Eth1Data{},
				Deposits: []*types.Deposit{},
			},
		)

		vals, err := sp.Transition(ctx, st, blk)
		require.NoError(t, err)
		require.Empty(t, vals) // no vals changes expected before next epoch
	}
	return blk
}
