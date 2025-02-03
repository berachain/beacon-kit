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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package genesis_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestSetDepositStorageCmd(t *testing.T) {
	t.Parallel()
	t.Run("command should be available and have correct use", func(t *testing.T) {
		t.Parallel()
		chainSpec, err := spec.DevnetChainSpec()
		require.NoError(t, err)
		cmd := genesis.SetDepositStorageCmd(chainSpec)
		require.Equal(t, "set-deposit-storage [eth/genesis/file.json]", cmd.Use)
	})

	t.Run("should set deposit storage correctly", func(t *testing.T) {
		t.Parallel()
		// Create a temporary directory for test files
		tmpDir := t.TempDir()

		// Setup test files
		mockGenesisPath := setupMockGenesis(t, tmpDir)
		_ = setupMockCLGenesis(t, tmpDir)

		// Setup client context
		clientCtx := client.Context{
			HomeDir: tmpDir,
		}
		ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

		// Create and execute the command
		chainSpec, err := spec.DevnetChainSpec()
		require.NoError(t, err)
		cmd := genesis.SetDepositStorageCmd(chainSpec)
		cmd.SetContext(ctx)
		// Change working directory to tmpDir for the test
		currentDir, err := os.Getwd()
		require.NoError(t, err)
		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		defer func() {
			err = os.Chdir(currentDir)
			require.NoError(t, err)
		}()

		cmd.SetArgs([]string{mockGenesisPath})
		require.NoError(t, cmd.Execute())

		verifyStorageOutput(t, mockGenesisPath)
	})
}

func setupMockGenesis(t *testing.T, tmpDir string) string {
	t.Helper()
	chainSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	depositAddr := common.Address(chainSpec.DepositContractAddress())

	mockGenesisPath := filepath.Join(tmpDir, "genesis.json")
	mockGenesis := map[string]interface{}{
		"alloc": map[string]interface{}{
			depositAddr.Hex(): map[string]interface{}{
				"balance": "0x0",
				"code":    "0x1234",
			},
		},
	}
	genesisBz, err := json.MarshalIndent(mockGenesis, "", "  ")
	require.NoError(t, err)
	err = afero.WriteFile(afero.NewOsFs(), mockGenesisPath, genesisBz, 0o644)
	require.NoError(t, err)
	return mockGenesisPath
}

func setupMockCLGenesis(t *testing.T, tmpDir string) string {
	t.Helper()
	// Create config directory in the root of tmpDir
	configDir := filepath.Join(tmpDir, "config")
	require.NoError(t, os.MkdirAll(configDir, 0o755))
	mockCLGenesisPath := filepath.Join(configDir, "genesis.json")

	mockCLGenesis := map[string]interface{}{
		"app_state": map[string]interface{}{
			"beacon": map[string]interface{}{
				"deposits": []interface{}{
					map[string]interface{}{
						"data": map[string]interface{}{
							"amount":               "32000000000",
							"pubkey":               "0x1234",
							"withdrawal_address":   "0x5678",
							"signature":            "0x9abc",
							"deposit_message_root": "0xdef0",
						},
					},
				},
			},
		},
	}
	clGenesisBz, err := json.MarshalIndent(mockCLGenesis, "", "  ")
	require.NoError(t, err)
	err = afero.WriteFile(afero.NewOsFs(), mockCLGenesisPath, clGenesisBz, 0o644)
	require.NoError(t, err)
	return mockCLGenesisPath
}

func verifyStorageOutput(t *testing.T, genesisPath string) {
	t.Helper()
	chainSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	depositAddr := common.Address(chainSpec.DepositContractAddress())

	// Verify the output genesis file
	outputBz, err := afero.ReadFile(afero.NewOsFs(), genesisPath)
	require.NoError(t, err)

	var output map[string]interface{}
	err = json.Unmarshal(outputBz, &output)
	require.NoError(t, err)

	// Check that the deposit contract storage was set correctly
	alloc, ok := output["alloc"].(map[string]interface{})
	require.True(t, ok)
	depositContract, ok := alloc[depositAddr.Hex()].(map[string]interface{})
	require.True(t, ok)
	storage, ok := depositContract["storage"].(map[string]interface{})
	require.True(t, ok)

	// Verify storage slots
	slot0 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	slot1 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")
	require.NotEmpty(t, storage[slot0.Hex()]) // Deposit count
	require.NotEmpty(t, storage[slot1.Hex()]) // Deposit root
}
