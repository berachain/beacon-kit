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

package ssz_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestVectorSizeSSZ(t *testing.T) {
	t.Run("uint8 vector", func(t *testing.T) {
		vector := ssz.Vector[ssz.Byte]{1, 2, 3, 4, 5}
		require.Len(t, vector, 5)
		require.Equal(t, 5, vector.SizeSSZ())
	})

	t.Run("byte slice vector", func(t *testing.T) {
		vector := ssz.Vector[math.U8]{1, 2, 3, 4, 5, 6, 7, 8}
		require.Len(t, vector, 8)
		require.Equal(t, 8, vector.SizeSSZ())
	})

	t.Run("uint64 vector", func(t *testing.T) {
		vector := ssz.Vector[math.U64]{1, 2, 3, 4, 5}
		require.Len(t, vector, 5)
		require.Equal(t, 40, vector.SizeSSZ())
	})

	t.Run("bool vector", func(t *testing.T) {
		vector := ssz.Vector[math.Bool]{true, false, true}
		require.Len(t, vector, 3)
		require.Equal(t, 3, vector.SizeSSZ())
	})

	t.Run("empty vector", func(t *testing.T) {
		vector := ssz.Vector[math.U64]{}
		require.Empty(t, vector)
		require.Equal(t, 0, vector.SizeSSZ())
	})
}

func TestVectorHashTreeRoot(t *testing.T) {
	t.Run("uint8 vector", func(t *testing.T) {
		vector := ssz.Vector[math.U8]{1, 2, 3, 4, 5}
		root, err := vector.HashTreeRoot()
		require.NoError(t, err)
		require.NotEqual(t, [32]byte{}, root)
	})

	t.Run("bool vector", func(t *testing.T) {
		vector := ssz.Vector[math.Bool]{true, false, true, false}
		root, err := vector.HashTreeRoot()
		require.NoError(t, err)
		require.NotEqual(t, [32]byte{}, root)
	})

	t.Run("uint64 vector", func(t *testing.T) {
		vector := ssz.Vector[math.U64]{1, 2, 3, 4, 5}
		root, err := vector.HashTreeRoot()
		require.NoError(t, err)
		require.NotEqual(t, [32]byte{}, root)
	})

	t.Run("consistency", func(t *testing.T) {
		vector1 := ssz.Vector[math.U8]{1, 2, 3, 4, 5}
		vector2 := ssz.Vector[math.U8]{1, 2, 3, 4, 5}
		root1, err1 := vector1.HashTreeRoot()
		root2, err2 := vector2.HashTreeRoot()
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Equal(t, root1, root2)
	})
}
