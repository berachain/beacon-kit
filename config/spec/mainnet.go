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

// MainnetChainSpec is the ChainSpec for the Berachain mainnet.
//
//nolint:mnd // okay to specify values here.
func MainnetChainSpec() (chain.Spec, error) {
	mainnetSpec := BaseSpec()

	// Chain ID is ???.
	mainnetSpec.DepositEth1ChainID = MainnetEth1ChainID

	// Target for block time is 2 seconds on Berachain mainnet.
	mainnetSpec.TargetSecondsPerEth1Block = 2

	// BGT contract address. // TODO: CONFIRM WITH SC TEAM!!!
	mainnetSpec.EVMInflationAddress = common.NewExecutionAddressFromHex(
		"0x289274787bAF083C15A45a174b7a8e44F0720660",
	)

	// 5.75 BERA is minted to the BGT contract per block as the upper bound of redeemable BGT.
	//
	// TODO: CONFIRM WITH QUANTUM TEAM!!!
	mainnetSpec.EVMInflationPerBlock = 5.75e9

	// ValidatorSetCap is 69 on Mainnet for version Deneb.
	mainnetSpec.ValidatorSetCap = 69 // TODO: FIXME!!!

	// MaxValidatorsPerWithdrawalsSweep is 31 because we expect at least 36
	// validators in the total validators set. We choose a prime number smaller
	// than the minimum amount of total validators possible.
	mainnetSpec.MaxValidatorsPerWithdrawalsSweep = 31 // TODO: FIXME!!!

	// MaxEffectiveBalance (or max stake) is 10 million BERA.
	mainnetSpec.MaxEffectiveBalance = 10_000_000 * 1e9

	// Ejection balance (or min stake) is 250k BERA.
	mainnetSpec.EjectionBalance = 250_000 * 1e9

	// Effective balance increment is 10k BERA
	// (equivalent to the Deposit Contract's MIN_DEPOSIT_AMOUNT).
	mainnetSpec.EffectiveBalanceIncrement = 10_000 * 1e9

	// Slots per epoch is 192 to mirror the time of epochs on Ethereum mainnet.
	//
	// TODO: FIXME!!! I really like 192 over 32 :)))
	mainnetSpec.SlotsPerEpoch = 192

	return chain.NewSpec(mainnetSpec)
}
