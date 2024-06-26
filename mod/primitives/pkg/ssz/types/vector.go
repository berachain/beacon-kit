package types

import (
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

type SSZMarshallable interface {
	SizeSSZ() int
}

type Basic interface {
	~bool | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	MarshalSSZ() ([]byte, error)
}

type SSZVectorBasic[T Basic] []T

// SizeSSZ returns the size of the list in bytes.
func (l SSZVectorBasic[T]) SizeSSZ() int {
	elementSize := reflect.TypeOf((*T)(nil)).Elem().Size()
	fmt.Println(elementSize)
	return int(elementSize) * len(l)
}

// HashTreeRoot returns the Merkle root of the SSZVectorBasic.
func (l SSZVectorBasic[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	m := ssz.NewMerkleizer[common.Spec, any, math.U64, [32]byte]()
	packedBytes := make([]byte, l.SizeSSZ())
	for _, v := range l {
		v, err := v.MarshalSSZ()
		if err != nil {
			return [32]byte{}, err
		}

		packedBytes = append(packedBytes, v...)
	}
	return m.MerkleizeByteSlice(packedBytes)
}
