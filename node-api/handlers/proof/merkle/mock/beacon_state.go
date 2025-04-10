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

package mock

import (
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

// BeaconState is a mock implementation of the StateDB BeaconState
type BeaconState struct {
	internal types.BeaconState
}

// NewBeaconState creates a new mock beacon state, with only the given slot,
// validators, execution number, and execution fee recipient.
func NewBeaconState(
	slot math.Slot,
	vals types.Validators,
	executionNumber math.U64,
	executionFeeRecipient common.ExecutionAddress,
	forkVersion common.Version,
) (BeaconState, error) {
	// If no validators are provided, create an empty slice.
	if len(vals) == 0 {
		vals = make(types.Validators, 0)
	}

	// Create an empty execution payload header with the given execution number and fee recipient.
	execPayloadHeader := types.NewEmptyExecutionPayloadHeaderWithVersion(forkVersion)
	execPayloadHeader.Number = executionNumber
	execPayloadHeader.FeeRecipient = executionFeeRecipient

	bsm := types.NewEmptyBeaconStateWithVersion(forkVersion)
	bsm.Slot = slot
	bsm.GenesisValidatorsRoot = common.Root{}
	bsm.Fork = &types.Fork{}
	bsm.LatestBlockHeader = types.NewEmptyBeaconBlockHeader()
	bsm.BlockRoots = []common.Root{}
	bsm.StateRoots = []common.Root{}
	bsm.LatestExecutionPayloadHeader = execPayloadHeader
	bsm.Eth1Data = &types.Eth1Data{}
	bsm.Eth1DepositIndex = 0
	bsm.Validators = vals
	bsm.Balances = []uint64{}
	bsm.RandaoMixes = []common.Bytes32{}
	bsm.NextWithdrawalIndex = 0
	bsm.NextWithdrawalValidatorIndex = 0
	bsm.Slashings = []math.Gwei{}
	bsm.TotalSlashing = 0

	return BeaconState{*bsm}, nil
}

// GetMarshallable implements proof BeaconState.
func (m *BeaconState) GetMarshallable() (
	*types.BeaconState, error,
) {
	beaconState := types.NewEmptyBeaconStateWithVersion(m.internal.GetForkVersion())
	beaconState.Slot = m.internal.Slot
	beaconState.GenesisValidatorsRoot = m.internal.GenesisValidatorsRoot
	beaconState.Fork = m.internal.Fork
	beaconState.LatestBlockHeader = m.internal.LatestBlockHeader
	beaconState.BlockRoots = m.internal.BlockRoots
	beaconState.StateRoots = m.internal.StateRoots
	beaconState.LatestExecutionPayloadHeader = m.internal.LatestExecutionPayloadHeader
	beaconState.Eth1Data = m.internal.Eth1Data
	beaconState.Eth1DepositIndex = m.internal.Eth1DepositIndex
	beaconState.Validators = m.internal.Validators
	beaconState.Balances = m.internal.Balances
	beaconState.RandaoMixes = m.internal.RandaoMixes
	beaconState.NextWithdrawalIndex = m.internal.NextWithdrawalIndex
	beaconState.NextWithdrawalValidatorIndex = m.internal.NextWithdrawalValidatorIndex
	beaconState.Slashings = m.internal.Slashings
	beaconState.TotalSlashing = m.internal.TotalSlashing

	if version.EqualsOrIsAfter(beaconState.GetForkVersion(), version.Electra()) {
		beaconState.PendingPartialWithdrawals = m.internal.PendingPartialWithdrawals
	}
	return beaconState, nil
}

// HashTreeRoot is the interface for the beacon store.
func (m *BeaconState) HashTreeRoot() common.Root {
	st, err := m.GetMarshallable()
	if err != nil {
		panic(err)
	}
	return st.HashTreeRoot()
}
