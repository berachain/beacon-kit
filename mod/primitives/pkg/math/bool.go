package math

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

// Ensure type implements types.SSZType.
var _ types.SSZType[Bool] = (*Bool)(nil)

type Bool bool

// SizeSSZ returns the size of the bool in bytes.
func (Bool) SizeSSZ() int {
	return constants.BoolSize
}

// MarshalSSZ marshals the bool into SSZ format.
func (b Bool) MarshalSSZ() ([]byte, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// NewFromSSZ creates a new Bool from SSZ format.
func (Bool) NewFromSSZ(buf []byte) (Bool, error) {
	if len(buf) != constants.BoolSize {
		return false, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.BoolSize,
			len(buf),
		)
	}
	return Bool(buf[0] != 0), nil
}

// HashTreeRoot returns the hash tree root of the bool.
func (b Bool) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	if b {
		buf[0] = 1
	}
	return [constants.BytesPerChunk]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (Bool) IsFixed() bool {
	return true
}

// Type returns the type of the bool.
func (Bool) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the bool.
func (Bool) ChunkCount() uint64 {
	return 1
}
