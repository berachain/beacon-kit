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

package spec

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Spec defines an interface for accessing chain-specific parameters.
type Chain[CometBFTConfigT any] interface {
	// Gwei value constants.

	// MinDepositAmount returns the minimum amount of Gwei required for a
	// deposit.
	GetMinDepositAmount() uint64

	// MaxEffectiveBalance returns the maximum balance counted in rewards
	// calculations in Gwei.
	GetMaxEffectiveBalance() uint64

	// EjectionBalance returns the balance below which a validator is ejected.
	GetEjectionBalance() uint64

	// EffectiveBalanceIncrement returns the increment of balance used in reward
	// calculations.
	GetEffectiveBalanceIncrement() uint64

	// HysteresisQuotient returns the quotient used in effective balance
	// calculations to create hysteresis. This provides resistance to small
	// balance changes triggering effective balance updates.
	GetHysteresisQuotient() uint64

	// HysteresisDownwardMultiplier returns the multiplier used when checking
	// if the effective balance should be decreased.
	GetHysteresisDownwardMultiplier() uint64

	// HysteresisUpwardMultiplier returns the multiplier used when checking
	// if the effective balance should be increased.
	GetHysteresisUpwardMultiplier() uint64

	// Time parameters constants.

	// SlotsPerEpoch returns the number of slots in an epoch.
	GetSlotsPerEpoch() uint64

	// SlotsPerHistoricalRoot returns the number of slots per historical root.
	GetSlotsPerHistoricalRoot() uint64

	// MinEpochsToInactivityPenalty returns the minimum number of epochs before
	// an inactivity penalty is applied.
	GetMinEpochsToInactivityPenalty() uint64

	// Signature Domains

	// DomainTypeProposer returns the domain for proposer signatures.
	GetDomainTypeProposer() common.DomainType

	// DomainTypeAttester returns the domain for attester signatures.
	GetDomainTypeAttester() common.DomainType

	// DomainTypeRandao returns the domain for RANDAO reveal signatures.
	GetDomainTypeRandao() common.DomainType

	// DomainTypeDeposit returns the domain for deposit signatures.
	GetDomainTypeDeposit() common.DomainType

	// DomainTypeVoluntaryExit returns the domain for voluntary exit signatures.
	GetDomainTypeVoluntaryExit() common.DomainType

	// DomainTypeSelectionProof returns the domain for selection proof
	GetDomainTypeSelectionProof() common.DomainType

	// DomainTypeAggregateAndProof returns the domain for aggregate and proof
	GetDomainTypeAggregateAndProof() common.DomainType

	// DomainTypeApplicationMask returns the domain for application signatures.
	GetDomainTypeApplicationMask() common.DomainType

	// Eth1-related values.

	// DepositContractAddress returns the deposit contract address.
	GetDepositContractAddress() common.ExecutionAddress

	// MaxDepositsPerBlock returns the maximum number of deposit operations per
	// block.
	GetMaxDepositsPerBlock() uint64

	// DepositEth1ChainID returns the chain ID of the deposit contract.
	GetDepositEth1ChainID() uint64

	// Eth1FollowDistance returns the distance between the eth1 chain and the
	// beacon chain for eth1 data.
	GetEth1FollowDistance() uint64

	// TargetSecondsPerEth1Block returns the target time between eth1 blocks.
	GetTargetSecondsPerEth1Block() uint64

	// Fork-related values.
	// DenebPlusForkEpoch returns the epoch at which the Deneb+ fork takes
	GetDenebPlusForkEpoch() math.Epoch
	// ElectraForkEpoch returns the epoch at which the Electra fork takes
	// effect.
	GetElectraForkEpoch() math.Epoch

	// State list lengths

	// EpochsPerHistoricalVector returns the length of the historical vector.
	GetEpochsPerHistoricalVector() uint64

	// EpochsPerSlashingsVector returns the length of the slashing vector.
	GetEpochsPerSlashingsVector() uint64

	// HistoricalRootsLimit returns the maximum number of historical root
	// entries.
	GetHistoricalRootsLimit() uint64

	// ValidatorRegistryLimit returns the maximum number of validators in the
	// registry.
	GetValidatorRegistryLimit() uint64

	// Rewards and Penalties

	// InactivityPenaltyQuotient returns the inactivity penalty quotient.
	GetInactivityPenaltyQuotient() uint64

	// ProportionalSlashingMultiplier returns the multiplier for calculating
	// slashing penalties.
	GetProportionalSlashingMultiplier() uint64

	// Capella Values

	// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
	// payload.
	GetMaxWithdrawalsPerPayload() uint64

	// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators
	// per withdrawal sweep.
	GetMaxValidatorsPerWithdrawalsSweep() uint64

	// Deneb Values

	// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
	// blob sidecar requests.
	GetMinEpochsForBlobsSidecarsRequest() uint64

	// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments
	// per block.
	GetMaxBlobCommitmentsPerBlock() uint64

	// MaxBlobsPerBlock returns the maximum number of blobs per block.
	GetMaxBlobsPerBlock() uint64

	// FieldElementsPerBlob returns the number of field elements per blob.
	GetFieldElementsPerBlob() uint64

	// BytesPerBlob returns the number of bytes per blob.
	GetBytesPerBlob() uint64

	// Helpers for ChainSpecData

	// ActiveForkVersionForSlot returns the active fork version for a given
	// slot.
	GetActiveForkVersionForSlot(slot math.Slot) uint32

	// ActiveForkVersionForEpoch returns the active fork version for a given
	// epoch.
	GetActiveForkVersionForEpoch(epoch math.Epoch) uint32

	// GetSlotToEpoch converts a slot number to an epoch number.
	GetSlotToEpoch(slot math.Slot) math.Epoch

	// WithinDAPeriod checks if a given block slot is within the data
	// availability period relative to the current slot.
	GetWithinDAPeriod(block, current math.Slot) bool

	// GetCometBFTConfigForSlot retrieves the CometBFT config for a specific
	// slot.
	GetCometBFTConfigForSlot(slot math.Slot) CometBFTConfigT

	// Berachain Values

	// EVMInflationAddress returns the address on the EVM which will receive
	// the inflation amount of native EVM balance through a withdrawal every
	// block.
	GetEVMInflationAddress() common.ExecutionAddress

	// EVMInflationPerBlock returns the amount of native EVM balance (in Gwei)
	// to be minted to the EVMInflationAddress via a withdrawal every block.
	GetEVMInflationPerBlock() uint64
}
