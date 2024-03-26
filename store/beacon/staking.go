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

package beacon

import (
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
)

// GetEth1DepositIndex retrieves the eth1 deposit index from the beacon state.
func (s *Store) GetEth1DepositIndex() (uint64, error) {
	return s.eth1DepositIndex.Get(s.ctx)
}

// SetEth1DepositIndex sets the eth1 deposit index in the beacon state.
func (s *Store) SetEth1DepositIndex(index uint64) error {
	return s.eth1DepositIndex.Set(s.ctx, index)
}

// ExpectedDeposits returns the first numPeek deposits in the queue.
func (s *Store) ExpectedDeposits(
	numView uint64,
) ([]*beacontypes.Deposit, error) {
	return s.depositQueue.PeekMulti(s.ctx, numView)
}

// EnqueueDeposits pushes the deposits to the queue.
func (s *Store) EnqueueDeposits(
	deposits []*beacontypes.Deposit,
) error {
	return s.depositQueue.PushMulti(s.ctx, deposits)
}

// DequeueDeposits returns the first numDequeue deposits in the queue.
func (s *Store) DequeueDeposits(
	numDequeue uint64,
) ([]*beacontypes.Deposit, error) {
	return s.depositQueue.PopMulti(s.ctx, numDequeue)
}

// ExpectedWithdrawals returns the first numView withdrawals in the queue.
func (s *Store) ExpectedWithdrawals(
	numView uint64,
) ([]*enginetypes.Withdrawal, error) {
	withdrawals, err := s.withdrawalQueue.PeekMulti(s.ctx, numView)
	if err != nil {
		return nil, err
	}
	return s.handleNilWithdrawals(withdrawals), nil
}

// EnqueueWithdrawals pushes the withdrawals to the queue.
func (s *Store) EnqueueWithdrawals(
	withdrawals []*enginetypes.Withdrawal,
) error {
	return s.withdrawalQueue.PushMulti(s.ctx, withdrawals)
}

// EnqueueWithdrawals pushes the withdrawals to the queue.
func (s *Store) DequeueWithdrawals(
	numDequeue uint64,
) ([]*enginetypes.Withdrawal, error) {
	withdrawals, err := s.withdrawalQueue.PopMulti(s.ctx, numDequeue)
	if err != nil {
		return nil, err
	}
	return s.handleNilWithdrawals(withdrawals), nil
}

// handleNilWithdrawals returns an empty slice if the input is nil.
func (Store) handleNilWithdrawals(
	withdrawals []*enginetypes.Withdrawal,
) []*enginetypes.Withdrawal {
	if withdrawals == nil {
		return make([]*enginetypes.Withdrawal, 0)
	}
	return withdrawals
}
