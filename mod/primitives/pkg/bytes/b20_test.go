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

package bytes_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/stretchr/testify/require"
)

func TestBytes20MarshalText(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B20
		want  string
	}{
		{
			name: "valid bytes",
			input: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
			want: "0x0102030405060708090a0b0c0d0e0f1011121314",
		},
		{
			name: "all zeros",
			input: bytes.B20{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: "0x0000000000000000000000000000000000000000",
		},
		{
			name: "all ones",
			input: bytes.B20{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want: "0xffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name: "mixed bytes",
			input: bytes.B20{
				0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55,
				0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE},
			want: "0xaabbccddeeff112233445566778899aabbccddee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalText()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, string(got),
				"Test case: %s", tt.name)
		})
	}
}

func TestBytes20SizeSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B20
		want  int
	}{
		{
			name: "size of B20",
			input: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
			want: bytes.B20Size,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.SizeSSZ()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes20MarshalSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B20
		want  []byte
	}{
		{
			name: "marshal B20",
			input: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
			want: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalSSZ()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes20IsFixed(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B20
		want  bool
	}{
		{
			name: "is fixed",
			input: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.IsFixed()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes20Type(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B20
		want  schema.SSZType
	}{
		{
			name: "type of B20",
			input: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
			want: schema.B20(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Type()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes20HashTreeRoot(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B20
		want  [32]byte
	}{
		{
			name: "hash tree root",
			input: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
			want: [32]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.HashTreeRoot()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestBytes20UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B20
		wantErr bool
	}{
		{
			name:  "valid hex",
			input: "0x0102030405060708090a0b0c0d0e0f1011121314",
			want: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
			wantErr: false,
		},
		{
			name:    "invalid hex",
			input:   "0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
			want:    bytes.B20{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B20
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

func TestBytes20UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bytes.B20
		wantErr bool
	}{
		{
			name:  "valid JSON",
			input: "\"0x0102030405060708090a0b0c0d0e0f1011121314\"",
			want: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   "\"0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ\"",
			want:    bytes.B20{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.B20
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

func TestToBytes20(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bytes.B20
	}{
		{
			name: "exact 20 bytes",
			input: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
				0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
			want: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
			},
		},
		{
			name:  "less than 20 bytes",
			input: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			want: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name:  "less than 20 bytes",
			input: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			want: bytes.B20{
				0x01,
				0x02,
				0x03,
				0x04,
				0x05,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
				0x00,
			},
		},
		{
			name: "more than 20 bytes",
			input: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
			want: bytes.B20{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bytes.ToBytes20(tt.input)
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}
