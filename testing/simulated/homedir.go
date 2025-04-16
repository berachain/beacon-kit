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
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/cli/commands/initialize"
	"github.com/berachain/beacon-kit/cli/commands/server/types"
	genesisutils "github.com/berachain/beacon-kit/cli/utils/genesis"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
)

const TestnetBeaconChainID = "testnet-beacon-80069"

// InitializeHomeDir sets up a temporary home directory with the necessary genesis state
// and configuration files for testing. It returns the configured CometBFT config along with
// the computed genesis validators root.
func InitializeHomeDir(t *testing.T, chainSpec chain.Spec, tempHomeDir string, elGenesisPath string) (*cmtcfg.Config, common.Root) {
	t.Helper()

	t.Logf("Initializing home directory: %s", tempHomeDir)
	// Create the default CometBFT configuration using the temporary home directory.
	cometConfig := createCometConfig(t, tempHomeDir)

	// Run initialization command to mimic 'beacond init'
	initCommand(t, chainSpec, cometConfig.RootDir)

	// Retrieve the BLS signer from the configured home directory.
	blsSigner := GetBlsSigner(cometConfig.RootDir)

	// Set the deposit amount to the maximum effective balance.
	depositAmount := math.Gwei(chainSpec.MaxEffectiveBalance())
	// Define an arbitrary withdrawal address.
	withdrawalAddress := common.NewExecutionAddressFromHex("0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4")

	// Add a genesis deposit.
	err := genesis.AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAddress, "")
	require.NoError(t, err, "failed to add genesis deposit")

	// Collect the genesis deposit.
	err = genesis.CollectGenesisDeposits(cometConfig)
	require.NoError(t, err, "failed to collect genesis deposits")

	// Update the execution layer deposit storage with the eth-genesis file.
	err = genesis.SetDepositStorage(chainSpec, cometConfig, elGenesisPath, false)
	require.NoError(t, err, "failed to set deposit storage")

	// Add the execution payload to the genesis configuration.
	err = genesis.AddExecutionPayload(chainSpec, path.Join(cometConfig.RootDir, filepath.Base(elGenesisPath)), cometConfig)
	require.NoError(t, err, "failed to add execution payload")

	// Compute the validators root from the genesis file.
	genesisValidatorsRoot, err := genesisutils.ComputeValidatorsRootFromFile(
		path.Join(cometConfig.RootDir, "config/genesis.json"),
		chainSpec,
	)
	require.NoError(t, err, "failed to compute validators root")

	return cometConfig, genesisValidatorsRoot
}

func CopyHomeDir(t *testing.T, sourceHomeDir, targetHomeDir string) {
	t.Logf("Copying home directory to: %s", targetHomeDir)
	srcPath := filepath.Join(filepath.Clean(sourceHomeDir), ".")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cp -r %s/* %s", srcPath, targetHomeDir))
	err := cmd.Run()
	require.NoError(t, err)
}

// createCometConfig creates a default CometBFT configuration with the home directory set.
func createCometConfig(t *testing.T, tempHomeDir string) *cmtcfg.Config {
	t.Helper()
	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.RootDir = tempHomeDir
	return cmtCfg
}

// genesisCreator implements the required interface for the beacond init command, while allowing for a custom
// fork version in the genesis file.
type genesisCreator struct {
	chainSpec    chain.Spec
	cometService *cometbft.Service
}

// initCommand runs the initialization command to prepare the home directory.
// This function emulates the 'beacond init' command.
func initCommand(t *testing.T, spec chain.Spec, homeDir string) {
	t.Helper()

	// Create a Cosmos SDK client context with the provided home directory and chain ID.
	clientCtx := client.Context{}.
		WithHomeDir(homeDir).
		WithChainID(TestnetBeaconChainID)

	// Ensure necessary directories exist.
	err := os.MkdirAll(path.Join(homeDir, "config"), 0700)
	require.NoError(t, err, "failed to create config directory")

	err = os.MkdirAll(path.Join(homeDir, "data"), 0700)
	require.NoError(t, err, "failed to create data directory")

	// Initialize the command to set up the home directory.
	initCMD := initialize.InitCmd(func(_ types.AppOptions) (chain.Spec, error) { return spec, nil }, &cometbft.Service{})
	initCMD.SetArgs([]string{"test-moniker"})

	// Set the command context; required to work around a Cosmos SDK issue.
	initCMD.SetContext(context.Background())
	err = client.SetCmdClientContextHandler(clientCtx, initCMD)
	require.NoError(t, err, "failed to set client context handler")

	// Allow unknown flags to enable running tests from IDEs without extra configuration.
	initCMD.FParseErrWhitelist.UnknownFlags = true

	// Execute the initialization command.
	err = initCMD.Execute()
	require.NoError(t, err, "failed to execute init command")
}
