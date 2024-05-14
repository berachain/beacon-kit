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

package hex

import (
	"math/big"
)

// decodeNibble decodes a single hexadecimal nibble (half-byte) into uint64.
func decodeNibble(in byte) uint64 {
	// uint64 conversion here is safe
	switch {
	case in >= '0' && in <= '9' && in >= hexBaseOffset:
		return uint64(in - hexBaseOffset) //#nosec G701
	case in >= 'A' && in <= 'F' && in >= hexAlphaOffsetUpper:
		return uint64(in - hexAlphaOffsetUpper) //#nosec G701
	case in >= 'a' && in <= 'f' && in >= hexAlphaOffsetLower:
		return uint64(in - hexAlphaOffsetLower) //#nosec G701
	default:
		return badNibble
	}
}

// getBigWordNibbles returns the number of nibbles required for big.Word.
//
//nolint:mnd // this is fine xD
func getBigWordNibbles() (int, error) {
	// This is a weird way to compute the number of nibbles required for
	// big.Word. The usual way would be to use constant arithmetic but go vet
	// can't handle that
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		return 16, nil
	case 2:
		return 8, nil
	default:
		return 0, ErrInvalidBigWordSize
	}
}
