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
	"github.com/berachain/beacon-kit/config"
	serverconfig "github.com/berachain/beacon-kit/config/config"
	"github.com/berachain/beacon-kit/config/template"
)

// DefaultAppConfigTemplate returns the default configuration template for the
// application.
func DefaultAppConfigTemplate() string {
	return serverconfig.DefaultConfigTemplate +
		"\n" + template.TomlTemplate
}

// DefaultAppConfig returns the default configuration for the application.
func DefaultAppConfig() any {
	// Define a struct for the custom app configuration.
	type CustomAppConfig struct {
		serverconfig.Config
		BeaconKit *config.Config `mapstructure:"beacon-kit"`
	}

	// Start with the default server configuration.
	cfg := serverconfig.DefaultConfig()
	cfg.Telemetry.Enabled = true

	// BeaconKit forces PebbleDB as the database backend.
	cfg.Pruning = "everything"

	// IAVL FastNode should ALWAYS be disabled on IAVL v1.x.
	cfg.IAVLDisableFastNode = true
	cfg.IAVLCacheSize = 2500

	// Create the custom app configuration.
	customAppConfig := CustomAppConfig{
		Config:    *cfg,
		BeaconKit: config.DefaultConfig(),
	}

	return customAppConfig
}
