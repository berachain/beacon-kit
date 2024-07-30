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

package cache

import "time"

const (
	defaultQueryContextSize = 20
	defaultQueryContextTTL  = 0
)

// Config is the configuration for an EngineCache.
type Config struct {
	// QueryContextSize is the size of the query context cache.
	QueryContextSize int `mapstructure:"query-context-size"`
	// QueryContextTTL is the time-to-live for query contexts in the cache.
	QueryContextTTL time.Duration `mapstructure:"query-context-ttl"`
}

// DefaultConfig returns the default configuration for an QueryCache.
func DefaultConfig() Config {
	return Config{
		QueryContextSize: defaultQueryContextSize,
		QueryContextTTL:  defaultQueryContextTTL,
	}
}
