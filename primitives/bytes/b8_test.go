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
	"testing"

	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/stretchr/testify/require"
)

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
			require.Equal(t, tt.want, got)
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
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
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
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToBytes8(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantRes bytes.B8
		wantErr error
	}{
		{
			name:    "Exact 8 bytes",
			input:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
			wantRes: bytes.B8{1, 2, 3, 4, 5, 6, 7, 8},
			wantErr: nil,
		},
		{
			name:    "Less than 8 bytes",
			input:   []byte{1, 2, 3, 4},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "Two bytes",
			input:   []byte{1, 2},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "Empty input",
			input:   []byte{},
			wantErr: bytes.ErrIncorrectLength,
		},
		{
			name:    "More than 8 bytes",
			input:   []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantErr: bytes.ErrIncorrectLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bytes.ToBytes8(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantRes, result)
			}
		})
	}
}

func TestBytes8MarshalSSZ(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B8
		want  []byte
	}{
		{
			name:  "marshal B8",
			input: bytes.B8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
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

func TestBytes8HashTreeRoot(t *testing.T) {
	tests := []struct {
		name  string
		input bytes.B8
		want  bytes.B32
	}{
		{
			name:  "hash tree root",
			input: bytes.B8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want:  [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
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
