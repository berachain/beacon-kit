package schema

import (
	"fmt"

	"github.com/karalabe/ssz"
)

type Codec struct {
	enc    *ssz.Encoder[*Codec]
	dec    *ssz.Decoder[*Codec]
	has    *ssz.Hasher[*Codec]
	schema *Builder
}

type CodecI[SelfT any] interface {
	ssz.CodecI[SelfT]
	Builder() *Builder
}

var _ ssz.CodecI[*Codec] = (*Codec)(nil)

func (c *Codec) Builder() *Builder {
	return c.schema
}

// Dec implements ssz.CodecI.
func (c *Codec) Dec() *ssz.Decoder[*Codec] {
	return c.dec
}

// Enc implements ssz.CodecI.
func (c *Codec) Enc() *ssz.Encoder[*Codec] {
	return c.enc
}

// Has implements ssz.CodecI.
func (c *Codec) Has() *ssz.Hasher[*Codec] {
	return c.has
}

// DefineDecoder implements ssz.CodecI.
func (c *Codec) DefineDecoder(impl func(dec *ssz.Decoder[*Codec])) {
	if c.dec != nil {
		impl(c.dec)
	}
}

// DefineEncoder implements ssz.CodecI.
func (c *Codec) DefineEncoder(impl func(enc *ssz.Encoder[*Codec])) {
	if c.enc != nil {
		impl(c.enc)
	}
}

// DefineHasher implements ssz.CodecI.
func (c *Codec) DefineHasher(impl func(has *ssz.Hasher[*Codec])) {
	if c.has != nil {
		impl(c.has)
	}
}

func (c *Codec) DefineSchema(impl func(builder *Builder)) {
	if c.schema != nil {
		impl(c.schema)
	}
}

func Build(obj interface{ DefineSSZ(codec *Codec) }) (SSZType, error) {
	builder := new(Builder)
	builder.stack = append(builder.stack, newContainer())
	c := &Codec{schema: builder}
	obj.DefineSSZ(c)
	if len(builder.stack) != 1 {
		return nil, fmt.Errorf("unbalanced stack: %d", len(builder.stack))
	}
	return builder.stack[0], nil
}

func DecodeFromBytes(data []byte, obj ssz.Object[*Codec]) error {
	codec := &Codec{dec: new(ssz.Decoder[*Codec])}
	return ssz.DecodeFromBytesWithCodec(codec, data, obj)
}
