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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	karalabessz "github.com/karalabe/ssz"
	zrntcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	zspec "github.com/protolambda/zrnt/eth2/configs"
	ztree "github.com/protolambda/ztyp/tree"
	"github.com/stretchr/testify/require"
)

func TestWithdrawals(t *testing.T) {
	t.Parallel()
	t.Run("SizeSSZ", func(t *testing.T) {
		t.Parallel()
		withdrawals := engineprimitives.Withdrawals{
			{Index: 1, Validator: 2, Address: [20]byte{1, 2, 3}, Amount: 100},
			{Index: 3, Validator: 4, Address: [20]byte{4, 5, 6}, Amount: 200},
		}
		expectedSize := uint32(
			len(withdrawals),
		) * engineprimitives.WithdrawalSize
		require.Equal(t, expectedSize, karalabessz.Size(withdrawals))
	})

	t.Run("HashTreeRoot", func(t *testing.T) {
		t.Parallel()
		withdrawals := engineprimitives.Withdrawals{
			{Index: 1, Validator: 2, Address: [20]byte{1, 2, 3}, Amount: 100},
			{Index: 3, Validator: 4, Address: [20]byte{4, 5, 6}, Amount: 200},
		}
		root := withdrawals.HashTreeRoot()
		require.NotEmpty(t, root)
	})

	t.Run("HashTreeRoot Comparison", func(t *testing.T) {
		t.Parallel()
		withdrawals := engineprimitives.Withdrawals{
			{Index: 1, Validator: 2, Address: [20]byte{1, 2, 3}, Amount: 100},
		}
		zwithdrawals := zrntcommon.Withdrawals{
			{
				Index:          zrntcommon.WithdrawalIndex(withdrawals[0].Index),
				ValidatorIndex: zrntcommon.ValidatorIndex(withdrawals[0].Validator),
				Address:        zrntcommon.Eth1Address(withdrawals[0].Address),
				Amount:         zrntcommon.Gwei(withdrawals[0].Amount),
			},
		}
		root := withdrawals.HashTreeRoot()
		hFn := ztree.GetHashFn()
		spec := zspec.Mainnet
		zroot := zwithdrawals.HashTreeRoot(spec, hFn)
		require.NotEmpty(t, root)
		require.Equal(t, root[:], zroot[:])
	})

	t.Run("HashTreeRoot", func(t *testing.T) {
		t.Parallel()
		withdrawals := engineprimitives.Withdrawals{
			{
				Index:     math.U64(1),
				Validator: math.ValidatorIndex(2),
				Address:   common.ExecutionAddress{1, 2, 3},
				Amount:    math.Gwei(100),
			},
			{
				Index:     math.U64(3),
				Validator: math.ValidatorIndex(4),
				Address:   common.ExecutionAddress{4, 5, 6},
				Amount:    math.Gwei(200),
			},
		}

		root := withdrawals.HashTreeRoot()
		require.NotEmpty(t, root)

		// Verify that the root changes when the withdrawals change
		withdrawals[0].Amount = math.Gwei(150)
		newRoot := withdrawals.HashTreeRoot()
		require.NotEqual(t, root, newRoot)

		// Verify that the order of withdrawals matters
		reversedWithdrawals := engineprimitives.Withdrawals{
			withdrawals[1],
			withdrawals[0],
		}
		reversedRoot := reversedWithdrawals.HashTreeRoot()
		require.NotEqual(t, newRoot, reversedRoot)
	})

	t.Run("HashTreeRoot of Empty List", func(t *testing.T) {
		t.Parallel()
		emptyWithdrawals := engineprimitives.Withdrawals{}
		emptyRoot := emptyWithdrawals.HashTreeRoot()
		require.NotEmpty(t, emptyRoot)

		// Verify that the root of an empty list is different from a non-empty
		// list
		nonEmptyWithdrawals := engineprimitives.Withdrawals{
			{
				Index:     math.U64(1),
				Validator: math.ValidatorIndex(2),
				Address:   common.ExecutionAddress{1, 2, 3},
				Amount:    math.Gwei(100),
			},
		}
		nonEmptyRoot := nonEmptyWithdrawals.HashTreeRoot()
		require.NotEqual(t, emptyRoot, nonEmptyRoot)
	})
}
