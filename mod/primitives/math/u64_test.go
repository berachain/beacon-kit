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

package math_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/stretchr/testify/require"
)

func TestU64_MarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected []byte
	}{
		{
			name:     "zero",
			value:    0,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "max uint64",
			value:    math.U64(^uint64(0)),
			expected: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "arbitrary number",
			value:    math.U64(123456789),
			expected: []byte{21, 205, 91, 7, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.value.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_UnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected math.U64
		err      error
	}{
		{
			name:     "valid data",
			data:     []byte{21, 205, 91, 7, 0, 0, 0, 0},
			expected: math.U64(123456789),
		},
		{
			name: "invalid data - short buffer",
			data: []byte{0, 0, 0},
			err:  math.ErrUnexpectedInputLengthBase,
		},
		{
			name:     "valid data - max uint64",
			data:     []byte{255, 255, 255, 255, 255, 255, 255, 255},
			expected: math.U64(^uint64(0)),
		},
		{
			name:     "valid data - zero",
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0},
			expected: math.U64(0),
		},
		{
			name: "invalid data - long buffer",
			data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 1},
			err:  math.ErrUnexpectedInputLengthBase,
		},
		{
			name:     "valid data - one",
			data:     []byte{1, 0, 0, 0, 0, 0, 0, 0},
			expected: math.U64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u math.U64
			err := u.UnmarshalSSZ(tt.data)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, u)
			}
		})
	}
}

func TestU64_RoundTripSSZ(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected []byte
	}{
		{
			name:     "zero value",
			value:    math.U64(0),
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "max uint64",
			value:    math.U64(^uint64(0)),
			expected: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "arbitrary number",
			value:    math.U64(123456789),
			expected: []byte{21, 205, 91, 7, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test MarshalSSZ
			marshaled, err := tt.value.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.expected, marshaled)

			// Test UnmarshalSSZ
			var unmarshaled math.U64
			err = unmarshaled.UnmarshalSSZ(tt.expected)
			require.NoError(t, err)
			require.Equal(t, tt.value, unmarshaled)
		})
	}
}

func TestU64_NextPowerOfTwo(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected math.U64
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: math.U64(0),
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: math.U64(1),
		},
		{
			name:     "already a power of two",
			value:    math.U64(8),
			expected: math.U64(8),
		},
		{
			name:     "not a power of two",
			value:    math.U64(9),
			expected: math.U64(16),
		},
		{
			name:     "large number",
			value:    math.U64(1<<63 - 1),
			expected: math.U64(1 << 63),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.NextPowerOfTwo()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_NextPowerOfTwoPanic(t *testing.T) {
	u := ^math.U64(0)
	require.Panics(t, func() {
		_ = u.NextPowerOfTwo()
	})
}

func TestU64_ILog2Ceil(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected uint8
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: 0,
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: 0,
		},
		{
			name:     "power of two",
			value:    math.U64(8),
			expected: 3,
		},
		{
			name:     "not a power of two",
			value:    math.U64(9),
			expected: 4,
		},
		{
			name:     "max uint64",
			value:    math.U64(1<<64 - 1),
			expected: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.ILog2Ceil()
			require.Equal(t, tt.expected, result)
		})
	}
}
func TestU64_PrevPowerOfTwo(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected math.U64
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: 1,
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: 1,
		},
		{
			name:     "two",
			value:    math.U64(2),
			expected: 2,
		},
		{
			name:     "three",
			value:    math.U64(3),
			expected: 2,
		},
		{
			name:     "four",
			value:    math.U64(4),
			expected: 4,
		},
		{
			name:     "five",
			value:    math.U64(5),
			expected: 4,
		},
		{
			name:     "eight",
			value:    math.U64(8),
			expected: 8,
		},
		{
			name:     "nine",
			value:    math.U64(9),
			expected: 8,
		},
		{
			name:     "thirty-two",
			value:    math.U64(32),
			expected: 32,
		},
		{
			name:     "thirty-three",
			value:    math.U64(33),
			expected: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.PrevPowerOfTwo()
			require.Equal(t, tt.expected, result)
		})
	}
}

// func TestU64List_HashTreeRoot(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		vector   math.U64List
// 		expected [32]byte
// 	}{
// 		// {
// 		// 	name:     "empty vector",
// 		// 	vector:   math.U64List{},
// 		// 	expected: [32]byte{}, // Assuming the hash of an empty vector is
// 		// zeroed
// 		// },
// 		// {
// 		// 	name:     "single element",
// 		// 	vector:   math.U64List{1},
// 		// 	expected: [32]byte{0x01}, // Simplified expected result
// 		// },
// 		{
// 			name: "multiple elements",
// 			vector: math.U64List{
// 				1,
// 				2,
// 				3,
// 				4,
// 				5,
// 				6,
// 				9,
// 				1000,
// 				34,
// 				334,
// 				33,
// 			},
// 			expected: [32]byte{0x0e}, // Simplified expected result
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := tt.vector.HashTreeRoot()
// 			require.NoError(t, err)
// 			list2 := make([]uint64, len(tt.vector))
// 			for i, v := range tt.vector {
// 				list2[i] = uint64(v)
// 			}
// 			result2, err := (&math.U64List2{Data: list2}).HashTreeRoot()
// 			require.Equal(t, result, result2)
// 		})
// 	}
// }

// func TestU64Vector_HashTreeRoot(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		container math.U64Container
// 		expected  [32]byte
// 	}{
// 		// {
// 		// 	name:     "empty vector",
// 		// 	vector:   math.U64List{},
// 		// 	expected: [32]byte{}, // Assuming the hash of an empty vector is
// 		// zeroed
// 		// },
// 		// {
// 		// 	name:     "single element",
// 		// 	vector:   math.U64List{1},
// 		// 	expected: [32]byte{0x01}, // Simplified expected result
// 		// },
// 		{
// 			name: "multiple elements",
// 			container: math.U64Container{
// 				// Field0: 1,
// 				Field1: 2,
// 				Field2: math.U64List{
// 					1,
// 					2,
// 					3,
// 					4,
// 					5,
// 					6,
// 					9,
// 					1000,
// 					34,
// 					334,
// 					33,
// 				},
// 			},
// 			expected: [32]byte{0x0e}, // Simplified expected result
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := tt.container.HashTreeRoot()
// 			require.NoError(t, err)

// 			result2, err := (&math.U64Container2{
// 				// Field0: 1,
// 				Field1: 2,
// 				Field2: []uint64{1, 2, 3, 4, 5, 6, 9, 1000, 34, 334, 33},
// 			}).HashTreeRoot()
// 			require.Equal(t, result, result2)
// 		})
// 	}
// }
