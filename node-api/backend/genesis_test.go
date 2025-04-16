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
//go:build test
// +build test

package backend_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/storage/beacondb"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	cmtcfg "github.com/cometbft/cometbft/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/require"
)

func setupStateWithGenesisValues(
	t *testing.T, cms storetypes.CommitMultiStore, kvStore *beacondb.KVStore,
) {
	t.Helper()

	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), false, log.NewNopLogger())
	kvStore = kvStore.WithContext(sdkCtx)
	require.NoError(t, kvStore.SetSlot(0))
	require.NoError(t, kvStore.SetGenesisValidatorsRoot(common.Root{0x1, 0x2, 0x3}))

	//nolint:errcheck // false positive as this has no return value
	sdkCtx.MultiStore().(storetypes.CacheMultiStore).Write()
}

func TestGetGenesisData(t *testing.T) {
	t.Parallel()

	// Build backend to test
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	cms, kvStore, depositStore, err := statetransition.BuildTestStores()
	require.NoError(t, err)

	// Setup state for genesis tests.
	setupStateWithGenesisValues(t, cms, kvStore)
	sb := storage.NewBackend(cs, nil, kvStore, depositStore, nil)

	// Create a temporary directory for CometBFT config
	tmpDir := t.TempDir()

	// Create CometBFT config with temporary directory
	cmtCfg := cmtcfg.DefaultConfig()
	cmtCfg.SetRoot(tmpDir)

	// Create config directory
	configDir := filepath.Join(tmpDir, "config")
	err = os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	// Create app genesis with version of Deneb 0x04000000.
	appGenesis := genutiltypes.NewAppGenesisWithVersion("test-chain", []byte(`
	{
		"beacon": {
			"fork_version": "0x04000000"
		}
	}
	`))
	appGenesis.GenesisTime = time.Unix(1737410400, 0)

	// Save genesis file to the config directory
	genesisFile := filepath.Join(configDir, "genesis.json")
	err = appGenesis.SaveAs(genesisFile)
	require.NoError(t, err)

	b, err := backend.New(sb, cs, cmtCfg)
	require.NoError(t, err)
	tcs := &testConsensusService{
		cms:     cms,
		kvStore: kvStore,
		cs:      cs,
	}
	b.AttachQueryBackend(tcs)

	// Test all genesis data.
	genesisTime, err := b.GenesisTime()
	require.NoError(t, err)
	require.Equal(t, math.U64(1737410400), genesisTime)

	genesisForkVersion, err := b.GenesisForkVersion()
	require.NoError(t, err)
	require.Equal(t, version.Deneb(), genesisForkVersion) // Deneb 0x04000000

	genesisValidatorsRoot, err := b.GenesisValidatorsRoot()
	require.NoError(t, err)
	require.Equal(t, common.Root{0x1, 0x2, 0x3}, genesisValidatorsRoot)
}
