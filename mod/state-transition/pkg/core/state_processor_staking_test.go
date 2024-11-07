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

// TestTransitionUpdateValidators shows that when validator is
// updated (increasing amount), corrensponding balance is updated.
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
		minBalance       = math.Gwei(cs.EffectiveBalanceIncrement())
		emptyAddress     = common.ExecutionAddress{}
		emptyCredentials = types.NewCredentialsFromExecutionAddress(
			emptyAddress,
		)
	)

	// Setup initial state via genesis
	// TODO: consider instead setting state artificially
	var (
		genDeposits = []*types.Deposit{
			{
				Pubkey:      [48]byte{0x01},
				Credentials: emptyCredentials,
				Amount:      maxBalance - 3*minBalance,
				Index:       uint64(0),
			},
			{
				Pubkey:      [48]byte{0x02},
				Credentials: emptyCredentials,
				Amount:      maxBalance - 6*minBalance,
				Index:       uint64(1),
			},
		}
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	_, err = sp.InitializePreminedBeaconStateFromEth1(
		beaconState,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)

	// create test inputs
	var (
		ctx = &transition.Context{
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			ProposerAddress:         dummyProposerAddr,
		}
		blkDeposits = []*types.Deposit{
			{
				Pubkey:      genDeposits[0].Pubkey,
				Credentials: emptyCredentials,
				Amount:      minBalance, // avoid breaching maxBalance
				Index:       genDeposits[0].Index,
			},
		}
	)

	blk := buildNextBlock(
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
			Deposits: blkDeposits,
		},
	)

	// run the test
	_, err = sp.Transition(ctx, beaconState, blk)
	require.NoError(t, err)

	// check validator is duly updated
	expectedValBalance := genDeposits[0].Amount + blkDeposits[0].Amount
	idx, err := beaconState.ValidatorIndexByPubkey(genDeposits[0].Pubkey)
	require.NoError(t, err)
	require.Equal(t, math.U64(genDeposits[0].Index), idx)

	val, err := beaconState.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, genDeposits[0].Pubkey, val.Pubkey)
	require.Equal(t, expectedValBalance, val.EffectiveBalance)

	// check validator balance is updated
	valBal, err := beaconState.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedValBalance, valBal)

	// check that validator index is duly set (1-indexed here, to be fixed)
	latestValIdx, err := beaconState.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)), latestValIdx)
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

	// STEP 1: Setup genesis with GetValidatorSetCapSize validators
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

	// STEP 2: Try and add an extra validator
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
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:     10,
				ExtraData:     []byte("testing"),
				Transactions:  [][]byte{},
				Withdrawals:   []*engineprimitives.Withdrawal{}, // no withdrawals
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{extraValDeposit},
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

	// STEP 3: show that following block must contain withdrawals for
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
		minBalance = math.Gwei(cs.EffectiveBalanceIncrement())
		rndSeed    = 2024 // seed used to generate unique random value
	)

	// STEP 1: Setup genesis with GetValidatorSetCapSize validators
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
	genDeposits[0].Amount = maxBalance - minBalance
	smallestVal := genDeposits[0]
	smallestValAddr := genAddresses[0]

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

	// STEP 2: Add an extra validator
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
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:     10,
				ExtraData:     []byte("testing"),
				Transactions:  [][]byte{},
				Withdrawals:   []*engineprimitives.Withdrawal{}, // no withdrawals
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{extraValDeposit},
		},
	)

	// run the test
	_, err = sp.Transition(ctx, bs, blk1)
	require.NoError(t, err)

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

	// STEP 3: show that following block must contain withdrawals for
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
	_, err = sp.Transition(ctx, bs, blk2)
	require.NoError(t, err)
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
