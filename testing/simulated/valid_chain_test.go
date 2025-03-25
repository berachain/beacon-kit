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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	dablob "github.com/berachain/beacon-kit/da/blob"
	datypes "github.com/berachain/beacon-kit/da/types"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/crypto"
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
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// iterate through the core loop `coreLoopIterations` times, i.e. Propose, Process, Finalize and Commit.
	proposals := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner)

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
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(lph.GetTimestamp().Unwrap()),
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
	consensusTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            consensusTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(uint64(consensusTime.Unix())),
	)
	s.Require().NoError(err)

	// Sign a valid transaction that is expected to pass
	recipientAddress := gethcommon.HexToAddress("0x56898d1aFb10cad584961eb96AcD476C6826e41E")
	validTx, err := gethtypes.SignNewTx(
		simulated.GetTestKey(s.T()),
		gethtypes.NewCancunSigner(big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))),
		&gethtypes.DynamicFeeTx{
			Nonce:     0,
			To:        &recipientAddress,
			Value:     big.NewInt(0),
			Gas:       21016,
			GasTipCap: big.NewInt(765625000),
			GasFeeCap: big.NewInt(765625000),
			Data:      []byte{},
		},
	)

	validTxs := []*gethprimitives.Transaction{validTx}
	// Create a new beacon block with the valid transaction.
	// Note: The beacon block returned here has an incorrect beacon state root, which is fixed in `ComputeAndSetStateRoot`.
	unsignedBlock := simulated.ComputeAndSetValidExecutionBlock(s.T(), proposedBlock.GetBeaconBlock(), s.SimulationClient, s.TestNode.ChainSpec, validTxs)

	proposerAddress, err := crypto.GetAddressFromPubKey(blsSigner.PublicKey())
	s.Require().NoError(err)

	// Finalize the block by applying the state transition to update its state root.
	queryCtx, err := s.SimComet.CreateQueryContext(currentHeight-1, false)
	s.Require().NoError(err)
	finalBlock, err := simulated.ComputeAndSetStateRoot(queryCtx, consensusTime, proposerAddress, s.TestNode.StateProcessor, s.TestNode.StorageBackend, unsignedBlock)
	s.Require().NoError(err)

	newSignedBlock, err := ctypes.NewSignedBeaconBlock(
		finalBlock,
		&ctypes.ForkData{
			CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(unsignedBlock.GetTimestamp().Unwrap()),
			GenesisValidatorsRoot: s.GenesisValidatorsRoot,
		},
		s.TestNode.ChainSpec,
		blsSigner,
	)
	s.Require().NoError(err)

	newBlockBytes, err := newSignedBlock.MarshalSSZ()
	s.Require().NoError(err)

	// Replace the old block with the new block in the proposal.
	proposal.Txs[0] = newBlockBytes

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the valid block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Finalize the block.
	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
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
	consensusTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            consensusTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(uint64(consensusTime.Unix())),
	)
	s.Require().NoError(err)

	// Create the Blobs, with proofs and commitments
	// Each blob will go into 1 transaction.
	blobs := []*eip4844.Blob{{1, 2, 3}, {4, 5, 6}}
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(require.New(s.T()), blobs, s.TestNode.KZGVerifier)
	s.Require().Len(proofs, len(blobs))
	s.Require().Len(commitments, len(blobs))

	// Sign blob transactions
	blobTxs := make([]*gethtypes.Transaction, len(blobs))
	for i := range blobs {
		blobCommitment := commitments[i]
		blobHash := blobCommitment.ToVersionedHash()
		txSidecar := &gethtypes.BlobTxSidecar{
			Blobs:       []kzg4844.Blob{kzg4844.Blob(blobs[i][:])},
			Commitments: []kzg4844.Commitment{kzg4844.Commitment(blobCommitment)},
			Proofs:      []kzg4844.Proof{kzg4844.Proof(proofs[i])},
		}
		blobTx, err := gethtypes.SignNewTx(
			simulated.GetTestKey(s.T()),
			gethtypes.NewCancunSigner(big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))),
			&gethtypes.BlobTx{
				Nonce: uint64(i),
				// Set to 875000000 as that is the tx base fee
				GasTipCap: uint256.NewInt(875000000),
				GasFeeCap: uint256.NewInt(875000000),
				// Set to 21000 for minimum intrinsic gas
				Gas:        210000,
				Value:      uint256.NewInt(0),
				Data:       []byte{},
				AccessList: nil,
				BlobFeeCap: uint256.NewInt(10),
				// If we have 1 tx with multiple blobs, we must add the blob hashes here.
				BlobHashes: []gethcommon.Hash{blobHash},
				// Sidecar must be set to nil here or Geth will error with "unexpected blob sidecar in transaction"
				Sidecar: nil,
			},
		)
		s.Require().NoError(err)
		// Once we've signed the Tx, we tag the blob with the tx purely for association between tx and sidecars.
		// In this case, each 1 tx has a sidecar with 1 blob, even though 1 tx could have more than 1 blob.
		blobTx = blobTx.WithBlobTxSidecar(txSidecar)
		blobTxs[i] = blobTx
	}

	proposedBlockMessage := simulated.ComputeAndSetValidExecutionBlock(
		s.T(),
		proposedBlock.GetBeaconBlock(),
		s.SimulationClient,
		s.TestNode.ChainSpec,
		blobTxs,
	)
	proposedBlockMessage.GetBody().SetBlobKzgCommitments(commitments)

	// Finalize the block by applying the state transition to update its state root.
	queryCtx, err := s.SimComet.CreateQueryContext(currentHeight-1, false)
	s.Require().NoError(err)

	// Retrieve the BLS signer and proposer address.
	proposerAddress, err := crypto.GetAddressFromPubKey(blsSigner.PublicKey())
	s.Require().NoError(err)

	proposedBlockMessage, err = simulated.ComputeAndSetStateRoot(queryCtx, consensusTime, proposerAddress, s.TestNode.StateProcessor, s.TestNode.StorageBackend, proposedBlockMessage)
	s.Require().NoError(err)

	newSignedBlock, err := ctypes.NewSignedBeaconBlock(
		proposedBlockMessage,
		&ctypes.ForkData{
			CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(proposedBlockMessage.GetTimestamp().Unwrap()),
			GenesisValidatorsRoot: s.GenesisValidatorsRoot,
		},
		s.TestNode.ChainSpec,
		blsSigner,
	)
	s.Require().NoError(err)

	// Inject the new block
	newSignedBlockBytes, err := newSignedBlock.MarshalSSZ()
	s.Require().NoError(err)
	proposal.Txs[0] = newSignedBlockBytes

	// Create the beaconBlock Header for the sidecar
	blockWithCommitmentsSignedHeader := ctypes.NewSignedBeaconBlockHeader(
		newSignedBlock.GetHeader(),
		newSignedBlock.GetSignature(),
	)

	sidecarsSlice := make([]*datypes.BlobSidecar, len(blobs))
	// Build Inclusion Proofs for Sidecars
	sidecarFactory := dablob.NewSidecarFactory(metrics.NewNoOpTelemetrySink())
	for i := range blobs {
		inclusionProof, err := sidecarFactory.BuildKZGInclusionProof(proposedBlockMessage.GetBody(), math.U64(i))
		s.Require().NoError(err)
		sidecar := datypes.BuildBlobSidecar(
			math.U64(i),
			blockWithCommitmentsSignedHeader,
			blobs[i],
			commitments[i],
			proofs[i],
			inclusionProof,
		)
		sidecarsSlice[i] = sidecar
	}
	sidecars := datypes.BlobSidecars(sidecarsSlice)
	// Inject the valid sidecar
	sidecarBytes, err := sidecars.MarshalSSZ()
	s.Require().NoError(err)

	proposal.Txs[1] = sidecarBytes

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the valid block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Finalize the block.
	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(finalizeResp)

	// Commit the block.
	_, err = s.SimComet.Comet.Commit(s.CtxComet, &types.CommitRequest{})
	s.Require().NoError(err)
}
