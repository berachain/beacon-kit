// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package mocks

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/lib/ssz"
)

type Uint8 uint8

const (
	BitsPerByte = 8
)

func (v Uint8) SizeSSZ() int {
	return 8 / BitsPerByte
}

func (v Uint8) MarshalSSZ() []byte {
	return []byte{byte(v)}
}

func (v Uint8) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	bz[0] = byte(v)
	return [32]byte(bz), nil
}

type Uint16 uint16

func (v Uint16) SizeSSZ() int {
	return 16 / BitsPerByte
}

func (v Uint16) MarshalSSZ() []byte {
	bz := make([]byte, v.SizeSSZ())
	binary.LittleEndian.PutUint16(bz, uint16(v))
	return bz
}

func (v Uint16) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	binary.LittleEndian.PutUint16(bz, uint16(v))
	return [32]byte(bz), nil
}

type Uint32 uint32

func (v Uint32) SizeSSZ() int {
	return 32 / BitsPerByte
}

func (v Uint32) MarshalSSZ() []byte {
	bz := make([]byte, v.SizeSSZ())
	binary.LittleEndian.PutUint32(bz, uint32(v))
	return bz
}

func (v Uint32) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	binary.LittleEndian.PutUint32(bz, uint32(v))
	return [32]byte(bz), nil
}

type Uint64 uint64

func (v Uint64) SizeSSZ() int {
	return 64 / BitsPerByte
}

func (v Uint64) MarshalSSZ() []byte {
	bz := make([]byte, v.SizeSSZ())
	binary.LittleEndian.PutUint64(bz, uint64(v))
	return bz
}

func (v Uint64) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	binary.LittleEndian.PutUint64(bz, uint64(v))
	return [32]byte(bz), nil
}

type Byte byte

func (v Byte) HashTreeRoot() ([32]byte, error) {
	return (Uint8(v)).HashTreeRoot()
}

type Bool bool

func (v Bool) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	if v {
		bz[0] = 1
	}
	return [32]byte(bz), nil
}

type UintN interface {
	MarshalSSZ() []byte
}

type Vector[T UintN] []T

func (v Vector[T]) HashTreeRoot() ([32]byte, error) {
	length := len(v)
	bz := make([]byte, 0)
	for i := 0; i < length; i++ {
		bz = append(bz, v[i].MarshalSSZ()...)
	}
	return ssz.MerkleizeByteSliceSSZ(bz)
}

type MockSingleFieldContainer[T ssz.Hashable] struct {
	Field T
}

// We can have a generator for this.
func (c *MockSingleFieldContainer[T]) Fields() []ssz.Hashable {
	return []ssz.Hashable{
		c.Field,
	}
}
