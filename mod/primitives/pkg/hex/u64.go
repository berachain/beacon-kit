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
