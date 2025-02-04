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
	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
)

// initBeaconTest initializes the any tests for the beacon node api.
func (s *BeaconKitE2ESuite) initBeaconTest() *types.ConsensusClient {
	// Wait for execution block 5.
	err := s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)

	// Get the consensus client.
	client := s.ConsensusClients()[config.ClientValidator0]
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
