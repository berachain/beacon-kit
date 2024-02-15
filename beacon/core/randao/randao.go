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

package randao

import "github.com/itsdevbear/bolaris/crypto/sha256"

// This is the internal representation of the randao reveal
// Although it is 32 bytes now, it can change
// We use the same size as Ed25519 sig
// TODO: update to 96 bytes when moving to BLS
type Reveal [32]byte

// This is a hashed value of the signed reveal.
type Mix [32]byte

// This is the external representation of the randao random number
// We fix this to 32 bytes.
type RandomValue [32]byte

func (m *Mix) MixInRandao(newReveal Reveal) error {
	hash := sha256.Hash(newReveal[:])

	if len(hash) != len(m) {
		return ErrMixHashRevealLengthMismatch
	}

	for idx, b := range hash {
		m[idx] ^= b
	}

	return nil
}
