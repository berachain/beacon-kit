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
	"errors"

	"cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/mod/primitives"
)

func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetSlashings() ([]uint64, error) {
	var slashings []uint64
	iter, err := kv.slashings.Iterate(kv.ctx, nil)
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

// GetSlashingAtIndex retrieves the slashing amount by index from the store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetSlashingAtIndex(
	index uint64,
) (primitives.Gwei, error) {
	amount, err := kv.slashings.Get(kv.ctx, index)
	if errors.Is(err, collections.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return primitives.Gwei(amount), nil
}

// SetSlashingAtIndex sets the slashing amount in the store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) SetSlashingAtIndex(
	index uint64,
	amount primitives.Gwei,
) error {
	return kv.slashings.Set(kv.ctx, index, uint64(amount))
}

// TotalSlashing retrieves the total slashing amount from the store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetTotalSlashing() (primitives.Gwei, error) {
	total, err := kv.totalSlashing.Get(kv.ctx)
	if errors.Is(err, collections.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return primitives.Gwei(total), nil
}

// SetTotalSlashing sets the total slashing amount in the store.
func (kv *KVStore[
	DepositT, ForkT, BeaconBlockHeaderT,
	ExecutionPayloadT, Eth1DataT, ValidatorT,
]) SetTotalSlashing(
	amount primitives.Gwei,
) error {
	return kv.totalSlashing.Set(kv.ctx, uint64(amount))
}
