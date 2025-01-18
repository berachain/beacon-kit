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

	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestInvalidDeposits(t *testing.T) {
	cs := setupChain(t, components.BoonetChainSpecType)
	sp, st, ds, ctx := setupState(t, cs)

	var (
		minBalance   = math.Gwei(cs.EjectionBalance() + cs.EffectiveBalanceIncrement())
		maxBalance   = math.Gwei(cs.MaxEffectiveBalance())
		credentials0 = types.NewCredentialsFromExecutionAddress(
			common.ExecutionAddress{},
		)
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
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)
	require.NoError(t, ds.EnqueueDeposits(ctx, genDeposits))
	_, err := sp.InitializePreminedBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, genVersion,
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
			Eth1Data: types.NewEth1Data(depRoot),
			Deposits: []*types.Deposit{invalidDeposit},
		},
	)

	// Add correct deposit to local store (honest validator will see this locally).
	require.NoError(t, ds.EnqueueDeposits(ctx, types.Deposits{correctDeposit}))

	// Run transition - should fail due to invalid deposit amount.
	_, err = sp.Transition(ctx, st, blk)
	require.Error(t, err)
	require.ErrorContains(t, err, "deposit mismatched")
}
