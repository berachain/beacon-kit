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
}

// DefaultConfig returns the default preconfirmation configuration.
func DefaultConfig() Config {
	return Config{
		Enabled:       false,
		SequencerMode: false,
		WhitelistPath: "",
	}
}

// IsSequencer returns true if this node is configured as the sequencer.
func (c *Config) IsSequencer() bool {
	return c.Enabled && c.SequencerMode
}
