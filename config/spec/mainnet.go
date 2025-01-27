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
	// BGT contract address.
	//
	// A hard fork will occur to set this value as the BGT contract address
	// when BGT beings to be minted.
	MainnetEVMInflationAddressDeneb = DefaultEVMInflationAddress

	// 0 BERA is minted to the BGT contract per block at genesis.
	//
	// A hard fork will occur to set this value as the upper bound of redeemable BGT per
	// block when BGT begins to be minted.
	MainnetEVMInflationPerBlockDeneb = DefaultEVMInflationPerBlock

	// MainnetValidatorSetCap is 69 on Mainnet for genesis version Deneb.
	MainnetValidatorSetCapDeneb = 69

	// MaxValidatorsPerWithdrawalsSweep is 31 because we expect at least 36
	// validators in the total validators set. We choose a prime number smaller
	// than the minimum amount of total validators possible.
	MainnetMaxValidatorsPerWithdrawalsSweep = 31

	// MainnetMaxEffectiveBalance is the max stake of 10 million BERA for genesis version Deneb.
	MainnetMaxEffectiveBalanceDeneb = 10_000_000 * 1e9

	// MainnetEffectiveBalanceIncrementDeneb is 10k BERA for genesis version Deneb
	// (equivalent to the Deposit Contract's MIN_DEPOSIT_AMOUNT).
	MainnetEffectiveBalanceIncrementDeneb = 10_000 * 1e9

	// MainnetEjectionBalance is 240k BERA, calculated as:
	// activation_balance - effective_balance_increment = 250k - 10k = 240k BERA.
	// Activation balance is the min stake of 250k BERA for genesis version Deneb.
	MainnetEjectionBalanceDeneb = 250_000*1e9 - MainnetEffectiveBalanceIncrementDeneb

	// Slots per epoch is 192 to mirror the time of epochs on Ethereum mainnet.
	MainnetSlotsPerEpochDeneb = 192

	// MainnetMinEpochsForBlobsSidecarsRequestDeneb is 4096 for genesis version Deneb
	// to match Ethereum mainnet.
	MainnetMinEpochsForBlobsSidecarsRequestDeneb = 4096

	// MainnetMaxBlobCommitmentsPerBlock is 4096 for genesis version Deneb
	// to match Ethereum mainnet.
	MainnetMaxBlobCommitmentsPerBlock = 4096

	// The deposit contract address on mainnet at Deneb is the same as the
	// default deposit contract address.
	MainnetDepositContractAddressDeneb = DefaultDepositContractAddress

	// MainnetDeneb1ForkEpoch is the epoch at which the Deneb1 fork occurs.
	//
	// TODO: set to the correct epoch.
	MainnetDeneb1ForkEpoch = 5000

	// MainnetEVMInflationAddressDeneb1 is the address on the EVM which will receive the
	// inflation amount of native EVM balance through a withdrawal every block in the Deneb1 fork.
	MainnetEVMInflationAddressDeneb1 = "0x656b95E550C07a9ffe548bd4085c72418Ceb1dba"

	// MainnetEVMInflationPerBlockDeneb1 is the amount of native EVM balance (in Gwei) to be
	// minted to the EVMInflationAddressDeneb1 via a withdrawal every block in the Deneb1 fork.
	MainnetEVMInflationPerBlockDeneb1 = 5.75 * 1e9
)

// MainnetChainSpec is the ChainSpec for the Berachain mainnet.
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
		DomainTypeVoluntaryExit:     bytes.FromUint32(DefaultDomainTypeVoluntaryExit),
		DomainTypeSelectionProof:    bytes.FromUint32(DefaultDomainTypeSelectionProof),
		DomainTypeAggregateAndProof: bytes.FromUint32(DefaultDomainTypeAggregateAndProof),
		DomainTypeApplicationMask:   bytes.FromUint32(DefaultDomainTypeApplicationMask),

		// Eth1-related values.
		DepositContractAddress:    common.NewExecutionAddressFromHex(MainnetDepositContractAddressDeneb),
		MaxDepositsPerBlock:       DefaultMaxDepositsPerBlock,
		DepositEth1ChainID:        MainnetEth1ChainID,
		Eth1FollowDistance:        DefaultEth1FollowDistance,
		TargetSecondsPerEth1Block: DefaultTargetSecondsPerEth1Block,

		// Fork-related values.
		Deneb1ForkEpoch:  MainnetDeneb1ForkEpoch,
		ElectraForkEpoch: DefaultElectraForkEpoch,

		// State list length constants.
		EpochsPerHistoricalVector: DefaultEpochsPerHistoricalVector,
		EpochsPerSlashingsVector:  DefaultEpochsPerSlashingsVector,
		HistoricalRootsLimit:      DefaultHistoricalRootsLimit,
		ValidatorRegistryLimit:    DefaultValidatorRegistryLimit,

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

		// Berachain values at genesis.
		ValidatorSetCap:      MainnetValidatorSetCapDeneb,
		EVMInflationAddress:  common.NewExecutionAddressFromHex(MainnetEVMInflationAddressDeneb),
		EVMInflationPerBlock: MainnetEVMInflationPerBlockDeneb,

		// Deneb1 values.
		EVMInflationAddressDeneb1:  common.NewExecutionAddressFromHex(MainnetEVMInflationAddressDeneb1),
		EVMInflationPerBlockDeneb1: MainnetEVMInflationPerBlockDeneb1,
	}

	return chain.NewSpec(mainnetSpec)
}
