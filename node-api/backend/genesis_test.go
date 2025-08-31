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
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/backend/mocks"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
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

func TestGetGenesisData_SunnyPath(t *testing.T) {
	t.Parallel()

	var (
		genesisTime              = int64(1737410400)
		testGenesisValidatorRoot = common.Root{0x1, 0x2, 0x3}
		testFailedLoadingState   = errors.New("test failed loading state")
	)

	testCases := []struct {
		name                string
		setMockExpectations func(storetypes.CommitMultiStore, *beacondb.KVStore, *mocks.ConsensusService)
		check               func(t *testing.T, b *backend.Backend)
	}{
		{
			name: "sunny path",
			setMockExpectations: func(
				cms storetypes.CommitMultiStore,
				kvStore *beacondb.KVStore,
				tcs *mocks.ConsensusService,
			) {
				t.Helper()

				setupStateWithGenesisValues(t, cms, kvStore, testGenesisValidatorRoot)

				tcs.EXPECT().CreateQueryContext(int64(utils.Head), false).RunAndReturn(
					func(int64, bool) (sdk.Context, error) {
						sdkCtx := sdk.NewContext(cms.CacheMultiStore(), false, log.NewNopLogger())
						return sdkCtx, nil
					},
				).Once()
			},
			check: func(t *testing.T, b *backend.Backend) {
				t.Helper()

				gotGenesisTime, err := b.GenesisTime()
				require.NoError(t, err)
				require.Equal(t, math.U64(genesisTime), gotGenesisTime)

				genesisForkVersion, err := b.GenesisForkVersion()
				require.NoError(t, err)
				require.Equal(t, version.Deneb(), genesisForkVersion)

				genesisValidatorsRoot, err := b.GenesisValidatorsRoot()
				require.NoError(t, err)
				require.Equal(t, testGenesisValidatorRoot, genesisValidatorsRoot)

				// reloading genesis validator root will provide the cached value.
				// CreateQueryContext above can be called only once and will err otherwise.
				genesisValidatorsRoot, err = b.GenesisValidatorsRoot()
				require.NoError(t, err)
				require.Equal(t, testGenesisValidatorRoot, genesisValidatorsRoot)
			},
		},
		{
			name: "app not ready",
			setMockExpectations: func(
				cms storetypes.CommitMultiStore,
				kvStore *beacondb.KVStore,
				tcs *mocks.ConsensusService,
			) {
				t.Helper()

				setupStateWithGenesisValues(t, cms, kvStore, testGenesisValidatorRoot)

				tcs.EXPECT().CreateQueryContext(int64(utils.Head), false).RunAndReturn(
					func(int64, bool) (sdk.Context, error) {
						// cometbft.ErrAppNotReady signals that consensus has not state
						// i.e. genesis has not been processed yet
						return sdk.Context{}, cometbft.ErrAppNotReady
					},
				).Once()
			},
			check: func(t *testing.T, b *backend.Backend) {
				t.Helper()

				gotGenesisTime, err := b.GenesisTime()
				require.NoError(t, err)
				require.Equal(t, math.U64(genesisTime), gotGenesisTime)

				genesisForkVersion, err := b.GenesisForkVersion()
				require.NoError(t, err)
				require.Equal(t, version.Deneb(), genesisForkVersion)

				// cometbft.ErrAppNotReady does not get propagated. Instead
				// we return an empty genesis validator root.
				genesisValidatorsRoot, err := b.GenesisValidatorsRoot()
				require.NoError(t, err)
				require.Empty(t, genesisValidatorsRoot)
			},
		},
		{
			name: "failed loading state",
			setMockExpectations: func(
				cms storetypes.CommitMultiStore,
				kvStore *beacondb.KVStore,
				tcs *mocks.ConsensusService,
			) {
				t.Helper()

				setupStateWithGenesisValues(t, cms, kvStore, testGenesisValidatorRoot)

				tcs.EXPECT().CreateQueryContext(int64(utils.Head), false).RunAndReturn(
					func(int64, bool) (sdk.Context, error) {
						return sdk.Context{}, testFailedLoadingState
					},
				).Once()
			},
			check: func(t *testing.T, b *backend.Backend) {
				t.Helper()

				gotGenesisTime, err := b.GenesisTime()
				require.NoError(t, err)
				require.Equal(t, math.U64(genesisTime), gotGenesisTime)

				genesisForkVersion, err := b.GenesisForkVersion()
				require.NoError(t, err)
				require.Equal(t, version.Deneb(), genesisForkVersion)

				genesisValidatorsRoot, err := b.GenesisValidatorsRoot()
				require.ErrorIs(t, err, testFailedLoadingState)
				require.Empty(t, genesisValidatorsRoot)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cs, err := spec.MainnetChainSpec()
			require.NoError(t, err)

			cmtCfg := buildTestCometConfig(t, genesisTime)

			cms, kvStore, depositStore, err := statetransition.BuildTestStores()
			require.NoError(t, err)
			sb := storage.NewBackend(
				cs, nil, kvStore, depositStore, nil, log.NewNopLogger(), metrics.NewNoOpTelemetrySink(),
			)

			tcs := mocks.NewConsensusService(t)

			// 2- Setup expectations before backend construction
			// (loading operations are carried out in backend.New())
			tc.setMockExpectations(cms, kvStore, tcs)

			// 3- Build backend
			b, err := backend.New(sb, cs, cmtCfg)
			require.NoError(t, err)
			b.AttachQueryBackend(tcs)

			// 4- Checks
			tc.check(t, b)
		})
	}
}

func setupStateWithGenesisValues(
	t *testing.T,
	cms storetypes.CommitMultiStore,
	kvStore *beacondb.KVStore,
	testGenesisValidatorRoot common.Root,
) {
	t.Helper()

	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), false, log.NewNopLogger())
	kvStore = kvStore.WithContext(sdkCtx)
	require.NoError(t, kvStore.SetSlot(0))
	require.NoError(t, kvStore.SetGenesisValidatorsRoot(testGenesisValidatorRoot))

	//nolint:errcheck // false positive as this has no return value
	sdkCtx.MultiStore().(storetypes.CacheMultiStore).Write()
}

func buildTestCometConfig(t *testing.T, genesisTime int64) *cmtcfg.Config {
	t.Helper()

	// Create a temporary directory for CometBFT config
	tmpDir := t.TempDir()
	cmtCfg := cmtcfg.DefaultConfig()
	cmtCfg.SetRoot(tmpDir)

	// Create config directory
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	// Create app genesis with version of Deneb 0x04000000.
	appGenesis := genutiltypes.NewAppGenesisWithVersion("test-chain", []byte(`
	{
		"beacon": {
			"fork_version": "0x04000000"
		}
	}
	`))
	appGenesis.GenesisTime = time.Unix(genesisTime, 0)

	// Save genesis file to the config directory
	genesisFile := filepath.Join(configDir, "genesis.json")
	err = appGenesis.SaveAs(genesisFile)
	require.NoError(t, err)

	return cmtCfg
}
