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

package indexes

import (
	"errors"
	"fmt"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

// Unique identifies an index that imposes uniqueness constraints on the reference key.
// It creates relationships between reference and primary key of the value.
type Unique[ReferenceKey, PrimaryKey, Value any] struct {
	getRefKey func(PrimaryKey, Value) (ReferenceKey, error)
	refKeys   collections.Map[ReferenceKey, PrimaryKey]
}

// NewUnique instantiates a new Unique index.
func NewUnique[ReferenceKey, PrimaryKey, Value any](
	storeKey []byte,
	keyPrefix []byte,
	refCodec codec.KeyCodec[ReferenceKey],
	pkCodec codec.KeyCodec[PrimaryKey],
	sa collections.StoreAccessor,
	getRefKeyFunc func(pk PrimaryKey, v Value) (ReferenceKey, error),
) *Unique[ReferenceKey, PrimaryKey, Value] {
	return &Unique[ReferenceKey, PrimaryKey, Value]{
		getRefKey: getRefKeyFunc,
		refKeys:   collections.NewMap(storeKey, keyPrefix, refCodec, codec.KeyToValueCodec(pkCodec), sa),
	}
}

func (unique *Unique[ReferenceKey, PrimaryKey, Value]) Reference(pk PrimaryKey, newValue Value, lazyOldValue func() (Value, error)) error {
	oldValue, err := lazyOldValue()
	switch {
	// if no error it means the value existed, and we need to remove the old indexes
	case err == nil:
		err = unique.unreference(pk, oldValue)
		if err != nil {
			return err
		}
	// if error is ErrNotFound, it means that the object does not exist, so we're creating indexes for the first time.
	// we do nothing.
	case errors.Is(err, collections.ErrNotFound):
	// default case means that there was some other error
	default:
		return err
	}
	// create new indexes, asserting no uniqueness constraint violation
	refKey, err := unique.getRefKey(pk, newValue)
	if err != nil {
		return err
	}
	has, err := unique.refKeys.Has(refKey)
	if err != nil {
		return err
	}
	if has {
		return fmt.Errorf("%w: index uniqueness constrain violation: %s", sdkcollections.ErrConflict, unique.refKeys.KeyCodec.Stringify(refKey))
	}
	return unique.refKeys.Set(refKey, pk)
}

func (unique *Unique[ReferenceKey, PrimaryKey, Value]) Unreference(pk PrimaryKey, getValue func() (Value, error)) error {
	value, err := getValue()
	if err != nil {
		return err
	}
	return unique.unreference(pk, value)
}

func (unique *Unique[ReferenceKey, PrimaryKey, Value]) unreference(pk PrimaryKey, value Value) error {
	refKey, err := unique.getRefKey(pk, value)
	if err != nil {
		return err
	}
	return unique.refKeys.Remove(refKey)
}

func (unique *Unique[ReferenceKey, PrimaryKey, Value]) MatchExact(ref ReferenceKey) (PrimaryKey, error) {
	return unique.refKeys.Get(ref)
}

func (unique *Unique[ReferenceKey, PrimaryKey, Value]) Iterate() (UniqueIterator[ReferenceKey, PrimaryKey], error) {
	iter, err := unique.refKeys.Iterate()
	return (UniqueIterator[ReferenceKey, PrimaryKey])(iter), err
}

// UniqueIterator is an Iterator wrapper, that exposes only the functionality needed to work with Unique keys.
type UniqueIterator[ReferenceKey, PrimaryKey any] sdkcollections.Iterator[ReferenceKey, PrimaryKey]

// PrimaryKey returns the iterator's current primary key.
func (unique UniqueIterator[ReferenceKey, PrimaryKey]) PrimaryKey() (PrimaryKey, error) {
	return (sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).Value()
}

// PrimaryKeys fully consumes the iterator, and returns all the primary keys.
func (unique UniqueIterator[ReferenceKey, PrimaryKey]) PrimaryKeys() ([]PrimaryKey, error) {
	return (sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).Values()
}

// FullKey returns the iterator's current full reference key as Pair[ReferenceKey, PrimaryKey].
func (unique UniqueIterator[ReferenceKey, PrimaryKey]) FullKey() (sdkcollections.Pair[ReferenceKey, PrimaryKey], error) {
	kv, err := (sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).KeyValue()
	return sdkcollections.Join(kv.Key, kv.Value), err
}

func (unique UniqueIterator[ReferenceKey, PrimaryKey]) FullKeys() ([]sdkcollections.Pair[ReferenceKey, PrimaryKey], error) {
	kvs, err := (sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).KeyValues()
	if err != nil {
		return nil, err
	}
	pairKeys := make([]sdkcollections.Pair[ReferenceKey, PrimaryKey], len(kvs))
	for index := range kvs {
		kv := kvs[index]
		pairKeys[index] = sdkcollections.Join(kv.Key, kv.Value)
	}
	return pairKeys, nil
}

func (unique UniqueIterator[ReferenceKey, PrimaryKey]) Next() {
	(sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).Next()
}

func (unique UniqueIterator[ReferenceKey, PrimaryKey]) Valid() bool {
	return (sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).Valid()
}

func (unique UniqueIterator[ReferenceKey, PrimaryKey]) Close() error {
	return (sdkcollections.Iterator[ReferenceKey, PrimaryKey])(unique).Close()
}
