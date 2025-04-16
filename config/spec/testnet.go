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

import "github.com/berachain/beacon-kit/chain"

// TestnetChainSpecData is the chain.SpecData for Berachain's public testnet.
func TestnetChainSpecData() *chain.SpecData {
	specData := MainnetChainSpecData()

	// Testnet uses chain ID of 80069.
	specData.DepositEth1ChainID = TestnetEth1ChainID

	// Timestamp of the genesis block of Bepolia testnet.
	specData.GenesisTime = 1739976735

	// Deneb1 fork timing on Bepolia. This is calculated based on the timestamp of the first bepolia
	// epoch, block 192, which was used to initiate the fork when beacon-kit forked by epoch instead
	// of by timestamp.
	specData.Deneb1ForkTime = 1740090694

	return specData
}

// TestnetChainSpec is the chain.Spec for Berachain's public testnet.
func TestnetChainSpec() (chain.Spec, error) {
	return chain.NewSpec(TestnetChainSpecData())
}
