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

package version_test

import (
	"encoding/binary"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

func TestFromUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected [4]byte
	}{
		{
			name:     "Phase0",
			input:    version.Phase0,
			expected: [4]byte{0, 0, 0, 0},
		},
		{
			name:     "Altair",
			input:    version.Altair,
			expected: [4]byte{1, 0, 0, 0},
		},
		{
			name:     "Bellatrix",
			input:    version.Bellatrix,
			expected: [4]byte{2, 0, 0, 0},
		},
		{
			name:     "Capella",
			input:    version.Capella,
			expected: [4]byte{3, 0, 0, 0},
		},
		{
			name:     "Deneb",
			input:    version.Deneb,
			expected: [4]byte{4, 0, 0, 0},
		},
		{
			name:     "Electra",
			input:    version.Electra,
			expected: [4]byte{5, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := version.FromUint32[[4]byte](tt.input)
			if result != tt.expected {
				t.Errorf(
					"FromUint32(%d) = %v, expected %v",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestFromUint32_CustomType(t *testing.T) {
	type CustomVersion [4]byte

	input := uint32(123456789)
	expected := CustomVersion{}
	binary.LittleEndian.PutUint32(expected[:], input)

	result := version.FromUint32[CustomVersion](input)
	if result != expected {
		t.Errorf("FromUint32(%d) = %v, expected %v", input, result, expected)
	}
}
