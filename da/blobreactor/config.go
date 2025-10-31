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

package blobreactor

import "time"

const (
	// defaultRequestTimeout is the maximum time to wait for a blob response from a peer
	defaultRequestTimeout = 5 * time.Second
	// defaultMaxMessagesPerPeerPerSecond limits the rate of incoming blob requests per peer
	defaultMaxMessagesPerPeerPerSecond = 10.0 / 60.0
	// defaultMaxGlobalRequestsPerSecond limits the total rate of incoming blob requests from all peers
	defaultMaxGlobalRequestsPerSecond = 30.0 / 60.0

	// BurstMultiplier determines how much burst capacity to allow relative to the rate limit.
	BurstMultiplier = 2.0
)

type Config struct {
	RequestTimeout              time.Duration    `mapstructure:"request-timeout"`
	MaxMessagesPerPeerPerSecond float64          `mapstructure:"max-messages-per-peer-per-second"`
	MaxGlobalRequestsPerSecond  float64          `mapstructure:"max-global-requests-per-second"`
	Reputation                  ReputationConfig `mapstructure:"reputation"`
}

func DefaultConfig() Config {
	return Config{
		RequestTimeout:              defaultRequestTimeout,
		MaxMessagesPerPeerPerSecond: defaultMaxMessagesPerPeerPerSecond,
		MaxGlobalRequestsPerSecond:  defaultMaxGlobalRequestsPerSecond,
		Reputation:                  DefaultReputationConfig(),
	}
}

// WithDefaults returns a new Config with zero values replaced by defaults
func (c Config) WithDefaults() Config {
	defaults := DefaultConfig()

	if c.RequestTimeout == 0 {
		c.RequestTimeout = defaults.RequestTimeout
	}
	if c.MaxMessagesPerPeerPerSecond == 0 {
		c.MaxMessagesPerPeerPerSecond = defaults.MaxMessagesPerPeerPerSecond
	}
	if c.MaxGlobalRequestsPerSecond == 0 {
		c.MaxGlobalRequestsPerSecond = defaults.MaxGlobalRequestsPerSecond
	}
	// Apply defaults to reputation config
	c.Reputation = c.Reputation.WithDefaults()

	return c
}
