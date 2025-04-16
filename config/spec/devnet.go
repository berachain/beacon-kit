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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// devnetEVMInflationAddress is the address of the EVM inflation contract.
	devnetEVMInflationAddress = "0x6942069420694206942069420694206942069420"

	// devnetEVMInflationPerBlock is the amount of native EVM balance (in units
	// of Gwei) to be minted per EL block.
	devnetEVMInflationPerBlock = 10 * params.GWei

	// devnetMaxStakeAmount is the maximum amount of native EVM balance (in units
	// of Gwei) that can be staked.
	devnetMaxStakeAmount = 4000 * params.GWei

	// devnetGenesisTime is the timestamp of devnet genesis.
	devnetGenesisTime = 0

	// devnetDeneb1ForkTime is the timestamp at which the Deneb1 fork occurs.
	devnetDeneb1ForkTime = 0

	// devnetElectraForkTime is the timestamp at which the Electra fork occurs.
	// devnet is configured to start on electra.
	devnetElectraForkTime = 0

	// devnetEVMInflationAddressDeneb1 is the address of the EVM inflation contract
	// after the Deneb1 fork.
	devnetEVMInflationAddressDeneb1 = "0x4206942069420694206942069420694206942069"

	// devnetEVMInflationPerBlockDeneb1 is the amount of native EVM balance (in units
	// of Gwei) to be minted per EL block after the Deneb1 fork.
	devnetEVMInflationPerBlockDeneb1 = 11 * params.GWei
)

// DevnetChainSpecData is the chain.SpecData for a devnet. It is similar to mainnet but
// has different values for testing EVM inflation, staking, and hard forks.
//
// TODO: remove modifications from mainnet spec to align with mainnet behavior.
func DevnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()
	specData.DepositEth1ChainID = DevnetEth1ChainID

	// Fork timings are set to facilitate local testing across fork versions.
	specData.GenesisTime = devnetGenesisTime
	specData.Deneb1ForkTime = devnetDeneb1ForkTime
	specData.ElectraForkTime = devnetElectraForkTime

	// EVM inflation is different from mainnet to test.
	specData.EVMInflationAddressGenesis = common.NewExecutionAddressFromHex(devnetEVMInflationAddress)
	specData.EVMInflationPerBlockGenesis = devnetEVMInflationPerBlock

	// EVM inflation is different from mainnet for now, after the Deneb1 fork.
	specData.EVMInflationAddressDeneb1 = common.NewExecutionAddressFromHex(devnetEVMInflationAddressDeneb1)
	specData.EVMInflationPerBlockDeneb1 = devnetEVMInflationPerBlockDeneb1

	// Staking is different from mainnet for now.
	specData.MaxEffectiveBalance = devnetMaxStakeAmount
	specData.EjectionBalance = defaultEjectionBalance
	specData.EffectiveBalanceIncrement = defaultEffectiveBalanceIncrement
	specData.SlotsPerEpoch = defaultSlotsPerEpoch

	return specData
}

// DevnetChainSpec is the chain.Spec for a devnet. Used by `make start` and unit tests.
func DevnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(DevnetChainSpecData())
}
