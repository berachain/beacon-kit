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

package spec_test

import (
	"testing"

	"github.com/berachain/beacon-kit/cli/flags"
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
	t.Parallel()

	// Set the app opts to force the devnet branch.
	opts := dummyAppOptions{values: map[string]interface{}{
		flags.ChainSpec: "devnet",
	}}
	cs, err := spec.Create(opts)
	require.NoError(t, err)
	require.NotNil(t, cs)
	devnetSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, devnetSpec, "expected devnet chain spec to match")
}

func TestCreateChainSpec_Testnet(t *testing.T) {
	t.Parallel()

	// Set the app opts to force the testnet branch.
	opts := dummyAppOptions{values: map[string]interface{}{
		flags.ChainSpec: "testnet",
	}}
	cs, err := spec.Create(opts)
	require.NoError(t, err)
	require.NotNil(t, cs)
	testnetSpec, err := spec.TestnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, testnetSpec, "expected testnet chain spec to match")
}

func TestCreateChainSpec_Mainnet(t *testing.T) {
	t.Parallel()

	// Set the app opts to force the mainnet branch.
	opts := dummyAppOptions{values: map[string]interface{}{
		flags.ChainSpec: "mainnet",
	}}
	cs, err := spec.Create(opts)
	require.NoError(t, err)
	require.NotNil(t, cs)
	mainnetSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, mainnetSpec, "expected mainnet chain spec to match")
}

func TestCreateChainSpec_Default_NoSpecFlag(t *testing.T) {
	t.Parallel()

	// Provide an empty app opts so that no spec flag is present.
	opts := dummyAppOptions{values: map[string]interface{}{}}
	cs, err := spec.Create(opts)
	require.NoError(t, err)
	mainnetSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, cs, mainnetSpec, "expected mainnet chain spec to match")
}

func TestCreateChainSpec_File(t *testing.T) {
	t.Parallel()

	// Provide a non-empty value for the custom spec file of mainnet.
	opts := dummyAppOptions{values: map[string]interface{}{
		flags.ChainSpec:         "file",
		flags.ChainSpecFilePath: "../../testing/networks/80094/spec.toml",
	}}
	mcs, err := spec.Create(opts)
	require.NoError(t, err)

	mainnetSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, mainnetSpec, mcs, "the chain spec loaded from TOML does not match the mainnet spec")

	// Provide a non-empty value for the custom spec file of testnet.
	opts.values[flags.ChainSpecFilePath] = "../../testing/networks/80069/spec.toml"
	tcs, err := spec.Create(opts)
	require.NoError(t, err)

	testnetSpec, err := spec.TestnetChainSpec()
	require.NoError(t, err)
	require.Equal(t, testnetSpec, tcs, "the chain spec loaded from TOML does not match the testnet spec")
}
