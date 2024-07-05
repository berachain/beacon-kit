package collections

import (
	"bytes"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
	"cosmossdk.io/core/store"
)

type Map[K, V any] struct {
	storeKey      []byte
	keyPrefix     []byte
	KeyCodec      codec.KeyCodec[K]
	ValueCodec    codec.ValueCodec[V]
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
		KeyCodec:      keyCodec,
		ValueCodec:    valueCodec,
		storeAccessor: storeAccessor,
	}
}

func (m *Map[K, V]) Get(key K) (V, error) {
	var result V
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return result, err
	}
	res, err := query(m.storeAccessor(), m.storeKey, prefixedKey)
	if err != nil {
		return result, err
	}

	return m.ValueCodec.Decode(res)
}

func (m *Map[K, V]) Set(key K, value V) error {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return err
	}
	store := m.storeAccessor()
	encodedValue, err := m.ValueCodec.Encode(value)
	if err != nil {
		return err
	}
	store.AddChange(m.storeKey, prefixedKey, encodedValue)
	return nil
}

// Has reports whether the key is present in storage or not.
// Errors with ErrEncoding if key encoding fails.
func (m *Map[K, V]) Has(key K) (bool, error) {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return false, err
	}
	res, err := query(m.storeAccessor(), m.storeKey, prefixedKey)
	if err != nil {
		return false, err
	}
	return res == nil, nil
}

// Remove removes the key from the storage.
// Errors with ErrEncoding if key encoding fails.
// If the key does not exist then this is a no-op.
func (m *Map[K, V]) Remove(key K) error {
	prefixedKey, err := sdkcollections.EncodeKeyWithPrefix(
		m.keyPrefix, m.KeyCodec, key,
	)
	if err != nil {
		return err
	}
	store := m.storeAccessor()
	store.AddChange(m.storeKey, prefixedKey, nil)
	return nil
}

// Iterate provides an Iterator over K and V. It accepts a Ranger interface.
// A nil ranger equals to iterate over all the keys in ascending order.
func (m *Map[K, V]) Iterate() (sdkcollections.Iterator[K, V], error) {
	return m.IterateRaw(m.keyPrefix, nil)
}

func (m *Map[K, V]) IterateRaw(
	start, end []byte,
) (sdkcollections.Iterator[K, V], error) {
	prefixedStart := append(m.keyPrefix, start...)
	var prefixedEnd []byte
	if end == nil {
		prefixedEnd = sdkcollections.NextBytesPrefixKey(m.keyPrefix)
	} else {
		prefixedEnd = append(m.keyPrefix, end...)
	}

	if bytes.Compare(prefixedStart, prefixedEnd) == 1 {
		return sdkcollections.Iterator[K, V]{}, sdkcollections.ErrInvalidIterator
	}

	var (
		iter   store.Iterator
		reader store.Reader
	)
	_, readerMap, err := m.storeAccessor().StateLatest()
	if err != nil {
		return sdkcollections.Iterator[K, V]{}, err
	}
	reader, err = readerMap.GetReader(m.storeKey)

	iter, err = reader.Iterator(start, end)
	if err != nil {
		return sdkcollections.Iterator[K, V]{}, err
	}

	return sdkcollections.Iterator[K, V]{
		KeyCodec:     m.KeyCodec,
		ValueCodec:   m.ValueCodec,
		Iter:         iter,
		PrefixLength: len(m.keyPrefix),
	}, nil
}
