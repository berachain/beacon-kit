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

package runtime

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	builder "github.com/itsdevbear/bolaris/beacon/builder"
	"github.com/itsdevbear/bolaris/beacon/sync"
	"github.com/itsdevbear/bolaris/runtime/abci/listener"
	"github.com/itsdevbear/bolaris/runtime/abci/preblock"
	"github.com/itsdevbear/bolaris/runtime/abci/proposal"
)

func (r *BeaconKitRuntime) BuildABCIComponents(
	nextPrepare sdk.PrepareProposalHandler,
	nextProcess sdk.ProcessProposalHandler,
	nextPreblocker sdk.PreBlocker,
	logger log.Logger,
) (
	sdk.PrepareProposalHandler, sdk.ProcessProposalHandler,
	sdk.PreBlocker, storetypes.StreamingManager,
) {
	var (
		chainService   *blockchain.Service
		builderService *builder.Service
		syncService    *sync.Service
	)
	if err := r.services.FetchService(&chainService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&syncService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&builderService); err != nil {
		panic(err)
	}

	proposalHandler := proposal.NewHandler(
		&r.cfg.ABCI,
		builderService,
		syncService,
		chainService,
		nextPrepare,
		nextProcess,
	)

	preBlocker := preblock.NewBeaconPreBlockHandler(
		&r.cfg.ABCI, r.logger, chainService, syncService, nextPreblocker,
	).PreBlocker()

	return proposalHandler.PrepareProposalHandler,
		proposalHandler.ProcessProposalHandler,
		preBlocker,
		storetypes.StreamingManager{
			ABCIListeners: []storetypes.ABCIListener{
				listener.NewBeaconListener(
					logger.With("module", "beacon-listener"),
					chainService,
				),
			},
		}
}
