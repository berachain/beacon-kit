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

package statedb

import (
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// ExpectedDeposits returns the first numPeek deposits in the queue.
func (s *StateDB) ExpectedDeposits(
	numView uint64,
) ([]*beacontypes.Deposit, error) {
	panic("light mode does not support writes")
}

// EnqueueDeposits pushes the deposits to the queue.
func (s *StateDB) EnqueueDeposits(
	deposits []*beacontypes.Deposit,
) error {
	panic("light mode does not support writes")
}

// DequeueDeposits returns the first numDequeue deposits in the queue.
func (s *StateDB) DequeueDeposits(
	numDequeue uint64,
) ([]*beacontypes.Deposit, error) {
	panic("light mode does not support writes")
}

// ExpectedWithdrawals returns the first numView withdrawals in the queue.
func (s *StateDB) ExpectedWithdrawals(
	numView uint64,
) ([]*primitives.Withdrawal, error) {
	panic("not implemented")
}

// EnqueueWithdrawals pushes the withdrawals to the queue.
func (s *StateDB) EnqueueWithdrawals(
	withdrawals []*primitives.Withdrawal,
) error {
	panic("light mode does not support writes")
}

// EnqueueWithdrawals pushes the withdrawals to the queue.
func (s *StateDB) DequeueWithdrawals(
	numDequeue uint64,
) ([]*primitives.Withdrawal, error) {
	panic("light mode does not support writes")
}

// // handleNilWithdrawals returns an empty slice if the input is nil.
// func handleNilWithdrawals(
// 	withdrawals []*primitives.Withdrawal,
// ) []*primitives.Withdrawal {
// 	if withdrawals == nil {
// 		return make([]*primitives.Withdrawal, 0)
// 	}
// 	return withdrawals
// }
