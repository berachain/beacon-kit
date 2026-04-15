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

// ProvidePreconfProposerTracker provides a ProposerTracker shared between the
// blockchain service (writer) and the preconf server (reader).
func ProvidePreconfProposerTracker() preconf.ProposerTracker {
	return preconf.NewProposerTracker()
}

// PreconfWhitelistInput is the input for the preconf whitelist provider.
type PreconfWhitelistInput struct {
	depinject.In
	Cfg    *config.Config
	Logger *phuslu.Logger
}

// ProvidePreconfWhitelist provides the preconf whitelist to the application.
// Returns an empty whitelist for non-sequencer nodes (all IsWhitelisted checks return false).
func ProvidePreconfWhitelist(in PreconfWhitelistInput) (preconf.Whitelist, error) {
	cfg := &in.Cfg.Preconf
	logger := in.Logger.With("service", "preconf")

	// Only the sequencer needs a populated whitelist (to know which proposers to build for).
	// Validators don't need it — they opt in to preconf by setting sequencer-url.
	if !cfg.Enabled || !cfg.SequencerMode {
		return preconf.EmptyWhitelist(), nil
	}

	// Sequencer mode requires a local payload builder so that round-change
	// rebuilds and optimistic FCU/payload-fetch flows can run. Fail fast at
	// startup rather than silently no-opping at runtime.
	if !in.Cfg.PayloadBuilder.Enabled {
		return nil, errors.New("preconf sequencer mode requires payload-builder.enabled=true")
	}

	if cfg.WhitelistPath == "" {
		return nil, errors.New("preconf sequencer mode enabled but whitelist-path is not set")
	}

	wl, err := preconf.NewWhitelist(cfg.WhitelistPath)
	if err != nil {
		return nil, err
	}

	if wl.Len() == 0 {
		return nil, errors.New("preconf whitelist is empty")
	}

	logger.Info(
		"Preconf sequencer mode enabled",
		"whitelist_count", wl.Len(),
		"whitelist_path", cfg.WhitelistPath,
	)

	return wl, nil
}
