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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//go:build e2e

package preconf_test

import (
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
)

const (
	// Log messages that indicate preconf is working.
	sequencerServingLog  = "GetPayloadBySlot completed"
	validatorFetchingLog = "Successfully fetched payload from sequencer"

	// Kurtosis service name for the dedicated sequencer CL node.
	sequencerCLService = "cl-sequencer-beaconkit-0"

	// Number of blocks to wait.
	blocksToWait        = 20
	blocksAfterFallback = 10
)

// TestSequencerFlow verifies the preconf pathway using ordered subtests.
// NOTE: This test removes the sequencer service in SequencerFallback,
// so it must run after TestPreconfTransactions (testify runs alphabetically).
func (s *PreconfE2ESuite) TestSequencerFlow() {
	// Wait for enough blocks so preconf has been exercised.
	err := s.WaitForFinalizedBlockNumber(blocksToWait)
	s.Require().NoError(err, "Network should reach finalized blocks")

	// Verify sequencer is serving payloads.
	s.Run("SequencerServesPayloads", func() {
		logs, err := s.GetServiceLogs(sequencerCLService)
		s.Require().NoError(err, "Should get sequencer logs")

		found := suite.ContainsLogMessage(logs, sequencerServingLog)
		s.Require().True(found,
			"Sequencer (%s) should serve payloads to validators. "+
				"Expected log message containing: %q",
			sequencerCLService, sequencerServingLog)
	})

	// Verify validators fetch from sequencer.
	s.Run("ValidatorFetches", func() {
		for _, validator := range []string{config.ClientValidator0} {
			validator := validator // capture for closure
			s.Run(validator, func() {
				logs, err := s.GetServiceLogs(validator)
				s.Require().NoError(err, "Should get validator logs for %s", validator)

				found := suite.ContainsLogMessage(logs, validatorFetchingLog)
				s.Require().True(found,
					"Validator (%s) should fetch payloads from sequencer. "+
						"Expected log message containing: %q",
					validator, validatorFetchingLog)
			})
		}
	})

	// Remove sequencer and verify validators fall back to local building.
	s.Run("SequencerFallback", func() {
		// Get current block before removing sequencer.
		currentBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get current block number")

		// Remove sequencer -- simulates crash/unavailability.
		s.T().Logf("Removing sequencer (%s)...", sequencerCLService)
		err = s.RemoveService(sequencerCLService)
		s.Require().NoError(err, "Should remove sequencer service")

		// Wait for more blocks -- validators must build locally now.
		targetBlock := currentBlock + blocksAfterFallback
		s.T().Logf("Waiting for %d more blocks (current: %d, target: %d)...",
			blocksAfterFallback, currentBlock, targetBlock)

		err = s.WaitForFinalizedBlockNumber(targetBlock)
		s.Require().NoError(err, "Network should continue producing blocks after sequencer removed")

		// Verify network continued.
		finalBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get final block number")
		s.Require().GreaterOrEqual(finalBlock, targetBlock,
			"Network should have produced blocks after sequencer removed")

		s.T().Logf("Network continued: block %d -> %d (fallback working)", currentBlock, finalBlock)
	})
}
