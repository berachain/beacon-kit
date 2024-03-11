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

// TODO: Consider consolidating BeaconState interface externally to x/beacon
// to facilitate withdrawals from x/beacon_staking.
// TODO: Explore constructing BeaconState from multiple sources beyond
// just x/beacon.
func (s *Store) ExpectedWithdrawals() ([]*enginetypes.Withdrawal, error) {
	return []*enginetypes.Withdrawal{}, nil
}
