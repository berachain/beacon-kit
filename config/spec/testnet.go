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

// TestnetChainSpecData is the chain.SpecData for Berachain's public testnet.
func TestnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()
	specData.DepositEth1ChainID = TestnetEth1ChainID

	// Genesis values of EVM inflation are consistent with Deneb1 to keep BERA minting on.
	specData.EVMInflationAddressGenesis = common.NewExecutionAddressFromHex(mainnetEVMInflationAddressDeneb1)
	specData.EVMInflationPerBlockGenesis = mainnetEVMInflationPerBlockDeneb1

	// Unlike mainnet, testnet activates Deneb1 at epoch 1.
	specData.Deneb1ForkEpoch = 1
	specData.EVMInflationAddressDeneb1 = common.NewExecutionAddressFromHex(mainnetEVMInflationAddressDeneb1)
	specData.EVMInflationPerBlockDeneb1 = mainnetEVMInflationPerBlockDeneb1

	return specData
}

// TestnetChainSpec is the chain.Spec for Berachain's public testnet.
func TestnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(TestnetChainSpecData())
}
