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
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func generateSigningData() *types.SigningData {
	return &types.SigningData{
		ObjectRoot: bytes.B32{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		Domain: bytes.B32{
			33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
			51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
	}
}
func TestSigningData_MarshalSSZ_UnmarshalSSZ(t *testing.T) {
	testCases := []struct {
		name     string
		data     *types.SigningData
		expected *types.SigningData
		err      error
	}{
		{
			name:     "Valid SigningData",
			data:     generateSigningData(),
			expected: generateSigningData(),
			err:      nil,
		},
		{
			name: "Empty SigningData",
			data: &types.SigningData{
				ObjectRoot: bytes.B32{},
				Domain:     bytes.B32{},
			},
			expected: &types.SigningData{
				ObjectRoot: bytes.B32{},
				Domain:     bytes.B32{},
			},
			err: nil,
		},
		{
			name:     "Invalid Buffer Size",
			data:     generateSigningData(),
			expected: nil,
			err:      ssz.ErrSize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.data.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			var unmarshalled types.SigningData
			if tc.name == "Invalid Buffer Size" {
				err = unmarshalled.UnmarshalSSZ(data[:32])
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

func TestSigningData_GetTree(t *testing.T) {
	data := generateSigningData()

	tree, err := data.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)

	expectedRoot, err := data.HashTreeRoot()
	require.NoError(t, err)

	// Compare the tree root with the expected root
	actualRoot := tree.Hash()
	require.Equal(t, string(expectedRoot[:]), string(actualRoot))
}
