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
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
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
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !stdbytes.Equal(got, tt.want) {
				t.Errorf("FromHex() = %v, want %v", got, tt.want)
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
				if !stdbytes.Equal(result, test.expected) {
					t.Errorf("Unexpected result for input %s. Expected: %v, Got: %v", test.input, test.expected, result)
				}
			}
		})
	}
}

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
			copied := bytes.SafeCopy(tt.original)

			if !stdbytes.Equal(tt.original, copied) {
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
			copied := bytes.SafeCopy2D(tt.original)

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
			result := bytes.CopyAndReverseEndianess(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf(
					"ReverseEndianness(%v) = %v, want %v",
					tt.name, result, tt.expected)
			}
		})
	}
}

func TestPrependExtendToSize(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		length   int
		expected []byte
	}{
		{name: "Extend smaller slice",
			input:    []byte{1, 2, 3},
			length:   5,
			expected: []byte{0, 0, 1, 2, 3}},
		{name: "Extend equal size slice",
			input:    []byte{4, 5, 6},
			length:   3,
			expected: []byte{4, 5, 6}},
		{name: "Extend larger slice, to smaller size does nothing",
			input:    []byte{7, 8, 9},
			length:   2,
			expected: []byte{7, 8, 9}},
		{name: "Extend empty slice",
			input:    []byte{},
			length:   4,
			expected: []byte{0, 0, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.PrependExtendToSize(tt.input, tt.length)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf(
					"PrependExtendToSize(%v, %v) = %v, want %v",
					tt.input, tt.length, result, tt.expected)
			}
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
			name:  "valid input",
			input: `"0x01020304"`,
			want:  bytes.B4{0x01, 0x02, 0x03, 0x04},
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes4.UnmarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes4.UnmarshalJSON() = %v, want %v", got, tt.want)
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
			if got := tt.h.String(); got != tt.want {
				t.Errorf("Bytes4.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes4MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		h       bytes.B4
		want    string
		wantErr bool
	}{
		{
			name: "valid bytes",
			h:    bytes.B4{0x01, 0x02, 0x03, 0x04},
			want: "0x01020304",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes4.MarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if string(got) != tt.want {
				t.Errorf(
					"Bytes4.MarshalText() = %v, want %v",
					string(got),
					tt.want,
				)
			}
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
			name:  "valid input",
			input: "0x01020304",
			want:  bytes.B4{0x01, 0x02, 0x03, 0x04},
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes4.UnmarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes4.UnmarshalText() = %v, want %v", got, tt.want)
			}
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes32.UnmarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes32.UnmarshalText() = %v, want %v", got, tt.want)
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
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes32.UnmarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes32.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes32MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   bytes.B32
		want    string
		wantErr bool
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes32.MarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if string(got) != tt.want {
				t.Errorf(
					"Bytes32.MarshalText() = %v, want %v",
					string(got),
					tt.want,
				)
			}
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
			if got := tt.input.String(); got != tt.want {
				t.Errorf("Bytes32.String() = %v, want %v", got, tt.want)
			}
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
			if got := tt.input.String(); got != tt.want {
				t.Errorf("Bytes48.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes48MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   bytes.B48
		want    string
		wantErr bool
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes48.MarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if string(got) != tt.want {
				t.Errorf(
					"Bytes48.MarshalText() = %s, want %s",
					string(got),
					tt.want,
				)
			}
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes48.UnmarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes48.UnmarshalText() = %v, want %v", got, tt.want)
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes96.UnmarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes96.UnmarshalText() = %v, want %v", got, tt.want)
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
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes96.UnmarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes96.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestBytes96MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		h       bytes.B96
		want    string
		wantErr bool
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Bytes96.MarshalText() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if string(got) != tt.want {
				t.Errorf(
					"Bytes96.MarshalText() = %v, want %v",
					string(got),
					tt.want,
				)
			}
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
			if got := tt.h.String(); got != tt.want {
				t.Errorf("Bytes96.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
