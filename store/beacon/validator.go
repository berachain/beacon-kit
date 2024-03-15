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

	"github.com/berachain/beacon-kit/primitives"
)

// AddValidator registers a new validator in the beacon state.
func (s *Store) AddValidator(
	ctx context.Context,
	pubkey []byte,
) error {
	idx, err := s.validatorIndex.Next(ctx)
	if err != nil {
		return err
	}
	return s.validatorIndexToPubkey.Set(ctx, idx, pubkey)
}

// UpdateValidator updates the pubkey of a validator.
func (s *Store) UpdateValidator(
	ctx context.Context,
	oldPubkey []byte,
	newPubkey []byte,
) error {
	// Get the index of the old pubkey.
	idx, err := s.validatorIndexToPubkey.Indexes.Pubkey.MatchExact(
		ctx,
		oldPubkey,
	)
	if err != nil {
		return err
	}

	// Set the new one
	return s.validatorIndexToPubkey.Set(ctx, idx, newPubkey)
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (s *Store) ValidatorIndexByPubkey(
	pubkey []byte,
) (primitives.ValidatorIndex, error) {
	idx, err := s.validatorIndexToPubkey.Indexes.Pubkey.MatchExact(
		s.ctx,
		pubkey,
	)
	if err != nil {
		return 0, err
	}
	return idx, nil
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (s *Store) ValidatorPubKeyByIndex(
	index primitives.ValidatorIndex,
) ([]byte, error) {
	pubkey, err := s.validatorIndexToPubkey.Get(s.ctx, index)
	if err != nil {
		return nil, err
	}
	return pubkey, err
}
