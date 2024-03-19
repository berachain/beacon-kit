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

func (v Uint8) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	binary.LittleEndian.PutUint16(bz, uint16(v))
	return [32]byte(bz), nil
}

type Uint16 uint16

func (v Uint16) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	binary.LittleEndian.PutUint16(bz, uint16(v))
	return [32]byte(bz), nil
}

type Uint32 uint32

func (v Uint32) HashTreeRoot() ([32]byte, error) {
	bz := make([]byte, 32)
	binary.LittleEndian.PutUint32(bz, uint32(v))
	return [32]byte(bz), nil
}

type Uint64 uint64

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

type MockUint64Container struct {
	Uint64Field Uint64
}

// We can have a generator for this.
func (c *MockUint64Container) Fields() []ssz.Hashable {
	return []ssz.Hashable{
		c.Uint64Field,
	}
}
