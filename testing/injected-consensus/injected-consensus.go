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

package injectedconsensus

import (
	"context"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/cli/commands/genesis"
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/flags"
	beaconkitconfig "github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/config/spec"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/log/phuslu"
	nodebuilder "github.com/berachain/beacon-kit/node-core/builder"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// DefaultComponents requires testing.T to avoid accidental misuse.
func DefaultComponents(_ *testing.T) []any {
	c := []any{
		components.ProvideAttributesFactory,
		components.ProvideAvailabilityStore,
		components.ProvideDepositContract,
		components.ProvideBlockStore,
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor,
		components.ProvideBlobProofVerifier,
		components.ProvideChainService,
		components.ProvideNode,
		components.ProvideChainSpec,
		components.ProvideConfig,
		components.ProvideServerConfig,
		components.ProvideDepositStore,
		components.ProvideEngineClient,
		components.ProvideExecutionEngine,
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder,
		components.ProvideReportingService,
		components.ProvideCometBFTService,
		components.ProvideServiceRegistry,
		components.ProvideSidecarFactory,
		components.ProvideStateProcessor,
		components.ProvideKVStore,
		components.ProvideStorageBackend,
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService,
		components.ProvideShutDownService,
	}
	c = append(c,
		components.ProvideNodeAPIServer,
		components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend,
	)
	//
	c = append(c, components.ProvideNodeAPIHandlers,
		components.ProvideNodeAPIBeaconHandler,
		components.ProvideNodeAPIBuilderHandler,
		components.ProvideNodeAPIConfigHandler,
		components.ProvideNodeAPIDebugHandler,
		components.ProvideNodeAPIEventsHandler,
		components.ProvideNodeAPINodeHandler,
		components.ProvideNodeAPIProofHandler,
	)

	return c
}

type TestNode struct {
	Node              nodetypes.Node
	CometService      *cometbft.Service
	BlockchainService *blockchain.Service
	CometConfig       *cmtcfg.Config
	Homedir           string
	Context           context.Context
	CancelFunc        context.CancelFunc
}

// createConfiguration creates the BeaconKit configuration and the CometBFT configuration.
func createConfiguration(t *testing.T, tempHomeDir string) (
	*beaconkitconfig.Config,
	*cmtcfg.Config,
) {
	t.Helper()
	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.RootDir = tempHomeDir
	// Forces Comet to Create it
	cmtCfg.NodeKey = "node_key.json"
	beaconCfg := beaconkitconfig.DefaultConfig()
	return beaconCfg, cmtCfg
}

// getAppOptions returns the Application Options we need to set for the Node Builder.
// Ideally we can avoid having to set the flags like this and just directly modify a config type.
func getAppOptions(t *testing.T, beaconKitConfig *beaconkitconfig.Config, tempHomeDir string) servertypes.AppOptions {
	t.Helper()
	appOpts := viper.New()
	// Execution Client Config
	appOpts.Set(flags.JWTSecretPath, "../files/jwt.hex")
	appOpts.Set(flags.RPCJWTRefreshInterval, beaconKitConfig.GetEngine().RPCJWTRefreshInterval)
	appOpts.Set(flags.RPCStartupCheckInterval, beaconKitConfig.GetEngine().RPCStartupCheckInterval)
	appOpts.Set(flags.RPCDialURL, beaconKitConfig.GetEngine().RPCDialURL)
	appOpts.Set(flags.RPCTimeout, beaconKitConfig.GetEngine().RPCTimeout)

	// BLS Config
	appOpts.Set(flags.PrivValidatorKeyFile, "./config/priv_validator_key.json")
	appOpts.Set(flags.PrivValidatorStateFile, "./data/priv_validator_state.json")

	// Beacon Config
	appOpts.Set(flags.BlockStoreServiceAvailabilityWindow, beaconKitConfig.GetBlockStoreService().AvailabilityWindow)
	appOpts.Set(flags.BlockStoreServiceEnabled, beaconKitConfig.GetBlockStoreService().Enabled)
	appOpts.Set(flags.KZGTrustedSetupPath, "../files/kzg-trusted-setup.json")
	appOpts.Set(flags.KZGImplementation, kzg.DefaultConfig().Implementation)

	// Payload Builder Config
	beaconKitConfig.GetPayloadBuilder().SuggestedFeeRecipient = common.NewExecutionAddressFromHex("0x981114102592310C347E61368342DDA67017bf84")
	appOpts.Set(flags.BuilderEnabled, beaconKitConfig.GetPayloadBuilder().Enabled)
	appOpts.Set(flags.BuildPayloadTimeout, beaconKitConfig.GetPayloadBuilder().PayloadTimeout)
	appOpts.Set(flags.SuggestedFeeRecipient, beaconKitConfig.GetPayloadBuilder().SuggestedFeeRecipient)

	// TODO: Cleanup this Set
	appOpts.Set("pruning", "default")
	appOpts.Set("home", tempHomeDir)
	return appOpts
}

// NewTestNode Uses the mainnet chainspec.
func NewTestNode(t *testing.T) *TestNode {
	t.Helper()

	ctx, cancelFunc := context.WithCancel(context.Background())
	logger := phuslu.NewLogger(os.Stdout, nil)

	tempHomeDir := t.TempDir()
	beaconKitConfig, cometConfig := createConfiguration(t, tempHomeDir)

	chainSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	appOpts := getAppOptions(t, beaconKitConfig, tempHomeDir)

	// Chain Spec
	t.Setenv(components.ChainSpecTypeEnvVar, components.MainnetChainSpecType)

	// Create the genesis deposit
	blsSigner := signer.BLSSigner{PrivValidator: types.NewMockPVWithKeyType(bls12381.KeyType)}

	// Make the deposit amount the Max effective balance - set arbitrarily higher than 250K BERA required for mainnet
	depositAmount := math.Gwei(chainSpec.MaxEffectiveBalance())
	withdrawalAddress := common.NewExecutionAddressFromHex("0x6Eb9C23e4c187452504Ef8c5fD8fA1a4b15BE162")
	err = genesis.AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAddress, "")
	require.NoError(t, err)

	// Collect the genesis deposit
	err = genesis.CollectGenesisValidators(cometConfig)
	require.NoError(t, err)

	// Update the EL Deposit Storage
	err = genesis.SetDepositStorage(chainSpec, cometConfig, "TBD", false)
	require.NoError(t, err)

	// 1. Build a node builder with your default or custom test components.
	nb := nodebuilder.New(
		nodebuilder.WithComponents[nodetypes.Node](DefaultComponents(t)),
	)

	database, err := db.OpenDB(tempHomeDir, dbm.PebbleDBBackend)
	require.NoError(t, err)

	node := nb.Build(
		logger,
		database,
		os.Stdout, // or some other writer
		cometConfig,
		appOpts,
	)

	// Fetch services we will want to query and interact with so they are easily accessible in testing
	var cometService *cometbft.Service
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
		CometConfig:       cometConfig,
		Homedir:           tempHomeDir,
		Context:           ctx,
		CancelFunc:        cancelFunc,
	}
}

// func genesisFromFile(t *testing.T, file string) *cosmosutil.AppGenesis {
//	t.Helper()
//	appGenesis, err := cosmosutil.AppGenesisFromFile(file)
//	require.NoError(t, err)
//	return appGenesis
//}
