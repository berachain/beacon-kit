// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
	if !codec.empty() {
		codec.peek().DefineField(name, Bool())
		return
	}
	ssz.DefineBool(codec.ssz, v)
}

func DefineUint8[T ~uint8](codec *Codec, name string, v *T) {
	if !codec.empty() {
		codec.peek().DefineField(name, U8())
		return
	}
	ssz.DefineUint8(codec.ssz, v)
}

func DefineUint16[T ~uint16](codec *Codec, name string, v *T) {
	if !codec.empty() {
		codec.peek().DefineField(name, U16())
		return
	}
	ssz.DefineUint16(codec.ssz, v)
}

func DefineUint32[T ~uint32](codec *Codec, name string, v *T) {
	if !codec.empty() {
		codec.peek().DefineField(name, U32())
		return
	}
	ssz.DefineUint32(codec.ssz, v)
}

func DefineUint64[T ~uint64](codec *Codec, name string, v *T) {
	if !codec.empty() {
		codec.peek().DefineField(name, U64())
		return
	}
	ssz.DefineUint64(codec.ssz, v)
}

func DefineUint256(codec *Codec, name string, n **uint256.Int) {
	if !codec.empty() {
		codec.peek().DefineField(name, U256())
		return
	}
	ssz.DefineUint256(codec.ssz, n)
}

func DefineUint256BigInt(c *Codec, name string, n **big.Int) {
	if !c.empty() {
		c.peek().DefineField(name, U256())
		return
	}
	ssz.DefineUint256BigInt(c.ssz, n)
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

func DefineCheckedStaticBytes(
	c *Codec,
	name string,
	blob *[]byte,
	size uint64,
) {
	if !c.empty() {
		c.peek().DefineField(name, DefineByteVector(size))
		return
	}
	ssz.DefineCheckedStaticBytes(c.ssz, blob, size)
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

func DefineDynamicBytesContent(
	c *Codec,
	name string,
	blob *[]byte,
	maxSize uint64,
) {
	if !c.empty() {
		c.peek().DefineField(name, DefineByteList(maxSize))
		return
	}
	ssz.DefineDynamicBytesContent(c.ssz, blob, maxSize)
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
	_ string,
	obj *T,
) {
	ssz.DefineDynamicObjectContent(c.ssz, obj)
}

func DefineArrayOfBits[T commonBitsLengths](
	c *Codec,
	_ string,
	bits *T,
	size uint64,
) {
	ssz.DefineArrayOfBits(c.ssz, bits, size)
}

func DefineSliceOfBitsOffset(
	c *Codec,
	_ string,
	bits *bitfield.Bitlist,
	maxBits uint64,
) {
	ssz.DefineSliceOfBitsOffset(c.ssz, bits, maxBits)
}

func DefineSliceOfBitsContent(
	c *Codec,
	_ string,
	bits *bitfield.Bitlist,
	maxBits uint64,
) {
	ssz.DefineSliceOfBitsContent(c.ssz, bits, maxBits)
}

func DefineArrayOfUint64s[T commonUint64sLengths](
	c *Codec,
	name string,
	ns *T,
) {
	if !c.empty() {
		c.peek().DefineField(name, DefineVector(U64(), uint64(len(*ns))))
		return
	}
	ssz.DefineArrayOfUint64s(c.ssz, ns)
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
	_ string,
	ns *[]T,
	maxItems uint64,
) {
	ssz.DefineSliceOfUint64sContent(c.ssz, ns, maxItems)
}

func DefineArrayOfStaticBytes[
	T commonBytesArrayLengths[U],
	U commonBytesLengths,
](
	c *Codec,
	name string,
	blobs *T,
) {
	if !c.empty() {
		length := uint64(len(*blobs))
		bytesVector := DefineByteVector(length)
		c.peek().DefineField(name, DefineVector(bytesVector, length))
		return
	}
	ssz.DefineArrayOfStaticBytes[T, U](c.ssz, blobs)
}

func DefineUnsafeArrayOfStaticBytes[T commonBytesLengths](
	c *Codec,
	name string,
	blobs []T,
) {
	if !c.empty() {
		length := uint64(len(blobs))
		bytesVector := DefineByteVector(uint64(len(blobs[0])))
		c.peek().DefineField(name, DefineVector(bytesVector, length))
		return
	}
	ssz.DefineUnsafeArrayOfStaticBytes(c.ssz, blobs)
}

func DefineCheckedArrayOfStaticBytes[T commonBytesLengths](
	c *Codec,
	name string,
	blobs *[]T,
	size uint64,
) {
	if !c.empty() {
		bytesVector := DefineByteVector(size)
		c.peek().DefineField(name, DefineVector(bytesVector, size))
		return
	}
	ssz.DefineCheckedArrayOfStaticBytes(c.ssz, blobs, size)
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
	_ string,
	blobs *[]T,
	maxItems uint64,
) {
	ssz.DefineSliceOfStaticBytesContent(c.ssz, blobs, maxItems)
}

func DefineSliceOfDynamicBytesOffset(
	codec *Codec,
	name string,
	blobs *[][]byte,
	maxItems uint64,
	maxSize uint64,
) {
	if !codec.empty() {
		codec.peek().
			DefineField(name, DefineList(DefineByteList(maxSize), maxItems))
		return
	}
	ssz.DefineSliceOfDynamicBytesOffset(codec.ssz, blobs, maxItems, maxSize)
}

func DefineSliceOfDynamicBytesContent(
	c *Codec,
	_ string,
	blobs *[][]byte,
	maxItems uint64,
	maxSize uint64,
) {
	ssz.DefineSliceOfDynamicBytesContent(c.ssz, blobs, maxItems, maxSize)
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
	_ string,
	objects *[]T,
	maxItems uint64,
) {
	ssz.DefineSliceOfStaticObjectsContent(c.ssz, objects, maxItems)
}

func DefineSliceOfDynamicObjectsOffset[T newableDynamicObject[U], U any](
	c *Codec,
	name string,
	objects *[]T,
	maxItems uint64,
) {
	if !c.empty() {
		obj := T(new(U))
		nc := newContainer()
		c.stack = append(c.stack, nc)
		obj.DefineSchema(c)
		c.stack = c.stack[:len(c.stack)-1]
		c.peek().DefineField(name, DefineList(nc, maxItems))
		return
	}
	ssz.DefineSliceOfDynamicObjectsOffset(c.ssz, objects, maxItems)
}

func DefineSliceOfDynamicObjectsContent[T newableDynamicObject[U], U any](
	c *Codec,
	_ string,
	objects *[]T,
	maxItems uint64,
) {
	ssz.DefineSliceOfDynamicObjectsContent(c.ssz, objects, maxItems)
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
