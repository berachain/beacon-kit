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
	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BoonetChainSpec is the ChainSpec for the localnet.
func BoonetChainSpec() (chain.Spec[
	common.DomainType,
	math.Epoch,
	math.Slot,
	any,
], error) {
	boonetSpec := BaseSpec()

	// Chain ID is 80000.
	boonetSpec.DepositEth1ChainID = BoonetEth1ChainID

	// BGT contract address.
	boonetSpec.EVMInflationAddress = common.NewExecutionAddressFromHex(
		"0x289274787bAF083C15A45a174b7a8e44F0720660",
	)

	// Basic parameters
	boonetSpec.EVMInflationPerBlock = 2.5e9
	boonetSpec.ValidatorSetCap = 256
	boonetSpec.MaxValidatorsPerWithdrawalsSweepPostUpgrade = 43
	boonetSpec.MaxEffectiveBalancePostUpgrade = 5_000_000 * 1e9

	spec, err := chain.NewChainSpec(boonetSpec)
	if err != nil {
		return nil, err
	}

	// Adding parameters for different block heights
	heightParams := []chain.HeightDependentParams{
		{
			Height:               0,
			MaxEffectiveBalance:  32e9,
			ValidatorSetCap:      256,
			EVMInflationPerBlock: 2.5e9,
		},
		{
			Height:               BoonetFork1Height,
			MaxEffectiveBalance:  64e9,
			ValidatorSetCap:      512,
			EVMInflationPerBlock: 3e9,
		},
		{
			Height:               BoonetFork2Height,
			MaxEffectiveBalance:  128e9,
			ValidatorSetCap:      1024,
			EVMInflationPerBlock: 3.5e9,
		},
		{
			Height:               BoonetFork3Height,
			MaxEffectiveBalance:  5_000_000 * 1e9,
			ValidatorSetCap:      2048,
			EVMInflationPerBlock: 4e9,
		},
	}

	spec.(*chain.ChainSpec[common.DomainType, math.Epoch, math.Slot, any]).HeightParams = heightParams

	return spec, nil
}
