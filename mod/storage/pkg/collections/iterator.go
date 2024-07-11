package collections

import (
	"fmt"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
	"cosmossdk.io/store"
)

type Iterator[K, V any] struct {
	storeIterator store.Iterator
	PrefixLength  int
	KeyPrefix     []byte
	StoreKey      []byte
	KeyCodec      codec.KeyCodec[K]
	ValueCodec    codec.ValueCodec[V]
}

// Keys returns all the keys in the iterator
func (i Iterator[K, V]) Keys() ([]K, error) {
	defer i.Close()
	var keys []K
	for ; i.storeIterator.Valid(); i.storeIterator.Next() {
		key, err := i.Key()
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (i Iterator[K, V]) Close() error {
	return i.storeIterator.Close()
}

// Key returns the current storetypes.Iterator decoded key.
func (i Iterator[K, V]) Key() (K, error) {
	// strip away key prefix
	bytesKey := i.storeIterator.Key()[i.PrefixLength:]

	read, key, err := i.KeyCodec.Decode(bytesKey)
	if err != nil {
		var k K
		return k, err
	}
	if read != len(bytesKey) {
		var k K
		return k, fmt.Errorf("%w: key decoder didn't fully consume the key: %T %x %d", sdkcollections.ErrEncoding, i.KeyCodec, bytesKey, read)
	}
	return key, nil
}
