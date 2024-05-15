// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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
	"cosmossdk.io/core/appmodule"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
)

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
//
// TODO: The keeper will eventually be dissolved.
type Keeper struct {
	storage.Backend[state.BeaconState]
}

// TODO: move this.
func DenebPayloadFactory() engineprimitives.ExecutionPayloadHeader {
	return &engineprimitives.ExecutionPayloadHeaderDeneb{}
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	fdb *filedb.DB,
	env appmodule.Environment,
	cs primitives.ChainSpec,
	ddb *deposit.KVStore,
) *Keeper {
	return &Keeper{
		Backend: *storage.NewBackend[state.BeaconState](
			cs,
			dastore.New[types.BeaconBlockBody](
				cs, filedb.NewRangeDB(fdb),
			),
			beacondb.New[
				*types.Fork,
				*types.BeaconBlockHeader,
				engineprimitives.ExecutionPayloadHeader,
				*types.Eth1Data,
				*types.Validator,
			](env.KVStoreService, DenebPayloadFactory), ddb),
	}
}
