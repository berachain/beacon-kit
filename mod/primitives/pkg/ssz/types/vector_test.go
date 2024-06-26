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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
	"github.com/stretchr/testify/require"
)

func TestSSZVectorBasic(t *testing.T) {
	t.Run("SizeSSZ for uint8 vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZByte]{1, 2, 3, 4, 5}
		require.Equal(t, 5, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for byte slice vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZUInt8]{1, 2, 3, 4, 5, 6, 7, 8}
		require.Equal(t, 8, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for uint64 vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZUInt64]{1, 2, 3, 4, 5}
		require.Equal(t, 40, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for bool vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZBool]{true, false, true}
		require.Equal(t, 3, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for empty vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZUInt64]{}
		require.Equal(t, 0, vector.SizeSSZ())
	})
}
