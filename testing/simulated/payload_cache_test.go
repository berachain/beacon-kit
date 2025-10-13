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

// This tests a reth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PectraForkSuite) TestReth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain2Validators(s.T()) // 1 reth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Reth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from reth.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}

		// Process the proposal.
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a geth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PectraForkSuite) TestGeth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Geth.InitializeChain2Validators(s.T()) // 1 geth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from geth.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}

		// Process the proposal.
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Geth.SimComet.Comet.FinalizeBlock(s.Geth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Geth.SimComet.Comet.Commit(s.Geth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a reth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PectraForkSuite) TestReth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain2Validators(s.T()) // 1 reth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Reth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            faultyConsensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		// As we reject our own built proposal, bkit should evict this payload from its cache.
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Reth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}

		// Process the proposal.
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a geth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PectraForkSuite) TestGeth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Geth.InitializeChain2Validators(s.T()) // 1 geth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            faultyConsensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		// As we reject our own built proposal, bkit should evict this payload from its cache.
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}

		// Process the proposal.
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Geth.SimComet.Comet.FinalizeBlock(s.Geth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Geth.SimComet.Comet.Commit(s.Geth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}
