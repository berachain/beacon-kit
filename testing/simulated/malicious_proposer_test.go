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
	"math/big"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/engine-primitives/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// TestProcessProposal_BadBlock_IsRejected effectively serves as a test for how a valid node would react to
// a malicious proposer proposing a block with an invalid EVM transaction.
func (s *SimulatedSuite) TestProcessProposal_BadBlock_IsRejected() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals := s.moveChainToHeight(blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)

	currentHeight := int64(blockHeight + coreLoopIterations)
	// Prepare a block proposal.
	proposalTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            proposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForSlot(math.Slot(currentHeight)),
	)
	s.Require().NoError(err)

	// Sign a malicious transaction that is expected to fail.
	recipientAddress := gethcommon.HexToAddress("0x56898d1aFb10cad584961eb96AcD476C6826e41E")
	maliciousTx, err := gethtypes.SignNewTx(
		simulated.GetTestKey(s.T()),
		gethtypes.NewCancunSigner(big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))),
		&gethtypes.DynamicFeeTx{
			Nonce:     0,
			To:        &recipientAddress,
			Value:     big.NewInt(0),
			Gas:       21016,
			GasTipCap: big.NewInt(10000000),
			GasFeeCap: big.NewInt(10000000),
			Data:      []byte{},
		},
	)

	// Initialize the slice with the malicious transaction.
	maliciousTxs := []*gethprimitives.Transaction{maliciousTx}

	// Create a malicious block by injecting an invalid transaction.
	maliciousBlock := simulated.ComputeAndSetInvalidExecutionBlock(s.T(), proposedBlock.GetMessage(), s.TestNode.ChainSpec, maliciousTxs)

	// Re-sign the block
	maliciousBlockSigned, err := ctypes.NewSignedBeaconBlock(
		maliciousBlock,
		&ctypes.ForkData{
			CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForSlot(maliciousBlock.GetSlot()),
			GenesisValidatorsRoot: s.GenesisValidatorsRoot,
		},
		s.TestNode.ChainSpec,
		blsSigner,
	)
	s.Require().NoError(err)

	maliciousBlockBytes, err := maliciousBlockSigned.MarshalSSZ()
	s.Require().NoError(err)

	// Replace the valid block with the malicious block in the proposal.
	proposal.Txs[0] = maliciousBlockBytes

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the malicious block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)

	// Verify that the log contains the expected error message.
	s.Require().Contains(s.LogBuffer.String(), errors.ErrInvalidPayloadStatus.Error())
	// Note this error message may change across execution clients. Base fee changes with number of core loop iterations.
	s.Require().Contains(s.LogBuffer.String(), "max fee per gas less than block base fee: address 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4, maxFeePerGas: 10000000, baseFee: 765625000")
}

// TestProcessProposal_InvalidTimestamps_Errors effectively serves as a test for how a valid node would react to
// a malicious proposer attempting to use a future timestamp in the block that does not match the consensus timestamp.
func (s *SimulatedSuite) TestProcessProposal_InvalidTimestamps_Errors() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals := s.moveChainToHeight(blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)
	currentHeight := int64(blockHeight + coreLoopIterations)

	// Prepare a block proposal, but 2 seconds in the future (i.e. attempt to roll timestamp forward)
	correctConsensusTime := time.Now()
	maliciousProposalTime := correctConsensusTime.Add(2 * time.Second)
	maliciousProposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            maliciousProposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(maliciousProposal)

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the malicious block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:             maliciousProposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
		// Use the correct time as the actual consensus time, which mismatches the proposal time.
		Time: correctConsensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), payloadtime.ErrTooFarInTheFuture.Error())
}
