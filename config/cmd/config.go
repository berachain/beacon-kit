// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

//nolint:gomnd // from sdk.
package cmd

import (
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	clientconfig "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	beaconconfig "github.com/itsdevbear/bolaris/config"
)

// InitCometBFTConfig customizes CometBFT Config values, falling back to
// defaults if no customization is needed.
func InitCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	consensus := cfg.Consensus
	consensus.TimeoutPropose = 3000 * time.Millisecond
	consensus.TimeoutPrevote = 2000 * time.Millisecond
	consensus.TimeoutPrecommit = 2000 * time.Millisecond
	consensus.TimeoutCommit = 3000 * time.Millisecond

	// BeaconKit forces PebbleDB as the database backend.
	cfg.DBBackend = "pebbledb"

	// Indexer is disabled to enhance performance.
	cfg.TxIndex.Indexer = "null"
	return cfg
}

// InitClientConfig sets up the default client configuration, allowing for
// overrides.
func InitClientConfig() (string, interface{}) {
	clientCfg := clientconfig.DefaultConfig()
	clientCfg.KeyringBackend = keyring.BackendTest
	return clientconfig.DefaultClientConfigTemplate, clientCfg
}

// InitAppConfig customizes the app configuration for BeaconKit, incorporating
// any necessary overrides.
func InitAppConfig() (string, interface{}) {
	// Define a struct for the custom app configuration.
	type CustomAppConfig struct {
		serverconfig.Config
		BeaconKit beaconconfig.Config `mapstructure:"beacon-kit"`
	}

	// Start with the default server configuration.
	cfg := serverconfig.DefaultConfig()
	cfg.MinGasPrices = "0stake"
	cfg.Telemetry.Enabled = true

	// BeaconKit forces PebbleDB as the database backend.
	cfg.AppDBBackend = "pebbledb"

	// Create the custom app configuration.
	customAppConfig := CustomAppConfig{
		Config:    *cfg,
		BeaconKit: *beaconconfig.DefaultConfig(),
	}

	// Combine the default template with the custom BeaconKit configuration.
	customAppTemplate := serverconfig.DefaultConfigTemplate +
		"\n" + customAppConfig.BeaconKit.Template()

	return customAppTemplate, customAppConfig
}
