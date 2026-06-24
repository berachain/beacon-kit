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
	"bytes"
	"math/big"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

// TestFullLifecycle_ValidBlock_IsSuccessful tests that a valid block proposal is processed, finalized, and committed.
// It loops through this core process `coreLoopIterations` times.
func (s *SimulatedSuite) TestFullLifecycle_ValidBlock_IsSuccessful() {
	const blockHeight = 1
	const coreLoopIterations = 10

	// Initialize the chain state.
	s.InitializeChain(s.T(), 1)
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Test happens post Deneb1 fork.
	startTime := time.Now()

	// iterate through the core loop `coreLoopIterations` times, i.e. Propose, Process, Finalize and Commit.
	proposals, _, _ := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, nodeAddress, startTime)

	// We expect that the number of proposals that were finalized should be `coreLoopIterations`.
	s.Require().Len(proposals, coreLoopIterations)

	currentHeight := int64(blockHeight + coreLoopIterations)
	// Validate post-commit state.
	queryCtx, err := s.SimComet.CreateQueryContext(currentHeight-1, false)
	s.Require().NoError(err)

	stateDB := s.TestNode.StorageBackend.StateFromContext(queryCtx)
	slot, err := stateDB.GetSlot()
	s.Require().NoError(err)
	s.Require().Equal(math.U64(currentHeight-1), slot)

	stateHeader, err := stateDB.GetLatestBlockHeader()
	s.Require().NoError(err)

	lph, err := stateDB.GetLatestExecutionPayloadHeader()
	s.Require().NoError(err)

	// Unmarshal the beacon block from the ABCI request.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposals[len(proposals)-1].Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(lph.GetTimestamp()),
	)
	s.Require().NoError(err)
	s.Require().Equal(proposedBlock.GetHeader().GetBodyRoot(), stateHeader.GetBodyRoot())
}

// TestFullLifecycle_ValidBlockWithInjectedTransaction_IsSuccessful effectively serves as a demonstration for how one can
// inject custom transactions and state transitions into the core loop.
func (s *SimulatedSuite) TestFullLifecycle_ValidBlockWithInjectedTransaction_IsSuccessful() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T(), 1)

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Build at the first block, on Deneb (pre Deneb1 fork at t=30). The proposer's EL build for the
	// first block is fresh (no cached optimistic payload from a prior round), so the injected tx
	// sitting in the txpool is included in the built block.
	currentHeight := int64(blockHeight)
	consensusTime := time.Unix(1, 0)

	// Since reth cannot use eth_simulateV1-based construction to build a valid payload,
	// submit a valid transaction to the EL txpool and let the proposer include it in the block.
	recipientAddress := gethcommon.HexToAddress("0x56898d1aFb10cad584961eb96AcD476C6826e41E")
	validTx, err := gethtypes.SignNewTx(
		simulated.GetTestKey(s.T()),
		gethtypes.NewCancunSigner(big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))),
		&gethtypes.DynamicFeeTx{
			Nonce:     0,
			To:        &recipientAddress,
			Value:     big.NewInt(0),
			Gas:       21016,
			GasTipCap: big.NewInt(10_000_000_000),
			GasFeeCap: big.NewInt(10_000_000_000),
			Data:      []byte{},
		},
	)
	s.Require().NoError(err)
	s.Require().NoError(s.TestNode.ContractBackend.SendTransaction(s.CtxApp, validTx))

	validTxBytes, err := validTx.MarshalBinary()
	s.Require().NoError(err)
	forkVersion := s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(consensusTime.Unix()))

	// Build the proposal, retrying until the proposer's EL build includes the injected tx.
	proposal := s.PrepareProposalUntil(s.T(), &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            consensusTime,
		ProposerAddress: nodeAddress,
	}, func(p *types.PrepareProposalResponse) bool {
		builtBlock, derr := encoding.UnmarshalBeaconBlockFromABCIRequest(
			p.Txs, blockchain.BeaconBlockTxIndex, forkVersion,
		)
		s.Require().NoError(derr)
		for _, txBytes := range builtBlock.GetBody().GetExecutionPayload().GetTransactions() {
			if bytes.Equal(txBytes, validTxBytes) {
				return true
			}
		}
		return false
	})

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the valid block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Finalize the block.
	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(finalizeResp)

	// Commit the block.
	_, err = s.SimComet.Comet.Commit(s.CtxComet, &types.CommitRequest{})
	s.Require().NoError(err)
}

// TestFullLifecycle_ValidBlockAndInjectedBlob_IsSuccessful tests that a valid block and blob and proposal is processed, finalized, and committed.
func (s *SimulatedSuite) TestFullLifecycle_ValidBlockAndInjectedBlob_IsSuccessful() {
	const blockHeight = 1

	// Initialize the chain state.
	s.InitializeChain(s.T(), 1)

	// Retrieve the proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Since reth cannot use eth_simulateV1-based construction to build a valid payload,
	// submit a valid transaction to the EL txpool and let the proposer include it in the block.
	currentHeight := int64(blockHeight)
	consensusTime := time.Unix(1, 0)

	// Create the blobs with proofs and commitments. Each blob goes into its own blob transaction.
	blobs := []*eip4844.Blob{{1, 2, 3}, {4, 5, 6}}
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(require.New(s.T()), blobs, s.TestNode.KZGVerifier)
	s.Require().Len(proofs, len(blobs))
	s.Require().Len(commitments, len(blobs))

	// Sign the blob transactions and submit them to the EL txpool.
	for i := range blobs {
		blobHash := commitments[i].ToVersionedHash()
		txSidecar := &gethtypes.BlobTxSidecar{
			Blobs:       []kzg4844.Blob{kzg4844.Blob(blobs[i][:])},
			Commitments: []kzg4844.Commitment{kzg4844.Commitment(commitments[i])},
			Proofs:      []kzg4844.Proof{kzg4844.Proof(proofs[i])},
		}
		blobTx, txErr := gethtypes.SignNewTx(
			simulated.GetTestKey(s.T()),
			gethtypes.NewCancunSigner(big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))),
			&gethtypes.BlobTx{
				Nonce:      uint64(i),
				GasTipCap:  uint256.NewInt(10_000_000_000),
				GasFeeCap:  uint256.NewInt(10_000_000_000),
				Gas:        210000,
				Value:      uint256.NewInt(0),
				Data:       []byte{},
				AccessList: nil,
				BlobFeeCap: uint256.NewInt(10_000_000_000),
				BlobHashes: []gethcommon.Hash{blobHash},
				Sidecar:    nil,
			},
		)
		s.Require().NoError(txErr)
		blobTx = blobTx.WithBlobTxSidecar(txSidecar)
		s.Require().NoError(s.TestNode.ContractBackend.SendTransaction(s.CtxApp, blobTx))
	}

	// Build the proposal, retrying until the proposer's EL build includes both blobs.
	proposal := s.PrepareProposalUntil(s.T(), &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            consensusTime,
		ProposerAddress: nodeAddress,
	}, func(p *types.PrepareProposalResponse) bool {
		sidecars, scErr := encoding.UnmarshalBlobSidecarsFromABCIRequest(p.Txs, blockchain.BlobSidecarsTxIndex)
		s.Require().NoError(scErr)
		return len(sidecars) == len(blobs)
	})

	// Reset the log buffer to discard old logs we don't care about.
	s.LogBuffer.Reset()
	// Process the proposal containing the valid block and blobs.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Finalize the block.
	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(finalizeResp)

	// Commit the block.
	_, err = s.SimComet.Comet.Commit(s.CtxComet, &types.CommitRequest{})
	s.Require().NoError(err)
}
