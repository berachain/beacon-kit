// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"encoding/hex"
	"strconv"

	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
)

// initBeaconTest initializes the any tests for the beacon node api.
func (s *BeaconKitE2ESuite) initBeaconTest() *types.ConsensusClient {
	// Wait for execution block 5.
	err := s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)

	// Get the consensus client.
	client := s.ConsensusClients()[config.DefaultClient]
	s.Require().NotNil(client)

	return client
}

// TestBeaconStateRoot tests the beacon node api for beacon state root.
func (s *BeaconKitE2ESuite) TestBeaconStateRoot() {
	client := s.initBeaconTest()

	// Ensure the state root is not nil.
	stateRootResp, err := client.BeaconStateRoot(
		s.Ctx(),
		&beaconapi.BeaconStateRootOpts{
			State: utils.StateIDHead,
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(stateRootResp)
	s.Require().False(stateRootResp.Data.IsZero())
}

// TestBeaconFork tests the beacon node api for beacon fork.
func (s *BeaconKitE2ESuite) TestBeaconFork() {
	client := s.initBeaconTest()

	stateForkResp, err := client.Fork(s.Ctx(), &beaconapi.ForkOpts{
		State: utils.StateIDHead,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(stateForkResp)

	fork := stateForkResp.Data
	s.Require().NotEmpty(fork.PreviousVersion)
	s.Require().NotEmpty(fork.CurrentVersion)
	expectedVersion := phase0.Version{0x04, 0x00, 0x00, 0x00}
	s.Require().Equal(
		expectedVersion,
		fork.PreviousVersion,
		"PreviousVersion does not match expected value",
	)
	s.Require().Equal(
		expectedVersion,
		fork.CurrentVersion,
		"CurrentVersion does not match expected value",
	)
	s.Require().Equal(
		phase0.Epoch(0),
		fork.Epoch,
		"Epoch does not match expected value",
	)
}

// TestBeaconValidators tests the beacon node api for beacon validators.
//
//nolint:lll
func (s *BeaconKitE2ESuite) TestBeaconValidators() {
	client := s.initBeaconTest()

	indices := []phase0.ValidatorIndex{0}
	// Ensure the validators are not nil.
	validatorsResp, err := client.Validators(
		s.Ctx(),
		&beaconapi.ValidatorsOpts{
			State:   utils.StateIDHead,
			Indices: indices,
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	validatorData := validatorsResp.Data
	s.Require().NotNil(validatorData, "Validator data should not be nil")
	s.Require().
		Len(validatorData, len(indices), "Number of validator responses should match number of requested indices")

	for _, validator := range validatorData {
		s.Require().NotNil(validator, "Validator should not be nil")

		s.Require().
			Contains(indices, validator.Index, "Validator index should be one of the requested indices")

		s.Require().
			NotEmpty(validator.Validator.PublicKey, "Validator public key should not be empty")
		s.Require().
			Len(validator.Validator.PublicKey, 48, "Validator public key should be 48 bytes long")

		s.Require().
			NotEmpty(validator.Validator.WithdrawalCredentials, "Withdrawal credentials should not be empty")
		s.Require().
			Len(validator.Validator.WithdrawalCredentials, 32, "Withdrawal credentials should be 32 bytes long")

		s.Require().
			True(validator.Validator.EffectiveBalance > 0, "Effective balance should be positive")
		s.Require().
			True(validator.Validator.EffectiveBalance <= 32e9, "Effective balance should not exceed 32 ETH")

		s.Require().
			False(validator.Validator.Slashed, "Slashed status should not be true")

		s.Require().
			True(validator.Validator.ActivationEpoch >= validator.Validator.ActivationEligibilityEpoch,
				"Activation epoch should be greater than or equal to activation eligibility epoch")

		s.Require().
			True(validator.Validator.WithdrawableEpoch >= validator.Validator.ExitEpoch,
				"Withdrawable epoch should be greater than or equal to exit epoch")

		s.Require().
			NotEmpty(validator.Status, "Validator status should not be empty")

		s.Require().
			True(validator.Balance > 0, "Validator balance should be positive")
		s.Require().
			True(validator.Balance <= 32e9, "Validator balance should not exceed 32 ETH")
	}
}

func (s *BeaconKitE2ESuite) TestBeaconValidatorBalances() {
	client := s.initBeaconTest()

	indices := []phase0.ValidatorIndex{0}
	// Ensure the validator balances are not nil.
	validatorBalancesResp, err := client.ValidatorBalances(
		s.Ctx(),
		&beaconapi.ValidatorBalancesOpts{
			State:   utils.StateIDHead,
			Indices: indices,
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(validatorBalancesResp)
	balanceMap := validatorBalancesResp.Data
	for _, index := range indices {
		balance, exists := balanceMap[index]
		s.Require().
			True(exists, "Balance should exist for validator index %d", index)
		s.Require().True(balance > 0, "Validator balance should be positive")
	}
}

func (s *BeaconKitE2ESuite) TestBeaconRandao() {
	client := s.initBeaconTest()
	stateRandaoResp, err := client.BeaconStateRandao(s.Ctx(),
		&beaconapi.BeaconStateRandaoOpts{
			State: utils.StateIDHead,
		})
	s.Require().NoError(err)
	s.Require().NotNil(stateRandaoResp)
	s.Require().NotEmpty(stateRandaoResp.Data)
	randao := stateRandaoResp.Data
	s.Require().Len(
		randao,
		32,
		"Randao should be 32 bytes long",
	)
	s.Require().NotEqual(
		make([]byte, 32),
		randao,
		"Randao should not be all zeros",
	)
}

func (s *BeaconKitE2ESuite) TestBeaconGenesis() {
	client := s.initBeaconTest()

	genesisResp, err := client.Genesis(s.Ctx(),
		&beaconapi.GenesisOpts{})
	s.Require().NoError(err)
	s.Require().NotNil(genesisResp)

	s.Require().NotZero(
		genesisResp.Data.GenesisTime,
		"Genesis time should not be zero",
	)

	s.Require().NotEmpty(
		genesisResp.Data.GenesisValidatorsRoot,
		"Genesis validators root should not be empty",
	)

	// s.Require().NotEmpty(
	//	genesisResp.Data.GenesisForkVersion,
	//	"Genesis fork version should be empty",
	// )
}

func (s *BeaconKitE2ESuite) TestBeaconBlockHeaderByID() {
	client := s.initBeaconTest()

	// Test getting the genesis block header.
	genesisResp, err := client.BeaconBlockHeader(
		s.Ctx(),
		&beaconapi.BeaconBlockHeaderOpts{
			Block: "genesis",
		},
	)

	s.Require().NoError(err)
	s.Require().NotNil(genesisResp)
	s.Require().NotNil(genesisResp.Data)
	s.Require().NotZero(genesisResp.Data.Root)
	s.Require().NotZero(genesisResp.Data.Header.Message.Slot)
	// Check slot to be equal than 1
	s.Require().Equal(uint64(1), uint64(genesisResp.Data.Header.Message.Slot))

	// Test getting the head block header.
	headResp, err := client.BeaconBlockHeader(
		s.Ctx(),
		&beaconapi.BeaconBlockHeaderOpts{
			Block: "head",
		},
	)

	s.Require().NoError(err)
	s.Require().NotNil(headResp)
	s.Require().NotNil(headResp.Data)
	s.Require().NotZero(headResp.Data.Root)
	s.Require().NotZero(headResp.Data.Header.Message.Slot)
	// Check slot to be greater than 1
	s.Require().Greater(uint64(headResp.Data.Header.Message.Slot), uint64(1))

	// Test getting a block header by slot.
	slot := headResp.Data.Header.Message.Slot - 1
	slotResp, err := client.BeaconBlockHeader(
		s.Ctx(),
		&beaconapi.BeaconBlockHeaderOpts{
			Block: strconv.FormatUint(uint64(slot), 10),
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(slotResp)
	s.Require().NotNil(slotResp.Data)
	s.Require().NotZero(slotResp.Data.Root)
	s.Require().Equal(uint64(slot), uint64(slotResp.Data.Header.Message.Slot))

	// Test getting a block header by block root.
	rootResp, err := client.BeaconBlockHeader(
		s.Ctx(),
		&beaconapi.BeaconBlockHeaderOpts{
			Block: "0x" + hex.EncodeToString(headResp.Data.Root[:]),
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(rootResp)
	s.Require().NotNil(rootResp.Data)
	s.Require().Equal(headResp.Data.Root, rootResp.Data.Root)
}
