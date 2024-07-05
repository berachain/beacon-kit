package ssz_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/tree"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
	"github.com/stretchr/testify/require"
)

type TopLevel struct {
	BytesField        common.Bytes32
	UintField         math.U64
	ListOfBasic       []math.U64
	ListOfContainer   []*Nested
	Nested            *Nested
	VectorOfBasic     []math.U64
	VectorOfContainer []*Nested
}

var _ ssz.HasSchema[*TopLevel] = (*TopLevel)(nil)
var _ ssz.HasSchema[*Nested] = (*Nested)(nil)

// Default implements ssz.HasSchema.
func (t *TopLevel) Default() *TopLevel {
	if t == nil {
		return &TopLevel{}
	}
	return t
}

// Schema implements ssz.HasSchema.
func (t *TopLevel) Schema() *ssz.Schema[*TopLevel] {
	s := &ssz.Schema[*TopLevel]{}
	s.DefineField("bytes_field", func(t *TopLevel) types.MinimalSSZType {
		return t.BytesField
	})
	s.DefineField("uint_field", func(t *TopLevel) types.MinimalSSZType {
		return ssz.U64(t.UintField)
	})
	s.DefineField("list_of_basic", func(t *TopLevel) types.MinimalSSZType {
		var list []ssz.U64
		for _, v := range t.ListOfBasic {
			list = append(list, ssz.U64(v))
		}
		return ssz.ListFromElements(1000, list...)
	})
	s.DefineField("list_of_container", func(t *TopLevel) types.MinimalSSZType {
		return ssz.ListFromSchema(1000, t.ListOfContainer)
	})
	s.DefineField("nested", func(t *TopLevel) types.MinimalSSZType {
		return ssz.ContainerFromSchema(t.Nested)
	})
	s.DefineField("vector_of_basic", func(t *TopLevel) types.MinimalSSZType {
		var elements []ssz.U64
		for _, v := range t.VectorOfBasic {
			elements = append(elements, ssz.U64(v))
		}
		return ssz.VectorFromElements(elements...)
	})
	return s
}

type Nested struct {
	BytesField  common.Bytes32
	UintField   math.U64
	ListOfBytes []common.Bytes32
}

// Default implements ssz.HasSchema.
func (n *Nested) Default() *Nested {
	if n == nil {
		return &Nested{}
	}
	return n
}

// Schema implements ssz.HasSchema.
func (n *Nested) Schema() *ssz.Schema[*Nested] {
	s := &ssz.Schema[*Nested]{}
	s.DefineField("bytes_field", func(n *Nested) types.MinimalSSZType {
		return n.BytesField
	})
	s.DefineField("uint_field", func(n *Nested) types.MinimalSSZType {
		return ssz.U64(n.UintField)
	})
	s.DefineField("list_of_bytes", func(n *Nested) types.MinimalSSZType {
		return ssz.ListFromElements(1024, n.ListOfBytes...)
	})
	return s
}

func assertGIndex(
	t *testing.T,
	gIndexed tree.GIndexed,
	path tree.ObjectPath,
	expect uint64,
) {
	t.Helper()
	gi := gIndexed.GIndex(math.U64(1), path)
	require.NotNil(t, gi)
	require.Equalf(
		t,
		expect,
		uint64(gi.GIndex),
		"expected %d, got %d",
		expect,
		gi.GIndex,
	)
}

func Test_Schema(t *testing.T) {
	state := &TopLevel{}
	container := ssz.ContainerFromSchema(state)

	assertGIndex(t, container, "bytes_field", 8)
	assertGIndex(t, container, "uint_field", 9)
	assertGIndex(t, container, "list_of_container", 11)
	assertGIndex(t, container, "list_of_container/12", 11*2*1024+12)
	assertGIndex(
		t,
		container,
		"list_of_container/12/uint_field",
		(11*2*1024+12)*4+1,
	)
}
