package schema

import (
	"math/big"

	"github.com/holiman/uint256"
	"github.com/karalabe/ssz"
	"github.com/prysmaticlabs/go-bitfield"
)

type commonBytesLengths interface {
	// fork | address | verkle-stem | hash | pubkey | committee | signature | bloom
	~[4]byte | ~[20]byte | ~[31]byte | ~[32]byte | ~[48]byte | ~[64]byte | ~[96]byte | ~[256]byte
}

type newableStaticObject[C ssz.CodecI[C], U any] interface {
	ssz.StaticObject[C]
	*U
}

type newableDynamicObject[C ssz.CodecI[C], U any] interface {
	ssz.DynamicObject[C]
	*U
}

type commonBitsLengths interface {
	// justification
	~[1]byte
}

type commonUint64sLengths interface {
	// slashing
	~[8192]uint64
}

type commonBytesArrayLengths[U commonBytesLengths] interface {
	// verkle IPA vectors | proof | committee | history | randao
	~[8]U | ~[33]U | ~[512]U | ~[8192]U | ~[65536]U
}

type Builder struct {
	stack []container
}

func (c *Builder) peek() container {
	if len(c.stack) == 0 {
		panic("stack is empty")
	}
	return c.stack[len(c.stack)-1]
}

func (c *Builder) empty() bool {
	return len(c.stack) == 0
}

func DefineBool[T ~bool](codec *Codec, name string, v *T) {
	if codec.schema != nil {
		codec.schema.peek().DefineField(name, Bool())
		return
	}
	ssz.DefineBool(codec, v)
}

func DefineUint8[T ~uint8, C ssz.CodecI[C]](codec *Builder, n *T) {
}

func DefineUint16[T ~uint16, C ssz.CodecI[C]](codec *Builder, n *T) {
}

func DefineUint32[T ~uint32, C ssz.CodecI[C]](codec *Builder, n *T) {
}

func DefineUint64[T ~uint64, C CodecI[C]](codec C, name string, n *T) {
	if codec.Builder() != nil {
		codec.Builder().peek().DefineField(name, U64())
		return
	}
	ssz.DefineUint64(codec, n)
}

func DefineUint256(codec *Codec, name string, n **uint256.Int) {
	if codec.schema != nil {
		codec.schema.peek().DefineField(name, U256())
		return
	}
	ssz.DefineUint256(codec, n)
}

func DefineUint256BigInt(c *Builder, n **big.Int) {
}

func DefineStaticBytes[T commonBytesLengths, C CodecI[C]](
	codec C,
	name string,
	blob *T,
) {
	if codec.Builder() != nil {
		codec.Builder().
			peek().
			DefineField(name, DefineByteVector(uint64(len(*blob))))
		return
	}
	ssz.DefineStaticBytes(codec, blob)
}

func DefineCheckedStaticBytes(c *Builder, blob *[]byte, size uint64) {
}

func DefineDynamicBytesOffset(
	c *Codec,
	name string,
	blob *[]byte,
	maxSize uint64,
) {
	if c.schema != nil {
		c.schema.peek().DefineField(name, DefineByteList(maxSize))
		return
	}
	ssz.DefineDynamicBytesOffset(c, blob, maxSize)
}

func DefineDynamicBytesContent(c *Builder, blob *[]byte, maxSize uint64) {
}

func DefineStaticObject[T newableStaticObject[C, U], U any, C CodecI[C]](
	codec C,
	name string,
	obj *T,
) {
	if codec.Builder() != nil {
		b := codec.Builder()
		nc := container{}
		b.peek().DefineField(name, nc)
		b.stack = append(b.stack, nc)
		o := *obj
		o.DefineSSZ(codec)
		b.stack = b.stack[:len(b.stack)-1]
		return
	}
	ssz.DefineStaticObject(codec, obj)
}

func DefineDynamicObjectOffset[T newableDynamicObject[C, U], U any, C CodecI[C]](
	codec C,
	name string,
	obj *T,
) {
	if codec.Builder() != nil {
		b := codec.Builder()
		nc := container{}
		b.peek().DefineField(name, nc)
		b.stack = append(b.stack, nc)
		o := *obj
		o.DefineSSZ(codec)
		b.stack = b.stack[:len(b.stack)-1]
		return
	}
	ssz.DefineDynamicObjectOffset(codec, obj)
}

func DefineDynamicObjectContent[T newableDynamicObject[C, U], U any, C CodecI[C]](
	c *Builder,
	obj *T,
) {
}

func DefineArrayOfBits[T commonBitsLengths](c *Builder, bits *T, size uint64) {
}

func DefineSliceOfBitsOffset(
	c *Builder,
	bits *bitfield.Bitlist,
	maxBits uint64,
) {
}

func DefineSliceOfBitsContent(
	c *Builder,
	bits *bitfield.Bitlist,
	maxBits uint64,
) {
}

func DefineArrayOfUint64s[T commonUint64sLengths](c *Builder, ns *T) {
}

func DefineSliceOfUint64sOffset[T ~uint64, C CodecI[C]](
	c C,
	name string,
	ns *[]T,
	maxItems uint64,
) {
	if c.Builder() != nil {
		c.Builder().peek().DefineField(name, DefineList(U64(), maxItems))
		return
	}
	ssz.DefineSliceOfUint64sOffset(c, ns, maxItems)
}

func DefineSliceOfUint64sContent[T ~uint64](
	c *Builder,
	ns *[]T,
	maxItems uint64,
) {
}

func DefineArrayOfStaticBytes[T commonBytesArrayLengths[U], U commonBytesLengths](
	c *Builder,
	blobs *T,
) {
}

func DefineUnsafeArrayOfStaticBytes[T commonBytesLengths](
	c *Builder,
	blobs []T,
) {
}

func DefineCheckedArrayOfStaticBytes[T commonBytesLengths](
	c *Builder,
	blobs *[]T,
	size uint64,
) {
}

func DefineSliceOfStaticBytesOffset[T commonBytesLengths, C CodecI[C]](
	c C,
	name string,
	bytes *[]T,
	maxItems uint64,
) {
	if c.Builder() != nil {
		var t T
		bytesVector := DefineByteVector(uint64(len(t)))
		c.Builder().peek().DefineField(name, DefineList(bytesVector, maxItems))
		return
	}
	ssz.DefineSliceOfStaticBytesOffset(c, bytes, maxItems)
}

func DefineSliceOfStaticBytesContent[T commonBytesLengths](
	c *Builder,
	blobs *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicBytesOffset[C CodecI[C]](
	codec C,
	name string,
	blobs *[][]byte,
	maxItems uint64,
	maxSize uint64,
) {
}

func DefineSliceOfDynamicBytesContent(
	c *Builder,
	blobs *[][]byte,
	maxItems uint64,
	maxSize uint64,
) {
}

func DefineSliceOfStaticObjectsOffset[
	T newableStaticObject[C, U],
	U any,
	C CodecI[C],
](
	codec C,
	name string,
	objects *[]T,
	maxItems uint64,
) {
	if codec.Builder() != nil {
		var t T
		c := container{}
		b := codec.Builder()
		b.stack = append(b.stack, c)
		t.DefineSSZ(codec)
		b.stack = b.stack[:len(b.stack)-1]
		b.peek().DefineField(name, DefineList(c, maxItems))
		return
	}
	ssz.DefineSliceOfStaticObjectsOffset(codec, objects, maxItems)
}

func DefineSliceOfStaticObjectsContent[
	T newableStaticObject[C, U],
	U any,
	C CodecI[C],
](
	c *Builder,
	objects *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicObjectsOffset[T newableDynamicObject[C, U], U any, C CodecI[C]](
	c *Builder,
	objects *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicObjectsContent[T newableDynamicObject[C, U], U any, C CodecI[C]](
	c *Builder,
	objects *[]T,
	maxItems uint64,
) {
}
