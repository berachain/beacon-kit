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

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

type multiOptions struct {
	uncheckedValue bool
}

// WithMultiUncheckedValue is an option that can be passed to NewMulti to
// ignore index values different from '[]byte{}' and continue with the operation.
// This should be used only to behave nicely in case you have used values different
// from '[]byte{}' in your storage before migrating to collections. Refer to
// WithKeySetUncheckedValue for more information.
func WithMultiUncheckedValue() func(*multiOptions) {
	return func(o *multiOptions) {
		o.uncheckedValue = true
	}
}

// Multi defines the most common index. It can be used to create a reference between
// a field of value and its primary key. Multiple primary keys can be mapped to the same
// reference key as the index does not enforce uniqueness constraints.
type Multi[ReferenceKey, PrimaryKey, Value any] struct {
	getRefKey func(pk PrimaryKey, value Value) (ReferenceKey, error)
	refKeys   collections.KeySet[sdkcollections.Pair[ReferenceKey, PrimaryKey]]
}

// NewMulti instantiates a new Multi instance given a schema,
// a Prefix, the humanized name for the index, the reference key key codec
// and the primary key key codec. The getRefKeyFunc is a function that
// given the primary key and value returns the referencing key.
func NewMulti[ReferenceKey, PrimaryKey, Value any](
	storeKey []byte,
	keyPrefix []byte,
	refCodec codec.KeyCodec[ReferenceKey],
	pkCodec codec.KeyCodec[PrimaryKey],
	sa collections.StoreAccessor,
	getRefKeyFunc func(pk PrimaryKey, value Value) (ReferenceKey, error),
) *Multi[ReferenceKey, PrimaryKey, Value] {
	return &Multi[ReferenceKey, PrimaryKey, Value]{
		getRefKey: getRefKeyFunc,
		refKeys:   collections.NewKeySet(storeKey, keyPrefix, sdkcollections.PairKeyCodec(refCodec, pkCodec), sa),
	}
}

func (multi *Multi[ReferenceKey, PrimaryKey, Value]) Reference(pk PrimaryKey, newValue Value, lazyOldValue func() (Value, error)) error {
	oldValue, err := lazyOldValue()
	switch {
	// if no error it means the value existed, and we need to remove the old indexes
	case err == nil:
		err = multi.unreference(pk, oldValue)
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
	// create new indexes
	refKey, err := multi.getRefKey(pk, newValue)
	if err != nil {
		return err
	}
	return multi.refKeys.Set(sdkcollections.Join(refKey, pk))
}

func (multi *Multi[ReferenceKey, PrimaryKey, Value]) Unreference(pk PrimaryKey, getValue func() (Value, error)) error {
	value, err := getValue()
	if err != nil {
		return err
	}
	return multi.unreference(pk, value)
}

func (multi *Multi[ReferenceKey, PrimaryKey, Value]) unreference(pk PrimaryKey, value Value) error {
	refKey, err := multi.getRefKey(pk, value)
	if err != nil {
		return err
	}
	return multi.refKeys.Remove(sdkcollections.Join(refKey, pk))
}

func (multi *Multi[ReferenceKey, PrimaryKey, Value]) Iterate() (MultiIterator[ReferenceKey, PrimaryKey], error) {
	iter, err := multi.refKeys.Iterate()
	return (MultiIterator[ReferenceKey, PrimaryKey])(iter), err
}

// MatchExact returns a MultiIterator containing all the primary keys referenced by the provided reference key.
func (multi *Multi[ReferenceKey, PrimaryKey, Value]) MatchExact(refKey ReferenceKey) (MultiIterator[ReferenceKey, PrimaryKey], error) {
	return multi.Iterate()
}

// MultiIterator is just a KeySetIterator with key as Pair[ReferenceKey, PrimaryKey].
type MultiIterator[ReferenceKey, PrimaryKey any] sdkcollections.KeySetIterator[sdkcollections.Pair[ReferenceKey, PrimaryKey]]

// PrimaryKey returns the iterator's current primary key.
func (i MultiIterator[ReferenceKey, PrimaryKey]) PrimaryKey() (PrimaryKey, error) {
	fullKey, err := i.FullKey()
	return fullKey.K2(), err
}

// PrimaryKeys fully consumes the iterator and returns the list of primary keys.
func (i MultiIterator[ReferenceKey, PrimaryKey]) PrimaryKeys() ([]PrimaryKey, error) {
	fullKeys, err := i.FullKeys()
	if err != nil {
		return nil, err
	}
	pks := make([]PrimaryKey, len(fullKeys))
	for i, fullKey := range fullKeys {
		pks[i] = fullKey.K2()
	}
	return pks, nil
}

// FullKey returns the current full reference key as Pair[ReferenceKey, PrimaryKey].
func (i MultiIterator[ReferenceKey, PrimaryKey]) FullKey() (sdkcollections.Pair[ReferenceKey, PrimaryKey], error) {
	return (sdkcollections.KeySetIterator[sdkcollections.Pair[ReferenceKey, PrimaryKey]])(i).Key()
}

// FullKeys fully consumes the iterator and returns all the list of full reference keys.
func (i MultiIterator[ReferenceKey, PrimaryKey]) FullKeys() ([]sdkcollections.Pair[ReferenceKey, PrimaryKey], error) {
	return (sdkcollections.KeySetIterator[sdkcollections.Pair[ReferenceKey, PrimaryKey]])(i).Keys()
}

// Next advances the iterator.
func (i MultiIterator[ReferenceKey, PrimaryKey]) Next() {
	(sdkcollections.KeySetIterator[sdkcollections.Pair[ReferenceKey, PrimaryKey]])(i).Next()
}

// Valid asserts if the iterator is still valid or not.
func (i MultiIterator[ReferenceKey, PrimaryKey]) Valid() bool {
	return (sdkcollections.KeySetIterator[sdkcollections.Pair[ReferenceKey, PrimaryKey]])(i).Valid()
}

// Close closes the iterator.
func (i MultiIterator[ReferenceKey, PrimaryKey]) Close() error {
	return (sdkcollections.KeySetIterator[sdkcollections.Pair[ReferenceKey, PrimaryKey]])(i).Close()
}
