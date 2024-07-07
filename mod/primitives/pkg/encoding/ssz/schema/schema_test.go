// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package schema_test

import (
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/tree"
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
			objectPath := tree.ObjectPath(tc.path)
			node, err := schema.GetTreeNode(root, objectPath.Split())
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

type Test struct {
	MyField1 uint64
	MyField2 Nested1
}

type Nested1 struct {
	MyField3 uint64
}

func (t *Test) DefineSSZSchema() schema.SSZType {
	return schema.Container(
		schema.Field("my_field1", schema.UInt64()),
		schema.Field("nested1", Nested1{}.DefineSSZSchema()),
	)
}

func (Nested1) DefineSSZSchema() schema.SSZType {
	return schema.Container(
		schema.Field("my_field3", schema.UInt64()),
	)
}

func TestNestedSchemas(t *testing.T) {
	testSchema := (&Test{}).DefineSSZSchema()

	t.Run("Test schema", func(t *testing.T) {
		node, err := schema.GetTreeNode(
			testSchema,
			tree.ObjectPath("my_field1").Split(),
		)
		require.NoError(t, err)
		require.Equal(t, uint64(2), node.GIndex)
		require.Equal(t, uint8(0), node.Offset)
	})

	t.Run("Nested1 schema", func(t *testing.T) {
		nested1Schema := Nested1{}.DefineSSZSchema()

		node, err := schema.GetTreeNode(
			nested1Schema,
			tree.ObjectPath("my_field3").Split(),
		)
		require.NoError(t, err)
		require.Equal(t, uint64(1), node.GIndex)
		require.Equal(t, uint8(0), node.Offset)
	})

	t.Run("Nested field access", func(t *testing.T) {
		node, err := schema.GetTreeNode(
			testSchema,
			tree.ObjectPath("nested1/my_field3").Split(),
		)
		require.NoError(t, err)
		// I think this something is up here.
		require.Equal(t, uint64(6), node.GIndex)
		require.Equal(t, uint8(0), node.Offset)
	})
}
