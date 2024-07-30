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

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

// QueryCache is a cache for query contexts.
type QueryCache struct {
	// queryCtxCache is an LRU cache that maps slots to their corresponding
	// query contexts.
	queryCtxCache *lru.LRU[math.Slot, context.Context]
}

// NewQueryCache creates a new QueryCache with the given config.
func NewQueryCache(config Config) *QueryCache {
	return &QueryCache{
		queryCtxCache: lru.NewLRU[math.Slot, context.Context](
			config.QueryContextSize,
			nil,
			config.QueryContextTTL,
		),
	}
}

// NewQueryCacheWithDefaultConfig creates a new QueryCache with the default
// configuration.
func NewQueryCacheWithDefaultConfig() *QueryCache {
	return NewQueryCache(DefaultConfig())
}

// GetQueryContext returns the query context for the given slot.
func (c *QueryCache) GetQueryContext(slot math.Slot) (context.Context, bool) {
	return c.queryCtxCache.Get(slot)
}

// AddQueryContext adds a query context to the cache for the given slot.
//
//nolint:revive // ctx is the value in cache.
func (c *QueryCache) AddQueryContext(slot math.Slot, ctx context.Context) {
	c.queryCtxCache.Add(slot, ctx)
}
