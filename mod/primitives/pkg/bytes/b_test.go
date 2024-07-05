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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:lll // long strings.
package bytes_test

import (
	stdbytes "bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/stretchr/testify/require"
)

func TestFromHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.Bytes
		wantErr bool
	}{
		{
			name:    "Valid hex string",
			input:   "0x48656c6c6f",
			want:    bytes.Bytes{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			wantErr: false,
		},
		{
			name:    "Empty hex string",
			input:   "0x",
			want:    bytes.Bytes{},
			wantErr: false,
		},
		{
			name:    "Invalid hex string - odd length",
			input:   "0x12345",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid hex string - no 0x prefix",
			input:   "12345",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty input string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bytes.FromHex(tt.input)
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.True(t, stdbytes.Equal(got, tt.want), "Test case: %s", tt.name)
			}
		})
	}
}

func TestMustFromHex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    bytes.Bytes
		shouldPanic bool
	}{
		{
			name:        "Valid hex string",
			input:       "0x48656c6c6f",
			expected:    bytes.Bytes{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			shouldPanic: false,
		},
		{
			name:        "Empty hex string",
			input:       "0x",
			expected:    bytes.Bytes{},
			shouldPanic: false,
		},
		{
			name:        "Invalid hex string",
			input:       "0x12345",
			expected:    nil,
			shouldPanic: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf(
							"MustFromHex did not panic for input: %s",
							test.input,
						)
					}
				}()
				_ = bytes.MustFromHex(test.input)
			} else {
				result := bytes.MustFromHex(test.input)
				require.True(t, stdbytes.Equal(result, test.expected), "Test case %s", test.name)
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
			result := bytes.CopyAndReverseEndianess(tt.input)
			require.Equal(t, tt.expected, result, "Test case %s", tt.name)
		})
	}
}

func TestBytes4UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B4
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   `"0x01020304"`,
			want:    bytes.B4{0x01, 0x02, 0x03, 0x04},
			wantErr: false,
		},
		{
			name:    "invalid input - not hex",
			input:   `"01020304"`,
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   `"0x010203"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B4
			err := got.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case %s", tt.name)
			} else {
				require.NoError(t, err, "Test case %s", tt.name)
				require.Equal(t, tt.want, got, "Test case %s", tt.name)
			}
		})
	}
}

func TestBytes4String(t *testing.T) {
	tests := []struct {
		name string
		h    bytes.B4
		want string
	}{
		{
			name: "non-empty bytes",
			h:    bytes.B4{0x01, 0x02, 0x03, 0x04},
			want: "0x01020304",
		},
		{
			name: "empty bytes",
			h:    bytes.B4{},
			want: "0x00000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.String()
			require.Equal(t, tt.want, got, "Test case %s", tt.name)
		})
	}
}

func TestBytes4MarshalText(t *testing.T) {
	tests := []struct {
		name string
		h    bytes.B4
		want string
	}{
		{
			name: "valid bytes",
			h:    bytes.B4{0x01, 0x02, 0x03, 0x04},
			want: "0x01020304",
		},
		{
			name: "all zeros",
			h:    bytes.B4{0x00, 0x00, 0x00, 0x00},
			want: "0x00000000",
		},
		{
			name: "all ones",
			h:    bytes.B4{0xFF, 0xFF, 0xFF, 0xFF},
			want: "0xffffffff",
		},
		{
			name: "mixed bytes",
			h:    bytes.B4{0xAA, 0xBB, 0xCC, 0xDD},
			want: "0xaabbccdd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.MarshalText()
			require.NoError(t, err, "Test case %s", tt.name)
			require.Equal(t, tt.want, string(got), "Test case %s", tt.name)
		})
	}
}

func TestBytes4UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B4
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   "0x01020304",
			want:    bytes.B4{0x01, 0x02, 0x03, 0x04},
			wantErr: false,
		},
		{
			name:    "invalid input - not hex",
			input:   "01020304",
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   "0x010203",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B4
			err := got.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case %s", tt.name)
			} else {
				require.NoError(t, err, "Test case %s", tt.name)
				require.Equal(t, tt.want, got, "Test case %s", tt.name)
			}
		})
	}
}
func TestToBytes4(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bytes.B4
	}{
		{
			name:     "Input less than 4 bytes",
			input:    []byte{0x01, 0x02},
			expected: bytes.B4{0x01, 0x02, 0x00, 0x00},
		},
		{
			name:     "Input exactly 4 bytes",
			input:    []byte{0x01, 0x02, 0x03, 0x04},
			expected: bytes.B4{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:     "Input more than 4 bytes",
			input:    []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expected: bytes.B4{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:     "Empty input",
			input:    []byte{},
			expected: bytes.B4{0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.ToBytes4(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
func TestBytes32UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B32
		wantErr bool
	}{
		{
			name:  "valid input",
			input: "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
			want: bytes.B32{
				0x01,
				0x02,
				0x03,
				0x04,
				0x05,
				0x06,
				0x07,
				0x08,
				0x09,
				0x0a,
				0x0b,
				0x0c,
				0x0d,
				0x0e,
				0x0f,
				0x10,
				0x11,
				0x12,
				0x13,
				0x14,
				0x15,
				0x16,
				0x17,
				0x18,
				0x19,
				0x1a,
				0x1b,
				0x1c,
				0x1d,
				0x1e,
				0x1f,
				0x20,
			},
			wantErr: false,
		},
		{
			name:    "invalid input - not hex",
			input:   "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B32
			err := got.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case %s", tt.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got, "Test case %s", tt.name)
			}
		})
	}
}

func TestBytes32UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B32
		wantErr bool
	}{
		{
			name:  "valid input",
			input: `"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"`,
			want: bytes.B32{
				0x01,
				0x02,
				0x03,
				0x04,
				0x05,
				0x06,
				0x07,
				0x08,
				0x09,
				0x0a,
				0x0b,
				0x0c,
				0x0d,
				0x0e,
				0x0f,
				0x10,
				0x11,
				0x12,
				0x13,
				0x14,
				0x15,
				0x16,
				0x17,
				0x18,
				0x19,
				0x1a,
				0x1b,
				0x1c,
				0x1d,
				0x1e,
				0x1f,
				0x20,
			},
			wantErr: false,
		},
		{
			name:    "invalid input - not hex",
			input:   `"0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"`,
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   `"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"`,
			wantErr: true,
		},
		{
			name:    "invalid input - extra characters",
			input:   `"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B32
			err := got.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got, "Test case: %s", tt.name)
			}
		})
	}
}

func TestBytes32MarshalText(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B32
		want  string
	}{
		{
			name: "valid input",
			input: bytes.B32{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12,
				0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
				0x1e, 0x1f, 0x20},
			want: "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
		},
		{
			name:  "empty input",
			input: bytes.B32{},
			want:  "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalText()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, string(got), "Test case: %s", tt.name)
		})
	}
}

func TestBytes32String(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B32
		want  string
	}{
		{
			name: "valid input",
			input: bytes.B32{0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
				0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
			want: "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
		},
		{
			name:  "empty input",
			input: bytes.B32{},
			want:  "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.String()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestToBytes32(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bytes.B32
	}{
		{
			name:     "Input less than 32 bytes",
			input:    []byte{1, 2, 3},
			expected: bytes.B32{1, 2, 3},
		},
		{
			name:     "Input exactly 32 bytes",
			input:    make([]byte, 32),
			expected: bytes.B32{},
		},
		{
			name:     "Input more than 32 bytes",
			input:    make([]byte, 40),
			expected: bytes.B32{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.ToBytes32(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestHashTreeRoot(t *testing.T) {
	tests := []struct {
		name     string
		input    bytes.B32
		expected [32]byte
	}{
		{
			name:     "Non-empty input",
			input:    bytes.B32{1, 2, 3},
			expected: [32]byte{1, 2, 3},
		},
		{
			name:     "Empty input",
			input:    bytes.B32{},
			expected: [32]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.HashTreeRoot()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestBytes48String(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B48
		want  string
	}{
		{
			name: "valid input",
			input: bytes.B48{
				0x01,
				0x02,
				0x03,
				0x04,
				0x05,
				0x06,
				0x07,
				0x08,
				0x09,
				0x0a,
				0x0b,
				0x0c,
				0x0d,
				0x0e,
				0x0f,
				0x10,
				0x11,
				0x12,
				0x13,
				0x14,
				0x15,
				0x16,
				0x17,
				0x18,
				0x19,
				0x1a,
				0x1b,
				0x1c,
				0x1d,
				0x1e,
				0x1f,
				0x20,
				0x21,
				0x22,
				0x23,
				0x24,
				0x25,
				0x26,
				0x27,
				0x28,
				0x29,
				0x2a,
				0x2b,
				0x2c,
				0x2d,
				0x2e,
				0x2f,
				0x30,
			},
			want: "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30",
		},
		{
			name:  "empty input",
			input: bytes.B48{},
			want:  "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.String()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes48MarshalText(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B48
		want  string
	}{
		{
			name: "valid input",
			input: bytes.B48{
				0x01,
				0x02,
				0x03,
				0x04,
				0x05,
				0x06,
				0x07,
				0x08,
				0x09,
				0x0a,
				0x0b,
				0x0c,
				0x0d,
				0x0e,
				0x0f,
				0x10,
				0x11,
				0x12,
				0x13,
				0x14,
				0x15,
				0x16,
				0x17,
				0x18,
				0x19,
				0x1a,
				0x1b,
				0x1c,
				0x1d,
				0x1e,
				0x1f,
				0x20,
				0x21,
				0x22,
				0x23,
				0x24,
				0x25,
				0x26,
				0x27,
				0x28,
				0x29,
				0x2a,
				0x2b,
				0x2c,
				0x2d,
				0x2e,
				0x2f,
				0x30,
			},
			want: "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30",
		},
		{
			name:  "empty input",
			input: bytes.B48{},
			want:  "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalText()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, string(got), "Test case: %s", tt.name)
		})
	}
}

func TestBytes48UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B48
		wantErr bool
	}{
		{
			name:  "valid input",
			input: "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30",
			want: bytes.B48{
				0x01,
				0x02,
				0x03,
				0x04,
				0x05,
				0x06,
				0x07,
				0x08,
				0x09,
				0x0a,
				0x0b,
				0x0c,
				0x0d,
				0x0e,
				0x0f,
				0x10,
				0x11,
				0x12,
				0x13,
				0x14,
				0x15,
				0x16,
				0x17,
				0x18,
				0x19,
				0x1a,
				0x1b,
				0x1c,
				0x1d,
				0x1e,
				0x1f,
				0x20,
				0x21,
				0x22,
				0x23,
				0x24,
				0x25,
				0x26,
				0x27,
				0x28,
				0x29,
				0x2a,
				0x2b,
				0x2c,
				0x2d,
				0x2e,
				0x2f,
				0x30,
			},
		},
		{
			name:    "invalid input - not hex",
			input:   "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30",
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B48
			err := got.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.want, got, "Test case: %s", tt.name)
			}
		})
	}
}

func TestToBytes48(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bytes.B48
	}{
		{
			name:     "Input less than 48 bytes",
			input:    []byte{1, 2, 3},
			expected: bytes.B48{1, 2, 3},
		},
		{
			name:     "Input exactly 48 bytes",
			input:    make([]byte, 48),
			expected: bytes.B48{},
		},
		{
			name:     "Input more than 48 bytes",
			input:    make([]byte, 60),
			expected: bytes.B48{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.ToBytes48(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bytes.B48
		wantErr  bool
	}{
		{
			name:     "Valid input",
			input:    `"0x010203000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
			expected: bytes.B48{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "Empty input",
			input:    `"0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
			expected: bytes.B48{},
			wantErr:  false,
		},
		{
			name:     "Invalid input - not hex",
			input:    `"invalid"`,
			expected: bytes.B48{},
			wantErr:  true,
		},
		{
			name:     "Invalid input - odd length",
			input:    `"0x010203"`,
			expected: bytes.B48{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B48
			err := got.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.expected, got, "Test case: %s", tt.name)
			}
		})
	}
}

func TestBytes96UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B96
		wantErr bool
	}{
		{
			name:  "valid input",
			input: "0x" + strings.Repeat("01", 96),
			want: func() bytes.B96 {
				var b bytes.B96
				for i := range b {
					b[i] = 0x01
				}
				return b
			}(),
		},
		{
			name:    "invalid input - not hex",
			input:   strings.Repeat("01", 96),
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   "0x" + strings.Repeat("01", 95),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B96
			err := got.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.want, got, "Test case: %s", tt.name)
			}
		})
	}
}

func TestBytes96UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B96
		wantErr bool
	}{
		{
			name:  "valid input",
			input: `"0x` + strings.Repeat("01", 96) + `"`,
			want: func() bytes.B96 {
				var b bytes.B96
				for i := range b {
					b[i] = 0x01
				}
				return b
			}(),
		},
		{
			name:    "invalid input - not hex",
			input:   `"` + strings.Repeat("01", 96) + `"`,
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   `"0x` + strings.Repeat("01", 95) + `"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B96
			err := got.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.want, got, "Test case: %s", tt.name)
			}
		})
	}
}
func TestBytes96MarshalText(t *testing.T) {
	tests := []struct {
		name string
		h    bytes.B96
		want string
	}{
		{
			name: "valid bytes",
			h: func() bytes.B96 {
				var b bytes.B96
				for i := range b {
					b[i] = 0x01
				}
				return b
			}(),
			want: "0x" + strings.Repeat("01", 96),
		},
		{
			name: "empty bytes",
			h:    bytes.B96{},
			want: "0x" + strings.Repeat("00", 96),
		},
		{
			name: "mixed bytes",
			h: func() bytes.B96 {
				var b bytes.B96
				for i := 0; i < len(b); i++ {
					b[i] = byte(i % 256)
				}
				return b
			}(),
			want: "0x" + func() string {
				var s string
				for i := range 96 {
					s += fmt.Sprintf("%02x", i%256)
				}
				return s
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.MarshalText()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, string(got), "Test case: %s", tt.name)
		})
	}
}

func TestBytes96String(t *testing.T) {
	tests := []struct {
		name string
		h    bytes.B96
		want string
	}{
		{
			name: "non-empty bytes",
			h: func() bytes.B96 {
				var b bytes.B96
				for i := range b {
					b[i] = 0x01
				}
				return b
			}(),
			want: "0x" + strings.Repeat("01", 96),
		},
		{
			name: "empty bytes",
			h:    bytes.B96{},
			want: "0x" + strings.Repeat("00", 96),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.String()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestToBytes96(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bytes.B96
	}{
		{
			name:     "Input less than 96 bytes",
			input:    []byte{1, 2, 3},
			expected: bytes.B96{1, 2, 3},
		},
		{
			name:     "Input exactly 96 bytes",
			input:    make([]byte, 96),
			expected: bytes.B96{},
		},
		{
			name:     "Input more than 96 bytes",
			input:    make([]byte, 100),
			expected: bytes.B96{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.ToBytes96(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestBytes8UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B8
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   `"0x0102030405060708"`,
			want:    bytes.B8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			wantErr: false,
		},
		{
			name:    "invalid input - not hex",
			input:   `"0102030405060708"`,
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   `"0x01020304"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B8
			err := got.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBytes8String(t *testing.T) {
	tests := []struct {
		name string
		h    bytes.B8
		want string
	}{
		{
			name: "non-empty bytes",
			h:    bytes.B8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: "0x0102030405060708",
		},
		{
			name: "empty bytes",
			h:    bytes.B8{},
			want: "0x0000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.String()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes8MarshalText(t *testing.T) {
	tests := []struct {
		name string
		h    bytes.B8
		want string
	}{
		{
			name: "valid bytes",
			h:    bytes.B8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: "0x0102030405060708",
		},
		{
			name: "empty bytes",
			h:    bytes.B8{},
			want: "0x0000000000000000",
		},
		{
			name: "all zeros",
			h:    bytes.B8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: "0x0000000000000000",
		},
		{
			name: "all ones",
			h:    bytes.B8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want: "0xffffffffffffffff",
		},
		{
			name: "mixed bytes",
			h:    bytes.B8{0xaa, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22},
			want: "0xaabbccddeeff1122",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.MarshalText()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, string(got), "Test case: %s", tt.name)
		})
	}
}

func TestBytes8UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B8
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   "0x0102030405060708",
			want:    bytes.B8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			wantErr: false,
		},
		{
			name:    "invalid input - not hex",
			input:   "0102030405060708",
			wantErr: true,
		},
		{
			name:    "invalid input - wrong length",
			input:   "0x01020304",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B8
			err := got.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
			}
		})
	}
}
func TestToBytes8(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bytes.B8
	}{
		{
			name:     "Exact 8 bytes",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
			expected: bytes.B8{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "Less than 8 bytes",
			input:    []byte{1, 2, 3, 4},
			expected: bytes.B8{1, 2, 3, 4, 0, 0, 0, 0},
		},
		{
			name:     "Two bytes",
			input:    []byte{1, 2},
			expected: bytes.B8{1, 2, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "Empty input",
			input:    []byte{},
			expected: bytes.B8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "More than 8 bytes",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: bytes.B8{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.ToBytes8(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestUnmarshalFixedJSON(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		input    []byte
		out      []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Valid input",
			typ:      reflect.TypeOf([4]byte{}),
			input:    []byte(`"0x01020304"`),
			out:      make([]byte, 4),
			expected: []byte{0x01, 0x02, 0x03, 0x04},
			wantErr:  false,
		},
		{
			name:     "Invalid input - not hex",
			typ:      reflect.TypeOf([4]byte{}),
			input:    []byte(`"01020304"`),
			out:      make([]byte, 4),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid input - wrong length",
			typ:      reflect.TypeOf([4]byte{}),
			input:    []byte(`"0x010203"`),
			out:      make([]byte, 4),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bytes.UnmarshalFixedJSON(tt.input, tt.out)
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.expected, tt.out, "Test case: %s", tt.name)
			}
		})
	}
}

func TestUnmarshalFixedText(t *testing.T) {
	tests := []struct {
		name     string
		typename string
		input    []byte
		out      []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Valid input",
			typename: "B4",
			input:    []byte("0x01020304"),
			out:      make([]byte, 4),
			expected: []byte{0x01, 0x02, 0x03, 0x04},
			wantErr:  false,
		},
		{
			name:     "Invalid input - not hex",
			typename: "B4",
			input:    []byte("01020304"),
			out:      make([]byte, 4),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid input - wrong length",
			typename: "B4",
			input:    []byte("0x010203"),
			out:      make([]byte, 4),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bytes.UnmarshalFixedText(tt.input, tt.out)
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.expected, tt.out, "Test case: %s", tt.name)
			}
		})
	}
}

func TestBytes_String(t *testing.T) {
	tests := []struct {
		name     string
		input    bytes.Bytes
		expected string
	}{
		{
			name:     "Empty bytes",
			input:    bytes.Bytes{},
			expected: "0x",
		},
		{
			name:     "Single byte",
			input:    bytes.Bytes{0x01},
			expected: "0x01",
		},
		{
			name:     "Multiple bytes",
			input:    bytes.Bytes{0x01, 0x02, 0x03, 0x04},
			expected: "0x01020304",
		},
		{
			name:     "Bytes with leading zeros",
			input:    bytes.Bytes{0x00, 0x00, 0x01, 0x02},
			expected: "0x00000102",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			require.Equal(
				t,
				tt.expected,
				string(result),
				"Test case: %s",
				tt.name,
			)
		})
	}
}

func TestBytes_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   bytes.Bytes
		want    string
		wantErr bool
	}{
		{
			name:    "Empty slice",
			input:   bytes.Bytes{},
			want:    "0x",
			wantErr: false,
		},
		{
			name:    "Single byte",
			input:   bytes.Bytes{0x01},
			want:    "0x01",
			wantErr: false,
		},
		{
			name:    "Multiple bytes",
			input:   bytes.Bytes{0x01, 0x02, 0x03},
			want:    "0x010203",
			wantErr: false,
		},
		{
			name:    "Nil slice",
			input:   nil,
			want:    "0x",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalText()
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.want, string(got), "Test case: %s", tt.name)
			}
		})
	}
}
