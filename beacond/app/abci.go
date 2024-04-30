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

package app

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
)

// PrepareProposal is called by the consensus engine to prepare a proposal for
// the next block.
func (app BeaconApp) PrepareProposal(
	req *abci.PrepareProposalRequest,
) (*abci.PrepareProposalResponse, error) {
	start := time.Now()
	defer func() {
		app.Logger().
			Info("prepareProposal executed",
				"duration", time.Since(start).String())
	}()
	return app.BaseApp.PrepareProposal(req)
}

// ProcessProposal is called by the consensus engine when a new proposal block
// is received.
func (app BeaconApp) ProcessProposal(
	req *abci.ProcessProposalRequest,
) (*abci.ProcessProposalResponse, error) {
	start := time.Now()
	defer func() {
		app.Logger().
			Info("processProposal executed",
				"duration", time.Since(start).String())
	}()
	return app.BaseApp.ProcessProposal(req)
}

// but before committing it to the consensus state.
func (app BeaconApp) FinalizeBlock(
	req *abci.FinalizeBlockRequest,
) (*abci.FinalizeBlockResponse, error) {
	start := time.Now()
	defer func() {
		app.Logger().
			Info("finalizedBlock executed",
				"duration", time.Since(start).String())
	}()
	return app.BaseApp.FinalizeBlock(req)
}

// Commit is our custom implementation of the ABCI method Commit.
func (app BeaconApp) Commit() (*abci.CommitResponse, error) {
	start := time.Now()
	defer func() {
		app.Logger().
			Info("commit executed",
				"duration", time.Since(start).String())
	}()
	return app.BaseApp.Commit()
}
