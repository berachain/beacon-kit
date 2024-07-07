package schema_test

import (
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/stretchr/testify/require"
)

func Test_Schema_Paths(t *testing.T) {
	nestedType := schema.Container(
		schema.Field("bytes32", schema.Bytes(32)),
		schema.Field("uint64", schema.UInt64()),
		schema.Field("list_bytes32", schema.List(schema.Bytes(32), 10)),
		schema.Field("bytes256", schema.Bytes(256)),
	)
	root := schema.Container(
		schema.Field("bytes32", schema.Bytes(32)),
		schema.Field("uint32", schema.UInt32()),
		schema.Field("list_uint64", schema.List(schema.UInt64(), 1000)),
		schema.Field("list_nested", schema.List(nestedType, 1000)),
		schema.Field("nested", nestedType),
		schema.Field("vector_uint128", schema.Vector(schema.UInt128(), 20)),
	)

	cases := []struct {
		path   string
		gindex uint64
		offset uint8
	}{
		{path: "bytes32", gindex: 8},
		{path: "bytes32/3", gindex: 8, offset: 3},
		{path: "uint32", gindex: 9},
		{path: "list_nested", gindex: 11},
		{path: "list_nested/12", gindex: 11*2*1024 + 12},
		{path: "list_nested/12/uint64", gindex: (11*2*1024+12)*4 + 1},
		{path: "nested", gindex: 12},
		{path: "nested/uint64", gindex: 12*4 + 1},
		{path: "nested/bytes256", gindex: 12*4 + 3},
		{path: "nested/bytes256/30", gindex: (12*4 + 3) * 8, offset: 30},
		{path: "vector_uint128", gindex: 13},
		// 20 128-bit ints occupy 320 bytes (10 chunks), nextPowerOfTwo(10) = 16
		{path: "vector_uint128/5", gindex: 13*16 + (5 / 2), offset: 16},
	}
	for _, tc := range cases {
		t.Run(strings.ReplaceAll(tc.path, "/", "."), func(t *testing.T) {
			objectPath := schema.Path(tc.path)
			node, err := schema.GetTreeNode(root, objectPath)
			require.NoError(t, err)
			require.Equalf(
				t,
				tc.gindex,
				node.GIndex,
				"expected %d, got %d",
				tc.gindex,
				node.GIndex)
			require.Equal(
				t,
				node.Offset,
				tc.offset,
				"expected %d, got %d",
				tc.offset,
				node.Offset,
			)
		})
	}
}
