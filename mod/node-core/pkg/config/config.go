// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package config

import (
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/errors"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/kzg"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config/flags"
	viperlib "github.com/berachain/beacon-kit/mod/node-core/pkg/config/viper"
	"github.com/berachain/beacon-kit/mod/payload/pkg/builder"
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

// GetEngine returns the execution client configuration.
func (c Config) GetEngine() engineclient.Config {
	return c.Engine
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
		return nil, errors.Newf("invalid application options type: %T", opts)
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
			viperlib.StringToConnectionURLFunc(),
		))); err != nil {
		return nil, errors.Newf(
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
	startCmd.Flags().String(flags.SuggestedFeeRecipient,
		defaultCfg.PayloadBuilder.SuggestedFeeRecipient.Hex(),
		"suggested fee recipient",
	)
	startCmd.Flags().String(flags.KZGTrustedSetupPath,
		defaultCfg.KZG.TrustedSetupPath,
		"kzg trusted setup path",
	)
	startCmd.Flags().String(flags.KZGImplementation,
		defaultCfg.KZG.Implementation,
		"kzg implementation")
}

// AddToSFlag adds the terms of service flag to the given command.
func AddToSFlag(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().Bool(
		flags.BeaconKitAcceptTos, false, "accept the terms of service")
}
