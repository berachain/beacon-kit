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

// UpdateBlockRootAtIndex sets a block root in the BeaconStore.
func (kv *KVStore) UpdateBlockRootAtIndex(
	index uint64,
	root primitives.Root,
) error {
	panic(writesNotSupported)
}

// GetBlockRoot retrieves the block root from the BeaconStore.
func (kv *KVStore) GetBlockRootAtIndex(
	index uint64,
) (primitives.Root, error) {
	key, err := kv.blockRoots.Key(index)
	if err != nil {
		return primitives.Root{}, err
	}

	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		key,
		0,
	)
	if err != nil {
		return primitives.Root{}, err
	}

	blockRoot, err := kv.blockRoots.Decode(res)
	if err != nil {
		return blockRoot, err
	}

	return blockRoot, nil
}

// SetLatestBlockHeader sets the latest block header in the BeaconStore.
func (kv *KVStore) SetLatestBlockHeader(
	header *primitives.BeaconBlockHeader,
) error {
	panic(writesNotSupported)
}

// GetLatestBlockHeader retrieves the latest block header from the BeaconStore.
func (kv *KVStore) GetLatestBlockHeader() (
	*primitives.BeaconBlockHeader, error,
) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.latestBlockHeader.Key(),
		0,
	)
	if err != nil {
		return &primitives.BeaconBlockHeader{}, err
	}

	latestBlockHeader, err := kv.latestBlockHeader.Decode(res)
	if err != nil {
		return latestBlockHeader, err
	}

	return latestBlockHeader, nil
}

// UpdateStateRootAtIndex updates the state root at the given slot.
func (kv *KVStore) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot primitives.Root,
) error {
	panic(writesNotSupported)
}

// StateRootAtIndex returns the state root at the given slot.
func (kv *KVStore) StateRootAtIndex(index uint64) (primitives.Root, error) {
	key, err := kv.stateRoots.Key(index)
	if err != nil {
		return primitives.Root{}, err
	}

	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		key,
		0,
	)

	if err != nil {
		return primitives.Root{}, err
	}

	stateRoot, err := kv.stateRoots.Decode(res)
	if err != nil {
		return stateRoot, err
	}

	return stateRoot, nil
}
