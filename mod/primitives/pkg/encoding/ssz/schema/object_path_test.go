package schema_test

import (
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/stretchr/testify/require"
)

func Test_ObjectPath(t *testing.T) {
	nested := schema.Container(
		schema.Field("bytes32", schema.B32()),
		schema.Field("uint64", schema.U64()),
		schema.Field("list_bytes32", schema.List(schema.B32(), 10)),
		schema.Field("bytes256", schema.Vector(schema.U8(), 256)),
	)
	root := schema.Container(
		schema.Field("bytes32", schema.B32()),
		schema.Field("uint32", schema.U32()),
		schema.Field("list_uint64", schema.List(schema.U64(), 1000)),
		schema.Field("list_nested", schema.List(nested, 1000)),
		schema.Field("nested", nested),
		schema.Field("vector_uint64", schema.Vector(schema.U64(), 40)),
	)

	cases := []struct {
		path   string
		gindex uint64
		offset uint8
		error  string
	}{
		// happy paths
		{path: "bytes32", gindex: 8},
		{path: "bytes32/3", gindex: 8, offset: 3},
		{path: "uint32", gindex: 9},
		{path: "list_nested", gindex: 11},
		{path: "list_nested/__len__", gindex: 2*11 + 1},
		{path: "list_nested/12", gindex: 11*2*1024 + 12},
		{path: "list_nested/12/uint64", gindex: (11*2*1024+12)*4 + 1},
		{path: "nested", gindex: 12},
		{path: "nested/uint64", gindex: 12*4 + 1},
		{path: "nested/bytes256", gindex: 12*4 + 3},
		{path: "nested/bytes256/30", gindex: (12*4 + 3) * 8, offset: 30},
		{path: "vector_uint64", gindex: 13},
		// 40 64-bit ints occupy 320 bytes (10 chunks), nextPowerOfTwo(10) = 16
		{path: "vector_uint64/5", gindex: 13*16 + (5 / 4), offset: 8},

		// error cases
		{path: "nested/__len__", error: "__len__ is only valid"},
	}
	for _, tc := range cases {
		t.Run(strings.ReplaceAll(tc.path, "/", "."), func(t *testing.T) {
			objectPath := schema.ObjectPath[[32]byte](tc.path)
			typ, gindex, offset, err := objectPath.GetGeneralizedIndex(root)
			if tc.error != "" {
				require.ErrorContains(t, err, tc.error)
				return
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, typ)
			require.NoError(t, err)
			require.Equalf(
				t,
				tc.gindex,
				uint64(gindex),
				"expected %d, got %d",
				tc.gindex,
				uint64(gindex))
			require.Equal(
				t,
				tc.offset,
				offset,
				"expected %d, got %d",
				tc.offset,
				offset,
			)
		})
	}
}
