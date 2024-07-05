package math

import (
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Ensure types implement types.SSZType.
var _ types.SSZType[U16] = (*U16)(nil)

// U16 represents a 16-bit unsigned integer that is both SSZ and JSON
type U16 uint16

/* -------------------------------------------------------------------------- */
/*                                     U16                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint16 in bytes.
func (U16) SizeSSZ() int {
	return constants.U16Size
}

// MarshalSSZ marshals the uint16 into SSZ format.
func (u U16) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U16Size)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf, nil
}

// NewFromSSZ creates a new U16 from SSZ format.
func (U16) NewFromSSZ(buf []byte) (U16, error) {
	if len(buf) != constants.U16Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U16Size,
			len(buf),
		)
	}
	return U16(binary.LittleEndian.Uint16(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint16.
func (u U16) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint16(buf[:constants.U16Size], uint16(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U16) IsFixed() bool {
	return true
}

// Type returns the type of the U16.
func (U16) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the uint16.
func (U16) ChunkCount() uint64 {
	return 1
}
