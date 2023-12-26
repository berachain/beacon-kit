// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package keeper

import (
	"context"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/cosmos/x/beacon/store"
	v1 "github.com/itsdevbear/bolaris/types/v1"
)

var LatestForkChoiceKey = []byte("latestForkChoice") //nolint:gochecknoglobals // fix later.

type (
	Keeper struct {
		executionClient engine.Caller
		storeKey        storetypes.StoreKey
	}
)

// NewKeeper creates new instances of the polaris Keeper.
func NewKeeper(
	executionClient engine.Caller,
	storeKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		executionClient: executionClient,
		storeKey:        storeKey,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	return sdk.UnwrapSDKContext(ctx).Logger()
}

// Setup initializes the polaris keeper.
func (k *Keeper) ForkChoiceStore(ctx context.Context) v1.ForkChoiceStore {
	return store.NewForkchoice(sdk.UnwrapSDKContext(ctx).KVStore(k.storeKey))
}
