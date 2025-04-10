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

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/errors"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/beacondb"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var errTestMemberNotImplemented = errors.New("not implemented")

// testConsensusService stubs consensus service
type testConsensusService struct {
	cms     storetypes.CommitMultiStore
	kvStore *beacondb.KVStore
	cs      chain.Spec
}

func (t *testConsensusService) CreateQueryContext(height int64, _ bool) (sdk.Context, error) {
	sdkCtx := sdk.NewContext(t.cms.CacheMultiStore(), false, log.NewNopLogger())

	// there validations mimics consensus service, not sure if they are necessary
	tmpState := statedb.NewBeaconStateFromDB(t.kvStore.WithContext(sdkCtx), t.cs)
	slot, err := tmpState.GetSlot()
	if err != nil {
		return sdk.Context{}, sdkerrors.ErrInvalidHeight
	}
	if height > int64(slot.Unwrap()) {
		return sdk.Context{}, sdkerrors.ErrInvalidHeight
	}
	// end of possibly unnecessary validations

	return sdkCtx, nil
}

func (t *testConsensusService) Start(_ context.Context) error {
	return errTestMemberNotImplemented
}

func (t *testConsensusService) Stop() error {
	return errTestMemberNotImplemented
}

func (t *testConsensusService) Name() string {
	panic(errTestMemberNotImplemented)
}

func (t *testConsensusService) LastBlockHeight() int64 {
	panic(errTestMemberNotImplemented)
}
