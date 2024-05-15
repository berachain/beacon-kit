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
	"encoding/hex"
	"reflect"

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

// UnmarshalFixedJSON decodes the input as a string with 0x prefix. The length
// of out determines the required input length. This function is commonly used
// to implement the UnmarshalJSON method for fixed-size types.

// UnmarshalFixedJSON decodes the input as a string with 0x prefix.
func DecodeFixedJSON(typ reflect.Type,
	bytesT reflect.Type,
	input,
	out []byte) error {
	if !isQuotedString(input) {
		return WrapUnmarshalError(ErrNonQuotedString, bytesT)
	}
	return WrapUnmarshalError(
		DecodeFixedText(typ.String(), input[1:len(input)-1], out), typ,
	)
}

// UnmarshalFixedText decodes the input as a string with 0x prefix. The length
// of out determines the required input length.
func DecodeFixedText(typename string, input, out []byte) error {
	raw, err := formatAndValidateText(input)
	if err != nil {
		return err
	}
	if len(raw)/encDecRatio != len(out) {
		return errors.Newf(
			"hex string has length %d, want %d for %s",
			len(raw), len(out)*encDecRatio, typename,
		)
	}
	// Pre-verify syntax before modifying out.
	for _, b := range raw {
		if decodeNibble(b) == badNibble {
			return ErrInvalidString
		}
	}
	if _, err = hex.Decode(out, raw); err != nil {
		return err
	}

	return nil
}
