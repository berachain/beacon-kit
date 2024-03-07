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

const (
	// cacheSize is the size of the cache.
	cacheSize = 20
	// cacheEvictionTime is the time after which the LRU is evicted.
	cacheEvictionTime = 5 * time.Minute
)

// LRUCache is a generic LRU cache to store client's objects.
type LRUCache[K comparable, V any] struct {
	*lru.LRU[K, V]
}

// NewLRUCache creates a new LRU cache.
func NewLRUCache[K comparable, V any]() *LRUCache[K, V] {
	return &LRUCache[K, V]{
		lru.NewLRU[K, V](
			cacheSize,
			nil,
			cacheEvictionTime,
		),
	}
}
