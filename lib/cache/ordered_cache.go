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
	"github.com/huandu/skiplist"
)

type Elem interface {
	Compare(other Elem) int
}

// OrderedCache is a set of elements that
// are maintained in an ascending order.
type OrderedCache[E Elem] struct {
	cache *skiplist.SkipList
}

// NewOrderedCache returns a new ordered cache.
func NewOrderedCache[E Elem]() *OrderedCache[E] {
	lessThanFunc := skiplist.LessThanFunc(func(lhs, rhs any) int {
		return lhs.(Elem).Compare(rhs.(Elem))
	})
	return &OrderedCache[E]{
		cache: skiplist.New(lessThanFunc),
	}
}

// Insert adds an element to the cache.
func (c *OrderedCache[E]) Insert(elem E) {
	c.cache.Set(elem, struct{}{})
}

// Remove removes an element from the cache.
func (c *OrderedCache[E]) Remove(elem E) {
	c.cache.Remove(elem)
}

// Contains returns true if the cache contains the element.
func (c *OrderedCache[E]) Contains(elem E) bool {
	return c.cache.Get(elem) != nil
}

// Front returns the first (smallest) element in the cache.
func (c *OrderedCache[E]) Front() (E, error) {
	elem := c.cache.Front()
	if elem == nil {
		var zero E
		return zero, ErrEmptyCache
	}
	return elem.Value.(E), nil
}

// RemoveFront removes the first element in the cache.
func (c *OrderedCache[E]) RemoveFront() (E, error) {
	elem := c.cache.RemoveFront()
	if elem == nil {
		var zero E
		return zero, ErrEmptyCache
	}
	return elem.Value.(E), nil
}

// Back returns the last (largest) element in the cache.
func (c *OrderedCache[E]) Back() (E, error) {
	elem := c.cache.Back()
	if elem == nil {
		var zero E
		return zero, ErrEmptyCache
	}
	return elem.Value.(E), nil
}

// RemoveBack removes the last element in the cache.
func (c *OrderedCache[E]) RemoveBack() (E, error) {
	elem := c.cache.RemoveBack()
	if elem == nil {
		var zero E
		return zero, ErrEmptyCache
	}
	return elem.Value.(E), nil
}

// Len returns the number of elements in the cache.
func (c *OrderedCache[E]) Len() int {
	return c.cache.Len()
}
