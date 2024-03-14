// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package beacon_test

import (
	"testing"

	sdklog "cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	beaconstore "github.com/berachain/beacon-kit/store/beacon"
	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidatorIndexes(t *testing.T) {
	testName := "test"
	logger := sdklog.NewNopLogger()
	keys := storetypes.NewKVStoreKeys(testName)
	cms := integration.CreateMultiStore(keys, logger)
	ctx := sdk.NewContext(cms, true, logger)
	storeKey := keys[testName]
	kvs := sdkruntime.NewKVStoreService(storeKey)

	beaconStore := beaconstore.NewStore(kvs)
	beaconStore = beaconStore.WithContext(ctx)

	t.Run("add validator and replace its pubkey", func(t *testing.T) {
		err := beaconStore.AddValidator(ctx, []byte("pubkey"))
		require.NoError(t, err)

		err = beaconStore.AddValidator(ctx, []byte("pubkey2"))
		require.NoError(t, err)

		// get the index
		index, err := beaconStore.ValidatorIndexByPubkey(ctx, []byte("pubkey2"))
		require.NoError(t, err)
		require.Equal(t, uint64(1), index)

		err = beaconStore.UpdateValidator(
			ctx,
			[]byte("pubkey2"),
			[]byte("newpubkey"),
		)
		require.NoError(t, err)

		// get the index again, it should be the same as before
		index, err = beaconStore.ValidatorIndexByPubkey(
			ctx,
			[]byte("newpubkey"),
		)
		require.NoError(t, err)
		require.Equal(t, uint64(1), index)
	})

	t.Run("add the same validator twice", func(t *testing.T) {
		err := beaconStore.AddValidator(ctx, []byte("pubkeyA"))
		require.NoError(t, err)

		err = beaconStore.AddValidator(ctx, []byte("pubkeyA"))
		require.Error(t, err)
	})
}
