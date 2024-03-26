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
	"github.com/berachain/beacon-kit/beacon/blockchain"
	builder "github.com/berachain/beacon-kit/beacon/builder"
	"github.com/berachain/beacon-kit/beacon/sync"
	"github.com/berachain/beacon-kit/health"
	"github.com/berachain/beacon-kit/runtime/abci/preblock"
	"github.com/berachain/beacon-kit/runtime/abci/proposal"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BuildABCIComponents returns the ABCI components for the beacon runtime.
func (r *BeaconKitRuntime) BuildABCIComponents(
	nextPrepare sdk.PrepareProposalHandler,
	nextProcess sdk.ProcessProposalHandler,
	nextPreblocker sdk.PreBlocker,
) (
	sdk.PrepareProposalHandler, sdk.ProcessProposalHandler,
	sdk.PreBlocker,
) {
	var (
		chainService   *blockchain.Service
		builderService *builder.Service
		healthService  *health.Service
		syncService    *sync.Service
	)
	if err := r.services.FetchService(&chainService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&healthService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&builderService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&syncService); err != nil {
		panic(err)
	}

	proposalHandler := proposal.NewHandler(
		&r.cfg.ABCI,
		builderService,
		healthService,
		chainService,
		nextPrepare,
		nextProcess,
	)

	preBlocker := preblock.NewBeaconPreBlockHandler(
		&r.cfg.ABCI,
		r.logger,
		chainService,
		syncService,
		nextPreblocker,
	).PreBlocker()

	return proposalHandler.PrepareProposalHandler,
		proposalHandler.ProcessProposalHandler,
		preBlocker
}
