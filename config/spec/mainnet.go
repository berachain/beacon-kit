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
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
)

const (
	// MainnetEVMInflationAddressDeneb is the BGT contract address for genesis version Deneb.
	//
	// TODO: CONFIRM WITH SC TEAM!!!
	MainnetEVMInflationAddressDeneb = "0x289274787bAF083C15A45a174b7a8e44F0720660"

	// MainnetEVMInflationPerBlock is 5.75 BERA minted to the BGT contract per block
	// as the upper bound of redeemable BGT for genesis version Deneb.
	//
	// TODO: CONFIRM WITH QUANTUM TEAM!!!
	MainnetEVMInflationPerBlockDeneb = 5.75e9

	// MainnetValidatorSetCap is 69 on Mainnet for version Deneb.
	//
	// TODO: FIXME!!!
	MainnetValidatorSetCapDeneb = 69

	// MainnetMaxValidatorsPerWithdrawalsSweep is 31 because we expect at least 36
	// validators in the total validators set at genesis. We choose a prime number smaller
	// than the minimum amount of total validators possible.
	//
	// TODO: FIXME!!!
	MainnetMaxValidatorsPerWithdrawalsSweep = 31

	// MainnetMaxEffectiveBalance is the max stake of 10 million BERA for genesis version Deneb.
	MainnetMaxEffectiveBalanceDeneb = 10_000_000 * 1e9

	// MainnetEffectiveBalanceIncrementDeneb is 10k BERA for genesis version Deneb
	// (equivalent to the Deposit Contract's MIN_DEPOSIT_AMOUNT).
	MainnetEffectiveBalanceIncrementDeneb = 10_000 * 1e9

	// MainnetEjectionBalance is 240k BERA, calculated as:
	// activation_balance - effective_balance_increment = 250k - 10k = 240k BERA.
	// Activation balance is the min stake of 250k BERA for genesis version Deneb.
	MainnetEjectionBalanceDeneb = 240_000 * 1e9

	// MainnetSlotsPerEpochDeneb is 192 for genesis version Deneb
	// to mirror the time of epochs on Ethereum mainnet.
	//
	// TODO: FIXME!!! I really like 192 over 32 :)))
	MainnetSlotsPerEpochDeneb = 192

	// MainnetMinEpochsForBlobsSidecarsRequestDeneb is 4096 for genesis version Deneb
	// to match Ethereum mainnet.
	MainnetMinEpochsForBlobsSidecarsRequestDeneb = 4096

	// MainnetMaxBlobCommitmentsPerBlock is 4096 for genesis version Deneb
	// to match Ethereum mainnet.
	MainnetMaxBlobCommitmentsPerBlock = 4096
)

// MainnetChainSpec is the ChainSpec for the Berachain mainnet.
//
//nolint:mnd // okay to specify values here.
func MainnetChainSpec() (chain.Spec, error) {
	mainnetSpec := &chain.SpecData{
		// Gwei values constants.
		MaxEffectiveBalance:       MainnetMaxEffectiveBalanceDeneb,
		EjectionBalance:           MainnetEjectionBalanceDeneb,
		EffectiveBalanceIncrement: MainnetEffectiveBalanceIncrementDeneb,

		HysteresisQuotient:           DefaultHysteresisQuotient,
		HysteresisDownwardMultiplier: DefaultHysteresisDownwardMultiplier,
		HysteresisUpwardMultiplier:   DefaultHysteresisUpwardMultiplier,

		// Time parameters constants.
		SlotsPerEpoch:                MainnetSlotsPerEpochDeneb,
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
		DepositEth1ChainID:        MainnetEth1ChainID,
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
		MaxValidatorsPerWithdrawalsSweep: MainnetMaxValidatorsPerWithdrawalsSweep,

		// Deneb values.
		MinEpochsForBlobsSidecarsRequest: MainnetMinEpochsForBlobsSidecarsRequestDeneb,
		MaxBlobCommitmentsPerBlock:       MainnetMaxBlobCommitmentsPerBlock,
		MaxBlobsPerBlock:                 DefaultMaxBlobsPerBlock,
		FieldElementsPerBlob:             DefaultFieldElementsPerBlob,
		BytesPerBlob:                     DefaultBytesPerBlob,
		KZGCommitmentInclusionProofDepth: DefaultKZGCommitmentInclusionProofDepth,

		// Berachain values.
		ValidatorSetCap:      MainnetValidatorSetCapDeneb,
		EVMInflationAddress:  common.NewExecutionAddressFromHex(MainnetEVMInflationAddressDeneb),
		EVMInflationPerBlock: MainnetEVMInflationPerBlockDeneb,
	}

	return chain.NewSpec(mainnetSpec)
}
