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

package store_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Mock implementations for BeaconBlock.
type MockBeaconBlock struct {
	slot math.U64
}

// Mock implementations GetSlot() for BeaconBlock.
func (b MockBeaconBlock) GetSlot() math.U64 {
	return b.slot
}

// Mock implementations for BlockEvent.
type MockBlockEvent struct {
	data MockBeaconBlock
}

// Mock implementations Data() for BlockEvent.
func (e MockBlockEvent) Data() MockBeaconBlock {
	return e.data
}

// Mock ChainSpec implementation.
type MockChainSpec struct {
	slotsPerEpoch                    uint64
	minEpochsForBlobsSidecarsRequest uint64
}

// SlotsPerEpoch implements chain.Spec.
func (cs MockChainSpec) SlotsPerEpoch() uint64 {
	return cs.slotsPerEpoch
}

// MinEpochsForBlobsSidecarsRequest implements chain.Spec.
func (cs MockChainSpec) MinEpochsForBlobsSidecarsRequest() uint64 {
	return cs.minEpochsForBlobsSidecarsRequest
}

// ActiveForkVersionForEpoch implements chain.Spec.
func (cs *MockChainSpec) ActiveForkVersionForEpoch(_ math.U64) uint32 {
	panic("unimplemented")
}

// ActiveForkVersionForSlot implements chain.Spec.
func (cs *MockChainSpec) ActiveForkVersionForSlot(_ math.U64) uint32 {
	panic("unimplemented")
}

// BytesPerBlob implements chain.Spec.
func (cs *MockChainSpec) BytesPerBlob() uint64 {
	panic("unimplemented")
}

// DepositContractAddress implements chain.Spec.
func (cs *MockChainSpec) DepositContractAddress() common.Address {
	panic("unimplemented")
}

// DepositEth1ChainID implements chain.Spec.
func (cs *MockChainSpec) DepositEth1ChainID() uint64 {
	panic("unimplemented")
}

// DomainTypeAggregateAndProof implements chain.Spec.
func (cs *MockChainSpec) DomainTypeAggregateAndProof() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeApplicationMask implements chain.Spec.
func (cs *MockChainSpec) DomainTypeApplicationMask() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeAttester implements chain.Spec.
func (cs *MockChainSpec) DomainTypeAttester() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeDeposit implements chain.Spec.
func (cs *MockChainSpec) DomainTypeDeposit() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeProposer implements chain.Spec.
func (cs *MockChainSpec) DomainTypeProposer() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeRandao implements chain.Spec.
func (cs *MockChainSpec) DomainTypeRandao() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeSelectionProof implements chain.Spec.
func (cs *MockChainSpec) DomainTypeSelectionProof() bytes.B4 {
	panic("unimplemented")
}

// DomainTypeVoluntaryExit implements chain.Spec.
func (cs *MockChainSpec) DomainTypeVoluntaryExit() bytes.B4 {
	panic("unimplemented")
}

// EffectiveBalanceIncrement implements chain.Spec.
func (cs *MockChainSpec) EffectiveBalanceIncrement() uint64 {
	panic("unimplemented")
}

// EjectionBalance implements chain.Spec.
func (cs *MockChainSpec) EjectionBalance() uint64 {
	panic("unimplemented")
}

// ElectraForkEpoch implements chain.Spec.
func (cs *MockChainSpec) ElectraForkEpoch() math.U64 {
	panic("unimplemented")
}

// EpochsPerHistoricalVector implements chain.Spec.
func (cs *MockChainSpec) EpochsPerHistoricalVector() uint64 {
	panic("unimplemented")
}

// EpochsPerSlashingsVector implements chain.Spec.
func (cs *MockChainSpec) EpochsPerSlashingsVector() uint64 {
	panic("unimplemented")
}

// Eth1FollowDistance implements chain.Spec.
func (cs *MockChainSpec) Eth1FollowDistance() uint64 {
	panic("unimplemented")
}

// FieldElementsPerBlob implements chain.Spec.
func (cs *MockChainSpec) FieldElementsPerBlob() uint64 {
	panic("unimplemented")
}

// GetCometBFTConfigForSlot implements chain.Spec.
func (cs *MockChainSpec) GetCometBFTConfigForSlot(_ math.U64) any {
	panic("unimplemented")
}

// HistoricalRootsLimit implements chain.Spec.
func (cs *MockChainSpec) HistoricalRootsLimit() uint64 {
	panic("unimplemented")
}

// InactivityPenaltyQuotient implements chain.Spec.
func (cs *MockChainSpec) InactivityPenaltyQuotient() uint64 {
	panic("unimplemented")
}

// MaxBlobCommitmentsPerBlock implements chain.Spec.
func (cs *MockChainSpec) MaxBlobCommitmentsPerBlock() uint64 {
	panic("unimplemented")
}

// MaxBlobsPerBlock implements chain.Spec.
func (cs *MockChainSpec) MaxBlobsPerBlock() uint64 {
	panic("unimplemented")
}

// MaxDepositsPerBlock implements chain.Spec.
func (cs *MockChainSpec) MaxDepositsPerBlock() uint64 {
	panic("unimplemented")
}

// MaxEffectiveBalance implements chain.Spec.
func (cs *MockChainSpec) MaxEffectiveBalance() uint64 {
	panic("unimplemented")
}

// MaxValidatorsPerWithdrawalsSweep implements chain.Spec.
func (cs *MockChainSpec) MaxValidatorsPerWithdrawalsSweep() uint64 {
	panic("unimplemented")
}

// MaxWithdrawalsPerPayload implements chain.Spec.
func (cs *MockChainSpec) MaxWithdrawalsPerPayload() uint64 {
	panic("unimplemented")
}

// MinDepositAmount implements chain.Spec.
func (cs *MockChainSpec) MinDepositAmount() uint64 {
	panic("unimplemented")
}

// MinEpochsToInactivityPenalty implements chain.Spec.
func (cs *MockChainSpec) MinEpochsToInactivityPenalty() uint64 {
	panic("unimplemented")
}

// ProportionalSlashingMultiplier implements chain.Spec.
func (cs *MockChainSpec) ProportionalSlashingMultiplier() uint64 {
	panic("unimplemented")
}

// SlotToEpoch implements chain.Spec.
func (cs *MockChainSpec) SlotToEpoch(_ math.U64) math.U64 {
	panic("unimplemented")
}

// SlotsPerHistoricalRoot implements chain.Spec.
func (cs *MockChainSpec) SlotsPerHistoricalRoot() uint64 {
	panic("unimplemented")
}

// TargetSecondsPerEth1Block implements chain.Spec.
func (cs *MockChainSpec) TargetSecondsPerEth1Block() uint64 {
	panic("unimplemented")
}

// ValidatorRegistryLimit implements chain.Spec.
func (cs *MockChainSpec) ValidatorRegistryLimit() uint64 {
	panic("unimplemented")
}

// WithinDAPeriod implements chain.Spec.
func (cs *MockChainSpec) WithinDAPeriod(_ math.U64, _ math.U64) bool {
	panic("unimplemented")
}

// TestBuildPruneRangeFn tests the BuildPruneRangeFn function.
func TestBuildPruneRangeFn(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		slotsPerEpoch uint64
		minEpochs     uint64
		eventSlot     math.U64
		expectedStart uint64
		expectedEnd   uint64
	}{
		{
			name:          "Slot greater than window",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(200),
			expectedStart: 0,
			expectedEnd:   40,
		},
		{
			name:          "Slot less than window",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(100),
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Exact boundary case",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(160),
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Greater than boundary just 1 uint case",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(161),
			expectedStart: 0,
			expectedEnd:   1,
		},
		{
			name:          "Zero slot case",
			slotsPerEpoch: 32,
			minEpochs:     5,
			eventSlot:     math.U64(0),
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "SlotsPerEpoch as one",
			slotsPerEpoch: 1,
			minEpochs:     5,
			eventSlot:     math.U64(50),
			expectedStart: 0,
			expectedEnd:   45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &MockChainSpec{
				slotsPerEpoch:                    tt.slotsPerEpoch,
				minEpochsForBlobsSidecarsRequest: tt.minEpochs,
			}

			pruneFn := store.BuildPruneRangeFn[MockBeaconBlock, MockBlockEvent](cs)

			event := MockBlockEvent{
				data: MockBeaconBlock{
					slot: tt.eventSlot,
				},
			}

			start, end := pruneFn(event)
			require.Equal(t, tt.expectedStart, start, "Test case : %s (expectedStart)", tt.name)
			require.Equal(t, tt.expectedEnd, end, "Test case : %s (expectedEnd)", tt.name)
		})
	}
}
