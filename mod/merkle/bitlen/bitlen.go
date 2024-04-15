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

package bitlen

const (
	mask0 = ^uint64((1 << (1 << iota)) - 1)
	mask1
	mask2
	mask3
	mask4
	mask5
)

const (
	bit0 = uint8(1 << iota)
	bit1
	bit2
	bit3
	bit4
	bit5
)

// bitmagic: binary search through a uint64, to find the BitIndex of next power
// of 2 (if not already a power of 2)
// Zero is a special case, it has a 0 depth.
// Example:
//
//	(in out): (0 0), (1 0), (2 1), (3 2), (4 2), (5 3), (6 3), (7 3), (8 3), (9
//
// 4).
func CoverDepth(v uint64) uint8 {
	if v == 0 || v == 1 {
		return 0
	}
	return BitIndex(v-1) + 1
}

// bitmagic: binary search through a uint64 to find the bit-length
// Zero is a special case, it has a 0 bit length.
// Example:
//
//	(in out): (0 0), (1 1), (2 2), (3 2), (4 3), (5 3), (6 3), (7 3), (8 4), (9
//
// 4).
func BitLength(v uint64) uint8 {
	if v == 0 {
		return 0
	}
	return BitIndex(v) + 1
}

// bitmagic: binary search through a uint64 to find the index (least bit being
// 0) of the first set bit.
// Zero is a special case, it has a 0 bit index.
// Example:
//
//	(in out): (0 0), (1 0), (2 1), (3 1), (4 2), (5 2), (6 2), (7 2), (8 3), (9
//
// 3).
func BitIndex(v uint64) uint8 {
	if v == 0 {
		return 0
	}
	var out uint8
	//nolint:gomnd // 5 is allowed.
	for shift := uint8(5); shift > 0; shift-- {
		mask := ^uint64((1 << (1 << shift)) - 1)
		bit := uint8(1 << shift)
		if v&mask != 0 {
			v >>= bit
			out |= bit
		}
	}
	if v&mask0 != 0 {
		out |= bit0
	}
	return out
}
