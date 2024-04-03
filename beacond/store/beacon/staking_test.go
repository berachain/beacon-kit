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

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	beaconstore "github.com/berachain/beacon-kit/beacond/store/beacon"
	"github.com/berachain/beacon-kit/mod/config/params"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/bytes"
	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDeposits(t *testing.T) {
	testName := "test"
	logger := log.NewTestLogger(t)
	keys := storetypes.NewKVStoreKeys(testName)
	cms := integration.CreateMultiStore(keys, logger)
	ctx := sdk.NewContext(cms, true, logger)
	storeKey := keys[testName]
	kvs := sdkruntime.NewKVStoreService(storeKey)
	env := sdkruntime.NewEnvironment(kvs, logger)

	beaconStore := beaconstore.NewStore(env, &params.BeaconChainConfig{})
	beaconStore = beaconStore.WithContext(ctx)
	t.Run("should work with deposit", func(t *testing.T) {
		cred := []byte("12345678901234567890123456789012")
		deposit := &beacontypes.Deposit{
			Pubkey: primitives.BLSPubkey(
				bytes.ToBytes48([]byte("pubkey")),
			),
			Credentials: beacontypes.WithdrawalCredentials(cred),
			Amount:      100,
			Signature: primitives.BLSSignature(
				bytes.ToBytes96([]byte("signature")),
			),
		}
		err := beaconStore.EnqueueDeposits([]*beacontypes.Deposit{deposit})
		require.NoError(t, err)
		deposits, err := beaconStore.DequeueDeposits(1)
		require.NoError(t, err)
		require.Equal(t, deposit, deposits[0])
	})
}
