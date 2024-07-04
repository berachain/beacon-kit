package collections

import (
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
	"cosmossdk.io/core/store"
)

type Map[K, V any] struct {
	storeKey      []byte
	keyPrefix     []byte
	keyCodec      codec.KeyCodec[K]
	valueCodec    codec.ValueCodec[V]
	storeAccessor StoreAccessor
}

func NewMap[K, V any](
	storeKey []byte,
	keyPrefix []byte,
	keyCodec codec.KeyCodec[K],
	valueCodec codec.ValueCodec[V],
	storeAccessor StoreAccessor,
) Map[K, V] {
	return Map[K, V]{
		storeKey:      storeKey,
		keyPrefix:     keyPrefix,
		keyCodec:      keyCodec,
		valueCodec:    valueCodec,
		storeAccessor: storeAccessor,
	}
}

func (m *Map[K, V]) Get(key K) (V, error) {
	var result V
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.keyCodec, key,
	)
	if err != nil {
		return result, err
	}
	res, err := query(m.storeAccessor(), m.storeKey, prefixedKey)
	if err != nil {
		return result, err
	}

	return m.valueCodec.Decode(res)
}

func (m *Map[K, V]) Set(key K, value V) error {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.keyCodec, key,
	)
	if err != nil {
		return err
	}
	store := m.storeAccessor()
	encodedValue, err := m.valueCodec.Encode(value)
	if err != nil {
		return err
	}
	store.AddChange(m.storeKey, prefixedKey, encodedValue)
	return nil
}

// Iterate provides an Iterator over K and V. It accepts a Ranger interface.
// A nil ranger equals to iterate over all the keys in ascending order.
func (m Map[K, V]) Iterate() (sdkcollections.Iterator[K, V], error) {
	var (
		iter   store.Iterator
		reader store.Reader
	)
	_, readerMap, err := m.storeAccessor().StateLatest()
	if err != nil {
		return sdkcollections.Iterator[K, V]{}, err
	}
	reader, err = readerMap.GetReader(m.storeKey)

	iter, err = reader.Iterator(m.keyPrefix, sdkcollections.NextBytesPrefixKey(m.keyPrefix))
	if err != nil {
		return sdkcollections.Iterator[K, V]{}, err
	}

	return sdkcollections.Iterator[K, V]{
		KeyCodec:     m.keyCodec,
		ValueCodec:   m.valueCodec,
		Iter:         iter,
		PrefixLength: len(m.keyPrefix),
	}, nil
}
