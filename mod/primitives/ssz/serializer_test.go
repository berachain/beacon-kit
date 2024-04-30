package ssz_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/ssz"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalU64Serializer(t *testing.T) {
	original := uint64(0x0102030405060708)
	s := ssz.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU64[uint64](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U64 failed")
}

func TestMarshalUnmarshalU32Serializer(t *testing.T) {
	original := uint32(0x01020304)
	s := ssz.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU32[uint32](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U32 failed")
}

func TestMarshalUnmarshalU16Serializer(t *testing.T) {
	original := uint16(0x0102)
	s := ssz.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU16[uint16](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U16 failed")
}

func TestMarshalUnmarshalU8Serializer(t *testing.T) {
	original := uint8(0x01)
	s := ssz.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalU8[uint8](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U8 failed")
}

func TestMarshalUnmarshalBoolSerializer(t *testing.T) {
	original := true
	s := ssz.NewSerializer()
	marshaled, _ := s.MarshalSSZ(original)
	unmarshaled := ssz.UnmarshalBool[bool](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal Bool failed")
}
