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

package ssz_test

import (
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
