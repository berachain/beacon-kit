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
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/cometbft/cometbft/abci/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// finalizeWaitDuration TODO: If we're too quick to finalize after process, we see the EL context get cancelled, so we wait. Figure out why.
const finalizeWaitDuration = 500 * time.Millisecond

type Simulated struct {
	suite.Suite
	simulated.SharedAccessors
	SimComet *simulated.SimComet
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

	go func() {
		// Node blocks on Start and hence we have to run in separate routine
		if err = s.TestNode.Start(s.Ctx); err != nil {
			s.T().Error(err)
		}
	}()
	// Wait for ~2 seconds for services to start
	<-time.After(2 * time.Second)
}

func (s *Simulated) TearDownTest() {
	err := s.ElHandle.Close()
	if err != nil {
		s.T().Error("Error closing EL handle")
	}
	s.CancelFunc()
}

func (s *Simulated) TestFullLifecycle_ValidBlock_IsSuccessful() {
	const height = 1
	appGenesis, err := genutiltypes.AppGenesisFromFile(
		s.HomeDir + "/config/genesis.json",
	)
	s.Require().NoError(err)
	initResponse, err := s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId:       "test-mainnet-chain",
		AppStateBytes: appGenesis.AppState,
	})
	s.Require().NoError(err)
	s.Require().Len(initResponse.GetValidators(), 1, "Expected 1 validator")
	deposits, err := s.TestNode.StorageBackend.DepositStore().GetDepositsByIndex(
		s.Ctx, constants.FirstDepositIndex,
		constants.FirstDepositIndex+s.TestNode.ChainSpec.MaxDepositsPerBlock(),
	)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1, "Expected 1 deposit")

	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          height,
		Time:            time.Now(),
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	processResponse, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:    proposal.Txs,
		Height: height,
		// If incorrect proposer address is used, we get a proposer mismatch error
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResponse.Status)

	<-time.After(finalizeWaitDuration)

	finalizeResponse, err := s.SimComet.Comet.FinalizeBlock(s.Ctx, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          height,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(finalizeResponse)

	_, err = s.SimComet.Comet.Commit(s.Ctx, &types.CommitRequest{})
	s.Require().NoError(err)

	// Post state checks
	queryContext, err := s.SimComet.CreateQueryContext(height, false)
	s.Require().NoError(err)
	stateDB := s.TestNode.StorageBackend.StateFromContext(queryContext)
	slot, err := stateDB.GetSlot()
	s.Require().NoError(err)
	s.Require().Equal(math.U64(height), slot)

	fetchedHeader, err := stateDB.GetLatestBlockHeader()
	s.Require().NoError(err)

	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForSlot(height),
	)
	s.Require().NoError(err)
	s.Require().Equal(proposedBlock.Message.GetHeader().GetBodyRoot(), fetchedHeader.GetBodyRoot())
}

func (s *Simulated) TestInitChain_InvalidChainID_MustError() {
	_, err := s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId: "henlo-im-invalid",
	})
	s.Require().ErrorContains(err, "invalid chain-id on InitChain; expected: test-mainnet-chain, got: henlo-im-invalid")
}
