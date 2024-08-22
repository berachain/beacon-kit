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

package block_test

import (
	"testing"

	storev2 "cosmossdk.io/store/v2/db"
	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
	"github.com/stretchr/testify/require"
)

type MockBeaconBlock struct {
	slot math.Slot
}

func (m MockBeaconBlock) GetSlot() math.Slot {
	return m.slot
}

func (m MockBeaconBlock) HashTreeRoot() common.Root {
	return [32]byte{byte(m.slot)}
}

func (m MockBeaconBlock) GetExecutionNumber() math.U64 {
	return m.slot
}

func (m MockBeaconBlock) GetStateRoot() common.Root {
	return [32]byte{byte(m.slot)}
}

func TestBlockStore(t *testing.T) {
	blockStore := block.NewStore[*MockBeaconBlock](
		storage.NewKVStoreProvider(storev2.NewMemDB()),
		noop.NewLogger[any](),
		uint64(5),
	)

	// Set 7 blocks.
	// The latest block is 7 and should hold the last 5 blocks in the window.
	for i := 1; i <= 7; i++ {
		blockStore.Set(&MockBeaconBlock{slot: math.Slot(i)})
	}

	// The window of slots in the store should look like this:
	// idx  [0, 1, 2, 3, 4] (actual values in maps)
	// slot [5, 6, 7, 3, 4] (corresponding true slots)
	slot, err := blockStore.GetSlotFromWindow(0)
	require.NoError(t, err)
	require.Equal(t, math.Slot(5), slot)
	slot, err = blockStore.GetSlotFromWindow(1)
	require.NoError(t, err)
	require.Equal(t, math.Slot(6), slot)
	slot, err = blockStore.GetSlotFromWindow(2)
	require.NoError(t, err)
	require.Equal(t, math.Slot(7), slot)
	slot, err = blockStore.GetSlotFromWindow(3)
	require.NoError(t, err)
	require.Equal(t, math.Slot(3), slot)
	slot, err = blockStore.GetSlotFromWindow(4)
	require.NoError(t, err)
	require.Equal(t, math.Slot(4), slot)

	// Get the slots by roots & execution numbers.
	for i := math.Slot(3); i <= 7; i++ {
		slot, err = blockStore.GetSlotByBlockRoot([32]byte{byte(i)})
		require.NoError(t, err)
		require.Equal(t, i, slot)

		slot, err = blockStore.GetSlotByExecutionNumber(i)
		require.NoError(t, err)
		require.Equal(t, i, slot)

		slot, err = blockStore.GetSlotByStateRoot([32]byte{byte(i)})
		require.NoError(t, err)
		require.Equal(t, i, slot)
	}

	// Try getting a slot that doesn't exist.
	_, err = blockStore.GetSlotByBlockRoot([32]byte{byte(8)})
	require.ErrorContains(t, err, "not found")
	_, err = blockStore.GetSlotByExecutionNumber(2)
	require.ErrorContains(t, err, "not found")

	// Try calling GetSlotFromWindow with an out of bounds index.
	_, err = blockStore.GetSlotFromWindow(5)
	require.ErrorContains(t, err, "index out of bounds")
}
