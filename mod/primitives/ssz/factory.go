package ssz

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/dgraph-io/ristretto"
)

var enableCache = false

// ToggleCache enables caching of ssz hash tree root. It is disabled by default.
func ToggleCache(val bool) {
	enableCache = val
}

// StructFactory exports an implementation of a interface
// containing helpers for marshaling/unmarshaling, and determining
// the hash tree root of struct values.
var StructFactory = newStructSSZ()
var basicFactory = newBasicSSZ()
var basicArrayFactory = newBasicArraySSZ()
var rootsArrayFactory = newRootsArraySSZ()
var compositeArrayFactory = newCompositeArraySSZ()
var basicSliceFactory = newBasicSliceSSZ()
var stringFactory = newStringSSZ()
var compositeSliceFactory = newCompositeSliceSSZ()

// SSZAble defines a type which can marshal/unmarshal and compute its
// hash tree root according to the Simple Serialize specification.
// See: https://github.com/ethereum/eth2.0-specs/blob/v0.8.2/specs/simple-serialize.md.
type SSZAble interface {
	// Root(val reflect.Value, typ reflect.Type, fieldName string, maxCapacity uint64) ([32]byte, error)
	// Marshal(val reflect.Value, typ reflect.Type, buf []byte, startOffset uint64) (uint64, error)
	// Unmarshal(val reflect.Value, typ reflect.Type, buf []byte, startOffset uint64) (uint64, error)
}

// SSZFactory recursively walks down a type and determines which SSZ-able
// core type it belongs to, and then returns and implementation of
// SSZ-able that contains marshal, unmarshal, and hash tree root related
// functions for use.
func SSZFactory(val reflect.Value, typ reflect.Type) (SSZAble, error) {
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
		return StructFactory, nil
	case kind == reflect.Ptr:
		return SSZFactory(val.Elem(), typ.Elem())
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
	lock      sync.Mutex
}

const BasicArraySizeCache = 100000

func newBasicArraySSZ() *basicArraySSZ {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: BasicArraySizeCache, // number of keys to track frequency of (1M).
		MaxCost:     1 << 22,             // maximum cost of cache (3MB).
		// 100,000 roots will take up approximately 3 MB in memory.
		BufferItems: 64, // number of keys per Get buffer.
	})
	return &basicArraySSZ{
		hashCache: cache,
	}
}

type basicSSZ struct {
	hashCache *ristretto.Cache
	lock      sync.Mutex
}

const BasicTypeCacheSize = 100000

func newBasicSSZ() *basicSSZ {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: BasicTypeCacheSize, // number of keys to track frequency of (100K).
		MaxCost:     1 << 23,            // maximum cost of cache (3MB).
		// 100,000 roots will take up approximately 3 MB in memory.
		BufferItems: 64, // number of keys per Get buffer.
	})
	return &basicSSZ{
		hashCache: cache,
	}
}

type rootsArraySSZ struct {
	hashCache    *ristretto.Cache
	lock         sync.Mutex
	cachedLeaves map[string][][]byte
	layers       map[string][][][]byte
}

const RootsArraySizeCache = 100000

func newRootsArraySSZ() *rootsArraySSZ {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: RootsArraySizeCache, // number of keys to track frequency of (100000).
		MaxCost:     1 << 23,             // maximum cost of cache (3MB).
		// 100,000 roots will take up approximately 3 MB in memory.
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
