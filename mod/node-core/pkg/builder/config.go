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

package builder

import (
	"time"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/config/pkg/template"
	beacon "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module"
	beaconv1alpha1 "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module/api/module/v1alpha1"
	cmtcfg "github.com/cometbft/cometbft/config"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

// DefaultAppConfig returns the default configuration for the application.
func DefaultAppConfig() any {
	// Define a struct for the custom app configuration.
	type CustomAppConfig struct {
		serverconfig.Config
		BeaconKit *config.Config `mapstructure:"beacon-kit"`
	}

	// Start with the default server configuration.
	cfg := serverconfig.DefaultConfig()
	cfg.MinGasPrices = "0stake"
	cfg.Telemetry.Enabled = true
	cfg.IAVLCacheSize = 25000

	// BeaconKit forces PebbleDB as the database backend.
	cfg.AppDBBackend = "pebbledb"
	cfg.Pruning = "everything"

	// IAVL FastNode should ALWAYS be disabled on IAVL v1.x.
	cfg.IAVLDisableFastNode = true

	// Create the custom app configuration.
	customAppConfig := CustomAppConfig{
		Config:    *cfg,
		BeaconKit: config.DefaultConfig(),
	}

	return customAppConfig
}

// DefaultAppConfigTemplate returns the default configuration template for the
// application.
func DefaultAppConfigTemplate() string {
	return serverconfig.DefaultConfigTemplate +
		"\n" + template.TomlTemplate
}

// DefaultCometConfig returns the default configuration for the CometBFT
// consensus engine.
//
//nolint:mnd // magic numbers are fine here.
func DefaultCometConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	consensus := cfg.Consensus
	consensus.TimeoutPropose = 1750 * time.Millisecond
	consensus.TimeoutPrecommit = 1000 * time.Millisecond
	consensus.TimeoutPrevote = 1000 * time.Millisecond
	consensus.TimeoutCommit = 1250 * time.Millisecond

	// BeaconKit forces PebbleDB as the database backend.
	cfg.DBBackend = "pebbledb"

	// These settings are set by default for performance reasons.
	cfg.TxIndex.Indexer = "null"
	cfg.Mempool.Type = "nop"
	cfg.Mempool.Size = 0
	cfg.Mempool.Recheck = false
	cfg.Mempool.Broadcast = false
	cfg.Storage.DiscardABCIResponses = true
	cfg.Storage.DiscardABCIResponses = true
	cfg.Instrumentation.Prometheus = true

	cfg.P2P.MaxNumInboundPeers = 100
	cfg.P2P.MaxNumOutboundPeers = 40
	return cfg
}

// DefaultDepInjectConfig returns the default configuration for the dependency
// injection framework.
func DefaultDepInjectConfig() depinject.Config {
	addrCdc := addresscodec.NewBech32Codec("bera")
	return depinject.Configs(
		appconfig.Compose(&appv1alpha1.Config{
			Modules: []*appv1alpha1.ModuleConfig{
				{
					Name: runtime.ModuleName,
					Config: appconfig.WrapAny(&runtimev1alpha1.Module{
						AppName:       DefaultAppName,
						PreBlockers:   []string{},
						BeginBlockers: []string{},
						EndBlockers:   []string{beacon.ModuleName},
						InitGenesis:   []string{beacon.ModuleName},
					}),
				},
				{
					Name:   beacon.ModuleName,
					Config: appconfig.WrapAny(&beaconv1alpha1.Module{}),
				},
			},
		}),
		depinject.Supply(
			func() address.Codec { return addrCdc },
			func() address.ValidatorAddressCodec { return addrCdc },
			func() address.ConsensusAddressCodec { return addrCdc },
		),
	)
}
