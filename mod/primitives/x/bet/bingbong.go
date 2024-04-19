package bet

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
)

//go:generate go run github.com/ferranbt/fastssz/sszgen -path bingbong.go -include ../../math -objs ItemB

type ItemB struct {
	Val1 math.U64
	Val2 []uint64 `ssz-max:"2"`
}

type ItemA struct {
	Val1 math.U64
	Val2 []math.U64
}

func (i ItemA) SizeSSZ() int {
	return 16
}

// HashTreeRoot
func (i ItemA) HashTreeRoot() (primitives.Root, error) {
	return ssz.MerkleizeContainer[math.U64, primitives.Root, any](i)
}
