package codec

import "cosmossdk.io/collections/codec"

type None struct{}

// Codec is a wrapper around a key and value codec.
type Codec[K, V any] struct {
	Key   codec.KeyCodec[K]
	Value codec.ValueCodec[V]
}
