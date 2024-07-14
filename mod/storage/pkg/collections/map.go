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

package collections

import (
	"bytes"
	"errors"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

// A MapKeeper works as an intermediary that holds a certain configuration,
// and uses that configuration to interact with the underlying store.
type MapKeeper[K, V any] struct {
	storeKey      []byte
	keyPrefix     []byte
	KeyCodec      codec.KeyCodec[K]
	ValueCodec    codec.ValueCodec[V]
	storeAccessor StoreAccessor
	Size          int // TODO: remove this, temp field for debugging
}

func NewMapKeeper[K, V any](
	storeKey []byte,
	keyPrefix []byte,
	keyCodec codec.KeyCodec[K],
	valueCodec codec.ValueCodec[V],
	storeAccessor StoreAccessor,
) MapKeeper[K, V] {
	return MapKeeper[K, V]{
		storeKey:      storeKey,
		keyPrefix:     keyPrefix,
		KeyCodec:      keyCodec,
		ValueCodec:    valueCodec,
		storeAccessor: storeAccessor,
	}
}

// Get retrieves the value from the store, and returns the decoded value.
func (m *MapKeeper[K, V]) Get(key K) (V, error) {
	var result V
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return result, err
	}
	res, err := m.storeAccessor().QueryState(m.storeKey, prefixedKey)
	if err != nil {
		return result, err
	}
	return m.ValueCodec.Decode(res)
}

// Set sets the value in the store with the encoded key and value.
func (m *MapKeeper[K, V]) Set(key K, value V) error {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return err
	}
	encodedValue, err := m.ValueCodec.Encode(value)
	if err != nil {
		return err
	}
	m.storeAccessor().AddChange(m.storeKey, prefixedKey, encodedValue)
	m.Size++
	return nil
}

// Has reports whether the key is present in storage or not.
// Errors with ErrEncoding if key encoding fails.
func (m *MapKeeper[K, V]) Has(key K) (bool, error) {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return false, err
	}
	_, err = m.storeAccessor().QueryState(m.storeKey, prefixedKey)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Remove removes the key from the storage.
// Errors with ErrEncoding if key encoding fails.
// If the key does not exist then this is a no-op.
func (m *MapKeeper[K, V]) Remove(key K) error {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return err
	}
	m.storeAccessor().AddChange(m.storeKey, prefixedKey, nil)
	m.Size--
	return nil
}

// Iterate provides an Iterator over K and V in ascending order.
func (m *MapKeeper[K, V]) Iterate() (sdkcollections.Iterator[K, V], error) {
	return m.IterateRaw(nil, nil)
}

// TODO: remove this when not needed
func (m *MapKeeper[K, V]) NumKeys() (int, error) {
	// get latest reader map
	_, readerMap, err := m.storeAccessor().StateLatest()
	if err != nil {
		return 0, err
	}
	// retrieve reader from reader map
	reader, err := readerMap.GetReader(m.storeKey)
	if err != nil {
		return 0, err
	}

	prefixedStart := m.keyPrefix
	prefixedEnd := sdkcollections.NextBytesPrefixKey(m.keyPrefix)

	// retrieve iterator from reader
	iter, err := reader.Iterator(prefixedStart, prefixedEnd)
	if err != nil {
		return 0, err
	}
	defer iter.Close()
	count := 0
	for ; iter.Valid(); iter.Next() {
		count++
	}
	return count, nil
}

// IterateRaw iterates over the collection. The iteration range is untyped, it uses raw
// bytes. The resulting Iterator is typed.
// A nil start iterates from the first key contained in the collection.
// A nil end iterates up to the last key contained in the collection.
// A nil start and a nil end iterates over every key contained in the collection.
func (m *MapKeeper[K, V]) IterateRaw(
	start, end []byte,
) (sdkcollections.Iterator[K, V], error) {
	// prepend start/end range with keyPrefix
	prefixedStart := append(m.keyPrefix, start...)
	var prefixedEnd []byte
	if end == nil {
		prefixedEnd = sdkcollections.NextBytesPrefixKey(m.keyPrefix)
	} else {
		prefixedEnd = append(m.keyPrefix, end...)
	}

	// input validation
	if bytes.Compare(prefixedStart, prefixedEnd) == 1 {
		return sdkcollections.Iterator[K, V]{}, sdkcollections.ErrInvalidIterator
	}

	iter, err := m.storeAccessor().Iterator(prefixedStart, prefixedEnd)
	if err != nil {
		return sdkcollections.Iterator[K, V]{}, err
	}

	// return iterator with key/value codecs and reader iterator
	return sdkcollections.Iterator[K, V]{
		KeyCodec:     m.KeyCodec,
		ValueCodec:   m.ValueCodec,
		Iter:         iter,
		PrefixLength: len(m.keyPrefix),
	}, nil
}
