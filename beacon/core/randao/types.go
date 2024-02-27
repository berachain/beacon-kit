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
type Reveal [96]byte

// This is a hashed value of the signed reveal.
type Mix [32]byte

// This is the external representation of the randao random number
// We fix this to 32 bytes.
type RandomValue [32]byte

// MixInRandao mixes in a new reveal into the Mix
func (m *Mix) MixInRandao(newReveal Reveal) error {
	// Hash the new reveal using SHA-256 to get a blockRandaoReveal
	blockRandaoReveal := sha256.Hash(newReveal[:])

	// Check if the length of the hashed reveal matches the length of the Mix
	if len(blockRandaoReveal) != len(m) {
		// Return an error if there is a length mismatch
		return ErrMixHashRevealLengthMismatch
	}

	// Iterate over the blockRandaoReveal bytes
	for idx, b := range blockRandaoReveal {
		// XOR each byte with the corresponding byte in the Mix
		m[idx] ^= b
	}

	// Return nil to indicate success
	return nil
}

// blockRandaoReveal := sha256.Hash(randaoReveal)
// if len(blockRandaoReveal) != len(latestMixSlice) {
// 	return nil, errors.New("blockRandaoReveal length doesn't match latestMixSlice length")
// }
// for i, x := range blockRandaoReveal {
// 	latestMixSlice[i] ^= x
// }
