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

package encoding

import (
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections/codec"
)

// Bytes32ValueCodec provides methods to encode and decode [32]byte values.
type Bytes32ValueCodec struct{}

// Assert that Bytes32ValueCodec implements codec.ValueCodec.
var _ codec.ValueCodec[[32]byte] = Bytes32ValueCodec{}

// Encode marshals the provided value into its [32]byte encoding.
func (Bytes32ValueCodec) Encode(value [32]byte) ([]byte, error) {
	return value[:], nil
}

// Decode unmarshals the provided bytes into a value of type [32]byte.
func (Bytes32ValueCodec) Decode(b []byte) ([32]byte, error) {
	var v [32]byte
	if len(b) != len(v) {
		return v, fmt.Errorf(
			"invalid length: expected %d, got %d",
			len(v),
			len(b),
		)
	}
	copy(v[:], b)
	return v, nil
}

// EncodeJSON is not implemented and will panic if called.
func (Bytes32ValueCodec) EncodeJSON(_ [32]byte) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (Bytes32ValueCodec) DecodeJSON(_ []byte) ([32]byte, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (Bytes32ValueCodec) Stringify(value [32]byte) string {
	return hex.EncodeToString(value[:])
}

// ValueType returns the name of the interface that this codec is intended for.
func (Bytes32ValueCodec) ValueType() string {
	return "Bytes32"
}
