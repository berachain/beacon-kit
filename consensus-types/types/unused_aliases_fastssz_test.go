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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

// TestUnusedTypeAliasesFastSSZ tests that all UnusedType aliases support fastssz methods.
func TestUnusedTypeAliasesFastSSZ(t *testing.T) {
	t.Parallel()

	// Test all the UnusedType aliases
	aliases := []struct {
		name   string
		create func() interface {
			MarshalSSZTo([]byte) ([]byte, error)
			UnmarshalSSZ([]byte) error
			SizeSSZ() int
			HashTreeRootWith(ssz.HashWalker) error
			GetTree() (*ssz.Node, error)
			HashTreeRoot() ([32]byte, error)
		}
	}{
		{
			name: "ProposerSlashing",
			create: func() interface {
				MarshalSSZTo([]byte) ([]byte, error)
				UnmarshalSSZ([]byte) error
				SizeSSZ() int
				HashTreeRootWith(ssz.HashWalker) error
				GetTree() (*ssz.Node, error)
				HashTreeRoot() ([32]byte, error)
			} {
				v := types.ProposerSlashing(0)
				return &v
			},
		},
		{
			name: "AttesterSlashing",
			create: func() interface {
				MarshalSSZTo([]byte) ([]byte, error)
				UnmarshalSSZ([]byte) error
				SizeSSZ() int
				HashTreeRootWith(ssz.HashWalker) error
				GetTree() (*ssz.Node, error)
				HashTreeRoot() ([32]byte, error)
			} {
				v := types.AttesterSlashing(0)
				return &v
			},
		},
		{
			name: "VoluntaryExit",
			create: func() interface {
				MarshalSSZTo([]byte) ([]byte, error)
				UnmarshalSSZ([]byte) error
				SizeSSZ() int
				HashTreeRootWith(ssz.HashWalker) error
				GetTree() (*ssz.Node, error)
				HashTreeRoot() ([32]byte, error)
			} {
				v := types.VoluntaryExit(0)
				return &v
			},
		},
		{
			name: "BlsToExecutionChange",
			create: func() interface {
				MarshalSSZTo([]byte) ([]byte, error)
				UnmarshalSSZ([]byte) error
				SizeSSZ() int
				HashTreeRootWith(ssz.HashWalker) error
				GetTree() (*ssz.Node, error)
				HashTreeRoot() ([32]byte, error)
			} {
				v := types.BlsToExecutionChange(0)
				return &v
			},
		},
		{
			name: "Attestation",
			create: func() interface {
				MarshalSSZTo([]byte) ([]byte, error)
				UnmarshalSSZ([]byte) error
				SizeSSZ() int
				HashTreeRootWith(ssz.HashWalker) error
				GetTree() (*ssz.Node, error)
				HashTreeRoot() ([32]byte, error)
			} {
				v := types.Attestation(0)
				return &v
			},
		},
	}

	for _, tt := range aliases {
		t.Run(tt.name, func(t *testing.T) {
			obj := tt.create()

			// Test MarshalSSZTo
			dst := make([]byte, 0)
			result, err := obj.MarshalSSZTo(dst)
			require.NoError(t, err)
			require.Equal(t, []byte{0}, result)

			// Test UnmarshalSSZ
			err = obj.UnmarshalSSZ([]byte{0})
			require.NoError(t, err)

			// Test SizeSSZ
			size := obj.SizeSSZ()
			require.Equal(t, 1, size)

			// Test HashTreeRootWith
			hh := ssz.NewHasher()
			err = obj.HashTreeRootWith(hh)
			require.NoError(t, err)

			// Test GetTree
			tree, err := obj.GetTree()
			require.NoError(t, err)
			require.NotNil(t, tree)

			// Compare hash tree roots
			sszRoot, err := obj.HashTreeRoot()
			require.NoError(t, err)
			hh = ssz.NewHasher()
			err = obj.HashTreeRootWith(hh)
			require.NoError(t, err)
			fastsszRoot, err := hh.HashRoot()
			require.NoError(t, err)
			require.Equal(t, sszRoot[:], fastsszRoot[:],
				"HashTreeRoot results should match between ssz and fastssz for %s", tt.name)
		})
	}
}

// TestUnusedTypeAliasesEnforcement tests that UnusedType aliases enforce zero values.
func TestUnusedTypeAliasesEnforcement(t *testing.T) {
	t.Parallel()

	// Test that unmarshaling non-zero values fails
	testCases := []struct {
		name      string
		unmarshal func([]byte) error
	}{
		{
			name: "ProposerSlashing",
			unmarshal: func(buf []byte) error {
				var v types.ProposerSlashing
				return v.UnmarshalSSZ(buf)
			},
		},
		{
			name: "AttesterSlashing",
			unmarshal: func(buf []byte) error {
				var v types.AttesterSlashing
				return v.UnmarshalSSZ(buf)
			},
		},
		{
			name: "VoluntaryExit",
			unmarshal: func(buf []byte) error {
				var v types.VoluntaryExit
				return v.UnmarshalSSZ(buf)
			},
		},
		{
			name: "BlsToExecutionChange",
			unmarshal: func(buf []byte) error {
				var v types.BlsToExecutionChange
				return v.UnmarshalSSZ(buf)
			},
		},
		{
			name: "Attestation",
			unmarshal: func(buf []byte) error {
				var v types.Attestation
				return v.UnmarshalSSZ(buf)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshaling non-zero value
			err := tc.unmarshal([]byte{1})
			require.Error(t, err)
			require.Contains(t, err.Error(), "UnusedType must be unused")

			// Test unmarshaling zero value succeeds
			err = tc.unmarshal([]byte{0})
			require.NoError(t, err)
		})
	}
}
