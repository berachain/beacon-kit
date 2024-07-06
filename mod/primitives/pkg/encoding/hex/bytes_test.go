// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:lll // long strings
package hex_test

import (
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestEncodeBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "typical byte slice",
			input:    []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expected: []byte("0x48656c6c6f"),
		},
		{
			name:     "empty byte slice",
			input:    []byte{},
			expected: []byte("0x"),
		},
		{
			name:     "single byte",
			input:    []byte{0x01},
			expected: []byte("0x01"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hex.EncodeBytes(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result, "Test case : %s", tt.name)
		})
	}
}

func TestUnmarshalByteText(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expected  []byte
		expectErr bool
	}{
		{
			name:      "valid hex string",
			input:     []byte("0x48656c6c6f"),
			expected:  []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expectErr: false,
		},
		{
			name:      "empty hex string",
			input:     []byte("0x"),
			expected:  []byte{},
			expectErr: false,
		},
		{
			name:      "invalid hex string",
			input:     []byte("0xZZZZ"),
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "invalid format",
			input:     []byte("invalid hex string"),
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hex.UnmarshalByteText(tt.input)
			if tt.expectErr {
				require.Error(t, err, "Test case : %s", tt.name)
			} else {
				require.NoError(t, err, "Test case : %s", tt.name)
				require.Equal(t, tt.expected, result, "Test case : %s", tt.name)
			}
		})
	}
}

func TestDecodeFixedText(t *testing.T) {
	tests := []struct {
		name      string
		typename  string
		input     []byte
		expected  []byte
		expectErr bool
	}{
		{
			name:      "valid hex string",
			typename:  "testType",
			input:     []byte("0x48656c6c6f"),
			expected:  []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expectErr: false,
		},
		{
			name:      "invalid hex string length",
			typename:  "testType",
			input:     []byte("0x48656c6c"),
			expected:  make([]byte, 5),
			expectErr: true,
		},
		{
			name:      "invalid hex characters",
			typename:  "testType",
			input:     []byte("0xZZZZZZZZZZ"),
			expected:  make([]byte, 5),
			expectErr: true,
		},
		{
			name:      "hex.Decode error",
			typename:  "testType",
			input:     []byte("0x123"), // Invalid length for hex.Decode
			expected:  make([]byte, 2),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := make([]byte, len(tt.expected))
			err := hex.DecodeFixedText(tt.input, out)
			if tt.expectErr {
				require.Error(t, err, "Test case : %s", tt.name)
			} else {
				require.NoError(t, err, "Test case : %s", tt.name)
				require.Equal(t, tt.expected, out, "Test case : %s", tt.name)
			}
		})
	}
}

func TestDecodeFixedJSON(t *testing.T) {
	tests := []struct {
		name      string
		typename  string
		input     []byte
		out       []byte
		expectErr bool
	}{
		{
			name:      "valid hex string",
			typename:  "testType",
			input:     []byte(`"0x48656c6c6f"`),
			out:       make([]byte, 5),
			expectErr: false,
		},
		{
			name:      "invalid hex string length",
			typename:  "testType",
			input:     []byte(`"0x48656c6c"`),
			out:       make([]byte, 5),
			expectErr: true,
		},
		{
			name:      "invalid hex characters",
			typename:  "testType",
			input:     []byte(`"0xZZZZZZZZZZ"`),
			out:       make([]byte, 5),
			expectErr: true,
		},
		{
			name:      "non-quoted string",
			typename:  "testType",
			input:     []byte(`0x48656c6c6f`),
			out:       make([]byte, 5),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hex.DecodeFixedJSON(
				tt.input,
				tt.out,
			)
			if tt.expectErr {
				require.Error(t, err, "Test case : %s", tt.name)
			} else {
				require.NoError(t, err, "Test case : %s", tt.name)
				require.Equal(t, []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, tt.out, "Test case : %s", tt.name)
			}
		})
	}
}

func BenchmarkDecodeFixedText(b *testing.B) {
	sizes := []int{100, 1000, 10000} // Different input sizes

	for _, size := range sizes {
		benchName := "Size" + strconv.Itoa(size)
		b.Run(benchName, func(b *testing.B) {
			input := make(
				[]byte,
				size*2+2,
			) // Each byte is represented by 2 hex characters + "0x" prefix
			input[0] = '0'
			input[1] = 'x'
			for i := 2; i < len(input); i += 2 {
				input[i] = 'a'
				input[i+1] = 'f'
			}
			out := make(
				[]byte,
				size,
			) // Adjust the size based on the expected output length

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := hex.DecodeFixedText(input, out)
				if err != nil {
					b.Fatalf("DecodeFixedText failed: %v", err)
				}
			}
		})
	}
}
