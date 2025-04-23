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
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// TestBlockProposerProof tests the block proposer proof endpoint by fetching and verifying
// the block proposer proofs against the SSZTest contract. Refer to
// beacon-kit/contracts/src/eip4788/SSZ.sol for details.
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

	// Get the block proposer proof for the parent block number.
	blockProposerResp, err := s.ConsensusClients()[config.ClientValidator0].BlockProposerProof(
		s.Ctx(), strconv.FormatUint(blockNumber-1, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(blockProposerResp)

	// Get the current block header.
	nextHeader, err := s.JSONRPCBalancer().HeaderByNumber(
		s.Ctx(), new(big.Int).SetUint64(blockNumber),
	)
	s.Require().NoError(err)
	s.Require().NotNil(nextHeader)

	// Get the block proposer proof for the current timestamp and enforce equality.
	blockProposerResp2, err := s.ConsensusClients()[config.ClientValidator0].BlockProposerProof(
		s.Ctx(), "t"+strconv.FormatUint(nextHeader.Time, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(blockProposerResp2)
	s.Require().Equal(*blockProposerResp.BeaconBlockHeader, *blockProposerResp2.BeaconBlockHeader)
	s.Require().Equal(blockProposerResp.BeaconBlockRoot, blockProposerResp2.BeaconBlockRoot)
	s.Require().Equal(blockProposerResp.ValidatorPubkey, blockProposerResp2.ValidatorPubkey)
	s.Require().ElementsMatch(
		blockProposerResp.ValidatorPubkeyProof, blockProposerResp2.ValidatorPubkeyProof,
	)
	s.Require().ElementsMatch(
		blockProposerResp.ProposerIndexProof, blockProposerResp2.ProposerIndexProof,
	)

	// Get the parent beacon block root of the current timestamp using EIP-4788 Beacon Roots
	// and verify equal to what is returned by the API proof/ endpoint.
	parentBlockRoot4788, err := sszTest.GetParentBlockRootAt(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		nextHeader.Time,
	)
	s.Require().NoError(err)
	s.Require().Equal(common.Root(parentBlockRoot4788), blockProposerResp.BeaconBlockRoot)

	// Verify the beacon block root is equal to HTR(BeaconBlockHeader).
	s.Require().Equal(
		blockProposerResp.BeaconBlockRoot, blockProposerResp.BeaconBlockHeader.HashTreeRoot(),
	)

	// Verify the slot is equal to the requested block number.
	s.Require().Equal(blockProposerResp.BeaconBlockHeader.Slot.Unwrap(), blockNumber-1)

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
		blockProposerResp.BeaconBlockRoot,
		blockProposerResp.BeaconBlockHeader.ProposerIndex.HashTreeRoot(),
		big.NewInt(merkle.ProposerIndexGIndexBlock),
	)
	s.Require().NoError(err)

	// If the proof or leaf is modified, the proof should not verify.
	isVerified, err := sszTest.VerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		proposerIndexProof,
		blockProposerResp.BeaconBlockRoot,
		(blockProposerResp.BeaconBlockHeader.ProposerIndex + 2).HashTreeRoot(), // malicious leaf
		big.NewInt(merkle.ProposerIndexGIndexBlock),
	)
	s.Require().NoError(err)
	s.Require().False(isVerified)

	// Get the chain spec to determine the fork version.
	// TODO: make test use configurable chain spec.
	cs, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	// Get the fork version based on the block's timestamp.
	header, err := s.JSONRPCBalancer().HeaderByNumber(
		s.Ctx(), new(big.Int).SetUint64(blockNumber-1),
	)
	s.Require().NoError(err)
	forkVersion := cs.ActiveForkVersionForTimestamp(math.U64(header.Time))
	zeroValidatorPubkeyGIndex, err := merkle.GetZeroValidatorPubkeyGIndexBlock(forkVersion)
	s.Require().NoError(err)

	// Calculate the validator pubkey GIndex based on fork version.
	gIndex := zeroValidatorPubkeyGIndex +
		(blockProposerResp.BeaconBlockHeader.ProposerIndex.Unwrap() * merkle.ValidatorPubkeyGIndexOffset)

	// Next verify the validator pubkey proof.
	validatorPubkeyProof := make([][32]byte, len(blockProposerResp.ValidatorPubkeyProof))
	for i, proofItem := range blockProposerResp.ValidatorPubkeyProof {
		validatorPubkeyProof[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		validatorPubkeyProof,
		blockProposerResp.BeaconBlockRoot,
		blockProposerResp.ValidatorPubkey.HashTreeRoot(),
		new(big.Int).SetUint64(gIndex),
	)
	s.Require().NoError(err)
}
