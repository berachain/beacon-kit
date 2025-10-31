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
	dablob "github.com/berachain/beacon-kit/da/blob"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/engine-primitives/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/observability/metrics/discard"
	"github.com/berachain/beacon-kit/primitives/common"
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

// TestProcessProposal_BadBlock_IsRejected effectively serves as a test for how a valid node would react to
// a malicious proposer proposing a block with an invalid EVM transaction.
func (s *SimulatedSuite) TestProcessProposal_BadBlock_IsRejected() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Test happens on Deneb, pre Deneb1 fork.
	startTime := time.Unix(0, 0)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals, _, proposalTime := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, nodeAddress, startTime)
	s.Require().Len(proposals, coreLoopIterations)

	// We expected this test to happen during Pre-Deneb1 fork.
	currentHeight := int64(blockHeight + coreLoopIterations)

	// Prepare a block proposal.
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            proposalTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(proposalTime.Unix())),
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
	maliciousBlock := simulated.ComputeAndSetInvalidExecutionBlock(
		s.T(), proposedBlock.GetBeaconBlock(), s.TestNode.ChainSpec, maliciousTxs, nil,
	)

	// Re-sign the block
	maliciousBlockSigned, err := ctypes.NewSignedBeaconBlock(
		maliciousBlock,
		&ctypes.ForkData{
			CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(maliciousBlock.GetTimestamp()),
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
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
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
	s.InitializeChain(s.T())
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Test happens post Deneb1 fork.
	startTime := time.Now()

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals, _, correctConsensusTime := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, nodeAddress, startTime)
	s.Require().Len(proposals, coreLoopIterations)
	currentHeight := int64(blockHeight + coreLoopIterations)

	// Prepare a block proposal. This will create a valid payload due to optimistic payload building.
	// It is called to flush the payload cache.
	validProposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            correctConsensusTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(validProposal)

	// Build a block with a malicious proposal time
	maliciousPayloadTime := correctConsensusTime.Add(2 * time.Second)
	maliciousProposalTxs := testBuildInvalidBlock(
		s.Require(),
		s.SharedAccessors,
		&types.PrepareProposalRequest{
			Txs:    validProposal.Txs,
			Height: currentHeight,
			Time:   correctConsensusTime,
		},
		func(sb *ctypes.SignedBeaconBlock) {
			blk := sb.BeaconBlock
			blk.Body.ExecutionPayload.Timestamp = math.U64(maliciousPayloadTime.Unix())
		},
	)

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the malicious block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             maliciousProposalTxs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		// Use the correct time as the actual consensus time, which mismatches the proposal time.
		Time: correctConsensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), payloadtime.ErrTooFarInTheFuture.Error())
}

// TestProcessProposal_InvalidBlobCommitment_Errors effectively serves as a test for a malicious blobs.
// Specifically, to bypass the commitment count check, we put 2 commitments into 1 tx and leave the
// other tx with zero commitments. This ultimately gets rejected by the execution client.
func (s *SimulatedSuite) TestProcessProposal_InvalidBlobCommitment_Errors() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Test happens on Deneb, pre Deneb1 fork.
	startTime := time.Unix(0, 0)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals, _, consensusTime := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, nodeAddress, startTime)
	s.Require().Len(proposals, coreLoopIterations)

	// We expected this test to happen during Pre-Deneb1 fork.
	currentHeight := int64(blockHeight + coreLoopIterations)

	// Prepare a block proposal.
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            consensusTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(consensusTime.Unix())),
	)
	s.Require().NoError(err)

	// Create the Blobs, with proofs and commitments
	// Each blob will go into 1 transaction.
	blobs := []*eip4844.Blob{{1, 2, 3}, {4, 5, 6}}
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(require.New(s.T()), blobs, s.TestNode.KZGVerifier)
	s.Require().Len(proofs, len(blobs))
	s.Require().Len(commitments, len(blobs))

	// Here is where we act malicious.
	blobVersionedHash0 := commitments[0].ToVersionedHash()
	blobVersionedHash1 := commitments[1].ToVersionedHash()
	versionedHashesForBlob := [][]gethcommon.Hash{
		{blobVersionedHash0, blobVersionedHash1}, // index 0
		nil,                                      // index 1 is nil
	}
	// Sign blob transactions
	blobTxs := make([]*gethtypes.Transaction, len(blobs))
	for i := range blobs {
		blobCommitment := commitments[i]
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
				BlobHashes: versionedHashesForBlob[i],
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

	proposedBlockMessage, err = simulated.ComputeAndSetStateRoot(queryCtx, consensusTime, nodeAddress, s.TestNode.StateProcessor, s.TestNode.StorageBackend, proposedBlockMessage)
	s.Require().NoError(err)

	newSignedBlock, err := ctypes.NewSignedBeaconBlock(
		proposedBlockMessage,
		&ctypes.ForkData{
			CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(proposedBlockMessage.GetTimestamp()),
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
	sidecarFactory := dablob.NewSidecarFactory(dablob.NewFactoryMetrics(discard.NewFactory()))
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
	// Inject the sidecar
	sidecarBytes, err := sidecars.MarshalSSZ()
	s.Require().NoError(err)

	proposal.Txs[1] = sidecarBytes

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), "could not apply tx 1 [0xdbbf691e9271a8bc3de5e64405337972fb4a5185cc3df160bac310c515f7d768]: blob transaction missing blob hashes")
}

// TestProcessProposal_InvalidBlobInclusionProof_Errors effectively serves as a test for a malicious blobs.
// Specifically, we tweak the KZG commitment inclusion proof such that it is invalid and must be rejected.
func (s *SimulatedSuite) TestProcessProposal_InvalidBlobInclusionProof_Errors() {
	const blockHeight = 1
	const coreLoopIterations = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Test happens on Deneb, pre Deneb1 fork.
	startTime := time.Unix(0, 0)

	// Go through 1 iteration of the core loop to bypass any startup specific edge cases such as sync head on startup.
	proposals, _, consensusTime := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, nodeAddress, startTime)
	s.Require().Len(proposals, coreLoopIterations)

	// We expected this test to happen during Pre-Deneb1 fork.
	currentHeight := int64(blockHeight + coreLoopIterations)

	// Prepare a block proposal.
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            consensusTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(consensusTime.Unix())),
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

	proposedBlockMessage, err = simulated.ComputeAndSetStateRoot(queryCtx, consensusTime, nodeAddress, s.TestNode.StateProcessor, s.TestNode.StorageBackend, proposedBlockMessage)
	s.Require().NoError(err)

	newSignedBlock, err := ctypes.NewSignedBeaconBlock(
		proposedBlockMessage,
		&ctypes.ForkData{
			CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(proposedBlockMessage.GetTimestamp()),
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
	sidecarFactory := dablob.NewSidecarFactory(dablob.NewFactoryMetrics(discard.NewFactory()))
	for i := range blobs {
		inclusionProof, err := sidecarFactory.BuildKZGInclusionProof(proposedBlockMessage.GetBody(), math.U64(i))
		s.Require().NoError(err)
		// Malicious point: We tweak the inclusion proof
		inclusionProof[len(inclusionProof)-1] = common.Root{}
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
	// Inject theI sidecar
	sidecarBytes, err := sidecars.MarshalSSZ()
	s.Require().NoError(err)

	proposal.Txs[1] = sidecarBytes

	// Reset the log buffer to discard old logs we don't care about
	s.LogBuffer.Reset()
	// Process the proposal containing the block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          currentHeight,
		ProposerAddress: nodeAddress,
		Time:            consensusTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), "invalid KZG commitment inclusion proof")
}

// testBuildInvalidBlock builds an invalid block, modifying it along what it is
// specified in modifyBlock function. It is then properly resigned.
func testBuildInvalidBlock(
	r *require.Assertions,
	builder simulated.SharedAccessors,
	PrepReq *types.PrepareProposalRequest,
	modifyBlock func(*ctypes.SignedBeaconBlock),
) [][]byte {
	blsSigner := simulated.GetBlsSigner(builder.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	r.NoError(err)

	signedBlk, sidecars, err := builder.SimComet.Comet.Blockchain.ParseProcessProposalRequest(
		&types.ProcessProposalRequest{
			Txs:             PrepReq.Txs,
			Height:          PrepReq.Height,
			ProposerAddress: pubkey.Address(),
			Time:            PrepReq.Time,
		},
	)
	r.NoError(err)

	modifyBlock(signedBlk)

	blk := signedBlk.BeaconBlock
	reSignedBlk, err := ctypes.NewSignedBeaconBlock( // resign to make sure signature checks pass
		blk,
		ctypes.NewForkData(
			builder.TestNode.ChainSpec.ActiveForkVersionForTimestamp(blk.GetTimestamp()),
			builder.GenesisValidatorsRoot,
		),
		builder.TestNode.ChainSpec,
		blsSigner,
	)
	r.NoError(err)

	signedBlkBytes, bbErr := reSignedBlk.MarshalSSZ()
	r.NoError(bbErr)

	res := make([][]byte, len(PrepReq.Txs))
	res[0] = signedBlkBytes

	sidecarsBytes, scErr := sidecars.MarshalSSZ()
	r.NoError(scErr)
	res[1] = sidecarsBytes

	return res
}
