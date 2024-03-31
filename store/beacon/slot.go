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

import "github.com/berachain/beacon-kit/primitives"

// GetSlot returns the current slot.
func (s *Store) GetSlot() (primitives.Slot, error) {
	slot, err := s.slot.Get(s.ctx)
	return primitives.Slot(slot), err
}

// SetSlot sets the current slot.
func (s *Store) SetSlot(slot primitives.Slot) error {
	return s.slot.Set(s.ctx, uint64(slot))
}

// GetCurrentEpoch returns the current epoch.
func (s *Store) GetCurrentEpoch(
	slotsPerEpoch uint64,
) (primitives.Epoch, error) {
	slots, err := s.slot.Get(s.ctx)
	if err != nil {
		return 0, err
	}
	return primitives.Epoch(slots / slotsPerEpoch), nil
}
