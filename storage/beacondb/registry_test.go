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

package beacondb_test

import (
	"context"
	"fmt"
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/db"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type testKVStoreService struct {
	ctx sdk.Context
}

func (kvs *testKVStoreService) OpenKVStore(context.Context) corestore.KVStore {
	//nolint:contextcheck // fine with tests
	store := sdk.UnwrapSDKContext(kvs.ctx).KVStore(testStoreKey)
	return storage.NewKVStore(store)
}

var testStoreKey = storetypes.NewKVStoreKey("storage-tests")

func TestBalances(t *testing.T) {
	t.Parallel()
	store, err := initTestStore()
	require.NoError(t, err)

	// no balance to start
	res, err := store.GetBalances()
	require.NoError(t, err)
	require.Empty(t, res)

	// add balances
	var (
		idx1, idx2     = math.U64(1_987), math.U64(1_989)
		inBal1, inBal2 = math.U64(8_992), math.U64(10_000)
	)
	require.NoError(t, store.SetBalance(idx1, inBal1))
	require.NoError(t, store.SetBalance(idx2, inBal2))

	// check we can query added balances
	outBal, err := store.GetBalance(idx1)
	require.NoError(t, err)
	require.Equal(t, inBal1, outBal)

	outBal, err = store.GetBalance(idx2)
	require.NoError(t, err)
	require.Equal(t, inBal2, outBal)

	res, err = store.GetBalances()
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, inBal1.Unwrap(), res[0])
	require.Equal(t, inBal2.Unwrap(), res[1])

	// update existing balances
	newInBal1, newInBal2 := math.U64(0), inBal2*2
	require.NoError(t, store.SetBalance(idx1, newInBal1))
	require.NoError(t, store.SetBalance(idx2, newInBal2))

	// check we can query updated balances
	outBal, err = store.GetBalance(idx1)
	require.NoError(t, err)
	require.Equal(t, newInBal1, outBal)

	outBal, err = store.GetBalance(idx2)
	require.NoError(t, err)
	require.Equal(t, newInBal2, outBal)

	res, err = store.GetBalances()
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, newInBal1.Unwrap(), res[0])
	require.Equal(t, newInBal2.Unwrap(), res[1])
}

func TestValidators(t *testing.T) {
	t.Parallel()
	store, err := initTestStore()
	require.NoError(t, err)

	// no validators to start
	res, err := store.GetValidators()
	require.NoError(t, err)
	require.Empty(t, res)

	// add validators
	var (
		inVal1 = &types.Validator{
			Pubkey:           bytes.B48{0x01},
			EffectiveBalance: 31e9,
		}
		inVal2 = &types.Validator{
			Pubkey:           bytes.B48{0x02},
			EffectiveBalance: 32e9,
		}
	)
	require.NoError(t, store.AddValidator(inVal1))
	require.NoError(t, store.AddValidator(inVal2))

	// check we can query added validators
	valIdx1, err := store.ValidatorIndexByPubkey(inVal1.GetPubkey())
	require.NoError(t, err)
	outVal, err := store.ValidatorByIndex(valIdx1)
	require.NoError(t, err)
	require.Equal(t, inVal1, outVal)

	valIdx2, err := store.ValidatorIndexByPubkey(inVal2.GetPubkey())
	require.NoError(t, err)
	outVal, err = store.ValidatorByIndex(valIdx2)
	require.NoError(t, err)
	require.Equal(t, inVal2, outVal)

	valCount, err := store.GetTotalValidators()
	require.NoError(t, err)
	require.Equal(t, uint64(2), valCount)

	res, err = store.GetValidators()
	require.NoError(t, err)
	require.Len(t, res, int(valCount))
	require.Equal(t, inVal1, res[0])
	require.Equal(t, inVal2, res[1])

	// update existing validators balances
	var (
		inUpdatedVal1 = &types.Validator{
			Pubkey:           inVal1.GetPubkey(),
			EffectiveBalance: inVal1.EffectiveBalance * 2,
		}
		inUpdatedVal2 = &types.Validator{
			Pubkey:           inVal2.GetPubkey(),
			EffectiveBalance: inVal1.EffectiveBalance / 2,
		}
	)
	require.NoError(t, store.UpdateValidatorAtIndex(valIdx1, inUpdatedVal1))
	require.NoError(t, store.UpdateValidatorAtIndex(valIdx2, inUpdatedVal2))

	// check we can query updated validators
	upValIdx1, err := store.ValidatorIndexByPubkey(inVal1.GetPubkey())
	require.NoError(t, err)
	require.Equal(t, valIdx1, upValIdx1)
	outVal, err = store.ValidatorByIndex(upValIdx1)
	require.NoError(t, err)
	require.Equal(t, inUpdatedVal1, outVal)

	upValIdx2, err := store.ValidatorIndexByPubkey(inVal2.GetPubkey())
	require.NoError(t, err)
	require.Equal(t, valIdx2, upValIdx2)
	outVal, err = store.ValidatorByIndex(upValIdx2)
	require.NoError(t, err)
	require.Equal(t, inUpdatedVal2, outVal)

	upValCount, err := store.GetTotalValidators()
	require.NoError(t, err)
	require.Equal(t, valCount, upValCount)

	res, err = store.GetValidators()
	require.NoError(t, err)
	require.Len(t, res, int(valCount))
	require.Equal(t, inUpdatedVal1, res[0])
	require.Equal(t, inUpdatedVal2, res[1])
}

// TestPendingPartialWithdrawals_Nil verifies that if no pending partial withdrawals
// have been set, then GetPendingPartialWithdrawals returns an error.
func TestPendingPartialWithdrawals_Nil(t *testing.T) {
	t.Parallel()
	store, err := initTestStore()
	require.NoError(t, err)

	ppw, err := store.GetPendingPartialWithdrawals()

	require.ErrorContains(t, err, "collections: not found: key 'no_key' of type SSZMarshallable")
	require.Nil(t, ppw)
}

// TestPendingPartialWithdrawals_EmptySlice verifies that when setting an empty slice, the Set and Get operations succeed.
func TestPendingPartialWithdrawals_EmptySlice(t *testing.T) {
	t.Parallel()
	store, err := initTestStore()
	require.NoError(t, err)

	// Attempt to set an empty slice.
	err = store.SetPendingPartialWithdrawals([]*types.PendingPartialWithdrawal{})
	var ppw []*types.PendingPartialWithdrawal
	require.NoError(t, err)
	ppw, err = store.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Empty(t, ppw)
}

// TestPendingPartialWithdrawals_SetAndGetNonEmpty verifies that setting a non-empty list of pending partial withdrawals succeeds.
func TestPendingPartialWithdrawals_SetAndGetNonEmpty(t *testing.T) {
	t.Parallel()
	store, err := initTestStore()
	require.NoError(t, err)

	// Create sample pending partial withdrawal entries.
	entry1 := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.U64(1),
		Amount:            math.U64(100),
		WithdrawableEpoch: math.U64(15),
	}
	entry2 := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.U64(2),
		Amount:            math.U64(200),
		WithdrawableEpoch: math.U64(20),
	}
	sampleWithdrawals := []*types.PendingPartialWithdrawal{entry1, entry2}

	var ppw []*types.PendingPartialWithdrawal
	// Attempt to set the non-empty slice.
	err = store.SetPendingPartialWithdrawals(sampleWithdrawals)

	require.NoError(t, err)
	ppw, err = store.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Equal(t, sampleWithdrawals, ppw)
}

// TestPendingPartialWithdrawals_Update verifies that updating the pending partial withdrawals works correctly.
func TestPendingPartialWithdrawals_Update(t *testing.T) {
	t.Parallel()
	store, err := initTestStore()
	require.NoError(t, err)

	require.NoError(t, err)
	// Create initial pending partial withdrawal entries.
	entry1 := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.U64(1),
		Amount:            math.U64(100),
		WithdrawableEpoch: math.U64(15),
	}
	entry2 := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.U64(2),
		Amount:            math.U64(200),
		WithdrawableEpoch: math.U64(20),
	}
	initialWithdrawals := []*types.PendingPartialWithdrawal{entry1, entry2}
	var ppw []*types.PendingPartialWithdrawal

	// Set the initial list.
	err = store.SetPendingPartialWithdrawals(initialWithdrawals)
	require.NoError(t, err)

	// Now update by modifying entry1's amount and dropping entry2.
	updatedEntry := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.U64(1),
		Amount:            math.U64(150), // Updated amount.
		WithdrawableEpoch: math.U64(15),
	}
	updatedWithdrawals := []*types.PendingPartialWithdrawal{updatedEntry}
	err = store.SetPendingPartialWithdrawals(updatedWithdrawals)
	require.NoError(t, err)

	ppw, err = store.GetPendingPartialWithdrawals()
	require.NoError(t, err)
	require.Equal(t, updatedWithdrawals, ppw)
}

func initTestStore() (*beacondb.KVStore, error) {
	db, err := db.OpenDB("", dbm.MemDBBackend)
	if err != nil {
		return nil, fmt.Errorf("failed opening mem db: %w", err)
	}
	var (
		nopLog     = log.NewNopLogger()
		nopMetrics = metrics.NewNoOpMetrics()
	)

	cms := store.NewCommitMultiStore(
		db,
		nopLog,
		nopMetrics,
	)

	cms.MountStoreWithDB(testStoreKey, storetypes.StoreTypeIAVL, nil)
	if err = cms.LoadLatestVersion(); err != nil {
		return nil, fmt.Errorf("failed to load latest version: %w", err)
	}

	ctx := sdk.NewContext(cms, true, nopLog)
	testStoreService := &testKVStoreService{
		ctx: ctx,
	}
	return beacondb.New(testStoreService), nil
}
