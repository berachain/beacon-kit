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
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/backend/mocks"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	cmtcfg "github.com/cometbft/cometbft/config"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetGenesisData_SunnyPath(t *testing.T) {
	t.Parallel()

	var (
		testGenesisTime          = int64(1737410400)
		testGenesisValidatorRoot = common.Root{0x1, 0x2, 0x3}
		testFailedLoadingState   = errors.New("test failed loading state")

		expectedGenesisTime    math.U64       // to be retrieved when setting genesis state and used in checks
		expectedGenesisVersion common.Version // to be retrieved when setting genesis state and used in checks
	)

	testCases := []struct {
		name                string
		setMockExpectations func(*mocks.ConsensusService, *mocks.GenesisStateProcessor)
		check               func(t *testing.T, b *backend.Backend, errLoad error)
	}{
		{
			name: "success",
			setMockExpectations: func(
				tcs *mocks.ConsensusService,
				sp *mocks.GenesisStateProcessor,
			) {
				t.Helper()

				tcs.EXPECT().IsAppReady().Return(nil)
				sp.EXPECT().InitializeBeaconStateFromEth1(
					mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(
					func(
						st *state.StateDB,
						_ types.Deposits,
						execPayloadHeader *types.ExecutionPayloadHeader,
						genVer bytes.B4,
					) {
						expectedGenesisVersion = genVer
						expectedGenesisTime = execPayloadHeader.GetTimestamp()

						require.NoError(t, st.SetFork(&types.Fork{
							PreviousVersion: genVer,
							CurrentVersion:  genVer,
						}))
						require.NoError(t, st.SetLatestExecutionPayloadHeader(execPayloadHeader))
						require.NoError(t, st.SetGenesisValidatorsRoot(testGenesisValidatorRoot))
					},
				).Return(nil, nil)
			},
			check: func(t *testing.T, b *backend.Backend, errLoad error) {
				t.Helper()

				require.NoError(t, errLoad)

				gotGenesisTime, err := b.GenesisTime()
				require.NoError(t, err)
				require.Equal(t, expectedGenesisTime, gotGenesisTime)

				gotGenesisVersion, err := b.GenesisForkVersion()
				require.NoError(t, err)
				require.Equal(t, expectedGenesisVersion, gotGenesisVersion)

				genesisValidatorsRoot, err := b.GenesisValidatorsRoot()
				require.NoError(t, err)
				require.Equal(t, testGenesisValidatorRoot, genesisValidatorsRoot)
			},
		},
		{
			name: "app not ready",
			setMockExpectations: func(
				tcs *mocks.ConsensusService,
				_ *mocks.GenesisStateProcessor,
			) {
				t.Helper()
				tcs.EXPECT().IsAppReady().Return(cometbft.ErrAppNotReady)
			},
			check: func(t *testing.T, b *backend.Backend, errLoad error) {
				t.Helper()

				require.NoError(t, errLoad)

				_, err := b.GenesisTime()
				require.ErrorIs(t, err, backend.ErrNodeAPINotReady)

				_, err = b.GenesisForkVersion()
				require.ErrorIs(t, err, backend.ErrNodeAPINotReady)

				_, err = b.GenesisValidatorsRoot()
				require.ErrorIs(t, err, backend.ErrNodeAPINotReady)
			},
		},
		{
			name: "failed loading state",
			setMockExpectations: func(
				tcs *mocks.ConsensusService,
				_ *mocks.GenesisStateProcessor,
			) {
				t.Helper()

				tcs.EXPECT().IsAppReady().Return(testFailedLoadingState)
			},
			check: func(t *testing.T, b *backend.Backend, errLoad error) {
				t.Helper()

				require.ErrorIs(t, errLoad, testFailedLoadingState)

				_, err := b.GenesisTime()
				require.ErrorIs(t, err, testFailedLoadingState)

				_, err = b.GenesisForkVersion()
				require.ErrorIs(t, err, testFailedLoadingState)

				_, err = b.GenesisValidatorsRoot()
				require.ErrorIs(t, err, testFailedLoadingState)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cs, err := spec.MainnetChainSpec()
			require.NoError(t, err)

			cmtCfg := buildTestCometConfig(t, testGenesisTime)

			_, kvStore, depositStore, err := statetransition.BuildTestStores()
			require.NoError(t, err)
			sb := storage.NewBackend(
				cs, nil, kvStore, depositStore, nil, log.NewNopLogger(), metrics.NewNoOpTelemetrySink(),
			)

			tcs := mocks.NewConsensusService(t)
			sp := mocks.NewGenesisStateProcessor(t)

			// 2- Setup expectations before backend construction
			tc.setMockExpectations(tcs, sp)

			// 3- Build backend
			b := backend.New(sb, sp, cs, cmtCfg, tcs)
			errLoad := b.LoadData(context.Background())

			// 4- Checks
			tc.check(t, b, errLoad)
		})
	}
}

//nolint:lll // long genesis
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

	appGenesis := genutiltypes.NewAppGenesisWithVersion("test-chain", []byte(`
	{
		"beacon": {
			"fork_version": "0x04000000",
			"deposits": [],
			"execution_payload_header": {
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"feeRecipient": "0x0000000000000000000000000000000000000000",
				"stateRoot": "0xaf7ac45ece564c84ee2451776587c548aebb91ba04eb6040fd2b26055539c8e3",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"prevRandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"blockNumber": "0x0",
				"gasLimit": "0x1c9c380",
				"gasUsed": "0x0",
				"timestamp": "0x67b5f01f",
				"extraData": "0x",
				"baseFeePerGas": "1000000000",
				"blockHash": "0x0207661de38f0e54ba91c8286096e72486784c79dc6a9681fc486b38335c042f",
				"transactionsRoot": "0x7ffe241ea60187fdb0187bfa22de35d1f9bed7ab061d9401fd47e34a54fbede1",
				"withdrawalsRoot": "0x792930bbd5baac43bcc798ee49aa8185ef76bb3b44ba62b91d86ae569e4bb535",
				"blobGasUsed": "0x0",
				"excessBlobGas": "0x0"
			}
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
