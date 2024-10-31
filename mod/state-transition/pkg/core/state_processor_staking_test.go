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
	)

	kvStore, err := initTestStore()
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
		}
		blkDeposits = []*types.Deposit{
			{
				Pubkey:      genDeposits[0].Pubkey,
				Credentials: emptyCredentials,
				Amount:      minBalance,
				Index:       genDeposits[0].Index,
			},
		}
	)

	// here we duly update state root, similarly to what we do in processSlot
	genBlockHeader, err := beaconState.GetLatestBlockHeader()
	require.NoError(t, err)
	genStateRoot := beaconState.HashTreeRoot()
	genBlockHeader.SetStateRoot(genStateRoot)

	blk, err := new(types.BeaconBlock).NewWithVersion(
		genBlockHeader.GetSlot()+1,
		genBlockHeader.GetProposerIndex(),
		genBlockHeader.HashTreeRoot(),
		version.Deneb,
	)
	require.NoError(t, err)
	blk.Body = &types.BeaconBlockBody{
		ExecutionPayload: &types.ExecutionPayload{
			Timestamp:    10,
			ExtraData:    []byte("testing"),
			Transactions: [][]byte{[]byte("tx1")},
			Withdrawals: []*engineprimitives.Withdrawal{
				{ // fill empty withdrawals
					Index:     0,
					Validator: 0,
					Address:   emptyAddress,
					Amount:    0,
				},
				{
					Index:     1,
					Validator: 1,
					Address:   emptyAddress,
					Amount:    0,
				},
			},
			BaseFeePerGas: math.NewU256(0),
		},
		Eth1Data: &types.Eth1Data{},
		Deposits: blkDeposits,
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
