// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package beacondb

import (
	"cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
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
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetSlashingAtIndex(
	index uint64,
) (math.Gwei, error) {
	amount, err := kv.slashings.Get(kv.ctx, index)
	if errors.Is(err, collections.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return math.Gwei(amount), nil
}

// SetSlashingAtIndex sets the slashing amount in the store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetSlashingAtIndex(
	index uint64,
	amount math.Gwei,
) error {
	return kv.slashings.Set(kv.ctx, index, uint64(amount))
}

// GetTotalSlashing retrieves the total slashing amount from the store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetTotalSlashing() (math.Gwei, error) {
	total, err := kv.totalSlashing.Get(kv.ctx)
	if errors.Is(err, collections.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return math.Gwei(total), nil
}

// SetTotalSlashing sets the total slashing amount in the store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetTotalSlashing(
	amount math.Gwei,
) error {
	return kv.totalSlashing.Set(kv.ctx, uint64(amount))
}

// IncreaseBalance increases the balance of a validator.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) IncreaseBalance(
	idx math.ValidatorIndex,
	delta math.Gwei,
) error {
	balance, err := kv.GetBalance(idx)
	if err != nil {
		return err
	}
	return kv.SetBalance(idx, balance+delta)
}

// DecreaseBalance decreases the balance of a validator.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) DecreaseBalance(
	idx math.ValidatorIndex,
	delta math.Gwei,
) error {
	balance, err := kv.GetBalance(idx)
	if err != nil {
		return err
	}
	return kv.SetBalance(idx, balance-min(balance, delta))
}

// UpdateSlashingAtIndex sets the slashing amount in the store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) UpdateSlashingAtIndex(
	index uint64,
	amount math.Gwei,
) error {
	// Update the total slashing amount before overwriting the old amount.
	total, err := kv.GetTotalSlashing()
	if err != nil {
		return err
	}

	oldValue, err := kv.GetSlashingAtIndex(index)
	if err != nil {
		return err
	}

	// Defensive check but total - oldValue should never underflow.
	if oldValue > total {
		return errors.New("count of total slashing is not up to date")
	} else if err = kv.SetTotalSlashing(
		total - oldValue + amount,
	); err != nil {
		return err
	}

	return kv.SetSlashingAtIndex(index, amount)
}
