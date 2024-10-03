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

package mock

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	ptypes "github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Compile time check to ensure BeaconState implements the methods
// required by the BeaconState for proofs.
var _ ptypes.BeaconState[
	*BeaconStateMarshallable,
	*types.ExecutionPayloadHeader,
	*types.Validator,
] = (*BeaconState)(nil)

// BeaconState is a mock implementation of the proof BeaconState interface
// using the default BeaconState type that is marshallable.
type (
	BeaconStateMarshallable = types.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.BeaconBlockHeader,
		types.Eth1Data,
		types.ExecutionPayloadHeader,
		types.Fork,
		types.Validator,
	]

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

	// Create an empty execution payload header with the given execution number
	// and fee recipient.
	execPayloadHeader := (&types.ExecutionPayloadHeader{}).Empty()
	execPayloadHeader.Number = executionNumber
	execPayloadHeader.FeeRecipient = executionFeeRecipient

	var (
		bsm = &BeaconStateMarshallable{}
		err error
	)
	bsm, err = bsm.New(
		0,
		common.Root{},
		slot,
		(&types.Fork{}).Empty(),
		(&types.BeaconBlockHeader{}).Empty(),
		[]common.Root{},
		[]common.Root{},
		(&types.Eth1Data{}).Empty(),
		0,
		execPayloadHeader,
		vals,
		[]uint64{},
		[]common.Bytes32{},
		0,
		0,
		[]math.Gwei{},
		0,
	)
	return &BeaconState{BeaconStateMarshallable: bsm}, err
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
