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

package beacondb

import "github.com/berachain/beacon-kit/mod/primitives"

// SetGenesisValidatorsRoot sets the genesis validators root in the beacon
// state.
func (kv *KVStore) SetGenesisValidatorsRoot(
	root primitives.Root,
) error {
	return kv.genesisValidatorsRoot.Set(kv.ctx, root)
}

// GetGenesisValidatorsRoot retrieves the genesis validators root from the
// beacon state.
func (kv *KVStore) GetGenesisValidatorsRoot() (primitives.Root, error) {
	return kv.genesisValidatorsRoot.Get(kv.ctx)
}

// GetSlot returns the current slot.
func (kv *KVStore) GetSlot() (primitives.Slot, error) {
	slot, err := kv.slot.Get(kv.ctx)
	return primitives.Slot(slot), err
}

// SetSlot sets the current slot.
func (kv *KVStore) SetSlot(slot primitives.Slot) error {
	return kv.slot.Set(kv.ctx, uint64(slot))
}
