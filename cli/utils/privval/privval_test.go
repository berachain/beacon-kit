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

package privval_test

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/berachain/beacon-kit/cli/utils/privval"
	"github.com/berachain/beacon-kit/primitives/crypto"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtbls12381 "github.com/cometbft/cometbft/crypto/bls12381"
	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/stretchr/testify/require"
)

func testConfig(t *testing.T) *cmtcfg.Config {
	t.Helper()
	cfg := cmtcfg.DefaultConfig()
	cfg.SetRoot(t.TempDir())
	// beacond init creates the config dir before initializing the files.
	require.NoError(t, os.MkdirAll(filepath.Dir(cfg.NodeKeyFile()), 0o750))
	return cfg
}

func TestInitializeNodeValidatorFiles_GeneratesKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, keyType, wantType string
	}{
		{"bls", crypto.CometBLSType, cmtbls12381.KeyType},
		{"ed25519", cmted25519.KeyType, cmted25519.KeyType},
		{"default is ed25519", "", cmted25519.KeyType},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg := testConfig(t)
			nodeID, pubKey, err := privval.InitializeNodeValidatorFiles(cfg, tc.keyType)
			require.NoError(t, err)
			require.NotEmpty(t, nodeID)
			require.Equal(t, tc.wantType, pubKey.Type())
			require.FileExists(t, cfg.PrivValidatorKeyFile())
			require.FileExists(t, cfg.PrivValidatorStateFile())
			require.FileExists(t, cfg.NodeKeyFile())
		})
	}
}

func TestInitializeNodeValidatorFiles_PreservesExistingKey(t *testing.T) {
	t.Parallel()
	cfg := testConfig(t)

	_, pubKey, err := privval.InitializeNodeValidatorFiles(cfg, crypto.CometBLSType)
	require.NoError(t, err)

	// A nonzero signing state must survive re-initialization.
	const state = `{"height":"42","round":0,"step":3}`
	require.NoError(t, os.WriteFile(cfg.PrivValidatorStateFile(), []byte(state), 0o600))
	keyBefore, err := os.ReadFile(cfg.PrivValidatorKeyFile())
	require.NoError(t, err)

	_, pubKey2, err := privval.InitializeNodeValidatorFiles(cfg, crypto.CometBLSType)
	require.NoError(t, err)
	require.True(t, pubKey.Equals(pubKey2))

	keyAfter, err := os.ReadFile(cfg.PrivValidatorKeyFile())
	require.NoError(t, err)
	require.Equal(t, keyBefore, keyAfter)
	stateAfter, err := os.ReadFile(cfg.PrivValidatorStateFile())
	require.NoError(t, err)
	require.JSONEq(t, state, string(stateAfter))
}

// A stat error other than "not exist" must not generate a replacement key.
func TestInitializeNodeValidatorFiles_StatErrorDoesNotGenerate(t *testing.T) {
	t.Parallel()
	cfg := testConfig(t)

	// A self-referencing symlink makes os.Stat fail with ELOOP.
	keyPath := cfg.PrivValidatorKeyFile()
	require.NoError(t, os.Symlink(keyPath, keyPath))

	_, _, err := privval.InitializeNodeValidatorFiles(cfg, crypto.CometBLSType)
	require.ErrorIs(t, err, syscall.ELOOP)
	require.NoFileExists(t, cfg.PrivValidatorStateFile())
}

// Test that InitializeNodeValidatorFilesFromMnemonic rejects bad input and leaves no validator key behind.
func TestInitializeNodeValidatorFiles_Errors(t *testing.T) {
	t.Parallel()
	const mnemonic = "abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon abandon abandon about"
	tests := []struct {
		name, mnemonic, keyType, wantErr string
	}{
		{"unsupported key type", "", "secp256k1", "unsupported consensus key type"},
		{"invalid mnemonic", "not a mnemonic", cmted25519.KeyType, "invalid mnemonic"},
		{"bls rejects mnemonic", mnemonic, crypto.CometBLSType, "does not support mnemonic"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg := testConfig(t)
			_, _, err := privval.InitializeNodeValidatorFilesFromMnemonic(
				cfg, tc.mnemonic, tc.keyType,
			)
			require.ErrorContains(t, err, tc.wantErr)
			require.NoFileExists(t, cfg.PrivValidatorKeyFile())
		})
	}
}
