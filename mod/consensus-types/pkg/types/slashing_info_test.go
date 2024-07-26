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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func generateSlashingInfo() *types.SlashingInfo {
	s := &types.SlashingInfo{}
	return s.New(12345, 67890)
}

func TestSlashingInfo_MarshalSSZ_UnmarshalSSZ(t *testing.T) {
	testCases := []struct {
		name     string
		data     *types.SlashingInfo
		expected *types.SlashingInfo
		err      error
	}{
		{
			name:     "Valid SlashingInfo",
			data:     generateSlashingInfo(),
			expected: generateSlashingInfo(),
			err:      nil,
		},
		{
			name: "Empty SlashingInfo",
			data: &types.SlashingInfo{
				Slot:  0,
				Index: 0,
			},
			expected: &types.SlashingInfo{
				Slot:  0,
				Index: 0,
			},
			err: nil,
		},
		{
			name:     "Invalid Buffer Size",
			data:     generateSlashingInfo(),
			expected: nil,
			err:      io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.data.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			var unmarshalled types.SlashingInfo
			if tc.name == "Invalid Buffer Size" {
				err = unmarshalled.UnmarshalSSZ(data[:8])
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				err = unmarshalled.UnmarshalSSZ(data)
				require.NoError(t, err)
				require.Equal(t, tc.expected, &unmarshalled)
			}
		})
	}
}

func TestSlashingInfo_GetTree(t *testing.T) {
	data := generateSlashingInfo()

	tree, err := data.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)

	expectedRoot, err := data.HashTreeRoot()
	require.NoError(t, err)

	// Compare the tree root with the expected root
	actualRoot := tree.Hash()
	require.Equal(t, string(expectedRoot[:]), string(actualRoot))
}

func TestSlashingInfo_MarshalSSZTo(t *testing.T) {
	testCases := []struct {
		name     string
		data     *types.SlashingInfo
		expected []byte
		err      error
	}{
		{
			name: "Valid SlashingInfo",
			data: generateSlashingInfo(),
			expected: func() []byte {
				bz, _ := generateSlashingInfo().MarshalSSZ()
				return bz
			}(),
			err: nil,
		},
		{
			name: "Empty SlashingInfo",
			data: &types.SlashingInfo{
				Slot:  0,
				Index: 0,
			},
			expected: func() []byte {
				bz, _ := (&types.SlashingInfo{Slot: 0, Index: 0}).MarshalSSZ()
				return bz
			}(),
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dst := make([]byte, 0)
			result, err := tc.data.MarshalSSZTo(dst)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestSlashingInfo_GettersSetters(t *testing.T) {
	// Create a new SlashingInfo instance
	si := &types.SlashingInfo{}

	// Test SetSlot and GetSlot
	expectedSlot := math.Slot(12345)
	si.SetSlot(expectedSlot)
	require.Equal(t, expectedSlot, si.GetSlot())

	// Test SetIndex and GetIndex
	expectedIndex := math.U64(67890)
	si.SetIndex(expectedIndex)
	require.Equal(t, expectedIndex, si.GetIndex())
}
