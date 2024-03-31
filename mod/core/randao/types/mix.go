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

package types

import (
	"github.com/berachain/beacon-kit/mod/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// Mix is a fixed-size array that stores the current state of the entropy mix.
type Mix primitives.Bytes32

// MixinNewReveal takes a new reveal (signature) and combines it with the
// current mix
// using a XOR operation, then returns the updated mix.
func (m Mix) MixinNewReveal(reveal primitives.BLSSignature) Mix {
	for idx, b := range sha256.Hash(reveal[:]) {
		m[idx] ^= b
	}
	return m
}

// MarshalText implements the encoding.TextMarshaler interface.
func (m Mix) MarshalText() ([]byte, error) {
	return primitives.Bytes32(m).MarshalText()
}
