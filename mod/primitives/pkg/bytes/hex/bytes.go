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

//nolint:gochecknoglobals // reflect.Type of Bytes set at runtime
var bytesT = reflect.TypeOf(Bytes(nil))

// Bytes marshals/unmarshals as a JSON string with 0x prefix.
// The empty slice marshals as "0x".
type Bytes []byte

// MarshalText implements encoding.TextMarshaler.
func (b Bytes) MarshalText() ([]byte, error) {
	result := make([]byte, len(b)*2+prefixLen)
	copy(result, prefix)
	hex.Encode(result[prefixLen:], b)
	return result, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	if !isQuotedString(input) {
		return WrapUnmarshalError(ErrNonQuotedString, bytesT)
	}

	return WrapUnmarshalError(b.UnmarshalText(input[1:len(input)-1]), bytesT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bytes) UnmarshalText(input []byte) error {
	raw, err := formatAndValidateText(input)
	if err != nil {
		return err
	}
	dec := make([]byte, len(raw)/encDecRatio)
	if _, err = hex.Decode(dec, raw); err != nil {
		return err
	}
	*b = dec

	return nil
}

// String returns the hex encoding of b.
func (b Bytes) String() String {
	return FromBytes(b)
}

// UnmarshalFixedJSON decodes the input as a string with 0x prefix. The length
// of out determines the required input length. This function is commonly used
// to implement the UnmarshalJSON method for fixed-size types.
func UnmarshalFixedJSON(typ reflect.Type, input, out []byte) error {
	if !isQuotedString(input) {
		return WrapUnmarshalError(ErrNonQuotedString, bytesT)
	}
	return WrapUnmarshalError(
		UnmarshalFixedText(typ.String(), input[1:len(input)-1], out), typ,
	)
}

// UnmarshalFixedText decodes the input as a string with 0x prefix. The length
// of out determines the required input length. This function is commonly used
// to implement the UnmarshalText method for fixed-size types.
func UnmarshalFixedText(typname string, input, out []byte) error {
	raw, err := formatAndValidateText(input)
	if err != nil {
		return err
	}
	if len(raw)/encDecRatio != len(out) {
		return errors.Newf(
			"hex string has length %d, want %d for %s",
			len(raw), len(out)*encDecRatio, typname,
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
