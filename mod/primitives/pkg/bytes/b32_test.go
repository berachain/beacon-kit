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

//nolint:lll // long strings.
package bytes_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/stretchr/testify/require"
)

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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
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
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
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
			require.Equal(t, tt.want, got)
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
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestB32MarshalSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B32
		want  []byte
	}{
		{
			name: "valid bytes",
			input: bytes.B32{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
				0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20},
			want: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
				0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
