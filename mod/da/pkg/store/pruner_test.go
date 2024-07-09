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

package store_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/chain-spec/pkg/chain"
	"github.com/berachain/beacon-kit/mod/da/pkg/store"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

// Mock implementations for BeaconBlock.
type MockBeaconBlock struct {
	slot math.U64
}

// Mock implementations GetSlot() for BeaconBlock.
func (b MockBeaconBlock) GetSlot() math.U64 {
	return b.slot
}

// Mock implementations for BlockEvent.
type MockBlockEvent struct {
	data MockBeaconBlock
}

// Mock implementations Data() for BlockEvent.
func (e MockBlockEvent) Data() MockBeaconBlock {
	return e.data
}

// TestBuildPruneRangeFn tests the BuildPruneRangeFn function.
//

func TestBuildPruneRangeFn(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		slotsPerEpoch uint64
		minEpochs     uint64
		eventSlot     math.U64
		expectedStart uint64
		expectedEnd   uint64
	}{
		{
			name:          "Slot greater than window",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(200),
			expectedStart: 0,
			expectedEnd:   40,
		},
		{
			name:          "Slot less than window",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(100),
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Exact boundary case",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(160),
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Greater than boundary just 1 uint case",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(161),
			expectedStart: 0,
			expectedEnd:   1,
		},
		{
			name:          "Zero slot case",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(0),
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "SlotsPerEpoch as one",
			slotsPerEpoch: 1,
			minEpochs:     5,
			eventSlot:     math.U64(50),
			expectedStart: 0,
			expectedEnd:   45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := chain.NewChainSpec(
				chain.SpecData[
					bytes.B4, math.U64, gethprimitives.ExecutionAddress, math.U64, any,
				]{
					SlotsPerEpoch:                    tt.slotsPerEpoch,
					MinEpochsForBlobsSidecarsRequest: tt.minEpochs,
				},
			)
			pruneFn := store.BuildPruneRangeFn[MockBeaconBlock, MockBlockEvent](
				cs,
			)
			event := MockBlockEvent{
				data: MockBeaconBlock{
					slot: tt.eventSlot,
				},
			}
			start, end := pruneFn(event)
			require.Equal(
				t,
				tt.expectedStart,
				start,
				"Test case : %s (expectedStart)",
				tt.name,
			)
			require.Equal(
				t,
				tt.expectedEnd,
				end,
				"Test case : %s (expectedEnd)",
				tt.name,
			)
		})
	}
}
