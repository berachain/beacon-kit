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
	updatedVals, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, updatedVals) // validators set updates only at epoch turn

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

	newEpochVals, err := sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, newEpochVals, 1) // just topped up one validator

	expectedBalance = genDeposits[2].Amount + blkDeposit.Amount
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
	updatedVals, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, updatedVals) // validators set updates only at epoch turn

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

	newEpochVals, err := sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, newEpochVals, 1) // just added 1 validator

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
func TestTransitionHittingValidatorsCap_ExtraSmall(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		rndSeed    = 2024 // seed used to generate unique random value
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
	_, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)

	// check extra validator is added with Withdraw epoch duly set
	extraValIdx, err := st.ValidatorIndexByPubkey(extraValDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, extraValDeposit.Pubkey, extraVal.Pubkey)
	require.Equal(t, math.Slot(1), extraVal.WithdrawableEpoch)

	extraValBalance, err := st.GetBalance(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, extraValDeposit.Amount, extraValBalance)

	// STEP 2: move the chain to the next epoch and show withdrawals
	// for rejected validator are enqueuued then
	_ = moveToEndOfEpoch(t, blk1, cs, sp, st, ctx)

	// finally the block turning epoch
	extraValAddr, err := extraValCreds.ToExecutionAddress()
	require.NoError(t, err)
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
func TestTransitionHittingValidatorsCap_ExtraBig(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
		rndSeed    = 2024 // seed used to generate unique random value
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
	genDeposits[0].Amount = minBalance + increment
	smallestVal := genDeposits[0]
	smallestValAddr, err := genDeposits[0].Credentials.ToExecutionAddress()
	require.NoError(t, err)

	var genVals transition.ValidatorUpdates
	genVals, err = sp.InitializePreminedBeaconStateFromEth1(
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
	var vals transition.ValidatorUpdates
	vals, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// check smallest validator is updated with Withdraw epoch duly set
	smallValIdx, err := st.ValidatorIndexByPubkey(smallestVal.Pubkey)
	require.NoError(t, err)
	smallVal, err := st.ValidatorByIndex(smallValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Slot(1), smallVal.WithdrawableEpoch)

	smallestValBalance, err := st.GetBalance(smallValIdx)
	require.NoError(t, err)
	require.Equal(t, smallestVal.Amount, smallestValBalance)

	// check that extra validator is added
	extraValIdx, err := st.ValidatorIndexByPubkey(extraValKey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t,
		math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch,
	)

	// STEP 2: move chain to next epoch to see extra validator
	// be activated and withdraws for evicted validator
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
					{
						Index:     0,
						Validator: smallValIdx,
						Address:   smallestValAddr,
						Amount:    smallestVal.Amount,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	vals, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.LessOrEqual(t, uint64(len(vals)), cs.ValidatorSetCap())
	require.Len(t, vals, 2) // just replaced one validator

	// check that we added the incoming validator at the epoch turn
	require.Equal(t, extraVal.Pubkey, vals[0].Pubkey)
	require.Equal(t, extraVal.EffectiveBalance, vals[0].EffectiveBalance)

	// check that we removed the smallest validator at the epoch turn
	require.Equal(t, smallVal.Pubkey, vals[1].Pubkey)
	require.Equal(t, math.Gwei(0), vals[1].EffectiveBalance)
}

// TestTransitionValidatorCap_DoubleEviction show that the
// eviction mechanism works fine even if multiple evictions
// happen in the same epoch.
//
// //nolint:maintidx // TODO: simplify
func TestTransitionValidatorCap_DoubleEviction(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 0: fill genesis with validators till cap. Let two of them
	// have smaller balance than others, so to be amenable for eviction.
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
	genDeposits[0].Amount = minBalance + increment
	smallest1Val := genDeposits[0]
	smallest1ValAddr, err := genDeposits[0].Credentials.ToExecutionAddress()
	require.NoError(t, err)

	genDeposits[1].Amount = minBalance + 2*increment
	smallestVal2 := genDeposits[1]
	smallestVal2Addr, err := genDeposits[1].Credentials.ToExecutionAddress()
	require.NoError(t, err)

	var genVals transition.ValidatorUpdates
	genVals, err = sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)
	require.Len(t, genVals, len(genDeposits))

	// STEP 1: Add an extra validator
	extraVal1Key, rndSeed := generateTestPK(t, rndSeed)
	extraVal1Creds, rndSeed := generateTestExecutionAddress(t, rndSeed)
	extraValDeposit1 := &types.Deposit{
		Pubkey:      extraVal1Key,
		Credentials: extraVal1Creds,
		Amount:      maxBalance - increment,
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
			Deposits: []*types.Deposit{extraValDeposit1},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

	// run the test
	vals, err := sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// check the smallest validator Withdraw epoch is updated
	smallVal1Idx, err := st.ValidatorIndexByPubkey(smallest1Val.Pubkey)
	require.NoError(t, err)
	smallVal1, err := st.ValidatorByIndex(smallVal1Idx)
	require.NoError(t, err)
	require.Equal(t, math.Slot(1), smallVal1.WithdrawableEpoch)

	smallVal1Balance, err := st.GetBalance(smallVal1Idx)
	require.NoError(t, err)
	require.Equal(t, smallest1Val.Amount, smallVal1Balance)

	// check that extra validator is added
	extraVal1Idx, err := st.ValidatorIndexByPubkey(extraVal1Key)
	require.NoError(t, err)
	extraVal1, err := st.ValidatorByIndex(extraVal1Idx)
	require.NoError(t, err)
	require.Equal(t,
		math.Epoch(constants.FarFutureEpoch), extraVal1.WithdrawableEpoch,
	)

	// STEP 2: add a second, large deposit to evict second smallest validator
	extraVal2Key, rndSeed := generateTestPK(t, rndSeed)
	extraVal2Creds, _ := generateTestExecutionAddress(t, rndSeed)
	extraVal2Deposit := &types.Deposit{
		Pubkey:      extraVal2Key,
		Credentials: extraVal2Creds,
		Amount:      maxBalance,
		Index:       uint64(len(genDeposits) + 1),
	}

	blk2 := buildNextBlock(
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
			Deposits: []*types.Deposit{extraVal2Deposit},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk2.Body.Deposits))

	// run the test
	vals, err = sp.Transition(ctx, st, blk2)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// check the second smallest validator Withdraw epoch is updated
	smallVal2Idx, err := st.ValidatorIndexByPubkey(smallestVal2.Pubkey)
	require.NoError(t, err)
	smallVal2, err := st.ValidatorByIndex(smallVal2Idx)
	require.NoError(t, err)
	require.Equal(t, math.Slot(1), smallVal2.WithdrawableEpoch)

	smallVal2Balance, err := st.GetBalance(smallVal2Idx)
	require.NoError(t, err)
	require.Equal(t, smallestVal2.Amount, smallVal2Balance)

	// check that extra validator is added
	extraVal2Idx, err := st.ValidatorIndexByPubkey(extraVal2Key)
	require.NoError(t, err)
	extraVal2, err := st.ValidatorByIndex(extraVal2Idx)
	require.NoError(t, err)
	require.Equal(t,
		math.Epoch(constants.FarFutureEpoch), extraVal2.WithdrawableEpoch,
	)

	// STEP 3: move to next epoch
	_ = moveToEndOfEpoch(t, blk2, cs, sp, st, ctx)

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
					{
						Index:     0,
						Validator: smallVal1Idx,
						Address:   smallest1ValAddr,
						Amount:    smallest1Val.Amount,
					},
					{
						Index:     1,
						Validator: smallVal2Idx,
						Address:   smallestVal2Addr,
						Amount:    smallestVal2.Amount,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	vals, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.LessOrEqual(t, uint64(len(vals)), cs.ValidatorSetCap())
	require.Len(t, vals, 4) // just replaced two validators

	// turn vals into map to avoid ordering issues
	valsSet := make(map[string]*transition.ValidatorUpdate)
	for _, v := range vals {
		valsSet[v.Pubkey.String()] = v
	}
	require.Equal(t, len(vals), len(valsSet)) // no duplicates

	// check that we added the incoming validator at the epoch turn
	addedVal1, found := valsSet[extraVal1.Pubkey.String()]
	require.True(t, found)
	require.Equal(t, extraVal1.EffectiveBalance, addedVal1.EffectiveBalance)

	addedVal2, found := valsSet[extraVal2.Pubkey.String()]
	require.True(t, found)
	require.Equal(t, extraVal2.EffectiveBalance, addedVal2.EffectiveBalance)

	// check that we removed the smallest validators at the epoch turn
	removedVal1, found := valsSet[smallVal1.Pubkey.String()]
	require.True(t, found)
	require.Equal(t, math.Gwei(0), removedVal1.EffectiveBalance)

	removeldVal2, found := valsSet[smallVal2.Pubkey.String()]
	require.True(t, found)
	require.Equal(t, math.Gwei(0), removeldVal2.EffectiveBalance)
}

// TestTransitionValidatorCap_IncreasingBalance_ExtraSmall shows that the
// eviction mechanism works fine even if stake required to validate is
// reached over multiple deposits.
func TestTransitionValidatorCap_IncreasingBalance_ExtraSmall(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		minBalance = math.Gwei(cs.EjectionBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with max number of validators
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

	// STEP 1: make a deposit below MinState to start
	extraValKey, rndSeed := generateTestPK(t, rndSeed)
	extraValCreds, _ := generateTestExecutionAddress(t, rndSeed)
	extraValAddr, err := extraValCreds.ToExecutionAddress()
	require.NoError(t, err)
	var (
		initialDeposit = minBalance / 2
		deposit1       = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      initialDeposit,
			Index:       uint64(len(genDeposits)),
		}

		topUpDeposit = (minBalance + increment) - initialDeposit
		deposit2     = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      topUpDeposit,
			Index:       deposit1.Index + 1,
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
			Deposits: []*types.Deposit{deposit1},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

	// run the test
	_, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)

	// STEP 2: top up the deposit to reach MinStake required to validate
	// before the epoch turns
	blk2 := buildNextBlock(
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
			Deposits: []*types.Deposit{deposit2},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk2.Body.Deposits))

	// run the test
	_, err = sp.Transition(ctx, st, blk2)
	require.NoError(t, err)

	// STEP 3: check that updated validator is dropped at epoch turn
	blk := moveToEndOfEpoch(t, blk2, cs, sp, st, ctx)

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
					&engineprimitives.Withdrawal{
						Index:     0,
						Validator: math.ValidatorIndex(len(genDeposits)),
						Address:   extraValAddr,
						Amount:    initialDeposit + topUpDeposit,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	valDiffs, err := sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Empty(t, valDiffs)
}

func TestTransitionValidatorCap_IncreasingBalance_ExtraBig(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		minBalance = math.Gwei(cs.EjectionBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with max number of validators
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
				Amount:      minBalance,
				Index:       idx,
			},
		)
	}
	genVal0Addr, err := genDeposits[0].Credentials.ToExecutionAddress()
	require.NoError(t, err)

	_, err = sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)

	// STEP 1: make a deposit below MinState to start
	extraValKey, rndSeed := generateTestPK(t, rndSeed)
	extraValCreds, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		initialDeposit = minBalance / 2
		deposit1       = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      initialDeposit,
			Index:       uint64(len(genDeposits)),
		}

		topUpDeposit = (maxBalance + increment) - initialDeposit
		deposit2     = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      topUpDeposit,
			Index:       deposit1.Index + 1,
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
			Deposits: []*types.Deposit{deposit1},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk1.Body.Deposits))

	// run the test
	_, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)

	// STEP 2: top up the deposit to reach MinStake required to validate
	// before the epoch turns
	blk2 := buildNextBlock(
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
			Deposits: []*types.Deposit{deposit2},
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blk2.Body.Deposits))

	// run the test
	_, err = sp.Transition(ctx, st, blk2)
	require.NoError(t, err)

	// STEP 3: check that updated validator is dropped at epoch turn
	blk := moveToEndOfEpoch(t, blk2, cs, sp, st, ctx)

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
					&engineprimitives.Withdrawal{
						Index:     0,
						Validator: math.ValidatorIndex(0), // drop first genesis validator I think
						Address:   genVal0Addr,
						Amount:    genDeposits[0].Amount,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{},
		},
	)

	valDiffs, err := sp.Transition(ctx, st, blk)
	require.NoError(t, err)
	require.Len(t, valDiffs, 2) // just replaced one validator

	// check that we added the incoming validator at the epoch turn
	extraValIdx, err := st.ValidatorIndexByPubkey(extraValKey)
	require.NoError(t, err)
	extraVal, err := st.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, extraVal.EffectiveBalance, minBalance+topUpDeposit)

	// TODO: CHECK THAT A GENESIS VALIDATOR IS DROPPED
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
