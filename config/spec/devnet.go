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
	// devnetGenesisTime is the timestamp of devnet genesis.
	devnetGenesisTime = 0

	// devnetDeneb1ForkTime is the timestamp at which the Deneb1 fork occurs.
	devnetDeneb1ForkTime = 0

	// devnetElectraForkTime is the timestamp at which the Electra fork occurs.
	devnetElectraForkTime = 0

	// devnetElectra1ForkTime is the timestamp at which the Electra1 fork occurs.
	// devnet is configured to start on electra1.
	devnetElectra1ForkTime = 0

	// devnetEVMInflationAddressDeneb1 is the address of the EVM inflation contract after the Deneb1 fork.
	devnetEVMInflationAddressDeneb1 = "0xEE0BD9569e41fA26A79305Fc31a663986Deb79FB"

	// devnetEVMInflationPerBlockDeneb1 is the amount of native EVM balance (in units
	// of Gwei) to be minted per EL block after the Deneb1 fork.
	devnetEVMInflationPerBlockDeneb1 = 11 * params.GWei

	// devnetMinValidatorWithdrawabilityDelay is the delay (in epochs) before a validator can withdraw their stake.
	devnetMinValidatorWithdrawabilityDelay = 32

	// devnetFuluForkTime is the timestamp at which the Fulu fork occurs on devnet.
	// Set to 0 so devnet starts with Fulu active.
	devnetFuluForkTime = 0

	// devnetEVMInflationAddressFulu is the address of the EVM inflation contract
	// after the Fulu fork on devnet.
	devnetEVMInflationAddressFulu = "0x4A2452Fd7e9FCA389d98063c5C3A8FC63838E451"

	// devnetEVMInflationPerBlockFulu is the amount of native EVM balance (in units
	// of Gwei) to be minted per EL block after the Fulu fork on devnet.
	devnetEVMInflationPerBlockFulu = 12 * params.GWei
)

// DevnetChainSpecData is the chain.SpecData for a devnet. We try to keep this
// as close to the mainnet spec as possible.
func DevnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()
	specData.DepositEth1ChainID = chain.DevnetEth1ChainID

	specData.Config.ConsensusUpdateHeight = 1
	specData.Config.ConsensusEnableHeight = 2

	// Fork timings are set to facilitate local testing across fork versions.
	specData.GenesisTime = devnetGenesisTime
	specData.Deneb1ForkTime = devnetDeneb1ForkTime
	specData.ElectraForkTime = devnetElectraForkTime
	specData.Electra1ForkTime = devnetElectra1ForkTime
	specData.FuluForkTime = devnetFuluForkTime

	// EVM inflation is different from mainnet for now, after the Deneb1 fork.
	specData.EVMInflationAddressDeneb1 = common.MustNewExecutionAddressFromHex(devnetEVMInflationAddressDeneb1)
	specData.EVMInflationPerBlockDeneb1 = devnetEVMInflationPerBlockDeneb1

	// Validator withdrawability delay is set to a low value to speed up testing.
	specData.MinValidatorWithdrawabilityDelay = devnetMinValidatorWithdrawabilityDelay

	// EVM inflation for the Fulu fork on devnet. The address remains the same as the Deneb1 fork.
	specData.EVMInflationAddressFulu = common.MustNewExecutionAddressFromHex(devnetEVMInflationAddressFulu)
	specData.EVMInflationPerBlockFulu = devnetEVMInflationPerBlockFulu

	// Use fewer slots per epoch for devnet to speed up testing.
	specData.SlotsPerEpoch = 32

	return specData
}

// DevnetChainSpec is the chain.Spec for a devnet. Used by `make start` and unit tests.
func DevnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(DevnetChainSpecData())
}
