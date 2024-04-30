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
	"github.com/berachain/beacon-kit/mod/beacon/staking/abi"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ProcessBlockEvents processes the logs from the deposit contract.
func (s *Service) ProcessBlockEvents(
	logs []engineprimitives.Log,
) error {
	var deposits []*primitives.Deposit
	for _, log := range logs {
		// We only care about logs from the deposit contract.
		if log.Address != s.ChainSpec().DepositContractAddress() {
			continue
		}

		// Switch statement to handle different log types.
		switch logSig := log.Topics[0]; {
		case logSig == DepositEventSig:
			deposit, err := s.unpackDepositLog(log)
			if err != nil {
				s.Logger().Error("failed to unpack deposit log", "err", err)
				return err
			}
			deposits = append(deposits, deposit)
			s.Logger().Info(
				"he was a sk8r boi 🛹", "deposit", deposit.Index, "amount", deposit.Amount,
			)
		default:
			continue
		}
	}

	return s.ds.EnqueueDeposits(deposits)
}

// unpackDepositLog takes a log from the deposit contract and unpacks it into a
// Deposit struct.
func (s *Service) unpackDepositLog(
	log engineprimitives.Log,
) (*primitives.Deposit, error) {
	d := &abi.BeaconDepositContractDeposit{}
	if err := s.abi.UnpackLogs(d, DepositEventName, log); err != nil {
		return nil, err
	}

	return primitives.NewDeposit(
		crypto.BLSPubkey(d.Pubkey),
		primitives.WithdrawalCredentials(d.Credentials),
		math.Gwei(d.Amount),
		crypto.BLSSignature(d.Signature),
		d.Index,
	), nil
}
