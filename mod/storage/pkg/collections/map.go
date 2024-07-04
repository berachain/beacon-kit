package collections

import (
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
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
