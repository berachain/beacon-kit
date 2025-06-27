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

package config

import (
	"time"

	"github.com/berachain/beacon-kit/beacon/validator"
	"github.com/berachain/beacon-kit/config/template"
	viperlib "github.com/berachain/beacon-kit/config/viper"
	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/errors"
	engineclient "github.com/berachain/beacon-kit/execution/client"
	log "github.com/berachain/beacon-kit/log/phuslu"
	blockstore "github.com/berachain/beacon-kit/node-api/block_store"
	"github.com/berachain/beacon-kit/node-api/server"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	DefaultChainSpec         = "mainnet"
	DefaultChainSpecFilePath = ""
	defaultShutdownTimeout   = 5 * time.Minute
)

// AppOptions is from the SDK, we should look to remove its usage.
type AppOptions interface {
	Get(string) interface{}
}

// DefaultConfig returns the default configuration for a BeaconKit chain.
func DefaultConfig() *Config {
	return &Config{
		ChainSpec:         DefaultChainSpec,
		ChainSpecFilePath: DefaultChainSpecFilePath,
		ShutdownTimeout:   defaultShutdownTimeout,
		Engine:            engineclient.DefaultConfig(),
		Logger:            log.DefaultConfig(),
		KZG:               kzg.DefaultConfig(),
		PayloadBuilder:    builder.DefaultConfig(),
		Validator:         validator.DefaultConfig(),
		BlockStoreService: blockstore.DefaultConfig(),
		NodeAPI:           server.DefaultConfig(),
	}
}

// Config is the main configuration struct for the BeaconKit chain.
type Config struct {
	// ChainSpec is the type of chain spec to use.
	ChainSpec string `mapstructure:"chain-spec"`
	// ChainSpecFilePath is the path to the chain spec file to use.
	ChainSpecFilePath string `mapstructure:"chain-spec-file"`
	// ShutdownTimeout is the maximum time to wait for the node to gracefully shutdown before
	// forcing an exit.
	ShutdownTimeout time.Duration `mapstructure:"shutdown-timeout"`
	// Engine is the configuration for the execution client.
	Engine engineclient.Config `mapstructure:"engine"`
	// Logger is the configuration for the logger.
	Logger log.Config `mapstructure:"logger"`
	// KZG is the configuration for the KZG blob verifier.
	KZG kzg.Config `mapstructure:"kzg"`
	// PayloadBuilder is the configuration for the local build payload timeout.
	PayloadBuilder builder.Config `mapstructure:"payload-builder"`
	// Validator is the configuration for the validator client.
	Validator validator.Config `mapstructure:"validator"`
	// BlockStoreService is the configuration for the block store service.
	BlockStoreService blockstore.Config `mapstructure:"block-store-service"`
	// NodeAPI is the configuration for the node API.
	NodeAPI server.Config `mapstructure:"node-api"`
}

// GetEngine returns the execution client configuration.
func (c Config) GetEngine() *engineclient.Config {
	return &c.Engine
}

// GetPayloadBuilder returns the block store configuration.
func (c Config) GetPayloadBuilder() *builder.Config {
	return &c.PayloadBuilder
}

// GetBlockStoreService returns the block store configuration.
func (c Config) GetBlockStoreService() *blockstore.Config {
	return &c.BlockStoreService
}

// GetLogger returns the logger configuration.
func (c Config) GetLogger() *log.Config {
	return &c.Logger
}

// Template returns the configuration template.
func (c Config) Template() string {
	return template.TomlTemplate
}

// ReadConfigFromAppOpts reads the configuration options from the given
// application options.
func ReadConfigFromAppOpts(opts AppOptions) (*Config, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, errors.New("invalid application options type")
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
		return nil, err
	}

	return &cfg.BeaconKit, nil
}
