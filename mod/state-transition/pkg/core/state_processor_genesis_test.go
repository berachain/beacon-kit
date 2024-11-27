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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	cryptomocks "github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	// Create state processor to test
	cs := spec.BetnetChainSpec()
	execEngine := mocks.NewExecutionEngine[
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		engineprimitives.Withdrawals,
	](t)
	mocksSigner := &cryptomocks.BLSSigner{}

	kvStore, depositStore, err := initTestStores()
	require.NoError(t, err)
	beaconState := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	sp := createStateProcessor(
		cs,
		execEngine,
		depositStore,
		mocksSigner,
		dummyProposerAddressVerifier,
	)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
	)

	// create test inputs
	var (
		genDeposits = []*types.Deposit{
			{
				Pubkey: [48]byte{0x01},
				Amount: maxBalance,
				Index:  uint64(0),
			},
			{
				Pubkey: [48]byte{0x02},
				Amount: minBalance + increment,
				Index:  uint64(1),
			},
			{
				Pubkey: [48]byte{0x03},
				Amount: minBalance,
				Index:  uint64(2),
			},
			{
				Pubkey: [48]byte{0x04},
				Amount: 2 * maxBalance,
				Index:  uint64(3),
			},
			{
				Pubkey: [48]byte{0x05},
				Amount: minBalance - increment,
				Index:  uint64(4),
			},
			{
				Pubkey: [48]byte{0x06},
				Amount: minBalance + increment*3/2,
				Index:  uint64(5),
			},
			{
				Pubkey: [48]byte{0x07},
				Amount: maxBalance + increment/10,
				Index:  uint64(6),
			},
			{
				Pubkey: [48]byte{0x08},
				Amount: minBalance + increment*99/100,
				Index:  uint64(7),
			},
		}
		goodDeposits = []*types.Deposit{
			genDeposits[0], genDeposits[1], genDeposits[3],
			genDeposits[5], genDeposits[6],
		}
		executionPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		fork                   = &types.Fork{
			PreviousVersion: version.FromUint32[common.Version](version.Deneb),
			CurrentVersion:  version.FromUint32[common.Version](version.Deneb),
			Epoch:           math.Epoch(constants.GenesisEpoch),
		}
	)

	// define mocks expectations
	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	// run test
	genVals, err := sp.InitializePreminedBeaconStateFromEth1(
		beaconState,
		genDeposits,
		executionPayloadHeader,
		fork.CurrentVersion,
	)

	// check outputs
	require.NoError(t, err)
	require.Len(t, genVals, len(goodDeposits))

	// check beacon state changes
	resSlot, err := beaconState.GetSlot()
	require.NoError(t, err)
	require.Equal(t, math.Slot(0), resSlot)

	resFork, err := beaconState.GetFork()
	require.NoError(t, err)
	require.Equal(t, fork, resFork)

	for _, dep := range goodDeposits {
		checkValidatorNonBartio(t, cs, beaconState, dep)
	}

	// check that validator index is duly set
	latestValIdx, err := beaconState.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)-1), latestValIdx)
}

func checkValidatorNonBartio(
	t *testing.T,
	cs chain.Spec[
		common.DomainType,
		math.Epoch,
		common.ExecutionAddress,
		math.Slot,
		any,
	],
	bs *TestBeaconStateT,
	dep *types.Deposit,
) {
	t.Helper()

	// checks on validators common to all networks
	commonChecksValidators(t, cs, bs, dep)

	// checks on validators for any network but Bartio
	idx, err := bs.ValidatorIndexByPubkey(dep.Pubkey)
	require.NoError(t, err)

	valBal, err := bs.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, dep.Amount, valBal)
}

func TestInitializeBartio(t *testing.T) {
	// Create state processor to test
	cs := spec.TestnetChainSpec()
	execEngine := mocks.NewExecutionEngine[
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		engineprimitives.Withdrawals,
	](t)
	mocksSigner := &cryptomocks.BLSSigner{}

	kvStore, depositStore, err := initTestStores()
	require.NoError(t, err)
	beaconState := new(TestBeaconStateT).NewFromDB(kvStore, cs)

	sp := createStateProcessor(
		cs,
		execEngine,
		depositStore,
		mocksSigner,
		dummyProposerAddressVerifier,
	)

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
		increment  = math.Gwei(cs.EffectiveBalanceIncrement())
		minBalance = math.Gwei(cs.EjectionBalance())
	)

	// create test inputs
	var (
		genDeposits = []*types.Deposit{
			{
				Pubkey: [48]byte{0x01},
				Amount: maxBalance,
				Index:  uint64(0),
			},
			{
				Pubkey: [48]byte{0x02},
				Amount: minBalance + increment,
				Index:  uint64(1),
			},
			{
				Pubkey: [48]byte{0x03},
				Amount: minBalance,
				Index:  uint64(2),
			},
			{
				Pubkey: [48]byte{0x04},
				Amount: 2 * maxBalance,
				Index:  uint64(3),
			},
			{
				Pubkey: [48]byte{0x05},
				Amount: minBalance - increment,
				Index:  uint64(4),
			},
			{
				Pubkey: [48]byte{0x06},
				Amount: minBalance + increment*3/2,
				Index:  uint64(5),
			},
			{
				Pubkey: [48]byte{0x07},
				Amount: maxBalance + increment/10,
				Index:  uint64(6),
			},
			{
				Pubkey: [48]byte{0x08},
				Amount: minBalance + increment*99/100,
				Index:  uint64(7),
			},
		}
		goodDeposits = []*types.Deposit{
			genDeposits[0], genDeposits[1], genDeposits[3],
			genDeposits[5], genDeposits[6],
		}
		executionPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		fork                   = &types.Fork{
			PreviousVersion: version.FromUint32[common.Version](version.Deneb),
			CurrentVersion:  version.FromUint32[common.Version](version.Deneb),
			Epoch:           math.Epoch(constants.GenesisEpoch),
		}
	)

	// define mocks expectations
	mocksSigner.On(
		"VerifySignature",
		mock.Anything, mock.Anything, mock.Anything,
	).Return(nil)

	// run test
	genVals, err := sp.InitializePreminedBeaconStateFromEth1(
		beaconState,
		genDeposits,
		executionPayloadHeader,
		fork.CurrentVersion,
	)

	// check outputs
	require.NoError(t, err)
	require.Len(t, genVals, len(goodDeposits))

	// check beacon state changes
	resSlot, err := beaconState.GetSlot()
	require.NoError(t, err)
	require.Equal(t, math.Slot(0), resSlot)

	resFork, err := beaconState.GetFork()
	require.NoError(t, err)
	require.Equal(t, fork, resFork)

	for _, dep := range goodDeposits {
		checkValidatorBartio(t, cs, beaconState, dep)
	}

	// check that validator index is duly set
	latestValIdx, err := beaconState.GetEth1DepositIndex()
	require.NoError(t, err)
	require.Equal(t, uint64(len(genDeposits)-1), latestValIdx)
}

func checkValidatorBartio(
	t *testing.T,
	cs chain.Spec[
		common.DomainType,
		math.Epoch,
		common.ExecutionAddress,
		math.Slot,
		any,
	],
	bs *TestBeaconStateT,
	dep *types.Deposit,
) {
	t.Helper()

	// checks on validators common to all networks
	commonChecksValidators(t, cs, bs, dep)

	// Bartio specific checks on validators
	idx, err := bs.ValidatorIndexByPubkey(dep.Pubkey)
	require.NoError(t, err)
	val, err := bs.ValidatorByIndex(idx)
	require.NoError(t, err)

	valBal, err := bs.GetBalance(idx)
	require.NoError(t, err)
	require.Equal(t, val.EffectiveBalance, valBal)
}

func commonChecksValidators(
	t *testing.T,
	cs chain.Spec[
		common.DomainType,
		math.Epoch,
		common.ExecutionAddress,
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

	var (
		maxBalance = math.Gwei(cs.MaxEffectiveBalance())
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

// in genesis UTs we don't need to verify proposer address
// (no one proposes genesis), hence the dummy implementation.
func dummyProposerAddressVerifier(bytes.B48) ([]byte, error) {
	return nil, nil
}
