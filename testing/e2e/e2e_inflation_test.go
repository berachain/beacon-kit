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
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TestEVMInflation checks that the EVM inflation address receives the correct
// amount of EVM inflation per block.
func (s *BeaconKitE2ESuite) TestEVMInflation() {
	// TODO: make test use configurable chain spec.
	chainspec, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	deneb1ForkSlot := chainspec.SlotsPerEpoch() * uint64(chainspec.Deneb1ForkEpoch())

	// Check over the first epoch before Deneb1, the balance of the Devnet EVM inflation address
	// increases by DevnetEVMInflationPerBlock.
	for i := range int64(deneb1ForkSlot) {
		err := s.WaitForFinalizedBlockNumber(uint64(i))
		s.Require().NoError(err)

		expectedBalance := new(big.Int).Mul(
			big.NewInt(int64(chainspec.EVMInflationPerBlock(math.Slot(i)))),
			big.NewInt(i),
		)

		balance, err := s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.Address(chainspec.EVMInflationAddress(math.Slot(i))),
			big.NewInt(i),
		)
		s.Require().NoError(err)
		s.Require().Zero(balance.Cmp(expectedBalance))
	}

	// Check over the first epoch after Deneb1, the balance of the Devnet EVM inflation address
	// post Deneb1 increases by DevnetEVMInflationPerBlockDeneb1.
	for i := deneb1ForkSlot; i < deneb1ForkSlot+chainspec.SlotsPerEpoch(); i++ {
		err := s.WaitForFinalizedBlockNumber(uint64(i))
		s.Require().NoError(err)

		expectedBalance := new(big.Int).Mul(
			big.NewInt(int64(chainspec.EVMInflationPerBlock(math.Slot(i)))),
			big.NewInt(int64(i)),
		)

		balance, err := s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.Address(chainspec.EVMInflationAddress(math.Slot(i))),
			big.NewInt(int64(i)),
		)
		s.Require().NoError(err)
		s.Require().Zero(balance.Cmp(expectedBalance))
	}

	// If the addresses are different, enforce that the balance of the EVM inflation address
	// prior to the hardfork is the same as it is now.
	priorEVMInflationAddress := chainspec.EVMInflationAddress(constants.GenesisSlot)
	postEVMInflationAddress := chainspec.EVMInflationAddress(math.Slot(deneb1ForkSlot))
	if priorEVMInflationAddress != postEVMInflationAddress {
		balanceRightBeforeHardfork, err := s.JSONRPCBalancer().BalanceAt(
			s.Ctx(), gethcommon.Address(priorEVMInflationAddress), big.NewInt(int64(deneb1ForkSlot)),
		)
		s.Require().NoError(err)
		balanceAfterHardfork, err := s.JSONRPCBalancer().BalanceAt(
			s.Ctx(), gethcommon.Address(postEVMInflationAddress), nil, // at the current block
		)
		s.Require().NoError(err)
		s.Require().Zero(balanceRightBeforeHardfork.Cmp(balanceAfterHardfork))
	}
}
