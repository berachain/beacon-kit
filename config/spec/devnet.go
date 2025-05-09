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
)

const (
	// devnetGenesisTime is the timestamp of devnet genesis.
	devnetGenesisTime = 0

	// devnetDeneb1ForkTime is the timestamp at which the Deneb1 fork occurs.
	devnetDeneb1ForkTime = 0

	// devnetElectraForkTime is the timestamp at which the Electra fork occurs.
	// devnet is configured to start on electra.
	devnetElectraForkTime = 0
)

// DevnetChainSpecData is the chain.SpecData for a devnet. We try to keep this
// as close to the mainnet spec as possible.
func DevnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()
	specData.DepositEth1ChainID = DevnetEth1ChainID

	// Fork timings are set to facilitate local testing across fork versions.
	specData.GenesisTime = devnetGenesisTime
	specData.Deneb1ForkTime = devnetDeneb1ForkTime
	specData.ElectraForkTime = devnetElectraForkTime

	// Use fewer slots per epoch for devnet to speed up testing.
	specData.SlotsPerEpoch = 32

	return specData
}

// DevnetChainSpec is the chain.Spec for a devnet. Used by `make start` and unit tests.
func DevnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(DevnetChainSpecData())
}
