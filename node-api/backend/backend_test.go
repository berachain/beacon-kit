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

package backend_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/config/spec"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/backend/mocks"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	coremocks "github.com/berachain/beacon-kit/node-core/types/mocks"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	cmtcfg "github.com/cometbft/cometbft/config"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBackendLoadData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                string
		setMockExpectations func(
			*coremocks.ConsensusService,
			*mocks.GenesisStateProcessor,
		)
		check func(t *testing.T, b *backend.Backend, errLoad error)
	}{
		{
			name: "success",
			setMockExpectations: func(
				cs *coremocks.ConsensusService,
				sp *mocks.GenesisStateProcessor,
			) {
				t.Helper()

				cs.EXPECT().IsAppReady().Return(nil) // mark the app as ready
				sp.EXPECT().InitializeBeaconStateFromEth1(
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil, nil) // duly process genesis once parsed
			},
			check: func(t *testing.T, _ *backend.Backend, errLoad error) {
				t.Helper()

				require.NoError(t, errLoad)
			},
		},
		{
			name: "app not ready",
			setMockExpectations: func(
				cs *coremocks.ConsensusService,
				_ *mocks.GenesisStateProcessor,
			) {
				t.Helper()

				cs.EXPECT().IsAppReady().Return(cometbft.ErrAppNotReady) // mark the app as not ready
			},
			check: func(t *testing.T, _ *backend.Backend, errLoad error) {
				t.Helper()

				// just keep going, it will load later on, as soon as possible
				require.NoError(t, errLoad)
			},
		},
		{
			name: "could not check app is ready",
			setMockExpectations: func(
				cs *coremocks.ConsensusService,
				_ *mocks.GenesisStateProcessor,
			) {
				t.Helper()

				cs.EXPECT().IsAppReady().Return(errors.New("unknown error"))
			},
			check: func(t *testing.T, _ *backend.Backend, errLoad error) {
				t.Helper()
				require.Error(t, errLoad)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1 - Build backend
			cs, err := spec.MainnetChainSpec()
			require.NoError(t, err)

			cmtCfg := buildTestCometConfig(t)

			_, kvStore, depositStore, err := statetransition.BuildTestStores()
			require.NoError(t, err)
			sb := storage.NewBackend(
				cs, nil, kvStore, depositStore, nil, log.NewNopLogger(), metrics.NewNoOpTelemetrySink(),
			)

			tcs := coremocks.NewConsensusService(t)
			sp := mocks.NewGenesisStateProcessor(t)

			b := backend.New(sb, sp, cs, cmtCfg, tcs)
			defer func() {
				require.NoError(t, b.Close())
			}()

			// 2- Setup expectations
			tc.setMockExpectations(tcs, sp)

			// 3 - Test
			errLoad := b.LoadData(t.Context())

			// 4- Checks
			tc.check(t, b, errLoad)
		})
	}
}

//nolint:lll // adapted genesis from mainnet
func buildTestCometConfig(t *testing.T) *cmtcfg.Config {
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
      "fork_version": "0x04000000",
      "deposits": [],
      "execution_payload_header": {
        "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "feeRecipient": "0x0000000000000000000000000000000000000000",
        "stateRoot": "0x2aace2f233f1ef6ca13e5fd8feae4cb1b0b580fa56c8ee081ab89d861eaf1515",
        "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "prevRandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "blockNumber": "0x0",
        "gasLimit": "0x1c9c380",
        "gasUsed": "0x0",
        "timestamp": "0x678e56e0",
        "extraData": "0x",
        "baseFeePerGas": "1000000000",
        "blockHash": "0xd57819422128da1c44339fc7956662378c17e2213e669b427ac91cd11dfcfb38",
        "transactionsRoot": "0x7ffe241ea60187fdb0187bfa22de35d1f9bed7ab061d9401fd47e34a54fbede1",
        "withdrawalsRoot": "0x792930bbd5baac43bcc798ee49aa8185ef76bb3b44ba62b91d86ae569e4bb535",
        "blobGasUsed": "0x0",
        "excessBlobGas": "0x0"
      }
    }
	}
	`))

	// Save genesis file to the config directory
	genesisFile := filepath.Join(configDir, "genesis.json")
	err = appGenesis.SaveAs(genesisFile)
	require.NoError(t, err)

	return cmtCfg
}
