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

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
)

// ProcessBlockEvents processes the logs from the deposit contract.
func (s *Service) ProcessBlockEvents(
	ctx context.Context,
	st state.BeaconState,
	logs []engineprimitives.Log,
) error {
	for _, log := range logs {
		// We only care about logs from the deposit contract.
		if log.Address != s.ChainSpec().DepositContractAddress() {
			continue
		}

		// Switch statement to handle different log types.
		var err error
		switch logSig := log.Topics[0]; {
		case logSig == DepositEventSig:
			err = s.processDepositLog(ctx, st, log)
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
	_ context.Context,
	st state.BeaconState,
	log engineprimitives.Log,
) error {
	d := &abi.BeaconDepositContractDeposit{}
	if err := s.abi.UnpackLogs(d, DepositEventName, log); err != nil {
		return err
	}

	s.Logger().Info(
		"he was a sk8r boi 🛹", "deposit", d.Index, "amount", d.Amount,
	)

	return st.EnqueueDeposits(
		[]*primitives.Deposit{primitives.NewDeposit(
			primitives.BLSPubkey(d.Pubkey),
			primitives.WithdrawalCredentials(d.Credentials),
			math.Gwei(d.Amount),
			primitives.BLSSignature(d.Signature),
			d.Index,
		)})
}
