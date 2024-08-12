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
	stateRoot, err := client.BeaconStateRoot(
		s.Ctx(),
		&beaconapi.BeaconStateRootOpts{
			State: utils.StateIDHead,
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(stateRoot)
	s.Require().False(stateRoot.Data.IsZero())
}
