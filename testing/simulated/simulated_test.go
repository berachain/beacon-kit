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
	"testing"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type Simulated struct {
	suite.Suite
	simulated.TestSuiteHandle
	SimComet          *simulated.SimComet
	BlockchainService *blockchain.Service
	// LogBuffer gives us a mechanism to access the reason for a comet rejection
	LogBuffer             *bytes.Buffer
	GenesisValidatorsRoot common.Root
}

// TestCustomCometComponent is a test suite with a custom comet driver can can control ourselves
//
//nolint:paralleltest // cannot be run in parallel due to use of environment variables.
func TestSimulatedCometComponent(t *testing.T) {
	suite.Run(t, new(Simulated))
}

func (s *Simulated) SetupTest() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.Ctx = ctx
	s.CancelFunc = cancelFunc

	tempHomeDir := s.T().TempDir()
	s.HomeDir = tempHomeDir
	// Initialize home directory
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), tempHomeDir)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the Geth node - needs to be done first as we need the auth rpc as input for the beacon node.
	elNode := simulated.NewGethNode(tempHomeDir, simulated.ValidGethImage())
	gethHandle, authRPC := elNode.Start(s.T())
	s.ElHandle = gethHandle

	// Build the Beacon node once we have the auth rpc url
	var logBuffer bytes.Buffer
	logger := phuslu.NewLogger(&logBuffer, nil)
	s.LogBuffer = &logBuffer

	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)

	testNode := simulated.NewTestNode(s.T(),
		simulated.TestNodeInput{
			TempHomeDir: tempHomeDir,
			CometConfig: cometConfig,
			AuthRPC:     authRPC,
			Logger:      logger,
			AppOpts:     viper.New(),
			Components:  components,
		})
	s.TestNode = testNode

	// Fetch services we will want to query and interact with so they are easily accessible in testing
	var noopCometService *simulated.SimComet
	err := testNode.FetchService(&noopCometService)
	s.Require().NoError(err)
	s.NotNil(noopCometService)
	s.SimComet = noopCometService

	var blockchainService *blockchain.Service
	err = testNode.FetchService(&blockchainService)
	s.Require().NoError(err)
	s.NotNil(blockchainService)
	s.BlockchainService = blockchainService
}

func (s *Simulated) TearDownTest() {
	err := s.ElHandle.Close()
	if err != nil {
		s.T().Error("Error closing geth handle")
	}
	s.CancelFunc()
}

func (s *Simulated) TestInitChain_InvalidChainID_MustError() {
	_, err := s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId: "henlo-im-invalid",
	})
	s.Require().ErrorContains(err, "invalid chain-id on InitChain; expected: test-mainnet-chain, got: henlo-im-invalid")
}

func (s *Simulated) TestInitChain_Valid_IsSuccessful() {
	appGenesis, err := genutiltypes.AppGenesisFromFile(
		s.HomeDir + "/config/genesis.json",
	)
	s.Require().NoError(err)
	res, err := s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId:       "test-mainnet-chain",
		AppStateBytes: appGenesis.AppState,
	})
	s.Require().NoError(err)
	// We expect 1 validator after initchain as there was only 1 deposit
	s.Require().Len(res.GetValidators(), 1)

	// The deposit store must have 1 deposit
	// TODO: Get depRange from chainspec
	deposits, err := s.TestNode.StorageBackend.DepositStore().GetDepositsByIndex(s.Ctx, constants.FirstDepositIndex, 100)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1)
}

func (s *Simulated) TestPrepareProposal_ValidRequest_IsSuccessful() {
	_, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height: 1,
	})
	s.Require().NoError(err)
}

func (s *Simulated) TestProcessProposal_HeightZero_MustError() {
	_, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Height: 0,
	})
	s.Require().ErrorContains(err, "processProposal at height 0: invalid height")
}

func (s *Simulated) TestProcessProposal_NilBeaconBlock_MustError() {
	res, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:    make([][]byte, 2),
		Height: 1,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, res.Status)
	s.Require().Contains(s.LogBuffer.String(), "nil beacon block in abci request")
}

func (s *Simulated) TestProcessProposal_ValidProposal_MustAccept() {
	// Initialize the chain correctly
	appGenesis, err := genutiltypes.AppGenesisFromFile(
		s.HomeDir + "/config/genesis.json",
	)
	s.Require().NoError(err)
	_, err = s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId:       "test-mainnet-chain",
		AppStateBytes: appGenesis.AppState,
	})
	s.Require().NoError(err)

	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	beaconChain := simulated.GenerateBeaconChain(s.T(), 2, func(block *gethprimitives.Block) (*ctypes.SignedBeaconBlock, error) {
		// Generate the valid proposal
		var beaconBlock *ctypes.BeaconBlock
		if beaconBlock, err = ctypes.NewBeaconBlockWithVersion(
			math.Slot(block.NumberU64()),
			math.ValidatorIndex(0),
			common.Root{1, 2, 3, 4, 5},
			version.Deneb1(),
		); err != nil {
			return nil, err
		}
		beaconBlock.StateRoot = common.Root{5, 4, 3, 2, 1}
		beaconBlock.Body = &ctypes.BeaconBlockBody{
			ExecutionPayload: simulated.BlockToExecutionPayload(block),
		}
		// Use propose index 0
		beaconBlock.ProposerIndex = 0

		body := beaconBlock.GetBody()
		body.SetProposerSlashings(ctypes.ProposerSlashings{})
		body.SetAttesterSlashings(ctypes.AttesterSlashings{})
		body.SetAttestations(ctypes.Attestations{})
		body.SetSyncAggregate(&ctypes.SyncAggregate{})
		body.SetVoluntaryExits(ctypes.VoluntaryExits{})
		body.SetBlsToExecutionChanges(ctypes.BlsToExecutionChanges{})

		var signedBeaconBlock *ctypes.SignedBeaconBlock
		if signedBeaconBlock, err = ctypes.NewSignedBeaconBlock(
			beaconBlock,
			ctypes.NewForkData(version.Deneb(), s.GenesisValidatorsRoot),
			s.TestNode.ChainSpec,
			blsSigner,
		); err != nil {
			return nil, err
		}
		return signedBeaconBlock, nil
	})

	blockBytes, err := beaconChain[0].MarshalSSZ()
	s.Require().NoError(err)

	blob := datypes.BlobSidecars{}
	blobBytes, err := blob.MarshalSSZ()
	s.Require().NoError(err)

	txs := make([][]byte, 2)
	txs[0] = blockBytes
	txs[1] = blobBytes

	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	res, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:    txs,
		Height: 1,
		// If incorrect proposer address is used, we get a proposer mismatch error
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, res.Status)
}
