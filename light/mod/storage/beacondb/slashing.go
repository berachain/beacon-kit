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

func (kv *KVStore) GetSlashings() ([]uint64, error) {
	// var slashings []uint64
	// iter, err := kv.slashings.Iterate(kv.ctx, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// for iter.Valid() {
	// 	var slashing uint64
	// 	slashing, err = iter.Value()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	slashings = append(slashings, slashing)
	// 	iter.Next()
	// }
	// return slashings, nil
	panic("not implemented")
}

// GetSlashingAtIndex retrieves the slashing amount by index from the store.
func (kv *KVStore) GetSlashingAtIndex(index uint64) (primitives.Gwei, error) {
	key, err := kv.slashings.Key(index)
	if err != nil {
		return 0, err
	}

	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		key,
		0,
	)
	if err != nil {
		return 0, err
	}

	amount, err := kv.slashings.Decode(res)
	if err != nil {
		return 0, err
	}

	return primitives.Gwei(amount), nil
}

// SetSlashingAtIndex sets the slashing amount in the store.
func (kv *KVStore) SetSlashingAtIndex(
	index uint64,
	amount primitives.Gwei,
) error {
	panic(writesNotSupported)
}

// TotalSlashing retrieves the total slashing amount from the store.
func (kv *KVStore) GetTotalSlashing() (primitives.Gwei, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.totalSlashing.Key(),
		0,
	)
	if err != nil {
		return 0, err
	}

	total, err := kv.totalSlashing.Decode(res)
	if err != nil {
		return 0, err
	}

	return primitives.Gwei(total), nil
}

// SetTotalSlashing sets the total slashing amount in the store.
func (kv *KVStore) SetTotalSlashing(amount primitives.Gwei) error {
	panic(writesNotSupported)
}
