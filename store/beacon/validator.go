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
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
)

// AddValidator registers a new validator in the beacon state.
func (s *Store) AddValidator(
	val *beacontypes.Validator,
) error {
	idx, err := s.validatorIndex.Next(s.ctx)
	if err != nil {
		return err
	}

	return s.validators.Set(s.ctx, idx, val)
}

// UpdateValidatorAtIndex updates a validator at a specific index.
func (s *Store) UpdateValidatorAtIndex(
	index primitives.ValidatorIndex,
	val *beacontypes.Validator,
) error {
	return s.validators.Set(s.ctx, index, val)
}

// AddValidator registers a new validator in the beacon state.
func (s *Store) RemoveValidatorAtIndex(
	idx primitives.ValidatorIndex,
) error {
	return s.validators.Remove(s.ctx, idx)
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
	return idx, nil
}

// ValidatorByIndex returns the validator address by index.
func (s *Store) ValidatorByIndex(
	index primitives.ValidatorIndex,
) (*beacontypes.Validator, error) {
	val, err := s.validators.Get(s.ctx, index)
	if err != nil {
		return nil, err
	}
	return val, err
}

// GetValidatorsByEffectiveBalance retrieves all validators from the
// beacon state.
func (s *Store) GetValidatorsByEffectiveBalance(
	limit int,
) ([]uint64, []*beacontypes.Validator, error) {
	var (
		vals     []*beacontypes.Validator
		valsIdxs []uint64
	)

	iter, err := s.validators.Indexes.EffectiveBalance.Iterate(
		s.ctx,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	// Iterate over all validators and collect them.
	err = indexes.ScanKeyValues(s.ctx, s.validators, iter, func(kv collections.KeyValue[uint64, *beacontypes.Validator]) (stop bool) {
		vals = append(vals, kv.Value)
		valsIdxs = append(valsIdxs, kv.Key)
		return limit > 0 && len(vals) == limit
	})
	if err != nil {
		return nil, nil, err
	}

	return nil, vals, nil
}

func (s *Store) GetLastValidatorSet() (*beacontypes.ValidatorSet, error) {
	return s.latestValidatorSet.Get(s.ctx)
}

func (s *Store) SetLastValidatorSet(
	valSet *beacontypes.ValidatorSet,
) error {
	return s.latestValidatorSet.Set(s.ctx, valSet)
}
