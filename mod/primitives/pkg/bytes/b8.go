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

package bytes

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/hex"
)

// B8 represents a 4-byte array.
type B8 [8]byte

// UnmarshalJSON implements the json.Unmarshaler interface for B8.
func (h *B8) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// ToBytes8 is a utility function that transforms a byte slice into a fixed
// 8-byte array. If the input exceeds 4 bytes, it gets truncated.
func ToBytes8(input []byte) B8 {
	//nolint:mnd // 8 bytes.
	return [8]byte(ExtendToSize(input, 8))
}

// String returns the hex string representation of B8.
func (h B8) String() string {
	return hex.FromBytes(h[:]).Unwrap()
}

// MarshalText implements the encoding.TextMarshaler interface for B8.
func (h B8) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B8.
func (h *B8) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}
