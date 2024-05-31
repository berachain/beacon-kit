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

// ------------------------------ B48 ------------------------------

// B48 represents a 48-byte array.
type B48 [48]byte

// ToBytes48 is a utility function that transforms a byte slice into a fixed
// 48-byte array. If the input exceeds 48 bytes, it gets truncated.
func ToBytes48(input []byte) B48 {
	//nolint:mnd // 48 bytes.
	return B48(ExtendToSize(input, 48))
}

// UnmarshalJSON implements the json.Unmarshaler interface for B48.
func (h *B48) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// String returns the hex string representation of B48.
func (h B48) String() string {
	return hex.FromBytes(h[:]).Unwrap()
}

// MarshalText implements the encoding.TextMarshaler interface for B48.
func (h B48) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for B48.
func (h *B48) UnmarshalText(text []byte) error {
	return UnmarshalTextHelper(h[:], text)
}
