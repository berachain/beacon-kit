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
	"cosmossdk.io/collections/indexes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// AddValidator registers a new validator in the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) AddValidator(val ValidatorT) error {
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
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) AddValidatorBartio(val ValidatorT) error {
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
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) UpdateValidatorAtIndex(
	index math.ValidatorIndex,
	val ValidatorT,
) error {
	return kv.validators.Set(kv.ctx, index.Unwrap(), val)
}

// ValidatorIndexByPubkey returns the validator address by index.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) ValidatorIndexByPubkey(
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
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) ValidatorIndexByCometBFTAddress(
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
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) ValidatorByIndex(
	index math.ValidatorIndex,
) (ValidatorT, error) {
	val, err := kv.validators.Get(kv.ctx, index.Unwrap())
	if err != nil {
		var t ValidatorT
		return t, err
	}
	return val, err
}

// GetValidators retrieves all validators from the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetValidators() (
	ValidatorsT, error,
) {
	registrySize, err := kv.validatorIndex.Peek(kv.ctx)
	if err != nil {
		return nil, err
	}

	var (
		vals = make([]ValidatorT, registrySize)
		val  ValidatorT
	)

	iter, err := kv.validators.Iterate(kv.ctx, nil)
	if err != nil {
		return nil, err
	}

	i := 0
	for iter.Valid() {
		val, err = iter.Value()
		if err != nil {
			return nil, err
		}
		vals[i] = val
		iter.Next()
		i++
	}

	return vals, nil
}

// GetTotalValidators returns the total number of validators.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetTotalValidators() (uint64, error) {
	validators, err := kv.GetValidators()
	if err != nil {
		return 0, err
	}
	return uint64(len(validators)), nil
}

// GetValidatorsByEffectiveBalance retrieves all validators sorted by
// effective balance from the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetValidatorsByEffectiveBalance() (
	[]ValidatorT, error,
) {
	var (
		vals []ValidatorT
		v    ValidatorT
		idx  uint64
	)

	iter, err := kv.validators.Indexes.EffectiveBalance.Iterate(
		kv.ctx,
		nil,
	)
	if err != nil {
		return nil, err
	}

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
	return vals, nil
}

// GetBalance returns the balance of a validator.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetBalance(
	idx math.ValidatorIndex,
) (math.Gwei, error) {
	balance, err := kv.balances.Get(kv.ctx, idx.Unwrap())
	return math.Gwei(balance), err
}

// SetBalance sets the balance of a validator.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetBalance(
	idx math.ValidatorIndex,
	balance math.Gwei,
) error {
	return kv.balances.Set(kv.ctx, idx.Unwrap(), balance.Unwrap())
}

// GetBalances returns the balancse of all validator.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetBalances() ([]uint64, error) {
	var balances []uint64
	iter, err := kv.balances.Iterate(kv.ctx, nil)
	if err != nil {
		return nil, err
	}

	var balance uint64
	for iter.Valid() {
		balance, err = iter.Value()
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
		iter.Next()
	}
	return balances, nil
}

// GetTotalActiveBalances returns the total active balances of all validatorkv.
// TODO: unhood this and probably store this as just a value changed on writekv.
// TODO: this shouldn't live in KVStore
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (math.Gwei, error) {
	iter, err := kv.validators.Indexes.EffectiveBalance.Iterate(kv.ctx, nil)
	if err != nil {
		return 0, err
	}

	slot, err := kv.slot.Get(kv.ctx)
	if err != nil {
		return 0, err
	}

	totalActiveBalances := math.Gwei(0)
	epoch := math.Epoch(slot / slotsPerEpoch)
	return totalActiveBalances, indexes.ScanValues(
		kv.ctx, kv.validators, iter, func(v ValidatorT,
		) bool {
			if v.IsActive(epoch) {
				totalActiveBalances += v.GetEffectiveBalance()
			}
			return false
		},
	)
}
