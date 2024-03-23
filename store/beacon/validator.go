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
	"context"

	"cosmossdk.io/collections"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
)

// AddValidator registers a new validator in the beacon state.
func (s *Store) AddValidator(
	ctx context.Context,
	val *beacontypes.Validator,
) error {
	idx, err := s.validatorIndex.Next(ctx)
	if err != nil {
		return err
	}

	return s.validatorByIndex.Set(ctx, idx, val)
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (s *Store) ValidatorIndexByPubkey(
	pubkey []byte,
) (primitives.ValidatorIndex, error) {
	idx, err := s.validatorByIndex.Indexes.Pubkey.MatchExact(
		s.ctx,
		pubkey,
	)
	if err != nil {
		return 0, err
	}
	return idx, nil
}

// ValidatorIndexByConsAddress returns the validator index by consensus address.
func (s *Store) ValidatorIndexByConsAddr(
	consAddress []byte,
) (primitives.ValidatorIndex, error) {
	idx, err := s.validatorByIndex.Indexes.ConsAddr.MatchExact(
		s.ctx,
		consAddress,
	)
	if err != nil {
		return 0, err
	}
	return idx, nil
}

// GetAllValidators retrieves all validators from the beacon state.
// TODO: Use the heap and limit the number of validators that will
// be pulled here, cause this could get ugly runtime wise.
func (s *Store) GetAllValidators(
	ctx context.Context,
) ([]*beacontypes.Validator, error) {
	iter, err := s.validatorByIndex.IterateRaw(
		ctx,
		nil,
		nil,
		collections.OrderAscending,
	)
	if err != nil {
		return nil, err
	}

	var (
		vals []*beacontypes.Validator
		v    *beacontypes.Validator
	)

	for ; iter.Valid(); iter.Next() {
		v, err = iter.Value()
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

// ValidatorByIndex returns the validator address by index.
func (s *Store) ValidatorByIndex(
	index primitives.ValidatorIndex,
) (*beacontypes.Validator, error) {
	val, err := s.validatorByIndex.Get(s.ctx, index)
	if err != nil {
		return nil, err
	}
	return val, err
}
