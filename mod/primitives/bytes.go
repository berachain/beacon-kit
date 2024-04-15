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

package primitives

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Bytes represents a byte array.
type Bytes = hexutil.Bytes

// ------------------------------ Bytes4 ------------------------------

// Bytes4 represents a 4-byte array.
type Bytes4 [4]byte

// UnmarshalJSON implements the json.Unmarshaler interface for Bytes4.
func (h *Bytes4) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// String returns the hex string representation of Bytes4.
func (h Bytes4) String() string {
	return hexutil.Encode(h[:])
}

// MarshalText implements the encoding.TextMarshaler interface for Bytes4.
func (h Bytes4) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes4.
func (h *Bytes4) UnmarshalText(text []byte) error {
	return unmarshalTextHelper(h[:], text)
}

// ------------------------------ Bytes32 ------------------------------

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

// MarshalText implements the encoding.TextMarshaler interface for Bytes32.
func (h Bytes32) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes32.
func (h *Bytes32) UnmarshalText(text []byte) error {
	return unmarshalTextHelper(h[:], text)
}

// ------------------------------ Bytes48 ------------------------------

// Bytes48 represents a 48-byte array.
type Bytes48 [48]byte

// UnmarshalJSON implements the json.Unmarshaler interface for Bytes48.
func (h *Bytes48) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// String returns the hex string representation of Bytes48.
func (h Bytes48) String() string {
	return hexutil.Encode(h[:])
}

// MarshalText implements the encoding.TextMarshaler interface for Bytes48.
func (h Bytes48) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes48.
func (h *Bytes48) UnmarshalText(text []byte) error {
	return unmarshalTextHelper(h[:], text)
}

// ------------------------------ Bytes96 ------------------------------

// Bytes96 represents a 96-byte array.
type Bytes96 [96]byte

// UnmarshalJSON implements the json.Unmarshaler interface for Bytes96.
func (h *Bytes96) UnmarshalJSON(input []byte) error {
	return unmarshalJSONHelper(h[:], input)
}

// String returns the hex string representation of Bytes96.
func (h Bytes96) String() string {
	return hexutil.Encode(h[:])
}

// MarshalText implements the encoding.TextMarshaler interface for Bytes96.
func (h Bytes96) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes96.
func (h *Bytes96) UnmarshalText(text []byte) error {
	return unmarshalTextHelper(h[:], text)
}

// ------------------------------ Helpers ------------------------------

// Helper function to unmarshal JSON for various byte types.
func unmarshalJSONHelper(target []byte, input []byte) error {
	bz := hexutil.Bytes{}
	if err := bz.UnmarshalJSON(input); err != nil {
		return err
	}
	if len(bz) != len(target) {
		return fmt.Errorf(
			"incorrect length, expected %d bytes but got %d",
			len(target), len(bz),
		)
	}
	copy(target, bz)
	return nil
}

// Helper function to unmarshal text for various byte types.
func unmarshalTextHelper(target []byte, text []byte) error {
	bz := hexutil.Bytes{}
	if err := bz.UnmarshalText(text); err != nil {
		return err
	}
	if len(bz) != len(target) {
		return fmt.Errorf(
			"incorrect length, expected %d bytes but got %d",
			len(target), len(bz),
		)
	}
	copy(target, bz)
	return nil
}
