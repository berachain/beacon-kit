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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blobreactor

import "time"

const (
	// defaultFetchTimeout is the overall deadline for fetching the sidecars of one proposal at the tip (all
	// tiers, all peers). It is sized to the ~2s consensus round cadence (TimeoutPropose/Prevote are capped at
	// 2s): a fetch that cannot finish in time fails the round rather than stalling the validator set, and the
	// retried round comes with larger timeouts while the sidecars keep propagating.
	defaultFetchTimeout = 2 * time.Second

	// defaultRequestTimeout is the per-peer timeout for a single request/response round trip. It must be a
	// fraction of defaultFetchTimeout so a single slow peer cannot consume the whole tip budget: at 750ms the
	// by-root tier can try two or three peers within the 2s window before giving up.
	defaultRequestTimeout = 750 * time.Millisecond
)

// Config is the node-local (config.toml) configuration of the blob reactor.
type Config struct {
	// RequestTimeout is the per-peer timeout for one request/response round trip.
	RequestTimeout time.Duration `mapstructure:"request-timeout"`
	// FetchTimeout is the overall deadline for retrieving the sidecars of one block at the tip of the chain before giving up.
	FetchTimeout time.Duration `mapstructure:"fetch-timeout"`
}

func DefaultConfig() Config {
	return Config{
		RequestTimeout: defaultRequestTimeout,
		FetchTimeout:   defaultFetchTimeout,
	}
}
