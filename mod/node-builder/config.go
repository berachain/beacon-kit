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

package nodebuilder

import (
	"time"

	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

// DefaultAppConfig returns the default configuration for the application.
func (nb *NodeBuilder[T]) DefaultAppConfig() any {
	// Define a struct for the custom app configuration.
	type CustomAppConfig struct {
		serverconfig.Config
		BeaconKit *config.Config `mapstructure:"beacon-kit"`
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
		BeaconKit: config.DefaultConfig(),
	}

	return customAppConfig
}

// DefaultAppConfigTemplate returns the default configuration template for the
// application.
func (nb *NodeBuilder[T]) DefaultAppConfigTemplate() string {
	return serverconfig.DefaultConfigTemplate +
		"\n" + config.Template
}

// DefaultCometConfig returns the default configuration for the CometBFT
// consensus engine.
//
//nolint:mnd // magic numbers are fine here.
func (nb *NodeBuilder[T]) DefaultCometConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	consensus := cfg.Consensus
	consensus.TimeoutPropose = 2800 * time.Millisecond
	consensus.TimeoutVote = 3000 * time.Millisecond
	consensus.TimeoutCommit = 0 * time.Millisecond

	// BeaconKit forces PebbleDB as the database backend.
	cfg.DBBackend = "pebbledb"

	// These settings are set by default for performance reasons.
	cfg.TxIndex.Indexer = "null"
	cfg.Mempool.Type = "nop"
	cfg.Storage.DiscardABCIResponses = true
	return cfg
}
