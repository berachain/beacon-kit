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

package components

import (
	"crypto/x509"
	"net"
	"net/url"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log/phuslu"
)

// PreconfClientInput is the input for the preconf client provider.
type PreconfClientInput struct {
	depinject.In

	Cfg    *config.Config
	Logger *phuslu.Logger
}

// ProvidePreconfClient provides the preconf client for fetching payloads from sequencer.
// Returns nil if preconf is disabled or sequencer URL is not configured.
//
//nolint:nilnil // nil client indicates preconf client is disabled
func ProvidePreconfClient(in PreconfClientInput) (*preconf.Client, error) {
	cfg := &in.Cfg.Preconf
	logger := in.Logger.With("service", "preconf-client")

	// Only create client if configured to fetch from sequencer
	if !cfg.ShouldFetchFromSequencer() {
		return nil, nil
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid preconf configuration")
	}

	// Load JWT secret
	if cfg.SequencerJWTPath == "" {
		return nil, errors.New("preconf enabled with sequencer-url but sequencer-jwt-path is not set")
	}

	jwtSecret, err := preconf.LoadJWTSecret(cfg.SequencerJWTPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load sequencer JWT from: %s", cfg.SequencerJWTPath)
	}

	var caCertPool *x509.CertPool
	if cfg.SequencerCACertPath != "" {
		caCertPool, err = preconf.LoadCACert(cfg.SequencerCACertPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load sequencer CA cert from: %s", cfg.SequencerCACertPath)
		}
		logger.Info("TLS CA certificate pinning enabled", "ca_cert", cfg.SequencerCACertPath)
	}

	// Plaintext to a remote sequencer leaks JWTs and payloads. Loopback is fine
	// (dev/devnet), so only warn for non-loopback http URLs.
	if u, perr := url.Parse(cfg.SequencerURL); perr == nil && u.Scheme == "http" {
		host := u.Hostname()
		ip := net.ParseIP(host)
		loopback := host == "localhost" || (ip != nil && ip.IsLoopback())
		if !loopback {
			logger.Warn("preconf sequencer-url uses plaintext http to a non-loopback host, "+
				"JWT tokens and payloads will transit unencrypted", "sequencer_url", cfg.SequencerURL,
			)
		}
	}

	logger.Info(
		"Preconf client configured",
		"sequencer_url", cfg.SequencerURL,
		"fetch_timeout", cfg.FetchTimeout,
	)

	return preconf.NewClient(
		logger,
		cfg.SequencerURL,
		jwtSecret,
		cfg.FetchTimeout,
		caCertPool,
		cfg.MaxResponseSize,
	), nil
}
