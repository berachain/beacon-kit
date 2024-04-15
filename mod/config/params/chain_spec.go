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

package params

import (
	"github.com/berachain/beacon-kit/mod/primitives"
)

var _ primitives.ChainSpec = (*chainSpec)(nil)

// chainSpec is a concrete implementation of the ChainSpec interface, holding
// the actual data.
type chainSpec struct {
	// Data contains the actual chain-specific parameter values.
	Data *ChainSpecData
}

// NewChainSpec creates a new instance of a ChainSpec with the provided data.
func NewChainSpec(data *ChainSpecData) primitives.ChainSpec {
	return &chainSpec{
		Data: data,
	}
}

// MinDepositAmount returns the minimum deposit amount required.
func (c *chainSpec) MinDepositAmount() uint64 {
	return c.Data.MinDepositAmount
}

// MaxEffectiveBalance returns the maximum effective balance.
func (c *chainSpec) MaxEffectiveBalance() uint64 {
	return c.Data.MaxEffectiveBalance
}

// EjectionBalance returns the balance below which a validator is ejected.
func (c *chainSpec) EjectionBalance() uint64 {
	return c.Data.EjectionBalance
}

// EffectiveBalanceIncrement returns the increment of effective balance.
func (c *chainSpec) EffectiveBalanceIncrement() uint64 {
	return c.Data.EffectiveBalanceIncrement
}

// SlotsPerEpoch returns the number of slots per epoch.
func (c *chainSpec) SlotsPerEpoch() uint64 {
	return c.Data.SlotsPerEpoch
}

// SlotsPerHistoricalRoot returns the number of slots per historical root.
func (c *chainSpec) SlotsPerHistoricalRoot() uint64 {
	return c.Data.SlotsPerHistoricalRoot
}

// DomainProposer returns the domain for beacon proposer signatures.
func (c *chainSpec) DomainTypeProposer() primitives.DomainType {
	return c.Data.DomainTypeProposer
}

// DomainAttester returns the domain for beacon attester signatures.
func (c *chainSpec) DomainTypeAttester() primitives.DomainType {
	return c.Data.DomainTypeAttester
}

// DomainRandao returns the domain for RANDAO reveal signatures.
func (c *chainSpec) DomainTypeRandao() primitives.DomainType {
	return c.Data.DomainTypeRandao
}

// DomainDeposit returns the domain for deposit contract signatures.
func (c *chainSpec) DomainTypeDeposit() primitives.DomainType {
	return c.Data.DomainTypeDeposit
}

// DomainVoluntaryExit returns the domain for voluntary exit signatures.
func (c *chainSpec) DomainTypeVoluntaryExit() primitives.DomainType {
	return c.Data.DomainTypeVoluntaryExit
}

// DomainSelectionProof returns the domain for selection proof signatures.
func (c *chainSpec) DomainTypeSelectionProof() primitives.DomainType {
	return c.Data.DomainTypeSelectionProof
}

// DomainAggregateAndProof returns the domain for aggregate and proof
// signatures.
func (c *chainSpec) DomainTypeAggregateAndProof() primitives.DomainType {
	return c.Data.DomainTypeAggregateAndProof
}

// DomainTypeApplicationMask returns the domain for the application mask.
func (c *chainSpec) DomainTypeApplicationMask() primitives.DomainType {
	return c.Data.DomainTypeApplicationMask
}

// DepositContractAddress returns the address of the deposit contract.
func (c *chainSpec) DepositContractAddress() primitives.ExecutionAddress {
	return c.Data.DepositContractAddress
}

// ElectraForkEpoch returns the epoch of the Electra fork.
func (c *chainSpec) ElectraForkEpoch() primitives.Epoch {
	return c.Data.ElectraForkEpoch
}

// EpochsPerHistoricalVector returns the number of epochs per historical vector.
func (c *chainSpec) EpochsPerHistoricalVector() uint64 {
	return c.Data.EpochsPerHistoricalVector
}

// EpochsPerSlashingsVector returns the number of epochs per slashings vector.
func (c *chainSpec) EpochsPerSlashingsVector() uint64 {
	return c.Data.EpochsPerSlashingsVector
}

// HistoricalRootsLimit returns the limit of historical roots.
func (c *chainSpec) HistoricalRootsLimit() uint64 {
	return c.Data.HistoricalRootsLimit
}

// ValidatorRegistryLimit returns the limit of the validator registry.
func (c *chainSpec) ValidatorRegistryLimit() uint64 {
	return c.Data.ValidatorRegistryLimit
}

// MaxDepositsPerBlock returns the maximum number of deposits per block.
func (c *chainSpec) MaxDepositsPerBlock() uint64 {
	return c.Data.MaxDepositsPerBlock
}

// ProportionalSlashingMultiplier returns the proportional slashing multiplier.
func (c *chainSpec) ProportionalSlashingMultiplier() uint64 {
	return c.Data.ProportionalSlashingMultiplier
}

// MaxWithdrawalsPerPayload returns the maximum number of withdrawals per
// payload.
func (c *chainSpec) MaxWithdrawalsPerPayload() uint64 {
	return c.Data.MaxWithdrawalsPerPayload
}

// MaxValidatorsPerWithdrawalsSweep returns the maximum number of validators per
// withdrawals sweep.
func (c *chainSpec) MaxValidatorsPerWithdrawalsSweep() uint64 {
	return c.Data.MaxValidatorsPerWithdrawalsSweep
}

// MinEpochsForBlobsSidecarsRequest returns the minimum number of epochs for
// blobs sidecars request.
func (c *chainSpec) MinEpochsForBlobsSidecarsRequest() uint64 {
	return c.Data.MinEpochsForBlobsSidecarsRequest
}

// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments per
// block.
func (c *chainSpec) MaxBlobCommitmentsPerBlock() uint64 {
	return c.Data.MaxBlobCommitmentsPerBlock
}

// MaxBlobsPerBlock returns the maximum number of blobs per block.
func (c *chainSpec) MaxBlobsPerBlock() uint64 {
	return c.Data.MaxBlobsPerBlock
}

// FieldElementsPerBlob returns the number of field elements per blob.
func (c *chainSpec) FieldElementsPerBlob() uint64 {
	return c.Data.FieldElementsPerBlob
}

// BytesPerBlob returns the number of bytes per blob.
func (c *chainSpec) BytesPerBlob() uint64 {
	return c.Data.BytesPerBlob
}
