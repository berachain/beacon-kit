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

	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

// LRUConfig is the configuration for an LRU cache.
type LRUConfig struct {
	// Capacity is the maximum number of items that the LRU can hold.
	Capacity int `mapstructure:"capacity"`
	// Expiry is the time in seconds after which the LRU is evicted.
	Expiry int `mapstructure:"expiry"`
}

// LRU is an LRU cache.
type LRU[K comparable, V any] struct {
	*lru.LRU[K, V]
}

// NewLRUWithConfig creates a new LRU with the given config.
func NewLRUWithConfig[K comparable, V any](
	config LRUConfig,
) *LRU[K, V] {
	if config == (LRUConfig{}) {
		return nil
	}
	return &LRU[K, V]{
		lru.NewLRU[K, V](
			config.Capacity,
			nil,
			time.Duration(config.Expiry)*time.Second,
		),
	}
}

// Template returns the TOML template for the LRU config.
func (c LRUConfig) Template() string {
	return `
capacity = {{.EngineCache.Capacity}}
expiry = {{.EngineCache.Expiry}}
`
}
