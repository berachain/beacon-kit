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
//go:build test
// +build test

package backend_test

import (
	"context"
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend"
	types "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestFilteredValidators(t *testing.T) {
	t.Parallel()

	// Build backend to test
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)
	sp, st, kvStore, _, _ := statetransition.SetupTestState(t, cs)
	sb := storage.NewBackend(cs, nil, kvStore, nil, nil)

	b := backend.New(sb, cs, sp)
	tcs := &testConsensusService{state: st}
	b.AttachQueryBackend(tcs)

	// Setup context
	stateSlot := math.Slot(10)
	require.NoError(t, st.SetSlot(stateSlot))

	tests := []struct {
		name        string
		inputsF     func() (math.Slot, []string /*ids*/, []string /*statuses*/)
		expectedErr error
		checkF      func(res []*types.ValidatorData) error
	}{
		{
			name: "height too high",
			inputsF: func() (math.Slot, []string, []string) {
				return stateSlot + 1, nil, nil
			},

			// this error really comes from testConsensusService.CreateQueryContext
			// so not really testing the implementation,but I gotta start somewhere
			expectedErr: sdkerrors.ErrInvalidHeight,
			checkF:      func(res []*types.ValidatorData) error { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot, ids, statuses := tt.inputsF()

			res, err := b.FilteredValidators(slot, ids, statuses)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.expectedErr)
			}
			require.NoError(t, tt.checkF(res))
		})
	}
}

var errTestMemberNotImplemented = errors.New("not implemented")

type testConsensusService struct {
	state *statedb.StateDB
}

func (t *testConsensusService) Start(ctx context.Context) error {
	return errTestMemberNotImplemented
}

func (t *testConsensusService) Stop() error {
	return errTestMemberNotImplemented
}

func (t *testConsensusService) Name() string {
	panic(errTestMemberNotImplemented)
}

func (t *testConsensusService) CreateQueryContext(height int64, prove bool) (sdk.Context, error) {
	slot, err := t.state.GetSlot()
	if err != nil {
		return sdk.Context{}, sdkerrors.ErrInvalidHeight
	}
	if height > int64(slot.Unwrap()) {
		return sdk.Context{}, sdkerrors.ErrInvalidHeight
	}
	return sdk.Context{}, nil
}

func (t *testConsensusService) LastBlockHeight() int64 {
	panic(errTestMemberNotImplemented)
}
