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

package simulated

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/cli/commands/initialize"
	"github.com/berachain/beacon-kit/cli/flags"
	beaconkitconfig "github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/config/spec"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/net/url"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ory/dockertest"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// ExecutionClient provides the interface for managing an ExecutionClient
type ExecutionClient interface {
	Start(t *testing.T) (*dockertest.Resource, *url.ConnectionURL)
}

type TestSuiteHandle struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	HomeDir    string
	TestNode   TestNode

	// Geth dockertest handles for closing
	ElHandle *dockertest.Resource
}

func InitializeHomeDir(t *testing.T, tempHomeDir string) *cmtcfg.Config {
	t.Helper()
	t.Logf("tempHomeDir=%s", tempHomeDir)
	cometConfig := createCometConfig(t, tempHomeDir)

	chainSpec, err := spec.TestnetChainSpec()
	require.NoError(t, err)

	t.Setenv(components.ChainSpecTypeEnvVar, components.TestnetChainSpecType)

	// Same as `beacond init`
	initCommand(t, tempHomeDir)

	// get the bls signer from the homedir
	blsSigner := getBlsSigner(tempHomeDir)

	// Make the deposit amount the Max effective balance - set arbitrarily higher than 250K BERA required for mainnet
	depositAmount := math.Gwei(chainSpec.MaxEffectiveBalance())
	// Arbitrary withdrawal address
	withdrawalAddress := common.NewExecutionAddressFromHex("0x6Eb9C23e4c187452504Ef8c5fD8fA1a4b15BE162")

	err = genesis.AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAddress, "")
	require.NoError(t, err)
	// Collect the genesis deposit
	err = genesis.CollectGenesisDeposits(cometConfig)
	require.NoError(t, err)
	// Update the EL Deposit Storage
	err = genesis.SetDepositStorage(chainSpec, cometConfig, "./eth-genesis.json", false)
	require.NoError(t, err)
	err = genesis.AddExecutionPayload(chainSpec, path.Join(tempHomeDir, "eth-genesis.json"), cometConfig)
	require.NoError(t, err)
	return cometConfig
}

// createConfiguration creates the BeaconKit configuration and the CometBFT configuration.
func createCometConfig(t *testing.T, tempHomeDir string) *cmtcfg.Config {
	t.Helper()
	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.RootDir = tempHomeDir
	return cmtCfg
}

func createBeaconKitConfig(_ *testing.T) *beaconkitconfig.Config {
	return beaconkitconfig.DefaultConfig()
}

// getAppOptions returns the Application Options we need to set for the Node Builder.
// Ideally we can avoid having to set the flags like this and just directly modify a config type.
func getAppOptions(t *testing.T, appOpts *viper.Viper, beaconKitConfig *beaconkitconfig.Config, tempHomeDir string) *viper.Viper {
	t.Helper()
	// Execution Client Config
	relativePathJwt := "../files/jwt.hex"
	jwtPath, err := filepath.Abs(relativePathJwt)
	require.NoError(t, err)
	appOpts.Set(flags.JWTSecretPath, jwtPath)
	appOpts.Set(flags.RPCJWTRefreshInterval, beaconKitConfig.GetEngine().RPCJWTRefreshInterval.String())
	appOpts.Set(flags.RPCStartupCheckInterval, beaconKitConfig.GetEngine().RPCStartupCheckInterval.String())
	appOpts.Set(flags.RPCDialURL, beaconKitConfig.GetEngine().RPCDialURL.String())
	appOpts.Set(flags.RPCTimeout, beaconKitConfig.GetEngine().RPCTimeout.String())

	appOpts.Set(flags.LogLevel, "debug")

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

func getBlsSigner(tempHomeDir string) *signer.BLSSigner {
	privValKeyFile := filepath.Join(tempHomeDir, "config/priv_validator_key.json")
	privValStateFile := filepath.Join(tempHomeDir, "data/priv_validator_state.json")
	return signer.NewBLSSigner(privValKeyFile, privValStateFile)
}

func initCommand(t *testing.T, tempHomeDir string) {
	t.Helper()

	clientCtx := client.Context{}.
		WithHomeDir(tempHomeDir).
		WithChainID("test-mainnet-chain")

	// This is necessary otherwise cosmos-sdk will see errors
	err := os.MkdirAll(tempHomeDir+"/config", 0700)
	require.NoError(t, err)

	err = os.MkdirAll(tempHomeDir+"/data", 0700)
	require.NoError(t, err)

	initCMD := initialize.InitCmd(&cometbft.Service{})
	// This is required due to a bug in cosmos sdk
	initCMD.SetContext(context.Background())

	err = client.SetCmdClientContextHandler(clientCtx, initCMD)
	require.NoError(t, err)

	// This is so that Goland can run the test from the IDE through test filtering
	initCMD.FParseErrWhitelist.UnknownFlags = true

	err = initCMD.Execute()
	require.NoError(t, err)
}
