// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package e2e_test

import (
	"github.com/berachain/beacon-kit/testing/e2e/suite"
)

// BeaconE2ESuite is a suite of tests simulating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.KurtosisE2ESuite
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
// TODO: Should check all clients, opposed to the load balancer.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	res, err := s.ConsensusClients()["cl-validator-beaconkit-0"].Stop(s.Ctx())
	s.Require().NoError(err)
	s.Logger().Info("Stopped validator \n", res)

	res, err = s.ConsensusClients()["cl-validator-beaconkit-0"].Start(s.Ctx(), s.Enclave())
	s.Require().NoError(err)
	s.Logger().Info("Started validator \n", res)

	s.Logger().Info("consensus client ports", "ports", s.ConsensusClients()["cl-validator-beaconkit-0"].GetPublicPorts()["cometbft-rpc"].GetNumber())

	err = s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)
}
