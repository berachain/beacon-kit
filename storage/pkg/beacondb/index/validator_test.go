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

package index_test

import (
	"testing"
)

func TestValidatorIndexes(_ *testing.T) {
	// testName := "test"
	// logger := log.NewTestLogger(t)
	// keys := storetypes.NewKVStoreKeys(testName)
	// cms := integration.CreateMultiStore(keys, logger)
	// ctx := sdk.NewContext(cms, true, logger)
	// storeKey := keys[testName]
	// kvs := sdkruntime.NewKVStoreService(storeKey)
	// env := sdkruntime.NewEnvironment(kvs, logger)

	// beaconStore := beaconstore.NewStore(env)
	// // beaconStore = beaconStore.WithContext(ctx)

	// t.Run("add validator and replace its pubkey", func(t *testing.T) {
	// 	err := beaconStore.AddValidator(ctx, []byte("pubkey"))
	// 	require.NoError(t, err)

	// 	err = beaconStore.AddValidator(ctx, []byte("pubkey2"))
	// 	require.NoError(t, err)

	// 	// get the index
	// 	index, err := beaconStore.ValidatorIndexByPubkey([]byte("pubkey2"))
	// 	require.NoError(t, err)
	// 	require.Equal(t, uint64(1), index)

	// 	err = beaconStore.UpdateValidator(
	// 		ctx,
	// 		[]byte("pubkey2"),
	// 		[]byte("newpubkey"),
	// 	)
	// 	require.NoError(t, err)

	// 	// get the index again, it should be the same as before
	// 	index, err = beaconStore.ValidatorIndexByPubkey(
	// 		[]byte("newpubkey"),
	// 	)
	// 	require.NoError(t, err)
	// 	require.Equal(t, uint64(1), index)
	// })

	// t.Run("add the same validator twice", func(t *testing.T) {
	// 	err := beaconStore.AddValidator(ctx, []byte("pubkeyA"))
	// 	require.NoError(t, err)

	// 	err = beaconStore.AddValidator(ctx, []byte("pubkeyA"))
	// 	require.Error(t, err)
	// })
}
