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

package deposit_test

import (
	"testing"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Mock implementations for Deposit.
type MockDeposit struct {
	index uint64
}

func (d MockDeposit) GetIndex() uint64 {
	return d.index
}

// Mock implementation of the Deposit interface.
func (d MockDeposit) New(
	_ crypto.BLSPubkey,
	_ any,
	_ math.U64,
	_ crypto.BLSSignature,
	index uint64,
) MockDeposit {
	return MockDeposit{index: index}
}

// Mock implementations for ExecutionPayload.
type MockExecutionPayload struct {
	number math.U64
}

func (ep MockExecutionPayload) GetNumber() math.U64 {
	return ep.number
}

// Mock implementations for BeaconBlockBody.
type MockBeaconBlockBody struct {
	deposits         []MockDeposit
	executionPayload MockExecutionPayload
}

func (b MockBeaconBlockBody) GetDeposits() []MockDeposit {
	return b.deposits
}

func (b MockBeaconBlockBody) GetExecutionPayload() MockExecutionPayload {
	return b.executionPayload
}

// Mock implementations for BeaconBlock.
type MockBeaconBlock struct {
	slot math.U64
	body MockBeaconBlockBody
}

func (b MockBeaconBlock) GetSlot() math.U64 {
	return b.slot
}

func (b MockBeaconBlock) GetBody() MockBeaconBlockBody {
	return b.body
}

// Mock implementations for BlockEvent.
type MockBlockEvent struct {
	data MockBeaconBlock
}

func (e MockBlockEvent) Type() asynctypes.EventID {
	return ""
}

func (e MockBlockEvent) Is(_ asynctypes.EventID) bool {
	// Mock implementation for Is method. Adjust logic as needed.
	return true
}

func (e MockBlockEvent) Data() MockBeaconBlock {
	return e.data
}

// Unit tests for BuildPruneRangeFn.
//
//nolint:lll
func TestBuildPruneRangeFn(t *testing.T) {
	tests := []struct {
		name          string
		maxDeposits   uint64
		deposits      []MockDeposit
		expectedStart uint64
		expectedEnd   uint64
	}{
		{
			name:          "No deposits",
			maxDeposits:   10,
			deposits:      []MockDeposit{},
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Less than max deposits",
			maxDeposits:   10,
			deposits:      []MockDeposit{{index: 5}},
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "Equal to max deposits",
			maxDeposits:   10,
			deposits:      []MockDeposit{{index: 10}},
			expectedStart: 0,
			expectedEnd:   10,
		},
		{
			name:          "More than max deposits",
			maxDeposits:   10,
			deposits:      []MockDeposit{{index: 15}},
			expectedStart: 5,
			expectedEnd:   10,
		},
		{
			name:          "Boundary case at max deposits",
			maxDeposits:   5,
			deposits:      []MockDeposit{{index: 5}},
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "Single deposit",
			maxDeposits:   3,
			deposits:      []MockDeposit{{index: 1}},
			expectedStart: 0,
			expectedEnd:   1,
		},
		{
			name:          "Zero max deposits",
			maxDeposits:   0,
			deposits:      []MockDeposit{{index: 1}},
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Multiple deposits with varying indices",
			maxDeposits:   10,
			deposits:      []MockDeposit{{index: 3}, {index: 6}, {index: 9}, {index: 12}},
			expectedStart: 2,
			expectedEnd:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := chain.NewChainSpec(
				chain.SpecData[
					bytes.B4, math.U64, common.Address, math.U64, any,
				]{
					MaxDepositsPerBlock: tt.maxDeposits,
				},
			)
			pruneFn := deposit.BuildPruneRangeFn[
				MockBeaconBlockBody,
				MockBeaconBlock,
				MockBlockEvent](cs)
			event := MockBlockEvent{
				data: MockBeaconBlock{
					body: MockBeaconBlockBody{
						deposits: tt.deposits,
					},
				},
			}
			start, end := pruneFn(event)
			require.Equal(t, tt.expectedStart, start, "Test case: %s (expectedStart)", tt.name)
			require.Equal(t, tt.expectedEnd, end, "Test case: %s (expectedEnd)", tt.name)
		})
	}
}
