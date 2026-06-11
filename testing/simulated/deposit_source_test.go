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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package simulated_test

import (
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/testing/simulated"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

// TestBodyDepositsAfterFuluRejected verifies that, once Fulu (Osaka) is active and the
// pre-Fulu deposit queue has been drained, a proposed block carrying deposits on the
// beacon block body is rejected. From Fulu onwards deposits must be sourced exclusively
// from EIP-6110 execution requests, so the only valid block body deposit source is the
// single first-Fulu catchup block.
//
// Chain spec (ProvideFuluDepositTestChainSpec): Electra1 at t=6, Fulu at t=7. Block 3
// (t=7) is the first Fulu block; block 4 (t=8) is the first block where deposits on the
// block body are no longer a valid source.
func (s *FuluDepositSuite) TestBodyDepositsAfterFuluRejected() {
	s.InitializeChain(s.T(), 1)

	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Advance through the first Fulu block (block 3 at t=7) so that the next block is a
	// non-first Fulu block where block body deposits are disallowed.
	const postFuluHeight = int64(4)
	_, _, nextBlockTime := s.MoveChainToHeight(s.T(), 1, postFuluHeight-1, nodeAddress, time.Unix(5, 0))
	s.Require().Equal(time.Unix(8, 0), nextBlockTime, "block 4 must be the first post-Fulu block")

	// Prepare a valid block proposal for the post-Fulu height.
	validProposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &cmtabci.PrepareProposalRequest{
		Height:          postFuluHeight,
		Time:            nextBlockTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(validProposal)

	// Inject a deposit onto the beacon block body. The execution payload is left untouched,
	// so the block is only invalid because it sources a deposit from the body after Fulu.
	maliciousTxs := testBuildInvalidBlock(
		s.Require(),
		s.SharedAccessors,
		&cmtabci.PrepareProposalRequest{
			Txs:    validProposal.Txs,
			Height: postFuluHeight,
			Time:   nextBlockTime,
		},
		func(sb *ctypes.SignedBeaconBlock) {
			sb.BeaconBlock.Body.SetDeposits(ctypes.Deposits{{Index: 99}})
		},
	)

	s.LogBuffer.Reset()
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &cmtabci.ProcessProposalRequest{
		Txs:             maliciousTxs,
		Height:          postFuluHeight,
		ProposerAddress: nodeAddress,
		Time:            nextBlockTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(cmtabci.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), core.ErrUnexpectedDepositSource.Error())
}
