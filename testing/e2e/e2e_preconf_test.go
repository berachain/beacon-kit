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

//go:build e2e_preconf

package e2e_test

import (
	"testing"

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

// PreconfE2ESuite is a test suite for preconfirmation functionality.
type PreconfE2ESuite struct {
	suite.KurtosisE2ESuite
}

// SetupSuite sets up the suite with preconf configuration.
func (s *PreconfE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions(suite.WithPreconfConfig())
}

// TestPreconfE2ESuite runs the preconf test suite.
func TestPreconfE2ESuite(t *testing.T) {
	suite.Run(t, new(PreconfE2ESuite))
}

// TestPreconfEndToEnd verifies the complete preconfirmation flow:
// 1. Sequencer serves payloads to whitelisted validators
// 2. Validators fetch payloads from sequencer
// 3. When sequencer is removed, validators fall back to local building
func (s *PreconfE2ESuite) TestPreconfEndToEnd() {
	// Wait for network to produce enough blocks for preconf to be exercised
	err := s.WaitForFinalizedBlockNumber(blocksToWait)
	s.Require().NoError(err, "Network should reach finalized blocks")

	sequencer := sequencerCLService
	fetchers := []string{config.ClientValidator0}

	// Step 1: Verify sequencer is serving payloads
	s.Run("SequencerServesPayloads", func() {
		logs, err := s.GetServiceLogs(sequencer)
		s.Require().NoError(err, "Should get sequencer logs")

		found := suite.ContainsLogMessage(logs, sequencerServingLog)
		s.Require().True(found,
			"Sequencer (%s) should serve payloads to validators. "+
				"Expected log message containing: %q",
			sequencer, sequencerServingLog)
	})

	// Step 2: Verify validators fetch from sequencer
	for _, validator := range fetchers {
		validator := validator // capture for closure
		s.Run("ValidatorFetches/"+validator, func() {
			logs, err := s.GetServiceLogs(validator)
			s.Require().NoError(err, "Should get validator logs for %s", validator)

			found := suite.ContainsLogMessage(logs, validatorFetchingLog)
			s.Require().True(found,
				"Validator (%s) should fetch payloads from sequencer. "+
					"Expected log message containing: %q",
				validator, validatorFetchingLog)
		})
	}

	// Step 3: Remove sequencer and verify fallback to local building
	s.Run("FallbackAfterSequencerRemoved", func() {
		// Get current block before removing sequencer
		currentBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get current block number")

		// Remove sequencer - simulates crash/unavailability
		s.T().Logf("Removing sequencer (%s)...", sequencer)
		err = s.RemoveService(sequencer)
		s.Require().NoError(err, "Should remove sequencer service")

		// Wait for more blocks - validators must build locally now
		targetBlock := currentBlock + blocksAfterFallback
		s.T().Logf("Waiting for %d more blocks (current: %d, target: %d)...",
			blocksAfterFallback, currentBlock, targetBlock)

		err = s.WaitForFinalizedBlockNumber(targetBlock)
		s.Require().NoError(err, "Network should continue producing blocks after sequencer removed")

		// Verify network continued
		finalBlock, err := s.RPCClient().BlockNumber(s.Ctx())
		s.Require().NoError(err, "Should get final block number")
		s.Require().GreaterOrEqual(finalBlock, targetBlock,
			"Network should have produced blocks after sequencer removed")

		s.T().Logf("Network continued: block %d → %d (fallback working)", currentBlock, finalBlock)
	})
}
