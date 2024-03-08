// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package staking

import (
	"context"

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	stakinglogs "github.com/berachain/beacon-kit/beacon/staking/logs"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// ProcessBlockEvents processes the logs from the deposit contract.
func (s *Service) ProcessBlockEvents(
	ctx context.Context,
	logs []coretypes.Log,
) error {
	for _, log := range logs {
		// We only care about logs from the deposit contract.
		if log.Address != s.BeaconCfg().Execution.DepositContractAddress {
			continue
		}

		// Switch statement to handle different log types.
		var err error
		switch logSig := log.Topics[0]; {
		case logSig == stakinglogs.DepositSig:
			err = s.processDepositLog(ctx, log)
		case logSig == stakinglogs.RedirectSig:
			err = s.processRedirectLog(ctx, log)
		case logSig == stakinglogs.WithdrawalSig:
			err = s.processWithdrawalLog(ctx, log)
		default:
			continue
		}
		if err != nil {
			s.Logger().Error("failed to process log", "err", err)
			return err
		}
	}
	return nil
}

// processDepositLog adds a deposit to the queue.
func (s *Service) processDepositLog(
	ctx context.Context,
	log coretypes.Log,
) error {
	deposit := new(beacontypes.Deposit)
	if err := deposit.UnmarshalEthLog(log); err != nil {
		return err
	}
	s.Logger().
		Info("he was a sk8r boi 🛹",
			"deposit", deposit.Index, "amount", deposit.Amount)
	return s.BeaconState(ctx).EnqueueDeposits([]*beacontypes.Deposit{deposit})
}

// processRedirectLog adds a redirect to the queue.
func (s *Service) processRedirectLog(_ context.Context, _ coretypes.Log) error {
	return nil
}

// processWithdrawalLog adds a withdrawal to the queue.
func (s *Service) processWithdrawalLog(
	_ context.Context,
	_ coretypes.Log,
) error {
	return nil
}
