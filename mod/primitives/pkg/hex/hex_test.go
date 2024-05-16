// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

//nolint:lll // long strings
package hex_test

import (
	"bytes"
	"math/big"
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/hex"
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
			if (err != nil) != test.expectErr {
				t.Errorf(
					"NewStringStrict() error = %v, expectErr %v",
					err,
					test.expectErr,
				)
			} else if err == nil {
				verifyInvariants("NewStringStrict()", t, str)
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
			verifyInvariants("NewString()", t, str)
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

			verifyInvariants("FromBytes()", t, result)

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
			verifyInvariants("FromUint64()", t, result)
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

			verifyInvariants("FromBigInt()", t, result)

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

func verifyInvariants(invoker string, t *testing.T, s hex.String) {
	if !s.Has0xPrefix() {
		t.Errorf(invoker+"result does not have 0x prefix: %v", s)
	}
	if s.IsEmpty() {
		t.Errorf(invoker+"result is empty: %v", s)
	}
}
