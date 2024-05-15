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

// nolint:lll // long strings
package hex_test

import (
	"math/big"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/hex"
)

// ====================== String Invariants Testing ===========================.
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
				t.Errorf("NewStringStrict() error = %v, expectErr %v", err, test.expectErr)
			} else if err == nil {
				if !str.Has0xPrefix() {
					t.Errorf("NewStringStrict() result does not have 0x prefix: %v", str)
				}
				if str.IsEmpty() {
					t.Errorf("NewStringStrict() result is empty: %v", str)
				}
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
			if !str.Has0xPrefix() {
				t.Errorf("NewString() result does not have 0x prefix: %v", str)
			}
			if str.IsEmpty() {
				t.Errorf("NewString() result is empty: %v", str)
			}
		})
	}
}

func TestFromBytesInvariant(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "Valid byte slice",
			input: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			name:  "Empty byte slice",
			input: []byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := hex.FromBytes(test.input)
			if !str.Has0xPrefix() {
				t.Errorf("FromBytes() result does not have 0x prefix: %v", str)
			}
			if str.IsEmpty() {
				t.Errorf("FromBytes() result is empty: %v", str)
			}
		})
	}
}

func TestFromUint64Invariant(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
	}{
		{
			name:  "Zero value",
			input: 0,
		},
		{
			name:  "Positive value",
			input: 12345,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := hex.FromUint64(test.input)
			if !str.Has0xPrefix() {
				t.Errorf("FromUint64() result does not have 0x prefix: %v", str)
			}
			if str.IsEmpty() {
				t.Errorf("FromUint64() result is empty: %v", str)
			}
		})
	}
}

func TestFromBigIntInvariant(t *testing.T) {
	tests := []struct {
		name  string
		input *big.Int
	}{
		{
			name:  "Zero value",
			input: big.NewInt(0),
		},
		{
			name:  "Positive value",
			input: big.NewInt(12345),
		},
		{
			name:  "Negative value",
			input: big.NewInt(-12345),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := hex.FromBigInt(test.input)
			if !str.Has0xPrefix() {
				t.Errorf("FromBigInt() result does not have 0x prefix: %v", str)
			}
			if str.IsEmpty() {
				t.Errorf("FromBigInt() result is empty: %v", str)
			}
		})
	}
}
