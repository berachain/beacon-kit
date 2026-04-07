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
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/ethereum/go-ethereum/common"
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

	// Interval between txs in the round-change tx burst.
	roundChangeTxInterval = 100 * time.Millisecond
)

// TestRoundRotationRebuild verifies that when a validator misses its proposal
// turn (causing a consensus round rotation), the sequencer detects the round
// change via the CometBFT EventBus and triggers a fresh payload build that
// picks up new mempool transactions.
//
// Strategy:
//  1. Let the preconf network warm up and produce blocks normally.
//  2. Stop one validator's CL node so it misses its next proposal.
//  3. While waiting, continuously submit txs to the preconf RPC so the
//     mempool always has fresh content for the rebuilt payload to pick up.
//  4. CometBFT times out on that round and selects a new proposer (round > 0).
//  5. The sequencer should detect this via EventNewRound and rebuild the payload.
//  6. Verify (a) sequencer logs show the rebuild, (b) flashblock monitor saw
//     multiple payload IDs at index 0, and (c) at least one of our txs landed
//     in a block that was rebuilt.
//  7. Restart the stopped validator and confirm the network continues.
func (s *PreconfE2ESuite) TestRoundRotationRebuild() {
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

	// Step 3: Start a background tx sender so the mempool stays populated
	// across the round-timeout window. We collect every successfully sent
	// hash for later inclusion-block verification.
	preconfClient := s.PreconfRPCClients(0)
	sender := s.TestAccounts()[0]
	receiver := s.TestAccounts()[1]
	startNonce, err := preconfClient.NonceAt(s.Ctx(), sender.Address(), nil)
	s.Require().NoError(err, "Should get sender nonce")
	gasTipCap, gasFeeCap := s.suggestGasCaps(preconfClient.Client)

	txCtx, cancelTxs := context.WithCancel(s.Ctx())
	defer cancelTxs()
	var sentHashes sync.Map // common.Hash -> struct{}
	var senderWG sync.WaitGroup
	senderWG.Add(1)
	go func() {
		defer senderWG.Done()
		nonce := startNonce
		ticker := time.NewTicker(roundChangeTxInterval)
		defer ticker.Stop()
		for {
			select {
			case <-txCtx.Done():
				return
			case <-ticker.C:
				tx, sErr := s.sendETHTransfer(transferParams{
					client:    preconfClient.Client,
					sender:    sender,
					to:        receiver.Address(),
					nonce:     nonce,
					amount:    big.NewInt(1),
					gasTipCap: gasTipCap,
					gasFeeCap: gasFeeCap,
				})
				if sErr != nil {
					continue // retry same nonce on next tick
				}
				sentHashes.Store(tx.Hash(), struct{}{})
				nonce++
			}
		}
	}()

	// Step 4: Wait for enough blocks so the stopped validator misses at
	// least one proposal turn.
	elClient := s.ExecutionClients(0)
	currentBlock, err := elClient.BlockNumber(s.Ctx())
	s.Require().NoError(err)

	targetBlock := currentBlock + blocksAfterStop
	s.T().Logf("Waiting for %d blocks (current: %d, target: %d)...",
		blocksAfterStop, currentBlock, targetBlock)

	err = s.WaitForFinalizedBlockNumber(targetBlock)
	s.Require().NoError(err, "Network should continue producing blocks with one validator down")

	// Stop the tx sender and let in-flight receipts settle.
	cancelTxs()
	senderWG.Wait()

	// Step 5: Verify the sequencer detected the round change and rebuilt.
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

	// Step 5b: Verify the flashblock monitor received flashblocks from the
	// rebuild. A rebuild produces new flashblocks with a different payload_id
	// for the same block_number, so we look for any block_number that has
	// multiple distinct payload_id values at index 0.
	var rebuiltBlocks map[uint64]struct{}
	s.Run("FlashblockMonitorReceivedRebuild", func() {
		logs, err := s.GetServiceLogsWithOptions(flashblockMonitorService, 10000, suite.DefaultLogCollectionTimeout)
		s.Require().NoError(err, "Should get flashblock monitor logs")

		rebuiltBlocks = flashblockRebuiltBlocks(logs)
		s.T().Logf("Flashblock monitor saw %d block(s) with multiple payload IDs (rebuild)", len(rebuiltBlocks))
		s.Require().NotEmpty(rebuiltBlocks,
			"Flashblock monitor should have seen at least one block_number "+
				"with multiple payload_id values, indicating a round-change rebuild")
	})

	// Step 5c: Verify at least one of our submitted txs landed in a rebuilt
	// block, proving the rebuilt payload picks up new mempool content rather
	// than just re-emitting the stale round-0 payload.
	s.Run("RebuiltPayloadIncludesNewTxs", func() {
		s.Require().NotEmpty(rebuiltBlocks, "rebuilt block set must be populated by previous subtest")

		var sentList []common.Hash
		sentHashes.Range(func(k, _ any) bool {
			sentList = append(sentList, k.(common.Hash))
			return true
		})
		s.Require().NotEmpty(sentList, "Background sender should have submitted at least one tx")
		s.T().Logf("Submitted %d txs during the round-change window", len(sentList))

		var includedInRebuild int
		for _, h := range sentList {
			receipt, rErr := elClient.TransactionReceipt(s.Ctx(), h)
			if rErr != nil || receipt == nil {
				continue
			}
			if _, ok := rebuiltBlocks[receipt.BlockNumber.Uint64()]; ok {
				includedInRebuild++
			}
		}
		s.T().Logf("%d/%d submitted txs landed in a rebuilt block", includedInRebuild, len(sentList))
		s.Require().Greater(includedInRebuild, 0,
			"At least one submitted tx should have landed in a block that was rebuilt on a round change")
	})

	// Step 6: Restart the stopped validator and confirm the network recovers.
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

// flashblockRebuiltBlocks parses the flashblock monitor logs and returns the
// set of block_numbers that had multiple distinct payload_id values at index
// 0. Each entry in the result represents a block where the sequencer issued
// more than one base flashblock — i.e. a round-change rebuild.
func flashblockRebuiltBlocks(logs []string) map[uint64]struct{} {
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

	rebuilt := make(map[uint64]struct{})
	for blockNum, payloadIDs := range basePayloads {
		if len(payloadIDs) > 1 {
			rebuilt[blockNum] = struct{}{}
		}
	}
	return rebuilt
}
