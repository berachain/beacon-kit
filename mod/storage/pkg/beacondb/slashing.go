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
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	storage "github.com/berachain/beacon-kit/mod/storage/pkg"
)

func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetSlashings() ([]math.Gwei, error) {
	var slashings []math.Gwei
	iter, err := kv.slashings.Iterate(kv.ctx, nil)
	err = storage.MapError(err)
	if err != nil {
		return nil, fmt.Errorf(
			"failed iterating slashings: %w",
			err,
		)
	}

	for iter.Valid() {
		var slashing uint64
		slashing, err = iter.Value()
		if err != nil {
			return nil, err
		}
		slashings = append(slashings, math.Gwei(slashing))
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
	return kv.slashings.Set(kv.ctx, index, amount.Unwrap())
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
	return kv.totalSlashing.Set(kv.ctx, amount.Unwrap())
}
