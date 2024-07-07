package schema

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

const chunkSize = constants.BytesPerChunk

type SSZType interface {
	ID() types.Type
	ItemLength() uint64
	Chunks() uint64

	position(p string) (uint64, uint8, error)
	child(p string) SSZType
}

// Container Type

// Enumerable Type (vectors and lists)

type Node struct {
	SSZType

	GIndex uint64
	Offset uint8
}

// GetTreeNode locates a node in the SSZ merkle tree by its path and a root
// schema node to begin traversal from with gindex 1.
//
//nolint:mnd // binary math
func GetTreeNode(typ SSZType, path []string) (Node, error) {
	var (
		gindex = uint64(1)
		offset uint8
	)
	for _, p := range path {
		if p == "__len__" {
			if _, ok := typ.(list); !ok {
				if _, ok := typ.(vector); !ok {
					return Node{}, fmt.Errorf("type %T is not a list or vector", typ)
				}
			}
			gindex = 2*gindex + 1
			offset = 0
		} else {
			pos, off, err := typ.position(p)
			if err != nil {
				return Node{}, err
			}
			i := uint64(1)
			if l, ok := typ.(list); ok && l.IsList() {
				i = 2
			}
			gindex = gindex*i*nextPowerOfTwo(typ.Chunks()) + pos
			typ = typ.child(p)
			offset = off
		}
	}
	return Node{SSZType: typ, GIndex: gindex, Offset: offset}, nil
}

//nolint:mnd // binary math
func nextPowerOfTwo(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
