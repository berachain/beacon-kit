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

package primitives_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/stretchr/testify/require"
)

func TestU64_MarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		value    primitives.U64
		expected []byte
	}{
		{
			name:     "zero",
			value:    0,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "max uint64",
			value:    primitives.U64(^uint64(0)),
			expected: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "arbitrary number",
			value:    primitives.U64(123456789),
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
		expected primitives.U64
		err      error
	}{
		{
			name:     "valid data",
			data:     []byte{21, 205, 91, 7, 0, 0, 0, 0},
			expected: primitives.U64(123456789),
		},
		{
			name: "invalid data - short buffer",
			data: []byte{0, 0, 0},
			err:  primitives.ErrInvalidSSZLength,
		},
		{
			name:     "valid data - max uint64",
			data:     []byte{255, 255, 255, 255, 255, 255, 255, 255},
			expected: primitives.U64(^uint64(0)),
		},
		{
			name:     "valid data - zero",
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0},
			expected: primitives.U64(0),
		},
		{
			name: "invalid data - long buffer",
			data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 1},
			err:  primitives.ErrInvalidSSZLength,
		},
		{
			name:     "valid data - one",
			data:     []byte{1, 0, 0, 0, 0, 0, 0, 0},
			expected: primitives.U64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u primitives.U64
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
		value    primitives.U64
		expected []byte
	}{
		{
			name:     "zero value",
			value:    primitives.U64(0),
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "max uint64",
			value:    primitives.U64(^uint64(0)),
			expected: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "arbitrary number",
			value:    primitives.U64(123456789),
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
			var unmarshaled primitives.U64
			err = unmarshaled.UnmarshalSSZ(tt.expected)
			require.NoError(t, err)
			require.Equal(t, tt.value, unmarshaled)
		})
	}
}
