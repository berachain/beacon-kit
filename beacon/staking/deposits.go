// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

// ProcessDeposit processes a deposit log from the execution layer
// and puts the deposit to the beacon state.
func (s *Service) ProcessDeposit(
	_ context.Context,
	deposit *consensusv1.Deposit,
) error {
	// Cache the deposit to be pushed to the queue later in batch.
	s.depositCache = append(s.depositCache, deposit)
	s.Logger().Info("delegating from execution layer",
		"validatorPubkey", deposit.GetPubkey(), "amount", deposit.GetAmount())
	return nil
}

// PersistDeposits persists the queued deposists to the keeper.
func (s *Service) PersistDeposits(ctx context.Context) error {
	beaconState := s.BeaconState(ctx)

	// Push the cached deposits to the beacon state's queue.
	err := beaconState.StoreDeposits(s.depositCache)
	if err != nil {
		return err
	}
	s.depositCache = nil

	// Get deposits, up to MaxDepositsPerBlock, from the queue.
	deposits, err := beaconState.PopDeposits(s.BeaconCfg().Limits.MaxDepositsPerBlock)
	if err != nil {
		return err
	}

	// Apply deposists to the staking keeper.
	return s.vcp.ApplyChanges(ctx, deposits, nil)
}
