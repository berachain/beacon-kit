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

	// BERA per block minting.
	boonetSpec.EVMInflationPerBlock = 2.5e9

	// ValidatorSetCap is 256 on the Boonet chain.
	boonetSpec.ValidatorSetCap = 256

	// MaxValidatorsPerWithdrawalsSweep is 43 because we expect at least 46
	// validators in the total validators set. We choose a prime number smaller
	// than the minimum amount of total validators possible.
	boonetSpec.MaxValidatorsPerWithdrawalsSweepPostUpgrade = 43

	// MaxEffectiveBalancePostUpgrade is 5 million BERA after the boonet
	// upgrade.
	//
	//nolint:mnd // ok.
	boonetSpec.MaxEffectiveBalancePostUpgrade = 5_000_000 * 1e9

	return chain.NewChainSpec(boonetSpec)
}
