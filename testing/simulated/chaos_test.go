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
	"context"
	"time"

	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

// TestProcessProposal_CrashedExecutionClient_Errors effectively serves as a test for how a valid node would react to
// a valid block being proposed but the execution client has crashed.
func (s *SimulatedSuite) TestProcessProposal_CrashedExecutionClient_Errors() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)

	currentHeight := int64(blockHeight + coreLoopIterations)
	// Prepare a valid block proposal.
	proposalTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
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
	// Process the proposal containing the valid block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), client.ErrBadConnection.Error())
}

// TestContextHandling_SIGINT_SafeShutdown mimicks the expected outcome of a SIGINT by calling context cancel and stop services.
func (s *SimulatedSuite) TestContextHandling_SIGINT_SafeShutdown() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Run through core loop iterations to bypass any startup edge cases.
	proposals := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)

	currentHeight := int64(blockHeight + coreLoopIterations)

	s.LogBuffer.Reset()
	// Kill the EL (execution layer)
	err = s.ElHandle.Close()
	s.Require().NoError(err)

	type proposalResult struct {
		proposal *cmtabci.PrepareProposalResponse
		err      error
	}
	// Capture result of prepare proposal
	resultCh := make(chan proposalResult, 1)
	// Prepare proposal in a separate goroutine since it will block due to retrying on the crashed EL.
	proposalTime := time.Now()
	go func() {
		proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
			Height:          currentHeight,
			Time:            proposalTime,
			ProposerAddress: pubkey.Address(),
		})
		resultCh <- proposalResult{
			proposal: proposal,
			err:      err,
		}
	}()

	// Mimic the behavior of the shutdown function when a SIGINT is observed.
	s.CtxAppCancelFn()
	s.TestNode.ServiceRegistry.StopAll()

	// Wait 2 seconds for PrepareProposal to return its result.
	select {
	case res := <-resultCh:
		s.Require().NoError(res.err)
		s.Require().Empty(res.proposal)
		// Shutdown is the last service that is completed and indicates
		s.Require().Contains(s.LogBuffer.String(), "All services stopped")
	case <-time.After(2 * time.Second):
		s.T().Error("PrepareProposal did not finish within 2 seconds after shutdown")
	}
}

// TestContextHandling_CancelledContext_Rejected tests that ABCI requests are rejected if the context is cancelled
func (s *SimulatedSuite) TestContextHandling_CancelledContext_Rejected() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)

	currentHeight := int64(blockHeight + coreLoopIterations)

	// Kill the EL
	err = s.ElHandle.Close()
	s.Require().NoError(err)

	// Cancel the App
	s.CtxAppCancelFn()

	s.LogBuffer.Reset()
	proposalTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            proposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().Empty(proposal)

	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().Error(err, context.Canceled)
	s.Require().Nil(processResp)

	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().Error(err, context.Canceled)
	s.Require().Nil(finalizeResp)

	_, err = s.SimComet.Comet.Commit(s.CtxComet, &types.CommitRequest{})
	s.Require().Error(err, context.Canceled)
}
