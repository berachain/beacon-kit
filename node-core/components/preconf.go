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
)

// PreconfWhitelistInput is the input for the preconf whitelist provider.
type PreconfWhitelistInput struct {
	depinject.In
	Cfg    *config.Config
	Logger *phuslu.Logger
}

// ProvidePreconfWhitelist is a function that provides the module to the
// application. Returns an empty whitelist if preconf is disabled or not
// in sequencer mode.
//
//nolint:nilnil // nil whitelist indicates preconf is disabled
func ProvidePreconfWhitelist(in PreconfWhitelistInput) (preconf.Whitelist, error) {
	cfg := &in.Cfg.Preconf
	logger := in.Logger.With("service", "preconf")

	if !cfg.Enabled {
		logger.Info("Preconfirmation support is disabled")
		return nil, nil
	}

	// Sequencer mode requires whitelist
	if cfg.SequencerMode && cfg.WhitelistPath == "" {
		return nil, errors.New("preconf sequencer mode enabled but whitelist-path is not set")
	}

	// Load whitelist if path is provided (both sequencer and validators need it)
	// - Sequencer: to know which proposers to build for
	// - Validators: to know if they should fetch from sequencer
	if cfg.WhitelistPath == "" {
		logger.Info("Preconfirmation enabled but no whitelist configured")
		return nil, nil
	}

	pubkeys, err := preconf.LoadWhitelist(cfg.WhitelistPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load preconf whitelist from: %s", cfg.WhitelistPath)
	}

	if len(pubkeys) == 0 {
		return nil, errors.New("preconf whitelist is empty")
	}

	if cfg.SequencerMode {
		logger.Info(
			"Preconf sequencer mode enabled",
			"whitelist_count", len(pubkeys),
			"whitelist_path", cfg.WhitelistPath,
		)
	} else {
		logger.Info(
			"Preconf whitelist loaded for validator mode",
			"whitelist_count", len(pubkeys),
		)
	}

	return preconf.NewWhitelist(pubkeys), nil
}
