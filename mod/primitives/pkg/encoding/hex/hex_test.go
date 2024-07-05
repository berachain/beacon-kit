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
	"bytes"
	"encoding"
	"math/big"
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

// ====================== Constructors ===========================.
func TestNewStringStrictInvariants(t *testing.T) {
	// NewStringStrict constructor should error if the input is invalid
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "Valid hex string",
			input:     "0x48656c6c6f",
			expectErr: false,
		},
		{
			name:      "Empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "No 0x prefix",
			input:     "48656c6c6f",
			expectErr: true,
		},
		{
			name:      "Valid single hex character",
			input:     "0x0",
			expectErr: false,
		},
		{
			name:      "Empty hex string",
			input:     "0x",
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str, err := hex.NewStringStrict(test.input)
			if test.expectErr {
				require.Error(t, err, "Test case: %s", test.name)
			} else {
				require.NoError(t, err, "Test case: %s", test.name)
				verifyInvariants(t, "NewStringStrict()", str)
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
			verifyInvariants(t, "NewString()", str)
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
			result := hex.FromBytes(tt.input)

			if result.Unwrap() != tt.expected {
				t.Errorf(
					"FromBytes() = %v, want %v",
					result.Unwrap(),
					tt.expected,
				)
			}

			verifyInvariants(t, "FromBytes()", result)

			decoded, err := result.ToBytes()
			if err != nil {
				t.Errorf("ToBytes() error = %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("ToBytes() = %v, want %v", decoded, tt.input)
			}
		})
	}
}

// ====================== Numeric ===========================.

// FromUint64, then ToUint64.
func TestUint64RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "Zero value",
			input:    0,
			expected: "0x0",
		},
		{
			name:     "Positive value",
			input:    12345,
			expected: "0x3039",
		},
		{
			name:     "Max uint64 value",
			input:    ^uint64(0), // 2^64 - 1
			expected: "0xffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hex.FromUint64(tt.input)

			if result.Unwrap() != tt.expected {
				t.Errorf(
					"FromUint64() = %v, want %v",
					result.Unwrap(),
					tt.expected,
				)
			}
			verifyInvariants(t, "FromUint64()", result)
			decoded, err := strconv.ParseUint(result.Unwrap()[2:], 16, 64)
			if err != nil {
				t.Errorf("ParseUint() error = %v", err)
			}
			if decoded != tt.input {
				t.Errorf("ParseUint() = %v, want %v", decoded, tt.input)
			}
		})
	}
}

// FromBigInt, then ToBigInt.
func TestBigIntRoundTrip(t *testing.T) {
	// assume FromBigInt only called on non-negative big.Int
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{
			name:     "zero value",
			input:    big.NewInt(0),
			expected: "0x0",
		},
		{
			name:     "positive value",
			input:    big.NewInt(12345),
			expected: "0x3039",
		},
		{
			name:     "large positive value",
			input:    new(big.Int).SetBytes(bytes.Repeat([]byte{0xff}, 32)),
			expected: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hex.FromBigInt(tt.input)

			if result.Unwrap() != tt.expected {
				t.Errorf(
					"FromBigInt() = %v, want %v",
					result.Unwrap(),
					tt.expected,
				)
			}

			verifyInvariants(t, "FromBigInt()", result)

			var dec *big.Int
			var err error

			if tt.input.Sign() >= 0 {
				dec, err = hex.NewString(result.Unwrap()).ToBigInt()
			} else {
				dec, err = hex.NewString(result.Unwrap()).ToBigInt()
				dec = dec.Neg(dec)
			}

			if err != nil {
				t.Errorf("ToBigInt() error = %v", err)
			}
			if dec.Cmp(tt.input) != 0 {
				t.Errorf("ToBigInt() = %v, want %v", dec, tt.input)
			}
		})
	}
}

// ====================== Helpers ===========================.

func verifyInvariants(t *testing.T, invoker string, s hex.String) {
	t.Helper()
	if !s.Has0xPrefix() {
		t.Errorf(invoker+"result does not have 0x prefix: %v", s)
	}
	if s.IsEmpty() {
		t.Errorf(invoker+"result is empty: %v", s)
	}
}

func TestUnmarshalJSONText(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		unmarshaler encoding.TextUnmarshaler
		expectErr   bool
	}{
		{
			name:        "Valid JSON text",
			input:       []byte(`"0x48656c6c6f"`),
			unmarshaler: new(hex.String),
			expectErr:   false,
		},
		{
			name:        "Invalid JSON text",
			input:       []byte(`"invalid"`),
			unmarshaler: new(hex.String),
			expectErr:   true,
		},
		{
			name:        "Invalid quoted JSON text",
			input:       []byte(`"0x`),
			unmarshaler: new(hex.String),
			expectErr:   true,
		},
		{
			name:        "Empty JSON text",
			input:       []byte(`""`),
			unmarshaler: new(hex.String),
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hex.UnmarshalJSONText(
				tt.input,
				tt.unmarshaler,
			)
			if tt.expectErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
			}
		})
	}
}

func TestString_MustToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    hex.String
		expected []byte
		panics   bool
	}{
		{"Valid hex string", "0x68656c6c6f", []byte("hello"), false},
		{"Another valid hex string", "0x776f726c64", []byte("world"), false},
		{"Invalid hex string", "0xinvalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustToBytes did not panic on invalid input")
					}
				}()
				tt.input.MustToBytes()
			} else {
				result := tt.input.MustToBytes()
				require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
			}
		})
	}
}

func TestString_ToUint64(t *testing.T) {
	tests := []struct {
		name      string
		input     hex.String
		expected  uint64
		expectErr bool
	}{
		{"Single digit", "0x1", 1, false},
		{"Two digits", "0x10", 16, false},
		{"Mixed digits and letters", "0x1a", 26, false},
		{"Invalid hex string", "0xinvalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.ToUint64()
			if tt.expectErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
			}
		})
	}
}

func TestString_MustToUInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    hex.String
		expected uint64
		panics   bool
	}{
		{"Single digit", "0x1", 1, false},
		{"Two digits", "0x10", 16, false},
		{"Mixed digits and letters", "0x1a", 26, false},
		{"Invalid hex string", "0xinvalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustToUInt64 did not panic on invalid input")
					}
				}()
				tt.input.MustToUInt64()
			} else {
				result := tt.input.MustToUInt64()
				require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
			}
		})
	}
}

func TestString_MustToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    hex.String
		expected *big.Int
		panics   bool
	}{
		{"Valid hex string", "0x1", big.NewInt(1), false},
		{"Another valid hex string", "0x10", big.NewInt(16), false},
		{"Large valid hex string", "0x1a", big.NewInt(26), false},
		{"Invalid hex string", "0xinvalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustToBigInt did not panic on invalid input")
					}
				}()
				tt.input.MustToBigInt()
			} else {
				result := tt.input.MustToBigInt()
				require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
			}
		})
	}
}
