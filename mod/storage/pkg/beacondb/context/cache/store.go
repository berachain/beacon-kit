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
	"errors"

	"cosmossdk.io/core/store"
)

var _ store.Writer = (*Store[store.Reader])(nil)

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store[T store.Reader] struct {
	cache     *tree // always ascending sorted
	changeSet *tree
	parent    T

	iteratorRangeCache *iteratorRangeCache[T]
}

// NewStore creates a new Store object
func NewStore[T store.Reader](parent T) *Store[T] {
	s := &Store[T]{
		cache:     newTree(),
		changeSet: newTree(),
		parent:    parent,
	}
	s.iteratorRangeCache = newIteratorRangeCache(parent, s.cache)
	return s
}

// Get implements types.KVStore.
func (s *Store[T]) Get(key []byte) (value []byte, err error) {
	value, found := s.cache.get(key)
	if found {
		return
	}
	value, err = s.parent.Get(key)
	if err != nil {
		return nil, err
	}

	// add the value into the cache.
	s.cache.set(key, value)
	return value, nil
}

// Set implements types.KVStore.
func (s *Store[T]) Set(key, value []byte) error {
	if value == nil {
		return errors.New("cannot set a nil value")
	}

	s.cache.set(key, value)
	s.changeSet.set(key, value)
	return nil
}

// Has implements types.KVStore.
func (s *Store[T]) Has(key []byte) (bool, error) {
	tmpValue, found := s.cache.get(key)
	if found {
		return tmpValue != nil, nil
	}
	return s.parent.Has(key)
}

// Delete implements types.KVStore.
func (s *Store[T]) Delete(key []byte) error {
	s.cache.delete(key)
	s.changeSet.delete(key)
	return nil
}

// ----------------------------------------
// Iteration

// Iterator implements types.KVStore.
func (s *Store[T]) Iterator(start, end []byte) (store.Iterator, error) {
	return s.iterator(start, end, true)
}

// ReverseIterator implements types.KVStore.
func (s *Store[T]) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return s.iterator(start, end, false)
}

func (s *Store[T]) iterator(
	start, end []byte,
	ascending bool,
) (store.Iterator, error) {
	// If the range has not been synced yet, sync it.
	if !s.iteratorRangeCache.Seen(start, end) {
		if err := s.iteratorRangeCache.SyncForRange(start, end); err != nil {
			return nil, err
		}
	}

	// Return the appropriate iterator.
	if ascending {
		return s.cache.Iterator(start, end)
	}
	return s.cache.ReverseIterator(start, end)
}

func (s *Store[T]) ApplyChangeSets(changes []store.KVPair) error {
	for _, c := range changes {
		if c.Remove {
			err := s.Delete(c.Key)
			if err != nil {
				return err
			}
		} else {
			err := s.Set(c.Key, c.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store[T]) ChangeSets() (cs []store.KVPair, err error) {
	cs = make([]store.KVPair, s.changeSet.size())
	iter, err := s.changeSet.Iterator(nil, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	i := 0
	for ; iter.Valid(); iter.Next() {
		k, v := iter.Key(), iter.Value()
		cs[i] = store.KVPair{
			Key:    k,
			Value:  v,
			Remove: v == nil, // maybe we can optimistically compute size.
		}
		i++
	}
	return cs, nil
}
