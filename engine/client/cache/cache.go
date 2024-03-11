// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package cache

import (
	"github.com/berachain/beacon-kit/primitives"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

// EngineCache is a cache for data retrieved by the EngineClient.
type EngineCache struct {
	headerByNumberCache *lru.LRU[uint64, *ethcoretypes.Header]
	headerByHashCache   *lru.LRU[primitives.ExecutionHash, *ethcoretypes.Header]
}

// NewEngineCacheWithConfig creates a new EngineCache with the given config.
func NewEngineCache(
	config Config,
) *EngineCache {
	return &EngineCache{
		headerByNumberCache: lru.NewLRU[uint64, *ethcoretypes.Header](
			config.HeaderSize,
			nil,
			config.HeaderTTL,
		),
		headerByHashCache: lru.NewLRU[primitives.ExecutionHash, *ethcoretypes.Header](
			config.HeaderSize,
			nil,
			config.HeaderTTL,
		),
	}
}

// NewEngineCache creates a new EngineCache.
func NewEngineCacheWithDefaultConfig() *EngineCache {
	return NewEngineCache(DefaultConfig())
}

// HeaderByNumber returns the header with the given number.
func (c *EngineCache) HeaderByNumber(
	number uint64,
) (*ethcoretypes.Header, bool) {
	return c.headerByNumberCache.Get(number)
}

// HeaderByHash returns the header with the given hash.
func (c *EngineCache) HeaderByHash(
	hash primitives.ExecutionHash,
) (*ethcoretypes.Header, bool) {
	return c.headerByHashCache.Get(hash)
}

// AddHeader adds the given header to the cache.
func (c *EngineCache) AddHeader(
	header *ethcoretypes.Header,
) {
	number := header.Number.Uint64()
	if oldHeader, ok := c.headerByNumberCache.Get(number); ok {
		c.headerByHashCache.Remove(oldHeader.Hash())
	}
	c.headerByNumberCache.Add(number, header)
	c.headerByHashCache.Add(header.Hash(), header)
}
