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
	"cosmossdk.io/collections/indexes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// AddValidator registers a new validator in the beacon state.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) AddValidator(
	val ValidatorT,
) error {
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
	return kv.balances.Set(kv.ctx, idx, uint64(val.GetEffectiveBalance()))
}

// UpdateValidatorAtIndex updates a validator at a specific index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) UpdateValidatorAtIndex(
	index math.ValidatorIndex,
	val ValidatorT,
) error {
	return kv.validators.Set(kv.ctx, uint64(index), val)
}

// RemoveValidatorAtIndex removes a validator at a specified index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) RemoveValidatorAtIndex(
	idx math.ValidatorIndex,
) error {
	return kv.validators.Remove(kv.ctx, uint64(idx))
}

// ValidatorPubKeyByIndex returns the validator address by index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
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

// ValidatorByIndex returns the validator address by index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
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
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) ValidatorByIndex(
	index math.ValidatorIndex,
) (ValidatorT, error) {
	val, err := kv.validators.Get(kv.ctx, uint64(index))
	if err != nil {
		var t ValidatorT
		return t, err
	}
	return val, err
}

// GetValidators retrieves all validators from the beacon state.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetValidators() (
	[]ValidatorT, error,
) {
	var (
		vals []ValidatorT
		val  ValidatorT
	)

	iter, err := kv.validators.Iterate(kv.ctx, nil)
	if err != nil {
		return nil, err
	}

	for iter.Valid() {
		val, err = iter.Value()
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
		iter.Next()
	}

	return vals, nil
}

// GetTotalValidators returns the total number of validators.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
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
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
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
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetBalance(
	idx math.ValidatorIndex,
) (math.Gwei, error) {
	balance, err := kv.balances.Get(kv.ctx, uint64(idx))
	return math.Gwei(balance), err
}

// SetBalance sets the balance of a validator.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) SetBalance(
	idx math.ValidatorIndex,
	balance math.Gwei,
) error {
	return kv.balances.Set(kv.ctx, uint64(idx), uint64(balance))
}

// GetBalances returns the balancse of all validator.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
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
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
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
