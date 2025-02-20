//go:build simulated

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

package simulated_test

import (
	"time"

	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
)

// TestProcessProposal_CrashedExecutionClient_Errors effectively serves as a test for how a valid node would react to
// a valid block being proposed but the execution client has crashed.
func (s *SimulatedSuite) TestProcessProposal_CrashedExecutionClient_Errors() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals := s.CoreLoop(blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)

	// Prepare a valid block proposal.
	proposalTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          blockHeight + coreLoopIterations,
		Time:            proposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Reset the log buffer to discard old logs we don't care about.
	s.LogBuffer.Reset()
	// Kill the execution client.
	err = s.ElHandle.Close()
	s.Require().NoError(err)
	// Process the proposal containing the malicious block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          blockHeight + coreLoopIterations,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), "got an unexpected server error in JSON-RPC response failed to convert from jsonrpc.Error")
}
