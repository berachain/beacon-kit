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

import "github.com/ethereum/go-ethereum/common/hexutil"

// Bytes32 represents a 32-byte array.
type Bytes32 [32]byte

// UnmarshalJSON implements the json.Unmarshaler interface for Bytes32.
func (h *Bytes32) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// String returns the hex string representation of Bytes32.
func (h Bytes32) String() string {
	return hexutil.Encode(h[:])
}

// HashTreeRoot returns the hash tree root of the Bytes32.
func (h Bytes32) HashTreeRoot() ([32]byte, error) {
	return h, nil
}

// MarshalText implements the encoding.TextMarshaler interface for Bytes32.
func (h Bytes32) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes32.
func (h *Bytes32) UnmarshalText(text []byte) error {
	return unmarshalTextHelper(h[:], text)
}

// SizeSSZ returns the size of its SSZ encoding in bytes.
func (h Bytes32) SizeSSZ() int {
	//nolint:mnd // vibes.
	return 32
}
