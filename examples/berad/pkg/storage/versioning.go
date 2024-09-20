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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SetGenesisValidatorsRoot sets the genesis validators root in the beacon
// state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) SetGenesisValidatorsRoot(
	root common.Root,
) error {
	return kv.genesisValidatorsRoot.Set(kv.ctx, root[:])
}

// GetGenesisValidatorsRoot retrieves the genesis validators root from the
// beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetGenesisValidatorsRoot() (common.Root, error) {
	bz, err := kv.genesisValidatorsRoot.Get(kv.ctx)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

// GetSlot returns the current slot.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetSlot() (math.Slot, error) {
	slot, err := kv.slot.Get(kv.ctx)
	return math.Slot(slot), err
}

// SetSlot sets the current slot.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) SetSlot(
	slot math.Slot,
) error {
	return kv.slot.Set(kv.ctx, slot.Unwrap())
}

// SetFork sets the fork version for the given epoch.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) SetFork(
	fork ForkT,
) error {
	return kv.fork.Set(kv.ctx, fork)
}

// GetFork gets the fork version for the given epoch.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetFork() (ForkT, error) {
	return kv.fork.Get(kv.ctx)
}

// GetGenesisTime retrieves the genesis time from the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetGenesisTime() (uint64, error) {
	genesisTime, err := kv.genesisTime.Get(kv.ctx)
	if err != nil {
		return 0, err
	}
	return genesisTime, nil
}

// SetGenesisTime sets the genesis time in the beacon state.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) SetGenesisTime(
	genesisTime uint64,
) error {
	return kv.genesisTime.Set(kv.ctx, genesisTime)
}
