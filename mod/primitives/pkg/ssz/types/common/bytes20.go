package common

import (
	"encoding/hex"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

var _ types.MinimalSSZType = Bytes20{}
var _ types.SSZType[Bytes20] = Bytes20{}

// Bytes20 represents a 20-byte array.
type Bytes20 [20]byte

// NewBytes20 creates a new Bytes20 from a byte slice.
func NewBytes20(b []byte) Bytes20 {
	var arr Bytes20
	copy(arr[:], b)
	return arr
}

// Bytes returns the byte slice representation of the Bytes20.
func (b Bytes20) Bytes() []byte {
	return b[:]
}

// String returns the hexadecimal string representation of the Bytes20.
func (b Bytes20) String() string {
	return hex.EncodeToString(b[:])
}

// SizeSSZ returns the SSZ encoded size of the Bytes20.
func (b Bytes20) SizeSSZ() int {
	return 20
}

// MarshalSSZ marshals the Bytes20 to SSZ encoded bytes.
func (b Bytes20) MarshalSSZ() ([]byte, error) {
	return b[:], nil
}

// UnmarshalSSZ unmarshals the Bytes20 from SSZ encoded bytes.
func (b Bytes20) UnmarshalSSZ(buf []byte) error {
	copy(b[:], buf)
	return nil
}

// HashTreeRoot returns the hash tree root of the Bytes20.
func (b Bytes20) HashTreeRoot() ([32]byte, error) {
	var root [32]byte
	copy(root[:], b[:])
	return root, nil
}

func (b Bytes20) Type() types.Type {
	return types.Composite
}

// IsFixed returns whether Bytes20 is a fixed-length type.
func (b Bytes20) IsFixed() bool {
	return true
}

// ChunkCount returns the number of chunks in the Bytes20.
func (b Bytes20) ChunkCount() uint64 {
	return 1
}

// NewFromSSZ creates a new Bytes20 from SSZ encoded bytes.
func (b Bytes20) NewFromSSZ(buf []byte) (Bytes20, error) {
	if err := b.UnmarshalSSZ(buf); err != nil {
		return Bytes20{}, err
	}
	return b, nil
}
