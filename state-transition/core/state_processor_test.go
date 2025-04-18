//go:build test
// +build test

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

package core_test

import (
	"testing"
	"time"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

func TestStateProcessor_ProcessSlots(t *testing.T) {
	t.Parallel()
	cs := setupChain(t)
	//nolint:dogsled // used for testing
	sp, st, _, _, _, _ := statetransition.SetupTestState(t, cs)

	// Initialize state with genesis.
	genesisTime := time.Now().Truncate(time.Second)
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{}),
				Amount:      math.Gwei(cs.MaxEffectiveBalance()),
				Index:       0,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(cs.GenesisForkVersion()),
		}
	)
	genPayloadHeader.Timestamp = math.U64(genesisTime.Unix())
	_, err := sp.InitializeBeaconStateFromEth1(
		st, genDeposits, genPayloadHeader, cs.GenesisForkVersion(),
	)
	require.NoError(t, err)

	// Grab slot that will result in epoch processing in ProcessSlots.
	epochSlot := math.Slot(cs.SlotsPerEpoch()) - 1

	// Setup test cases
	tests := []struct {
		name      string
		stateSlot math.Slot
		slot      math.Slot
		wantErr   bool
	}{
		{"equal-slot", epochSlot, epochSlot, false},
		{"correct-slot", epochSlot, epochSlot + 1, false},
		{"too-large-slot", epochSlot, epochSlot + 2, true},
		{"too-small-slot", epochSlot, epochSlot - 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = st.SetSlot(tt.stateSlot)
			require.NoError(t, err)
			_, err = sp.ProcessSlots(st, tt.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
