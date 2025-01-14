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
	"github.com/berachain/beacon-kit/primitives/common"
)

// BoonetChainSpec is the ChainSpec for the localnet.
func BoonetChainSpec() (chain.Spec, error) {
	boonetSpec := BaseSpec()

	// Chain ID is 80094.
	boonetSpec.DepositEth1ChainID = BoonetEth1ChainID

	// Target for block time is 2 seconds on Berachain mainnet.
	boonetSpec.TargetSecondsPerEth1Block = 2

	// BGT contract address.
	boonetSpec.EVMInflationAddress = common.NewExecutionAddressFromHex(
		"0x289274787bAF083C15A45a174b7a8e44F0720660",
	)

	// 0 BERA is minted to the BGT contract per block at genesis.
	//
	// A hard fork will occur to set this value as the upper bound of redeemable BGT per
	// block when BGT begins to be minted.
	boonetSpec.EVMInflationPerBlock = 0

	// ValidatorSetCap is 69 on Mainnet for version Deneb at genesis.
	boonetSpec.ValidatorSetCap = 6

	// MaxValidatorsPerWithdrawalsSweep is 31 because we expect at least 36
	// validators in the total validators set. We choose a prime number smaller
	// than the minimum amount of total validators possible.
	boonetSpec.MaxValidatorsPerWithdrawalsSweep = 7

	// MaxEffectiveBalance (or max stake) is 10 million BERA.
	boonetSpec.MaxEffectiveBalance = 10_000_000 * 1e9

	// Effective balance increment is 10k BERA
	// (equivalent to the Deposit Contract's MIN_DEPOSIT_AMOUNT).
	boonetSpec.EffectiveBalanceIncrement = 10_000 * 1e9

	// Since the activation balance (min stake) is 250k BERA, we set the ejection balance be
	// activation_balance - effective_balance_increment = 250k - 10k = 240k BERA.
	boonetSpec.EjectionBalance = 250_000*1e9 - boonetSpec.EffectiveBalanceIncrement

	// Slots per epoch is 192 to mirror the time of epochs on Ethereum mainnet.
	boonetSpec.SlotsPerEpoch = 192

	// MinEpochsForBlobsSidecarsRequest is 4096 to match Ethereum mainnet.
	boonetSpec.MinEpochsForBlobsSidecarsRequest = 4096

	// MaxBlobCommitmentsPerBlock is 4096 to match Ethereum mainnet.
	boonetSpec.MaxBlobCommitmentsPerBlock = 4096

	return chain.NewSpec(boonetSpec)
}
