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

package chain

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/assert"
)

// TestActiveForkVersionForEpoch tests the ActiveForkVersionForEpoch method
func TestActiveForkVersionForEpoch(t *testing.T) {
	// Define concrete types for the generic parameters
	type domainTypeT [4]byte
	type epochT uint64
	type executionAddressT [20]byte
	type slotT uint64
	type cometBFTConfigT struct{}

	// Create an instance of chainSpec with test data
	spec := chainSpec[
		domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
	]{
		Data: SpecData[
			domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
		]{
			ElectraForkEpoch: 10,
		},
	}

	// Define test cases
	testCases := []struct {
		epoch    epochT
		expected uint32
	}{
		{epoch: 9, expected: version.Deneb},
		{epoch: 10, expected: version.Electra},
		{epoch: 11, expected: version.Electra},
	}

	// Run test cases
	for _, tc := range testCases {
		result := spec.ActiveForkVersionForEpoch(tc.epoch)
		assert.Equal(t, tc.expected, result, "unexpected active fork version for epoch %d", tc.epoch)
	}
}

// TestSlotToEpoch tests the SlotToEpoch method.
func TestSlotToEpoch(t *testing.T) {
	// Define concrete types for the generic parameters
	type domainTypeT [4]byte
	type epochT uint64
	type executionAddressT [20]byte
	type slotT uint64
	type cometBFTConfigT struct{}

	// Create an instance of chainSpec with test data
	spec := chainSpec[
		domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
	]{
		Data: SpecData[
			domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
		]{
			SlotsPerEpoch: 32,
		},
	}

	// Define test cases
	testCases := []struct {
		slot     slotT
		expected epochT
	}{
		{slot: 0, expected: 0},
		{slot: 31, expected: 0},
		{slot: 32, expected: 1},
		{slot: 63, expected: 1},
		{slot: 64, expected: 2},
		{slot: 95, expected: 2},
	}

	// Run test cases
	for _, tc := range testCases {
		result := spec.SlotToEpoch(tc.slot)
		assert.Equal(t, tc.expected, result, "unexpected epoch for slot %d", tc.slot)
	}
}

// TestActiveForkVersionForSlot tests the ActiveForkVersionForSlot method.
func TestActiveForkVersionForSlot(t *testing.T) {
	// Define concrete types for the generic parameters
	type domainTypeT [4]byte
	type epochT uint64
	type executionAddressT [20]byte
	type slotT uint64
	type cometBFTConfigT struct{}

	// Create an instance of chainSpec with test data
	spec := chainSpec[
		domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
	]{
		Data: SpecData[
			domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
		]{
			SlotsPerEpoch:    32,
			ElectraForkEpoch: 10,
		},
	}

	// Define test cases
	testCases := []struct {
		slot     slotT
		expected uint32
	}{
		{slot: 0, expected: version.Deneb},
		{slot: 319, expected: version.Deneb},
		{slot: 320, expected: version.Electra},
		{slot: 640, expected: version.Electra},
	}

	// Run test cases
	for _, tc := range testCases {
		result := spec.ActiveForkVersionForSlot(tc.slot)
		assert.Equal(t, tc.expected, result, "unexpected fork version for slot %d", tc.slot)
	}
}

// TestWithinDAPeriod tests the WithinDAPeriod method.
func TestWithinDAPeriod(t *testing.T) {
	// Define concrete types for the generic parameters
	type domainTypeT [4]byte
	type epochT uint64
	type executionAddressT [20]byte
	type slotT uint64
	type cometBFTConfigT struct{}

	// Create an instance of chainSpec with test data
	spec := chainSpec[
		domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
	]{
		Data: SpecData[
			domainTypeT, epochT, executionAddressT, slotT, cometBFTConfigT,
		]{
			SlotsPerEpoch:                    32,
			MinEpochsForBlobsSidecarsRequest: 5,
		},
	}

	// Define test cases
	testCases := []struct {
		block    slotT
		current  slotT
		expected bool
	}{
		{block: 0, current: 160, expected: true},    // Block is within DA period (5 epochs)
		{block: 0, current: 192, expected: false},   // Block is outside DA period (>5 epochs)
		{block: 160, current: 320, expected: true},  // Block is within DA period
		{block: 160, current: 352, expected: false}, // Block is outside DA period
	}

	// Run test cases
	for _, tc := range testCases {
		result := spec.WithinDAPeriod(tc.block, tc.current)
		assert.Equal(t, tc.expected, result, "unexpected DA period result for block %d and current %d", tc.block, tc.current)
	}
}
