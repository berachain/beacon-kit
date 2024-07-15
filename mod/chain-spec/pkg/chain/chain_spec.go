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

// Spec defines an interface for accessing chain-specific parameters.
type Spec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
	CometBFTConfigT any,
] interface {

	// Gwei value constants.

	// MinDepositAmount returns the minimum amount of Gwei required for a
	// deposit.
	MinDepositAmount() uint64

	// MaxEffectiveBalance returns the maximum balance counted in rewards
	// calculations in Gwei.
	MaxEffectiveBalance() uint64

	// EjectionBalance returns the balance below which a validator is ejected.
	EjectionBalance() uint64

	// EffectiveBalanceIncrement returns the increment of balance used in reward
	// calculations.
	EffectiveBalanceIncrement() uint64

	// Time parameters constants.

	// SlotsPerEpoch returns the number of slots in an epoch.
	SlotsPerEpoch() uint64

	// SlotsPerHistoricalRoot returns the number of slots per historical root.
	SlotsPerHistoricalRoot() uint64

	// MinEpochsToInactivityPenalty returns the minimum number of epochs before
	// an inactivity penalty is applied.
	MinEpochsToInactivityPenalty() uint64

	// Signature Domains

	// DomainTypeProposer returns the domain for proposer signatures.
	DomainTypeProposer() DomainTypeT

	// DomainTypeAttester returns the domain for attester signatures.
	DomainTypeAttester() DomainTypeT

	// DomainTypeRandao returns the domain for RANDAO reveal signatures.
	DomainTypeRandao() DomainTypeT

	// DomainTypeDeposit returns the domain for deposit signatures.
	DomainTypeDeposit() DomainTypeT

	// DomainTypeVoluntaryExit returns the domain for voluntary exit signatures.
	DomainTypeVoluntaryExit() DomainTypeT

	// DomainTypeSelectionProof returns the domain for selection proof
	DomainTypeSelectionProof() DomainTypeT

	// DomainTypeAggregateAndProof returns the domain for aggregate and proof
	DomainTypeAggregateAndProof() DomainTypeT

	// DomainTypeApplicationMask returns the domain for application signatures.
	DomainTypeApplicationMask() DomainTypeT

	// Eth1-related values.

	// DepositContractAddress returns the deposit contract address.
	DepositContractAddress() ExecutionAddressT

	// MaxDepositsPerBlock returns the maximum number of deposit operations per
	// block.
	MaxDepositsPerBlock() uint64

	// DepositEth1ChainID returns the chain ID of the deposit contract.
	DepositEth1ChainID() uint64

	// Eth1FollowDistance returns the distance between the eth1 chain and the
	// beacon chain for eth1 data.
	Eth1FollowDistance() uint64

	// TargetSecondsPerEth1Block returns the target time between eth1 blocks.
	TargetSecondsPerEth1Block() uint64

	// Fork-related values.

	// ElectraForkEpoch returns the epoch at which the Electra fork takes
	// effect.
	ElectraForkEpoch() EpochT

	// State list lengths

	// EpochsPerHistoricalVector returns the length of the historical vector.
	EpochsPerHistoricalVector() uint64

	// EpochsPerSlashingsVector returns the length of the slashing vector.
	EpochsPerSlashingsVector() uint64

	// HistoricalRootsLimit returns the maximum number of historical root
	// entries.
	HistoricalRootsLimit() uint64

	// ValidatorRegistryLimit returns the maximum number of validators in the
	// registry.
	ValidatorRegistryLimit() uint64

	// Rewards and Penalties

	// InactivityPenaltyQuotient returns the inactivity penalty quotient.
	InactivityPenaltyQuotient() uint64

	// ProportionalSlashingMultiplier returns the multiplier for calculating
	// slashing penalties.
	ProportionalSlashingMultiplier() uint64

	// Capella Values

	// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
	// payload.
	MaxWithdrawalsPerPayload() uint64

	// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators
	// per withdrawal sweep.
	MaxValidatorsPerWithdrawalsSweep() uint64

	// Deneb Values

	// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
	// blob sidecar requests.
	MinEpochsForBlobsSidecarsRequest() uint64

	// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments
	// per block.
	MaxBlobCommitmentsPerBlock() uint64

	// MaxBlobsPerBlock returns the maximum number of blobs per block.
	MaxBlobsPerBlock() uint64

	// FieldElementsPerBlob returns the number of field elements per blob.
	FieldElementsPerBlob() uint64

	// BytesPerBlob returns the number of bytes per blob.
	BytesPerBlob() uint64

	// Helpers for ChainSpecData

	// ActiveForkVersionForSlot returns the active fork version for a given
	// slot.
	ActiveForkVersionForSlot(slot SlotT) uint32

	// ActiveForkVersionForEpoch returns the active fork version for a given
	// epoch.
	ActiveForkVersionForEpoch(epoch EpochT) uint32

	// SlotToEpoch converts a slot number to an epoch number.
	SlotToEpoch(slot SlotT) EpochT

	// WithinDAPeriod checks if a given block slot is within the data
	// availability period relative to the current slot.
	WithinDAPeriod(block, current SlotT) bool

	// GetCometBFTConfigForSlot retrieves the CometBFT config for a specific
	// slot.
	GetCometBFTConfigForSlot(slot SlotT) CometBFTConfigT
}

// chainSpec is a concrete implementation of the ChainSpec interface, holding
// the actual data.
type chainSpec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
	CometBFTConfigT any,
] struct {
	// Data contains the actual chain-specific parameter values.
	Data SpecData[DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT]
}

// NewChainSpec creates a new instance of a ChainSpec with the provided data.
func NewChainSpec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
	CometBFTConfigT any,
](data SpecData[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) Spec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
] {
	return &chainSpec[
		DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
	]{
		Data: data,
	}
}

// MinDepositAmount returns the minimum deposit amount required.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MinDepositAmount() uint64 {
	return c.Data.MinDepositAmount
}

// MaxEffectiveBalance returns the maximum effective balance.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MaxEffectiveBalance() uint64 {
	return c.Data.MaxEffectiveBalance
}

// EjectionBalance returns the balance below which a validator is ejected.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) EjectionBalance() uint64 {
	return c.Data.EjectionBalance
}

// EffectiveBalanceIncrement returns the increment of effective balance.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) EffectiveBalanceIncrement() uint64 {
	return c.Data.EffectiveBalanceIncrement
}

// SlotsPerEpoch returns the number of slots per epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) SlotsPerEpoch() uint64 {
	return c.Data.SlotsPerEpoch
}

// SlotsPerHistoricalRoot returns the number of slots per historical root.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) SlotsPerHistoricalRoot() uint64 {
	return c.Data.SlotsPerHistoricalRoot
}

// MinEpochsToInactivityPenalty returns the minimum number of epochs before an
// inactivity penalty is applied.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MinEpochsToInactivityPenalty() uint64 {
	return c.Data.MinEpochsToInactivityPenalty
}

// DomainTypeProposer returns the domain for beacon proposer signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeProposer() DomainTypeT {
	return c.Data.DomainTypeProposer
}

// DomainTypeAttester returns the domain for beacon attester signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeAttester() DomainTypeT {
	return c.Data.DomainTypeAttester
}

// DomainTypeRandao returns the domain for RANDAO reveal signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeRandao() DomainTypeT {
	return c.Data.DomainTypeRandao
}

// DomainTypeDeposit returns the domain for deposit contract signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeDeposit() DomainTypeT {
	return c.Data.DomainTypeDeposit
}

// DomainTypeVoluntaryExit returns the domain for voluntary exit signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeVoluntaryExit() DomainTypeT {
	return c.Data.DomainTypeVoluntaryExit
}

// DomainTypeSelectionProof returns the domain for selection proof signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeSelectionProof() DomainTypeT {
	return c.Data.DomainTypeSelectionProof
}

// DomainTypeAggregateAndProof returns the domain for aggregate and proof
// signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeAggregateAndProof() DomainTypeT {
	return c.Data.DomainTypeAggregateAndProof
}

// DomainTypeApplicationMask returns the domain for the application mask.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DomainTypeApplicationMask() DomainTypeT {
	return c.Data.DomainTypeApplicationMask
}

// DepositContractAddress returns the address of the deposit contract.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DepositContractAddress() ExecutionAddressT {
	return c.Data.DepositContractAddress
}

// MaxDepositsPerBlock returns the maximum number of deposits per block.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MaxDepositsPerBlock() uint64 {
	return c.Data.MaxDepositsPerBlock
}

// DepositEth1ChainID returns the chain ID of the execution chain.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) DepositEth1ChainID() uint64 {
	return c.Data.DepositEth1ChainID
}

// Eth1FollowDistance returns the distance between the eth1 chain and the beacon
// chain.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) Eth1FollowDistance() uint64 {
	return c.Data.Eth1FollowDistance
}

// TargetSecondsPerEth1Block returns the target time between eth1 blocks.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) TargetSecondsPerEth1Block() uint64 {
	return c.Data.TargetSecondsPerEth1Block
}

// ElectraForkEpoch returns the epoch of the Electra fork.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) ElectraForkEpoch() EpochT {
	return c.Data.ElectraForkEpoch
}

// EpochsPerHistoricalVector returns the number of epochs per historical vector.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) EpochsPerHistoricalVector() uint64 {
	return c.Data.EpochsPerHistoricalVector
}

// EpochsPerSlashingsVector returns the number of epochs per slashings vector.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) EpochsPerSlashingsVector() uint64 {
	return c.Data.EpochsPerSlashingsVector
}

// HistoricalRootsLimit returns the limit of historical roots.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) HistoricalRootsLimit() uint64 {
	return c.Data.HistoricalRootsLimit
}

// ValidatorRegistryLimit returns the limit of the validator registry.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) ValidatorRegistryLimit() uint64 {
	return c.Data.ValidatorRegistryLimit
}

// InactivityPenaltyQuotient returns the inactivity penalty quotient.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) InactivityPenaltyQuotient() uint64 {
	return c.Data.InactivityPenaltyQuotient
}

// ProportionalSlashingMultiplier returns the proportional slashing multiplier.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) ProportionalSlashingMultiplier() uint64 {
	return c.Data.ProportionalSlashingMultiplier
}

// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
// payload.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MaxWithdrawalsPerPayload() uint64 {
	return c.Data.MaxWithdrawalsPerPayload
}

// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators per
// withdrawals sweep.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MaxValidatorsPerWithdrawalsSweep() uint64 {
	return c.Data.MaxValidatorsPerWithdrawalsSweep
}

// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
// blobs sidecars request.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MinEpochsForBlobsSidecarsRequest() uint64 {
	return c.Data.MinEpochsForBlobsSidecarsRequest
}

// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments per
// block.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MaxBlobCommitmentsPerBlock() uint64 {
	return c.Data.MaxBlobCommitmentsPerBlock
}

// MaxBlobsPerBlock returns the maximum number of blobs per block.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) MaxBlobsPerBlock() uint64 {
	return c.Data.MaxBlobsPerBlock
}

// FieldElementsPerBlob returns the number of field elements per blob.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) FieldElementsPerBlob() uint64 {
	return c.Data.FieldElementsPerBlob
}

// BytesPerBlob returns the number of bytes per blob.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) BytesPerBlob() uint64 {
	return c.Data.BytesPerBlob
}

// GetCometBFTConfigForSlot returns the CometBFT configuration for the given
// slot.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) GetCometBFTConfigForSlot(_ SlotT) CometBFTConfigT {
	return c.Data.CometValues
}
