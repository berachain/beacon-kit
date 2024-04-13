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

// UInt256ByteLength is the length of a uint256 in bytes.
const UInt256ByteLength = 32

// LittleEndian represents a uint256 number. It
// is designed to marshal and unmarshal JSON in little-endian
// format, while under the hood storing the value as big-endian
// for compatibility with.
type LittleEndian []byte

// LittleFromBigEndian creates a new LittleEndian from a big-endian
// byte slice.
func LittleFromBigEndian(b []byte) LittleEndian {
	return LittleEndian(byteslib.CopyAndReverseEndianess(b))
}

// UInt256 converts an LittleEndian to a uint256.Int.
func (s LittleEndian) UInt256() *uint256.Int {
	return new(uint256.Int).SetBytes([]byte(s))
}

// Big converts an LittleEndian to a big.Int.
func (s LittleEndian) Big() *big.Int {
	return new(big.Int).SetBytes([]byte(s))
}

// MarshalJSON marshals a LittleEndian to JSON, it flips the endianness
// before encoding it to hex such that it is marshalled as big-endian.
func (s LittleEndian) MarshalJSON() ([]byte, error) {
	baseFee := new(big.Int).
		SetBytes(byteslib.CopyAndReverseEndianess(s))
	return []byte("\"" + hexutil.EncodeBig(baseFee) + "\""), nil
}

// UnmarshalJSON unmarshals a LittleEndian from JSON by decoding the hex
// string and flipping the endianness, such that it is unmarshalled as
// big-endian.
func (s *LittleEndian) UnmarshalJSON(input []byte) error {
	input = bytes.Trim(input, "\"")
	baseFee, err := hexutil.DecodeBig(string(input))
	if err != nil {
		return err
	}
	*s = LittleEndian(byteslib.ExtendToSize(
		byteslib.CopyAndReverseEndianess(baseFee.Bytes()),
		UInt256ByteLength,
	))
	return nil
}
