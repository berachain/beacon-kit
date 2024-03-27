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
	stakingabi "github.com/berachain/beacon-kit/contracts/abi"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// ProcessBlockEvents processes the logs from the deposit contract.
func (s *Service) ProcessBlockEvents(
	ctx context.Context,
	logs []coretypes.Log,
) error {
	for _, log := range logs {
		// We only care about logs from the deposit contract.
		if log.Address != s.BeaconCfg().DepositContractAddress {
			continue
		}

		// Switch statement to handle different log types.
		var err error
		switch logSig := log.Topics[0]; {
		case logSig == DepositEventSig:
			err = s.processDepositLog(ctx, log)
		case logSig == WithdrawalEventSig:
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
	d := &stakingabi.BeaconDepositContractDeposit{}
	if err := s.abi.UnpackLogs(d, DepositEventName, log); err != nil {
		return err
	}

	s.Logger().Info(
		"he was a sk8r boi 🛹", "deposit", d.Index, "amount", d.Amount,
	)

	return s.BeaconState(ctx).EnqueueDeposits([]*beacontypes.Deposit{{
		Index:       d.Index,
		Pubkey:      primitives.BLSPubkey(d.Pubkey),
		Credentials: beacontypes.WithdrawalCredentials(d.Credentials),
		Amount:      primitives.Gwei(d.Amount),
		Signature:   primitives.BLSSignature(d.Signature),
	}})
}

// processWithdrawalLog adds a withdrawal to the queue.
func (s *Service) processWithdrawalLog(
	ctx context.Context,
	log coretypes.Log,
) error {
	w := &stakingabi.BeaconDepositContractWithdrawal{}
	if err := s.abi.UnpackLogs(w, WithdrawalEventName, log); err != nil {
		return err
	}

	// Get the validator index from the pubkey.
	valIdx, err := s.BeaconState(ctx).ValidatorIndexByPubkey(w.FromPubkey)
	if err != nil {
		return err
	}

	executionAddr, err := beacontypes.
		WithdrawalCredentials(w.Credentials).ToExecutionAddress()
	if err != nil {
		return err
	}

	s.Logger().Info(
		"she said, \"see you later, boi\" 💅", "deposit", w.Index, "amount", w.Amount,
	)

	return s.BeaconState(ctx).EnqueueWithdrawals([]*enginetypes.Withdrawal{{
		Index:     w.Index,
		Validator: valIdx,
		Address:   executionAddr,
		Amount:    primitives.Gwei(w.Amount),
	}})
}
