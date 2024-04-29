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

package config

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/beacon/validator"
	engineclient "github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/node-builder/components/kzg"
	"github.com/berachain/beacon-kit/mod/node-builder/config/flags"
	viperlib "github.com/berachain/beacon-kit/mod/node-builder/config/viper"
	"github.com/berachain/beacon-kit/mod/payload/builder"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DefaultConfig returns the default configuration for a BeaconKit chain.
func DefaultConfig() *Config {
	return &Config{
		Engine:         engineclient.DefaultConfig(),
		KZG:            kzg.DefaultConfig(),
		PayloadBuilder: builder.DefaultConfig(),
		Validator:      validator.DefaultConfig(),
	}
}

// Config is the main configuration struct for the BeaconKit chain.
type Config struct {
	// Engine is the configuration for the execution client.
	Engine engineclient.Config `mapstructure:"engine"`

	// KZG is the configuration for the KZG blob verifier.
	KZG kzg.Config `mapstructure:"kzg"`

	// PayloadBuilder is the configuration for the local build payload timeout.
	PayloadBuilder builder.Config `mapstructure:"payload-builder"`

	// Validator is the configuration for the validator client.
	Validator validator.Config `mapstructure:"validator"`
}

// Template returns the configuration template.
func (c Config) Template() string {
	return Template
}

// MustReadConfigFromAppOpts reads the configuration options from the given
// application options.
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
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, fmt.Errorf("invalid application options type: %T", opts)
	}

	type cfgUnmarshaller struct {
		BeaconKit Config `mapstructure:"beacon-kit"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDialURLFunc(),
		))); err != nil {
		return nil, fmt.Errorf(
			"failed to decode beacon-kit configuration: %w",
			err,
		)
	}

	return &cfg.BeaconKit, nil
}

// AddBeaconKitFlags implements servertypes.ModuleInitFlags interface.
func AddBeaconKitFlags(startCmd *cobra.Command) {
	defaultCfg := DefaultConfig()
	startCmd.Flags().String(
		flags.JWTSecretPath, defaultCfg.Engine.JWTSecretPath,
		"path to the execution client secret")
	startCmd.Flags().String(
		flags.RPCDialURL, defaultCfg.Engine.RPCDialURL.String(), "rpc dial url")
	startCmd.Flags().Uint64(
		flags.RPCRetries, defaultCfg.Engine.RPCRetries, "rpc retries")
	startCmd.Flags().Duration(
		flags.RPCTimeout, defaultCfg.Engine.RPCTimeout, "rpc timeout")
	startCmd.Flags().Duration(
		flags.RPCStartupCheckInterval,
		defaultCfg.Engine.RPCStartupCheckInterval,
		"rpc startup check interval")
	startCmd.Flags().Duration(flags.RPCJWTRefreshInterval,
		defaultCfg.Engine.RPCJWTRefreshInterval,
		"rpc jwt refresh interval")
	startCmd.Flags().Uint64(
		flags.RequiredChainID, defaultCfg.Engine.RequiredChainID,
		"required chain id")
	startCmd.Flags().String(flags.SuggestedFeeRecipient,
		defaultCfg.PayloadBuilder.SuggestedFeeRecipient.Hex(),
		"suggested fee recipient",
	)
	startCmd.Flags().String(flags.KZGTrustedSetupPath,
		defaultCfg.KZG.TrustedSetupPath,
		"kzg trusted setup path",
	)
}

// AddToSFlag adds the terms of service flag to the given command.
func AddToSFlag(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().Bool(
		flags.BeaconKitAcceptTos, false, "accept the terms of service")
}
