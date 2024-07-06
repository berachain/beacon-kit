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

type Monolith struct {
	BytesField        common.Bytes32
	UintField         math.U64
	ListOfBasic       []math.U64
	ListOfContainer   []*Nested
	Nested            *Nested
	VectorOfBasic     []math.U64
	VectorOfContainer []*Nested
}

type Nested struct {
	BytesField  common.Bytes32
	UintField   math.U64
	ListOfBytes []common.Bytes32
}

// --------------------------------------------------------------------------------
// This section shows an approach for schema builder by way of the ssz.HasSchema
// interface.  No reflection is used!
// --------------------------------------------------------------------------------

var _ ssz.HasSchema[*Monolith] = (*Monolith)(nil)
var _ ssz.HasSchema[*Nested] = (*Nested)(nil)

// Default implements ssz.HasSchema.
func (t *Monolith) Default() *Monolith {
	if t == nil {
		return &Monolith{}
	}
	return t
}

// Schema implements ssz.HasSchema.
func (t *Monolith) Schema() *ssz.Schema[*Monolith] {
	s := &ssz.Schema[*Monolith]{}
	s.DefineField("bytes_field", func(t *Monolith) types.MinimalSSZType {
		return t.BytesField
	})
	s.DefineField("uint_field", func(t *Monolith) types.MinimalSSZType {
		return ssz.U64(t.UintField)
	})
	s.DefineField("list_of_basic", func(t *Monolith) types.MinimalSSZType {
		var list []ssz.U64
		for _, v := range t.ListOfBasic {
			list = append(list, ssz.U64(v))
		}
		return ssz.ListFromElements(1000, list...)
	})
	s.DefineField("list_of_container", func(t *Monolith) types.MinimalSSZType {
		return ssz.ListFromSchema(1000, t.ListOfContainer)
	})
	s.DefineField("nested", func(t *Monolith) types.MinimalSSZType {
		return ssz.ContainerFromSchema(t.Nested)
	})
	s.DefineField("vector_of_basic", func(t *Monolith) types.MinimalSSZType {
		var elements []ssz.U64
		for _, v := range t.VectorOfBasic {
			elements = append(elements, ssz.U64(v))
		}
		return ssz.VectorFromElements(elements...)
	})
	return s
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
		return ssz.ListFromElements(10, n.ListOfBytes...)
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
	state := &Monolith{}
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

// --------------------------------------------------------------------------------
// An Active Record pattern for building containers, the dev must build the
// field order themselves.
// --------------------------------------------------------------------------------

type MonolithRecord struct {
	*ssz.Container
}

func (m *MonolithRecord) BytesField(gIndex math.U64) *tree.Node {
	return m.Container.GIndex2(gIndex, []uint64{0})
}

func (m *MonolithRecord) UintField(gIndex math.U64) *tree.Node {
	return m.Container.GIndex2(gIndex, []uint64{1})
}

func (m *MonolithRecord) ListOfBasic(
	gIndex math.U64,
) (*tree.Node, *ssz.List[ssz.U64]) {
	list := m.Container.Get(2).(*ssz.List[ssz.U64])
	return m.Container.GIndex2(
		gIndex,
		[]uint64{2},
	), list
}

func NewMonolithRecord(m *Monolith) *MonolithRecord {
	return &MonolithRecord{
		Container: ssz.ContainerFromElements(
			m.BytesField,
			ssz.U64(m.UintField),
			ssz.ListSelect(
				func(u math.U64) types.MinimalSSZType { return ssz.U64(u) },
				m.ListOfBasic,
				1000,
			),
			ssz.ListSelect(NewNestedRecord, m.ListOfContainer, 1000),
			NewNestedRecord(m.Nested),
		),
	}
}

func NewNestedRecord(n *Nested) types.MinimalSSZType {
	return ssz.ContainerFromElements(
		n.BytesField,
		ssz.U64(n.UintField),
		ssz.ListFromElements(10, n.ListOfBytes...),
	)
}

func Test_Record(t *testing.T) {
	state := &Monolith{Nested: &Nested{ListOfBytes: []common.Bytes32{}}}
	record := NewMonolithRecord(state)
	require.Equal(t, math.U64(8), record.BytesField(1).GIndex)
	require.Equal(t, math.U64(9), record.UintField(1).GIndex)
}

// --------------------------------------------------------------------------------
// This section shows an approach for schema building without the ssz.HasSchema
// interface. Containers still register fields by names but the schema is built
// by passing Key-Value pairs around where the Key is the field name and the Value
// is the SSZTyped field value.
// --------------------------------------------------------------------------------

// --------------------------------------------------------------------------------
// This section shows an approach for schema building with reflection.
// --------------------------------------------------------------------------------
