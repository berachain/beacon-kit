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

package primitives

import "github.com/prysmaticlabs/prysm/v5/math"

// Slot is just a nice alias for SSZUint64.
type Slot = SSZUint64

const SlotsPerEpoch = 1

// Div divides slot by x.
// In case of arithmetic issues (overflow/underflow/div by zero) panic is thrown.
func (s Slot) Div(x uint64) Slot {
	res, err := s.SafeDiv(x)
	if err != nil {
		panic(err.Error())
	}
	return res
}

// SafeDiv divides slot by x.
// In case of arithmetic issues (overflow/underflow/div by zero) error is returned.
func (s Slot) SafeDiv(x uint64) (Slot, error) {
	res, err := math.Div64(uint64(s), x)
	return Slot(res), err
}

// DivSlot divides slot by another slot.
// In case of arithmetic issues (overflow/underflow/div by zero) panic is thrown.
func (s Slot) DivSlot(x Slot) Slot {
	return s.Div(uint64(x))
}

// ToEpoch returns the epoch number of the input slot.
//
// Spec pseudocode definition:
//
//	def compute_epoch_at_slot(slot: Slot) -> Epoch:
//	  """
//	  Return the epoch number at ``slot``.
//	  """
//	  return Epoch(slot // SLOTS_PER_EPOCH)
func ToEpoch(slot Slot) Epoch {
	return Epoch(slot.DivSlot(SlotsPerEpoch))
}
