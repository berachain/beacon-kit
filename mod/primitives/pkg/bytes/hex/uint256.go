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

package hex

import (
	"encoding/json"
	"reflect"

	"github.com/holiman/uint256"
)

// U256 marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type U256 uint256.Int

// MarshalText implements encoding.TextMarshaler.
func (b U256) MarshalText() ([]byte, error) {
	u256 := (*uint256.Int)(&b)
	return []byte(u256.Hex()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *U256) UnmarshalJSON(input []byte) error {
	// The uint256.Int.UnmarshalJSON method accepts "dec", "0xhex"; we must be
	// more strict, hence we check string and invoke SetFromHex directly.
	if !isQuotedString(input) {
		return ErrNonQuotedString
	}
	// The hex decoder needs to accept empty string ("") as '0', which uint256.Int
	// would reject.
	if len(input) == prefixLen {
		(*uint256.Int)(b).Clear()
		return nil
	}
	err := (*uint256.Int)(b).SetFromHex(string(input[1 : len(input)-1]))
	if err != nil {
		u256T := reflect.TypeOf((*uint256.Int)(nil))
		return &json.UnmarshalTypeError{Value: err.Error(), Type: u256T}
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *U256) UnmarshalText(input []byte) error {
	// The uint256.Int.UnmarshalText method accepts "dec", "0xhex"; we must be
	// more strict, hence we check string and invoke SetFromHex directly.
	return (*uint256.Int)(b).SetFromHex(string(input))
}

// XString returns the hex encoding of b.
func (b *U256) String() string {
	return (*uint256.Int)(b).Hex()
}
