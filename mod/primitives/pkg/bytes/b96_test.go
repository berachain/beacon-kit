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
	"fmt"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle/zero"
	"github.com/stretchr/testify/require"
)

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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
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
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
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
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToBytes96(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantRes bytes.B96
		wantErr error
	}{
		{
			name:    "Input less than 96 bytes",
			input:   []byte{1, 2, 3},
			wantRes: bytes.B96{},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "Input exactly 96 bytes",
			input:   make([]byte, 96),
			wantRes: bytes.B96{},
		},
		{
			name:    "Input more than 96 bytes",
			input:   make([]byte, 100),
			wantRes: bytes.B96{},
			wantErr: bytes.ErrIncorrectLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bytes.ToBytes96(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantRes, result)
			}
		})
	}
}

func TestB96_HashTreeRoot(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B96
		want  bytes.B32
	}{
		{
			name:  "Zero bytes",
			input: bytes.B96{},
			want:  zero.Hashes[2],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.HashTreeRoot()
			require.Equal(t, tt.want, result)
		})
	}
}

func BenchmarkB96_MarshalJSON(b *testing.B) {
	data := bytes.B96{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		require.NoError(b, err)
	}
}

func BenchmarkB96_UnmarshalJSON(b *testing.B) {
	//nolint:lll // its a test.
	jsonData := []byte(
		`"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5fdd"`,
	)
	var data bytes.B96
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := data.UnmarshalJSON(jsonData)
		require.NoError(b, err)
	}
}

func TestB96MarshalSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B96
		want  []byte
	}{
		{
			name: "valid bytes",
			input: bytes.B96{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
				0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, 0x21,
				0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
				0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
				0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40, 0x41, 0x42,
				0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D,
				0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
				0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60},
			want: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
				0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, 0x21,
				0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
				0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
				0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40, 0x41, 0x42,
				0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D,
				0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
				0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60},
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
