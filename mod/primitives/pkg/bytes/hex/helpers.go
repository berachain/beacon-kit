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
	"math/big"
	"reflect"
)

// has0xPrefix returns true if s has a 0x prefix.
func has0xPrefix[T []byte | string](s T) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// isQuotedString returns true if input has quotes.
func isQuotedString[T []byte | string](input T) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// validateText validates the input text for a hex string.
func validateText(input []byte, wantPrefix bool) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if has0xPrefix(input) {
		input = input[2:]
	} else if wantPrefix {
		return nil, ErrMissingPrefix
	}
	if len(input)%2 != 0 {
		return nil, ErrOddLength
	}
	return input, nil
}

// validateNumber checks the input text for a hex number.
func validateNumber[T []byte | string](input T) (T, error) {
	if len(input) == 0 {
		return *new(T), nil // empty strings are allowed
	}
	if !has0xPrefix(input) {
		return *new(T), ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return *new(T), ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return *new(T), ErrLeadingZero
	}
	return input, nil
}

// wrapUnmarshalError wraps an error occurring during JSON unmarshaling.
func wrapUnmarshalError(err error, t reflect.Type) error {
	if err != nil {
		err = &json.UnmarshalTypeError{Value: err.Error(), Type: t}
	}

	return err
}

func decodeNibble(in byte) uint64 {
	offset := 10
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + byte(offset))
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + byte(offset))
	default:
		return badNibble
	}
}

//nolint:mnd // this is fine xD
func getBigWordNibbles() int {
	// This is a weird way to compute the number of nibbles required for
	// big.Word. The usual way would be to use constant arithmetic but go vet
	// can't handle that

	var bigWordNibbles int
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		bigWordNibbles = 16
	case 2:
		bigWordNibbles = 8
	default:
		panic("weird big.Word size")
	}
	return bigWordNibbles
}
