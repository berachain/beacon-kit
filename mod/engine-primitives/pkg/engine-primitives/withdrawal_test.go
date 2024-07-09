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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestWithdrawal(t *testing.T) {
	withdrawal := (&engineprimitives.Withdrawal{}).New(
		math.U64(1),
		math.ValidatorIndex(1),
		gethprimitives.ExecutionAddress{1, 2, 3, 4, 5},
		math.Gwei(1000),
	)

	require.Equal(t, math.U64(1), withdrawal.GetIndex())
	require.Equal(t, math.ValidatorIndex(1), withdrawal.GetValidatorIndex())
	require.Equal(t,
		gethprimitives.ExecutionAddress{1, 2, 3, 4, 5},
		withdrawal.GetAddress(),
	)
	require.Equal(t, math.Gwei(1000), withdrawal.GetAmount())
}

func TestWithdrawal_Equals(t *testing.T) {
	withdrawal1 := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(1),
		Address:   gethprimitives.ExecutionAddress{1, 2, 3, 4, 5},
		Amount:    math.Gwei(1000),
	}

	withdrawal2 := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(1),
		Address:   gethprimitives.ExecutionAddress{1, 2, 3, 4, 5},
		Amount:    math.Gwei(1000),
	}

	withdrawal3 := &engineprimitives.Withdrawal{
		Index:     math.U64(2),
		Validator: math.ValidatorIndex(2),
		Address:   gethprimitives.ExecutionAddress{2, 3, 4, 5, 6},
		Amount:    math.Gwei(2000),
	}

	// Test that Equals returns true for two identical withdrawals
	require.True(t, withdrawal1.Equals(withdrawal2))

	// Test that Equals returns false for two different withdrawals
	require.False(t, withdrawal1.Equals(withdrawal3))
}

func TestWithdrawal_HashTreeRoot(t *testing.T) {
	withdrawal := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(2),
		Address: gethprimitives.ExecutionAddress{
			1,
			2,
			3,
			4,
			5,
			6,
			7,
			8,
			9,
			10,
			11,
			12,
			13,
			14,
			15,
			16,
			17,
			18,
			19,
			20,
		},
		Amount: math.Gwei(1000),
	}

	// Get the hash tree root using the built-in method
	builtInRoot, err := withdrawal.HashTreeRoot()
	require.NoError(t, err)

	// Create a Container with the same elements
	container := ssz.ContainerFromElements(
		withdrawal.Index,
		withdrawal.Validator,
		ssz.ByteVectorFromBytes(withdrawal.Address.Bytes()),
		withdrawal.Amount,
	)

	// Get the hash tree root using the Container
	containerRoot, err := container.HashTreeRoot()
	require.NoError(t, err)

	// Compare the results
	require.Equal(
		t,
		builtInRoot,
		containerRoot,
		"Hash tree roots should be equal",
	)
}

func TestWithdrawalMethods(t *testing.T) {
	withdrawal := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(2),
		Address:   [20]byte{1, 2, 3},
		Amount:    math.Gwei(100),
	}

	t.Run("IsFixed", func(t *testing.T) {
		require.True(t, withdrawal.IsFixed())
	})

	t.Run("Type", func(t *testing.T) {
		require.True(t, withdrawal.Type().ID().IsContainer())
	})

	t.Run("ItemLength", func(t *testing.T) {
		require.Equal(t, uint64(constants.RootLength), withdrawal.ItemLength())
	})

	t.Run("Getters", func(t *testing.T) {
		require.Equal(t, math.U64(1), withdrawal.GetIndex())
		require.Equal(t, math.ValidatorIndex(2), withdrawal.GetValidatorIndex())
		require.Equal(t, common.Address([20]byte{0x01, 0x02, 0x03}),
			withdrawal.GetAddress())
		require.Equal(t, math.U64(100), withdrawal.GetAmount())
	})
}
