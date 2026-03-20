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
	"fmt"
	"strings"
	"time"

	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
)

const (
	sequencerWhitelistReloadLog = "Preconf whitelist reloaded"
	sequencerNotWhitelistedLog  = "is_whitelisted=false"

	whitelistPath = "/root/preconf/whitelist.json"

	// How long to wait for one round-rotation block (Commit.Round > 0) to appear.
	roundRotationTimeout = 90 * time.Second

	// The number of additional blocks required to confirm the chain is healthy.
	blocksAfterRotation = 10
)

// TestPreconf_RoundRotation_NonWhitelistedValidator verifies the scenario
// where a NON-whitelisted validator fails to propose and the round skips to a
// whitelisted validator instead.
func (s *PreconfE2ESuite) TestPreconf_RoundRotation_NonWhitelistedValidator() {
	err := s.WaitForFinalizedBlockNumber(blocksToWait)
	s.Require().NoError(err, "Network should reach initial finalized blocks")

	var originalWhitelist []string
	stoppedValidator := config.ValidatorConsensusClientName(0)

	// Unwhitelist first validator.
	s.Run("UnwhitelistValidator", func() {
		val0Pubkey, err := s.getValidatorBLSPubkey(stoppedValidator)
		s.Require().NoError(err, "Should get validator 0 BLS pubkey from key file")
		s.T().Logf("Validator 0 BLS pubkey: %s", val0Pubkey)

		originalWhitelist, err = s.readSequencerWhitelist()
		s.Require().NoError(err, "Should read current whitelist from sequencer")
		s.Require().NotEmpty(originalWhitelist, "Whitelist must not be empty at start")
		s.T().Logf("Original whitelist has %d entries", len(originalWhitelist))

		// Remove validator 0's pubkey.
		reducedWhitelist := make([]string, 0, len(originalWhitelist)-1)
		for _, pk := range originalWhitelist {
			if pk != val0Pubkey {
				reducedWhitelist = append(reducedWhitelist, pk)
			}
		}
		s.Require().Len(reducedWhitelist, len(originalWhitelist)-1,
			"Validator 0 pubkey must have been found in and removed from the whitelist")

		s.Require().NoError(s.writeSequencerWhitelist(reducedWhitelist))
		s.Require().NoError(s.sighupSequencer())
		s.T().Logf("Whitelist narrowed to %d entries; validator 0 is now non-whitelisted", len(reducedWhitelist))

		s.Require().Eventually(func() bool {
			logs, logErr := s.GetServiceLogs(sequencerCLService)
			return logErr == nil && suite.ContainsLogMessage(logs, sequencerWhitelistReloadLog)
		}, 15*time.Second, 500*time.Millisecond,
			"Sequencer must log a successful whitelist reload after SIGHUP")
	})

	// Stop the non-whitelisted validator.
	s.Run("StopValidator", func() {
		s.T().Logf("Stopping %s to induce round rotation...", stoppedValidator)
		s.Require().NoError(s.StopService(stoppedValidator))
	})

	currentBlock, err := s.ExecutionClients(0).BlockNumber(s.Ctx())
	s.Require().NoError(err)

	// When stopped non-whitelisted validator is selected as proposer, round rotation must occur.
	s.Run("RoundRotationWithNonWhitelistedProposer", func() {
		rotated, rotatedHeight := s.waitForRoundRotationBlock(currentBlock+1, roundRotationTimeout)
		s.Require().True(rotated,
			"At least one block must be committed at Round > 0 after stopping validator 0")
		s.T().Logf("Round rotation observed at block height %d", rotatedHeight)

		// Poll for the log: Kurtosis log propagation may lag the block commit by a few seconds.
		s.Require().Eventually(func() bool {
			seqLogs, logErr := s.GetServiceLogs(sequencerCLService)
			return logErr == nil && suite.ContainsLogMessage(seqLogs, sequencerNotWhitelistedLog)
		}, 15*time.Second, 500*time.Millisecond,
			"Sequencer must log %q for the non-whitelisted proposer that caused the rotation",
			sequencerNotWhitelistedLog)
	})

	// Restore original environment: all validators are whitelisted and the stopped one is restarted.
	s.Run("RestoreWhitelistAndRestart", func() {
		// Snapshot the non-whitelisted log count before restoring so we can assert it stops growing.
		logsBefore, logErr := s.GetServiceLogs(sequencerCLService)
		s.Require().NoError(logErr)
		falseCountBefore := suite.CountLogMessages(logsBefore, sequencerNotWhitelistedLog)

		s.Require().NoError(s.writeSequencerWhitelist(originalWhitelist))
		s.Require().NoError(s.sighupSequencer())

		s.T().Logf("Restarting %s...", stoppedValidator)
		s.Require().NoError(s.StartService(stoppedValidator))

		s.Require().NoError(s.WaitForNBlockNumbers(blocksAfterRotation),
			"Chain must continue after full whitelist is restored and validator restarts")

		// After blocksAfterRotation new blocks, the count must be the same.
		logsAfter, logErr := s.GetServiceLogs(sequencerCLService)
		s.Require().NoError(logErr)
		falseCountAfter := suite.CountLogMessages(logsAfter, sequencerNotWhitelistedLog)
		s.Require().Equal(falseCountBefore, falseCountAfter,
			"No new %q entries must appear after the whitelist is restored "+
				"(got %d before restore, %d after — delta indicates non-whitelisted validators remain)",
			sequencerNotWhitelistedLog, falseCountBefore, falseCountAfter)
	})
}

// getValidatorBLSPubkey returns the "0x<hex>" formatted BLS pubkey of the given validator service.
func (s *PreconfE2ESuite) getValidatorBLSPubkey(validatorService string) (string, error) {
	// The premined deposit JSON has a "pubkey" field in 0x<hex> format used by the sequencer whitelist.
	output, err := s.execOnService(validatorService,
		[]string{"sh", "-c", "jq -r '.pubkey' /root/.beacond/config/premined-deposits/premined-deposit-*.json | head -1"})
	if err != nil {
		return "", fmt.Errorf("read premined deposit pubkey from %s: %w", validatorService, err)
	}
	pubkey := strings.TrimSpace(output)
	if pubkey == "" {
		return "", fmt.Errorf("empty pubkey from premined deposit on %s", validatorService)
	}
	return pubkey, nil
}

// readSequencerWhitelist returns the list of BLS pubkeys in the sequencer whitelist.
func (s *PreconfE2ESuite) readSequencerWhitelist() ([]string, error) {
	output, err := s.execOnService(sequencerCLService, []string{"cat", whitelistPath})
	if err != nil {
		return nil, err
	}
	var pubkeys []string
	if err = json.Unmarshal([]byte(output), &pubkeys); err != nil {
		return nil, fmt.Errorf("parse whitelist JSON: %w", err)
	}
	return pubkeys, nil
}

// writeSequencerWhitelist serialises pubkeys to the sequencer whitelist.
func (s *PreconfE2ESuite) writeSequencerWhitelist(pubkeys []string) error {
	data, err := json.Marshal(pubkeys)
	if err != nil {
		return fmt.Errorf("marshal whitelist: %w", err)
	}
	// Use printf to avoid shell interpretation of the JSON content.
	cmd := fmt.Sprintf("printf '%%s' '%s' > %s", string(data), whitelistPath)
	_, err = s.execOnService(sequencerCLService, []string{"sh", "-c", cmd})
	return err
}

// sighupSequencer SIGHUPs the sequencer beacond process to trigger whitelist reload.
func (s *PreconfE2ESuite) sighupSequencer() error {
	_, err := s.execOnService(sequencerCLService, []string{"sh", "-c", "kill -HUP $(pgrep beacond)"})
	return err
}

// execOnService runs cmd inside the named Kurtosis service container.
// Returns the combined stdout output.
func (s *PreconfE2ESuite) execOnService(serviceName string, cmd []string) (string, error) {
	sCtx, err := s.Enclave().GetServiceContext(serviceName)
	if err != nil {
		return "", fmt.Errorf("get service context for %s: %w", serviceName, err)
	}
	exitCode, output, err := sCtx.ExecCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("exec %v on %s: %w", cmd, serviceName, err)
	}
	if exitCode != 0 {
		return "", fmt.Errorf("exec %v on %s exited %d: %s", cmd, serviceName, exitCode, output)
	}
	return output, nil
}

// waitForRoundRotationBlock polls CometBFT commit data starting at startHeight
// until it finds a block committed at Round > 0 or the timeout expires.
func (s *PreconfE2ESuite) waitForRoundRotationBlock(
	startHeight uint64,
	timeout time.Duration,
) (found bool, height uint64) {
	cometClient := s.ConsensusClients(1) // validator 1 stays alive throughout
	s.Require().NotNil(cometClient, "Consensus client 1 must be available")

	ctx, cancel := context.WithTimeout(s.Ctx(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	nextHeight := startHeight
	for {
		select {
		case <-ctx.Done():
			return false, 0
		case <-ticker.C:
		}

		h := int64(nextHeight) //nolint:gosec // heights in tests won't overflow int64
		result, err := cometClient.Commit(ctx, &h)
		if err != nil {
			continue // block not yet finalized
		}

		round := result.SignedHeader.Commit.Round
		s.T().Logf("Height %d committed at Round %d", nextHeight, round)
		if round > 0 {
			return true, nextHeight
		}
		nextHeight++
	}
}
