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
	"errors"

	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/mod/primitives"
)

func (s *Store) GetSlashings() ([]uint64, error) {
	var slashings []uint64
	iter, err := s.slashings.Iterate(s.ctx, nil)
	if err != nil {
		return nil, err
	}
	for iter.Valid() {
		var slashing uint64
		slashing, err = iter.Value()
		if err != nil {
			return nil, err
		}
		slashings = append(slashings, slashing)
		iter.Next()
	}
	return slashings, nil
}

// UpdateSlashingAtIndex sets the slashing amount in the store.
func (s *Store) UpdateSlashingAtIndex(
	index uint64,
	amount primitives.Gwei,
) error {
	// Update the total slashing amount before overwriting the old amount.
	total, err := s.totalSlashing.Get(s.ctx)
	if !errors.Is(err, sdkcollections.ErrNotFound) && err != nil {
		return err
	}

	oldValue, err := s.GetSlashingAtIndex(index)
	if !errors.Is(err, sdkcollections.ErrNotFound) && err != nil {
		return err
	}

	// Defensive check but total - oldValue should never underflow.
	if uint64(oldValue) > total {
		return errors.New("count of total slashing is not up to date")
	}
	if err = s.totalSlashing.Set(
		s.ctx,
		total-uint64(oldValue)+uint64(amount),
	); err != nil {
		return err
	}

	return s.slashings.Set(s.ctx, index, uint64(amount))
}

// GetSlashingAtIndex retrieves the slashing amount by index from the store.
func (s *Store) GetSlashingAtIndex(index uint64) (primitives.Gwei, error) {
	amount, err := s.slashings.Get(s.ctx, index)
	if err != nil {
		return 0, err
	}
	return primitives.Gwei(amount), nil
}

// TotalSlashing retrieves the total slashing amount from the store.
func (s *Store) GetTotalSlashing() (primitives.Gwei, error) {
	total, err := s.totalSlashing.Get(s.ctx)
	if err != nil {
		return 0, err
	}
	return primitives.Gwei(total), nil
}

// SetTotalSlashing sets the total slashing amount in the store.
func (s *Store) SetTotalSlashing(amount primitives.Gwei) error {
	return s.totalSlashing.Set(s.ctx, uint64(amount))
}
