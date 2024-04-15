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

package beacondb_test

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	consensusprimitives "github.com/berachain/beacon-kit/mod/primitives-consensus"
	"github.com/berachain/beacon-kit/mod/primitives/bytes"
	"github.com/berachain/beacon-kit/mod/storage/beacondb"
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

	sdb := beacondb.New(
		sdkruntime.NewKVStoreService(storeKey),
	)
	sdb = sdb.WithContext(ctx)
	t.Run("should work with deposit", func(t *testing.T) {
		cred := []byte("12345678901234567890123456789012")
		deposit := &consensusprimitives.Deposit{
			Pubkey: primitives.BLSPubkey(
				bytes.ToBytes48([]byte("pubkey")),
			),
			Credentials: consensusprimitives.WithdrawalCredentials(cred),
			Amount:      100,
			Signature: primitives.BLSSignature(
				bytes.ToBytes96([]byte("signature")),
			),
		}
		err := sdb.EnqueueDeposits(consensusprimitives.Deposits{deposit})
		require.NoError(t, err)
		deposits, err := sdb.DequeueDeposits(1)
		require.NoError(t, err)
		require.Equal(t, deposit, deposits[0])
	})
}
