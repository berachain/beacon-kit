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

import "errors"

// isQuotedString returns true if input has quotes.
func isQuotedString[T []byte | string](input T) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
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
