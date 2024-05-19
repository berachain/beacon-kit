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
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
)

// PrepareProposal is called by the consensus engine to prepare a proposal for
// the next block.
func (r *BeaconKitRuntime[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
	BeaconBlockBodyT,
	StorageBackendT,
]) PrepareProposal(
	req *abci.PrepareProposalRequest,
	// TODO: This nextPrepareProposal function is temporary.
	nextPrepareProposal func(
		*abci.PrepareProposalRequest,
	) (*abci.PrepareProposalResponse, error),
) (*abci.PrepareProposalResponse, error) {
	start := time.Now()
	defer func() {
		r.logger.
			Info("prepare-proposal executed",
				"duration", time.Since(start).String())
	}()
	return nextPrepareProposal(req)
}

// ProcessProposal is called by the consensus engine when a new proposal block
// is received.
func (r *BeaconKitRuntime[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
	BeaconBlockBodyT,
	StorageBackendT,
]) ProcessProposal(
	req *abci.ProcessProposalRequest,
	// TODO: This nextProcessProposal function is temporary.
	nextProcessProposal func(
		*abci.ProcessProposalRequest,
	) (*abci.ProcessProposalResponse, error),
) (*abci.ProcessProposalResponse, error) {
	start := time.Now()
	defer func() {
		r.logger.
			Info("process-proposal executed",
				"duration", time.Since(start).String())
	}()
	return nextProcessProposal(req)
}

// but before committing it to the consensus state.
func (r *BeaconKitRuntime[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
	BeaconBlockBodyT,
	StorageBackendT,
]) FinalizeBlock(
	req *abci.FinalizeBlockRequest,
	// TODO: This nextFinalizeBlock function is temporary.
	nextFinalizeBlock func(
		*abci.FinalizeBlockRequest,
	) (*abci.FinalizeBlockResponse, error),
) (*abci.FinalizeBlockResponse, error) {
	start := time.Now()
	defer func() {
		r.logger.
			Info("finalized-block executed",
				"duration", time.Since(start).String())
	}()
	return nextFinalizeBlock(req)
}

// Commit is our custom implementation of the ABCI method Commit.
func (r *BeaconKitRuntime[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
	BeaconBlockBodyT,
	StorageBackendT,
	// TODO: This nextCommit function is temporary.
]) Commit(nextCommit func() (
	*abci.CommitResponse, error,
)) (*abci.CommitResponse, error) {
	start := time.Now()
	defer func() {
		r.logger.
			Info("commit executed",
				"duration", time.Since(start).String())
	}()
	return nextCommit()
}
