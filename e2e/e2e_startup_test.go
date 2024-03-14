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
	"encoding/hex"
	"log"
	"os/exec"

	"github.com/berachain/beacon-kit/e2e/suite"
	"github.com/ethereum/go-ethereum/crypto"
)

// BeaconE2ESuite is a suite of tests simulating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.KurtosisE2ESuite
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	err := s.WaitForFinalizedBlockNumber(6)
	s.Require().NoError(err)
}

// TestForgeScriptExecution tests the execution of a forge script
// against the beacon-kit network.
func (s *BeaconKitE2ESuite) TestForgeScriptExecution() {
	url := s.KurtosisE2ESuite.JSONRPCBalancer().URL()
	pk := hex.EncodeToString(crypto.FromECDSA(s.GenesisAccount().PrivateKey()))

	// Change directory to /contracts/ before executing the command
	cmdStr := "cd ../contracts && " +
		"forge build && " +
		"forge script ./script/DeployAndCallERC20.s.sol " +
		"--broadcast --rpc-url=" + url + " " +
		"--private-key=" + pk

	// Execute the command
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf(
			"Failed to execute command: %s, with the error: %s",
			cmdStr,
			err,
		)
	}

	log.Printf("Output: %s\n", string(output))
}
