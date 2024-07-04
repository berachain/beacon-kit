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

func EncodeBytes[B ~[]byte](b B) ([]byte, error) {
	result := make([]byte, len(b)*2+prefixLen)
	copy(result, prefix)
	hex.Encode(result[prefixLen:], b)
	return result, nil
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
		return errors.Newf(
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
