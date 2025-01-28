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
)

const (
	// DevnetEVMInflationAddress is the address of the EVM inflation contract.
	DevnetEVMInflationAddress = "0x6942069420694206942069420694206942069420"

	// DevnetEVMInflationPerBlock is the amount of native EVM balance (in units
	// of Gwei) to be minted per EL block.
	DevnetEVMInflationPerBlock = 10e9

	// MaxStakeAmount is the maximum amount of native EVM balance (in units
	// of Gwei) that can be staked.
	DevnetMaxStakeAmount = 4000e9
)

// DevnetChainSpecData is the chain.SpecData for a devnet. It is similar to mainnet but
// has different values for testing EVM inflation and staking.
//
// TODO: remove modifications from mainnet spec to align with mainnet behavior.
func DevnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()
	specData.DepositEth1ChainID = DevnetEth1ChainID

	// EVM inflation is different from mainnet to test.
	specData.EVMInflationAddress = common.NewExecutionAddressFromHex(DevnetEVMInflationAddress)
	specData.EVMInflationPerBlock = DevnetEVMInflationPerBlock

	// Staking is different from mainnet for now.
	specData.MaxEffectiveBalance = DevnetMaxStakeAmount
	specData.EjectionBalance = defaultEjectionBalance
	specData.EffectiveBalanceIncrement = defaultEffectiveBalanceIncrement
	specData.SlotsPerEpoch = defaultSlotsPerEpoch

	return specData
}

// DevnetChainSpec is the chain.Spec for a devnet. Used by `make start` and unit tests.
func DevnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(DevnetChainSpecData())
}
