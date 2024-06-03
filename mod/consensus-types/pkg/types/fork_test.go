// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func TestFork_Serialization(t *testing.T) {
	original := &types.Fork{
		PreviousVersion: common.Version{1, 2, 3, 4},
		CurrentVersion:  common.Version{5, 6, 7, 8},
		Epoch:           math.Epoch(1000),
	}

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
	require.Equal(t, 16, size)
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

	require.ErrorIs(t, err, ssz.ErrSize)
}
