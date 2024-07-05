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

package math

import (
	"math/big"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/holiman/uint256"
)

// U256NumBytes is the number of bytes in a uint256.
const U256NumBytes = 32

// Wei is the smallest unit of Ether, we store the value as LittleEndian for
// the best compatibility with the SSZ spec.
type Wei = U256L

// U256L represents a uint256 number stored as little-endian. It
// is designed to marshal and unmarshal JSON in big-endian
// format, while under the hood storing the value as little-endian
// for compatibility with the SSZ spec.
type U256L [32]byte

// --------------------------- Constructors ----------------------------

// NewU256L creates a new U256L from a byte slice.
func NewU256L(bz []byte) (U256L, error) {
	// Ensure that we are not silently truncating the input.
	if len(bz) > U256NumBytes {
		return U256L{}, ErrUnexpectedInputLength(U256NumBytes, len(bz))
	}
	return U256L(byteslib.ExtendToSize(bz, U256NumBytes)), nil
}

// MustNewU256L creates a new U256L from a byte slice.
// Panics if the input is invalid.
func MustNewU256L(bz []byte) U256L {
	n, err := NewU256L(bz)
	if err != nil {
		panic(err)
	}
	return n
}

// NewU256LFromBigEndian creates a new U256L from a big-endian
// byte slice.
func NewU256LFromBigEndian(b []byte) (U256L, error) {
	return NewU256L(byteslib.CopyAndReverseEndianess(b))
}

// MustNewU256LFromBigEndian creates a new U256L from a big-endian
// byte slice. Panics if the input is invalid.
func MustNewU256LFromBigEndian(b []byte) U256L {
	n, err := NewU256L(byteslib.CopyAndReverseEndianess(b))
	if err != nil {
		panic(err)
	}
	return n
}

// NewU256LFromBigInt creates a new U256L from a big.Int.
func NewU256LFromBigInt(b *big.Int) (U256L, error) {
	if b == nil {
		return U256L{}, ErrNilBigInt
	} else if b.Sign() < 0 {
		return U256L{}, ErrNegativeBigInt(b)
	}
	return NewU256LFromBigEndian(b.Bytes())
}

// MustNewU256LFromBigInt creates a new U256L from a big.Int.
// Panics if the input is invalid.
func MustNewU256LFromBigInt(b *big.Int) U256L {
	n, err := NewU256LFromBigInt(b)
	if err != nil {
		panic(err)
	}
	return n
}

// ------------------------------ Unwraps ------------------------------

// Unwrap converts an U256L to a raw [32]byte chunk.
func (s U256L) Unwrap() [32]byte {
	return s
}

// UnwrapU256 converts an U256L to a *U256.
func (s U256L) UnwrapU256() *U256 {
	return (*U256)(
		(uint256.NewInt(0).SetBytes(byteslib.CopyAndReverseEndianess(s[:]))),
	)
}

// UnwrapBig converts a U256 to a non-negative big.Int.
func (s U256L) UnwrapBig() *big.Int {
	// SetBytes treats byte slice as unsigned int in big-endian, so result
	// will always be non-negative.
	return new(big.Int).SetBytes(byteslib.CopyAndReverseEndianess(s[:]))
}

// -------------------------- JSONMarshallable -------------------------

// MarshalJSON marshals a U256L to JSON, it flips the endianness
// before encoding it to hex such that it is marshalled as big-endian.
func (s U256L) MarshalJSON() ([]byte, error) {
	return []byte(hex.FromBigInt(s.UnwrapBig()).AddQuotes().Unwrap()), nil
}

// UnmarshalJSON unmarshals a U256L from JSON by decoding the hex
// string and flipping the endianness, such that it is unmarshalled as
// big-endian.
func (s *U256L) UnmarshalJSON(input []byte) error {
	baseFee, err := hex.FromJSONString(input).ToBigInt()
	if err != nil {
		return err
	}
	*s = U256L(
		byteslib.ExtendToSize(
			byteslib.CopyAndReverseEndianess(
				baseFee.Bytes()), U256NumBytes),
	)
	return nil
}

// -------------------------- SSZMarshallable --------------------------

// MarshalSSZTo serializes the U64 into a byte slice.
func (s U256L) MarshalSSZTo(buf []byte) ([]byte, error) {
	copy(buf, s[:])
	return buf, nil
}

// MarshalSSZ serializes a U256L into a byte slice.
func (s U256L) MarshalSSZ() ([]byte, error) {
	return s[:], nil
}

// UnmarshalSSZ deserializes a U256L from a byte slice.
func (s *U256L) UnmarshalSSZ(buf []byte) error {
	if len(buf) != U256NumBytes {
		return ErrUnexpectedInputLength(U256NumBytes, len(buf))
	}
	copy(s[:], buf)
	return nil
}

// SizeSSZ returns the size of the U256L in bytes.
func (s U256L) SizeSSZ() int {
	return U256NumBytes
}

// String returns the string representation of a U256L.
func (s *U256L) String() string {
	return s.UnwrapU256().Unwrap().String()
}
