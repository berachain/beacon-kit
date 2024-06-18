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

package ssz_test

import (
	"reflect"
	"testing"

	ssz "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	sszv2 "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/v2/lib"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalU64Serializer(t *testing.T) {
	original := uint64(0x0102030405060708)
	s := sszv2.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU64[uint64](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U64 failed")
}

func TestMarshalUnmarshalU32Serializer(t *testing.T) {
	original := uint32(0x01020304)
	s := sszv2.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU32[uint32](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U32 failed")
}

func TestMarshalUnmarshalU16Serializer(t *testing.T) {
	original := uint16(0x0102)
	s := sszv2.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU16[uint16](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U16 failed")
}

func TestMarshalUnmarshalU8Serializer(t *testing.T) {
	original := uint8(0x01)
	s := sszv2.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU8[uint8](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U8 failed")
}

func TestMarshalUnmarshalBoolSerializer(t *testing.T) {
	original := true
	s := sszv2.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalBool[bool](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal Bool failed")
}

func TestGetArrayDimensionality(t *testing.T) {
	oneDimensionalTests := []struct {
		name     string
		input    any
		expected int
	}{
		{"1D empty array", [0]int32{}, 1},
		{"1D empty slice", []int32{}, 1},
		{"1D non-empty array", [3]int32{1, 2, 3}, 1},
		{"1D non-empty slice", []int32{1, 2, 3}, 1},
		{"1D empty byte array", [0]byte{}, 1},
		{"1D empty byte slice", []byte{}, 1},
		{"Byte array", [3]byte{1, 2, 3}, 1},
		{"Byte slice", []byte{1, 2, 3}, 1},
	}

	twoDimensionalTests := []struct {
		name     string
		input    any
		expected int
	}{
		{"2D empty array", [0][0]int32{}, 2},
		{"2D empty slice", [][]int32{}, 2},
		{"2D non-empty array", [2][3]int32{{1, 2, 3}, {4, 5, 6}}, 2},
		{"2D non-empty slice", [][]int32{{1, 2, 3}, {4, 5, 6}}, 2},
		{"2D empty bytes", [][]byte{}, 2},
		{"2D non-empty bytes", [][]byte{{1, 2, 3}, {4, 5, 6}}, 2},
	}

	threeDimensionalTests := []struct {
		name     string
		input    any
		expected int
	}{
		{"3D empty array", [0][0][0]int32{}, 3},
		{"3D empty slice", [][][]int32{}, 3},
		{
			"3D non-empty array",
			[2][2][2]int32{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			3,
		},
		{
			"3D non-empty slice",
			[][][]int32{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			3,
		},
	}

	t.Run("1D", func(t *testing.T) {
		for _, tt := range oneDimensionalTests {
			t.Run(tt.name, func(t *testing.T) {
				val := reflect.ValueOf(tt.input)
				result := sszv2.GetArrayDimensionality(val)
				if result != tt.expected {
					t.Errorf(
						"Expected dimensionality %d, but got %d",
						tt.expected,
						result,
					)
				}
			})
		}
	})

	t.Run("2D", func(t *testing.T) {
		for _, tt := range twoDimensionalTests {
			t.Run(tt.name, func(t *testing.T) {
				val := reflect.ValueOf(tt.input)
				result := sszv2.GetArrayDimensionality(val)
				if result != tt.expected {
					t.Errorf(
						"Expected dimensionality %d, but got %d",
						tt.expected,
						result,
					)
				}
			})
		}
	})

	t.Run("3D", func(t *testing.T) {
		for _, tt := range threeDimensionalTests {
			t.Run(tt.name, func(t *testing.T) {
				val := reflect.ValueOf(tt.input)
				result := sszv2.GetArrayDimensionality(val)
				if result != tt.expected {
					t.Errorf(
						"Expected dimensionality %d, but got %d",
						tt.expected,
						result,
					)
				}
			})
		}
	})
}
