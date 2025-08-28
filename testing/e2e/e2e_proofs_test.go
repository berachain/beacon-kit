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
	s.Require().NotNil(blockProposerResp.BeaconBlockHeader)

	// Get the next block header.
	nextHeader, err := s.JSONRPCBalancer().HeaderByNumber(
		s.Ctx(), new(big.Int).SetUint64(blockNumber),
	)
	s.Require().NoError(err)
	s.Require().NotNil(nextHeader)

	// Get the block proposer proof for the next timestamp and enforce equality.
	blockProposerResp2, err := s.ConsensusClients()[config.ClientValidator0].BlockProposerProof(
		s.Ctx(), "t"+strconv.FormatUint(nextHeader.Time, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(blockProposerResp2)
	s.Require().NotNil(blockProposerResp2.BeaconBlockHeader)
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
		(blockProposerResp.BeaconBlockHeader.ProposerIndex.Unwrap() * merkle.ValidatorGIndexOffset)

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

// TestValidatorBalanceProof tests the validator balance proof endpoint by fetching and verifying
// validator balance proofs against the SSZTest contract.
func (s *BeaconKitE2ESuite) TestValidatorBalanceProof() {
	// Sender account
	sender := s.TestAccounts()[0]

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Deploy the SSZTest contract to verify the validator balance proof.
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

	// Get the validator balance proof for validator 0 at the parent block number.
	validatorIndex := uint64(0)
	balanceResp, err := s.ConsensusClients()[config.ClientValidator0].ValidatorBalanceProof(
		s.Ctx(), strconv.FormatUint(blockNumber-1, 10), strconv.FormatUint(validatorIndex, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(balanceResp)
	s.Require().NotNil(balanceResp.BeaconBlockHeader)
	s.Require().Equal(balanceResp.BeaconBlockRoot, balanceResp.BeaconBlockHeader.HashTreeRoot())
	s.Require().Equal(balanceResp.BeaconBlockHeader.Slot.Unwrap(), blockNumber-1)

	// Get the next block header.
	nextHeader, err := s.JSONRPCBalancer().HeaderByNumber(s.Ctx(), new(big.Int).SetUint64(blockNumber))
	s.Require().NoError(err)
	s.Require().NotNil(nextHeader)

	// Get the block proposer proof for the next timestamp and enforce equality.
	balanceResp2, err := s.ConsensusClients()[config.ClientValidator0].ValidatorBalanceProof(
		s.Ctx(), "t"+strconv.FormatUint(nextHeader.Time, 10), strconv.FormatUint(validatorIndex, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(balanceResp2)
	s.Require().NotNil(balanceResp2.BeaconBlockHeader)
	s.Require().Equal(*balanceResp.BeaconBlockHeader, *balanceResp2.BeaconBlockHeader)
	s.Require().Equal(balanceResp.BeaconBlockRoot, balanceResp2.BeaconBlockRoot)
	s.Require().Equal(balanceResp.BalanceLeaf, balanceResp2.BalanceLeaf)
	s.Require().ElementsMatch(
		balanceResp.BalanceProof, balanceResp2.BalanceProof,
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
	s.Require().Equal(common.Root(parentBlockRoot4788), balanceResp.BeaconBlockRoot)

	// Get the chain spec to determine the fork version.
	cs, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	// Get the fork version based on the block's timestamp.
	header, err := s.JSONRPCBalancer().HeaderByNumber(
		s.Ctx(), new(big.Int).SetUint64(blockNumber-1),
	)
	s.Require().NoError(err)
	forkVersion := cs.ActiveForkVersionForTimestamp(math.U64(header.Time))
	zeroValidatorBalanceGIndex, err := merkle.GetZeroValidatorBalanceGIndexBlock(forkVersion)
	s.Require().NoError(err)

	// Calculate the balance GIndex based on fork version.
	// Balances are packed 4 per leaf, so we need to divide by 4.
	leafIndex := validatorIndex / 4
	gIndex := zeroValidatorBalanceGIndex + leafIndex

	// Verify the validator balance proof.
	balanceProof := make([][32]byte, len(balanceResp.BalanceProof))
	for i, proofItem := range balanceResp.BalanceProof {
		balanceProof[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		balanceProof,
		balanceResp.BeaconBlockRoot,
		balanceResp.BalanceLeaf, // The leaf contains 4 packed balances
		new(big.Int).SetUint64(gIndex),
	)
	s.Require().NoError(err)
}

// TestValidatorCredentialsProof tests the validator withdrawal credentials proof endpoint by fetching
// and verifying withdrawal credentials proofs against the SSZTest contract.
func (s *BeaconKitE2ESuite) TestValidatorCredentialsProof() {
	// Sender account
	sender := s.TestAccounts()[0]

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Deploy the SSZTest contract to verify the validator credentials proof.
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

	// Get the validator credentials proof for validator 0 at the parent block number.
	validatorIndex := uint64(0)
	credsResp, err := s.ConsensusClients()[config.ClientValidator0].ValidatorCredentialsProof(
		s.Ctx(), strconv.FormatUint(blockNumber-1, 10), strconv.FormatUint(validatorIndex, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(credsResp)
	s.Require().NotNil(credsResp.BeaconBlockHeader)
	s.Require().Equal(credsResp.BeaconBlockRoot, credsResp.BeaconBlockHeader.HashTreeRoot())
	s.Require().Equal(credsResp.BeaconBlockHeader.Slot.Unwrap(), blockNumber-1)

	// Get the next block header.
	nextHeader, err := s.JSONRPCBalancer().HeaderByNumber(s.Ctx(), new(big.Int).SetUint64(blockNumber))
	s.Require().NoError(err)
	s.Require().NotNil(nextHeader)

	// Get the block proposer proof for the next timestamp and enforce equality.
	credsResp1, err := s.ConsensusClients()[config.ClientValidator0].ValidatorCredentialsProof(
		s.Ctx(), "t"+strconv.FormatUint(nextHeader.Time, 10), strconv.FormatUint(validatorIndex, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(credsResp1)
	s.Require().NotNil(credsResp1.BeaconBlockHeader)
	s.Require().Equal(*credsResp.BeaconBlockHeader, *credsResp1.BeaconBlockHeader)
	s.Require().Equal(credsResp.BeaconBlockRoot, credsResp1.BeaconBlockRoot)
	s.Require().Equal(credsResp.ValidatorWithdrawalCredentials, credsResp1.ValidatorWithdrawalCredentials)
	s.Require().ElementsMatch(credsResp.WithdrawalCredentialsProof, credsResp1.WithdrawalCredentialsProof)

	// Get the parent beacon block root of the current timestamp using EIP-4788 Beacon Roots
	// and verify equal to what is returned by the API proof/ endpoint.
	parentBlockRoot4788, err := sszTest.GetParentBlockRootAt(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		nextHeader.Time,
	)
	s.Require().NoError(err)
	s.Require().Equal(common.Root(parentBlockRoot4788), credsResp.BeaconBlockRoot)

	// Get the chain spec to determine the fork version.
	cs, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	// Get the fork version based on the block's timestamp.
	header, err := s.JSONRPCBalancer().HeaderByNumber(
		s.Ctx(), new(big.Int).SetUint64(blockNumber-1),
	)
	s.Require().NoError(err)
	forkVersion := cs.ActiveForkVersionForTimestamp(math.U64(header.Time))
	zeroValidatorCredentialsGIndex, err := merkle.GetZeroValidatorCredentialsGIndexBlock(forkVersion)
	s.Require().NoError(err)

	// Calculate the credentials GIndex based on fork version.
	gIndex := zeroValidatorCredentialsGIndex + (validatorIndex * merkle.ValidatorGIndexOffset)

	// Verify the validator withdrawal credentials proof.
	credentialsProof := make([][32]byte, len(credsResp.WithdrawalCredentialsProof))
	for i, proofItem := range credsResp.WithdrawalCredentialsProof {
		credentialsProof[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		credentialsProof,
		credsResp.BeaconBlockRoot,
		common.Root(credsResp.ValidatorWithdrawalCredentials),
		new(big.Int).SetUint64(gIndex),
	)
	s.Require().NoError(err)

	// Test with a different validator index
	validatorIndex2 := uint64(1)
	credsResp2, err := s.ConsensusClients()[config.ClientValidator0].ValidatorCredentialsProof(
		s.Ctx(), strconv.FormatUint(blockNumber-1, 10), strconv.FormatUint(validatorIndex2, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(credsResp2)

	// Calculate the credentials GIndex for validator 1
	gIndex2 := zeroValidatorCredentialsGIndex + (validatorIndex2 * merkle.ValidatorGIndexOffset)

	// Verify the validator withdrawal credentials proof for validator 1
	credentialsProof2 := make([][32]byte, len(credsResp2.WithdrawalCredentialsProof))
	for i, proofItem := range credsResp2.WithdrawalCredentialsProof {
		credentialsProof2[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{
			Context: s.Ctx(),
		},
		credentialsProof2,
		credsResp2.BeaconBlockRoot,
		common.Root(credsResp2.ValidatorWithdrawalCredentials),
		new(big.Int).SetUint64(gIndex2),
	)
	s.Require().NoError(err)
}

// TestValidatorPubkeyProof tests the validator pubkey proof endpoint by fetching and verifying
// validator pubkey proofs against the SSZTest contract.
func (s *BeaconKitE2ESuite) TestValidatorPubkeyProof() {
	// Sender account
	sender := s.TestAccounts()[0]

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)

	// Deploy the SSZTest contract to verify the validator pubkey proof.
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

	// Get the validator pubkey proof for validator 0 at the parent block number.
	validatorIndex := uint64(0)
	pubkeyResp, err := s.ConsensusClients()[config.ClientValidator0].ValidatorPubkeyProof(
		s.Ctx(), strconv.FormatUint(blockNumber-1, 10), strconv.FormatUint(validatorIndex, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(pubkeyResp)
	s.Require().NotNil(pubkeyResp.BeaconBlockHeader)
	s.Require().Equal(pubkeyResp.BeaconBlockRoot, pubkeyResp.BeaconBlockHeader.HashTreeRoot())
	s.Require().Equal(pubkeyResp.BeaconBlockHeader.Slot.Unwrap(), blockNumber-1)

	// Get the next block header.
	nextHeader, err := s.JSONRPCBalancer().HeaderByNumber(s.Ctx(), new(big.Int).SetUint64(blockNumber))
	s.Require().NoError(err)
	s.Require().NotNil(nextHeader)

	// Get the pubkey proof for the next timestamp and enforce equality.
	pubkeyResp2, err := s.ConsensusClients()[config.ClientValidator0].ValidatorPubkeyProof(
		s.Ctx(), "t"+strconv.FormatUint(nextHeader.Time, 10), strconv.FormatUint(validatorIndex, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(pubkeyResp2)
	s.Require().NotNil(pubkeyResp2.BeaconBlockHeader)
	s.Require().Equal(*pubkeyResp.BeaconBlockHeader, *pubkeyResp2.BeaconBlockHeader)
	s.Require().Equal(pubkeyResp.BeaconBlockRoot, pubkeyResp2.BeaconBlockRoot)
	s.Require().Equal(pubkeyResp.ValidatorPubkey, pubkeyResp2.ValidatorPubkey)
	s.Require().ElementsMatch(pubkeyResp.ValidatorPubkeyProof, pubkeyResp2.ValidatorPubkeyProof)

	// Get the parent beacon block root of the current timestamp using EIP-4788 Beacon Roots
	// and verify equal to what is returned by the API proof/ endpoint.
	parentBlockRoot4788, err := sszTest.GetParentBlockRootAt(
		&bind.CallOpts{Context: s.Ctx()},
		nextHeader.Time,
	)
	s.Require().NoError(err)
	s.Require().Equal(common.Root(parentBlockRoot4788), pubkeyResp.BeaconBlockRoot)

	// Get the chain spec to determine the fork version.
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

	// Calculate the pubkey GIndex based on fork version.
	gIndex := zeroValidatorPubkeyGIndex + (validatorIndex * merkle.ValidatorGIndexOffset)

	// Verify the validator pubkey proof.
	validatorPubkeyProof := make([][32]byte, len(pubkeyResp.ValidatorPubkeyProof))
	for i, proofItem := range pubkeyResp.ValidatorPubkeyProof {
		validatorPubkeyProof[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{Context: s.Ctx()},
		validatorPubkeyProof,
		pubkeyResp.BeaconBlockRoot,
		pubkeyResp.ValidatorPubkey.HashTreeRoot(),
		new(big.Int).SetUint64(gIndex),
	)
	s.Require().NoError(err)

	// Test with a different validator index
	validatorIndex2 := uint64(1)
	pubkeyResp3, err := s.ConsensusClients()[config.ClientValidator0].ValidatorPubkeyProof(
		s.Ctx(), strconv.FormatUint(blockNumber-1, 10), strconv.FormatUint(validatorIndex2, 10),
	)
	s.Require().NoError(err)
	s.Require().NotNil(pubkeyResp3)

	// Calculate the pubkey GIndex for validator 1
	gIndex2 := zeroValidatorPubkeyGIndex + (validatorIndex2 * merkle.ValidatorGIndexOffset)

	// Verify the validator pubkey proof for validator 1
	validatorPubkeyProof2 := make([][32]byte, len(pubkeyResp3.ValidatorPubkeyProof))
	for i, proofItem := range pubkeyResp3.ValidatorPubkeyProof {
		validatorPubkeyProof2[i] = proofItem
	}
	err = sszTest.MustVerifyProof(
		&bind.CallOpts{Context: s.Ctx()},
		validatorPubkeyProof2,
		pubkeyResp3.BeaconBlockRoot,
		pubkeyResp3.ValidatorPubkey.HashTreeRoot(),
		new(big.Int).SetUint64(gIndex2),
	)
	s.Require().NoError(err)
}
