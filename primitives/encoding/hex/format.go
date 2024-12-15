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

import (
	"errors"
	"strings"
)

// IsValidHex performs basic validations that every hex string
// must pass (there may be extra ones depending on the type encoded)
// It returns the suffix (dropping 0x prefix) in the hope to appease nilaway.
func IsValidHex[T ~[]byte | ~string](s T) (T, error) {
	if len(s) == 0 {
		return *new(T), ErrEmptyString
	}
	if len(s) < prefixLen {
		return *new(T), ErrMissingPrefix
	}
	if strings.ToLower(string(s[:prefixLen])) != Prefix {
		return *new(T), ErrMissingPrefix
	}
	return s[prefixLen:], nil
}

// ValidateQuotedString errs if input has no quotes.
// For convenience it returns the unstrip content if it does not err.
func ValidateQuotedString(input []byte) ([]byte, error) {
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		return input[1 : len(input)-1], nil
	}
	return nil, ErrNonQuotedString
}

// formatAndValidateText validates the input text for a hex string.
func formatAndValidateText(input []byte) ([]byte, error) {
	input, err := IsValidHex(input)
	if errors.Is(err, ErrEmptyString) {
		return nil, nil // empty strings are allowed
	} else if err != nil {
		return nil, err
	}

	if len(input)%2 != 0 {
		return nil, ErrOddLength
	}
	return input, nil
}

// formatAndValidateNumber checks the input text for a hex number.
func formatAndValidateNumber[T []byte | string](input T) (T, error) {
	input, err := IsValidHex(input)
	if err != nil {
		return *new(T), err
	}

	if len(input) == 0 {
		return *new(T), ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return *new(T), ErrLeadingZero
	}
	return input, nil
}
