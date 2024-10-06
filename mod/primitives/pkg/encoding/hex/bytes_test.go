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
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestEncodeBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "typical byte slice",
			input:    []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expected: "0x48656c6c6f",
		},
		{
			name:     "empty byte slice",
			input:    []byte{},
			expected: "0x",
		},
		{
			name:     "single byte",
			input:    []byte{0x01},
			expected: "0x01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hex.FromBytes(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "typical byte slice",
			input:    []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expected: "0x48656c6c6f",
		},
		{
			name:     "empty byte slice",
			input:    []byte{},
			expected: "0x",
		},
		{
			name:     "single byte",
			input:    []byte{0x01},
			expected: "0x01",
		},
		{
			name: "long byte slice",
			input: []byte{
				0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad,
				0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef,
				0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe,
				0xba, 0xbe},
			expected: "0xdeadbeefcafebabe" + "deadbeefcafebabe" + "deadbeefcafebabe" + "deadbeefcafebabe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hex.FromBytes(tt.input)
			require.Equal(t, tt.expected, result)

			decoded, err := hex.ToBytes(result)
			require.NoError(t, err)
			require.Equal(t, tt.input, decoded)
		})
	}
}

func TestString_MustToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		panics   bool
	}{
		{"Valid hex string", "0x68656c6c6f", []byte("hello"), false},
		{"Another valid hex string", "0x776f726c64", []byte("world"), false},
		{"Invalid hex string", "0xinvalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				res []byte
				f   = func() {
					res = hex.MustToBytes(tt.input)
				}
			)
			if tt.panics {
				require.Panics(t, f, "Test case: %s", tt.name)
			} else {
				require.NotPanics(t, f)
				require.Equal(t, tt.expected, res, "Test case: %s", tt.name)
			}
		})
	}
}
