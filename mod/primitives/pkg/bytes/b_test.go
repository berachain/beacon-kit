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

package bytes_test

import (
	stdhex "encoding/hex"
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestFromHex(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput bytes.Bytes
		wantErr    error
	}{
		{
			name:       "Valid hex string",
			input:      "0x48656c6c6f",
			wantOutput: bytes.Bytes{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			wantErr:    nil,
		},
		{
			name:       "Empty hex string",
			input:      "0x",
			wantOutput: bytes.Bytes{},
			wantErr:    nil,
		},
		{
			name:       "Invalid hex string - odd length",
			input:      "0x12345",
			wantOutput: nil,
			wantErr:    stdhex.ErrLength,
		},
		{
			name:       "Invalid hex string - no 0x prefix",
			input:      "12345",
			wantOutput: nil,
			wantErr:    hex.ErrMissingPrefix,
		},
		{
			name:       "Empty input string",
			input:      "",
			wantOutput: nil,
			wantErr:    hex.ErrEmptyString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hex.ToBytes(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantOutput, bytes.Bytes(got))
			}
		})
	}
}

func TestToBytesSafe(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []byte
		shouldPanic bool
	}{
		{
			name:        "Valid hex string",
			input:       "0x68656c6c6f",
			expected:    bytes.Bytes("hello"),
			shouldPanic: false,
		},
		{
			name:        "Another valid hex string",
			input:       "0x776f726c64",
			expected:    bytes.Bytes("world"),
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
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				res []byte
				f   = func() {
					res = hex.ToBytesSafe(tt.input)
				}
			)
			if tt.shouldPanic {
				require.Panics(t, f)
			} else {
				require.NotPanics(t, f)
				require.Equal(t, tt.expected, res)
			}
		})
	}
}

func TestBytesUnmarshalJSONText(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expectErr bool
	}{
		{
			name:      "Valid JSON text",
			input:     []byte(`"0x48656c6c6f"`),
			expectErr: false,
		},
		{
			name:      "Invalid JSON text",
			input:     []byte(`"invalid"`),
			expectErr: true,
		},
		{
			name:      "Invalid quoted JSON text",
			input:     []byte(`"0x`),
			expectErr: true,
		},
		{
			name:      "Empty JSON text",
			input:     []byte(`""`),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Bytes{}
			err := b.UnmarshalJSON(tt.input)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHashTreeRoot(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B32
		want  bytes.B32
	}{
		{
			name:  "Non-empty input",
			input: bytes.B32{1, 2, 3},
			want:  [32]byte{1, 2, 3},
		},
		{
			name:  "Empty input",
			input: bytes.B32{},
			want:  [32]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.HashTreeRoot()
			require.Equal(t, tt.want, result)
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, tt.out)
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, tt.out)
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
			require.Equal(t, tt.expected, result)
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, string(got))
			}
		})
	}
}
