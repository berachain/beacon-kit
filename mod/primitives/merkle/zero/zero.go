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

package zero

import sha256 "github.com/minio/sha256-simd"

// NumZeroHashes is the number of pre-computed zero-hashes.
const NumZeroHashes = 64

// Hashes is a pre-computed list of zero-hashes for each depth level.
//
//nolint:gochecknoglobals // saves recomputing.
var Hashes [NumZeroHashes + 1][32]byte

// initialize the zero-hashes pre-computed data with the given hash-function.
func InitZeroHashes(zeroHashesLevels int) {
	for i := range zeroHashesLevels {
		v := [64]byte{}
		copy(v[:32], Hashes[i][:])
		copy(v[32:], Hashes[i][:])
		Hashes[i+1] = sha256.Sum256(v[:])
	}
}

//nolint:init // saves recomputing.
func init() {
	InitZeroHashes(NumZeroHashes)
}
