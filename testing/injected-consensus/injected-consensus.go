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

package injected_consensus

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/cli/flags"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/config/spec"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/da/kzg"
	executionconfig "github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	nodebuilder "github.com/berachain/beacon-kit/node-core/builder"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	cosmosutil "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/params"
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
	Node              nodetypes.Node
	CometService      *cometbft.Service[*phuslu.Logger]
	BlockchainService *blockchain.Service
	CometConfig       *cmtcfg.Config
	Homedir           string
	Context           context.Context
	CancelFunc        context.CancelFunc
}

// func copyFile(t *testing.T, src, dst string) error {
//	t.Helper()
//	// Open the source file
//	srcFile, err := os.Open(src)
//	if err != nil {
//		return fmt.Errorf("failed to open source file %s: %w", src, err)
//	}
//	defer srcFile.Close()
//
//	// Create the destination file
//	dstFile, err := os.Create(dst)
//	if err != nil {
//		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
//	}
//	defer dstFile.Close()
//
//	// Copy the file contents
//	_, err = io.Copy(dstFile, srcFile)
//	if err != nil {
//		return fmt.Errorf("failed to copy from %s to %s: %w", src, dst, err)
//	}
//
//	return nil
//}

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run command %s %v: %w", command, args, err)
	}
	return nil
}

func initializeBeaconState(t *testing.T, homedir string) {
	ethGenesis := "./eth-genesis.json"

	commands := [][]string{
		{"./build/bin/beacond", "genesis", "add-premined-deposit", "--home", homedir, "32000000000", "0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"},
		{"./build/bin/beacond", "genesis", "collect-premined-deposits", "--home", homedir},
		{"./build/bin/beacond", "genesis", "set-deposit-storage", ethGenesis, "--home", homedir},
		{"./build/bin/beacond", "genesis", "execution-payload", ethGenesis, "--home", homedir},
	}

	for _, cmdArgs := range commands {
		err := runCommand(cmdArgs[0], cmdArgs[1:]...)
		require.NoError(t, err)
	}
}

// Uses the mainnet chainspec.
func NewTestNode(t *testing.T) *TestNode {
	t.Helper()

	// Create a test node that
	ctx, cancelFunc := context.WithCancel(context.Background())
	// 1. Build a node builder with your default or custom test components.
	nb := nodebuilder.New(
		nodebuilder.WithComponents[
			nodetypes.Node,
			*phuslu.Logger,
			*phuslu.Config,
		](DefaultComponents(t)),
	)

	tempHomeDir := t.TempDir()

	// initializeBeaconState(t, tempHomeDir)

	logger := phuslu.NewLogger(os.Stdout, nil)

	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.RootDir = tempHomeDir

	beaconCfg := config.DefaultConfig()
	executionClientConfig := executionconfig.DefaultConfig()
	payloadBuilderCfg := payloadbuilder.DefaultConfig()

	chainSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	// Ideally we can avoid having to set the flags like this and just directly modify a config type
	appOpts := viper.New()
	// Execution Client Config
	appOpts.Set(flags.JWTSecretPath, "../files/jwt.hex")
	appOpts.Set(flags.RPCJWTRefreshInterval, executionClientConfig.RPCJWTRefreshInterval)
	appOpts.Set(flags.RPCStartupCheckInterval, executionClientConfig.RPCStartupCheckInterval)
	appOpts.Set(flags.RPCDialURL, executionClientConfig.RPCDialURL)
	appOpts.Set(flags.RPCTimeout, executionClientConfig.RPCTimeout)

	// BLS Config
	appOpts.Set(flags.PrivValidatorKeyFile, "./config/priv_validator_key.json")
	appOpts.Set(flags.PrivValidatorStateFile, "./data/priv_validator_state.json")

	// Beacon Config
	appOpts.Set(flags.BlockStoreServiceAvailabilityWindow, beaconCfg.BlockStoreService.AvailabilityWindow)
	appOpts.Set(flags.BlockStoreServiceEnabled, beaconCfg.BlockStoreService.Enabled)
	appOpts.Set(flags.KZGTrustedSetupPath, "../files/kzg-trusted-setup.json")
	appOpts.Set(flags.KZGImplementation, kzg.DefaultConfig().Implementation)

	// Payload Builder Config
	payloadBuilderCfg.SuggestedFeeRecipient = common.NewExecutionAddressFromHex("0x981114102592310C347E61368342DDA67017bf84")
	appOpts.Set(flags.BuilderEnabled, payloadBuilderCfg.Enabled)
	appOpts.Set(flags.BuildPayloadTimeout, payloadBuilderCfg.PayloadTimeout)
	appOpts.Set(flags.SuggestedFeeRecipient, payloadBuilderCfg.SuggestedFeeRecipient)

	// Chain Spec
	t.Setenv(components.ChainSpecTypeEnvVar, components.MainnetChainSpecType)

	// Create the genesis deposit
	blsSigner := signer.BLSSigner{PrivValidator: types.NewMockPVWithKeyType(bls12381.KeyType)}
	depositAmount := math.Gwei(250_000 * params.GWei)
	withdrawalAddress := common.NewExecutionAddressFromHex("0x6Eb9C23e4c187452504Ef8c5fD8fA1a4b15BE162")
	err = genesis.AddGenesisDeposit(chainSpec, cmtCfg, blsSigner, depositAmount, withdrawalAddress, "")
	require.NoError(t, err)

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
		Node:              node,
		CometService:      cometService,
		BlockchainService: blockchainService,
		CometConfig:       cmtCfg,
		Homedir:           tempHomeDir,
		Context:           ctx,
		CancelFunc:        cancelFunc,
	}
}

func genesisFromFile(t *testing.T, file string) *cosmosutil.AppGenesis {
	t.Helper()
	appGenesis, err := cosmosutil.AppGenesisFromFile(file)
	require.NoError(t, err)
	return appGenesis
}
