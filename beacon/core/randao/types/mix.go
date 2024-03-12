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

import "github.com/berachain/beacon-kit/crypto/sha256"

// Mix represents the current state of the RANDAO's entropy mixing process.
// RANDAO can be conceptualized as a deck of cards being passed and shuffled
// by each participant, thereby continuously re-randomizing the deck.
// This process ensures that even if an individual's contribution to the
// randomness is weak, the overall entropy of the system remains high. Mix keeps
// track of this "current mix" or the state of the shuffled deck as it
// circulates among participants.
const MixLength = 32

// Mix is a fixed-size array that stores the current state of the entropy mix.
type Mix [MixLength]byte

// MixinNewReveal takes a new reveal (signature) and combines it with the
// current mix
// using a XOR operation, then returns the updated mix.
func (m Mix) MixinNewReveal(reveal Reveal) Mix {
	for idx, b := range sha256.Hash(reveal[:]) {
		m[idx] ^= b
	}
	return m
}
