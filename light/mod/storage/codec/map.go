package codec

import (
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

type Map[K, V any] struct {
	prefix []byte

	key   codec.KeyCodec[K]
	value codec.ValueCodec[V]
}

// NewMap creates a new map with the given prefix, key codec, and value codec.
func NewMap[K, V any](
	prefix string, key codec.KeyCodec[K], value codec.ValueCodec[V],
) Map[K, V] {
	return Map[K, V]{
		prefix: []byte(prefix),
		key:    key,
		value:  value,
	}
}

// Key encodes the key with the prefix and key codec.
func (m Map[K, V]) Key(key K) ([]byte, error) {
	bytesKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.prefix, m.key, key,
	)
	if err != nil {
		return nil, err
	}

	return bytesKey, nil
}

// DecodeKey decodes the key with the prefix and key codec.
func (m Map[K, V]) Decode(value []byte) (V, error) {
	var out V
	var err error

	if value == nil {
		return out, nil
	}

	out, err = m.value.Decode(value)
	if err != nil {
		return out, err
	}

	return out, nil
}
