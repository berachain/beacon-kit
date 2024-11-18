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
	"math/big"

	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// BeaconE2ESuite is a suite of tests simulating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.KurtosisE2ESuite
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
// TODO: Should check all clients, opposed to just the load balancer.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	err := s.WaitForFinalizedBlockNumber(10)
	s.Require().NoError(err)
}

// TestEVMInflation checks that the EVM inflation address receives the correct
// amount of EVM inflation per block.
func (s *BeaconKitE2ESuite) TestEVMInflation() {
	evmInflationPerBlockWei, _ := big.NewFloat(
		config.EVMInflationPerBlockWei).Int(nil)

	// Check over the next 10 EVM blocks, that after every block, the balance
	// of the EVM inflation address increases by EVMInflationPerBlockWei.
	for i := int64(1); i <= 10; i++ {
		err := s.WaitForFinalizedBlockNumber(uint64(i))
		s.Require().NoError(err)

		balance, err := s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.HexToAddress(config.EVMInflationAddress),
			big.NewInt(i),
		)
		s.Require().NoError(err)
		s.Require().Equal(
			balance.Cmp(new(big.Int).Mul(
				evmInflationPerBlockWei, big.NewInt(i)),
			), 0,
		)
	}
}
