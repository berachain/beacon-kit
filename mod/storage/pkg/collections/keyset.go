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
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

// KeySet builds on top of a Map and represents a collection retaining only a set of keys and no value.
type KeySet[K any] struct {
	m Map[K, sdkcollections.NoValue]
}

// NewKeySet returns a KeySet given a schema, prefix, name, and key codec.
func NewKeySet[K any](
	storeKey []byte,
	keyPrefix []byte,
	keyCodec codec.KeyCodec[K],
	storeAccessor StoreAccessor,
) KeySet[K] {
	return KeySet[K]{
		m: NewMap(storeKey, keyPrefix, keyCodec, NoValueCodec(), storeAccessor),
	}
}

// NoValueCodec returns the codec for NoValue.
func NoValueCodec() codec.ValueCodec[sdkcollections.NoValue] {
	return sdkcollections.NoValue{}
}

// Set adds the key to the KeySet. Errors on encoding problems.
func (k *KeySet[K]) Set(key K) error {
	return k.m.Set(key, sdkcollections.NoValue{})
}

// Has returns if the key is present in the KeySet.
// An error is returned only in case of encoding problems.
func (k *KeySet[K]) Has(key K) (bool, error) {
	return k.m.Has(key)
}

// Remove removes the key from the KeySet. An error is returned in case of encoding error.
func (k *KeySet[K]) Remove(key K) error {
	return k.m.Remove(key)
}

// Iterate iterates over the keys in the KeySet.
func (k *KeySet[K]) Iterate() (sdkcollections.KeySetIterator[K], error) {
	iter, err := (*Map[K, sdkcollections.NoValue])(&k.m).Iterate()
	if err != nil {
		return sdkcollections.KeySetIterator[K]{}, err
	}

	return (sdkcollections.KeySetIterator[K])(iter), nil
}

// IterateRaw iterates over the raw byte keys in the KeySet.
func (k *KeySet[K]) IterateRaw(start, end []byte) (sdkcollections.Iterator[K, sdkcollections.NoValue], error) {
	return k.m.IterateRaw(start, end)
}
