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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package chain

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

type BalancesSpec interface {
	// MaxEffectiveBalance returns the maximum balance counted in rewards calculations in Gwei.
	MaxEffectiveBalance() uint64

	// EjectionBalance returns the balance below which a validator is ejected.
	EjectionBalance() uint64

	// EffectiveBalanceIncrement returns the increment of balance used in reward
	// calculations.
	EffectiveBalanceIncrement() uint64
}

type HysteresisSpec interface {
	// HysteresisQuotient returns the quotient used in effective balance
	// calculations to create hysteresis. This provides resistance to small
	// balance changes triggering effective balance updates.
	HysteresisQuotient() uint64

	// HysteresisDownwardMultiplier returns the multiplier used when checking
	// if the effective balance should be decreased.
	HysteresisDownwardMultiplier() uint64

	// HysteresisUpwardMultiplier returns the multiplier used when checking
	// if the effective balance should be increased.
	HysteresisUpwardMultiplier() uint64
}

type DepositSpec interface {
	// MaxDepositsPerBlock returns the maximum number of deposit operations per
	// block.
	MaxDepositsPerBlock() uint64

	// DepositEth1ChainID returns the chain ID of the deposit contract.
	DepositEth1ChainID() uint64
}

type DomainTypeSpec interface {
	// Signature Domains

	// DomainTypeProposer returns the domain for proposer signatures.
	DomainTypeProposer() common.DomainType

	// DomainTypeAttester returns the domain for attester signatures.
	DomainTypeAttester() common.DomainType

	// DomainTypeRandao returns the domain for RANDAO reveal signatures.
	DomainTypeRandao() common.DomainType

	// DomainTypeDeposit returns the domain for deposit signatures.
	DomainTypeDeposit() common.DomainType

	// DomainTypeVoluntaryExit returns the domain for voluntary exit signatures.
	DomainTypeVoluntaryExit() common.DomainType

	// DomainTypeSelectionProof returns the domain for selection proof
	DomainTypeSelectionProof() common.DomainType

	// DomainTypeAggregateAndProof returns the domain for aggregate and proof
	DomainTypeAggregateAndProof() common.DomainType

	// DomainTypeApplicationMask returns the domain for application signatures.
	DomainTypeApplicationMask() common.DomainType
}

// Fork-related values.
type ForkSpec interface {
	// GenesisTime returns the time at which the genesis block was created.
	GenesisTime() uint64

	// Deneb1ForkTime returns the time at which the Deneb1 fork takes effect.
	Deneb1ForkTime() uint64

	// ElectraForkTime returns the time at which the Electra fork takes effect.
	ElectraForkTime() uint64
}

type BlobSpec interface {
	// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments
	// per block.
	MaxBlobCommitmentsPerBlock() uint64

	// MaxBlobsPerBlock returns the maximum number of blobs per block.
	MaxBlobsPerBlock() uint64

	// FieldElementsPerBlob returns the number of field elements per blob.
	FieldElementsPerBlob() uint64

	// WithinDAPeriod checks if a given block slot is within the data
	// availability period relative to the current slot.
	WithinDAPeriod(block, current math.Slot) bool

	// BytesPerBlob returns the number of bytes per blob.
	BytesPerBlob() uint64

	// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
	// blob sidecar requests.
	MinEpochsForBlobsSidecarsRequest() math.Epoch
}

// Helpers for Fork Version
type ForkVersionSpec interface {
	// GenesisForkVersion returns the fork version at genesis.
	GenesisForkVersion() common.Version

	// ActiveForkVersionForTimestamp returns the active fork version for a given timestamp.
	ActiveForkVersionForTimestamp(timestamp math.U64) common.Version
}

type EVMInflationSpec interface {
	// EVMInflationAddress returns the address on the EVM which will receive
	// the inflation amount of native EVM balance through a withdrawal every
	// block.
	EVMInflationAddress(timestamp math.U64) common.ExecutionAddress

	// EVMInflationPerBlock returns the amount of native EVM balance (in Gwei)
	// to be minted to the EVMInflationAddress via a withdrawal every block.
	EVMInflationPerBlock(timestamp math.U64) uint64
}

type WithdrawalsSpec interface {
	// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
	// payload.
	MaxWithdrawalsPerPayload() uint64

	// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators
	// per withdrawal sweep.
	MaxValidatorsPerWithdrawalsSweep() uint64
}

// Spec defines an interface for accessing chain-specific parameters.
type Spec interface {
	DepositSpec
	BalancesSpec
	HysteresisSpec
	DomainTypeSpec
	ForkSpec
	BlobSpec
	ForkVersionSpec
	EVMInflationSpec
	WithdrawalsSpec

	// Time parameters constants.

	// SlotsPerEpoch returns the number of slots in an epoch.
	SlotsPerEpoch() uint64

	// SlotsPerHistoricalRoot returns the number of slots per historical root.
	SlotsPerHistoricalRoot() uint64

	// MinEpochsToInactivityPenalty returns the minimum number of epochs before
	// an inactivity penalty is applied.
	MinEpochsToInactivityPenalty() uint64

	// Eth1-related values.

	// DepositContractAddress returns the deposit contract address.
	DepositContractAddress() common.ExecutionAddress

	// Eth1FollowDistance returns the distance between the eth1 chain and the
	// beacon chain for eth1 data.
	Eth1FollowDistance() uint64

	// TargetSecondsPerEth1Block returns the target time between eth1 blocks.
	TargetSecondsPerEth1Block() uint64

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

	// SlotToEpoch converts a slot number to an epoch number.
	SlotToEpoch(slot math.Slot) math.Epoch

	// Berachain Values

	// ValidatorSetCap retrieves the maximum number of validators allowed in the active set.
	ValidatorSetCap() uint64
}

// spec is a concrete implementation of the Spec interface, holding the actual data.
type spec struct {
	// Data contains the actual chain-specific parameter values.
	Data *SpecData
}

// NewSpec creates a new instance of a Spec with the provided data.
func NewSpec(data *SpecData) (Spec, error) {
	s := spec{Data: data}
	return s, s.validate()
}

// validate ensures that the chain spec is valid, returning error if it is not.
func (s spec) validate() error {
	if s.Data.MaxWithdrawalsPerPayload <= 1 {
		return ErrInsufficientMaxWithdrawalsPerPayload
	}

	if s.Data.ValidatorSetCap > s.Data.ValidatorRegistryLimit {
		return ErrInvalidValidatorSetCap
	}

	// EVM Inflation values can be zero or non-zero, no validation needed.

	// TODO: Add more validation rules here.
	return nil
}

// MaxEffectiveBalance returns the maximum effective balance.
func (s spec) MaxEffectiveBalance() uint64 {
	return s.Data.MaxEffectiveBalance
}

// EjectionBalance returns the balance below which a validator is ejected.
func (s spec) EjectionBalance() uint64 {
	return s.Data.EjectionBalance
}

// EffectiveBalanceIncrement returns the increment of effective balance.
func (s spec) EffectiveBalanceIncrement() uint64 {
	return s.Data.EffectiveBalanceIncrement
}

func (s spec) HysteresisQuotient() uint64 {
	return s.Data.HysteresisQuotient
}

func (s spec) HysteresisDownwardMultiplier() uint64 {
	return s.Data.HysteresisDownwardMultiplier
}

func (s spec) HysteresisUpwardMultiplier() uint64 {
	return s.Data.HysteresisUpwardMultiplier
}

// SlotsPerEpoch returns the number of slots per epoch.
func (s spec) SlotsPerEpoch() uint64 {
	return s.Data.SlotsPerEpoch
}

// SlotsPerHistoricalRoot returns the number of slots per historical root.
func (s spec) SlotsPerHistoricalRoot() uint64 {
	return s.Data.SlotsPerHistoricalRoot
}

// MinEpochsToInactivityPenalty returns the minimum number of epochs before an
// inactivity penalty is applied.
func (s spec) MinEpochsToInactivityPenalty() uint64 {
	return s.Data.MinEpochsToInactivityPenalty
}

// DomainTypeProposer returns the domain for beacon proposer signatures.
func (s spec) DomainTypeProposer() common.DomainType {
	return s.Data.DomainTypeProposer
}

// DomainTypeAttester returns the domain for beacon attester signatures.
func (s spec) DomainTypeAttester() common.DomainType {
	return s.Data.DomainTypeAttester
}

// DomainTypeRandao returns the domain for RANDAO reveal signatures.
func (s spec) DomainTypeRandao() common.DomainType {
	return s.Data.DomainTypeRandao
}

// DomainTypeDeposit returns the domain for deposit contract signatures.
func (s spec) DomainTypeDeposit() common.DomainType {
	return s.Data.DomainTypeDeposit
}

// DomainTypeVoluntaryExit returns the domain for voluntary exit signatures.
func (s spec) DomainTypeVoluntaryExit() common.DomainType {
	return s.Data.DomainTypeVoluntaryExit
}

// DomainTypeSelectionProof returns the domain for selection proof signatures.
func (s spec) DomainTypeSelectionProof() common.DomainType {
	return s.Data.DomainTypeSelectionProof
}

// DomainTypeAggregateAndProof returns the domain for aggregate and proof
// signatures.
func (s spec) DomainTypeAggregateAndProof() common.DomainType {
	return s.Data.DomainTypeAggregateAndProof
}

// DomainTypeApplicationMask returns the domain for the application mask.
func (s spec) DomainTypeApplicationMask() common.DomainType {
	return s.Data.DomainTypeApplicationMask
}

// DepositContractAddress returns the address of the deposit contract.
func (s spec) DepositContractAddress() common.ExecutionAddress {
	return s.Data.DepositContractAddress
}

// MaxDepositsPerBlock returns the maximum number of deposits per block.
func (s spec) MaxDepositsPerBlock() uint64 {
	return s.Data.MaxDepositsPerBlock
}

// DepositEth1ChainID returns the chain ID of the execution chain.
func (s spec) DepositEth1ChainID() uint64 {
	return s.Data.DepositEth1ChainID
}

// Eth1FollowDistance returns the distance between the eth1 chain and the beacon
// chain.
func (s spec) Eth1FollowDistance() uint64 {
	return s.Data.Eth1FollowDistance
}

// TargetSecondsPerEth1Block returns the target time between eth1 blocks.
func (s spec) TargetSecondsPerEth1Block() uint64 {
	return s.Data.TargetSecondsPerEth1Block
}

// GenesisTime returns the time at which the genesis block was created.
func (s spec) GenesisTime() uint64 {
	return s.Data.GenesisTime
}

// Deneb1ForkTime returns the epoch of the Deneb1 fork.
func (s spec) Deneb1ForkTime() uint64 {
	return s.Data.Deneb1ForkTime
}

// ElectraForkTime returns the epoch of the Electra fork.
func (s spec) ElectraForkTime() uint64 {
	return s.Data.ElectraForkTime
}

// EpochsPerHistoricalVector returns the number of epochs per historical vector.
func (s spec) EpochsPerHistoricalVector() uint64 {
	return s.Data.EpochsPerHistoricalVector
}

// EpochsPerSlashingsVector returns the number of epochs per slashings vector.
func (s spec) EpochsPerSlashingsVector() uint64 {
	return s.Data.EpochsPerSlashingsVector
}

// HistoricalRootsLimit returns the limit of historical roots.
func (s spec) HistoricalRootsLimit() uint64 {
	return s.Data.HistoricalRootsLimit
}

// ValidatorRegistryLimit returns the limit of the validator registry.
func (s spec) ValidatorRegistryLimit() uint64 {
	return s.Data.ValidatorRegistryLimit
}

// InactivityPenaltyQuotient returns the inactivity penalty quotient.
func (s spec) InactivityPenaltyQuotient() uint64 {
	return s.Data.InactivityPenaltyQuotient
}

// ProportionalSlashingMultiplier returns the proportional slashing multiplier.
func (s spec) ProportionalSlashingMultiplier() uint64 {
	return s.Data.ProportionalSlashingMultiplier
}

// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
// payload.
func (s spec) MaxWithdrawalsPerPayload() uint64 {
	return s.Data.MaxWithdrawalsPerPayload
}

// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators per withdrawals sweep.
func (s spec) MaxValidatorsPerWithdrawalsSweep() uint64 {
	return s.Data.MaxValidatorsPerWithdrawalsSweep
}

// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
// blobs sidecars request.
func (s spec) MinEpochsForBlobsSidecarsRequest() math.Epoch {
	return math.Epoch(s.Data.MinEpochsForBlobsSidecarsRequest)
}

// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments per
// block.
func (s spec) MaxBlobCommitmentsPerBlock() uint64 {
	return s.Data.MaxBlobCommitmentsPerBlock
}

// MaxBlobsPerBlock returns the maximum number of blobs per block.
func (s spec) MaxBlobsPerBlock() uint64 {
	return s.Data.MaxBlobsPerBlock
}

// FieldElementsPerBlob returns the number of field elements per blob.
func (s spec) FieldElementsPerBlob() uint64 {
	return s.Data.FieldElementsPerBlob
}

// BytesPerBlob returns the number of bytes per blob.
func (s spec) BytesPerBlob() uint64 {
	return s.Data.BytesPerBlob
}

// ValidatorSetCap retrieves the maximum number of validators allowed in the active set.
func (s spec) ValidatorSetCap() uint64 {
	return s.Data.ValidatorSetCap
}

// EVMInflationAddress returns the address on the EVM which will receive the
// inflation amount of native EVM balance through a withdrawal every block.
func (s spec) EVMInflationAddress(timestamp math.U64) common.ExecutionAddress {
	fv := s.ActiveForkVersionForTimestamp(timestamp)
	switch fv {
	case version.Deneb1(), version.Electra():
		return s.Data.EVMInflationAddressDeneb1
	case version.Deneb():
		return s.Data.EVMInflationAddressGenesis
	default:
		panic(fmt.Sprintf("EVMInflationAddress not supported for this fork version: %d", fv))
	}
}

// EVMInflationPerBlock returns the amount of native EVM balance (in Gwei) to
// be minted to the EVMInflationAddress via a withdrawal every block.
func (s spec) EVMInflationPerBlock(timestamp math.U64) uint64 {
	fv := s.ActiveForkVersionForTimestamp(timestamp)
	switch fv {
	case version.Deneb1(), version.Electra():
		return s.Data.EVMInflationPerBlockDeneb1
	case version.Deneb():
		return s.Data.EVMInflationPerBlockGenesis
	default:
		panic(fmt.Sprintf("EVMInflationPerBlock not supported for this fork version: %d", fv))
	}
}
