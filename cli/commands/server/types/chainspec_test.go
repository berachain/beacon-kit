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

package types_test

import (
	"os"
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/stretchr/testify/require"
)

// dummyAppOptions is a simple implementation of the AppOptions interface for testing.
type dummyAppOptions struct {
	values map[string]interface{}
}

func (d dummyAppOptions) Get(key string) interface{} {
	return d.values[key]
}

func TestCreateChainSpec_Devnet(t *testing.T) {
	// Set the env variable to force the devnet branch.
	t.Setenv(types.ChainSpecTypeEnvVar, types.DevnetChainSpecType)
	opts := dummyAppOptions{values: map[string]interface{}{}}
	cs, err := types.CreateChainSpec(opts)
	require.NoError(t, err)
	require.NotNil(t, cs)
	devnetSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, devnetSpec, "expected devnet chain spec to match")
}

func TestCreateChainSpec_Testnet(t *testing.T) {
	// Set the env variable to force the testnet branch.
	t.Setenv(types.ChainSpecTypeEnvVar, types.TestnetChainSpecType)
	opts := dummyAppOptions{values: map[string]interface{}{}}
	cs, err := types.CreateChainSpec(opts)
	require.NoError(t, err)
	require.NotNil(t, cs)
	testnetSpec, err := spec.TestnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, testnetSpec, "expected testnet chain spec to match")
}

func TestCreateChainSpec_Mainnet(t *testing.T) {
	// Set the env variable to force the mainnet branch.
	t.Setenv(types.ChainSpecTypeEnvVar, types.MainnetChainSpecType)
	opts := dummyAppOptions{values: map[string]interface{}{}}
	cs, err := types.CreateChainSpec(opts)
	require.NoError(t, err)
	require.NotNil(t, cs)
	mainnetSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, mainnetSpec, "expected mainnet chain spec to match")
}

//nolint:paralleltest // uses envars
func TestCreateChainSpec_Default_NoSpecFlag(t *testing.T) {
	// Ensure the env variable is unset so that the default branch is taken.
	err := os.Unsetenv(types.ChainSpecTypeEnvVar)
	require.NoError(t, err)
	// Provide an empty AppOptions so that no spec flag is present.
	opts := dummyAppOptions{values: map[string]interface{}{}}
	cs, err := types.CreateChainSpec(opts)
	mainnetSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, mainnetSpec, "expected mainnet chain spec to match")
}

//nolint:paralleltest // uses envars
func TestCreateChainSpec_ConfigurableEnvar_WithSpecFlag(t *testing.T) {
	// Ensure the env variable is unset so that the default branch is taken.
	err := os.Unsetenv(types.ChainSpecTypeEnvVar)
	require.NoError(t, err)
	// Provide a non-empty value for the configurable spec flag.
	opts := dummyAppOptions{values: map[string]interface{}{
		types.FlagConfigurableChainSpecPath: "mainnet_spec.toml",
	}}
	cs, err := types.CreateChainSpec(opts)
	require.NoError(t, err)

	mainnetSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, mainnetSpec, "the chain spec loaded from TOML does not match the mainnet spec")
}
