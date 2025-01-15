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

	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	cs := setupChain(t, components.BetnetChainSpecType)
	sp, st, _, _ := setupState(t, cs)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance(false))
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
	)

	// create test inputs
	var (
		genDeposits = []*types.Deposit{
			{
				Pubkey: [48]byte{0x01},
				Amount: maxBalance,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x01},
				),
				Index: uint64(0),
			},
			{
				Pubkey: [48]byte{0x02},
				Amount: minBalance + increment,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x02},
				),
				Index: uint64(1),
			},
			{
				Pubkey: [48]byte{0x03},
				Amount: minBalance,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x03},
				),
				Index: uint64(2),
			},
			{
				Pubkey: [48]byte{0x04},
				Amount: 2 * maxBalance,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x04},
				),
				Index: uint64(3),
			},
			{
				Pubkey: [48]byte{0x05},
				Amount: minBalance - increment,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x05},
				),
				Index: uint64(4),
			},
			{
				Pubkey: [48]byte{0x06},
				Amount: minBalance + increment*3/2,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x06},
				),
				Index: uint64(5),
			},
			{
				Pubkey: [48]byte{0x07},
				Amount: maxBalance + increment/10,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x07},
				),
				Index: uint64(6),
			},
			{
				Pubkey: [48]byte{0x08},
				Amount: minBalance + increment*99/100,
				Credentials: types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x08},
				),
				Index: uint64(7),
			},
		}
		goodDeposits = []*types.Deposit{
			genDeposits[0], genDeposits[1], genDeposits[3],
			genDeposits[5], genDeposits[6],
		}
		executionPayloadHeader = &types.ExecutionPayloadHeader{}
		fork                   = &types.Fork{
			PreviousVersion: version.FromUint32[common.Version](version.Deneb),
			CurrentVersion:  version.FromUint32[common.Version](version.Deneb),
			Epoch:           math.Epoch(constants.GenesisEpoch),
		}
	)

	// run test
	genVals, err := sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		executionPayloadHeader,
		fork.CurrentVersion,
	)

	// check outputs
	require.NoError(t, err)
	require.Len(t, genVals, len(goodDeposits))

	// check beacon state changes
	resSlot, err := st.GetSlot()
	require.NoError(t, err)
	require.Equal(t, math.Slot(0), resSlot)

	resFork, err := st.GetFork()
	require.NoError(t, err)
	require.Equal(t, fork, resFork)

	for _, dep := range goodDeposits {
		checkValidator(t, cs, st, dep)
	}

	// check that deposit index is duly set. On betnet
	// deposit index is set to the last accepted deposit.
	latestValIdx, err := st.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)), latestValIdx)
}

func checkValidator(
	t *testing.T,
	cs chain.Spec[
		common.DomainType,
		math.Epoch,
		math.Slot,
		any,
	],
	bs *TestBeaconStateT,
	dep *types.Deposit,
) {
	t.Helper()

	idx, err := bs.ValidatorIndexByPubkey(dep.Pubkey)
	require.NoError(t, err)

	val, err := bs.ValidatorByIndex(idx)
	require.NoError(t, err)
	require.Equal(t, dep.Pubkey, val.Pubkey)

	// checks on validators common to all networks
	commonChecksValidators(t, cs, val, dep)

	// checks on validators for any network but Bartio
	require.Equal(t, math.Epoch(0), val.GetActivationEligibilityEpoch())
	require.Equal(t, math.Epoch(0), val.GetActivationEpoch())

	valBal, err := bs.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, dep.Amount, valBal)
}

func commonChecksValidators(
	t *testing.T,
	cs chain.Spec[
		common.DomainType,
		math.Epoch,
		math.Slot,
		any,
	],
	val *types.Validator,
	dep *types.Deposit,
) {
	t.Helper()
	require.Equal(t, dep.Pubkey, val.Pubkey)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance(false))
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
	)
	switch {
	case dep.Amount >= maxBalance:
		require.Equal(t, maxBalance, val.EffectiveBalance)
	case dep.Amount > minBalance && dep.Amount < maxBalance:
		// Effective balance must be a multiple of increment.
		// If balance is not, effective balance is rounded down
		if dep.Amount%increment == 0 {
			require.Equal(t, dep.Amount, val.EffectiveBalance)
		} else {
			require.Less(t, val.EffectiveBalance, dep.Amount)
			require.Greater(t, val.EffectiveBalance, dep.Amount-increment)
			require.Zero(t, val.EffectiveBalance%increment)
		}
	case dep.Amount <= minBalance:
		require.Equal(t, math.Gwei(0), val.EffectiveBalance)
	}
}
