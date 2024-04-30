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

package bytes

// SafeCopy creates a copy of the provided byte slice. If the input slice is
// non-nil and has a length of 32 bytes, it assumes the slice represents a hash
// and copies it into a fixed-size array before returning a slice of that array.
// For other non-nil slices, it returns a dynamically allocated copy. If the
// input slice is nil, it returns nil.
func SafeCopy(src []byte) []byte {
	if src == nil {
		return nil
	}

	//nolint:mnd // 32 bytes.
	if len(src) == 32 {
		var copied [32]byte
		copy(copied[:], src)
		return copied[:]
	}

	copied := make([]byte, len(src))
	copy(copied, src)
	return copied
}

// SafeCopy2D creates a copy of a two-dimensional byte slice. It iterates over
// the outer slice, copying each inner slice using SafeCopy. If the input is
// non-nil, it returns a copy of the
// two-dimensional slice. If the input is nil, it returns nil.
func SafeCopy2D(src [][]byte) [][]byte {
	if src == nil {
		return nil
	}

	copied := make([][]byte, len(src))
	for i, s := range src {
		copied[i] = SafeCopy(s)
	}
	return copied
}

// CopyAndReverseEndianess will copy the input byte slice and return the
// flipped version of it.
func CopyAndReverseEndianess(input []byte) []byte {
	copied := make([]byte, len(input))
	copy(copied, input)
	for i, j := 0, len(copied)-1; i < j; i, j = i+1, j-1 {
		copied[i], copied[j] = copied[j], copied[i]
	}
	return copied
}

// ToBytes32 is a utility function that transforms a byte slice into a fixed
// 32-byte array. If the input exceeds 32 bytes, it gets truncated.
func ToBytes32(input []byte) [32]byte {
	//nolint:mnd // 32 bytes.
	return [32]byte(ExtendToSize(input, 32))
}

// ToBytes48 is a utility function that transforms a byte slice into a fixed
// 48-byte array. If the input exceeds 48 bytes, it gets truncated.
func ToBytes48(input []byte) [48]byte {
	//nolint:mnd // 32 bytes.
	return [48]byte(ExtendToSize(input, 48))
}

// ToBytes96 is a utility function that transforms a byte slice into a fixed
// 96-byte array. If the input exceeds 96 bytes, it gets truncated.
func ToBytes96(input []byte) [96]byte {
	//nolint:mnd // 32 bytes.
	return [96]byte(ExtendToSize(input, 96))
}

// ExtendToSize extends a byte slice to a specified length. It returns the
// original slice if it's already larger.
func ExtendToSize(slice []byte, length int) []byte {
	if len(slice) >= length {
		return slice
	}
	return append(slice, make([]byte, length-len(slice))...)
}

// PrependExtendToSize extends a byte slice to a specified length by
// prepending zero bytes. It returns the original slice if it's
// already larger.
func PrependExtendToSize(slice []byte, length int) []byte {
	if len(slice) >= length {
		return slice
	}
	return append(make([]byte, length-len(slice)), slice...)
}
