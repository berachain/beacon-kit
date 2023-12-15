// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package keeper

import (
	"context"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/beacon/prysm"
	"github.com/itsdevbear/bolaris/cosmos/x/evm/types"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

var LatestForkChoiceKey = []byte("latestForkChoice")

type (
	Keeper struct {
		// consensusAPI is the consensus API
		executionClient prysm.EngineCaller
		storeKey        storetypes.StoreKey
	}
)

// NewKeeper creates new instances of the polaris Keeper.
func NewKeeper(
	executionClient prysm.EngineCaller,
	storeKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		executionClient: executionClient,
		storeKey:        storeKey,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	return sdk.UnwrapSDKContext(ctx).Logger().With(types.ModuleName)
}

// Forkchoice Updates

func (k *Keeper) LatestForkChoice(ctx context.Context) (*enginev1.ForkchoiceState, error) {
	bz := sdk.UnwrapSDKContext(ctx).KVStore(k.storeKey).Get(LatestForkChoiceKey)
	if bz == nil {
		return nil, nil
	}

	var state enginev1.ForkchoiceState
	if err := state.UnmarshalJSON(bz); err != nil {
		return nil, err
	}
	return &state, nil
}

func (k *Keeper) SetLatestForkChoice(ctx context.Context, state *enginev1.ForkchoiceState) error {
	bz, err := state.MarshalJSON()
	if err != nil {
		return err
	}

	sdk.UnwrapSDKContext(ctx).KVStore(k.storeKey).Set(LatestForkChoiceKey, bz)
	return nil
}
