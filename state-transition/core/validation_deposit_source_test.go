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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/state-transition/core"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

// TestDepositSourceRejectsEIP6110DepositsBeforeFulu verifies that the consensus layer
// rejects a pre-Fulu block carrying EIP-6110 deposit requests. Before Fulu, deposits
// must be sourced exclusively from the deposit-contract events on the block body.
//
//nolint:paralleltest // uses envars
func TestDepositSourceRejectsEIP6110DepositsBeforeFulu(t *testing.T) {
	cs := setupPreFuluChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	credentials := types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})

	// Setup initial state with one genesis validator.
	genDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x00},
			Credentials: credentials,
			Amount:      cs.MaxEffectiveBalance(),
			Index:       0,
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

	// Use the genesis deposit root so the block's only defect is the stray EIP-6110 deposit.
	_, depRoot, err := ds.GetDepositsByIndex(
		ctx.ConsensusCtx(), constants.FirstDepositIndex, uint64(len(genDeposits)),
	)
	require.NoError(t, err)

	// Build a pre-Fulu block with no body deposits but a stray EIP-6110 deposit request.
	strayDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: credentials,
		Amount:      cs.MinActivationBalance(),
		Index:       1,
	}
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{},
		&types.ExecutionRequests{Deposits: []*types.Deposit{strayDeposit}},
		st.EVMInflationWithdrawal(10),
	)

	_, err = sp.Transition(ctx, st, blk)
	require.ErrorIs(t, err, core.ErrUnexpectedDepositSource)
}

// TestDepositSourceRejectsBodyDepositsAfterFulu verifies that the consensus layer
// rejects a post-Fulu block carrying beacon block body deposits. From Fulu onwards
// (after the first Fulu block) deposits must be sourced exclusively from EIP-6110
// execution requests.
//
//nolint:paralleltest // uses envars
func TestDepositSourceRejectsBodyDepositsAfterFulu(t *testing.T) {
	// The devnet chain spec activates every fork (including Fulu) at genesis, so any
	// block built here is a Fulu block whose predecessor is also Fulu, i.e. not the
	// first Fulu block.
	cs := setupChain(t)
	sp, st, ds, ctx, _, _ := statetransition.SetupTestState(t, cs)

	credentials := types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{})

	// Setup initial state with one genesis validator.
	genDeposits := types.Deposits{
		{
			Pubkey:      [48]byte{0x00},
			Credentials: credentials,
			Amount:      cs.MaxEffectiveBalance(),
			Index:       0,
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

	_, depRoot, err := ds.GetDepositsByIndex(
		ctx.ConsensusCtx(), constants.FirstDepositIndex, uint64(len(genDeposits)),
	)
	require.NoError(t, err)

	// Build a post-Fulu block carrying a deposit on the block body, which is no longer
	// a valid deposit source once EIP-6110 is active.
	strayDeposit := &types.Deposit{
		Pubkey:      [48]byte{0x01},
		Credentials: credentials,
		Amount:      cs.MinActivationBalance(),
		Index:       1,
	}
	blk := buildNextBlock(
		t,
		cs,
		st,
		types.NewEth1Data(depRoot),
		10,
		[]*types.Deposit{strayDeposit},
		&types.ExecutionRequests{},
		st.EVMInflationWithdrawal(10),
	)

	_, err = sp.Transition(ctx, st, blk)
	require.ErrorIs(t, err, core.ErrUnexpectedDepositSource)
}
