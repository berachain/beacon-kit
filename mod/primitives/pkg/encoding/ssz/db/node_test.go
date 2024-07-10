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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package db_test

import (
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/db"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/stretchr/testify/require"
)

func Test_Schema_Paths(t *testing.T) {
	nestedType := schema.DefineContainer(
		schema.Field("bytes32", schema.DefineByteVector(32)),
		schema.Field("uint64", schema.U64()),
		schema.Field(
			"list_bytes32",
			schema.DefineList(schema.DefineByteVector(32), 10),
		),
		schema.Field("bytes256", schema.DefineByteVector(256)),
	)
	root := schema.DefineContainer(
		schema.Field("bytes32", schema.DefineByteVector(32)),
		schema.Field("uint32", schema.U32()),
		schema.Field("list_uint64", schema.DefineList(schema.U64(), 1000)),
		schema.Field("list_nested", schema.DefineList(nestedType, 1000)),
		schema.Field("nested", nestedType),
		schema.Field("vector_uint128", schema.DefineVector(schema.U128(), 20)),
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
			objectPath := merkle.ObjectPath[uint64, [32]byte](tc.path)
			node, err := db.NewTreeNode(root, objectPath)
			require.NoError(t, err)
			require.Equalf(
				t,
				tc.gindex,
				node.GIndex(),
				"expected %d, got %d",
				tc.gindex,
				node.GIndex())
			require.Equal(
				t,
				node.Offset(),
				tc.offset,
				"expected %d, got %d",
				tc.offset,
				node.Offset(),
			)
		})
	}
}

func TestNewTreeNodeEdgeCases(t *testing.T) {
	nestedType := schema.DefineContainer(
		schema.Field("uint64", schema.U64()),
		schema.Field("bytes32", schema.B32()),
		schema.Field("bytes256", schema.DefineVector(schema.U8(), 256)),
	)

	root := schema.DefineContainer(
		schema.Field("bytes32", schema.B32()),
		schema.Field("uint32", schema.U32()),
		schema.Field("list_uint64", schema.DefineList(schema.U64(), 1000)),
		schema.Field("list_nested", schema.DefineList(nestedType, 1000)),
		schema.Field("nested", nestedType),
		schema.Field("vector_uint128", schema.DefineVector(schema.U128(), 20)),
	)

	cases := []struct {
		name        string
		path        string
		expectError bool
	}{
		{name: "Invalid field", path: "nonexistent", expectError: true},
		{
			name:        "Invalid nested field",
			path:        "nested/nonexistent",
			expectError: true,
		},
		{
			name:        "Too deep nesting",
			path:        "nested/uint64/extra",
			expectError: true,
		},
		{
			name:        "Valid deep nesting",
			path:        "list_nested/5/bytes256/31",
			expectError: false,
		},
		{name: "Empty path", path: "", expectError: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			objectPath := merkle.ObjectPath[uint64, [32]byte](tc.path)
			_, err := db.NewTreeNode(root, objectPath)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
