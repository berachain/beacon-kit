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
	"cosmossdk.io/collections/indexes"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// AddValidator registers a new validator in the beacon state.
func (s *StateDB) AddValidator(
	val *beacontypes.Validator,
) error {
	// Get the ne
	idx, err := s.validatorIndex.Next(s.ctx)
	if err != nil {
		return err
	}

	// Push onto the validators list.
	if err = s.validators.Set(s.ctx, idx, val); err != nil {
		return err
	}

	// Push onto the balances list.
	return s.balances.Set(s.ctx, idx, uint64(val.EffectiveBalance))
}

// UpdateValidatorAtIndex updates a validator at a specific index.
func (s *StateDB) UpdateValidatorAtIndex(
	index primitives.ValidatorIndex,
	val *beacontypes.Validator,
) error {
	return s.validators.Set(s.ctx, uint64(index), val)
}

// RemoveValidatorAtIndex removes a validator at a specified index.
func (s *StateDB) RemoveValidatorAtIndex(
	idx primitives.ValidatorIndex,
) error {
	return s.validators.Remove(s.ctx, uint64(idx))
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (s *StateDB) ValidatorIndexByPubkey(
	pubkey primitives.BLSPubkey,
) (primitives.ValidatorIndex, error) {
	idx, err := s.validators.Indexes.Pubkey.MatchExact(
		s.ctx,
		pubkey[:],
	)
	if err != nil {
		return 0, err
	}
	return primitives.ValidatorIndex(idx), nil
}

// ValidatorByIndex returns the validator address by index.
func (s *StateDB) ValidatorByIndex(
	index primitives.ValidatorIndex,
) (*beacontypes.Validator, error) {
	val, err := s.validators.Get(s.ctx, uint64(index))
	if err != nil {
		return nil, err
	}
	return val, err
}

// GetValidators retrieves all validators from the beacon state.
func (s *StateDB) GetValidators() (
	[]*beacontypes.Validator, error,
) {
	var (
		vals []*beacontypes.Validator
		val  *beacontypes.Validator
	)

	iter, err := s.validators.Iterate(s.ctx, nil)
	if err != nil {
		return nil, err
	}

	for iter.Valid() {
		val, err = iter.Value()
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
		iter.Next()
	}

	return vals, nil
}

// GetValidatorsByEffectiveBalance retrieves all validators sorted by
// effective balance from the beacon state.
func (s *StateDB) GetValidatorsByEffectiveBalance() (
	[]*beacontypes.Validator, error,
) {
	var (
		vals []*beacontypes.Validator
		v    *beacontypes.Validator
		idx  uint64
	)

	iter, err := s.validators.Indexes.EffectiveBalance.Iterate(
		s.ctx,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Iterate over all validators and collect them.
	for ; iter.Valid(); iter.Next() {
		idx, err = iter.PrimaryKey()
		if err != nil {
			return nil, err
		}
		if v, err = s.validators.Get(s.ctx, idx); err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

// IncreaseBalance increases the balance of a validator.
func (s *StateDB) IncreaseBalance(
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
func (s *StateDB) DecreaseBalance(
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

// GetBalances returns the balancse of all validator.
func (s *StateDB) GetBalances() ([]uint64, error) {
	var balances []uint64
	iter, err := s.balances.Iterate(s.ctx, nil)
	if err != nil {
		return nil, err
	}

	var balance uint64
	for iter.Valid() {
		balance, err = iter.Value()
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
		iter.Next()
	}
	return balances, nil
}

// GetTotalActiveBalances returns the total active balances of all validators.
// TODO: unhood this and probably store this as just a value changed on writes.
func (s *StateDB) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (primitives.Gwei, error) {
	iter, err := s.validators.Indexes.EffectiveBalance.Iterate(s.ctx, nil)
	if err != nil {
		return 0, err
	}

	slot, err := s.slot.Get(s.ctx)
	if err != nil {
		return 0, err
	}

	totalActiveBalances := primitives.Gwei(0)
	epoch := primitives.Epoch(slot / slotsPerEpoch)
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
