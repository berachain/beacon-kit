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
	"github.com/berachain/beacon-kit/mod/primitives"
)

func (s *StateDB) GetSlashings() ([]uint64, error) {
	panic("not implemented")
}

// UpdateSlashingAtIndex sets the slashing amount in the store.
func (s *StateDB) UpdateSlashingAtIndex(
	index uint64,
	amount primitives.Gwei,
) error {
	panic("light mode does not support writes")
}

// GetSlashingAtIndex retrieves the slashing amount by index from the store.
func (s *StateDB) GetSlashingAtIndex(index uint64) (primitives.Gwei, error) {
	panic("not implemented")
}

// TotalSlashing retrieves the total slashing amount from the store.
func (s *StateDB) GetTotalSlashing() (primitives.Gwei, error) {
	panic("not implemented")
}

// SetTotalSlashing sets the total slashing amount in the store.
func (s *StateDB) SetTotalSlashing(amount primitives.Gwei) error {
	panic("light mode does not support writes")
}
