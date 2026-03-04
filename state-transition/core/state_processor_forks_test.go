//go:build test

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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core_test

import (
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

// setupElectra1Chain creates a chain spec where Electra2 is activated at a future timestamp,
// so the chain starts in Electra1. This is used to test the Electra1 -> Electra2 fork transition.
func setupElectra1Chain(t *testing.T, electra2ForkTime uint64) chain.Spec {
	t.Helper()
	specData := spec.DevnetChainSpecData()
	// Keep all forks at genesis except Electra2.
	specData.Electra2ForkTime = electra2ForkTime
	cs, err := chain.NewSpec(specData)
	require.NoError(t, err)
	return cs
}

// TestElectra2ForkCatchupDeposits verifies that during the Electra2 fork upgrade,
// any deposits remaining in the deposit store (fetched from EL but not yet included
// in a beacon block) are properly caught up and applied to state.
//
//nolint:paralleltest // uses envars
func TestElectra2ForkCatchupDeposits(t *testing.T) {
	// Electra2 activates at timestamp 20, so blocks with timestamp < 20 are Electra1.
	electra2ForkTime := uint64(20)
	cs := setupElectra1Chain(t, electra2ForkTime)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance       = cs.MaxEffectiveBalance()
		emptyCredentials = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// STEP 0: Setup genesis with 2 validators (pre-Electra2).
	genDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: emptyCredentials,
			Amount:      maxBalance,
			Index:       uint64(0),
		},
		{
			Pubkey:      [48]byte{0x02},
			Credentials: emptyCredentials,
			Amount:      maxBalance,
			Index:       uint64(1),
		},
	}
	genPayloadHeader := &types.ExecutionPayloadHeader{
		Versionable: types.NewVersionable(cs.GenesisForkVersion()),
	}

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	valDiff, err := sp.InitializeBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		cs.GenesisForkVersion(),
	)
	require.NoError(t, err)
	require.Len(t, valDiff, len(genDeposits))

	// Verify initial eth1DepositIndex matches genesis deposits count.
	eth1DepIndex, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)), eth1DepIndex)

	// STEP 1: Build a few blocks in Electra1 (timestamp < electra2ForkTime).
	// Add a deposit via block body, like the pre-Electra2 flow.
	newDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x03},
		Credentials: emptyCredentials,
		Amount:      maxBalance,
		Index:       uint64(len(genDeposits)),
	}
	totalDepositsCount := uint64(len(genDeposits)) + 1
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), []*types.Deposit{newDeposit}))

	var depRoot common.Root
	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), constants.FirstDepositIndex, totalDepositsCount)
	require.NoError(t, err)

	// Build block at timestamp 10 (pre-Electra2), including the new deposit.
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10, // timestamp before Electra2
		[]*types.Deposit{newDeposit},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)

	valDiff, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)
	require.Empty(t, valDiff)

	// Verify the new deposit was processed.
	eth1DepIndex, err = st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, totalDepositsCount, eth1DepIndex)

	// STEP 2: Simulate deposits fetched from EL but not yet included in a beacon block.
	// These represent the "gap" deposits that would be lost without the catchup mechanism.
	catchupDeposit1 := &types.Deposit{
		Pubkey:      [48]byte{0x04},
		Credentials: emptyCredentials,
		Amount:      maxBalance,
		Index:       totalDepositsCount,
	}
	catchupDeposit2 := &types.Deposit{
		Pubkey:      [48]byte{0x05},
		Credentials: emptyCredentials,
		Amount:      maxBalance,
		Index:       totalDepositsCount + 1,
	}
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), []*types.Deposit{catchupDeposit1, catchupDeposit2}))
	totalDepositsCount += 2

	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), constants.FirstDepositIndex, totalDepositsCount)
	require.NoError(t, err)

	// Verify the catchup deposits are in the store but NOT yet in state.
	eth1DepIndex, err = st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(3), eth1DepIndex) // only genesis + block1 deposit applied

	// Verify the new validator pubkeys are NOT yet known.
	_, err = st.ValidatorIndexByPubkey(catchupDeposit1.Pubkey)
	require.Error(t, err)
	_, err = st.ValidatorIndexByPubkey(catchupDeposit2.Pubkey)
	require.Error(t, err)

	// STEP 3: Build the first Electra2 block (timestamp >= electra2ForkTime).
	// The Electra2 fork upgrade in Transition -> ProcessFork -> upgradeToElectra2
	// should catch up the 2 queued deposits during the fork transition.
	blk2 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		math.U64(electra2ForkTime), // exactly at fork activation
		[]*types.Deposit{},         // no deposits in block body (Electra2 uses EIP-6110)
		&types.ExecutionRequests{},  // no execution requests for simplicity
		st.EVMInflationWithdrawal(math.U64(electra2ForkTime)),
	)

	valDiff, err = sp.Transition(ctx, st, blk2)
	require.NoError(t, err)

	// STEP 4: Verify the catchup deposits were applied to state.
	eth1DepIndex, err = st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, totalDepositsCount, eth1DepIndex) // all deposits should now be processed

	// Verify the new validators exist in state.
	idx1, err := st.ValidatorIndexByPubkey(catchupDeposit1.Pubkey)
	require.NoError(t, err)
	balance1, err := st.GetBalance(idx1)
	require.NoError(t, err)
	require.Equal(t, maxBalance, balance1)

	idx2, err := st.ValidatorIndexByPubkey(catchupDeposit2.Pubkey)
	require.NoError(t, err)
	balance2, err := st.GetBalance(idx2)
	require.NoError(t, err)
	require.Equal(t, maxBalance, balance2)
}

// TestElectra2ForkCatchupDepositsEmpty verifies that the Electra2 fork upgrade
// handles the case where there are no queued deposits gracefully.
//
//nolint:paralleltest // uses envars
func TestElectra2ForkCatchupDepositsEmpty(t *testing.T) {
	// Electra2 activates at timestamp 20.
	electra2ForkTime := uint64(20)
	cs := setupElectra1Chain(t, electra2ForkTime)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance       = cs.MaxEffectiveBalance()
		emptyCredentials = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup genesis.
	genDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: emptyCredentials,
			Amount:      maxBalance,
			Index:       uint64(0),
		},
	}
	genPayloadHeader := &types.ExecutionPayloadHeader{
		Versionable: types.NewVersionable(cs.GenesisForkVersion()),
	}

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	totalDepositsCount := uint64(len(genDeposits))
	_, depRoot, err := ds.GetDepositsByIndex(ctx.ConsensusCtx(), constants.FirstDepositIndex, totalDepositsCount)
	require.NoError(t, err)

	// Build one Electra1 block.
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)
	_, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)

	// No additional deposits in store. Build the first Electra2 block.
	blk2 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		math.U64(electra2ForkTime),
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(math.U64(electra2ForkTime)),
	)

	// Should succeed without error - no catchup deposits to process.
	_, err = sp.Transition(ctx, st, blk2)
	require.NoError(t, err)

	// eth1DepositIndex should remain unchanged.
	eth1DepIndex, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, totalDepositsCount, eth1DepIndex)
}

// TestElectra2ForkCatchupDepositsTopUp verifies that the Electra2 fork upgrade
// correctly handles catchup deposits that top up an existing validator's balance
// rather than creating a new validator.
//
//nolint:paralleltest // uses envars
func TestElectra2ForkCatchupDepositsTopUp(t *testing.T) {
	electra2ForkTime := uint64(20)
	cs := setupElectra1Chain(t, electra2ForkTime)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance       = cs.MaxEffectiveBalance()
		increment        = cs.EffectiveBalanceIncrement()
		emptyCredentials = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup genesis with one validator.
	genDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: emptyCredentials,
			Amount:      maxBalance - 2*increment,
			Index:       uint64(0),
		},
	}
	genPayloadHeader := &types.ExecutionPayloadHeader{
		Versionable: types.NewVersionable(cs.GenesisForkVersion()),
	}

	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	totalDepositsCount := uint64(len(genDeposits))
	_, depRoot, err := ds.GetDepositsByIndex(ctx.ConsensusCtx(), constants.FirstDepositIndex, totalDepositsCount)
	require.NoError(t, err)

	// Build one Electra1 block.
	blk1 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)
	_, err = sp.Transition(ctx, st, blk1)
	require.NoError(t, err)

	// Record initial balance.
	idx, err := st.ValidatorIndexByPubkey(genDeposits[0].Pubkey)
	require.NoError(t, err)
	initialBalance, err := st.GetBalance(idx)
	require.NoError(t, err)

	// Add a top-up deposit for the existing validator to the deposit store
	// (simulating a deposit fetched from EL but not yet included).
	topUpAmount := math.Gwei(increment)
	topUpDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01}, // same validator
		Credentials: emptyCredentials,
		Amount:      topUpAmount,
		Index:       totalDepositsCount,
	}
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), []*types.Deposit{topUpDeposit}))
	totalDepositsCount++

	_, depRoot, err = ds.GetDepositsByIndex(ctx.ConsensusCtx(), constants.FirstDepositIndex, totalDepositsCount)
	require.NoError(t, err)

	// Build the first Electra2 block - the catchup should apply the top-up deposit.
	blk2 := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		math.U64(electra2ForkTime),
		[]*types.Deposit{},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(math.U64(electra2ForkTime)),
	)

	_, err = sp.Transition(ctx, st, blk2)
	require.NoError(t, err)

	// Verify the top-up was applied.
	newBalance, err := st.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, initialBalance+topUpAmount, newBalance)

	// Verify eth1DepositIndex advanced.
	eth1DepIndex, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, totalDepositsCount, eth1DepIndex)
}
