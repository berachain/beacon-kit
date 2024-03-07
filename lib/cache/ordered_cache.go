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
	"sync"

	"github.com/huandu/skiplist"
)

type Comparable[T any] interface {
	Compare(lhs, rhs T) int
}

// OrderedCache is a set of elements that
// are maintained in an ascending order.
type OrderedCache[T any] struct {
	store *skiplist.SkipList
	// mu is a mutex that protects the cache.
	mu sync.RWMutex
}

// NewOrderedCache returns a new ordered cache.
func NewOrderedCache[T any](c Comparable[T]) *OrderedCache[T] {
	ascendingOrder := skiplist.GreaterThanFunc(func(lhs, rhs any) int {
		return c.Compare(lhs.(T), rhs.(T))
	})
	return &OrderedCache[T]{
		store: skiplist.New(ascendingOrder),
	}
}

// Insert adds an element to the cache.
func (c *OrderedCache[T]) Insert(elem T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Set(elem, struct{}{})
}

// Remove removes an element from the cache.
func (c *OrderedCache[T]) Remove(elem T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Remove(elem)
}

// Contains returns true if the cache contains the element.
func (c *OrderedCache[T]) Contains(elem T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store.Get(elem) != nil
}

// Front returns the first (smallest) element in the cache.
func (c *OrderedCache[T]) Front() (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem := c.store.Front()
	if elem == nil {
		var zero T
		return zero, ErrEmptyCache
	}
	return elem.Key().(T), nil
}

// RemoveFront removes the first element in the cache.
func (c *OrderedCache[T]) RemoveFront() (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem := c.store.RemoveFront()
	if elem == nil {
		var zero T
		return zero, ErrEmptyCache
	}
	return elem.Key().(T), nil
}

// Back returns the last (largest) element in the cache.
func (c *OrderedCache[T]) Back() (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem := c.store.Back()
	if elem == nil {
		var zero T
		return zero, ErrEmptyCache
	}
	return elem.Key().(T), nil
}

// RemoveBack removes the last element in the cache.
func (c *OrderedCache[T]) RemoveBack() (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem := c.store.RemoveBack()
	if elem == nil {
		var zero T
		return zero, ErrEmptyCache
	}
	return elem.Key().(T), nil
}

// Len returns the number of elements in the cache.
func (c *OrderedCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.store.Len()
}
