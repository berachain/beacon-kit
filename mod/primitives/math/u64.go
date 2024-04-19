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

package math

import (
	"encoding/binary"
	"math/bits"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	// U64NumBytes is the number of bytes in a U64.
	U64NumBytes = 8
	// U64NumBits is the number of bits in a U64.
	U64NumBits = U64NumBytes * 8
)

// U64 represents a 64-bit unsigned integer that is both SSZ and JSON
// marshallable. We marshal U64 as hex strings in JSON in order to keep the
// execution client apis happy, and we marshal U64 as little-endian in SSZ to be
// compatible with the spec.
type U64 uint64

//nolint:lll // links.
type (
	// Gwei is a denomination of 1e9 Wei represented as a U64.
	Gwei = U64

	// Slot as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Slot = U64

	// CommitteeIndex as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	CommitteeIndex = U64

	// ValidatorIndex as per the Ethereum 2.0  Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	ValidatorIndex = U64

	// Epoch as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Epoch = U64
)

// -------------------------- SSZMarshallable --------------------------

// MarshalSSZTo serializes the U64 into a byte slice.
func (u U64) MarshalSSZTo(buf []byte) ([]byte, error) {
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// MarshalSSZ serializes the U64 into a byte slice.
func (u U64) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, U64NumBytes)
	if _, err := u.MarshalSSZTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// UnmarshalSSZ deserializes the U64 from a byte slice.
func (u *U64) UnmarshalSSZ(buf []byte) error {
	if len(buf) != U64NumBytes {
		return ErrInvalidSSZLength
	}
	if u == nil {
		u = new(U64)
	}
	*u = U64(binary.LittleEndian.Uint64(buf))
	return nil
}

// SizeSSZ returns the size of the U64 in bytes.
func (u U64) SizeSSZ() int {
	return U64NumBytes
}

// HashTreeRoot computes the Merkle root of the U64 using SSZ hashing rules.
func (u U64) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, U64NumBytes)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	var hashRoot [32]byte
	copy(hashRoot[:], buf)
	return hashRoot, nil
}

// -------------------------- JSONMarshallable -------------------------

// UnmarshalJSON parses a blob in hex syntax.
func (u *U64) UnmarshalJSON(input []byte) error {
	return (*hexutil.Uint64)(u).UnmarshalJSON(input)
}

// MarshalText returns the hex representation of b.
func (u U64) MarshalText() ([]byte, error) {
	return hexutil.Uint64(u).MarshalText()
}

// ---------------------------- U64 Methods ----------------------------

// Unwrap returns the underlying uint64 value of U64.
func (u U64) Unwrap() uint64 {
	return uint64(u)
}

// NextPowerOfTwo returns the next power of two greater than or equal to the.
//
//nolint:mnd // powers of 2.
func (u U64) NextPowerOfTwo() U64 {
	u--
	u |= u >> 1
	u |= u >> 2
	u |= u >> 4
	u |= u >> 8
	u |= u >> 16
	u++
	return u
}

// ILog2Ceil returns the ceiling of the base 2 logarithm of the U64.
func (u U64) ILog2Ceil() uint8 {
	// Log2(0) is undefined, should we panic?
	if u == 0 {
		return 0
	}
	//#nosec:G701 // we handle the case of u == 0 above, so this is safe.
	return U64NumBits - uint8(bits.LeadingZeros64(uint64(u-1)))
}

// type U64Vector[U ~uint64] []U

// func (v U64Vector[U]) HashTreeRoot() ([32]byte, error) {
// 	return ssz.MerkleizeVecBasic[U64, U, [32]byte](v)
// }

// func (v U64Vector[U]) SizeSSZ() int {
// 	return int(ssz.SizeOfComposite[[32]byte, U64Vector[U]](v))
// }

// type U64List []U64

// func (v U64List) HashTreeRoot() ([32]byte, error) {
// 	return ssz.MerkleizeListBasic[U64, U64, [32]byte](v, 16)
// }

// func (v U64List) SizeSSZ() int {
// 	return int(ssz.SizeOfComposite[[32]byte, U64List](v))
// }

// type U64Container struct {
// 	Field2 U64List
// 	Field1 U64
// }

// func (c U64Container) SizeSSZ() int {
// 	return c.Field1.SizeSSZ() + c.Field2.SizeSSZ()
// }

// func (c U64Container) HashTreeRoot() ([32]byte, error) {
// 	return ssz.MerkleizeContainer[U64, U64Container, [32]byte](c)
// }

// //go:generate go run github.com
// /ferranbt/fastssz/sszgen -objs U64List2,U64Container2 --path ./u64
// .go -output bet.ssz.go
// type U64List2 struct {
// 	Data []uint64 `ssz-max:"16"`
// }

// type U64Container2 struct {
// 	Field2 []uint64 `ssz-max:"16"`
// 	Field1 U64
// }
