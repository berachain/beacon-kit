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
	"reflect"
	"strconv"
)

// U64 marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type U64 uint64

// MarshalText implements encoding.TextMarshaler.
func (b U64) MarshalText() ([]byte, error) {
	buf := make([]byte, prefixLen, initialCapacity)
	copy(buf, prefix)
	buf = strconv.AppendUint(buf, uint64(b), hexBase)
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *U64) UnmarshalJSON(input []byte) error {
	uint64T := reflect.TypeOf(U64(0))
	if !isQuotedString(input) {
		return wrapUnmarshalError(ErrNonQuotedString, uint64T)
	}
	return wrapUnmarshalError(b.UnmarshalText(input[1:len(input)-1]), uint64T)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *U64) UnmarshalText(input []byte) error {
	raw, err := validateNumber(input)
	if err != nil {
		return err
	}
	if len(raw) > bytesIn64Bits {
		return ErrUint64Range
	}
	var dec uint64
	for _, byte := range raw {
		nib := decodeNibble(byte)
		if nib == badNibble {
			return ErrInvalidString
		}
		dec *= hexBase // hex shift left :D
		dec += nib
	}
	*b = U64(dec)
	return nil
}

// XString returns the hex encoding of b.
func (b U64) String() XString {
	return FromUint64(uint64(b))
}
