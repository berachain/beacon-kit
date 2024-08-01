package schema

import (
	"math/big"

	"github.com/holiman/uint256"
	"github.com/karalabe/ssz"
	"github.com/prysmaticlabs/go-bitfield"
)

type commonBytesLengths interface {
	// fork | address | verkle-stem | hash | pubkey | committee | signature | bloom | blob & tx (nice fuckup dev kek)
	~[4]byte | ~[20]byte | ~[31]byte | ~[32]byte | ~[48]byte | ~[64]byte | ~[96]byte | ~[256]byte | ~[131072]byte | ~[]byte
}

type newableStaticObject[U any] interface {
	ssz.StaticObject
	*U
	DefineSchema(*Builder)
}

type newableDynamicObject[U any] interface {
	ssz.DynamicObject
	*U
	DefineSchema(*Builder)
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
	sc    *ssz.Codec
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

func DefineBool[T ~bool](c *Builder, name string, v *T) {
	if !c.empty() {
		c.peek().DefineField(name, Bool())
		return
	}
	ssz.DefineBool(c.sc, v)
}

func DefineUint8[T ~uint8](c *Builder, n *T) {
}

func DefineUint16[T ~uint16](c *Builder, n *T) {
}

func DefineUint32[T ~uint32](c *Builder, n *T) {
}

func DefineUint64[T ~uint64](c *Builder, name string, n *T) {
	if !c.empty() {
		c.peek().DefineField(name, U64())
		return
	}
	ssz.DefineUint64(c.sc, n)
}

func DefineUint256(b *Builder, name string, n **uint256.Int) {
	if !b.empty() {
		b.peek().DefineField(name, U256())
		return
	}
	ssz.DefineUint256(b.sc, n)
}

func DefineUint256BigInt(c *Builder, n **big.Int) {
}

func DefineStaticBytes[T commonBytesLengths](c *Builder, name string, blob *T) {
}

func DefineCheckedStaticBytes(c *Builder, blob *[]byte, size uint64) {
}

func DefineDynamicBytesOffset(
	c *Builder,
	name string,
	blob *[]byte,
	maxSize uint64,
) {
	if !c.empty() {
		c.peek().DefineField(name, DefineByteList(maxSize))
		return
	}
	ssz.DefineDynamicBytesOffset(c.sc, blob, maxSize)
}

func DefineDynamicBytesContent(c *Builder, blob *[]byte, maxSize uint64) {
}

func DefineStaticObject[T newableStaticObject[U], U any](
	b *Builder,
	name string,
	obj *T,
) {
	if !b.empty() {
		nc := container{}
		b.peek().DefineField(name, nc)
		b.stack = append(b.stack, nc)
		o := *obj
		o.DefineSchema(b)
		b.stack = b.stack[:len(b.stack)-1]
		return
	}
	ssz.DefineStaticObject(b.sc, obj)
}

func DefineDynamicObjectOffset[T newableDynamicObject[U], U any](
	b *Builder,
	name string,
	obj *T,
) {
	if !b.empty() {
		nc := container{}
		b.peek().DefineField(name, nc)
		b.stack = append(b.stack, nc)
		o := *obj
		o.DefineSchema(b)
		b.stack = b.stack[:len(b.stack)-1]
		return
	}
	ssz.DefineDynamicObjectOffset(b.sc, obj)
}

func DefineDynamicObjectContent[T newableDynamicObject[U], U any](
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

func DefineSliceOfUint64sOffset[T ~uint64](
	b *Builder,
	name string,
	ns *[]T,
	maxItems uint64,
) {
	if !b.empty() {
		b.peek().DefineField(name, DefineList(U64(), maxItems))
		return
	}
	ssz.DefineSliceOfUint64sOffset(b.sc, ns, maxItems)
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

func DefineSliceOfStaticBytesOffset[T commonBytesLengths](
	b *Builder,
	name string,
	bytes *[]T,
	maxItems uint64,
) {
	if !b.empty() {
		var t T
		bytesVector := DefineByteVector(uint64(len(t)))
		b.peek().DefineField(name, DefineList(bytesVector, maxItems))
		return
	}
	ssz.DefineSliceOfStaticBytesOffset(b.sc, bytes, maxItems)
}

func DefineSliceOfStaticBytesContent[T commonBytesLengths](
	c *Builder,
	blobs *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicBytesOffset(
	c *Builder,
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

func DefineSliceOfStaticObjectsOffset[T newableStaticObject[U], U any](
	b *Builder,
	name string,
	objects *[]T,
	maxItems uint64,
) {
	if !b.empty() {
		var t T
		c := container{}
		b.stack = append(b.stack, c)
		t.DefineSchema(b)
		b.stack = b.stack[:len(b.stack)-1]
		b.peek().DefineField(name, DefineList(c, maxItems))
		return
	}
	ssz.DefineSliceOfStaticObjectsOffset(b.sc, objects, maxItems)
}

func DefineSliceOfStaticObjectsContent[T newableStaticObject[U], U any](
	c *Builder,
	objects *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicObjectsOffset[T newableDynamicObject[U], U any](
	c *Builder,
	objects *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicObjectsContent[T newableDynamicObject[U], U any](
	c *Builder,
	objects *[]T,
	maxItems uint64,
) {
}
