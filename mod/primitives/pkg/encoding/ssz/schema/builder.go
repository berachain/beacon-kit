package schema

import (
	"errors"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/karalabe/ssz"
	"github.com/prysmaticlabs/go-bitfield"
)

type commonBytesLengths interface {
	// fork | address | verkle-stem | hash | pubkey | committee | signature | bloom
	~[4]byte | ~[20]byte | ~[31]byte | ~[32]byte | ~[48]byte | ~[64]byte | ~[96]byte | ~[256]byte
}

type newableStaticObject[U any] interface {
	ssz.StaticObject
	*U
	DefineSchema(*Codec)
}

type newableDynamicObject[U any] interface {
	ssz.DynamicObject
	*U
	DefineSchema(*Codec)
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

type Codec struct {
	ssz   *ssz.Codec
	stack []*container
}

func (c *Codec) peek() *container {
	if len(c.stack) == 0 {
		panic("stack is empty")
	}
	return c.stack[len(c.stack)-1]
}

func (c *Codec) empty() bool {
	return len(c.stack) == 0
}

func DefineBool[T ~bool](codec *Codec, name string, v *T) {
	if codec.empty() {
		codec.peek().DefineField(name, Bool())
		return
	}
	ssz.DefineBool(codec.ssz, v)
}

func DefineUint8[T ~uint8](codec *Codec, n *T) {
}

func DefineUint16[T ~uint16](codec *Codec, n *T) {
}

func DefineUint32[T ~uint32](codec *Codec, n *T) {
}

func DefineUint64[T ~uint64](codec *Codec, name string, n *T) {
	if !codec.empty() {
		codec.peek().DefineField(name, U64())
		return
	}
	ssz.DefineUint64(codec.ssz, n)
}

func DefineUint256(codec *Codec, name string, n **uint256.Int) {
	if !codec.empty() {
		codec.peek().DefineField(name, U256())
		return
	}
	ssz.DefineUint256(codec.ssz, n)
}

func DefineUint256BigInt(c *Codec, n **big.Int) {
}

func DefineStaticBytes[T commonBytesLengths](
	codec *Codec,
	name string,
	blob *T,
) {
	if !codec.empty() {
		codec.peek().DefineField(name, DefineByteVector(uint64(len(*blob))))
		return
	}
	ssz.DefineStaticBytes(codec.ssz, blob)
}

func DefineCheckedStaticBytes(c *Codec, blob *[]byte, size uint64) {

}

func DefineDynamicBytesOffset(
	c *Codec,
	name string,
	blob *[]byte,
	maxSize uint64,
) {
	if !c.empty() {
		c.peek().DefineField(name, DefineByteList(maxSize))
		return
	}
	ssz.DefineDynamicBytesOffset(c.ssz, blob, maxSize)
}

func DefineDynamicBytesContent(c *Codec, blob *[]byte, maxSize uint64) {
}

func DefineStaticObject[T newableStaticObject[U], U any](
	codec *Codec,
	name string,
	obj *T,
) {
	if !codec.empty() {
		nc := newContainer()
		codec.peek().DefineField(name, nc)
		codec.stack = append(codec.stack, nc)
		o := *obj
		o.DefineSchema(codec)
		codec.stack = codec.stack[:len(codec.stack)-1]
		return
	}
	ssz.DefineStaticObject(codec.ssz, obj)
}

func DefineDynamicObjectOffset[T newableDynamicObject[U], U any](
	codec *Codec,
	name string,
	obj *T,
) {
	if !codec.empty() {
		nc := newContainer()
		codec.peek().DefineField(name, nc)
		codec.stack = append(codec.stack, nc)
		o := *obj
		o.DefineSchema(codec)
		codec.stack = codec.stack[:len(codec.stack)-1]
		return
	}
	ssz.DefineDynamicObjectOffset(codec.ssz, obj)
}

func DefineDynamicObjectContent[T newableDynamicObject[U], U any](
	c *Codec,
	obj *T,
) {
}

func DefineArrayOfBits[T commonBitsLengths](c *Codec, bits *T, size uint64) {
}

func DefineSliceOfBitsOffset(
	c *Codec,
	bits *bitfield.Bitlist,
	maxBits uint64,
) {
}

func DefineSliceOfBitsContent(
	c *Codec,
	bits *bitfield.Bitlist,
	maxBits uint64,
) {
}

func DefineArrayOfUint64s[T commonUint64sLengths](c *Codec, ns *T) {
}

func DefineSliceOfUint64sOffset[T ~uint64](
	c *Codec,
	name string,
	ns *[]T,
	maxItems uint64,
) {
	if !c.empty() {
		c.peek().DefineField(name, DefineList(U64(), maxItems))
		return
	}
	ssz.DefineSliceOfUint64sOffset(c.ssz, ns, maxItems)
}

func DefineSliceOfUint64sContent[T ~uint64](
	c *Codec,
	ns *[]T,
	maxItems uint64,
) {
}

func DefineArrayOfStaticBytes[T commonBytesArrayLengths[U], U commonBytesLengths](
	c *Codec,
	blobs *T,
) {
}

func DefineUnsafeArrayOfStaticBytes[T commonBytesLengths](
	c *Codec,
	blobs []T,
) {
}

func DefineCheckedArrayOfStaticBytes[T commonBytesLengths](
	c *Codec,
	blobs *[]T,
	size uint64,
) {
}

func DefineSliceOfStaticBytesOffset[T commonBytesLengths](
	c *Codec,
	name string,
	bytes *[]T,
	maxItems uint64,
) {
	if !c.empty() {
		var t T
		bytesVector := DefineByteVector(uint64(len(t)))
		c.peek().DefineField(name, DefineList(bytesVector, maxItems))
		return
	}
	ssz.DefineSliceOfStaticBytesOffset(c.ssz, bytes, maxItems)
}

func DefineSliceOfStaticBytesContent[T commonBytesLengths](
	c *Codec,
	blobs *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicBytesOffset(
	codec *Codec,
	name string,
	blobs *[][]byte,
	maxItems uint64,
	maxSize uint64,
) {
}

func DefineSliceOfDynamicBytesContent(
	c *Codec,
	blobs *[][]byte,
	maxItems uint64,
	maxSize uint64,
) {
}

func DefineSliceOfStaticObjectsOffset[
	T newableStaticObject[U],
	U any,
](
	codec *Codec,
	name string,
	objects *[]T,
	maxItems uint64,
) {
	if !codec.empty() {
		obj := T(new(U))
		c := newContainer()
		codec.stack = append(codec.stack, c)
		obj.DefineSchema(codec)
		codec.stack = codec.stack[:len(codec.stack)-1]
		codec.peek().DefineField(name, DefineList(c, maxItems))
		return
	}
	ssz.DefineSliceOfStaticObjectsOffset(codec.ssz, objects, maxItems)
}

func DefineSliceOfStaticObjectsContent[
	T newableStaticObject[U],
	U any,
](
	c *Codec,
	objects *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicObjectsOffset[T newableDynamicObject[U], U any](
	c *Codec,
	objects *[]T,
	maxItems uint64,
) {
}

func DefineSliceOfDynamicObjectsContent[T newableDynamicObject[U], U any](
	c *Codec,
	objects *[]T,
	maxItems uint64,
) {
}

func Build(obj interface{ DefineSchema(*Codec) }) (SSZType, error) {
	nc := newContainer()
	codec := &Codec{stack: []*container{nc}}
	obj.DefineSchema(codec)
	if len(codec.stack) != 1 {
		return nil, errors.New("stack is not empty")
	}
	return nc, nil
}
