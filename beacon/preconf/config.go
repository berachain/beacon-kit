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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package preconf

import "time"

const (
	// DefaultAPIPort is the default port for the preconf API server.
	DefaultAPIPort = 9090

	// DefaultFetchTimeout is the default timeout for the HTTP client when fetching payloads from sequencer.
	DefaultFetchTimeout = 500 * time.Millisecond
)

// Config holds preconfirmation configuration.
type Config struct {
	// Enabled is a global toggle for preconfirmation support.
	// If false, all preconf functionality is disabled regardless of other settings.
	Enabled bool `mapstructure:"enabled"`

	// SequencerMode indicates this node runs as the sequencer.
	// When true, triggers optimistic payload builds for whitelisted proposers.
	// Requires Enabled to be true.
	SequencerMode bool `mapstructure:"sequencer-mode"`

	// WhitelistPath is the path to the whitelist JSON file containing validator pubkeys.
	// Required when Enabled and SequencerMode are both true.
	WhitelistPath string `mapstructure:"whitelist-path"`

	// === Sequencer-side settings ===

	// ValidatorJWTsPath is the path to JSON file mapping validator pubkeys to JWT secrets.
	// Required when SequencerMode is true.
	ValidatorJWTsPath string `mapstructure:"validator-jwts-path"`

	// APIPort is the port for the preconf API server that validators connect to.
	// Required when SequencerMode is true.
	APIPort int `mapstructure:"api-port"`

	// === Validator-side settings ===

	// SequencerURL is the URL of the sequencer's preconf API endpoint.
	// Required when this validator wants to fetch payloads from sequencer.
	SequencerURL string `mapstructure:"sequencer-url"`

	// SequencerJWTPath is the path to this validator's JWT secret for authenticating
	// with the sequencer.
	// Required when SequencerURL is set.
	SequencerJWTPath string `mapstructure:"sequencer-jwt-path"`

	// FetchTimeout is the timeout for fetching payloads from sequencer.
	FetchTimeout time.Duration `mapstructure:"fetch-timeout"`
}

// DefaultConfig returns the default preconfirmation configuration.
func DefaultConfig() Config {
	return Config{
		APIPort:      DefaultAPIPort,
		FetchTimeout: DefaultFetchTimeout,
	}
}

// IsSequencer returns true if this node is configured as the sequencer.
func (c *Config) IsSequencer() bool {
	return c != nil && c.Enabled && c.SequencerMode
}

// ShouldFetchFromSequencer returns true if this node should fetch payloads from sequencer.
func (c *Config) ShouldFetchFromSequencer() bool {
	return c != nil && c.Enabled && c.SequencerURL != ""
}
