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

	"github.com/berachain/beacon-kit/mod/chain-spec/pkg/chain"
	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

// TestTransitionUpdateValidators shows that when validator is
// updated (increasing amount), corrensponding balance is updated.
func TestTransitionUpdateValidators(t *testing.T) {
	cs := setupChain(t, components.BoonetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		maxBalance       = math.Gwei(cs.MaxEffectiveBalance())
		minBalance       = math.Gwei(cs.EffectiveBalanceIncrement())
		emptyCredentials = types.NewCredentialsFromExecutionAddress(
			common.ExecutionAddress{},
		)
	)

	// Setup initial state via genesis
	// TODO: consider instead setting state artificially
	var (
		genDeposits = []*types.Deposit{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: emptyCredentials,
				Amount:      maxBalance - 3*minBalance,
				Index:       0,
			},
			{
				Pubkey:      [48]byte{0x01},
				Credentials: emptyCredentials,
				Amount:      maxBalance - 6*minBalance,
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
	// progressStateToSlot(t, st, math.U64(spec.BoonetFork2Height))
	// _, err = sp.ProcessSlots(st, math.U64(spec.BoonetFork2Height))
	// require.NoError(t, err)
	// blkHeader, err := st.GetLatestBlockHeader()
	// require.NoError(t, err)
	// blkHeader.SetSlot(math.U64(spec.BoonetFork2Height) - 1)
	// require.NoError(t, st.SetLatestBlockHeader(blkHeader))

	// create test inputs
	blkDeposits := []*types.Deposit{
		{
			Pubkey:      genDeposits[0].Pubkey,
			Credentials: emptyCredentials,
			Amount:      minBalance, // avoid breaching maxBalance
			Index:       uint64(len(genDeposits)),
		},
	}
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
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &types.Eth1Data{},
			Deposits: blkDeposits,
		},
	)

	// make sure included deposit is already available in deposit store
	require.NoError(t, ds.EnqueueDeposits(blkDeposits))

	// run the test
	vals, err := sp.Transition(ctx, st, blk)

	// check outputs
	require.NoError(t, err)
	require.Zero(t, vals) // just update, no new validators

	// check validator is duly updated
	expectedValBalance := genDeposits[0].Amount + blkDeposits[0].Amount
	idx, err := st.ValidatorIndexByPubkey(genDeposits[0].Pubkey)
	require.NoError(t, err)
	require.Equal(t, math.U64(genDeposits[0].Index), idx)

	val, err := st.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, genDeposits[0].Pubkey, val.Pubkey)
	require.Equal(t, expectedValBalance, val.EffectiveBalance)

	// check validator balance is updated
	valBal, err := st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, expectedValBalance, valBal)

	// check that validator index is duly set (1-indexed here, to be fixed)
	latestValIdx, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)), latestValIdx)
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
	vals, err := sp.Transition(ctx, st, blk)

	// Check outputs and ensure withdrawals in payload is consistent with
	// statedb expected withdrawals.
	require.NoError(t, err)
	require.Zero(t, vals)

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
	vals, err := sp.Transition(ctx, st, blk)

	// Check outputs and ensure withdrawals in payload is consistent with
	// statedb expected withdrawals.
	require.NoError(t, err)
	require.Zero(t, vals)

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
	vals, err = sp.Transition(ctx, st, blk)

	// Check outputs and ensure withdrawals in payload is consistent with
	// statedb expected withdrawals.
	require.NoError(t, err)
	require.Zero(t, vals)

	// Validator 1 is now withdrawn from.
	val1BalAfter, err = st.GetBalance(math.U64(1))
	require.NoError(t, err)
	require.Equal(t, maxBalance, val1BalAfter)
}
