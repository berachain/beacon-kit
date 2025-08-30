// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package common_test

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

// Verify that UnmarshalSSZ properly enforces the UnusedType constraint
func TestDecodeUnusedTypeEquality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
		errMsg  string
	}{
		{name: "decode-unused-type-empty", buf: []byte{0x00}, wantErr: false},
		{name: "decode-unused-type-one", buf: []byte{0x01}, wantErr: true, errMsg: "UnusedType must be unused"},
		{name: "decode-unused-type-max", buf: []byte{0xff}, wantErr: true, errMsg: "UnusedType must be unused"},
		{name: "decode-unused-type-too-long", buf: []byte{0xff, 0xff}, wantErr: true, errMsg: "expected buffer of length 1"},
		{name: "decode-unused-type-too-short", buf: []byte{}, wantErr: true, errMsg: "expected buffer of length 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(common.UnusedType)
			err := got.UnmarshalSSZ(tt.buf)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				want := common.UnusedType(0)
				require.Equal(t, &want, got)
			}
		})
	}
}

// Verify that MarshalSSZ produces the same bytes as the previous implementation
// defined by:
// []byte{uint8(*ut)}
func TestEncodeUnusedTypeEquality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ut   common.UnusedType
	}{
		{name: "encode-unused-type-empty", ut: common.UnusedType(0)},
		{name: "encode-unused-type-one", ut: common.UnusedType(1)},
		{name: "encode-unused-type-max", ut: ^common.UnusedType(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ut.MarshalSSZ()
			if err != nil {
				t.Errorf("MarshalSSZ() error = %v", err)
				return
			}
			want := []byte{uint8(tt.ut)}
			require.Equal(t, want, got)
		})
	}
}

// TestUnusedTypeFastSSZ tests the fastssz methods of UnusedType.
func TestUnusedTypeFastSSZ(t *testing.T) {
	t.Run("ValidUnusedType", func(t *testing.T) {
		ut := common.UnusedType(0)

		// Test MarshalSSZTo
		dst := make([]byte, 0)
		result, err := ut.MarshalSSZTo(dst)
		require.NoError(t, err)
		require.Equal(t, []byte{0}, result)

		// Test UnmarshalSSZ
		var ut2 common.UnusedType
		err = ut2.UnmarshalSSZ([]byte{0})
		require.NoError(t, err)
		require.Equal(t, ut, ut2)

		// Test SizeSSZ
		size := ut.SizeSSZ()
		require.Equal(t, 1, size)

		// Test HashTreeRootWith
		hh := fastssz.NewHasher()
		err = ut.HashTreeRootWith(hh)
		require.NoError(t, err)

		// Test GetTree
		tree, err := ut.GetTree()
		require.NoError(t, err)
		require.NotNil(t, tree)

		// Test HashTreeRoot
		root, err := ut.HashTreeRoot()
		require.NoError(t, err)
		var expectedRoot [32]byte
		require.Equal(t, expectedRoot, root)
	})

	t.Run("InvalidUnusedType", func(t *testing.T) {
		// Test unmarshaling non-zero value
		var ut common.UnusedType
		err := ut.UnmarshalSSZ([]byte{1})
		require.Error(t, err)
		require.Contains(t, err.Error(), "UnusedType must be unused")

		// Test unmarshaling wrong size
		err = ut.UnmarshalSSZ([]byte{0, 0})
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected buffer of length 1")

		// Test empty buffer
		err = ut.UnmarshalSSZ([]byte{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected buffer of length 1")
	})

	t.Run("CompareHashTreeRoot", func(t *testing.T) {
		ut := common.UnusedType(0)

		// Compare HashTreeRoot with fastssz HashTreeRootWith
		sszRoot, err := ut.HashTreeRoot()
		require.NoError(t, err)

		hh := fastssz.NewHasher()
		err = ut.HashTreeRootWith(hh)
		require.NoError(t, err)
		fastsszRoot, err := hh.HashRoot()
		require.NoError(t, err)

		require.Equal(t, sszRoot[:], fastsszRoot[:],
			"HashTreeRoot results should match between ssz and fastssz")

		// Also compare with HashTreeRoot
		directRoot, err := ut.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, sszRoot[:], directRoot[:],
			"HashTreeRoot results should match with direct fastssz method")
	})
}

// TestEnforceAllUnused tests the EnforceAllUnused helper function.
func TestEnforceAllUnused(t *testing.T) {
	t.Run("AllUnused", func(t *testing.T) {
		ut1 := common.UnusedType(0)
		ut2 := common.UnusedType(0)
		err := common.EnforceAllUnused(&ut1, &ut2)
		require.NoError(t, err)
	})

	t.Run("SomeUsed", func(t *testing.T) {
		ut1 := common.UnusedType(0)
		ut2 := common.UnusedType(1)
		err := common.EnforceAllUnused(&ut1, &ut2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "UnusedType must be unused")
	})
}
