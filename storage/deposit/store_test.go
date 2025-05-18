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

package deposit_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/storage/db"
	"github.com/berachain/beacon-kit/storage/deposit"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
)

func TestDataMigration(t *testing.T) {
	t.Parallel()

	dbV1, err := db.OpenDB("testV1", dbm.MemDBBackend)
	require.NoError(t, err)

	dbV2, err := db.OpenDB("testV2", dbm.MemDBBackend)
	require.NoError(t, err)

	nopLog := log.NewNopLogger()
	dummyCtx := context.Background()

	var store deposit.StoreManager
	require.NotPanics(t, func() {
		store = deposit.NewStore(dbV1, dbV2, nopLog)
	})

	// add data to migrate
	ins := []*types.Deposit{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x01}),
			Amount:      2025,
			Signature:   crypto.BLSSignature{0x01},
			Index:       0,
		},
		{
			Pubkey:      [48]byte{0x02},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x02}),
			Amount:      2026,
			Signature:   crypto.BLSSignature{0x02},
			Index:       1,
		},
		{
			Pubkey:      [48]byte{0x03},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x03}),
			Amount:      2027,
			Signature:   crypto.BLSSignature{0x03},
			Index:       2,
		},
	}

	require.NoError(t, store.SelectVersion(deposit.V1))
	require.NoError(t, store.EnqueueDeposits(dummyCtx, ins))

	// carry out migration
	require.NoError(t, store.MigrateV1ToV2())

	require.NoError(t, store.SelectVersion(deposit.V2))
	outs, root, err := store.GetDepositsByIndex(dummyCtx, ins[0].Index, uint64(len(ins)))
	require.NoError(t, err)

	// inputs and outputs have slightly different types, so we compare them explicitly
	require.Equal(t, len(ins), len(outs))
	for i, d := range outs {
		require.Equal(t, ins[i], d)
	}
	require.NotEmpty(t, root)

	// close stores and check that when reopening them migration is still there
	require.NoError(t, store.Close())

	require.NotPanics(t, func() {
		store = deposit.NewStore(dbV1, dbV2, nopLog)
	})

	require.NoError(t, store.SelectVersion(deposit.V2))
	outs2, root2, err := store.GetDepositsByIndex(dummyCtx, ins[0].Index, uint64(len(ins)))
	require.NoError(t, err)
	require.Equal(t, outs, outs2)
	require.Equal(t, root, root2)
	require.NoError(t, store.Close())
}

// Try migrating storeV1 content to storeV2 twice and show that
// data migration is carried out only once
func TestDataMigrationIsIdempotent(t *testing.T) {
	t.Parallel()

	dbV1, err := db.OpenDB("testV1", dbm.MemDBBackend)
	require.NoError(t, err)

	dbV2, err := db.OpenDB("testV2", dbm.MemDBBackend)
	require.NoError(t, err)

	nopLog := log.NewNopLogger()
	dummyCtx := context.Background()

	var store deposit.StoreManager
	require.NotPanics(t, func() {
		store = deposit.NewStore(dbV1, dbV2, nopLog)
	})

	// add data to migrate
	ins0 := []*types.Deposit{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x01}),
			Amount:      2025,
			Signature:   crypto.BLSSignature{0x01},
			Index:       0,
		},
	}

	require.NoError(t, store.SelectVersion(deposit.V1))
	require.NoError(t, store.EnqueueDeposits(dummyCtx, ins0))
	require.NoError(t, store.MigrateV1ToV2())

	require.NoError(t, store.SelectVersion(deposit.V2))
	outs0, root0, err := store.GetDepositsByIndex(dummyCtx, ins0[0].Index, uint64(len(ins0)))
	require.NoError(t, err)
	require.Len(t, outs0, len(ins0))
	require.Equal(t, ins0[0], outs0[0])

	// try to migrate again with different data
	ins1 := []*types.Deposit{
		{
			Pubkey:      [48]byte{0xff},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0xff}),
			Amount:      2052,
			Signature:   crypto.BLSSignature{0xff},
			Index:       1987,
		},
	}
	require.NoError(t, store.SelectVersion(deposit.V1))
	require.NoError(t, store.EnqueueDeposits(dummyCtx, ins1))

	require.NoError(t, store.MigrateV1ToV2())

	// show that new data are not migrated
	require.NoError(t, store.SelectVersion(deposit.V2))
	outs1, root1, err := store.GetDepositsByIndex(dummyCtx, ins0[0].Index, uint64(len(ins0)+len(ins1)))
	require.NoError(t, err)
	require.Equal(t, outs0, outs1)
	require.Equal(t, root0, root1)
}
