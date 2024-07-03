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

package json

import (
	"encoding"
	"reflect"
)

// UnmarshalJSONText unmarshals a JSON string with a 0x prefix into a given
// TextUnmarshaler. It validates the input and then removes the surrounding
// quotes before passing the inner content to the UnmarshalText method.
func UnmarshalJSONText(input []byte,
	u encoding.TextUnmarshaler,
	t reflect.Type,
) error {
	if err := ValidateUnmarshalInput(input); err != nil {
		return WrapUnmarshalError(err, t)
	}
	return WrapUnmarshalError(u.UnmarshalText(input[1:len(input)-1]), t)
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

// isQuotedString returns true if input has quotes.
func isQuotedString[T []byte | string](input T) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}
