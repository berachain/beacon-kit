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
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

	sp := core.NewStateProcessor[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		*TestBeaconStateT,
		*transition.Context,
		*types.Deposit,
		*types.Eth1Data,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*TestKVStoreT,
		*types.Validator,
		types.Validators,
		*engineprimitives.Withdrawal,
		engineprimitives.Withdrawals,
		types.WithdrawalCredentials,
	](
		cs,
		execEngine,
		mocksSigner,
		func(bytes.B48) ([]byte, error) {
			return dummyProposerAddr, nil
		},
	)

	kvStore, err := initTestStore()
	require.NoError(t, err)
	beaconState := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		minBalance = math.Gwei(cs.EffectiveBalanceIncrement())
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

	// here we duly update state root, similarly to what we do in processSlot
	genBlockHeader, err := beaconState.GetLatestBlockHeader()
	require.NoError(t, err)
	genStateRoot := beaconState.HashTreeRoot()
	genBlockHeader.SetStateRoot(genStateRoot)

	blk := &types.BeaconBlock{
		Slot:          genBlockHeader.GetSlot() + 1,
		ProposerIndex: genBlockHeader.GetProposerIndex(),
		ParentRoot:    genBlockHeader.HashTreeRoot(),
		StateRoot:     common.Root{},
		Body: &types.BeaconBlockBody{
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
	}

	// run the test
	vals, err := sp.Transition(ctx, beaconState, blk)

	// check outputs
	require.NoError(t, err)
	require.Zero(t, vals) // just update, no new validators

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

// NOTE: not sure this is any helpful. The point would be testing
// what happens when deposits made are not multiple of minAmount
// and how activation happens henceforth

// TestTransitionPartialDeposits checks that stake needed to
// activate a validator can be accrued across multiple deposits.
// Also the test checks that any deposit in excess of MaxEffectiveBalance is
// enqueued for withdrawal

// TODO: The test should check that deposit signature is checked only
// the first time.
func TestTransitionPartialDeposits(t *testing.T) {
	// Create state processor to test
	cs := spec.BetnetChainSpec()
	execEngine := mocks.NewExecutionEngine[
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		engineprimitives.Withdrawals,
	](t)
	mocksSigner := &cryptomocks.BLSSigner{}
	dummyProposerAddr := []byte{0xff}

	sp := core.NewStateProcessor[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		*TestBeaconStateT,
		*transition.Context,
		*types.Deposit,
		*types.Eth1Data,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*TestKVStoreT,
		*types.Validator,
		types.Validators,
		*engineprimitives.Withdrawal,
		engineprimitives.Withdrawals,
		types.WithdrawalCredentials,
	](
		cs,
		execEngine,
		mocksSigner,
		func(bytes.B48) ([]byte, error) {
			return dummyProposerAddr, nil
		},
	)

	kvStore, err := initTestStore()
	require.NoError(t, err)
	beaconState := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		minBalance = math.Gwei(cs.EffectiveBalanceIncrement())
	)

	// STEP 1: init state with genesis with a single validator
	// TODO: consider instead setting state artificially
	var (
		genDeposit = &types.Deposit{
			Pubkey:      [48]byte{0x01},
			Credentials: emptyCredentials,
			Amount:      maxBalance - 3*minBalance,
			Index:       uint64(0),
		}
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)

	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	vals0, err := sp.InitializePreminedBeaconStateFromEth1(
		beaconState,
		[]*types.Deposit{genDeposit},
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)
	require.Len(t, vals0, 1)
	require.Equal(t, genDeposit.Pubkey, vals0[0].Pubkey)

	// STEP 2: Deposit an amount below minimum required.
	var (
		ctx = &transition.Context{
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			ProposerAddress:         dummyProposerAddr,
		}
		blkDeposit1 = &types.Deposit{
			Pubkey:      [48]byte{0xff},
			Credentials: emptyCredentials,
			Amount:      minBalance,
			Index:       1,
		}
	)

	// here we duly update state root, similarly to what we do in processSlot
	genBlockHeader, err := beaconState.GetLatestBlockHeader()
	require.NoError(t, err)
	genStateRoot := beaconState.HashTreeRoot()
	genBlockHeader.SetStateRoot(genStateRoot)

	blk1 := &types.BeaconBlock{
		Slot:          genBlockHeader.GetSlot() + 1,
		ProposerIndex: genBlockHeader.GetProposerIndex(),
		ParentRoot:    genBlockHeader.HashTreeRoot(),
		StateRoot:     common.Root{},
		Body: &types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:     10,
				ExtraData:     []byte("testing"),
				Transactions:  [][]byte{},
				Withdrawals:   []*engineprimitives.Withdrawal{}, // no withdrawals
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{blkDeposit1},
		},
	}

	// run the test
	vals1, err := sp.Transition(ctx, beaconState, blk1)

	// check outputs
	require.NoError(t, err)
	require.Zero(t, vals1) // added validator has not enough stake to be active

	// check that the validator is added with the expected EffectiveBalance 0
	idx, err := beaconState.ValidatorIndexByPubkey(blkDeposit1.Pubkey)
	require.NoError(t, err)

	val1, err := beaconState.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, blkDeposit1.Amount, val1.EffectiveBalance)

	// also check that validator balance matched with deposited amount
	bal1, err := beaconState.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, blkDeposit1.Amount, bal1)

	// STEP 3: Deposit the amount required to reach the required stake to validate
	var blkDeposit2 = &types.Deposit{
		Pubkey:      blkDeposit1.Pubkey,
		Credentials: emptyCredentials,
		Amount:      maxBalance - blkDeposit1.Amount,
		Index:       2,
	}

	// here we duly update state root, similarly to what we do in processSlot
	blk1Header, err := beaconState.GetLatestBlockHeader()
	require.NoError(t, err)
	blk1StateRoot := beaconState.HashTreeRoot()
	blk1Header.SetStateRoot(blk1StateRoot)

	blk2 := &types.BeaconBlock{
		Slot:          blk1.GetSlot() + 1,
		ProposerIndex: blk1.GetProposerIndex(),
		ParentRoot:    blk1Header.HashTreeRoot(),
		StateRoot:     common.Root{},
		Body: &types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Timestamp:     blk1.Body.ExecutionPayload.Number + 1,
				ExtraData:     []byte("testing"),
				Transactions:  [][]byte{},
				Withdrawals:   []*engineprimitives.Withdrawal{}, // no withdrawals
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: []*types.Deposit{blkDeposit2},
		},
	}

	// run the test
	vals2, err := sp.Transition(ctx, beaconState, blk2)

	// check outputs
	require.NoError(t, err)
	require.Empty(t, vals2) // validator will be returned next block

	// check that the validator is added with the expected EffectiveBalance 0
	idx, err = beaconState.ValidatorIndexByPubkey(blkDeposit2.Pubkey)
	require.NoError(t, err)

	val2, err := beaconState.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, blkDeposit1.Amount+blkDeposit2.Amount, val2.EffectiveBalance)

	// also check that validator balance matched with deposited amount
	bal2, err := beaconState.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, blkDeposit1.Amount+blkDeposit2.Amount, bal2)
}
