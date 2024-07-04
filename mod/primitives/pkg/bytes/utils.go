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

package bytes

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

// ------------------------------ Helpers ------------------------------

// Helper function to unmarshal JSON for various byte types.
func unmarshalJSONHelper(target []byte, input []byte) error {
	bz := Bytes{}
	if err := bz.UnmarshalJSON(input); err != nil {
		return err
	}
	if len(bz) != len(target) {
		return errors.Newf(
			"incorrect length, expected %d bytes but got %d",
			len(target), len(bz),
		)
	}
	copy(target, bz)
	return nil
}

// UnmarshalTextHelper function to unmarshal text for various byte types.
func UnmarshalTextHelper(target []byte, text []byte) error {
	bz := Bytes{}
	if err := bz.UnmarshalText(text); err != nil {
		return err
	}
	if len(bz) != len(target) {
		return errors.Newf(
			"incorrect length, expected %d bytes but got %d",
			len(target), len(bz),
		)
	}
	copy(target, bz)
	return nil
}

// MustFromHex returns the bytes represented by the given hex string.
// It panics if the input is not a valid hex string.
func MustFromHex(input string) []byte {
	bz, err := FromHex(input)
	if err != nil {
		panic(err)
	}
	return bz
}

// FromHex returns the bytes represented by the given hex string.
// An error is returned if the input is not a valid hex string.
func FromHex(input string) ([]byte, error) {
	s, err := hex.NewStringStrict(input)
	if err != nil {
		return nil, err
	}
	h, err := s.ToBytes()
	if err != nil {
		return nil, err
	}
	return h, nil
}

// CopyAndReverseEndianess will copy the input byte slice and return the
// flipped version of it.
func CopyAndReverseEndianess(input []byte) []byte {
	copied := make([]byte, len(input))
	copy(copied, input)
	for i, j := 0, len(copied)-1; i < j; i, j = i+1, j-1 {
		copied[i], copied[j] = copied[j], copied[i]
	}
	return copied
}

// ExtendToSize extends a byte slice to a specified length. It returns the
// original slice if it's already larger.
func ExtendToSize(slice []byte, length int) []byte {
	if len(slice) >= length {
		return slice
	}
	return append(slice, make([]byte, length-len(slice))...)
}

// UnmarshalFixedJSON decodes the input as a string with 0x prefix. The length
// of out determines the required input length. This function is commonly used
// to implement the UnmarshalJSON method for fixed-size types.
func UnmarshalFixedJSON(input, out []byte) error {
	return hex.DecodeFixedJSON(input, out)
}

// UnmarshalFixedText decodes the input as a string with 0x prefix. The length
// of out determines the required input length. This function is commonly used
// to implement the UnmarshalText method for fixed-size types.
func UnmarshalFixedText(input, out []byte) error {
	return hex.DecodeFixedText(input, out)
}
