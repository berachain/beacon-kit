//go:build simulated

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
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/cli/commands/initialize"
	genesisutils "github.com/berachain/beacon-kit/cli/utils/genesis"
	"github.com/berachain/beacon-kit/config/spec"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
)

func InitializeHomeDir(t *testing.T, tempHomeDir string) (*cmtcfg.Config, common.Root) {
	t.Helper()
	t.Logf("tempHomeDir=%s", tempHomeDir)
	cometConfig := createCometConfig(t, tempHomeDir)

	chainSpec, err := spec.TestnetChainSpec()
	require.NoError(t, err)

	t.Setenv(components.ChainSpecTypeEnvVar, components.TestnetChainSpecType)

	// Same as `beacond init`
	initCommand(t, cometConfig.RootDir)

	// get the bls signer from the homedir
	blsSigner := GetBlsSigner(cometConfig.RootDir)

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
	err = genesis.AddExecutionPayload(chainSpec, path.Join(cometConfig.RootDir, "eth-genesis.json"), cometConfig)
	require.NoError(t, err)

	genesisValidatorsRoot, err := genesisutils.ComputeValidatorsRootFromFile(path.Join(cometConfig.RootDir, "config/genesis.json"), chainSpec)
	require.NoError(t, err)
	return cometConfig, genesisValidatorsRoot
}

// createConfiguration creates the BeaconKit configuration and the CometBFT configuration.
func createCometConfig(t *testing.T, tempHomeDir string) *cmtcfg.Config {
	t.Helper()
	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.RootDir = tempHomeDir
	return cmtCfg
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
	initCMD.SetArgs([]string{
		"test-moniker",
	})

	// This is required due to a bug in cosmos sdk
	initCMD.SetContext(context.Background())
	err = client.SetCmdClientContextHandler(clientCtx, initCMD)
	require.NoError(t, err)

	// This is so that Goland can run the test from the IDE through test filtering
	initCMD.FParseErrWhitelist.UnknownFlags = true
	err = initCMD.Execute()
	require.NoError(t, err)
}
