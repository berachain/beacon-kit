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
