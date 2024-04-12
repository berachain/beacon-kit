package codec

import (
	"cosmossdk.io/collections/codec"
)

type Item[V any] struct {
	prefix []byte

	value codec.ValueCodec[V]
}

func NewItem[V any](prefix string, value codec.ValueCodec[V]) Item[V] {
	return Item[V]{
		prefix: []byte(prefix),
		value:  value,
	}
}

func (i Item[V]) Key() []byte {
	return i.prefix
}

func (i Item[V]) Decode(value []byte) (V, error) {
	var out V
	var err error

	if value == nil {
		return out, nil
	}

	out, err = i.value.Decode(value)
	if err != nil {
		return out, err
	}

	return out, nil
}
