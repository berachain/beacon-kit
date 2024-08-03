package schema

import "github.com/karalabe/ssz"

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
	panic("unimplemented")
}

// Has implements ssz.CodecI.
func (c *Codec) Has() *ssz.Hasher[*Codec] {
	panic("unimplemented")
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
