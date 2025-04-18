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
	"github.com/ethereum/go-ethereum/params"
)

const (
	// BGT contract address.
	//
	// A hard fork will occur to set this value as the BGT contract address
	// when BGT beings to be minted.
	mainnetEVMInflationAddress = defaultEVMInflationAddress

	// 0 BERA is minted to the BGT contract per block at genesis.
	//
	// A hard fork will occur to set this value as the upper bound of redeemable BGT per
	// block when BGT begins to be minted.
	mainnetEVMInflationPerBlock = defaultEVMInflationPerBlock

	// mainnetValidatorSetCap is 69 on Mainnet at genesis.
	mainnetValidatorSetCap = 69

	// mainnetMaxValidatorsPerWithdrawalsSweep is 31 because we expect at least 36
	// validators in the total validators set. We choose a prime number smaller
	// than the minimum amount of total validators possible.
	mainnetMaxValidatorsPerWithdrawalsSweep = 31

	// mainnetMaxEffectiveBalance is the max stake of 10 million BERA at genesis.
	mainnetMaxEffectiveBalance = 10_000_000 * params.GWei

	// mainnetEffectiveBalanceIncrement is 10k BERA at genesis
	// (equivalent to the Deposit Contract's MIN_DEPOSIT_AMOUNT).
	mainnetEffectiveBalanceIncrement = 10_000 * params.GWei

	// mainnetEjectionBalance is 240k BERA, calculated as:
	// activation_balance - effective_balance_increment = 250k - 10k = 240k BERA.
	// Activation balance is the min stake of 250k BERA at genesis.
	mainnetEjectionBalance = 250_000*params.GWei - mainnetEffectiveBalanceIncrement

	// mainnetSlotsPerEpoch is 192 to mirror the time of epochs on Ethereum mainnet.
	mainnetSlotsPerEpoch = 192

	// mainnetMinEpochsForBlobsSidecarsRequest is 4096 at genesis to match Ethereum mainnet.
	mainnetMinEpochsForBlobsSidecarsRequest = defaultMinEpochsForBlobsSidecarsRequest

	// mainnetMaxBlobCommitmentsPerBlock is 4096 at genesis to match Ethereum mainnet.
	mainnetMaxBlobCommitmentsPerBlock = defaultMaxBlobCommitmentsPerBlock

	// The deposit contract address on mainnet at genesis is the same as the
	// default deposit contract address.
	mainnetDepositContractAddress = defaultDepositContractAddress

	// mainnetGenesisTime is the timestamp of the Berachain mainnet genesis block.
	mainnetGenesisTime = 1737381600

	// mainnetDeneb1ForkTime is the timestamp at which the Deneb1 fork occurs.
	// This is calculated based on the timestamp of the 2855th mainnet epoch, block 548160, which
	// was used to initiate the fork when beacon-kit forked by epoch instead of by timestamp.
	mainnetDeneb1ForkTime = 1738415507

	// mainnetEVMInflationAddressDeneb1 is the address on the EVM which will receive the
	// inflation amount of native EVM balance through a withdrawal every block in the Deneb1 fork.
	mainnetEVMInflationAddressDeneb1 = "0x656b95E550C07a9ffe548bd4085c72418Ceb1dba"

	// mainnetEVMInflationPerBlockDeneb1 is the amount of native EVM balance (in Gwei) to be
	// minted to the EVMInflationAddressDeneb1 via a withdrawal every block in the Deneb1 fork.
	mainnetEVMInflationPerBlockDeneb1 = 5.75 * params.GWei
)

// MainnetChainSpecData is the chain.SpecData for the Berachain mainnet.
func MainnetChainSpecData() *chain.SpecData {
	return &chain.SpecData{
		// Gwei values constants.
		MaxEffectiveBalance:       mainnetMaxEffectiveBalance,
		EjectionBalance:           mainnetEjectionBalance,
		EffectiveBalanceIncrement: mainnetEffectiveBalanceIncrement,

		HysteresisQuotient:           defaultHysteresisQuotient,
		HysteresisDownwardMultiplier: defaultHysteresisDownwardMultiplier,
		HysteresisUpwardMultiplier:   defaultHysteresisUpwardMultiplier,

		// Time parameters constants.
		SlotsPerEpoch:                mainnetSlotsPerEpoch,
		SlotsPerHistoricalRoot:       defaultSlotsPerHistoricalRoot,
		MinEpochsToInactivityPenalty: defaultMinEpochsToInactivityPenalty,

		// Signature domains.
		DomainTypeProposer:          bytes.FromUint32(defaultDomainTypeProposer),
		DomainTypeAttester:          bytes.FromUint32(defaultDomainTypeAttester),
		DomainTypeRandao:            bytes.FromUint32(defaultDomainTypeRandao),
		DomainTypeDeposit:           bytes.FromUint32(defaultDomainTypeDeposit),
		DomainTypeVoluntaryExit:     bytes.FromUint32(defaultDomainTypeVoluntaryExit),
		DomainTypeSelectionProof:    bytes.FromUint32(defaultDomainTypeSelectionProof),
		DomainTypeAggregateAndProof: bytes.FromUint32(defaultDomainTypeAggregateAndProof),
		DomainTypeApplicationMask:   bytes.FromUint32(defaultDomainTypeApplicationMask),

		// Eth1-related values.
		DepositContractAddress:    common.NewExecutionAddressFromHex(mainnetDepositContractAddress),
		MaxDepositsPerBlock:       defaultMaxDepositsPerBlock,
		DepositEth1ChainID:        MainnetEth1ChainID,
		Eth1FollowDistance:        defaultEth1FollowDistance,
		TargetSecondsPerEth1Block: defaultTargetSecondsPerEth1Block,

		// Fork-related values.
		GenesisTime:     mainnetGenesisTime,
		Deneb1ForkTime:  mainnetDeneb1ForkTime,
		ElectraForkTime: defaultElectraForkTime,

		// State list length constants.
		EpochsPerHistoricalVector: defaultEpochsPerHistoricalVector,
		EpochsPerSlashingsVector:  defaultEpochsPerSlashingsVector,
		HistoricalRootsLimit:      defaultHistoricalRootsLimit,
		ValidatorRegistryLimit:    defaultValidatorRegistryLimit,

		// Capella values.
		MaxWithdrawalsPerPayload:         defaultMaxWithdrawalsPerPayload,
		MaxValidatorsPerWithdrawalsSweep: mainnetMaxValidatorsPerWithdrawalsSweep,

		// Deneb values.
		MinEpochsForBlobsSidecarsRequest: mainnetMinEpochsForBlobsSidecarsRequest,
		MaxBlobCommitmentsPerBlock:       mainnetMaxBlobCommitmentsPerBlock,
		MaxBlobsPerBlock:                 defaultMaxBlobsPerBlock,
		FieldElementsPerBlob:             defaultFieldElementsPerBlob,
		BytesPerBlob:                     defaultBytesPerBlob,
		KZGCommitmentInclusionProofDepth: defaultKZGCommitmentInclusionProofDepth,

		// Berachain values at genesis.
		ValidatorSetCap:             mainnetValidatorSetCap,
		EVMInflationAddressGenesis:  common.NewExecutionAddressFromHex(mainnetEVMInflationAddress),
		EVMInflationPerBlockGenesis: mainnetEVMInflationPerBlock,

		// Deneb1 values.
		EVMInflationAddressDeneb1:  common.NewExecutionAddressFromHex(mainnetEVMInflationAddressDeneb1),
		EVMInflationPerBlockDeneb1: mainnetEVMInflationPerBlockDeneb1,
	}
}

// MainnetChainSpec is the ChainSpec for the Berachain mainnet.
func MainnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(MainnetChainSpecData())
}
