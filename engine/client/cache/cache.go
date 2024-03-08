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
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcoretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// HeaderKey is the key for the header cache.
	HeaderKey = "header"
)

// EngineCache is a cache for data retrieved by the EngineClient.
type EngineCache struct {
	headerByNumberCache *LRU[uint64, *gethcoretypes.Header]
	headerByHashCache   *LRU[gethcommon.Hash, *gethcoretypes.Header]
}

// NewEngineCacheWithConfig creates a new EngineCache with the given config.
func NewEngineCacheWithConfig(
	config EngineCacheConfig,
) *EngineCache {
	// headerByNumberCache and headerByHashCache share the same config.
	headerByNumberCache := NewLRUWithConfig[uint64, *gethcoretypes.Header](
		config.Cfgs[HeaderKey],
	)
	headerByHashCache := NewLRUWithConfig[gethcommon.Hash, *gethcoretypes.Header](
		config.Cfgs[HeaderKey],
	)
	return &EngineCache{
		headerByNumberCache,
		headerByHashCache,
	}
}

// HeaderByNumber returns the header with the given number.
func (c *EngineCache) HeaderByNumber(
	number uint64,
) (*gethcoretypes.Header, bool) {
	return c.headerByNumberCache.Get(number)
}

// HeaderByHash returns the header with the given hash.
func (c *EngineCache) HeaderByHash(
	hash gethcommon.Hash,
) (*gethcoretypes.Header, bool) {
	return c.headerByHashCache.Get(hash)
}

// AddHeader adds the given header to the cache.
func (c *EngineCache) AddHeader(
	header *gethcoretypes.Header,
) {
	number := header.Number.Uint64()
	if oldHeader, ok := c.headerByNumberCache.Get(number); ok {
		c.headerByHashCache.Remove(oldHeader.Hash())
	}
	c.headerByNumberCache.Add(number, header)
	c.headerByHashCache.Add(header.Hash(), header)
}

// EngineCacheConfig is the configuration for an EngineCache.
type EngineCacheConfig struct {
	Cfgs map[string]LRUConfig
}

// Template returns the TOML template for the EngineCacheConfig.
func (t *EngineCacheConfig) Template() string {
	template := ""
	for k, v := range t.Cfgs {
		template += fmt.Sprintf(`
[engine-cache.%s]
%s
`, k, v.Template())
	}
	return template
}
