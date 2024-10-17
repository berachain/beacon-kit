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
	"encoding/hex"

	"github.com/berachain/beacon-kit/mod/errors"
)

var ErrInvalidHexStringLength = errors.New("invalid hex string length")

// EncodeBytes creates a hex string with 0x prefix.
// Inverse operation is ToBytes or MustToBytes.
func EncodeBytes(b []byte) string {
	hexStr := make([]byte, len(b)*2+prefixLen)
	copy(hexStr, Prefix)
	hex.Encode(hexStr[prefixLen:], b)
	return string(hexStr)
}

// MustToBytes returns the bytes represented by the given hex string.
// It panics if the input is not a valid hex string.
func MustToBytes(input string) []byte {
	bz, err := ToBytes(input)
	if err != nil {
		panic(err)
	}
	return bz
}

// ToBytes returns the bytes represented by the given hex string.
// An error is returned if the input is not a valid hex string.
func ToBytes(input string) ([]byte, error) {
	strippedInput, err := IsValidHex(input)
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(strippedInput)
}

func UnmarshalByteText(input []byte) ([]byte, error) {
	raw, err := formatAndValidateText(input)
	if err != nil {
		return []byte{}, err
	}
	dec := make([]byte, len(raw)/encDecRatio)
	if _, err = hex.Decode(dec, raw); err != nil {
		return []byte{}, err
	}
	return dec, nil
}

// DecodeFixedJSON decodes the input as a string with 0x prefix. The length
// of out determines the required input length. This function is commonly used
// to implement the UnmarshalJSON method for fixed-size types.
func DecodeFixedJSON(input, out []byte) error {
	if !isQuotedString(input) {
		return ErrNonQuotedString
	}
	return DecodeFixedText(input[1:len(input)-1], out)
}

// DecodeFixedText decodes the input as a string with 0x prefix. The length
// of out determines the required input length.
func DecodeFixedText(input, out []byte) error {
	raw, err := formatAndValidateText(input)
	if err != nil {
		return err
	}
	if len(raw)/encDecRatio != len(out) {
		return errors.Wrapf(
			ErrInvalidHexStringLength,
			"hex string has length %d, want %d",
			len(raw), len(out)*encDecRatio,
		)
	}
	// Pre-verify syntax and decode in a single pass
	for i := 0; i < len(raw); i += 2 {
		highNibble := decodeNibble(raw[i])
		lowNibble := decodeNibble(raw[i+1])
		if highNibble == badNibble || lowNibble == badNibble {
			return ErrInvalidString
		}
		out[i/2] = byte((highNibble << nibbleShift) | lowNibble)
	}

	return nil
}
