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

// has0xPrefix returns true if s has a 0x prefix.
func has0xPrefix[T []byte | string](s T) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// ensure0xPrefix ensures that s has a 0x prefix. If it doesn't, it adds it.
func ensure0xPrefix[T []byte | string](s T) T {
	if has0xPrefix(s) {
		return s
	}
	switch v := any(s).(type) {
	case string:
		return T("0x" + v)
	case []byte:
		return T(append([]byte("0x"), v...))
	default:
		return s
	}
}

// ensureStringInvariants ensures that String invariants are met by appending
// 0x prefix if missing, and converting empty string to "0x0".
func ensureStringInvariants(s string) string {
	if len(s) == 0 {
		s = "0"
	}
	return ensure0xPrefix(s)
}

// isQuotedString returns true if input has quotes.
func isQuotedString[T []byte | string](input T) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// formatAndValidateText validates the input text for a hex string.
func formatAndValidateText(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if !has0xPrefix(input) {
		return nil, ErrMissingPrefix
	}
	input = input[2:]
	if len(input)%2 != 0 {
		return nil, ErrOddLength
	}
	return input, nil
}

// formatAndValidateNumber checks the input text for a hex number.
func formatAndValidateNumber[T []byte | string](input T) (T, error) {
	// realistically, this shouldn't rarely error if called on
	// unwrapped hex.String
	if len(input) == 0 {
		return *new(T), ErrEmptyString
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
