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
	"errors"
	"fmt"
	"reflect"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

// Indexes represents a type which groups multiple Index
// of one Value saved with the provided PrimaryKey.
// Indexes is just meant to be a struct containing all
// the indexes to maintain relationship for.
type Indexes[PrimaryKey, Value any] interface {
	// IndexesList is implemented by the Indexes type
	// and returns all the grouped Index of Value.
	IndexesList() []Index[PrimaryKey, Value]
}

// Index represents an index of the Value indexed using the type PrimaryKey.
type Index[PrimaryKey, Value any] interface {
	// Reference creates a reference between the provided primary key and value.
	// It provides a lazyOldValue function that if called will attempt to fetch
	// the previous old value, returns ErrNotFound if no value existed.
	Reference(pk PrimaryKey, newValue Value, lazyOldValue func() (Value, error)) error
	// Unreference removes the reference between the primary key and value.
	// If error is ErrNotFound then it means that the value did not exist before.
	Unreference(pk PrimaryKey, lazyOldValue func() (Value, error)) error
}

// IndexedMap works like a Map but creates references between fields of Value and its PrimaryKey.
// These relationships are expressed and maintained using the Indexes type.
// Internally IndexedMap can be seen as a partitioned collection, one partition
// is a Map[PrimaryKey, Value], that maintains the object, the second
// are the Indexes.
type IndexedMap[PrimaryKey, Value, Idx any] struct {
	Indexes         Idx
	computedIndexes []Index[PrimaryKey, Value]
	primaryMap      Map[PrimaryKey, Value]
}

// NewIndexedMapSafe behaves like NewIndexedMap but returns errors.
func NewIndexedMapSafe[K, V, I any](
	storeKey []byte,
	keyPrefix []byte,
	pkCodec codec.KeyCodec[K],
	valueCodec codec.ValueCodec[V],
	indexes I,
	sa StoreAccessor,
) (im *IndexedMap[K, V, I], err error) {
	var indexesList []Index[K, V]
	indexesImpl, ok := any(indexes).(Indexes[K, V])
	if ok {
		indexesList = indexesImpl.IndexesList()
	} else {
		// if does not implement Indexes, then we try to infer using reflection
		indexesList, err = tryInferIndexes[I, K, V](indexes)
		if err != nil {
			return nil, fmt.Errorf("unable to infer indexes using reflection, consider implementing Indexes interface: %w", err)
		}
	}

	return &IndexedMap[K, V, I]{
		computedIndexes: indexesList,
		Indexes:         indexes,
		primaryMap:      NewMapKeeper(storeKey, keyPrefix, pkCodec, valueCodec, sa),
	}, nil
}

var (
	// testing sentinel errors
	errNotStruct = errors.New("wanted struct or pointer to a struct")
	errNotIndex  = errors.New("field is not an index implementation")
)

func tryInferIndexes[I, K, V any](indexes I) ([]Index[K, V], error) {
	typ := reflect.TypeOf(indexes)
	v := reflect.ValueOf(indexes)
	// check if struct or pointer to a struct
	if typ.Kind() != reflect.Struct && (typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct) {
		return nil, fmt.Errorf("%w: type %v", errNotStruct, typ)
	}
	// dereference
	if typ.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	indexesImpl := make([]Index[K, V], v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		index, ok := field.Interface().(Index[K, V])
		if !ok {
			return nil, fmt.Errorf("%w: field number %d", errNotIndex, i)
		}
		indexesImpl[i] = index
	}
	return indexesImpl, nil
}

// NewIndexedMap instantiates a new IndexedMap. Accepts a SchemaBuilder, a Prefix,
// a humanized name that defines the name of the collection, the primary key codec
// which is basically what IndexedMap uses to encode the primary key to bytes,
// the value codec which is what the IndexedMap uses to encode the value.
// Then it expects the initialized indexes. Reflection is used to infer the
// indexes, Indexes can optionally be implemented to be explicit. Panics
// on failure to create indexes. If you want an erroring API use NewIndexedMapSafe.
func NewIndexedMap[K, V, I any](
	storeKey []byte,
	keyPrefix []byte,
	pkCodec codec.KeyCodec[K],
	valueCodec codec.ValueCodec[V],
	indexes I,
	sa StoreAccessor,
) *IndexedMap[K, V, I] {
	im, err := NewIndexedMapSafe(storeKey, keyPrefix, pkCodec, valueCodec, indexes, sa)
	if err != nil {
		panic(err)
	}
	return im
}

// Get gets the object given its primary key.
func (im *IndexedMap[PrimaryKey, Value, Idx]) Get(pk PrimaryKey) (Value, error) {
	return im.primaryMap.Get(pk)
}

// Iterate allows to iterate over the objects given a Ranger of the primary key.
func (im *IndexedMap[PrimaryKey, Value, Idx]) Iterate() (sdkcollections.Iterator[PrimaryKey, Value], error) {
	return im.primaryMap.Iterate()
}

func (im *IndexedMap[PrimaryKey, Value, Idx]) NumKeys() (uint64, error) {
	numKeys, err := im.primaryMap.NumKeys()
	if err != nil {
		return 0, err
	}
	return uint64(numKeys), nil
}

// Has reports if exists a value with the provided primary key.
func (im *IndexedMap[PrimaryKey, Value, Idx]) Has(pk PrimaryKey) (bool, error) {
	return im.primaryMap.Has(pk)
}

// Set maps the value using the primary key. It will also iterate every index and instruct them to
// add or update the indexes.
func (im *IndexedMap[PrimaryKey, Value, Idx]) Set(pk PrimaryKey, value Value) error {
	err := im.ref(pk, value)
	if err != nil {
		return err
	}
	return im.primaryMap.Set(pk, value)
}

// Remove removes the value associated with the primary key from the map. Then
// it iterates over all the indexes and instructs them to remove all the references
// associated with the removed value.
func (im *IndexedMap[PrimaryKey, Value, Idx]) Remove(pk PrimaryKey) error {
	err := im.unref(pk)
	if err != nil {
		return err
	}
	return im.primaryMap.Remove(pk)
}

// // Walk applies the same semantics as Map.Walk.
// func (m *IndexedMap[PrimaryKey, Value, Idx]) Walk(ranger Ranger[PrimaryKey], walkFunc func(key PrimaryKey, value Value) (stop bool, err error)) error {
// 	return m.m.Walk(ranger, walkFunc)
// }

// IterateRaw iterates the IndexedMap using raw bytes keys. Follows the same semantics as Map.IterateRaw
func (im *IndexedMap[PrimaryKey, Value, Idx]) IterateRaw(start, end []byte) (sdkcollections.Iterator[PrimaryKey, Value], error) {
	return im.primaryMap.IterateRaw(start, end)
}

func (im *IndexedMap[PrimaryKey, Value, Idx]) KeyCodec() codec.KeyCodec[PrimaryKey] {
	return im.primaryMap.KeyCodec
}

func (im *IndexedMap[PrimaryKey, Value, Idx]) ValueCodec() codec.ValueCodec[Value] {
	return im.primaryMap.ValueCodec
}

func (im *IndexedMap[PrimaryKey, Value, Idx]) ref(pk PrimaryKey, value Value) error {
	for _, index := range im.computedIndexes {
		err := index.Reference(pk, value, cachedGet[PrimaryKey, Value](im, pk))
		if err != nil {
			return err
		}
	}
	return nil
}

func (im *IndexedMap[PrimaryKey, Value, Idx]) unref(pk PrimaryKey) error {
	for _, index := range im.computedIndexes {
		err := index.Unreference(pk, cachedGet[PrimaryKey, Value](im, pk))
		if err != nil {
			return err
		}
	}
	return nil
}

// cachedGet returns a function that gets the value V, given the key K but
// returns always the same result on multiple calls.
func cachedGet[K, V any, M interface {
	Get(key K) (V, error)
}](m M, key K,
) func() (V, error) {
	var (
		value      V
		err        error
		calledOnce bool
	)

	return func() (V, error) {
		if calledOnce {
			return value, err
		}
		value, err = m.Get(key)
		calledOnce = true
		return value, err
	}
}
