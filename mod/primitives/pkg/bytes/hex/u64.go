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

type HexMarshaler interface {
	MarshalHex() ([]byte, error)
	UnmarshalHex(data []byte) error
}

// MarshalText returns a byte slice containing the hexadecimal representation
// of input
func MarshalText(b uint64) ([]byte, error) {
	buf := make([]byte, prefixLen, initialCapacity)
	copy(buf, prefix)
	buf = strconv.AppendUint(buf, b, hexBase)
	return buf, nil
}

// ValidateUnmarshalInput returns true if input is a valid JSON string.
func ValidateUnmarshalInput(input []byte) error {
	if isQuotedString(string(input)) {
		return ErrNonQuotedString
	} else {
		return nil
	}
}

// GetReflectType returns the reflect.Type of i.
func GetReflectType(i any) reflect.Type {
	return reflect.TypeOf(i)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func UnmarshalText(b uint64, input []byte) (uint64, error) {
	raw, err := validateNumber(input)
	if err != nil {
		return 0, err
	}
	if len(raw) > bytesIn64Bits {
		return 0, ErrUint64Range
	}
	var dec uint64
	for _, byte := range raw {
		nib := decodeNibble(byte)
		if nib == badNibble {
			return dec, ErrInvalidString
		}
		dec *= hexBase // hex shift left :D
		dec += nib
	}
	return dec, nil
}
