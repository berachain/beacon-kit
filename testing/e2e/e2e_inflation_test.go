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
	"sync"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// TestEVMInflation checks that the EVM inflation address receives the correct
// amount of EVM inflation per block.
func (s *BeaconKitE2ESuite) TestEVMInflation() {
	// TODO: make test use configurable chain spec.
	chainspec, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	var (
		inflationPerBlock          uint64
		inflationAddress           common.ExecutionAddress
		oldInflationAddress        common.ExecutionAddress
		preForkAddressFinalBalance *big.Int
		preForkLatestBalance       *big.Int
		balance                    *big.Int
		expectedBalance            *big.Int
		forkSlot                   int64
		onceOnFork                 sync.Once
	)
	// Arbitrarily run test for 2 epochs.
	for blkNum := range int64(2 * chainspec.SlotsPerEpoch()) {
		err = s.WaitForFinalizedBlockNumber(uint64(blkNum))
		s.Require().NoError(err)
		payload, errBlk := s.JSONRPCBalancer().BlockByNumber(s.Ctx(), big.NewInt(blkNum))
		s.Require().NoError(errBlk)

		payloadTime := payload.Time()
		inflationPerBlock = chainspec.EVMInflationPerBlock(math.U64(payloadTime)).Unwrap()
		inflationAddress = chainspec.EVMInflationAddress(math.U64(payloadTime))
		if chainspec.Deneb1ForkTime() > 0 && payloadTime >= chainspec.Deneb1ForkTime() {
			// If we have passed the Deneb1 fork, do some verifications and update inflation values.
			onceOnFork.Do(func() {
				oldInflationPerBlock := chainspec.EVMInflationPerBlock(math.U64(chainspec.Deneb1ForkTime() - 1))
				oldInflationAddress = chainspec.EVMInflationAddress(math.U64(chainspec.Deneb1ForkTime() - 1))

				// Verify the post fork inflation changes
				s.Require().NotEqual(oldInflationPerBlock, inflationPerBlock)
				s.Require().NotEqual(oldInflationAddress, inflationAddress)
				forkSlot = blkNum

				// take the snapshot of balance right before the fork and check it won't change anymore
				preForkAddressFinalBalance, err = s.JSONRPCBalancer().BalanceAt(
					s.Ctx(), gethcommon.Address(oldInflationAddress), big.NewInt(blkNum-1),
				)
				s.Require().NoError(err)
			})

			// Enforce that the balance of the EVM inflation address
			// prior to the hardfork is the same as it is now.
			preForkLatestBalance, err = s.JSONRPCBalancer().BalanceAt(
				s.Ctx(), gethcommon.Address(oldInflationAddress), nil, // at the current block
			)
			s.Require().NoError(err)
			s.Require().Zero(preForkAddressFinalBalance.Cmp(preForkLatestBalance))

			expectedBalance = new(big.Int).Mul(
				new(big.Int).SetUint64(inflationPerBlock*params.GWei),
				big.NewInt(blkNum-forkSlot+1),
			)
		} else {
			// Pre-Deneb1
			expectedBalance = new(big.Int).Mul(
				new(big.Int).SetUint64(inflationPerBlock*params.GWei),
				big.NewInt(blkNum),
			)
		}

		balance, err = s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.Address(inflationAddress),
			big.NewInt(blkNum),
		)
		s.Require().NoError(err)
		s.Require().Zero(balance.Cmp(expectedBalance),
			"height", blkNum,
			"balance", balance,
			"expectedBalance", expectedBalance,
		)
	}
}
