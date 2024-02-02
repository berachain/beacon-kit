// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package config

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/cmd/flags"
	"github.com/spf13/cobra"
)

// Config is the main configuration struct for the Polaris chain.
type Config struct {
	// ExecutionClient is the configuration for the execution client.
	ExecutionClient ExecutionClient

	// Beacon is the configuration for the fork epochs.
	Beacon Beacon

	// Proposal is the configuration for the proposal handler.
	Proposal Proposal
}

// DefaultConfig returns the default configuration for a polaris chain.
func DefaultConfig() *Config {
	return &Config{
		ExecutionClient: ExecutionClient{
			RPCDialURL:      "http://localhost:8551",
			RPCTimeout:      5, //nolint:gomnd // default config.
			RPCRetries:      3, //nolint:gomnd // default config.
			JWTSecretPath:   "./app/jwt.hex",
			RequiredChainID: 7, //nolint:gomnd // default config.
		},
		Beacon: DefaultBeaconConfig(),
	}
}

// Client is the configuration struct for the execution client.
type ExecutionClient struct {
	// RPCDialURL is the HTTP url of the execution client JSON-RPC endpoint.
	RPCDialURL string
	// RPCTimeout is the RPC timeout for execution client requests.
	RPCTimeout uint64
	// RPCRetries is the number of retries before shutting down consensus client.
	RPCRetries uint64
	// JWTSecretPath is the path to the JWT secret.
	JWTSecretPath string
	// RequiredChainID is the chain id that the consensus client must be connected to.
	RequiredChainID uint64
}

// SetupCosmosConfig sets up the Cosmos SDK configuration to be compatible with the
// semantics of etheruem.
func SetupCosmosConfig() {
	// set the address prefixes
	config := sdk.GetConfig()

	// We use CoinType == 60 to match Ethereum.
	// This is not strictly necessary, though highly recommended.
	config.SetCoinType(60) //nolint:gomnd // its okay.
	config.SetPurpose(sdk.Purpose)
	config.Seal()
}

// MustReadConfigFromAppOpts reads the configuration options from the given
// application options. Panics if the configuration cannot be read.
func MustReadConfigFromAppOpts(opts servertypes.AppOptions) *Config {
	cfg, err := ReadConfigFromAppOpts(opts)
	if err != nil {
		panic(err)
	}
	return cfg
}

// ReadConfigFromAppOpts reads the configuration options from the given
// application options.
func ReadConfigFromAppOpts(opts servertypes.AppOptions) (*Config, error) {
	return readConfigFromAppOptsParser(AppOptionsParser{AppOptions: opts})
}

// TODO: cleanup parsing logic.
func readConfigFromAppOptsParser(parser AppOptionsParser) (*Config, error) {
	var err error
	conf := &Config{}

	if conf.ExecutionClient.RPCDialURL, err = parser.GetString(flags.RPCDialURL); err != nil {
		return nil, err
	}
	if conf.ExecutionClient.RPCRetries, err = parser.GetUint64(flags.RPCRetries); err != nil {
		return nil, err
	}
	if conf.ExecutionClient.RPCTimeout, err = parser.GetUint64(
		flags.RPCTimeout,
	); err != nil {
		return nil, err
	}
	if conf.ExecutionClient.JWTSecretPath, err = parser.GetString(
		flags.JWTSecretPath,
	); err != nil {
		return nil, err
	}
	if conf.ExecutionClient.RequiredChainID, err = parser.GetUint64(
		flags.RequiredChainID,
	); err != nil {
		return nil, err
	}

	if conf.Beacon.AltairForkEpoch, err = parser.GetEpoch(
		flags.AltairForkEpoch,
	); err != nil {
		return nil, err
	}

	if conf.Beacon.BellatrixForkEpoch, err = parser.GetEpoch(
		flags.BellatrixForkEpoch,
	); err != nil {
		return nil, err
	}

	if conf.Beacon.CapellaForkEpoch, err = parser.GetEpoch(
		flags.CapellaForkEpoch,
	); err != nil {
		return nil, err
	}

	if conf.Beacon.DenebForkEpoch, err = parser.GetEpoch(
		flags.DenebForkEpoch,
	); err != nil {
		return nil, err
	}

	if conf.Proposal.BeaconKitBlockPosition, err = parser.GetUint(
		flags.BeaconKitBlockPosition,
	); err != nil {
		return nil, err
	}

	return conf, nil
}

// AddBeaconKitFlags implements servertypes.ModuleInitFlags interface.
func AddBeaconKitFlags(startCmd *cobra.Command) {
	defaultCfg := DefaultConfig().ExecutionClient
	startCmd.Flags().String(flags.JWTSecretPath, defaultCfg.JWTSecretPath,
		"path to the execution client secret")
	startCmd.Flags().String(flags.RPCDialURL, defaultCfg.RPCDialURL, "rpc dial url")
	startCmd.Flags().Uint64(flags.RPCRetries, defaultCfg.RPCRetries, "rpc retries")
	startCmd.Flags().Uint64(flags.RPCTimeout, defaultCfg.RPCTimeout, "rpc timeout")
	startCmd.Flags().Uint64(flags.RequiredChainID, defaultCfg.RequiredChainID, "required chain id")
}
