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
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/runtime/abci"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BuildABCIComponents returns the ABCI components for the beacon runtime.
func (r *BeaconKitRuntime) BuildABCIComponents(
	nextPrepare sdk.PrepareProposalHandler,
	nextProcess sdk.ProcessProposalHandler,
) (
	sdk.PrepareProposalHandler, sdk.ProcessProposalHandler,
	sdk.PreBlocker,
) {
	var (
		chainService   *blockchain.Service
		builderService *validator.Service
	)
	if err := r.services.FetchService(&chainService); err != nil {
		panic(err)
	}

	if err := r.services.FetchService(&builderService); err != nil {
		panic(err)
	}

	handler := abci.NewHandler(
		&r.cfg.ABCI,
		builderService,
		chainService,
		nextPrepare,
		nextProcess,
	)

	return handler.PrepareProposalHandler,
		handler.ProcessProposalHandler,
		handler.PreBlocker
}
