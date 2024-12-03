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
	"testing"

	"github.com/berachain/beacon-kit/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/primitives/pkg/merkle/zero"
	"github.com/stretchr/testify/require"
)

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
			require.Equal(t, tt.want, got)
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
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestToBytes48(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantRes bytes.B48
		wantErr error
	}{
		{
			name:    "Input less than 48 bytes",
			input:   []byte{1, 2, 3},
			wantRes: bytes.B48{},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "Input exactly 48 bytes",
			input:   make([]byte, 48),
			wantRes: bytes.B48{},
		},
		{
			name:    "Input more than 48 bytes",
			input:   make([]byte, 60),
			wantRes: bytes.B48{},
			wantErr: bytes.ErrIncorrectLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bytes.ToBytes48(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantRes, result)
			}
		})
	}
}

func TestB48UnmarshalJSON(t *testing.T) {
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestB48_HashTreeRoot(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B48
		want  bytes.B32
	}{
		{
			name:  "Zero bytes",
			input: bytes.B48{},
			want:  zero.Hashes[1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.HashTreeRoot()
			require.Equal(t, tt.want, result)
		})
	}
}

func TestB48MarshalSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B48
		want  []byte
	}{
		{
			name: "valid bytes",
			input: bytes.B48{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
				0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, 0x21,
				0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
				0x2D, 0x2E, 0x2F, 0x30},
			want: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
				0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, 0x21,
				0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
				0x2D, 0x2E, 0x2F, 0x30},
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
