package types

import "encoding/binary"

type SSZBool bool

// SizeSSZ returns the size of the bool in bytes.
func (b SSZBool) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the bool into SSZ format.
func (b SSZBool) MarshalSSZ() ([]byte, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// -----------------------------

// -----------------------------

type SSZUInt8 uint8

// SizeSSZ returns the size of the uint8 in bytes.
func (u SSZUInt8) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the uint8 into SSZ format.
func (u SSZUInt8) MarshalSSZ() ([]byte, error) {
	return []byte{byte(u)}, nil
}

// -----------------------------

type SSZUInt16 uint16

// SizeSSZ returns the size of the uint16 in bytes.
func (u SSZUInt16) SizeSSZ() int {
	return 2
}

// MarshalSSZ marshals the uint16 into SSZ format.
func (u SSZUInt16) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf, nil
}

// -----------------------------

type SSZUInt32 uint32

// SizeSSZ returns the size of the uint32 in bytes.
func (u SSZUInt32) SizeSSZ() int {
	return 4
}

// MarshalSSZ marshals the uint32 into SSZ format.
func (u SSZUInt32) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf, nil
}

// -----------------------------

type SSZUInt64 uint64

// SizeSSZ returns the size of the uint64 in bytes.
func (u SSZUInt64) SizeSSZ() int {
	return 8
}

// MarshalSSZ marshals the uint64 into SSZ format.
func (u SSZUInt64) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// -----------------------------

type SSZByte byte

// SizeSSZ returns the size of the byte slice in bytes.
func (b SSZByte) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the byte into SSZ format.
func (b SSZByte) MarshalSSZ() ([]byte, error) {
	return []byte{byte(b)}, nil
}
