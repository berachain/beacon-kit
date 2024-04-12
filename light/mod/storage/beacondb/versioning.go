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

import (
	"github.com/berachain/beacon-kit/light/mod/storage/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// SetGenesisValidatorsRoot sets the genesis validators root in the beacon
// state.
func (kv *KVStore) SetGenesisValidatorsRoot(
	root primitives.Root,
) error {
	panic(writesNotSupported)
}

// GetGenesisValidatorsRoot retrieves the genesis validators root from the
// beacon state.
func (kv *KVStore) GetGenesisValidatorsRoot() (primitives.Root, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.genesisValidatorsRoot.Key(),
		0,
	)
	if err != nil {
		return primitives.Root{}, err
	}

	genesisValidatorsRoot, err := kv.genesisValidatorsRoot.Decode(res)
	if err != nil {
		return primitives.Root{}, err
	}

	return genesisValidatorsRoot, nil
}

// GetSlot returns the current slot.
func (kv *KVStore) GetSlot() (primitives.Slot, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.slot.Key(),
		0,
	)
	if err != nil {
		return 0, err
	}

	slot, err := kv.slot.Decode(res)
	if err != nil {
		return 0, err
	}

	return primitives.Slot(slot), nil
}

// SetSlot sets the current slot.
func (kv *KVStore) SetSlot(slot primitives.Slot) error {
	panic(writesNotSupported)
}
