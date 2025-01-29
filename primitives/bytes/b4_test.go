// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"encoding/binary"
	"testing"

	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestFromUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected bytes.B4
	}{
		{
			name:     "Phase0",
			input:    version.Phase0,
			expected: bytes.B4{0, 0, 0, 0},
		},
		{
			name:     "Altair",
			input:    version.Altair,
			expected: bytes.B4{1, 0, 0, 0},
		},
		{
			name:     "Bellatrix",
			input:    version.Bellatrix,
			expected: bytes.B4{2, 0, 0, 0},
		},
		{
			name:     "Capella",
			input:    version.Capella,
			expected: bytes.B4{3, 0, 0, 0},
		},
		{
			name:     "Deneb",
			input:    version.Deneb,
			expected: bytes.B4{4, 0, 0, 0},
		},
		{
			name:     "Deneb1",
			input:    version.Deneb1,
			expected: bytes.B4{4, 1, 0, 0},
		},
		{
			name:     "Electra",
			input:    version.Electra,
			expected: bytes.B4{5, 0, 0, 0},
		},
		{
			name:     "Electra1",
			input:    version.Electra1,
			expected: bytes.B4{5, 1, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.FromUint32(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
func TestFromUint32_CustomType(t *testing.T) {
	input := uint32(123456789)
	expected := bytes.B4{}
	binary.LittleEndian.PutUint32(expected[:], input)

	result := bytes.FromUint32(input)
	require.Equal(t, expected, result)
}

func TestToUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    bytes.B4
		expected uint32
	}{
		{
			name:     "Phase0",
			input:    bytes.B4{0, 0, 0, 0},
			expected: version.Phase0,
		},
		{
			name:     "Altair",
			input:    bytes.B4{1, 0, 0, 0},
			expected: version.Altair,
		},
		{
			name:     "Bellatrix",
			input:    bytes.B4{2, 0, 0, 0},
			expected: version.Bellatrix,
		},
		{
			name:     "Capella",
			input:    bytes.B4{3, 0, 0, 0},
			expected: version.Capella,
		},
		{
			name:     "Deneb",
			input:    bytes.B4{4, 0, 0, 0},
			expected: version.Deneb,
		},
		{
			name:     "Deneb1",
			input:    bytes.B4{4, 1, 0, 0},
			expected: version.Deneb1,
		},
		{
			name:     "Electra",
			input:    bytes.B4{5, 0, 0, 0},
			expected: version.Electra,
		},
		{
			name:     "Electra1",
			input:    bytes.B4{5, 1, 0, 0},
			expected: version.Electra1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToUint32()
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestToUint32_CustomType(t *testing.T) {
	input := bytes.B4{0x15, 0xCD, 0x5B, 0x07}
	expected := uint32(123456789)

	result := input.ToUint32()
	require.Equal(t, expected, result)
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
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
			require.Equal(t, tt.want, got)
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
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestToBytes4(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantRes bytes.B4
		wantErr error
	}{
		{
			name:    "Input less than 4 bytes",
			input:   []byte{0x01, 0x02},
			wantRes: bytes.B4{},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "Input exactly 4 bytes",
			input:   []byte{0x01, 0x02, 0x03, 0x04},
			wantRes: bytes.B4{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:    "Input more than 4 bytes",
			input:   []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			wantRes: bytes.B4{},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "Empty input",
			input:   []byte{},
			wantRes: bytes.B4{},
			wantErr: bytes.ErrIncorrectLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bytes.ToBytes4(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantRes, result)
			}
		})
	}
}

func TestBytes4MarshalSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B4
		want  []byte
	}{
		{
			name:  "marshal B4",
			input: bytes.B4{0x01, 0x02, 0x03, 0x04},
			want:  []byte{0x01, 0x02, 0x03, 0x04},
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

func TestBytes4HashTreeRoot(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B4
		want  bytes.B32
	}{
		{
			name:  "hash tree root",
			input: bytes.B4{0x01, 0x02, 0x03, 0x04},
			want:  bytes.B32{0x01, 0x02, 0x03, 0x04},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
