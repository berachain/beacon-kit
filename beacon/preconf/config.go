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

import (
	"net/url"
	"time"

	"github.com/berachain/beacon-kit/errors"
)

const (
	// DefaultAPIPort is the default port for the preconf API server.
	DefaultAPIPort = 9090

	// DefaultFetchTimeout is the default timeout for the HTTP client when fetching payloads from sequencer.
	DefaultFetchTimeout = 500 * time.Millisecond

	// DefaultHealthCheckInterval is how often the client probes the sequencer when it is unavailable.
	DefaultHealthCheckInterval = 10 * time.Second
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

	// TLSCertPath is the path to the TLS certificate file for the preconf API server.
	// When set together with TLSKeyPath, the server uses HTTPS.
	TLSCertPath string `mapstructure:"tls-cert-path"`

	// TLSKeyPath is the path to the TLS private key file for the preconf API server.
	// Must be set if TLSCertPath is set.
	TLSKeyPath string `mapstructure:"tls-key-path"`

	// === Validator-side settings ===

	// SequencerURL is the URL of the sequencer's preconf API endpoint.
	// Required when this validator wants to fetch payloads from sequencer.
	SequencerURL string `mapstructure:"sequencer-url"`

	// SequencerJWTPath is the path to this validator's JWT secret for authenticating
	// with the sequencer.
	// Required when SequencerURL is set.
	SequencerJWTPath string `mapstructure:"sequencer-jwt-path"`

	// SequencerCACertPath is the optional path to a CA certificate for verifying
	// the sequencer's TLS certificate. When set, only this CA is trusted.
	SequencerCACertPath string `mapstructure:"sequencer-ca-cert-path"`

	// FetchTimeout is the timeout for fetching payloads from sequencer.
	FetchTimeout time.Duration `mapstructure:"fetch-timeout"`

	// HealthCheckInterval is how often to probe the sequencer health endpoint when it becomes unavailable.
	HealthCheckInterval time.Duration `mapstructure:"health-check-interval"`
}

// DefaultConfig returns the default preconfirmation configuration.
func DefaultConfig() Config {
	return Config{
		APIPort:             DefaultAPIPort,
		FetchTimeout:        DefaultFetchTimeout,
		HealthCheckInterval: DefaultHealthCheckInterval,
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

// TLSEnabled returns true if both TLS cert and key paths are configured.
func (c *Config) TLSEnabled() bool {
	return c != nil && c.TLSCertPath != "" && c.TLSKeyPath != ""
}

// Validate checks the config for structural consistency.
func (c *Config) Validate() error {
	if (c.TLSCertPath != "") != (c.TLSKeyPath != "") {
		return errors.New("tls-cert-path and tls-key-path must both be set or both be empty")
	}
	if c.SequencerCACertPath != "" {
		if c.SequencerURL == "" {
			return errors.New("sequencer-ca-cert-path requires sequencer-url to be set")
		}
		u, err := url.Parse(c.SequencerURL)
		if err != nil {
			return errors.Wrapf(err, "invalid sequencer-url: %s", c.SequencerURL)
		}
		if u.Scheme != "https" {
			return errors.New("sequencer-url must use https:// scheme when sequencer-ca-cert-path is set")
		}
	}
	return nil
}
