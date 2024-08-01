// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	blockstore "github.com/berachain/beacon-kit/mod/beacon/block_store"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/config/pkg/template"
	cometbft "github.com/berachain/beacon-kit/mod/consensus/pkg/comet"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	log "github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-api/server"
	"github.com/berachain/beacon-kit/mod/payload/pkg/builder"
)

// TODO: remove the comet import here so config is generalizable.
// For now we're using it to handle consensus config with the rest
// of the beacon kit app config. (does this make sense?? maybe not)

// DefaultConfig returns the default configuration for a BeaconKit chain.
func DefaultConfig() *Config {
	return &Config{
		Engine:            engineclient.DefaultConfig(),
		Logger:            log.DefaultConfig(),
		KZG:               kzg.DefaultConfig(),
		PayloadBuilder:    builder.DefaultConfig(),
		Validator:         validator.DefaultConfig(),
		BlockStoreService: blockstore.DefaultConfig(),
		NodeAPI:           server.DefaultConfig(),
		CometBFT:          *cometbft.DefaultConfig(),
	}
}

// Config is the main configuration struct for the BeaconKit chain.
type Config struct {
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
	// CometBFT is the configuration for the CometBFT node.
	CometBFT cometbft.Config `mapstructure:"cometbft"`
}

func (c Config) Default() *Config {
	return DefaultConfig()
}

// GetEngine returns the execution client configuration.
func (c Config) GetEngine() *engineclient.Config {
	return &c.Engine
}

// GetLogger returns the logger configuration.
func (c Config) GetLogger() *log.Config {
	return &c.Logger
}

// Template returns the configuration template.
func (c Config) Template() string {
	return template.TomlTemplate
}
