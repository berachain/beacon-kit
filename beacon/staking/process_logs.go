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
	stakinglogs "github.com/berachain/beacon-kit/beacon/staking/logs"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// ProcessBlockEvents processes the logs from the deposit contract.
func (s *Service) ProcessBlockEvents(
	logs []coretypes.Log,
) error {
	for _, log := range logs {
		// We only care about logs from the deposit
		if log.Address != s.BeaconCfg().Execution.DepositContractAddress {
			continue
		}

		// Switch statement to handle different log types.
		var err error
		switch logSig := log.Topics[0]; {
		case logSig == stakinglogs.DepositSig:
			err = s.addDepositToQueue()
		case logSig == stakinglogs.RedirectSig:
			err = s.addRedirectToQueue()
		case logSig == stakinglogs.WithdrawalSig:
			err = s.addWithdrawalToQueue()
		default:
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// addDepositToQueue adds a deposit to the queue.
func (s *Service) addDepositToQueue() error {
	return nil
}

// addRedirectToQueue adds a redirect to the queue.
func (s *Service) addRedirectToQueue() error {
	return nil
}

// addWithdrawalToQueue adds a withdrawal to the queue.
func (s *Service) addWithdrawalToQueue() error {
	return nil
}
