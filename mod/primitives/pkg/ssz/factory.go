// SPDX-License-Identifier: MIT
//
// # Copyright (c) 2024 Berachain Foundation
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
package ssz

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/dgraph-io/ristretto"
)

// Serializable defines a type which can marshal/unmarshal and compute its
// hash tree root according to the Simple Serialize specification.
// See: https://github.com/ethereum/eth2.0-specs/blob/v0.8.2/specs/simple-serialize.md.
//
//nolint:lll
type Serializable interface {
	//nolint:lll
	// Root(val reflect.Value, typ reflect.Type, fieldName string, maxCapacity uint64) ([32]byte, error)
	//nolint:lll
	// Marshal(val reflect.Value, typ reflect.Type, buf []byte, startOffset uint64) (uint64, error)
	//nolint:lll
	// Unmarshal(val reflect.Value, typ reflect.Type, buf []byte, startOffset uint64) (uint64, error)
}

// SSZFactory recursively walks down a type and determines which SSZ-able
// core type it belongs to, and then returns and implementation of
// SSZ-able that contains marshal, unmarshal, and hash tree root related
// functions for use.
func Factory(val reflect.Value, typ reflect.Type) (Serializable, error) {
	// StructFactory exports an implementation of a interface
	// containing helpers for marshaling/unmarshaling, and determining
	// the hash tree root of struct values.
	var structFactory = newStructSSZ()
	var basicFactory = newBasicSSZ()
	var basicArrayFactory = newBasicArraySSZ()
	var rootsArrayFactory = newRootsArraySSZ()
	var compositeArrayFactory = newCompositeArraySSZ()
	var basicSliceFactory = newBasicSliceSSZ()
	var stringFactory = newStringSSZ()
	var compositeSliceFactory = newCompositeSliceSSZ()

	kind := typ.Kind()
	switch {
	case isBasicType(kind) || isBasicTypeArray(typ, typ.Kind()):
		return basicFactory, nil
	case kind == reflect.String:
		return stringFactory, nil
	case kind == reflect.Slice:
		switch {
		case isBasicType(typ.Elem().Kind()):
			return basicSliceFactory, nil
		case !isVariableSizeType(typ.Elem()):
			return basicSliceFactory, nil
		default:
			return compositeSliceFactory, nil
		}
	case kind == reflect.Array:
		switch {
		case isRootsArray(val, typ):
			return rootsArrayFactory, nil
		case isBasicTypeArray(typ.Elem(), typ.Elem().Kind()):
			return basicArrayFactory, nil
		case !isVariableSizeType(typ.Elem()):
			return basicArrayFactory, nil
		default:
			return compositeArrayFactory, nil
		}
	case kind == reflect.Struct:
		return structFactory, nil
	case kind == reflect.Ptr:
		return Factory(val.Elem(), typ.Elem())
	default:
		return nil, fmt.Errorf("unsupported kind: %v", kind)
	}
}

type structSSZ struct{}

func newStructSSZ() *structSSZ {
	return &structSSZ{}
}

type basicArraySSZ struct {
	hashCache *ristretto.Cache
	_         sync.Mutex
}

const BasicArraySizeCache = 100000

func newBasicArraySSZ() *basicArraySSZ {
	//#nosec:G703 // WIP. Error from cache can be ignored
	cache, _ := ristretto.NewCache(&ristretto.Config{
		// number of keys to track frequency of (1M).
		NumCounters: BasicArraySizeCache,
		//nolint:mnd // explained in comment or not a magic number
		MaxCost: 1 << 22, // maximum cost of cache (3MB).
		// 100,000 roots will take up approximately 3 MB in memory.
		//nolint:mnd // explained in comment or not a magic number
		BufferItems: 64, // number of keys per Get buffer.
	})
	return &basicArraySSZ{
		hashCache: cache,
	}
}

type basicSSZ struct {
	hashCache *ristretto.Cache
	_         sync.Mutex
}

const BasicTypeCacheSize = 100000

func newBasicSSZ() *basicSSZ {
	//#nosec:G703 // WIP. Error from cache can be ignored
	cache, _ := ristretto.NewCache(&ristretto.Config{
		// number of keys to track frequency of (100K).
		NumCounters: BasicTypeCacheSize,
		//nolint:mnd // explained in comment or not a magic number
		MaxCost: 1 << 23, // maximum cost of cache (3MB).
		// 100,000 roots will take up approximately 3 MB in memory.
		//nolint:mnd // explained in comment or not a magic number
		BufferItems: 64, // number of keys per Get buffer.
	})
	return &basicSSZ{
		hashCache: cache,
	}
}

type rootsArraySSZ struct {
	hashCache    *ristretto.Cache
	_            sync.Mutex
	cachedLeaves map[string][][]byte
	layers       map[string][][][]byte
}

const RootsArraySizeCache = 100000

func newRootsArraySSZ() *rootsArraySSZ {
	//#nosec:G703 // WIP. Error from cache can be ignored
	cache, _ := ristretto.NewCache(&ristretto.Config{
		// number of keys to track frequency of (100000).
		NumCounters: RootsArraySizeCache,
		//nolint:mnd // explained in comment or not a magic number
		MaxCost: 1 << 23, // maximum cost of cache (3MB).
		// 100,000 roots will take up approximately 3 MB in memory.
		//nolint:mnd // explained in comment or not a magic number
		BufferItems: 64, // number of keys per Get buffer.
	})
	return &rootsArraySSZ{
		hashCache:    cache,
		cachedLeaves: make(map[string][][]byte),
		layers:       make(map[string][][][]byte),
	}
}

type compositeArraySSZ struct{}

func newCompositeArraySSZ() *compositeArraySSZ {
	return &compositeArraySSZ{}
}

type basicSliceSSZ struct{}

func newBasicSliceSSZ() *basicSliceSSZ {
	return &basicSliceSSZ{}
}

type stringSSZ struct{}

func newStringSSZ() *stringSSZ {
	return &stringSSZ{}
}

type compositeSliceSSZ struct{}

func newCompositeSliceSSZ() *compositeSliceSSZ {
	return &compositeSliceSSZ{}
}
