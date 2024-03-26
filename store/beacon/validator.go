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
	"github.com/berachain/beacon-kit/primitives"
)

// SetGenesisValidatorsRoot sets the genesis validators root in the beacon
// state.
func (s *Store) SetGenesisValidatorsRoot(
	root primitives.Root,
) error {
	return s.genesisValidatorsRoot.Set(s.ctx, root)
}

// GetGenesisValidatorsRoot retrieves the genesis validators root from the
// beacon state.
func (s *Store) GetGenesisValidatorsRoot() (primitives.Root, error) {
	return s.genesisValidatorsRoot.Get(s.ctx)
}

// AddValidator registers a new validator in the beacon state.
func (s *Store) AddValidator(
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
func (s *Store) UpdateValidatorAtIndex(
	index primitives.ValidatorIndex,
	val *beacontypes.Validator,
) error {
	return s.validators.Set(s.ctx, uint64(index), val)
}

// RemoveValidatorAtIndex removes a validator at a specified index.
func (s *Store) RemoveValidatorAtIndex(
	idx primitives.ValidatorIndex,
) error {
	return s.validators.Remove(s.ctx, uint64(idx))
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (s *Store) ValidatorIndexByPubkey(
	pubkey []byte,
) (primitives.ValidatorIndex, error) {
	idx, err := s.validators.Indexes.Pubkey.MatchExact(
		s.ctx,
		pubkey,
	)
	if err != nil {
		return 0, err
	}
	return primitives.ValidatorIndex(idx), nil
}

// ValidatorByIndex returns the validator address by index.
func (s *Store) ValidatorByIndex(
	index primitives.ValidatorIndex,
) (*beacontypes.Validator, error) {
	val, err := s.validators.Get(s.ctx, uint64(index))
	if err != nil {
		return nil, err
	}
	return val, err
}

// GetValidatorsByEffectiveBalance retrieves all validators sorted by
// effective balance from the beacon state.
func (s *Store) GetValidatorsByEffectiveBalance() (
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
