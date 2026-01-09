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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log/phuslu"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
)

// PreconfServerInput is the input for the preconf server provider.
type PreconfServerInput struct {
	depinject.In

	Cfg          *config.Config
	Logger       *phuslu.Logger
	Whitelist    preconf.Whitelist
	LocalBuilder *payloadbuilder.PayloadBuilder
}

// ProvidePreconfServer provides the preconf API server for sequencer mode.
// Returns nil if preconf is disabled or not in sequencer mode.
//
//nolint:nilnil // nil server indicates preconf server is disabled
func ProvidePreconfServer(in PreconfServerInput) (*preconf.Server, error) {
	cfg := &in.Cfg.Preconf
	logger := in.Logger.With("service", "preconf-server")

	// Only start server in sequencer mode
	if !cfg.IsSequencer() {
		return nil, nil
	}

	// Load validator JWTs
	if cfg.ValidatorJWTsPath == "" {
		return nil, errors.New("preconf sequencer mode enabled but validator-jwts-path is not set")
	}

	validatorJWTs, err := preconf.LoadValidatorJWTs(cfg.ValidatorJWTsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load validator JWTs from: %s", cfg.ValidatorJWTsPath)
	}

	if len(validatorJWTs) == 0 {
		return nil, errors.New("validator JWTs file is empty")
	}

	// Local builder is required for the server to retrieve payloads
	if in.LocalBuilder == nil {
		return nil, errors.New("preconf server requires a LocalBuilder")
	}

	logger.Info(
		"Preconf API server configuration loaded",
		"port", cfg.APIPort,
		"validator_count", len(validatorJWTs),
	)

	return preconf.NewServer(
		logger,
		validatorJWTs,
		in.Whitelist,
		in.LocalBuilder,
		cfg.APIPort,
	), nil
}
