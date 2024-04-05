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
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// AddValidator registers a new validator in the beacon state.
func (s *StateDB) AddValidator(
	val *beacontypes.Validator,
) error {
	panic("light mode does not support writes")
}

// UpdateValidatorAtIndex updates a validator at a specific index.
func (s *StateDB) UpdateValidatorAtIndex(
	index primitives.ValidatorIndex,
	val *beacontypes.Validator,
) error {
	panic("light mode does not support writes")
}

// RemoveValidatorAtIndex removes a validator at a specified index.
func (s *StateDB) RemoveValidatorAtIndex(
	idx primitives.ValidatorIndex,
) error {
	panic("light mode does not support writes")
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (s *StateDB) ValidatorIndexByPubkey(
	pubkey primitives.BLSPubkey,
) (primitives.ValidatorIndex, error) {
	panic("not implemented")
}

// ValidatorByIndex returns the validator address by index.
func (s *StateDB) ValidatorByIndex(
	index primitives.ValidatorIndex,
) (*beacontypes.Validator, error) {
	panic("not implemented")
}

// GetValidators retrieves all validators from the beacon state.
func (s *StateDB) GetValidators() (
	[]*beacontypes.Validator, error,
) {
	panic("not implemented")
}

// GetValidatorsByEffectiveBalance retrieves all validators sorted by
// effective balance from the beacon state.
func (s *StateDB) GetValidatorsByEffectiveBalance() (
	[]*beacontypes.Validator, error,
) {
	panic("not implemented")
}

// IncreaseBalance increases the balance of a validator.
func (s *StateDB) IncreaseBalance(
	idx primitives.ValidatorIndex,
	delta primitives.Gwei,
) error {
	panic("light mode does not support writes")
}

// DecreaseBalance decreases the balance of a validator.
func (s *StateDB) DecreaseBalance(
	idx primitives.ValidatorIndex,
	delta primitives.Gwei,
) error {
	panic("light mode does not support writes")
}

// GetBalances returns the balancse of all validator.
func (s *StateDB) GetBalances() ([]uint64, error) {
	panic("not implemented")
}

// GetTotalActiveBalances returns the total active balances of all validators.
// TODO: unhood this and probably store this as just a value changed on writes.
func (s *StateDB) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (primitives.Gwei, error) {
	panic("not implemented")
}
