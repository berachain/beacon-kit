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
	"time"

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
	client := s.ConsensusClients()[config.AlternateClient]
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
	s.Require().Equal(fork.Epoch, phase0.Epoch(0))
}

func (s *BeaconKitE2ESuite) TestBeaconValidators() {
	client := s.initBeaconTest()

	// Ensure the validators are not nil.
	validatorsResp, err := client.Validators(
		s.Ctx(),
		&beaconapi.ValidatorsOpts{
			Common: beaconapi.CommonOpts{
				Timeout: 5 * time.Minute,
			},
			State:   utils.StateIDHead,
			Indices: []phase0.ValidatorIndex{0},
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)
	s.Require().NotEmpty(validatorsResp.Data)
}

func (s *BeaconKitE2ESuite) TestBeaconValidatorBalances() {
	client := s.initBeaconTest()

	// Ensure the validator balances are not nil.
	validatorBalancesResp, err := client.ValidatorBalances(
		s.Ctx(),
		&beaconapi.ValidatorBalancesOpts{
			Common: beaconapi.CommonOpts{
				Timeout: 5 * time.Minute,
			},
			State:   utils.StateIDHead,
			Indices: []phase0.ValidatorIndex{0},
			//PubKeys: []phase0.BLSPubKey{},
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(validatorBalancesResp)
}

func (s *BeaconKitE2ESuite) TestBeaconRandao() {
	client := s.initBeaconTest()
	stateRandaoResp, err := client.BeaconStateRandao(s.Ctx(),
		&beaconapi.BeaconStateRandaoOpts{
			State: utils.StateIDHead,
		})
	s.Require().NoError(err)
	s.Require().NotNil(stateRandaoResp)
}
