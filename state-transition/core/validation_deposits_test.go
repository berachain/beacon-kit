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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // uses envars
func TestInvalidDeposits(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		minBalance = math.Gwei(
			cs.EjectionBalance() +
				cs.EffectiveBalanceIncrement(),
		)
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		credentials0 = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup initial state with one validator
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: credentials0,
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

	// Create the correct deposit for pubkey 1.
	correctDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: credentials0,
		Amount:      minBalance,
		Index:       1,
	}

	// Create an invalid deposit with extra balance going to pubkey 1
	invalidDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: credentials0,
		Amount:      maxBalance, // Invalid - should be minBalance
		Index:       1,
	}

	// Create test block with invalid deposit, BUT the correct deposit for pubkey 1.
	depRoot := append(genDeposits, correctDeposit).HashTreeRoot()
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{invalidDeposit},
		st.EVMInflationWithdrawal(10),
	)

	// Add correct deposit to local store (honest validator will see this locally).
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), types.Deposits{correctDeposit}))

	// Run transition - should fail due to invalid deposit amount.
	_, err = sp.Transition(ctx, st, blk)
	require.Error(t, err)
	require.ErrorContains(t, err, "deposit mismatched")
}

//nolint:paralleltest // uses envars
func TestInvalidDepositsCount(t *testing.T) {
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		credentials0 = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup initial state with one validator
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: credentials0,
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

	// Create the correct deposits.
	correctDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: credentials0,
			Amount:      maxBalance,
			Index:       1,
		},
		{
			Pubkey:      [48]byte{0x02},
			Credentials: credentials0,
			Amount:      maxBalance,
			Index:       2,
		},
	}

	// Create test block with the correct deposits.
	depRoot := append(genDeposits, correctDeposits...).HashTreeRoot()
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		correctDeposits,
		st.EVMInflationWithdrawal(10),
	)

	// Add JUST 1 correct deposit to local store. This node SHOULD fail to verify.
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), types.Deposits{correctDeposits[0]}))

	// Run transition.
	_, err = sp.Transition(ctx, st, blk)
	require.Error(t, err)
	require.ErrorContains(t, err, "deposits lengths mismatched")
}

func TestLocalDepositsExceedBlockDeposits(t *testing.T) {
	t.Parallel()
	csData := spec.DevnetChainSpecData()
	csData.MaxDepositsPerBlock = 1 // Set only 1 deposit allowed per block.
	cs, err := chain.NewSpec(csData)
	require.NoError(t, err)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		credentials0 = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup initial state with one validator
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: credentials0,
				Amount:      maxBalance,
				Index:       0,
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

	// Create the block deposits.
	blockDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: credentials0,
			Amount:      maxBalance,
			Index:       1,
		},
	}

	// Create test block with the correct deposits.
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

	extraLocalDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: credentials0,
		Amount:      maxBalance,
		Index:       2,
	}

	// Add both deposits to local store (which includes more than what's in the block).
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), append(blockDeposits, extraLocalDeposit)))

	// Run transition.
	_, err = sp.Transition(ctx, st, blk)
	require.NoError(t, err)
}

func TestLocalDepositsExceedBlockDepositsBadRoot(t *testing.T) {
	t.Parallel()
	csData := spec.DevnetChainSpecData()
	csData.MaxDepositsPerBlock = 1 // Set only 1 deposit allowed per block.
	cs, err := chain.NewSpec(csData)
	require.NoError(t, err)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	var (
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		credentials0 = types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})
	)

	// Setup initial state with one validator
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: credentials0,
				Amount:      maxBalance,
				Index:       0,
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

	// Create the block deposits.
	blockDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: credentials0,
			Amount:      maxBalance,
			Index:       1,
		},
	}

	extraLocalDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: credentials0,
		Amount:      maxBalance,
		Index:       2,
	}

	// Now, the block proposer ends up adding the correct 1 deposit per block, BUT spoofs the
	// deposits root to use the entire deposits list.
	badDepRoot := append(genDeposits, append(blockDeposits, extraLocalDeposit)...).HashTreeRoot()
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(badDepRoot),
		10,
		blockDeposits,
		st.EVMInflationWithdrawal(10),
	)

	// Add both deposits to local store (which includes more than what's in the block).
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), append(blockDeposits, extraLocalDeposit)))

	// Run transition.
	_, err = sp.Transition(ctx, st, blk)
	require.Error(t, err)
	require.ErrorContains(t, err, "deposits root mismatch")
}
