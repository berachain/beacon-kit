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

	cmdtypes "github.com/berachain/beacon-kit/beacond/cmd/types"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	ptypes "github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// mockBeaconState is a mock implementation of the BeaconState interface.
type mockBeaconState struct {
	vals types.Validators
	bsm  *cmdtypes.BeaconStateMarshallable
}

// NewBeaconState creates a new mock beacon state, with only the given slot and
// validators set.
func NewBeaconState(
	slot math.Slot,
	vals types.Validators,
) (
	ptypes.BeaconState[
		*cmdtypes.BeaconStateMarshallable,
		*types.ExecutionPayloadHeader,
		*types.Validator,
	],
	error,
) {
	var (
		bsm = &cmdtypes.BeaconStateMarshallable{}
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
		(&types.ExecutionPayloadHeader{}).Empty(),
		vals,
		[]uint64{},
		[]common.Bytes32{},
		0,
		0,
		[]uint64{},
		0,
	)
	if err != nil {
		return nil, err
	}

	return &mockBeaconState{
		vals: vals,
		bsm:  bsm,
	}, nil
}

func (*mockBeaconState) GetLatestExecutionPayloadHeader() (
	*types.ExecutionPayloadHeader, error,
) {
	return &types.ExecutionPayloadHeader{}, nil
}

func (m *mockBeaconState) GetMarshallable() (
	*cmdtypes.BeaconStateMarshallable, error,
) {
	return m.bsm, nil
}

func (m *mockBeaconState) ValidatorByIndex(
	index math.ValidatorIndex,
) (*types.Validator, error) {
	if index >= math.ValidatorIndex(len(m.vals)) {
		return nil, errors.New("validator index out of range")
	}

	return m.vals[index], nil
}

func (m *mockBeaconState) HashTreeRoot() common.Root {
	bsm, err := m.GetMarshallable()
	if err != nil {
		return common.Root{}
	}

	return bsm.HashTreeRoot()
}
