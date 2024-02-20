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

package primitives_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// NOTE: SSZ is Big Endian.

func TestSSZUint64_UnmarshalSSZ_InvalidBufferLength(t *testing.T) {
	sszType := primitives.SSZUint64(0)
	var serializedObj [7]byte // Incorrect buffer size to trigger error
	err := sszType.UnmarshalSSZ(serializedObj[:])
	if err == nil {
		t.Fatal("Expected an error due to invalid buffer length, but got nil")
	}
	expectedErrMsg := "expected buffer of length"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain %q, got: %v", expectedErrMsg, err)
	}
}

func TestSSZUint64_MarshalUnmarshal_RoundTrip(t *testing.T) {
	testCases := []uint64{
		0,                    // test zero value
		8,                    // test a small number
		18446744073709551615, // test max uint64 value
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Value_%d", tc), func(t *testing.T) {
			original := primitives.SSZUint64(tc)
			marshalled, err := original.MarshalSSZ()
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			var unmarshalled primitives.SSZUint64
			if err = unmarshalled.UnmarshalSSZ(marshalled); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if original != unmarshalled {
				t.Errorf("Round-trip mismatch: original %d, after %d", original, unmarshalled)
			}
		})
	}
}

func TestSSZUint64_Serialization(t *testing.T) {
	tests := []struct {
		name            string
		value           uint64
		serializedBytes []byte
		root            []byte
	}{
		{
			name:  "test maximum value",
			value: 18446744073709551615,
			serializedBytes: hexutil.MustDecode(
				"0xffffffffffffffff"),
			root: hexutil.MustDecode(
				"0xffffffffffffffff000000000000000000000000000000000000000000000000"),
		},
		{
			name:  "random",
			value: 12345678901234567890,
			// NOTE hexutil is little endian, SSZ is big endian, so we flip hex here to compensate.
			serializedBytes: hexutil.MustDecode(
				"0xd20a1feb8ca954ab"),
			root: hexutil.MustDecode(
				"0xd20a1feb8ca954ab000000000000000000000000000000000000000000000000"),
		},
		{
			name:  "zero",
			value: 0,
			serializedBytes: hexutil.MustDecode(
				"0x0000000000000000"),
			root: hexutil.MustDecode(
				"0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := primitives.SSZUint64(tt.value)
			serializedBytes, err := s.MarshalSSZ()
			if err != nil {
				t.Fatalf("SSZUint64.MarshalSSZ() unexpected error = %v", err)
			}
			if !reflect.DeepEqual(tt.serializedBytes, serializedBytes) {
				t.Errorf("SSZUint64.MarshalSSZ() = %v, want %v",
					serializedBytes, tt.serializedBytes)
			}

			htr, err := s.HashTreeRoot()
			if err != nil {
				t.Fatalf("SSZUint64.HashTreeRoot() unexpected error = %v", err)
			}
			if !reflect.DeepEqual(tt.root, htr[:]) {
				t.Errorf("SSZUint64.HashTreeRoot() = %v, want %v", htr[:], tt.root)
			}
		})
	}
}
