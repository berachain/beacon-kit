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

//nolint:govet,gomnd,lll // from sdk.
package root

import (
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"

	serverconfig "github.com/cosmos/cosmos-sdk/server/config"

	beaconconfig "github.com/itsdevbear/bolaris/config"
)

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	consensus := cfg.Consensus
	consensus.TimeoutPropose = time.Second * 5
	consensus.TimeoutPrevote = time.Second * 1
	consensus.TimeoutPrecommit = time.Second * 1
	consensus.TimeoutCommit = time.Second * 3

	cfg.DBBackend = "pebbledb"

	// Disable the indexer
	cfg.TxIndex.Indexer = "null"
	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.

	type CustomAppConfig struct {
		serverconfig.Config
		BeaconKit beaconconfig.Config `mapstructure:"beacon-kit"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In BeaconApp, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0stake"
	// srvCfg.BaseConfig.IAVLDisableFastNode = true // disable fastnode by default
	srvCfg.IAVLCacheSize = 10000

	srvCfg.Telemetry.Enabled = true
	srvCfg.API.Enable = true
	srvCfg.Telemetry.MetricsSink = "mem"

	srvCfg.AppDBBackend = "pebbledb"
	customAppConfig := CustomAppConfig{
		Config:    *srvCfg,
		BeaconKit: *beaconconfig.DefaultConfig(),
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate + customAppConfig.BeaconKit.Template()

	return customAppTemplate, customAppConfig
}
