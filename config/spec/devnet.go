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

package spec

import (
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
)

const (
	// DevnetEVMInflationAddress is the address of the EVM inflation contract.
	DevnetEVMInflationAddress = "0x6942069420694206942069420694206942069420"

	// DevnetEVMInflationPerBlock is the amount of native EVM balance (in units
	// of Gwei) to be minted per EL block.
	DevnetEVMInflationPerBlock = 10e9

	// DevnetMaxStakeAmount is the maximum amount of native EVM balance (in units
	// of Gwei) that can be staked.
	DevnetMaxStakeAmount = 4000e9
)

// DevnetChainSpec is the chain.Spec for a devnet. Used by `make start` and unit tests.
func DevnetChainSpec() (chain.Spec, error) {
	devnetSpecData := &chain.SpecData{
		// Gwei values constants.
		MaxEffectiveBalance:       DevnetMaxStakeAmount,
		EjectionBalance:           DefaultEjectionBalance,
		EffectiveBalanceIncrement: DefaultEffectiveBalanceIncrement,

		HysteresisQuotient:           DefaultHysteresisQuotient,
		HysteresisDownwardMultiplier: DefaultHysteresisDownwardMultiplier,
		HysteresisUpwardMultiplier:   DefaultHysteresisUpwardMultiplier,

		// Time parameters constants.
		SlotsPerEpoch:                DefaultSlotsPerEpoch,
		SlotsPerHistoricalRoot:       DefaultSlotsPerHistoricalRoot,
		MinEpochsToInactivityPenalty: DefaultMinEpochsToInactivityPenalty,

		// Signature domains.
		DomainTypeProposer:          bytes.FromUint32(DefaultDomainTypeProposer),
		DomainTypeAttester:          bytes.FromUint32(DefaultDomainTypeAttester),
		DomainTypeRandao:            bytes.FromUint32(DefaultDomainTypeRandao),
		DomainTypeDeposit:           bytes.FromUint32(DefaultDomainTypeDeposit),
		DomainTypeVoluntaryExit:     bytes.FromUint32(DefaultDomainTypeVoluntaryExit),
		DomainTypeSelectionProof:    bytes.FromUint32(DefaultDomainTypeSelectionProof),
		DomainTypeAggregateAndProof: bytes.FromUint32(DefaultDomainTypeAggregateAndProof),
		DomainTypeApplicationMask:   bytes.FromUint32(DefaultDomainTypeApplicationMask),

		// Eth1-related values.
		DepositContractAddress:    common.NewExecutionAddressFromHex(DefaultDepositContractAddress),
		MaxDepositsPerBlock:       DefaultMaxDepositsPerBlock,
		DepositEth1ChainID:        DevnetEth1ChainID,
		Eth1FollowDistance:        DefaultEth1FollowDistance,
		TargetSecondsPerEth1Block: DefaultTargetSecondsPerEth1Block,

		// Fork-related values.
		Deneb1ForkEpoch:  DefaultDeneb1ForkEpoch,
		ElectraForkEpoch: DefaultElectraForkEpoch,

		// State list length constants.
		EpochsPerHistoricalVector: DefaultEpochsPerHistoricalVector,
		EpochsPerSlashingsVector:  DefaultEpochsPerSlashingsVector,
		HistoricalRootsLimit:      DefaultHistoricalRootsLimit,
		ValidatorRegistryLimit:    DefaultValidatorRegistryLimit,

		// Slashing.
		InactivityPenaltyQuotient:      DefaultInactivityPenaltyQuotient,
		ProportionalSlashingMultiplier: DefaultProportionalSlashingMultiplier,

		// Capella values.
		MaxWithdrawalsPerPayload:         DefaultMaxWithdrawalsPerPayload,
		MaxValidatorsPerWithdrawalsSweep: DefaultMaxValidatorsPerWithdrawalsSweep,

		// Deneb values.
		MinEpochsForBlobsSidecarsRequest: DefaultMinEpochsForBlobsSidecarsRequest,
		MaxBlobCommitmentsPerBlock:       DefaultMaxBlobCommitmentsPerBlock,
		MaxBlobsPerBlock:                 DefaultMaxBlobsPerBlock,
		FieldElementsPerBlob:             DefaultFieldElementsPerBlob,
		BytesPerBlob:                     DefaultBytesPerBlob,
		KZGCommitmentInclusionProofDepth: DefaultKZGCommitmentInclusionProofDepth,

		// Berachain values.
		ValidatorSetCap:      DefaultValidatorSetCap,
		EVMInflationAddress:  common.NewExecutionAddressFromHex(DevnetEVMInflationAddress),
		EVMInflationPerBlock: DevnetEVMInflationPerBlock,
	}

	return chain.NewSpec(devnetSpecData)
}
