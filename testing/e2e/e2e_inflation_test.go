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

package e2e_test

import (
	"math/big"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/math"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TestEVMInflation checks that the EVM inflation address receives the correct
// amount of EVM inflation per block.
func (s *BeaconKitE2ESuite) runEVMInflation() {
	s.Logger().Info("Running TestEVMInflation")
	evmInflationPerBlockWei, _ := big.NewFloat(
		spec.DevnetEVMInflationPerBlock * math.GweiPerWei,
	).Int(nil)

	s.Logger().Info("EVM Inflation Per Block Wei", "evmInflationPerBlockWei", evmInflationPerBlockWei)

	// Check over the next 10 EVM blocks, that after every block, the balance
	// of the EVM inflation address increases by DevnetEVMInflationPerBlock.
	for i := range int64(10) {
		err := s.WaitForFinalizedBlockNumber(uint64(i))
		s.Require().NoError(err)
		s.Logger().Info("Waiting for finalized block number", "blockNumber", i)

		s.Logger().Info("Balance at", "blockNumber", i, "address", spec.DevnetEVMInflationAddress)
		s.Logger().Info("jsonrpc balancer", "jsonrpcBalancer", s.JSONRPCBalancer())
		balance, err := s.JSONRPCBalancer().BalanceAt(
			s.Ctx(),
			gethcommon.HexToAddress(spec.DevnetEVMInflationAddress),
			big.NewInt(i),
		)
		s.Logger().Info("Balance at", "balance", balance)
		s.Require().NoError(err)
		s.Require().Zero(balance.Cmp(new(big.Int).Mul(
			evmInflationPerBlockWei, big.NewInt(i)),
		))
	}
}
