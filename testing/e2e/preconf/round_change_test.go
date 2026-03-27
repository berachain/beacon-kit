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
	"encoding/json"
	"strings"

	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
)

const (
	// Log message emitted by the sequencer's blockchain service when
	// a round change triggers a fresh payload build.
	roundChangeRebuildLog = "Round change detected: rebuilding payload for whitelisted proposer"

	// Log message emitted by the CometBFT service layer when it
	// receives an EventNewRound with round > 0.
	roundChangeDetectedLog = "CometBFT round change detected"

	// Number of blocks to let the network run before stopping a validator.
	warmupBlocks = 5

	// Number of blocks to wait after stopping a validator. Must be enough
	// for the stopped validator to miss at least one proposal turn.
	blocksAfterStop = 10

	// Kurtosis service name for the flashblock monitor.
	flashblockMonitorService = "flashblock-monitor"
)

// TestRoundChangeRebuild verifies that when a validator misses its proposal
// turn (causing a consensus round timeout), the sequencer detects the round
// change via the CometBFT EventBus and triggers a fresh payload build.
//
// Strategy:
//  1. Let the preconf network warm up and produce blocks normally.
//  2. Stop one validator's CL node so it misses its next proposal.
//  3. CometBFT times out on that round and selects a new proposer (round > 0).
//  4. The sequencer should detect this via EventNewRound and rebuild the payload.
//  5. Verify by scanning the sequencer's logs for the rebuild message.
//  6. Restart the stopped validator and confirm the network continues.
func (s *PreconfE2ESuite) TestRoundChangeRebuild() {
	// Step 1: Let the network warm up.
	err := s.WaitForFinalizedBlockNumber(warmupBlocks)
	s.Require().NoError(err, "Network should produce warmup blocks")

	// Step 2: Stop one validator's CL to force it to miss proposals.
	// When it's this validator's turn to propose, CometBFT will time out
	// and enter round > 0 with a different proposer.
	stoppedValidator := config.ValidatorConsensusClientName(0)
	s.T().Logf("Stopping validator %s to force round timeouts...", stoppedValidator)
	err = s.StopService(stoppedValidator)
	s.Require().NoError(err, "Should stop validator")

	// Step 3: Wait for enough blocks so the stopped validator misses at
	// least one proposal turn.
	elClient := s.ExecutionClients(0)
	currentBlock, err := elClient.BlockNumber(s.Ctx())
	s.Require().NoError(err)

	targetBlock := currentBlock + blocksAfterStop
	s.T().Logf("Waiting for %d blocks (current: %d, target: %d)...",
		blocksAfterStop, currentBlock, targetBlock)

	err = s.WaitForFinalizedBlockNumber(targetBlock)
	s.Require().NoError(err, "Network should continue producing blocks with one validator down")

	// Step 4: Verify the sequencer detected the round change and rebuilt.
	s.Run("SequencerDetectsRoundChange", func() {
		logs, err := s.GetServiceLogs(sequencerCLService)
		s.Require().NoError(err, "Should get sequencer logs")

		s.Require().True(suite.ContainsLogMessage(logs, roundChangeDetectedLog),
			"Sequencer should have detected a CometBFT round change. "+
				"Expected log: %q", roundChangeDetectedLog)

		s.Require().True(suite.ContainsLogMessage(logs, roundChangeRebuildLog),
			"Sequencer should have rebuilt payload on round change. "+
				"Expected log: %q", roundChangeRebuildLog)
	})

	// Step 4b: Verify the flashblock monitor received flashblocks from the
	// rebuild. A rebuild produces new flashblocks with a different payload_id
	// for the same block_number, so we look for any block_number that has
	// multiple distinct payload_id values at index 0.
	s.Run("FlashblockMonitorReceivedRebuild", func() {
		logs, err := s.GetServiceLogsWithOptions(flashblockMonitorService, 10000, suite.DefaultLogCollectionTimeout)
		s.Require().NoError(err, "Should get flashblock monitor logs")

		rebuilds := countFlashblockRebuilds(logs)
		s.T().Logf("Flashblock monitor saw %d block(s) with multiple payload IDs (rebuild)", rebuilds)
		s.Require().Greater(rebuilds, 0,
			"Flashblock monitor should have seen at least one block_number "+
				"with multiple payload_id values, indicating a round-change rebuild")
	})

	// Step 5: Restart the stopped validator and confirm the network recovers.
	s.Run("NetworkRecovery", func() {
		s.T().Logf("Restarting validator %s...", stoppedValidator)
		err := s.StartService(stoppedValidator)
		s.Require().NoError(err, "Should restart validator")

		currentBlock, err := elClient.BlockNumber(s.Ctx())
		s.Require().NoError(err)
		err = s.WaitForFinalizedBlockNumber(currentBlock + warmupBlocks)
		s.Require().NoError(err, "Network should continue producing blocks after validator restart")
	})
}

// flashblockEntry is the minimal structure we parse from the monitor's raw
// JSON output. Each line is a serialized BerachainFlashblockPayload.
type flashblockEntry struct {
	PayloadID string `json:"payload_id"`
	Index     uint64 `json:"index"`
	Metadata  struct {
		BlockNumber uint64 `json:"block_number"`
	} `json:"metadata"`
}

// countFlashblockRebuilds parses the flashblock monitor logs and counts how
// many block_numbers had multiple distinct payload_id values at index 0.
// A count > 0 means the sequencer rebuilt at least one payload (round change).
func countFlashblockRebuilds(logs []string) int {
	// Map block_number → set of payload_ids seen at index 0.
	basePayloads := make(map[uint64]map[string]struct{})

	for _, line := range logs {
		// The monitor logs contain both setup messages and raw JSON.
		// Only attempt to parse lines that look like JSON.
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "{") {
			continue
		}

		var entry flashblockEntry
		if err := json.Unmarshal([]byte(trimmed), &entry); err != nil {
			continue
		}

		// Only look at index-0 flashblocks (the "base" of each sequence).
		if entry.Index != 0 {
			continue
		}

		blockNum := entry.Metadata.BlockNumber
		if basePayloads[blockNum] == nil {
			basePayloads[blockNum] = make(map[string]struct{})
		}
		basePayloads[blockNum][entry.PayloadID] = struct{}{}
	}

	// Count blocks that had more than one payload_id at index 0.
	rebuilds := 0
	for _, payloadIDs := range basePayloads {
		if len(payloadIDs) > 1 {
			rebuilds++
		}
	}
	return rebuilds
}
