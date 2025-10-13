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
	"io"
	"os"
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
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
)

const TestnetBeaconChainID = "testnet-beacon-80069"

const WithdrawalExecutionAddress = "0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"

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
	depositAmount := chainSpec.MaxEffectiveBalance()
	// Define an arbitrary withdrawal address.
	withdrawalAddress, err := common.NewExecutionAddressFromHex(WithdrawalExecutionAddress)
	require.NoError(t, err, "failed to create withdrawal address")

	// Add a genesis deposit.
	err = genesis.AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAddress, "")
	require.NoError(t, err, "failed to add genesis deposit")

	// Collect the genesis deposit.
	err = genesis.CollectGenesisDeposits(cometConfig)
	require.NoError(t, err, "failed to collect genesis deposits")

	// Update the execution layer deposit storage with the eth-genesis file.
	err = genesis.SetDepositStorage(chainSpec, cometConfig, elGenesisPath)
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

// Initialize2HomeDirs sets up a temporary home directory with the necessary genesis state
// and configuration files for testing. It returns the configured CometBFT config along with
// the computed genesis validators root. This creates a 2 validator setup.
func Initialize2HomeDirs(
	t *testing.T,
	chainSpec chain.Spec,
	tempHomeDir1, tempHomeDir2, elGenesisPath string,
) (*cmtcfg.Config, *cmtcfg.Config, common.Root) {
	t.Helper()

	t.Logf("Initializing home directory: %s and %s", tempHomeDir1, tempHomeDir2)
	// Create the default CometBFT configuration using the temporary home directory.
	cmtCfg1 := createCometConfig(t, tempHomeDir1)
	// Create a new temp home dir for the second validator.
	cmtCfg2 := createCometConfig(t, tempHomeDir2)

	// Run initialization command to mimic 'beacond init'
	initCommand(t, chainSpec, cmtCfg1.RootDir)
	initCommand(t, chainSpec, cmtCfg2.RootDir)

	// Retrieve the BLS signer from the configured home directory.
	blsSigner1 := GetBlsSigner(cmtCfg1.RootDir)
	blsSigner2 := GetBlsSigner(cmtCfg2.RootDir)

	// Set the deposit amount to the maximum effective balance.
	depositAmount := chainSpec.MaxEffectiveBalance()
	// Define an arbitrary withdrawal address.
	withdrawalAddress, err := common.NewExecutionAddressFromHex(WithdrawalExecutionAddress)
	require.NoError(t, err, "failed to create withdrawal address")

	// Add a genesis deposit for the first validator.
	err = genesis.AddGenesisDeposit(chainSpec, cmtCfg1, blsSigner1, depositAmount, withdrawalAddress, "")
	require.NoError(t, err, "failed to add genesis deposit")
	// Add a genesis deposit for the second validator.
	err = genesis.AddGenesisDeposit(chainSpec, cmtCfg2, blsSigner2, depositAmount, withdrawalAddress, "")
	require.NoError(t, err, "failed to add genesis deposit 2")

	// cmtCfg1 contains premined deposit for both validators.
	dir1 := filepath.Join(cmtCfg1.RootDir, "config", "premined-deposits")
	dir2 := filepath.Join(cmtCfg2.RootDir, "config", "premined-deposits")
	copyMissingFiles(t, dir2, dir1)

	// Collect the genesis deposits in cmtCfg1.
	err = genesis.CollectGenesisDeposits(cmtCfg1)
	require.NoError(t, err, "failed to collect genesis deposits")

	// Copy the genesis file from cmtCfg1 to cmtCfg2 so both validators have the same genesis file.
	srcGenesis := filepath.Join(cmtCfg1.RootDir, "config", "genesis.json")
	dstGenesis := filepath.Join(cmtCfg2.RootDir, "config", "genesis.json")
	data, err := os.ReadFile(srcGenesis)
	require.NoError(t, err, "failed to read genesis file: %s", srcGenesis)
	err = os.WriteFile(dstGenesis, data, 0o600)
	require.NoError(t, err, "failed to write genesis file: %s", dstGenesis)

	// Update the execution layer deposit storage with the eth-genesis file.
	err = genesis.SetDepositStorage(chainSpec, cmtCfg1, elGenesisPath)
	require.NoError(t, err, "failed to set deposit storage")
	err = genesis.SetDepositStorage(chainSpec, cmtCfg2, elGenesisPath)
	require.NoError(t, err, "failed to set deposit storage 2")

	// Add the execution payload to the genesis configuration.
	err = genesis.AddExecutionPayload(chainSpec, path.Join(cmtCfg1.RootDir, filepath.Base(elGenesisPath)), cmtCfg1)
	require.NoError(t, err, "failed to add execution payload")
	err = genesis.AddExecutionPayload(chainSpec, path.Join(cmtCfg2.RootDir, filepath.Base(elGenesisPath)), cmtCfg2)
	require.NoError(t, err, "failed to add execution payload 2")

	// Compute the validators root from the genesis file.
	genesisValidatorsRoot, err := genesisutils.ComputeValidatorsRootFromFile(
		path.Join(cmtCfg1.RootDir, "config/genesis.json"),
		chainSpec,
	)
	require.NoError(t, err, "failed to compute validators root")
	genesisValidatorsRoot2, err := genesisutils.ComputeValidatorsRootFromFile(
		path.Join(cmtCfg2.RootDir, "config/genesis.json"),
		chainSpec,
	)
	require.NoError(t, err, "failed to compute validators root 2")
	require.Equal(t, genesisValidatorsRoot, genesisValidatorsRoot2, "validators' roots should be equal")

	return cmtCfg1, cmtCfg2, genesisValidatorsRoot
}

// createCometConfig creates a default CometBFT configuration with the home directory set.
func createCometConfig(t *testing.T, tempHomeDir string) *cmtcfg.Config {
	t.Helper()
	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.RootDir = tempHomeDir
	return cmtCfg
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

// copyMissingFiles copies files that exist in srcDir but not in dstDir.
func copyMissingFiles(t *testing.T, srcDir, dstDir string) {
	t.Helper()

	// Ensure destination directory exists.
	require.NoError(t, os.MkdirAll(dstDir, 0o700))

	srcEntries, err := os.ReadDir(srcDir)
	require.NoError(t, err, "failed to read source directory: %s", srcDir)

	for _, e := range srcEntries {
		if e.IsDir() {
			continue
		}
		srcPath := filepath.Join(srcDir, e.Name())
		dstPath := filepath.Join(dstDir, e.Name())

		if _, err := os.Stat(dstPath); err == nil {
			// Already exists.
			continue
		}

		srcFile, err := os.Open(srcPath)
		require.NoError(t, err, "failed to open src file: %s", srcPath)
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		require.NoError(t, err, "failed to create dst file: %s", dstPath)
		_, err = io.Copy(dstFile, srcFile)
		require.NoError(t, err, "failed to copy to dst file: %s", dstPath)
		require.NoError(t, dstFile.Close())
	}
}
