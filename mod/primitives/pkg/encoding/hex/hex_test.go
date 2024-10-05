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

//nolint:lll // long strings
package hex_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

// ====================== Constructors ===========================.
func TestNewStringStrictInvariants(t *testing.T) {
	// NewStringStrict constructor should error if the input is invalid
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "Valid hex string",
			input:       "0x48656c6c6f",
			expectedErr: nil,
		},
		{
			name:        "Empty string",
			input:       "",
			expectedErr: hex.ErrEmptyString,
		},
		{
			name:        "No 0x prefix",
			input:       "48656c6c6f",
			expectedErr: hex.ErrMissingPrefix,
		},
		{
			name:        "Valid single hex character",
			input:       "0x0",
			expectedErr: nil,
		},
		{
			name:        "Empty hex string",
			input:       "0x",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hex.IsValidHex(tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewStringInvariants(t *testing.T) {
	// NewString constructor should never error or panic
	// output should always satisfy the string invariants regardless of input
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Valid hex string",
			input: "0x48656c6c6f",
		},
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "No 0x prefix",
			input: "48656c6c6f",
		},
		{
			name:  "Valid single hex character",
			input: "0x0",
		},
		{
			name:  "Empty hex string",
			input: "0x",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := hex.NewString(test.input)
			_, err := hex.IsValidHex(str)
			require.NoError(t, err)
		})
	}
}

// ====================== Bytes ===========================.
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
			result := hex.EncodeBytes(tt.input)
			require.Equal(t, tt.expected, result)

			decoded, err := hex.DecodeToBytes(result)
			require.NoError(t, err)
			require.Equal(t, tt.input, decoded)
		})
	}
}

// ====================== Helpers ===========================.

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
					res = hex.MustDecodeToBytes(tt.input)
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
