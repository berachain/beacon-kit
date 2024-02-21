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

package store_test

import (
	"testing"

	sdklog "cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper/store"
	"github.com/stretchr/testify/require"
)

func TestBeaconStore(t *testing.T) {
	testName := "test"
	logger := sdklog.NewNopLogger()
	keys := storetypes.NewKVStoreKeys(testName)
	cms := integration.CreateMultiStore(keys, logger)
	ctx := sdk.NewContext(cms, cmtproto.Header{}, true, logger)
	kvs := sdkruntime.NewKVStoreService(keys[testName])

	beaconStore := store.NewBeaconStore(ctx, kvs, &config.DefaultConfig().Beacon)

	safeHash := common.HexToHash("0x123")
	beaconStore.SetSafeEth1BlockHash(safeHash)
	hash := beaconStore.GetSafeEth1BlockHash()
	require.Equal(t, safeHash, hash)
	hash.SetBytes([]byte("0x789"))
	require.Equal(t, safeHash, beaconStore.GetSafeEth1BlockHash())
	newSafeHash := common.HexToHash("0x456")
	beaconStore.SetSafeEth1BlockHash(newSafeHash)
	require.Equal(t, newSafeHash, beaconStore.GetSafeEth1BlockHash())

	finalHash := common.HexToHash("0x456")
	beaconStore.SetFinalizedEth1BlockHash(finalHash)
	require.Equal(t, finalHash, beaconStore.GetFinalizedEth1BlockHash())

	genesisHash := common.HexToHash("0x789")
	beaconStore.SetGenesisEth1Hash(genesisHash)
	require.Equal(t, genesisHash, beaconStore.GenesisEth1Hash())
}
