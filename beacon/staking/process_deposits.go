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

	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

// AcceptDepositsIntoQueue records a list
// of deposits in the beacon state's queue.
func (s *Service) AcceptDepositsIntoQueue(
	ctx context.Context,
	deposits []*consensusv1.Deposit,
) error {
	// Push the deposits to the beacon state's queue.
	err := s.BeaconState(ctx).EnqueueDeposits(deposits)
	if err != nil {
		return err
	}
	return nil
}

// DequeueDeposits returns up to MaxDepositsPerBlock deposits
// from the beacon state's queue.
func (s *Service) DequeueDeposits(
	ctx context.Context,
) ([]*consensusv1.Deposit, error) {
	return s.BeaconState(ctx).DequeueDeposits(
		s.BeaconCfg().Limits.MaxDepositsPerBlock,
	)
}

// ApplyDeposits processes the deposits in the beacon state's queue,
// up to MaxDepositsPerBlock, by applying them to the underlying staking module.
func (s *Service) ApplyDeposits(ctx context.Context) error {
	beaconState := s.BeaconState(ctx)

	// Get deposits, up to MaxDepositsPerBlock, from the queue
	// to apply to the underlying low-level staking module (e.g Cosmos SDK's
	// x/staking).
	deposits, err := beaconState.DequeueDeposits(
		s.BeaconCfg().Limits.MaxDepositsPerBlock,
	)
	if err != nil {
		return err
	}

	// Apply deposists to the underlying staking module.
	return s.vcp.ApplyChanges(ctx, deposits, nil)
}
