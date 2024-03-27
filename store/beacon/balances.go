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
	"cosmossdk.io/collections/indexes"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
)

// IncreaseBalance increases the balance of a validator.
func (s *Store) IncreaseBalance(
	idx primitives.ValidatorIndex,
	delta primitives.Gwei,
) error {
	balance, err := s.balances.Get(s.ctx, uint64(idx))
	if err != nil {
		return err
	}
	balance += uint64(delta)
	return s.balances.Set(s.ctx, uint64(idx), balance)
}

// DecreaseBalance decreases the balance of a validator.
func (s *Store) DecreaseBalance(
	idx primitives.ValidatorIndex,
	delta primitives.Gwei,
) error {
	balance, err := s.balances.Get(s.ctx, uint64(idx))
	if err != nil {
		return err
	}
	balance -= min(balance, uint64(delta))
	return s.balances.Set(s.ctx, uint64(idx), balance)
}

// GetTotalActiveBalances returns the total active balances of all validators.
// TODO: unhood this and probably store this as just a value changed on writes.
func (s *Store) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (primitives.Gwei, error) {
	iter, err := s.validators.Indexes.EffectiveBalance.Iterate(s.ctx, nil)
	if err != nil {
		return 0, err
	}

	epoch, err := s.GetEpoch(slotsPerEpoch)
	if err != nil {
		return 0, err
	}

	totalActiveBalances := primitives.Gwei(0)
	return totalActiveBalances, indexes.ScanValues(
		s.ctx, s.validators, iter, func(v *beacontypes.Validator,
		) bool {
			if v.IsActive(epoch) {
				totalActiveBalances += v.EffectiveBalance
			}
			return false
		},
	)
}
