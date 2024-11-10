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
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	cryptomocks "github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestTransitionUpdateValidator shows the lifecycle
// of a validator's balance updates.
func TestTransitionUpdateValidator(t *testing.T) {
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
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
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
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{blkDeposit},
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
	var blk *types.BeaconBlock
	for i := 1; i < int(cs.SlotsPerEpoch())-1; i++ {
		blk = buildNextBlock(
			t,
			beaconState,
			&types.BeaconBlockBody{
				ExecutionPayload: dummyExecutionPayload,
				Eth1Data:         &types.Eth1Data{},
				Deposits:         []*types.Deposit{},
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
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{},
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

// TestTransitionCreateValidator shows the lifecycle
// of a validator creation.
func TestTransitionCreateValidator(t *testing.T) {
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
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
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
			Pubkey:      [48]byte{0xff}, // a new key for a new validator
			Credentials: emptyCredentials,
			Amount:      maxBalance,
			Index:       uint64(len(genDeposits)),
		}
	)

	blk1 := buildNextBlock(
		t,
		beaconState,
		&types.BeaconBlockBody{
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{blkDeposit},
		},
	)

	// run the test
	updatedVals, err := sp.Transition(ctx, beaconState, blk1)
	require.NoError(t, err)
	require.Empty(t, updatedVals) // validators set updates only at epoch turn

	// check validator balances are duly updated
	var (
		expectedBalance          = blkDeposit.Amount
		expectedEffectiveBalance = expectedBalance
	)
	idx, err := beaconState.ValidatorIndexByPubkey(blkDeposit.Pubkey)
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
	var blk *types.BeaconBlock
	for i := 1; i < int(cs.SlotsPerEpoch())-1; i++ {
		blk = buildNextBlock(
			t,
			beaconState,
			&types.BeaconBlockBody{
				ExecutionPayload: dummyExecutionPayload,
				Eth1Data:         &types.Eth1Data{},
				Deposits:         []*types.Deposit{},
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
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{},
		},
	)

	newEpochVals, err := sp.Transition(ctx, beaconState, blk)
	require.NoError(t, err)
	require.Len(t, newEpochVals, len(genDeposits)+1)

	// Assuming genesis order is preserved here which is not necessary
	// TODO: remove this assumption

	// all genesis validators are unchanged
	for i := range len(genDeposits) {
		require.Equal(t, genVals[i], newEpochVals[i], fmt.Sprintf("idx: %d", i))
	}

	expectedBalance = blkDeposit.Amount
	expectedEffectiveBalance = expectedBalance

	balance, err = beaconState.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	val, err = beaconState.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, expectedEffectiveBalance, val.EffectiveBalance)
}

// TestTransitionHittingValidatorsCap shows that the extra
// validator added when validators set is at cap is immediately
// scheduled for withdrawal along with its deposit if it does not
// improve staked amount.
func TestTransitionHittingValidatorsCap_ExtraSmall(t *testing.T) {
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
	bs := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with GetValidatorSetCapSize validators
	// TODO: consider instead setting state artificially
	var (
		genDeposits      = make([]*types.Deposit, 0, cs.GetValidatorSetCapSize())
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	// let genesis define all available validators
	for idx := range cs.GetValidatorSetCapSize() {
		var (
			key   bytes.B48
			creds types.WithdrawalCredentials
		)
		key, rndSeed = generateTestPK(t, rndSeed)
		creds, _, rndSeed = generateTestExecutionAddress(t, rndSeed)

		genDeposits = append(genDeposits,
			&types.Deposit{
				Pubkey:      key,
				Credentials: creds,
				Amount:      maxBalance,
				Index:       uint64(idx),
			},
		)
	}

	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	_, err = sp.InitializePreminedBeaconStateFromEth1(
		bs,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)

	// STEP 1: Try and add an extra validator
	extraValKey, rndSeed := generateTestPK(t, rndSeed)
	extraValCreds, extraValAddr, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		ctx = &transition.Context{
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			ProposerAddress:         dummyProposerAddr,
		}
		extraValDeposit = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      maxBalance,
			Index:       uint64(len(genDeposits)),
		}
	)

	blk1 := buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{extraValDeposit},
		},
	)

	// run the test
	_, err = sp.Transition(ctx, bs, blk1)
	require.NoError(t, err)

	// check extra validator is added with Withdraw epoch duly set
	extraValIdx, err := bs.ValidatorIndexByPubkey(extraValDeposit.Pubkey)
	require.NoError(t, err)
	extraVal, err := bs.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, extraValDeposit.Pubkey, extraVal.Pubkey)
	require.Equal(t, math.Slot(0), extraVal.WithdrawableEpoch)

	extraValBalance, err := bs.GetBalance(extraValIdx)
	require.NoError(t, err)
	require.Equal(t, extraValDeposit.Amount, extraValBalance)

	// STEP 2: show that following block must contain withdrawals for
	// the rejected validator
	blk2 := buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
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
	_, err = sp.Transition(ctx, bs, blk2)
	require.NoError(t, err)
}

// TestTransitionHittingValidatorsCap shows that if the extra
// validator added when validators set is at cap improves amount staked
// an existing validator is immediately scheduled for withdrawal
// along with its deposit.
func TestTransitionHittingValidatorsCap_ExtraBig(t *testing.T) {
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
	bs := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 0: Setup genesis with GetValidatorSetCapSize validators
	// TODO: consider instead setting state artificially
	var (
		genDeposits      = make([]*types.Deposit, 0, cs.GetValidatorSetCapSize())
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	// let genesis define all available validators
	genAddresses := make([]common.ExecutionAddress, 0)
	for idx := range cs.GetValidatorSetCapSize() {
		var (
			key   bytes.B48
			creds types.WithdrawalCredentials
			addr  common.ExecutionAddress
		)
		key, rndSeed = generateTestPK(t, rndSeed)
		creds, addr, rndSeed = generateTestExecutionAddress(t, rndSeed)

		genDeposits = append(genDeposits,
			&types.Deposit{
				Pubkey:      key,
				Credentials: creds,
				Amount:      maxBalance,
				Index:       uint64(idx),
			},
		)
		genAddresses = append(genAddresses, addr)
	}
	// make a deposit small to be ready for eviction
	genDeposits[0].Amount = minBalance + increment
	smallestVal := genDeposits[0]
	smallestValAddr := genAddresses[0]

	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	var genVals transition.ValidatorUpdates
	genVals, err = sp.InitializePreminedBeaconStateFromEth1(
		bs,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)
	require.Len(t, genVals, len(genDeposits))

	// STEP 1: Add an extra validator
	extraValKey, rndSeed := generateTestPK(t, rndSeed)
	extraValCreds, _, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		ctx = &transition.Context{
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			ProposerAddress:         dummyProposerAddr,
		}
		extraValDeposit = &types.Deposit{
			Pubkey:      extraValKey,
			Credentials: extraValCreds,
			Amount:      maxBalance,
			Index:       uint64(len(genDeposits)),
		}
	)

	blk1 := buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{extraValDeposit},
		},
	)

	// run the test
	var vals transition.ValidatorUpdates
	vals, err = sp.Transition(ctx, bs, blk1)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// check smallest validator is updated with Withdraw epoch duly set
	smallValIdx, err := bs.ValidatorIndexByPubkey(smallestVal.Pubkey)
	require.NoError(t, err)
	smallVal, err := bs.ValidatorByIndex(smallValIdx)
	require.NoError(t, err)
	require.Equal(t, math.Slot(0), smallVal.WithdrawableEpoch)

	smallestValBalance, err := bs.GetBalance(smallValIdx)
	require.NoError(t, err)
	require.Equal(t, smallestVal.Amount, smallestValBalance)

	// check that extra validator is added
	extraValIdx, err := bs.ValidatorIndexByPubkey(extraValKey)
	require.NoError(t, err)
	extraVal, err := bs.ValidatorByIndex(extraValIdx)
	require.NoError(t, err)
	require.Equal(t,
		math.Epoch(constants.FarFutureEpoch), extraVal.WithdrawableEpoch,
	)

	// STEP 2: show that following block must contain withdrawals for
	// the evicted, smallest validator
	blk2 := buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
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

	// run the test
	vals, err = sp.Transition(ctx, bs, blk2)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// STEP 3: moving chain forward to next epoch to see extra validator
	// be activated, just replaced
	var blk *types.BeaconBlock
	for i := 2; i < int(cs.SlotsPerEpoch())-1; i++ {
		blk = buildNextBlock(
			t,
			bs,
			&types.BeaconBlockBody{
				ExecutionPayload: dummyExecutionPayload,
				Eth1Data:         &types.Eth1Data{},
				Deposits:         []*types.Deposit{},
			},
		)

		vals, err = sp.Transition(ctx, bs, blk)
		require.NoError(t, err)
		require.Empty(t, vals) // no vals changes expected before next epoch
	}

	blk = buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{},
		},
	)

	vals, err = sp.Transition(ctx, bs, blk)
	require.NoError(t, err)
	require.LessOrEqual(t, uint32(len(vals)), cs.GetValidatorSetCapSize())
	require.Len(t, vals, len(genDeposits)) // just replaced one validator

	// check that we removed the smallest validator at the epoch turn
	removedVals := ValUpdatesDiff(genVals, vals)
	require.Len(t, removedVals, 1)
	require.Equal(t, smallVal.Pubkey, removedVals[0].Pubkey)

	// check that we added the incoming validator at the epoch turn
	addedVals := ValUpdatesDiff(vals, genVals)
	require.Len(t, addedVals, 1)
	require.Equal(t, extraVal.Pubkey, addedVals[0].Pubkey)
}

// show that eviction mechanism works fine even if multiple evictions
// happen in the same epoch.
func TestTransitionValidatorCap_DoubleEviction(t *testing.T) {
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
	bs := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 0: fill genesis with validators till cap. Let two of them
	// have smaller balance than others, so to be amenable for eviction.
	var (
		genDeposits      = make([]*types.Deposit, 0, cs.GetValidatorSetCapSize())
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	// let genesis define all available validators
	genAddresses := make([]common.ExecutionAddress, 0)
	for idx := range cs.GetValidatorSetCapSize() {
		var (
			key   bytes.B48
			creds types.WithdrawalCredentials
			addr  common.ExecutionAddress
		)
		key, rndSeed = generateTestPK(t, rndSeed)
		creds, addr, rndSeed = generateTestExecutionAddress(t, rndSeed)

		genDeposits = append(genDeposits,
			&types.Deposit{
				Pubkey:      key,
				Credentials: creds,
				Amount:      maxBalance,
				Index:       uint64(idx),
			},
		)
		genAddresses = append(genAddresses, addr)
	}
	// make a deposit small to be ready for eviction
	genDeposits[0].Amount = minBalance + increment
	smallest1Val := genDeposits[0]
	smallest1ValAddr := genAddresses[0]

	genDeposits[1].Amount = minBalance + 2*increment
	smallest2Val := genDeposits[1]
	// smallest2ValAddr := genAddresses[1]

	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	var genVals transition.ValidatorUpdates
	genVals, err = sp.InitializePreminedBeaconStateFromEth1(
		bs,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)
	require.Len(t, genVals, len(genDeposits))

	// STEP 1: Add an extra validator
	extraVal1Key, rndSeed := generateTestPK(t, rndSeed)
	extraVal1Creds, _, _ := generateTestExecutionAddress(t, rndSeed)
	var (
		ctx = &transition.Context{
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			ProposerAddress:         dummyProposerAddr,
		}
		extraValDeposit1 = &types.Deposit{
			Pubkey:      extraVal1Key,
			Credentials: extraVal1Creds,
			Amount:      maxBalance,
			Index:       uint64(len(genDeposits)),
		}
	)

	blk1 := buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: dummyExecutionPayload,
			Eth1Data:         &types.Eth1Data{},
			Deposits:         []*types.Deposit{extraValDeposit1},
		},
	)

	// run the test
	vals, err := sp.Transition(ctx, bs, blk1)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// check the smallest validator Withdraw epoch is updated
	smallVal1Idx, err := bs.ValidatorIndexByPubkey(smallest1Val.Pubkey)
	require.NoError(t, err)
	smallVal1, err := bs.ValidatorByIndex(smallVal1Idx)
	require.NoError(t, err)
	require.Equal(t, math.Slot(0), smallVal1.WithdrawableEpoch)

	smallVal1Balance, err := bs.GetBalance(smallVal1Idx)
	require.NoError(t, err)
	require.Equal(t, smallest1Val.Amount, smallVal1Balance)

	// check that extra validator is added
	extraVal1Idx, err := bs.ValidatorIndexByPubkey(extraVal1Key)
	require.NoError(t, err)
	extraVal1, err := bs.ValidatorByIndex(extraVal1Idx)
	require.NoError(t, err)
	require.Equal(t,
		math.Epoch(constants.FarFutureEpoch), extraVal1.WithdrawableEpoch,
	)

	// STEP 2: add a second, large deposit to evict second smallest validator
	extraVal2Key, rndSeed := generateTestPK(t, rndSeed)
	extraVal2Creds, _, _ := generateTestExecutionAddress(t, rndSeed)
	extraVal2Deposit := &types.Deposit{
		Pubkey:      extraVal2Key,
		Credentials: extraVal2Creds,
		Amount:      maxBalance,
		Index:       uint64(len(genDeposits) + 1),
	}

	blk2 := buildNextBlock(
		t,
		bs,
		&types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:    blk1.Body.ExecutionPayload.Timestamp + 1,
				ExtraData:    []byte("testing"),
				Transactions: [][]byte{},
				Withdrawals: []*engineprimitives.Withdrawal{
					{
						Index:     0,
						Validator: smallVal1Idx,
						Address:   smallest1ValAddr,
						Amount:    smallest1Val.Amount,
					},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{extraVal2Deposit},
		},
	)

	// run the test
	vals, err = sp.Transition(ctx, bs, blk2)
	require.NoError(t, err)
	require.Empty(t, vals) // no vals changes expected before next epoch

	// check the second smallest validator Withdraw epoch is updated
	smallVal2Idx, err := bs.ValidatorIndexByPubkey(smallest2Val.Pubkey)
	require.NoError(t, err)
	smallVal2, err := bs.ValidatorByIndex(smallVal2Idx)
	require.NoError(t, err)
	require.Equal(t, math.Slot(0), smallVal2.WithdrawableEpoch)

	smallVal2Balance, err := bs.GetBalance(smallVal2Idx)
	require.NoError(t, err)
	require.Equal(t, smallest2Val.Amount, smallVal2Balance)

	// check that extra validator is added
	extraVal2Idx, err := bs.ValidatorIndexByPubkey(extraVal2Key)
	require.NoError(t, err)
	extraVal2, err := bs.ValidatorByIndex(extraVal2Idx)
	require.NoError(t, err)
	require.Equal(t,
		math.Epoch(constants.FarFutureEpoch), extraVal2.WithdrawableEpoch,
	)
}

func generateTestExecutionAddress(
	t *testing.T,
	rndSeed int,
) (types.WithdrawalCredentials, common.ExecutionAddress, int) {
	t.Helper()

	addrStr := strconv.Itoa(rndSeed)
	addrBytes := bytes.ExtendToSize([]byte(addrStr), bytes.B20Size)
	execAddr, err := bytes.ToBytes20(addrBytes)
	require.NoError(t, err)
	rndSeed++
	return types.NewCredentialsFromExecutionAddress(
		common.ExecutionAddress(execAddr),
	), common.ExecutionAddress(execAddr), rndSeed
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
