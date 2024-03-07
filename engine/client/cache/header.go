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
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
)

// HeaderCache is a specific cache for headers,
// by indexing both header number and hash.
type HeaderCache struct {
	// headerByNumberCache is the LRU cache for headers by number.
	headerByNumberCache *LRUCache[uint64, *ethcoretypes.Header]
	// headerByHashCache is the LRU cache for headers by hash.
	headerByHashCache *LRUCache[ethcommon.Hash, *ethcoretypes.Header]
}

// NewHeaderCache creates a new HeaderCache.
func NewHeaderCache() *HeaderCache {
	return &HeaderCache{
		headerByNumberCache: NewLRUCache[uint64, *ethcoretypes.Header](),
		headerByHashCache:   NewLRUCache[ethcommon.Hash, *ethcoretypes.Header](),
	}
}

// GetByNumber returns the header by number.
func (c *HeaderCache) GetByNumber(
	number uint64,
) (*ethcoretypes.Header, bool) {
	return c.headerByNumberCache.Get(number)
}

// GetByHash returns the header by hash.
func (c *HeaderCache) GetByHash(
	hash ethcommon.Hash,
) (*ethcoretypes.Header, bool) {
	return c.headerByHashCache.Get(hash)
}

// Add adds the header to the cache.
func (c *HeaderCache) Add(header *ethcoretypes.Header) {
	number := header.Number.Uint64()
	if oldHeader, ok := c.headerByNumberCache.Get(number); ok {
		c.headerByHashCache.Remove(oldHeader.Hash())
	}
	c.headerByNumberCache.Add(header.Number.Uint64(), header)
	c.headerByHashCache.Add(header.Hash(), header)
}
