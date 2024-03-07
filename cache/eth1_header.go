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
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

const (
	// cacheSize is the size of the cache.
	cacheSize = 20
	// cacheEvictionTime is the time after which the LRU is evicted.
	cacheEvictionTime = 5 * time.Minute
)

type Eth1HeaderCache struct {
	// headerByNumberCache is the LRU cache for headers by number.
	headerByNumberCache *lru.LRU[uint64, *ethcoretypes.Header]
	// headerByHashCache is the LRU cache for headers by hash.
	headerByHashCache *lru.LRU[ethcommon.Hash, *ethcoretypes.Header]
}

// NewEth1HeaderCache creates a new Eth1HeaderCache.
func NewEth1HeaderCache() *Eth1HeaderCache {
	headerByNumberCache := lru.NewLRU[uint64, *ethcoretypes.Header](
		cacheSize,
		nil,
		cacheEvictionTime,
	)
	headerByHashCache := lru.NewLRU[ethcommon.Hash, *ethcoretypes.Header](
		cacheSize,
		nil,
		cacheEvictionTime,
	)
	return &Eth1HeaderCache{
		headerByNumberCache: headerByNumberCache,
		headerByHashCache:   headerByHashCache,
	}
}

// GetByNumber returns the header by number.
func (c *Eth1HeaderCache) GetByNumber(
	number uint64,
) (*ethcoretypes.Header, bool) {
	return c.headerByNumberCache.Get(number)
}

// GetByHash returns the header by hash.
func (c *Eth1HeaderCache) GetByHash(
	hash ethcommon.Hash,
) (*ethcoretypes.Header, bool) {
	return c.headerByHashCache.Get(hash)
}

// Add adds the header to the cache.
func (c *Eth1HeaderCache) Add(header *ethcoretypes.Header) {
	c.headerByNumberCache.Add(header.Number.Uint64(), header)
	c.headerByHashCache.Add(header.Hash(), header)
}
