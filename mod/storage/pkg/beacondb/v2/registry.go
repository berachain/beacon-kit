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
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections/indexes"
)

// AddValidator registers a new validator in the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) AddValidator(
	val ValidatorT,
) error {
	// Get the next validator index.
	idx, err := s.validatorIndex.Next()
	if err != nil {
		return err
	}

	// Push onto the validators list.
	if err = s.validators.Set(idx, val); err != nil {
		return err
	}

	// Push onto the balances list.
	return s.balances.Set(idx, uint64(val.GetEffectiveBalance()))
}

// AddValidator registers a new validator in the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) AddValidatorBartio(val ValidatorT) error {
	// Get the next validator index.
	idx, err := s.validatorIndex.Next()
	if err != nil {
		return err
	}

	// Push onto the validators list.
	if err = s.validators.Set(idx, val); err != nil {
		return err
	}

	// Push onto the balances list.
	return s.balances.Set(idx, uint64(val.GetEffectiveBalance()))
}

// UpdateValidatorAtIndex updates a validator at a specific index.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) UpdateValidatorAtIndex(
	index math.ValidatorIndex,
	val ValidatorT,
) error {
	return s.validators.Set(uint64(index), val)
}

// RemoveValidatorAtIndex removes a validator at a specified index.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) RemoveValidatorAtIndex(
	idx math.ValidatorIndex,
) error {
	return s.validators.Remove(uint64(idx))
}

// ValidatorIndexByPubkey returns the validator address by index.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) ValidatorIndexByPubkey(
	pubkey crypto.BLSPubkey,
) (math.ValidatorIndex, error) {
	idx, err := s.validators.Indexes.Pubkey.MatchExact(pubkey[:])
	if err != nil {
		return 0, err
	}
	return math.ValidatorIndex(idx), nil
}

// ValidatorIndexByCometBFTAddress returns the validator address by index.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) ValidatorIndexByCometBFTAddress(
	cometBFTAddress []byte,
) (math.ValidatorIndex, error) {
	idx, err := s.validators.Indexes.CometBFTAddress.MatchExact(
		cometBFTAddress,
	)
	if err != nil {
		return 0, err
	}
	return math.ValidatorIndex(idx), nil
}

// ValidatorByIndex returns the validator address by index.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) ValidatorByIndex(
	index math.ValidatorIndex,
) (ValidatorT, error) {
	val, err := s.validators.Get(uint64(index))
	if err != nil {
		var t ValidatorT
		return t, err
	}
	return val, err
}

// GetValidators retrieves all validators from the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetValidators() (
	[]ValidatorT, error,
) {
	var (
		vals []ValidatorT
		val  ValidatorT
	)
	s.Save()
	iter, err := s.validators.Iterate()
	if err != nil {
		return nil, err
	}

	keys, err := iter.Keys()
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		val, err = s.validators.Get(k)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				continue
			}
			return nil, err
		}
		vals = append(vals, val)
	}
	iter.Close()
	return vals, nil
}

// GetTotalValidators returns the total number of validators.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetTotalValidators() (uint64, error) {
	// validators, err := s.GetValidators()
	// if err != nil {
	// 	return 0, err
	// }
	// return uint64(len(validators)), nil
	return s.validators.NumKeys()
}

// GetValidatorsByEffectiveBalance retrieves all validators sorted by
// effective balance from the beacon state.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetValidatorsByEffectiveBalance() (
	[]ValidatorT, error,
) {
	var (
		vals []ValidatorT
		v    ValidatorT
		idx  uint64
	)

	iter, err := s.validators.Indexes.EffectiveBalance.Iterate()
	if err != nil {
		return nil, err
	}

	// Iterate over all validators and collect them.
	for ; iter.Valid(); iter.Next() {
		idx, err = iter.PrimaryKey()
		if err != nil {
			return nil, err
		}
		if v, err = s.validators.Get(idx); err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

// GetBalance returns the balance of a validator.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetBalance(
	idx math.ValidatorIndex,
) (math.Gwei, error) {
	balance, err := s.balances.Get(uint64(idx))
	return math.Gwei(balance), err
}

// SetBalance sets the balance of a validator.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) SetBalance(
	idx math.ValidatorIndex,
	balance math.Gwei,
) error {
	return s.balances.Set(uint64(idx), uint64(balance))
}

// GetBalances returns the balancse of all validator.
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetBalances() ([]uint64, error) {
	var balances []uint64
	iter, err := s.balances.Iterate()
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
// TODO: this shouldn't live in Store
func (s *StateManager[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (math.Gwei, error) {
	iter, err := s.validators.Indexes.EffectiveBalance.Iterate()
	if err != nil {
		return 0, err
	}

	slot, err := s.slot.Get()
	if err != nil {
		return 0, err
	}

	totalActiveBalances := math.Gwei(0)
	epoch := math.Epoch(slot / slotsPerEpoch)
	return totalActiveBalances, indexes.ScanValues(
		s.validators, iter, func(v ValidatorT,
		) bool {
			if v.IsActive(epoch) {
				totalActiveBalances += v.GetEffectiveBalance()
			}
			return false
		},
	)
}
