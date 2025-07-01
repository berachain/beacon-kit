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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func TestFork_Serialization(t *testing.T) {
	t.Parallel()
	fork := types.NewFork(
		common.Version{1, 2, 3, 4},
		common.Version{5, 6, 7, 8},
		math.Epoch(1000),
	)

	data, err := fork.MarshalSSZ()
	require.NotNil(t, data)
	require.NoError(t, err)

	var unmarshalled types.Fork
	err = ssz.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)
	require.Equal(t, fork, &unmarshalled)

	var buf []byte
	buf, err = fork.MarshalSSZTo(buf)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestFork_SizeSSZ(t *testing.T) {
	t.Parallel()
	fork := &types.Fork{
		PreviousVersion: common.Version{1, 2, 3, 4},
		CurrentVersion:  common.Version{5, 6, 7, 8},
		Epoch:           math.Epoch(1000),
	}

	size := karalabessz.Size(fork)
	require.Equal(t, uint32(16), size)
}

func TestFork_HashTreeRoot(t *testing.T) {
	t.Parallel()
	fork := &types.Fork{
		PreviousVersion: common.Version{1, 2, 3, 4},
		CurrentVersion:  common.Version{5, 6, 7, 8},
		Epoch:           math.Epoch(1000),
	}

	require.NotPanics(t, func() {
		_ = fork.HashTreeRoot()
	})
}

func TestFork_GetTree(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	buf := make([]byte, 10) // size less than 16

	unmarshalled := new(types.Fork)
	err := ssz.Unmarshal(buf, unmarshalled)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestFork_SerializationFormat(t *testing.T) {
	t.Parallel()

	fork := &types.Fork{
		PreviousVersion: common.Version{0x01, 0x02, 0x03, 0x04},
		CurrentVersion:  common.Version{0x05, 0x06, 0x07, 0x08},
		Epoch:           math.Epoch(0x1234567890ABCDEF),
	}

	data, err := fork.MarshalSSZ()
	require.NoError(t, err)

	// Document the expected SSZ format:
	// PreviousVersion: 4 bytes
	// CurrentVersion: 4 bytes
	// Epoch: 8 bytes (little-endian)
	expectedBytes := []byte{
		0x01, 0x02, 0x03, 0x04, // PreviousVersion
		0x05, 0x06, 0x07, 0x08, // CurrentVersion
		0xEF, 0xCD, 0xAB, 0x90, 0x78, 0x56, 0x34, 0x12, // Epoch (little-endian)
	}

	require.Equal(t, expectedBytes, data, "Serialization format should match SSZ specification")
}
