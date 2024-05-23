// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package chain

// Spec defines an interface for accessing chain-specific parameters.
type Spec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
] interface {

	// Gwei value constants.
	//
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
	//
	// SlotsPerEpoch returns the number of slots in an epoch.
	SlotsPerEpoch() uint64
	// SlotsPerHistoricalRoot returns the number of slots per historical root.
	SlotsPerHistoricalRoot() uint64

	// Signature Domains
	//
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
	//
	// DepositContractAddress returns the deposit contract address.
	DepositContractAddress() ExecutionAddressT
	// DepositEth1ChainID returns the chain ID of the deposit contract.
	DepositEth1ChainID() uint64

	// Fork-related values.
	//
	// ElectraForkEpoch returns the epoch at which the Electra fork takes
	// effect.
	ElectraForkEpoch() EpochT

	// State list lengths
	//
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

	// MaxDepositsPerBlock returns the maximum number of deposit operations per
	// block.
	MaxDepositsPerBlock() uint64
	// ProportionalSlashingMultiplier returns the multiplier for calculating
	// slashing penalties.
	ProportionalSlashingMultiplier() uint64

	// Capella Values
	//
	// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
	// payload.
	MaxWithdrawalsPerPayload() uint64
	// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators
	// per withdrawal sweep.
	MaxValidatorsPerWithdrawalsSweep() uint64

	// Deneb Values
	//
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
	//
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
}

// chainSpec is a concrete implementation of the ChainSpec interface, holding
// the actual data.
type chainSpec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
] struct {
	// Data contains the actual chain-specific parameter values.
	Data SpecData[DomainTypeT, EpochT, ExecutionAddressT, SlotT]
}

// NewChainSpec creates a new instance of a ChainSpec with the provided data.
func NewChainSpec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
](data SpecData[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) Spec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
] {
	return &chainSpec[
		DomainTypeT, EpochT, ExecutionAddressT, SlotT,
	]{
		Data: data,
	}
}

// MinDepositAmount returns the minimum deposit amount required.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MinDepositAmount() uint64 {
	return c.Data.MinDepositAmount
}

// MaxEffectiveBalance returns the maximum effective balance.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MaxEffectiveBalance() uint64 {
	return c.Data.MaxEffectiveBalance
}

// EjectionBalance returns the balance below which a validator is ejected.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) EjectionBalance() uint64 {
	return c.Data.EjectionBalance
}

// EffectiveBalanceIncrement returns the increment of effective balance.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) EffectiveBalanceIncrement() uint64 {
	return c.Data.EffectiveBalanceIncrement
}

// SlotsPerEpoch returns the number of slots per epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) SlotsPerEpoch() uint64 {
	return c.Data.SlotsPerEpoch
}

// SlotsPerHistoricalRoot returns the number of slots per historical root.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) SlotsPerHistoricalRoot() uint64 {
	return c.Data.SlotsPerHistoricalRoot
}

// DomainProposer returns the domain for beacon proposer signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeProposer() DomainTypeT {
	return c.Data.DomainTypeProposer
}

// DomainAttester returns the domain for beacon attester signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeAttester() DomainTypeT {
	return c.Data.DomainTypeAttester
}

// DomainRandao returns the domain for RANDAO reveal signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeRandao() DomainTypeT {
	return c.Data.DomainTypeRandao
}

// DomainDeposit returns the domain for deposit contract signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeDeposit() DomainTypeT {
	return c.Data.DomainTypeDeposit
}

// DomainVoluntaryExit returns the domain for voluntary exit signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeVoluntaryExit() DomainTypeT {
	return c.Data.DomainTypeVoluntaryExit
}

// DomainSelectionProof returns the domain for selection proof signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeSelectionProof() DomainTypeT {
	return c.Data.DomainTypeSelectionProof
}

// DomainAggregateAndProof returns the domain for aggregate and proof
// signatures.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeAggregateAndProof() DomainTypeT {
	return c.Data.DomainTypeAggregateAndProof
}

// DomainTypeApplicationMask returns the domain for the application mask.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DomainTypeApplicationMask() DomainTypeT {
	return c.Data.DomainTypeApplicationMask
}

// DepositContractAddress returns the address of the deposit contract.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DepositContractAddress() ExecutionAddressT {
	return c.Data.DepositContractAddress
}

// DepositEth1ChainID returns the chain ID of the execution chain.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) DepositEth1ChainID() uint64 {
	return c.Data.DepositEth1ChainID
}

// ElectraForkEpoch returns the epoch of the Electra fork.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) ElectraForkEpoch() EpochT {
	return c.Data.ElectraForkEpoch
}

// EpochsPerHistoricalVector returns the number of epochs per historical vector.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) EpochsPerHistoricalVector() uint64 {
	return c.Data.EpochsPerHistoricalVector
}

// EpochsPerSlashingsVector returns the number of epochs per slashings vector.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) EpochsPerSlashingsVector() uint64 {
	return c.Data.EpochsPerSlashingsVector
}

// HistoricalRootsLimit returns the limit of historical roots.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) HistoricalRootsLimit() uint64 {
	return c.Data.HistoricalRootsLimit
}

// ValidatorRegistryLimit returns the limit of the validator registry.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) ValidatorRegistryLimit() uint64 {
	return c.Data.ValidatorRegistryLimit
}

// MaxDepositsPerBlock returns the maximum number of deposits per block.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MaxDepositsPerBlock() uint64 {
	return c.Data.MaxDepositsPerBlock
}

// ProportionalSlashingMultiplier returns the proportional slashing multiplier.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) ProportionalSlashingMultiplier() uint64 {
	return c.Data.ProportionalSlashingMultiplier
}

// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
// payload.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MaxWithdrawalsPerPayload() uint64 {
	return c.Data.MaxWithdrawalsPerPayload
}

// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators per
// withdrawals sweep.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MaxValidatorsPerWithdrawalsSweep() uint64 {
	return c.Data.MaxValidatorsPerWithdrawalsSweep
}

// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
// blobs sidecars request.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MinEpochsForBlobsSidecarsRequest() uint64 {
	return c.Data.MinEpochsForBlobsSidecarsRequest
}

// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments per
// block.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MaxBlobCommitmentsPerBlock() uint64 {
	return c.Data.MaxBlobCommitmentsPerBlock
}

// MaxBlobsPerBlock returns the maximum number of blobs per block.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) MaxBlobsPerBlock() uint64 {
	return c.Data.MaxBlobsPerBlock
}

// FieldElementsPerBlob returns the number of field elements per blob.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) FieldElementsPerBlob() uint64 {
	return c.Data.FieldElementsPerBlob
}

// BytesPerBlob returns the number of bytes per blob.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) BytesPerBlob() uint64 {
	return c.Data.BytesPerBlob
}
