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
var _ types.SSZType[U32] = (*U32)(nil)

// U32 represents a 32-bit unsigned integer that is both SSZ and JSON
type U32 uint32

/* -------------------------------------------------------------------------- */
/*                                     U8                                     */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                     U32                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint32 in bytes.
func (U32) SizeSSZ() int {
	return constants.U32Size
}

// MarshalSSZ marshals the uint32 into SSZ format.
func (u U32) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U32Size)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf, nil
}

// NewFromSSZ creates a new U32 from SSZ format.
func (U32) NewFromSSZ(buf []byte) (U32, error) {
	if len(buf) != constants.U32Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U32Size,
			len(buf),
		)
	}
	return U32(binary.LittleEndian.Uint32(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint32.
func (u U32) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint32(buf[:constants.U32Size], uint32(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U32) IsFixed() bool {
	return true
}

// Type returns the type of the U32.
func (U32) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the uint32.
func (U32) ChunkCount() uint64 {
	return 1
}
