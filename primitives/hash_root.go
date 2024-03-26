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

import "github.com/ethereum/go-ethereum/common/hexutil"

const (
	HashRootLength = 32
)

// HashRoot is a 32-byte root of a hash tree structure.
type HashRoot [HashRootLength]byte

// MarshalJSON implements the json.Marshaler interface.s.
func (h HashRoot) UnmarshalJSON(input []byte) error {
	bz := (&hexutil.Bytes{})
	if err := bz.UnmarshalJSON(input); err != nil {
		return err
	}
	copy(h[:], *bz)
	return nil
}

// String returns the hex string representation of the HashRoot.
func (h HashRoot) String() string {
	return hexutil.Encode(h[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (h HashRoot) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (h HashRoot) UnmarshalText(text []byte) error {
	bz := (&hexutil.Bytes{})
	if err := bz.UnmarshalText(text); err != nil {
		return err
	}
	copy(h[:], *bz)
	return nil
}
