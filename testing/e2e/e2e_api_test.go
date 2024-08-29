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
	"fmt"
	beaconapi "github.com/attestantio/go-eth2-client/api"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/testing/e2e/config"
)

// TestBeaconAPISuite tests that the api test suite is setup correctly with a
// working beacon node-api client.
func (s *BeaconKitE2ESuite) TestBeaconAPIStartup() {
	// Wait for execution block 5.
	err := s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)

	// Get the consensus client.
	client := s.ConsensusClients()[config.DefaultClient]
	s.Require().NotNil(client)

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
	fmt.Println("stateRootResp", stateRootResp.Data.String())

	// Ensure the state fork is not nil.
	//Error: Received unexpected error: failed to unmarshal data : previous version missing

	stateForkResp, err := client.Fork(s.Ctx(), &beaconapi.ForkOpts{
		State: stateRootResp.Data.String(),
	})
	fmt.Println("err", err)
	s.Require().NoError(err)
	s.Require().NotNil(stateForkResp)

	// Ensure the state validators are not nil.
	//stateValidatorsResp, err := client.Validators(s.Ctx(), &beaconapi.ValidatorsOpts{
	//	State:   utils.StateIDHead,
	//	Indices: make([]phase0.ValidatorIndex, 0),
	//PubKeys:         make([]phase0.BLSPubKey, 0),
	//ValidatorStates: make([]apiv1.ValidatorState, 0),
	//})
	//s.Require().NoError(err)
	//s.Require().NotNil(stateValidatorsResp)
	////s.Require().NotEmpty(stateValidatorsResp.Data)

	// Ensure the state validator are not nil.
	//stateValidatorResp, err := client.Validator(s.Ctx(), &beaconapi.ValidatorsOpts{
	//	State: utils.StateIDHead,
	//})
	//s.Require().NoError(err)
	//s.Require().NotNil(stateValidatorResp)
	//s.Require().NotNil(stateValidatorResp.Data)

	// {"level":"debug","service":"client",
	// "impl":"http","id":"36162dc1","address":"http://0.0.0.0:52501",
	//"endpoint":"/eth/v1/beacon/states/0/validator_balances",
	//"status_code":400,"status_code":400,
	//"response":{"code":400,"message":"invalid request"},"time":"2024-08-29T15:58:40+05:30","message":"POST failed"}
	// Ensure the state validator balances are not nil.
	//stateValidatorBalanceResp, err := client.ValidatorBalances(s.Ctx(), &beaconapi.ValidatorBalancesOpts{
	//	State: "0",
	//	//Indices: []phase0.ValidatorIndex{},
	//	//PubKeys: nil,
	//})
	//s.Require().NoError(err)
	//s.Require().NotNil(stateValidatorBalanceResp)

	// json: cannot unmarshal string into Go value of type http.beaconStateRandaoJSON
	// Ensure beacon randao is not nil.
	//stateRandaoResp, err := client.BeaconStateRandao(s.Ctx(), &beaconapi.BeaconStateRandaoOpts{
	//	State: utils.StateIDHead,
	//})
	//s.Require().NoError(err)
	//s.Require().NotNil(stateRandaoResp)

}
