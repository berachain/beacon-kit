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
	"context"
	"crypto/sha256"
	"math/big"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	dablob "github.com/berachain/beacon-kit/da/blob"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	mathpkg "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// SimulatedSuite defines our test suite for the simulated Comet component.
type SimulatedSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors

	SimComet              *simulated.SimComet
	LogBuffer             *bytes.Buffer
	GenesisValidatorsRoot common.Root
	SimulationClient      *execution.SimulationClient
}

// TestSimulatedCometComponent runs the test suite.
func TestSimulatedCometComponent(t *testing.T) {
	suite.Run(t, new(SimulatedSuite))
}

// SetupTest initializes the test environment.
func (s *SimulatedSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.Ctx, s.CancelFunc = context.WithCancel(context.Background())
	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), s.HomeDir)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Geth node.
	elNode := execution.NewGethNode(s.HomeDir, execution.ValidGethImage())
	elHandle, authRPC := elNode.Start(s.T())
	s.ElHandle = elHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.LogBuffer = new(bytes.Buffer)
	logger := phuslu.NewLogger(s.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component.
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	s.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     authRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})

	// Retrieve the simulated Comet service.
	var cometService *simulated.SimComet
	err := s.TestNode.FetchService(&cometService)
	s.Require().NoError(err)
	s.Require().NotNil(cometService)
	s.SimComet = cometService

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.TestNode.Start(s.Ctx)
	}()

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.EngineClient)
	// Allow a short period for services to fully initialize.
	time.Sleep(2 * time.Second)
}

// TearDownTest cleans up the test environment.
func (s *SimulatedSuite) TearDownTest() {
	if err := s.ElHandle.Close(); err != nil {
		s.T().Logf("Error closing EL handle: %s", err)
	}
	s.CancelFunc()
}

// initializeChain sets up the chain using the genesis file.
func (s *SimulatedSuite) initializeChain() {
	// Load the genesis state.
	appGenesis, err := genutiltypes.AppGenesisFromFile(s.HomeDir + "/config/genesis.json")
	s.Require().NoError(err)

	// Initialize the chain.
	initResp, err := s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId:       simulated.TestnetBeaconChainID,
		AppStateBytes: appGenesis.AppState,
	})
	s.Require().NoError(err)
	s.Require().Len(initResp.Validators, 1, "Expected 1 validator")

	// Verify that the deposit store contains the expected deposits.
	deposits, err := s.TestNode.StorageBackend.DepositStore().GetDepositsByIndex(
		s.Ctx,
		constants.FirstDepositIndex,
		constants.FirstDepositIndex+s.TestNode.ChainSpec.MaxDepositsPerBlock(),
	)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1, "Expected 1 deposit")
}

// TestFullLifecycle_ValidBlock_IsSuccessful tests that a valid block proposal is processed, finalized, and committed.
// It loops through this core process `coreLoopIterations` times.
func (s *SimulatedSuite) TestFullLifecycle_ValidBlock_IsSuccessful() {
	const blockHeight = 1
	const coreLoopIterations = 10

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// iterate through the core loop `coreLoopIterations` times, i.e. Propose, Process, Finalize and Commit.
	proposals := s.CoreLoop(blockHeight, coreLoopIterations, blsSigner)

	// We expect that the number of proposals that were finalized should be `coreLoopIterations`.
	s.Require().Len(proposals, coreLoopIterations)

	// Validate post-commit state.
	queryCtx, err := s.SimComet.CreateQueryContext(blockHeight+coreLoopIterations-1, false)
	s.Require().NoError(err)

	stateDB := s.TestNode.StorageBackend.StateFromContext(queryCtx)
	slot, err := stateDB.GetSlot()
	s.Require().NoError(err)
	s.Require().Equal(mathpkg.U64(blockHeight+coreLoopIterations-1), slot)

	stateHeader, err := stateDB.GetLatestBlockHeader()
	s.Require().NoError(err)

	// Unmarshal the beacon block from the ABCI request.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposals[len(proposals)-1].Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForSlot(slot),
	)
	s.Require().NoError(err)
	s.Require().Equal(proposedBlock.Message.GetHeader().GetBodyRoot(), stateHeader.GetBodyRoot())
}

// TestFullLifecycle_ValidBlock_IsSuccessful tests that a valid block and proposal is processed, finalized, and committed.
func (s *SimulatedSuite) TestFullLifecycle_ValidBlockAndBlob_IsSuccessful() {
	const blockHeight = 1

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Prepare a valid block proposal.
	proposalTime := time.Unix(0, 0).Add(1 * time.Second)
	proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          blockHeight,
		Time:            proposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForSlot(blockHeight),
	)
	s.Require().NoError(err)

	// Create the Blobs, with proofs and commitments
	blobs := []*eip4844.Blob{{1, 2, 3}, {4, 5, 6}}
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(require.New(s.T()), blobs, s.TestNode.KZGVerifier)
	s.Require().Len(proofs, 2)
	s.Require().Len(commitments, 2)

	// Sign blob transactions
	blobTxs := make([]*gethtypes.Transaction, len(blobs))
	blobTxSidecars := make([]*gethtypes.BlobTxSidecar, len(blobs))
	for i := range blobs {
		blobCommitment := commitments[i]
		blobHash := kzg4844.CalcBlobHashV1(sha256.New(), (*kzg4844.Commitment)(&blobCommitment))
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
				Gas:        21000,
				Value:      nil,
				Data:       nil,
				AccessList: nil,
				// Set to 875000000 as that is the blob base fee
				BlobFeeCap: uint256.NewInt(1),
				BlobHashes: []gethcommon.Hash{blobHash},
				// Sidecar must be set to nil here or Geth will error with "unexpected blob sidecar in transaction"
				Sidecar: nil,
				V:       nil,
				R:       nil,
				S:       nil,
			},
		)
		s.Require().NoError(err)
		blobTxs[i] = blobTx
		blobTxSidecars[i] = txSidecar
	}

	// These are magic value obtained from the geth logs and must be replaced with value from eth_simulateV1 API.
	receiptsRoot := gethcommon.HexToHash("10457e39b8c68ced2071538b4c7034fe68f9c666187fd6b2d6ddcc21149f0d10")
	stateRoot := gethcommon.HexToHash("0275f214e8a0b0ebd8ef427599a0ff339c5171716553b1c522fbd97ac9b108e8")
	proposedBlock = simulated.CreateBlockWithTransactions(
		require.New(s.T()),
		s.SimulationClient,
		proposedBlock,
		blsSigner,
		s.TestNode.ChainSpec,
		s.GenesisValidatorsRoot,
		blobTxs,
		blobTxSidecars,
		&receiptsRoot,
		&stateRoot,
	)

	blockWithCommitments := simulated.CreateBeaconBlockWithBlobs(
		require.New(s.T()),
		s.TestNode.ChainSpec,
		commitments,
		proposedBlock.GetMessage(),
		blsSigner,
		s.GenesisValidatorsRoot,
	)

	sidecarFactory := dablob.NewSidecarFactory(s.TestNode.ChainSpec, metrics.NewNoOpTelemetrySink())
	inclusionProofs := make([][]common.Root, len(blobs))
	for i := range blobs {
		inclusionProof, err := sidecarFactory.BuildKZGInclusionProof(blockWithCommitments.GetMessage().GetBody(), mathpkg.U64(i))
		s.Require().NoError(err)
		inclusionProofs[i] = inclusionProof
	}

	blockWithCommitmentBytes, err := blockWithCommitments.MarshalSSZ()
	s.Require().NoError(err)

	// Inject the new block
	proposal.Txs[0] = blockWithCommitmentBytes

	// Create the beaconBlock Header for the sidecar

	blockWithCommitmentsSignedHeader := ctypes.NewSignedBeaconBlockHeader(
		blockWithCommitments.GetMessage().GetHeader(),
		blockWithCommitments.GetSignature(),
	)

	sidecarsSlice := make([]*datypes.BlobSidecar, len(blobs))
	for i := range blobs {
		sidecar := datypes.BuildBlobSidecar(
			mathpkg.U64(i),
			blockWithCommitmentsSignedHeader,
			blobs[i],
			commitments[i],
			proofs[i],
			inclusionProofs[i],
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
	processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          blockHeight,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
}

// CoreLoop will iterate through the core loop `iterations` times, i.e. Propose, Process, Finalize and Commit.
// Returns the list of proposed comet blocks.
func (s *SimulatedSuite) CoreLoop(startHeight, iterations int64, proposer *signer.BLSSigner) []*types.PrepareProposalResponse {
	// Prepare a block proposal.
	pubkey, err := proposer.GetPubKey()
	s.Require().NoError(err)

	var proposedCometBlocks []*types.PrepareProposalResponse

	for currentHeight := startHeight; currentHeight < startHeight+iterations; currentHeight++ {
		proposalTime := time.Now()
		proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
			Height:          currentHeight,
			Time:            proposalTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposal)

		// Process the proposal.
		processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
			Time:            proposalTime,
		})
		s.Require().NoError(err)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block.
		finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.Ctx, &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(finalizeResp)

		// Commit the block.
		_, err = s.SimComet.Comet.Commit(s.Ctx, &types.CommitRequest{})
		s.Require().NoError(err)

		// Record the Commit Block
		proposedCometBlocks = append(proposedCometBlocks, proposal)
	}
	return proposedCometBlocks
}
