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

	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
)

func (s *PectraForkSuite) Test_BuildOptimisticallyAtFork() {
	// Initialize the chain state.
	client := s.Geth
	client.InitializeChain(s.T()) // init the validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(client.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// setup first payload and consensus timestamps so that
	// - it would be pre Electra fork
	// - the second block, built optimistically would be at the electra fork.
	specs := client.TestNode.ChainSpec
	firstBlkConsensusTime := specs.ElectraForkTime() - 1 // before fork
	firstBlkPayloadTime := firstBlkConsensusTime
	secondBlkConsensusTime := specs.ElectraForkTime()
	secondBlkPayloadTime := payloadtime.Next(
		math.U64(secondBlkConsensusTime),
		math.U64(firstBlkPayloadTime),
		true, // this is the formula used while setting second block timestamp optimistically
	)
	s.Require().GreaterOrEqual(secondBlkPayloadTime, math.U64(specs.ElectraForkTime())) // post fork

	nextBlockHeight := int64(1)
	{
		// 1- Build pre-fork block
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            time.Unix(int64(firstBlkConsensusTime), 0),
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := client.SimComet.Comet.PrepareProposal(client.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// 2- Process the proposal. This will trigger am optimistic payload build for block height 2.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(firstBlkConsensusTime), 0),
		}
		processResp, respErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// 3- finalize and commit the first block
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(firstBlkConsensusTime), 0),
		}
		_, finalizeErr := client.SimComet.Comet.FinalizeBlock(client.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := client.SimComet.Comet.Commit(client.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}

	// Now build the next block
	nextBlockHeight++
	{
		// 4- Build post-fork block. Make sure that the fork transition happens
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            time.Unix(int64(secondBlkConsensusTime), 0),
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := client.SimComet.Comet.PrepareProposal(client.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)
		s.Require().Contains(client.LogBuffer.String(), "invoked ProcessFork")

		// 5- Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(secondBlkConsensusTime), 0),
		}
		processResp, respErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// 6- finalize and commit the second block
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(secondBlkConsensusTime), 0),
		}
		_, finalizeErr := client.SimComet.Comet.FinalizeBlock(client.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := client.SimComet.Comet.Commit(client.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}
