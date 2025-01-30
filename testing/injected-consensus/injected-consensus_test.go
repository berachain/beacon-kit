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

package injected_consensus_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/cli/flags"
	"github.com/berachain/beacon-kit/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/da/kzg"
	executionconfig "github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	nodebuilder "github.com/berachain/beacon-kit/node-core/builder"
	"github.com/berachain/beacon-kit/node-core/components"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	cosmosutil "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// DefaultComponents requires testing.T to avoid accidental misuse.
func DefaultComponents(_ *testing.T) []any {
	c := []any{
		components.ProvideAttributesFactory[*phuslu.Logger],
		components.ProvideAvailabilityStore[*phuslu.Logger],
		components.ProvideDepositContract,
		components.ProvideBlockStore[*phuslu.Logger],
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor[*phuslu.Logger],
		components.ProvideBlobProofVerifier,
		components.ProvideChainService[*phuslu.Logger],
		components.ProvideNode,
		components.ProvideChainSpec,
		components.ProvideConfig,
		components.ProvideServerConfig,
		components.ProvideDepositStore[*phuslu.Logger],
		components.ProvideEngineClient[*phuslu.Logger],
		components.ProvideExecutionEngine[*phuslu.Logger],
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder[*phuslu.Logger],
		components.ProvideReportingService[*phuslu.Logger],
		components.ProvideCometBFTService[*phuslu.Logger],
		components.ProvideServiceRegistry[*phuslu.Logger],
		components.ProvideSidecarFactory,
		components.ProvideStateProcessor[*phuslu.Logger],
		components.ProvideKVStore,
		components.ProvideStorageBackend,
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService[*phuslu.Logger],
		components.ProvideShutDownService[*phuslu.Logger],
	}
	c = append(c,
		components.ProvideNodeAPIServer[*phuslu.Logger, echo.Context],
		components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend,
	)
	//
	c = append(c, components.ProvideNodeAPIHandlers[echo.Context],
		components.ProvideNodeAPIBeaconHandler[echo.Context],
		components.ProvideNodeAPIBuilderHandler[echo.Context],
		components.ProvideNodeAPIConfigHandler[echo.Context],
		components.ProvideNodeAPIDebugHandler[echo.Context],
		components.ProvideNodeAPIEventsHandler[echo.Context],
		components.ProvideNodeAPINodeHandler[echo.Context],
		components.ProvideNodeAPIProofHandler[echo.Context],
	)

	return c
}

type TestNode struct {
	node              nodetypes.Node
	cometService      *cometbft.Service[*phuslu.Logger]
	blockchainService *blockchain.Service
	cometConfig       *cmtcfg.Config
	homedir           string
}

func makeTempHomeDir(t *testing.T) string {
	t.Helper()
	// create random suffix to avoid conflicts
	const rndSuffixLen = 5
	bytes := make([]byte, rndSuffixLen)
	_, err := rand.Read(bytes)
	require.NoError(t, err)

	rndSuffix := hex.EncodeToString(bytes)

	tmpFolder := filepath.Join(os.TempDir(), "/injected-consensus", rndSuffix)
	require.NoError(t, os.MkdirAll(tmpFolder, os.ModePerm))
	t.Log("tmp folder:", tmpFolder)
	return tmpFolder
}

func copyFile(t *testing.T, src, dst string) error {
	t.Helper()
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	// Copy the file contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy from %s to %s: %w", src, dst, err)
	}

	return nil
}

// and does not start the CometBFT loop, allowing us to inject our own calls.
func newTestNode(t *testing.T) *TestNode {
	t.Helper()
	// 1. Build a node builder with your default or custom test components.
	nb := nodebuilder.New(
		nodebuilder.WithComponents[
			nodetypes.Node,
			*phuslu.Logger,
			*phuslu.Config,
		](DefaultComponents(t)),
	)

	// tempHomeDir := makeTempHomeDir(t)
	tempHomeDir := "/Users/rezbera/Code/beacon-kit/.tmp/beacond"

	logger := phuslu.NewLogger(os.Stdout, nil)

	// Use an in-memory DB
	// db := dbm.NewMemDB()
	cmtCfg := cometbft.DefaultConfig()
	beaconCfg := config.DefaultConfig()
	executionClientConfig := executionconfig.DefaultConfig()

	appOpts := viper.New()

	// err := copyFile(t, "./test_priv_validator_key.json", tempHomeDir+"/priv_validator_key.json")
	// require.NoError(t, err)
	// err = copyFile(t, "./test_priv_validator_state.json", tempHomeDir+"/priv_validator_state.json")
	// require.NoError(t, err)

	appOpts.Set(flags.JWTSecretPath, "../files/jwt.hex")
	appOpts.Set(flags.RPCJWTRefreshInterval, executionClientConfig.RPCJWTRefreshInterval)
	appOpts.Set(flags.RPCStartupCheckInterval, executionClientConfig.RPCStartupCheckInterval)
	appOpts.Set(flags.RPCDialURL, executionClientConfig.RPCDialURL)
	// appOpts.Set(flags.PrivValidatorKeyFile, "./config/priv_validator_key.json")
	// appOpts.Set(flags.PrivValidatorStateFile, "./data/priv_validator_state.json")

	appOpts.Set(flags.BlockStoreServiceAvailabilityWindow, beaconCfg.BlockStoreService.AvailabilityWindow)
	appOpts.Set(flags.BlockStoreServiceEnabled, beaconCfg.BlockStoreService.Enabled)
	appOpts.Set(flags.KZGTrustedSetupPath, "../files/kzg-trusted-setup.json")
	appOpts.Set(flags.KZGImplementation, kzg.DefaultConfig().Implementation)

	t.Setenv(components.ChainSpecTypeEnvVar, components.DevnetChainSpecType)

	// TODO: Cleanup this Set
	appOpts.Set("pruning", "default")
	appOpts.Set("home", tempHomeDir)

	database, err := db.OpenDB(tempHomeDir, dbm.PebbleDBBackend)
	require.NoError(t, err)

	node := nb.Build(
		logger,
		database,
		os.Stdout, // or some other writer
		cmtCfg,
		appOpts,
	)

	var cometService *cometbft.Service[*phuslu.Logger]
	err = node.FetchService(&cometService)
	require.NoError(t, err)
	require.NotNil(t, cometService)

	var blockchainService *blockchain.Service
	err = node.FetchService(&blockchainService)
	require.NoError(t, err)
	require.NotNil(t, blockchainService)

	return &TestNode{
		node:              node,
		cometService:      cometService,
		blockchainService: blockchainService,
		cometConfig:       cmtCfg,
		homedir:           tempHomeDir,
	}
}

func genesisFromFile(t *testing.T, file string) *cosmosutil.AppGenesis {
	t.Helper()
	appGenesis, err := cosmosutil.AppGenesisFromFile(file)
	require.NoError(t, err)
	return appGenesis
}
