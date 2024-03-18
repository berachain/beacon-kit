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

package bytes_test

import (
	"bytes"
	"reflect"
	"testing"

	byteslib "github.com/berachain/beacon-kit/lib/bytes"
)

func TestSafeCopy(t *testing.T) {
	tests := []struct {
		name     string
		original []byte
	}{
		{name: "Normal case", original: []byte{1, 2, 3, 4, 5}},
		{name: "Empty slice", original: []byte{}},
		{name: "Single element slice", original: []byte{9}},
		{name: "Large slice", original: make([]byte, 100)},
		{name: "Another normal case", original: []byte{6, 6, 6, 6, 6}},
		{name: "Another single element slice", original: []byte{5}},
		{name: "Another large slice", original: make([]byte, 200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := byteslib.SafeCopy(tt.original)

			if !bytes.Equal(tt.original, copied) {
				t.Errorf("SafeCopy did not copy the slice correctly")
			}

			// Modifying the copied slice should not affect the original slice
			if len(copied) > 0 {
				copied[0] = 10
				if tt.original[0] == copied[0] {
					t.Errorf(
						"Modifying the copied slice affected the original slice",
					)
				}
			}
		})
	}
}

func TestSafeCopy2D(t *testing.T) {
	tests := []struct {
		name     string
		original [][]byte
	}{
		{
			name: "Normal case",
			original: [][]byte{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			name:     "Empty slice",
			original: [][]byte{},
		},
		{
			name: "Single element slice",
			original: [][]byte{
				{9},
			},
		},
		{
			name: "Mixed lengths",
			original: [][]byte{
				{1, 2, 3},
				{4},
				{5, 6},
			},
		},
		{
			name: "Nil inner slice",
			original: [][]byte{
				nil,
				{1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := byteslib.SafeCopy2D(tt.original)

			if !reflect.DeepEqual(tt.original, copied) {
				t.Errorf("SafeCopy2D did not copy the slice correctly")
			}

			// Modifying the copied slice should not affect the original slice
			if len(copied) > 0 && len(copied[0]) > 0 {
				copied[0][0] = 10
				if tt.original[0][0] == copied[0][0] {
					t.Errorf(
						"Modifying the copied slice affected the original slice",
					)
				}
			}
		})
	}
}

func TestReverseEndianness(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{name: "Even length",
			input:    []byte{1, 2, 3, 4},
			expected: []byte{4, 3, 2, 1}},
		{name: "Odd length",
			input:    []byte{1, 2, 3, 4, 5},
			expected: []byte{5, 4, 3, 2, 1}},
		{name: "Empty slice",
			input:    []byte{},
			expected: []byte{}},
		{name: "Single element",
			input:    []byte{1},
			expected: []byte{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := byteslib.CopyAndReverseEndianess(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf(
					"ReverseEndianness(%v) = %v, want %v",
					tt.name, result, tt.expected)
			}
		})
	}
}
