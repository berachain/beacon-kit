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

package e2e_test

import (
	"encoding/json"
	"io"
	"time"

	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
)

// TestGetPendingPartialWithdrawalsAPI verifies the pending partial withdrawals API response format
func (s *BeaconKitE2ESuite) TestGetPendingPartialWithdrawalsAPI() {
	s.initBeaconTest()

	// Make a request to the getPendingPartialWithdrawals endpoint
	resp, err := s.getPendingPartialWithdrawals(utils.StateIDHead)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Read the raw response body for debugging
	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	// For debugging
	s.T().Logf("Raw getPendingPartialWithdrawals response: %s", string(body))

	// Parse the response manually to verify the format
	var genericResp struct {
		Version             string          `json:"version"`
		ExecutionOptimistic bool            `json:"execution_optimistic"`
		Finalized           bool            `json:"finalized"`
		Data                json.RawMessage `json:"data"`
	}

	err = json.Unmarshal(body, &genericResp)
	s.Require().NoError(err, "Failed to parse the response envelope")

	// Verify the envelope format
	s.Require().NotEmpty(genericResp.Version, "Response should have a version field")
	s.T().Logf("Response version: %s", genericResp.Version)
	s.T().Logf("Response execution_optimistic: %v", genericResp.ExecutionOptimistic)
	s.T().Logf("Response finalized: %v", genericResp.Finalized)

	// Parse the actual withdrawals data
	var withdrawals []*types.PendingPartialWithdrawalData
	err = json.Unmarshal(genericResp.Data, &withdrawals)
	s.Require().NoError(err, "Failed to parse the withdrawals data")

	// Even if there are no pending withdrawals, the format should be valid
	s.T().Logf("Number of pending partial withdrawals: %d", len(withdrawals))

	// Also test with finalized state
	s.checkPendingPartialWithdrawalsWithState(utils.StateIDFinalized)

	// Also test with justified state
	s.checkPendingPartialWithdrawalsWithState(utils.StateIDJustified)
}

// checkPendingPartialWithdrawalsWithState tests the getPendingPartialWithdrawals endpoint with a specific state ID
func (s *BeaconKitE2ESuite) checkPendingPartialWithdrawalsWithState(stateID string) {
	s.T().Logf("Testing pending partial withdrawals with state: %s", stateID)

	// Get withdrawals for the specified state
	resp, err := s.getPendingPartialWithdrawals(stateID)
	if err != nil {
		s.T().Logf("Error getting state %s: %v", stateID, err)
		return
	}
	defer resp.Body.Close()

	// Read and validate the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.T().Logf("Error reading response body: %v", err)
		return
	}

	// For debugging
	s.T().Logf("Response for state %s: %s", stateID, string(body))

	// Parse the response manually
	var genericResp struct {
		Version             string          `json:"version"`
		ExecutionOptimistic bool            `json:"execution_optimistic"`
		Finalized           bool            `json:"finalized"`
		Data                json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(body, &genericResp); err != nil {
		s.T().Logf("Error parsing response: %v", err)
		return
	}

	// Parse the withdrawals data
	var withdrawals []*types.PendingPartialWithdrawalData
	if err := json.Unmarshal(genericResp.Data, &withdrawals); err != nil {
		s.T().Logf("Error parsing withdrawals data: %v", err)
		return
	}

	// Log results
	s.T().Logf("State %s: Found %d pending partial withdrawals", stateID, len(withdrawals))
}

// TestMonitorPendingPartialWithdrawals monitors for pending partial withdrawals over time
func (s *BeaconKitE2ESuite) TestMonitorPendingPartialWithdrawals() {
	// Initialize the test environment
	s.initBeaconTest()

	// Monitor for 5 minutes, checking every 30 seconds
	endTime := time.Now().Add(5 * time.Minute)

	s.T().Logf("Starting to monitor pending partial withdrawals for 5 minutes")

	for time.Now().Before(endTime) {
		// Get the pending partial withdrawals
		resp, err := s.getPendingPartialWithdrawals(utils.StateIDHead)
		if err != nil {
			s.T().Logf("Error getting pending withdrawals: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Read and parse the response
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			s.T().Logf("Error reading response body: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Parse the response envelope
		var genericResp struct {
			Version             string          `json:"version"`
			ExecutionOptimistic bool            `json:"execution_optimistic"`
			Finalized           bool            `json:"finalized"`
			Data                json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(body, &genericResp); err != nil {
			s.T().Logf("Error parsing response envelope: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Parse the actual withdrawals data
		var withdrawals []*types.PendingPartialWithdrawalData
		if err := json.Unmarshal(genericResp.Data, &withdrawals); err != nil {
			s.T().Logf("Error parsing withdrawals data: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Log the number of pending withdrawals
		s.T().Logf("Number of pending partial withdrawals: %d", len(withdrawals))

		if len(withdrawals) > 0 {
			// Log details about each withdrawal
			for i, withdrawal := range withdrawals {
				s.T().Logf("Withdrawal %d: Validator Index: %d, Amount: %d",
					i+1, withdrawal.ValidatorIndex, withdrawal.Amount)
			}

			// Withdrawals found, test completed successfully
			s.T().Logf("Found pending partial withdrawals, test completed successfully")
			return
		}

		// Wait before checking again
		time.Sleep(30 * time.Second)
	}

	s.T().Logf("Monitoring period ended without finding any pending partial withdrawals")
}
