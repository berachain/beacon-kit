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
	"github.com/berachain/beacon-kit/storage/deposit/v2"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
)

func TestSimpleInsertionAndRetrieval(t *testing.T) {
	t.Parallel()

	baseDB, err := db.OpenDB("", dbm.MemDBBackend)
	require.NoError(t, err)

	nopLog := log.NewNopLogger()
	dummyCtx := context.Background()

	var store *deposit.KVStore
	require.NotPanics(t, func() {
		store = deposit.NewStore(baseDB, nopLog)
	})

	ins := []*types.Deposit{
		{
			Pubkey:      [48]byte{0x01},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x01}),
			Amount:      2025,
			Signature:   crypto.BLSSignature{0x01},
			Index:       1,
		},
		{
			Pubkey:      [48]byte{0x02},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x02}),
			Amount:      2026,
			Signature:   crypto.BLSSignature{0x02},
			Index:       2,
		},
		{
			Pubkey:      [48]byte{0x03},
			Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{0x03}),
			Amount:      2027,
			Signature:   crypto.BLSSignature{0x03},
			Index:       3,
		},
	}

	require.NoError(t, store.EnqueueDeposits(dummyCtx, ins))

	outs, root, err := store.GetDepositsByIndex(dummyCtx, ins[0].Index, uint64(len(ins)))
	require.NoError(t, err)

	// inputs and outputs have slightly different types, so we compare them explicitly
	require.Equal(t, len(ins), len(outs))
	for i, d := range outs {
		require.Equal(t, ins[i], d)
	}

	require.NotEmpty(t, root)
	require.NoError(t, store.Close())

	// repoen the store and check that data can be retrieved again
	var newStore *deposit.KVStore
	require.NotPanics(t, func() {
		newStore = deposit.NewStore(baseDB, nopLog)
	})

	outs2, root2, err := newStore.GetDepositsByIndex(dummyCtx, ins[0].Index, uint64(len(ins)))
	require.NoError(t, err)

	// inputs and outputs have slightly different types, so we compare them explicitly
	require.Equal(t, outs, outs2)
	require.Equal(t, root, root2)
	require.NoError(t, newStore.Close())
}
