package collections

import (
	"cosmossdk.io/collections/codec"
)

type Item[V any] struct {
	storeKey      []byte
	key           []byte
	valueCodec    codec.ValueCodec[V]
	storeAccessor StoreAccessor
}

func NewItem[V any](
	storeKey []byte,
	key []byte,
	valueCodec codec.ValueCodec[V],
	storeAccessor StoreAccessor,
) Item[V] {
	return Item[V]{
		storeKey:      storeKey,
		key:           key,
		valueCodec:    valueCodec,
		storeAccessor: storeAccessor,
	}
}

func (i *Item[V]) Get() (V, error) {
	var result V
	res, err := query(i.storeAccessor(), i.storeKey, i.key)
	if err != nil {
		return result, err
	}

	return i.valueCodec.Decode(res)
}

func (i *Item[V]) Set(value V) error {
	store := i.storeAccessor()
	encodedValue, err := i.valueCodec.Encode(value)
	if err != nil {
		return err
	}
	store.AddChange(i.storeKey, i.key, encodedValue)
	return nil
}
