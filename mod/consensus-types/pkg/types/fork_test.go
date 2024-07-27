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

func TestFork_Serialization(t *testing.T) {
	original := (&types.Fork{}).New(
		common.Version{1, 2, 3, 4},
		common.Version{5, 6, 7, 8},
		math.Epoch(1000),
	)

	data, err := original.MarshalSSZ()
	require.NotNil(t, data)
	require.NoError(t, err)

	var unmarshalled types.Fork
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original, &unmarshalled)
}

func TestFork_SizeSSZ(t *testing.T) {
	fork := &types.Fork{
		PreviousVersion: common.Version{1, 2, 3, 4},
		CurrentVersion:  common.Version{5, 6, 7, 8},
		Epoch:           math.Epoch(1000),
	}

	size := fork.SizeSSZ()
	require.Equal(t, uint32(16), size)
}

func TestFork_HashTreeRoot(t *testing.T) {
	fork := &types.Fork{
		PreviousVersion: common.Version{1, 2, 3, 4},
		CurrentVersion:  common.Version{5, 6, 7, 8},
		Epoch:           math.Epoch(1000),
	}

	_, err := fork.HashTreeRoot()
	require.NoError(t, err)
}

func TestFork_GetTree(t *testing.T) {
	fork := &types.Fork{
		PreviousVersion: common.Version{1, 2, 3, 4},
		CurrentVersion:  common.Version{5, 6, 7, 8},
		Epoch:           math.Epoch(1000),
	}

	tree, err := fork.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestFork_UnmarshalSSZ_ErrSize(t *testing.T) {
	buf := make([]byte, 10) // size less than 16

	var unmarshalledFork types.Fork
	err := unmarshalledFork.UnmarshalSSZ(buf)

	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
