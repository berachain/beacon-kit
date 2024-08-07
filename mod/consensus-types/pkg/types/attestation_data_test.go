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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func generateAttestationData() *types.AttestationData {
	return &types.AttestationData{
		Slot:  12345,
		Index: 67890,
		BeaconBlockRoot: [32]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
		},
	}
}

func TestAttestationData_MarshalSSZ_UnmarshalSSZ(t *testing.T) {
	testCases := []struct {
		name     string
		data     *types.AttestationData
		expected *types.AttestationData
		err      error
	}{
		{
			name:     "Valid AttestationData",
			data:     generateAttestationData(),
			expected: generateAttestationData(),
			err:      nil,
		},
		{
			name: "Empty AttestationData",
			data: &types.AttestationData{
				Slot:            0,
				Index:           0,
				BeaconBlockRoot: [32]byte{},
			},
			expected: &types.AttestationData{
				Slot:            0,
				Index:           0,
				BeaconBlockRoot: [32]byte{},
			},
			err: nil,
		},
		{
			name:     "Invalid Buffer Size",
			data:     generateAttestationData(),
			expected: nil,
			err:      io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.data.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			var unmarshalled types.AttestationData
			if tc.name == "Invalid Buffer Size" {
				err = unmarshalled.UnmarshalSSZ(data[:32])
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				err = unmarshalled.UnmarshalSSZ(data)
				require.NoError(t, err)
				require.Equal(t, tc.expected, &unmarshalled)

				var buf []byte
				buf, err = tc.data.MarshalSSZTo(buf)
				require.NoError(t, err)

				// The two byte slices should be equal
				require.Equal(t, data, buf)
			}
		})
	}
}

func TestAttestationData_GetTree(t *testing.T) {
	data := generateAttestationData()

	tree, err := data.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)

	expectedRoot := data.HashTreeRoot()

	// Compare the tree root with the expected root
	actualRoot := tree.Hash()
	require.Equal(t, string(expectedRoot[:]), string(actualRoot))
}

func TestAttestationData_Getters(t *testing.T) {
	data := generateAttestationData()
	beaconBlockRoot := common.Root{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
	}

	require.NotNil(t, data)

	require.Equal(t, math.U64(12345), data.GetSlot())
	require.Equal(t, math.U64(67890), data.GetIndex())
	require.Equal(t, beaconBlockRoot, data.GetBeaconBlockRoot())
}
