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

package e2e_test

import (
	"fmt"

	"github.com/berachain/beacon-kit/geth-primitives/ssztest"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// TestBlockProposerProof tests the block proposer proof endpoint.
func (s *BeaconKitE2ESuite) TestBlockProposerProof() {
	// Sender account
	sender := s.TestAccounts()[0]

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Deploy the SSZTest contract to verify the block proposer proof.
	addr, tx, _, err := ssztest.DeploySSZTest(&bind.TransactOpts{
		From:     sender.Address(),
		Signer:   sender.SignerFunc(chainID),
		GasLimit: 1000000,
		Context:  s.Ctx(),
	}, s.JSONRPCBalancer())
	s.Require().NoError(err)

	// Confirm deployment.
	receipt, err := bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)
	s.Logger().Info("SSZTest contract deployed successfully", "address", addr.Hex())

	// Get the current block number.
	blockNumber, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)

	// Get the block proposer proof for the current block number.
	blockProposerResp, err := s.ConsensusClients()[config.ClientValidator0].BlockProposerProof(
		s.Ctx(), fmt.Sprintf("%d", blockNumber),
	)
	s.Require().NoError(err)
	s.Require().NotNil(blockProposerResp)

	// Verify the beacon block root is equal to HTR(BeaconBlockHeader).
	beaconBlockHeaderRoot := blockProposerResp.BeaconBlockHeader.HashTreeRoot()
	s.Require().Equal(blockProposerResp.BeaconBlockRoot, beaconBlockHeaderRoot)

	// Verify the slot is equal to the requested block number.
	s.Require().Equal(blockProposerResp.BeaconBlockHeader.Slot.Unwrap(), blockNumber)
}
