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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
)

// This tests a reth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It currently fails here.
//
// BUG with optimistic ON:
//   - block height N, round i
//   - reth EL builds payload for height N
//   - reth accepts process proposal (N, i)
//   - marking N/N - 1 as "head"/"finalized" to the EL
//   - supermajority does not finalize (N, i)
//   - block height N, next round i + 1
//   - **BUG**: beaconKit purges payloadID for height N and reth EL does NOT store the already built payload.
//   - any EL who accepted proposal (N, i) returns empty (comet block has 0 txs when it should have 2)
//     for proposal (N, i + 1) because engine API returns "Received nil payload ID on VALID engine response"
//   - since N - 1 is "finalized" to the EL, asking to re-build payload N by sending FCU with
//     N - 1 as "head" block, EL always returns nil.
//
// TODO: remove from PectraForkSuite.
func (s *PectraForkSuite) TestReth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain(s.T()) // 1 reth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Reth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(1 * time.Second)
	consensusTime = time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2) // FAILS HERE

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
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

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
// validator is chosen to propose at a subsequent round and this block is finalized.
//
// TODO: remove from PectraForkSuite.
func (s *PectraForkSuite) TestGeth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Geth.InitializeChain(s.T()) // 1 geth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(1 * time.Second)
	consensusTime = time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2) // FAILS HERE

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
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

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

func (s *PectraForkSuite) TestReth_MultiplePayloadRebuilds() {
	// Initialize the chain state.
	s.Reth.InitializeChain(s.T()) // 1 reth validator
	s.Geth.InitializeChain(s.T()) // helper node
	helpBuilder := s.Geth

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Reth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	const blkHeight = int64(1)
	var (
		specs         = s.Reth.TestNode.ChainSpec
		consensusTime = time.Unix(int64(specs.ElectraForkTime()), 0)

		validTxsHeight1 [][]byte
	)

	{
		// 1- Build a valid block at height 1, via the helpBuilder
		prepareRequest := &types.PrepareProposalRequest{
			Height:          blkHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := helpBuilder.SimComet.Comet.PrepareProposal(helpBuilder.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)
		validTxsHeight1 = proposal.Txs

		// 2- Process the block via s.Reth node. The proposal is expected
		// to pass and start building payload for height 2, optimistically.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          blkHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	// For some reason, the supermajority does not finalize the block.
	// Another block comes, still at height 1, this time *invalid*. This would force
	// the node to rebuild height 1, which the EL cannot do since it has already received
	// an FCU(head == block_at_height_2)
	{
		invalidTxs := buildInvalidBlock(s, helpBuilder, validTxsHeight1, blkHeight, consensusTime)

		// 3- Process the invalid proposal proposal. It will be rejected
		// and attempt to build optimistically a block at height 1.
		processRequest := &types.ProcessProposalRequest{
			Txs:             invalidTxs,
			Height:          blkHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT.String(), processResp.Status.String())
	}

	{
		// 4- Finally let reth node build block at height 1, process and finalize it
		prepareRequest := &types.PrepareProposalRequest{
			Height:          blkHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal via s.Reth node. The proposal is expected
		// to pass and start building payload for height 2, optimistically.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          blkHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          blkHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// it's key to build an invalid block that fails VerifyIncomingBlock
// but passes other checks (signature, version, etc)
func buildInvalidBlock(
	s *PectraForkSuite,
	builder simulated.SharedAccessors,
	txs [][]byte,
	nextBlockHeight int64,
	consensusTime time.Time,
) [][]byte {
	blsSigner := simulated.GetBlsSigner(builder.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	signedBlk, sidecars, err := builder.SimComet.Comet.Blockchain.ParseBeaconBlock(
		&types.ProcessProposalRequest{
			Txs:             txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		},
	)
	s.Require().NoError(err)
	blk := signedBlk.BeaconBlock

	blk.Body.RandaoReveal = [96]byte{'t', 'e', 's', 't'}
	reSignedBlk, err := ctypes.NewSignedBeaconBlock(
		blk,
		ctypes.NewForkData(
			builder.TestNode.ChainSpec.ActiveForkVersionForTimestamp(blk.GetTimestamp()),
			builder.GenesisValidatorsRoot,
		),
		builder.TestNode.ChainSpec,
		blsSigner,
	)
	s.Require().NoError(err)

	signedBlkBytes, bbErr := reSignedBlk.MarshalSSZ()
	s.Require().NoError(bbErr)
	txs[0] = signedBlkBytes

	sidecarsBytes, scErr := sidecars.MarshalSSZ()
	s.Require().NoError(scErr)
	txs[1] = sidecarsBytes

	return txs
}
