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

package e2e_test

import (
	"math/big"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/math"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// runEVMInflation checks that the EVM inflation address receives the correct
// amount of EVM inflation per block.
func (s *BeaconKitE2ESuite) runEVMInflation() {
  s.Logger().Info("Running TestEVMInflation")
	// TODO: make test use configurable chain spec.
	chainspec, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	deneb1ForkSlot := chainspec.SlotsPerEpoch() * uint64(chainspec.Deneb1ForkEpoch())

	// Check over the first epoch before Deneb1, the balance of the Devnet EVM inflation address
	// increases by DevnetEVMInflationPerBlock.
	preForkInflation := chainspec.EVMInflationPerBlock(math.Slot(0))
	preForkAddress := chainspec.EVMInflationAddress(math.Slot(0))
	for blkNum := range int64(deneb1ForkSlot) {
		err = s.WaitForFinalizedBlockNumber(uint64(blkNum))
		s.Require().NoError(err)

		expectedBalance := new(big.Int).Mul(
			new(big.Int).SetUint64(preForkInflation*math.GweiPerWei),
			big.NewInt(blkNum),
		)

		var balance *big.Int
		balance, err = s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.Address(preForkAddress),
			big.NewInt(blkNum),
		)
		s.Require().NoError(err)
		s.Require().Zero(balance.Cmp(expectedBalance),
			"height", blkNum,
			"balance", balance,
			"expectedBalance", expectedBalance,
		)
	}

	// Check over the first epoch after Deneb1, the balance of the Devnet EVM inflation address
	// post Deneb1 increases by DevnetEVMInflationPerBlockDeneb1.
	postForkInflation := chainspec.EVMInflationPerBlock(math.Slot(deneb1ForkSlot))
	s.Require().NotEqual(preForkInflation, postForkInflation)

	postForkAddress := chainspec.EVMInflationAddress(math.Slot(deneb1ForkSlot))
	s.Require().NotEqual(preForkAddress, postForkAddress)

	// take the snapshot of balance right before the fork and check it won't change anymore
	var preForkAddressFinalBalance *big.Int
	preForkAddressFinalBalance, err = s.JSONRPCBalancer().BalanceAt(
		s.Ctx(), gethcommon.Address(preForkAddress), big.NewInt(int64(deneb1ForkSlot-1)),
	)
	s.Require().NoError(err)

	for blkNum := deneb1ForkSlot; blkNum < deneb1ForkSlot+chainspec.SlotsPerEpoch(); blkNum++ {
		err = s.WaitForFinalizedBlockNumber(blkNum)
		s.Require().NoError(err)

		expectedBalance := new(big.Int).Mul(
			new(big.Int).SetUint64(postForkInflation*math.GweiPerWei),
			big.NewInt(int64(blkNum-(deneb1ForkSlot-1))),
		)

		var balance *big.Int
		balance, err = s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.Address(postForkAddress),
			big.NewInt(int64(blkNum)),
		)
		s.Require().NoError(err)
		s.Require().Zero(balance.Cmp(expectedBalance),
			"height", blkNum,
			"balance", balance,
			"expectedBalance", expectedBalance,
		)

		// Enforce that the balance of the EVM inflation address
		// prior to the hardfork is the same as it is now.
		var preForkLatestBalance *big.Int
		preForkLatestBalance, err = s.JSONRPCBalancer().BalanceAt(
			s.Ctx(), gethcommon.Address(preForkAddress), nil, // at the current block
		)
		s.Require().NoError(err)
		s.Require().Zero(preForkAddressFinalBalance.Cmp(preForkLatestBalance))
	}
}
