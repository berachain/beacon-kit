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

package runtime

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	initialsync "github.com/itsdevbear/bolaris/beacon/initial-sync"
	"github.com/itsdevbear/bolaris/runtime/abci/preblock"
	"github.com/itsdevbear/bolaris/runtime/abci/proposal"
)

// CosmosApp is an interface that defines the methods needed for the Cosmos setup.
type CosmosApp interface {
	baseapp.ProposalTxVerifier
	SetPrepareProposal(sdk.PrepareProposalHandler)
	SetProcessProposal(sdk.ProcessProposalHandler)
	SetVerifyVoteExtensionHandler(sdk.VerifyVoteExtensionHandler)
	PreBlocker() sdk.PreBlocker
	SetPreBlocker(sdk.PreBlocker)
	Mempool() mempool.Mempool
}

func (r *BeaconKitRuntime) RegisterApp(app CosmosApp) error {
	var (
		chainService *blockchain.Service
		syncService  *initialsync.Service
	)
	if err := r.services.FetchService(&chainService); err != nil {
		return err
	}

	if err := r.services.FetchService(&syncService); err != nil {
		panic(err)
	}

	// Build and Register Prepare and Process Proposal Handlers.
	defaultProposalHandler := baseapp.NewDefaultProposalHandler(app.Mempool(), app)
	proposalHandler := proposal.NewHandler(
		chainService,
		defaultProposalHandler.PrepareProposalHandler(),
		defaultProposalHandler.ProcessProposalHandler(),
	)
	app.SetPrepareProposal(proposalHandler.PrepareProposalHandler)
	app.SetProcessProposal(proposalHandler.ProcessProposalHandler)

	// Build and Register Preblock Handler.
	app.SetPreBlocker(
		preblock.NewBeaconPreBlockHandler(r.logger, r.fscp, syncService, nil).PreBlocker(),
	)
	return nil
}
