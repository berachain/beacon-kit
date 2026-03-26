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
	"io"
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
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
)

const TestnetBeaconChainID = "testnet-beacon-80069"

const WithdrawalExecutionAddress = "0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"

// InitializeHomeDirs sets up one or more temporary home directories with the
// necessary genesis state and configuration files for testing. It returns a
// CometBFT config per home directory along with the computed genesis validators
// root. When multiple directories are provided, premined deposits are
// consolidated into the first directory and the resulting genesis file is
// replicated to all others.
func InitializeHomeDirs(
	t *testing.T,
	chainSpec chain.Spec,
	elGenesisPath string,
	tempHomeDirs ...string,
) ([]*cmtcfg.Config, common.Root) {
	t.Helper()
	require.NotEmpty(t, tempHomeDirs, "at least one home directory is required")

	n := len(tempHomeDirs)
	t.Logf("Initializing %d home director(ies): %v", n, tempHomeDirs)

	configs := make([]*cmtcfg.Config, n)
	for i, dir := range tempHomeDirs {
		configs[i] = createCometConfig(t, dir)
		initCommand(t, chainSpec, configs[i].RootDir)
	}

	depositAmount := chainSpec.MaxEffectiveBalance()
	withdrawalAddress, err := common.NewExecutionAddressFromHex(WithdrawalExecutionAddress)
	require.NoError(t, err, "failed to create withdrawal address")

	for i, cfg := range configs {
		signer := GetBlsSigner(cfg.RootDir)
		err = genesis.AddGenesisDeposit(chainSpec, cfg, signer, depositAmount, withdrawalAddress, "")
		require.NoError(t, err, "failed to add genesis deposit %d", i+1)
	}

	primaryDir := filepath.Join(configs[0].RootDir, "config", "premined-deposits")
	for i := 1; i < n; i++ {
		copyMissingFiles(t, filepath.Join(configs[i].RootDir, "config", "premined-deposits"), primaryDir)
	}

	err = genesis.CollectGenesisDeposits(configs[0])
	require.NoError(t, err, "failed to collect genesis deposits")

	if n > 1 {
		srcGenesis := filepath.Join(configs[0].RootDir, "config", "genesis.json")
		data, readErr := os.ReadFile(srcGenesis)
		require.NoError(t, readErr, "failed to read genesis file: %s", srcGenesis)
		for i := 1; i < n; i++ {
			dst := filepath.Join(configs[i].RootDir, "config", "genesis.json")
			require.NoError(t, os.WriteFile(dst, data, 0o600), "failed to write genesis file: %s", dst)
		}
	}

	for i, cfg := range configs {
		err = genesis.SetDepositStorage(chainSpec, cfg, elGenesisPath)
		require.NoError(t, err, "failed to set deposit storage %d", i+1)

		err = genesis.AddExecutionPayload(chainSpec, path.Join(cfg.RootDir, filepath.Base(elGenesisPath)), cfg)
		require.NoError(t, err, "failed to add execution payload %d", i+1)
	}

	genesisValidatorsRoot, err := genesisutils.ComputeValidatorsRootFromFile(
		path.Join(configs[0].RootDir, "config/genesis.json"),
		chainSpec,
	)
	require.NoError(t, err, "failed to compute validators root")
	for i := 1; i < n; i++ {
		root, rootErr := genesisutils.ComputeValidatorsRootFromFile(
			path.Join(configs[i].RootDir, "config/genesis.json"),
			chainSpec,
		)
		require.NoError(t, rootErr, "failed to compute validators root %d", i+1)
		require.Equal(t, genesisValidatorsRoot, root, "validators' roots should be equal")
	}

	return configs, genesisValidatorsRoot
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
