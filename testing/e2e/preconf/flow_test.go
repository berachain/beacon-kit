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
	sequencerServingLog   = "GetPayloadBySlot completed"
	validatorFetchingLog  = "Successfully fetched payload from sequencer"
	sequencerDownLog      = "Detected sequencer offline, starting health monitor"
	sequencerRecoveredLog = "Sequencer is healthy again, stopping health monitor"

	// Kurtosis service name for the dedicated sequencer CL node.
	sequencerCLService = "cl-sequencer-beaconkit-0"

	// Number of blocks to wait.
	blocksToWait        = 20
	blocksAfterFallback = 10
)

// TestSequencerFlow verifies the preconf pathway using ordered subtests.
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
		validators := []string{
			config.ClientValidator0, config.ClientValidator1, config.ClientValidator2,
			config.ClientValidator3, config.ClientValidator4,
		}
		for _, validator := range validators {
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

	// Stop sequencer and verify validators fall back to local building.
	s.Run("SequencerFallback", func() {
		// Get current block before stopping sequencer.
		currentBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get current block number")

		// Stop sequencer -- simulates crash/unavailability.
		s.T().Logf("Stopping sequencer (%s)...", sequencerCLService)
		err = s.StopService(sequencerCLService)
		s.Require().NoError(err, "Should stop sequencer service")

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

	// After sequencer have been stopped, verify each validator logs the circuit breaker
	// message at most once on their first proposal after the sequencer went down.
	s.Run("SequencerCircuitBreaker", func() {
		// Wait for enough additional blocks so each validator has had the chance
		// to propose at least twice since the sequencer was removed.
		currentBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get current block number")
		err = s.WaitForFinalizedBlockNumber(currentBlock + blocksAfterFallback)
		s.Require().NoError(err, "Network should continue producing blocks")

		validators := []string{
			config.ClientValidator0, config.ClientValidator1, config.ClientValidator2,
			config.ClientValidator3, config.ClientValidator4,
		}
		for _, validator := range validators {
			logs, err := s.GetServiceLogs(validator)
			s.Require().NoError(err, "Should get logs for %s", validator)

			count := suite.CountLogMessages(logs, sequencerDownLog)
			s.Require().LessOrEqual(count, 1,
				"Validator %s tripped circuit breaker %d times: should only detect sequencer down once",
				validator, count)
		}
	})

	// When sequencer gets restarted, verify each validator's monitor detects the recovery
	// and validators resume fetching payloads from the sequencer.
	s.Run("SequencerRecovery", func() {
		s.T().Logf("Restarting sequencer (%s)...", sequencerCLService)
		err := s.StartService(sequencerCLService)
		s.Require().NoError(err, "Should restart sequencer service")

		// Wait for enough blocks so each validator has had the chance to propose
		// at least once after the monitor detects recovery.
		currentBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get current block number")
		err = s.WaitForFinalizedBlockNumber(currentBlock + blocksAfterFallback)
		s.Require().NoError(err, "Network should continue producing blocks after sequencer restarted")

		// Validators should have detected recovery and resumed fetching from the sequencer.
		validators := []string{
			config.ClientValidator0, config.ClientValidator1, config.ClientValidator2,
			config.ClientValidator3, config.ClientValidator4,
		}
		for _, validator := range validators {
			logs, err := s.GetServiceLogs(validator)
			s.Require().NoError(err, "Should get logs for %s", validator)

			s.Require().True(suite.ContainsLogMessage(logs, sequencerRecoveredLog),
				"Validator %s health monitor should have detected sequencer recovery. "+
					"Expected log message containing: %q",
				validator, sequencerRecoveredLog)

			s.Require().True(suite.ContainsLogMessage(logs, validatorFetchingLog),
				"Validator %s should resume fetching from sequencer after recovery. "+
					"Expected log message containing: %q",
				validator, validatorFetchingLog)
		}
	})
}
