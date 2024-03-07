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

package primitives

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	byteslib "github.com/itsdevbear/bolaris/lib/bytes"
)

const thirtyTwo = 32

// SSZUInt256 represents a ssz-able uint64.
type SSZUInt256 []byte

// UInt256 converts an SSZUInt256 to a uint256.Int.
func (s *SSZUInt256) UInt256() *uint256.Int {
	return new(uint256.Int).SetBytes([]byte(*s))
}

// Big converts an SSZUInt256 to a big.Int.
func (s *SSZUInt256) Big() *big.Int {
	return new(big.Int).SetBytes([]byte(*s))
}

// UnmarshalJSON unmarshals a SSZUInt256 from JSON by decoding the hex string
// and flipping the endianness.
func (s *SSZUInt256) UnmarshalJSON(input []byte) error {
	input = bytes.Trim(input, "\"")
	baseFee, err := hexutil.DecodeBig(string(input))
	if err != nil {
		return err
	}
	*s = SSZUInt256(byteslib.ExtendToSize(
		byteslib.CopyAndReverseEndianess(baseFee.Bytes()),
		thirtyTwo,
	))
	return nil
}

// MarshalJSON marshals a SSZUInt256 to JSON, it flips the endianness
// before encoding it to hex.
func (s SSZUInt256) MarshalJSON() ([]byte, error) {
	baseFee := new(big.Int).
		SetBytes(byteslib.CopyAndReverseEndianess(s))
	return []byte("\"" + hexutil.EncodeBig(baseFee) + "\""), nil
}
