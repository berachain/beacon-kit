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

package uint256

import (
	"bytes"
	"math/big"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

// UInt256Bytes is the number of bytes in a uint256.
const UInt256Bytes = 32

// LittleEndian represents a uint256 number. It
// is designed to marshal and unmarshal JSON in big-endian
// format, while under the hood storing the value as little-endian
// for compatibility with the SSZ spec.
type LittleEndian [32]byte

// NewLittleEndian creates a new LittleEndian from a byte slice.
func NewLittleEndian(bz []byte) LittleEndian {
	return LittleEndian(byteslib.ExtendToSize(bz, UInt256Bytes))
}

// LittleFromBigEndian creates a new LittleEndian from a big-endian
// byte slice.
func LittleFromBigEndian(b []byte) LittleEndian {
	return LittleEndian(
		byteslib.ExtendToSize(
			byteslib.CopyAndReverseEndianess(b),
			UInt256Bytes,
		),
	)
}

// LittleFromBigInt creates a new LittleEndian from a big.Int.
func LittleFromBigInt(b *big.Int) LittleEndian {
	if b == nil {
		return LittleEndian{}
	}
	return LittleFromBigEndian(b.Bytes())
}

// UInt256 converts an LittleEndian to a uint256.Int.
func (s LittleEndian) ToUInt256() *uint256.Int {
	return new(uint256.Int).SetBytes(byteslib.CopyAndReverseEndianess(s[:]))
}

// Big converts an LittleEndian to a big.Int.
func (s LittleEndian) ToBig() *big.Int {
	return new(big.Int).SetBytes(byteslib.CopyAndReverseEndianess(s[:]))
}

// MarshalJSON marshals a LittleEndian to JSON, it flips the endianness
// before encoding it to hex such that it is marshalled as big-endian.
func (s LittleEndian) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hexutil.EncodeBig(s.ToBig()) + "\""), nil
}

// UnmarshalJSON unmarshals a LittleEndian from JSON by decoding the hex
// string and flipping the endianness, such that it is unmarshalled as
// big-endian.
func (s *LittleEndian) UnmarshalJSON(input []byte) error {
	baseFee, err := hexutil.DecodeBig(string(bytes.Trim(input, "\"")))
	if err != nil {
		return err
	}
	*s = LittleEndian(
		byteslib.ExtendToSize(
			byteslib.CopyAndReverseEndianess(
				baseFee.Bytes()), UInt256Bytes),
	)
	return nil
}

// String returns the string representation of a LittleEndian.
func (s *LittleEndian) String() string {
	return s.ToUInt256().String()
}
