package sszconstructors

import (
	"fmt"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

func NewFromSSZ[T interface {
	ssz.Object
	constraints.SSZUnmarshaler
}](buf []byte) (T, error) {
	var v T
	if v.IsUnusedFromSZZ() {
		// we special case construction of unused types, for efficiency
		if len(buf) != 1 {
			return v, fmt.Errorf("expected 1 byte got %d", len(buf))
		}
		//#nosec:G701 // UnusedType is uint8 and byte is uint8.
		tmp := types.UnusedType(buf[0])
		v, _ = any(tmp).(T) // TODO ABENEGIA: get rid of this cast
		return v, nil
	}

	if err := ssz.DecodeFromBytes(buf, v); err != nil {
		return v, err
	}
	return v, v.VerifySyntaxFromSSZ()
}
