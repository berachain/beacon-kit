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
	"math/big"
	"strconv"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/geth-primitives/ssztest"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	mlib "github.com/berachain/beacon-kit/primitives/merkle"
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
	addr, tx, sszTest, err := ssztest.DeploySSZTest(&bind.TransactOpts{
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
		s.Ctx(), strconv.FormatUint(blockNumber, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(blockProposerResp)

	// Verify the beacon block root is equal to HTR(BeaconBlockHeader).
	beaconBlockHeaderRoot := blockProposerResp.BeaconBlockHeader.HashTreeRoot()
	s.Require().Equal(blockProposerResp.BeaconBlockRoot, beaconBlockHeaderRoot)

	// Verify the slot is equal to the requested block number.
	s.Require().Equal(blockProposerResp.BeaconBlockHeader.Slot.Unwrap(), blockNumber)

	// First verify the proposer index proof.
	proposerIndexProof := make([][32]byte, len(blockProposerResp.ProposerIndexProof))
	for i, proofItem := range blockProposerResp.ProposerIndexProof {
		proposerIndexProof[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		proposerIndexProof,
		beaconBlockHeaderRoot,
		blockProposerResp.BeaconBlockHeader.ProposerIndex.HashTreeRoot(),
		big.NewInt(merkle.ProposerIndexGIndexBlock),
	)
	s.Require().NoError(err)

	// Get the header.
	header, err := s.JSONRPCBalancer().HeaderByNumber(s.Ctx(), new(big.Int).SetUint64(blockNumber))
	s.Require().NoError(err)
	s.Require().NotNil(header)

	// Get the chain spec to determine the fork version.
	// TODO: make test use configurable chain spec.
	cs, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	// Get validator pubkey GIndex for the fork version.
	zeroValidatorPubkeyGIndex, err := merkle.GetZeroValidatorPubkeyGIndexBlock(
		cs.ActiveForkVersionForTimestamp(math.U64(header.Time)),
	)
	s.Require().NoError(err)
	gIndex := zeroValidatorPubkeyGIndex +
		(blockProposerResp.BeaconBlockHeader.ProposerIndex.Unwrap() * merkle.ValidatorPubkeyGIndexOffset)

	// Next verify the validator pubkey proof.
	validatorPubkeyProof := make([][32]byte, len(blockProposerResp.ValidatorPubkeyProof))
	for i, proofItem := range blockProposerResp.ValidatorPubkeyProof {
		validatorPubkeyProof[i] = proofItem
	}

	if !mlib.VerifyProof(
		beaconBlockHeaderRoot,
		common.Root(blockProposerResp.ValidatorPubkey.HashTreeRoot()),
		gIndex,
		validatorPubkeyProof,
	) {
		s.FailNow("validator pubkey proof failed to verify against beacon root")
	}

	err = sszTest.MustVerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		validatorPubkeyProof,
		beaconBlockHeaderRoot,
		blockProposerResp.ValidatorPubkey.HashTreeRoot(),
		new(big.Int).SetUint64(gIndex),
	)
	s.Require().NoError(err)
}
