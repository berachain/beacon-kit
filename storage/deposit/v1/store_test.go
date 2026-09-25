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
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/storage/db"
	"github.com/berachain/beacon-kit/storage/deposit/v1"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
)

func TestPruneRemovesStartEndRange(t *testing.T) {
	baseDB, err := db.OpenDB("", dbm.MemDBBackend)
	require.NoError(t, err)

	nopLog := log.NewNopLogger()
	ctx := t.Context()

	var store *deposit.KVStore
	require.NotPanics(t, func() {
		store = deposit.NewStore(baseDB, nopLog)
	})

	deposits := make([]*types.Deposit, 0, 15)
	for i := range 15 {
		b := uint8(i % 255)
		deposits = append(deposits, &types.Deposit{
			Pubkey:      [48]byte{b},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{b}),
			Amount:      10_000,
			Signature:   crypto.BLSSignature{b},
			Index:       uint64(i),
		})
	}
	require.NoError(t, store.EnqueueDeposits(ctx, deposits))

	// Prune is documented as removing [start, end).
	require.NoError(t, store.Prune(ctx, 5, 10))

	// Deposits below start survive.
	got, _, err := store.GetDepositsByIndex(ctx, 0, 15)
	require.NoError(t, err)
	require.Len(t, got, 5)

	// Deposits at and above end survive (bug deleted [start, start+end)).
	got, _, err = store.GetDepositsByIndex(ctx, 10, 5)
	require.NoError(t, err)
	require.Len(t, got, 5)
	require.Equal(t, uint64(10), got[0].GetIndex().Unwrap())
}

func BenchmarkDepositsInsertion(b *testing.B) {
	baseDB, err := db.OpenDB("", dbm.MemDBBackend)
	require.NoError(b, err)

	nopLog := log.NewNopLogger()
	dummyCtx := b.Context()

	var store *deposit.KVStore
	require.NotPanics(b, func() {
		store = deposit.NewStore(baseDB, nopLog)
	})

	inputSize := 5_000
	inputs := make([][]*types.Deposit, 0, inputSize)
	for i := range inputSize {
		b := uint8(i % 255)
		d := []*types.Deposit{
			{ // typing just to ease up insertions
				Pubkey:      [48]byte{b},
				Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{b}),
				Amount:      10_000,
				Signature:   crypto.BLSSignature{b},
				Index:       uint64(i),
			},
		}
		inputs = append(inputs, d)
	}
	var root common.Root

	b.ResetTimer()
	for range b.N {
		for i, d := range inputs {
			require.NoError(b, store.EnqueueDeposits(dummyCtx, d))

			// this is why v1 is less efficient than v2. To get the historical root we need to
			// load all depoisits from genesis and hash them together. Here I allocate up to i+1
			// so as least as possible.
			_, root, err = store.GetDepositsByIndex(dummyCtx, constants.FirstDepositIndex, uint64(i+1))
			require.NoError(b, err)
			_ = root // an attempt to avoid compiler optimizations
		}
	}
}
