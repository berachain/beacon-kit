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

// Bytes32KeyCodec provides methods to encode and decode [32]byte keys.
type Bytes32KeyCodec struct{}

// Assert that Bytes32ValueCodec implements codec.ValueCodec.
var _ codec.KeyCodec[[32]byte] = Bytes32KeyCodec{}

var (
	Bytes32Key   = Bytes32KeyCodec{}
	Bytes32Value = codec.KeyToValueCodec(Bytes32Key)
)

// Encode writes the key 32-bytes into the buffer.
// Returns the number of bytes written.
func (Bytes32KeyCodec) Encode(buffer []byte, value [32]byte) (int, error) {
	return copy(buffer, value[:]), nil
}

// Decode unmarshals the provided bytes into a value of type [32]byte.
func (Bytes32KeyCodec) Decode(buffer []byte) (int, [32]byte, error) {
	var v [32]byte
	if len(buffer) != len(v) {
		return 0, v, fmt.Errorf(
			"invalid length: expected %d, got %d",
			len(v),
			len(buffer),
		)
	}
	copiedBytes := copy(v[:], buffer)
	return copiedBytes, v, nil
}

// Size returns the buffer size need to encode 32-byte key.
func (Bytes32KeyCodec) Size(key [32]byte) int {
	return len(key)
}

// EncodeJSON is not implemented and will panic if called.
func (Bytes32KeyCodec) EncodeJSON(_ [32]byte) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (Bytes32KeyCodec) DecodeJSON(_ []byte) ([32]byte, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (Bytes32KeyCodec) Stringify(key [32]byte) string {
	return hex.EncodeToString(key[:])
}

// KeyType returns the name of the interface that this codec is intended for.
func (Bytes32KeyCodec) KeyType() string {
	return "Bytes32"
}

// EncodeNonTerminal writes the key bytes into the buffer.
func (b Bytes32KeyCodec) EncodeNonTerminal(buffer []byte, key [32]byte) (int, error) {
	return b.Encode(buffer, key)
}

// DecodeNonTerminal reads the buffer provided and
// returns the 32-byte key.
func (b Bytes32KeyCodec) DecodeNonTerminal(buffer []byte) (int, [32]byte, error) {
	return b.Decode(buffer)
}

// SizeNonTerminal returns the buffer size need to encode 32-byte key.
func (Bytes32KeyCodec) SizeNonTerminal(key [32]byte) int {
	return len(key)
}
