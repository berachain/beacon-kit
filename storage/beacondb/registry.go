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

	"cosmossdk.io/collections/indexes"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// AddValidator registers a new validator in the beacon state.
func (kv *KVStore) AddValidator(val *ctypes.Validator) error {
	// Get the next validator index from the sequence.
	idx, err := kv.validatorIndex.Next(kv.ctx)
	if err != nil {
		return err
	}

	// Push onto the validators list.
	if err = kv.validators.Set(kv.ctx, idx, val); err != nil {
		return err
	}

	return kv.balances.Set(kv.ctx, idx, 0)
}

// AddValidator registers a new validator in the beacon state.
func (kv *KVStore) AddValidatorBartio(val *ctypes.Validator) error {
	// Get the ne
	idx, err := kv.validatorIndex.Next(kv.ctx)
	if err != nil {
		return err
	}

	// Push onto the validators list.
	if err = kv.validators.Set(kv.ctx, idx, val); err != nil {
		return err
	}

	// Push onto the balances list.
	return kv.balances.Set(kv.ctx, idx, val.GetEffectiveBalance().Unwrap())
}

// UpdateValidatorAtIndex updates a validator at a specific index.
func (kv *KVStore) UpdateValidatorAtIndex(
	index math.ValidatorIndex,
	val *ctypes.Validator,
) error {
	return kv.validators.Set(kv.ctx, index.Unwrap(), val)
}

// ValidatorIndexByPubkey returns the validator address by index.
func (kv *KVStore) ValidatorIndexByPubkey(
	pubkey crypto.BLSPubkey,
) (math.ValidatorIndex, error) {
	idx, err := kv.validators.Indexes.Pubkey.MatchExact(
		kv.ctx,
		pubkey[:],
	)
	if err != nil {
		return 0, err
	}
	return math.ValidatorIndex(idx), nil
}

// ValidatorIndexByCometBFTAddress returns the validator address by index.
func (kv *KVStore) ValidatorIndexByCometBFTAddress(
	cometBFTAddress []byte,
) (math.ValidatorIndex, error) {
	idx, err := kv.validators.Indexes.CometBFTAddress.MatchExact(
		kv.ctx,
		cometBFTAddress,
	)
	if err != nil {
		return 0, err
	}
	return math.ValidatorIndex(idx), nil
}

// ValidatorByIndex returns the validator address by index.
func (kv *KVStore) ValidatorByIndex(
	index math.ValidatorIndex,
) (*ctypes.Validator, error) {
	val, err := kv.validators.Get(kv.ctx, index.Unwrap())
	if err != nil {
		var t *ctypes.Validator
		return t, err
	}
	return val, err
}

// GetValidators retrieves all validators from the beacon state.
func (kv *KVStore) GetValidators() (
	ctypes.Validators, error,
) {
	registrySize, err := kv.validatorIndex.Peek(kv.ctx)
	if err != nil {
		return nil, err
	}

	var (
		vals = make([]*ctypes.Validator, 0, registrySize)
		val  *ctypes.Validator
	)

	iter, err := kv.validators.Iterate(kv.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, iter.Close())
	}()

	for ; iter.Valid(); iter.Next() {
		val, err = iter.Value()
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}

	return vals, err
}

// GetTotalValidators returns the total number of validators.
func (kv *KVStore) GetTotalValidators() (uint64, error) {
	validators, err := kv.GetValidators()
	if err != nil {
		return 0, err
	}
	return uint64(len(validators)), nil
}

// GetValidatorsByEffectiveBalance retrieves all validators sorted by
// effective balance from the beacon state.
func (kv *KVStore) GetValidatorsByEffectiveBalance() (
	[]*ctypes.Validator, error,
) {
	var (
		vals []*ctypes.Validator
		v    *ctypes.Validator
		idx  uint64
	)

	iter, err := kv.validators.Indexes.EffectiveBalance.Iterate(
		kv.ctx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, iter.Close())
	}()

	// Iterate over all validators and collect them.
	for ; iter.Valid(); iter.Next() {
		idx, err = iter.PrimaryKey()
		if err != nil {
			return nil, err
		}
		if v, err = kv.validators.Get(kv.ctx, idx); err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, err
}

// GetBalance returns the balance of a validator.
func (kv *KVStore) GetBalance(
	idx math.ValidatorIndex,
) (math.Gwei, error) {
	balance, err := kv.balances.Get(kv.ctx, idx.Unwrap())
	return math.Gwei(balance), err
}

// SetBalance sets the balance of a validator.
func (kv *KVStore) SetBalance(
	idx math.ValidatorIndex,
	balance math.Gwei,
) error {
	return kv.balances.Set(kv.ctx, idx.Unwrap(), balance.Unwrap())
}

// GetBalances returns the balancse of all validator.
func (kv *KVStore) GetBalances() ([]uint64, error) {
	var balances []uint64
	iter, err := kv.balances.Iterate(kv.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, iter.Close())
	}()

	var balance uint64
	for ; iter.Valid(); iter.Next() {
		balance, err = iter.Value()
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	return balances, err
}

// GetTotalActiveBalances returns the total active balances of all validatorkv.
// TODO: unhood this and probably store this as just a value changed on writekv.
// TODO: this shouldn't live in KVStore
func (kv *KVStore) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (math.Gwei, error) {
	slot, err := kv.slot.Get(kv.ctx)
	if err != nil {
		return 0, err
	}

	totalActiveBalances := math.Gwei(0)
	epoch := math.Epoch(slot / slotsPerEpoch)

	iter, err := kv.validators.Indexes.EffectiveBalance.Iterate(kv.ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = errors.Join(err, iter.Close())
	}()

	err = indexes.ScanValues(
		kv.ctx, kv.validators, iter, func(v *ctypes.Validator,
		) bool {
			if v.IsActive(epoch) {
				totalActiveBalances += v.GetEffectiveBalance()
			}
			return false
		},
	)
	return totalActiveBalances, err
}
