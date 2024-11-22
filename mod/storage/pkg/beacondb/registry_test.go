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

package beacondb_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/db"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/stretchr/testify/require"
)

var (
	testStoreKey = storetypes.NewKVStoreKey("storage-tests")
	testCodec    = &encoding.SSZInterfaceCodec[*types.ExecutionPayloadHeader]{}
)

func TestBalances(t *testing.T) {
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

func initTestStore() (
	*beacondb.KVStore[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		[]*types.Validator,
	], error) {
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

	return beacondb.New[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		[]*types.Validator,
	](
		runtime.NewKVStoreService(testStoreKey),
		testCodec,
	), nil
}
