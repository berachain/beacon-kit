package ssz_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type RootContainer struct {
	Bytes32       bytes.B32
	Uint32        uint32
	ListUint64    []uint64
	ListNested    []*NestedContainer
	Nested        *NestedContainer
	VectorUint128 [][16]byte
}

type NestedContainer struct {
	Bytes32     bytes.B32
	Uint64      uint64
	ListBytes32 []bytes.B32
	Bytes256    [256]byte
}

func (n *NestedContainer) DefineSSZ(c *ssz.Codec) {
	ssz.DefineFixedVector(c, ssz.ByteVectorFromBytes(n.Bytes32[:]))
	ssz.DefineBasic(c, math.U64(n.Uint64))
	ssz.DefineList(c, ssz.ListFromElements(10, n.ListBytes32...))
	ssz.DefineFixedVector(c, ssz.ByteVectorFromBytes(n.Bytes256[:]))
}

func Test_Codec_Encode(t *testing.T) {
	nested := &NestedContainer{}
	codec := ssz.SchemaCodec()
	nested.DefineSSZ(codec)
}
