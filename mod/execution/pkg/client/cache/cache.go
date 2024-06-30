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
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

// EngineCache is a cache for data retrieved by the EngineClient.
type EngineCache struct {
	// headerByNumberCache is an LRU cache that maps block numbers to their
	// corresponding headers.
	headerByNumberCache *lru.LRU[
		uint64, *engineprimitives.Header,
	]
	// headerByHashCache is an LRU cache that maps block hashes to their
	// corresponding headers.
	headerByHashCache *lru.LRU[
		gethprimitives.ExecutionHash, *engineprimitives.Header,
	]
}

// NewEngineCache creates a new EngineCache with the given config.
func NewEngineCache(
	config Config,
) *EngineCache {
	return &EngineCache{
		headerByNumberCache: lru.NewLRU[
			uint64, *engineprimitives.Header,
		](
			config.HeaderSize,
			nil,
			config.HeaderTTL,
		),
		headerByHashCache: lru.NewLRU[
			gethprimitives.ExecutionHash, *engineprimitives.Header,
		](
			config.HeaderSize,
			nil,
			config.HeaderTTL,
		),
	}
}

// NewEngineCacheWithDefaultConfig creates a new EngineCache.
func NewEngineCacheWithDefaultConfig() *EngineCache {
	return NewEngineCache(DefaultConfig())
}

// HeaderByNumber returns the header with the given number.
func (c *EngineCache) HeaderByNumber(
	number uint64,
) (*engineprimitives.Header, bool) {
	return c.headerByNumberCache.Get(number)
}

// HeaderByHash returns the header with the given hash.
func (c *EngineCache) HeaderByHash(
	hash gethprimitives.ExecutionHash,
) (*engineprimitives.Header, bool) {
	return c.headerByHashCache.Get(hash)
}

// AddHeader adds the given header to the cache.
func (c *EngineCache) AddHeader(
	header *engineprimitives.Header,
) {
	number := header.Number.Uint64()
	if oldHeader, ok := c.headerByNumberCache.Get(number); ok {
		c.headerByHashCache.Remove(oldHeader.Hash())
	}
	c.headerByNumberCache.Add(number, header)
	c.headerByHashCache.Add(header.Hash(), header)
}
