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
	"errors"

	"github.com/berachain/beacon-kit/consensus-types/types"
	ptypes "github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/holiman/uint256"
)

// Compile time check to ensure BeaconState implements the methods
// required by the BeaconState for proofs.
var _ ptypes.BeaconState[*BeaconStateMarshallable] = (*BeaconState)(nil)

// BeaconState is a mock implementation of the proof BeaconState interface
// using the default BeaconState type that is marshallable.
type (
	BeaconStateMarshallable = types.BeaconState

	BeaconState struct {
		*BeaconStateMarshallable
	}
)

// NewBeaconState creates a new mock beacon state, with only the given slot,
// validators, execution number, and execution fee recipient.
func NewBeaconState(
	slot math.Slot,
	vals types.Validators,
	executionNumber math.U64,
	executionFeeRecipient common.ExecutionAddress,
) (*BeaconState, error) {
	// If no validators are provided, create an empty slice.
	if len(vals) == 0 {
		vals = make(types.Validators, 0)
	}

	// Create an empty execution payload header with the given execution number and fee recipient.
	execPayloadHeader := &types.ExecutionPayloadHeader{
		Number:        executionNumber,
		FeeRecipient:  executionFeeRecipient,
		BaseFeePerGas: &uint256.Int{},
	}

	bsm := &BeaconStateMarshallable{
		// TODO(pectra): Change this to an argument.
		Versionable:                  types.NewVersionable(version.Deneb()),
		Slot:                         slot,
		GenesisValidatorsRoot:        common.Root{},
		Fork:                         &types.Fork{},
		LatestBlockHeader:            &types.BeaconBlockHeader{},
		BlockRoots:                   []common.Root{},
		StateRoots:                   []common.Root{},
		LatestExecutionPayloadHeader: execPayloadHeader,
		Eth1Data:                     &types.Eth1Data{},
		Eth1DepositIndex:             0,
		Validators:                   vals,
		Balances:                     []uint64{},
		RandaoMixes:                  []common.Bytes32{},
		NextWithdrawalIndex:          0,
		NextWithdrawalValidatorIndex: 0,
		Slashings:                    []math.Gwei{},
		TotalSlashing:                0,
	}

	return &BeaconState{BeaconStateMarshallable: bsm}, nil
}

// GetLatestExecutionPayloadHeader implements proof BeaconState.
func (m *BeaconState) GetLatestExecutionPayloadHeader() (
	*types.ExecutionPayloadHeader, error,
) {
	return m.BeaconStateMarshallable.LatestExecutionPayloadHeader, nil
}

// GetMarshallable implements proof BeaconState.
func (m *BeaconState) GetMarshallable() (
	*BeaconStateMarshallable, error,
) {
	return m.BeaconStateMarshallable, nil
}

// ValidatorByIndex implements proof BeaconState.
func (m *BeaconState) ValidatorByIndex(
	index math.ValidatorIndex,
) (*types.Validator, error) {
	vals := m.BeaconStateMarshallable.Validators
	if index >= math.ValidatorIndex(len(vals)) {
		return nil, errors.New("validator index out of range")
	}

	return vals[index], nil
}
