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
	"strconv"
)

// This file contains functions for encoding and decoding uint64 values to and
// from hexadecimal strings, and marshaling and unmarshaling uint64 values to
// and from byte slices representing hexadecimal strings.

// MarshalText returns a byte slice containing the hexadecimal representation
// of uint64 input.
func MarshalText(b uint64) ([]byte, error) {
	buf := make([]byte, prefixLen, initialCapacity)
	copy(buf, prefix)
	buf = strconv.AppendUint(buf, b, hexBase)
	return buf, nil
}

// ValidateUnmarshalInput validates the input byte slice for unmarshaling.
// It returns an error iff input is not a quoted string.
// This is used to prevent exposing validation logic to the caller.
func ValidateUnmarshalInput(input []byte) error {
	if !isQuotedString(string(input)) {
		return ErrNonQuotedString
	}
	return nil
}

// UnmarshalUint64Text parses a byte slice containing a hexadecimal string and
// returns the uint64 value it represents.
func UnmarshalUint64Text(input []byte) (uint64, error) {
	raw, err := formatAndValidateNumber(input)
	if err != nil {
		return 0, err
	}
	if len(raw) > bytesPer64Bits {
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
