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
	"github.com/berachain/beacon-kit/mod/primitives"
)

// UpdateBlockRootAtIndex sets a block root in the BeaconStore.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) UpdateBlockRootAtIndex(
	index uint64,
	root primitives.Root,
) error {
	return kv.blockRoots.Set(kv.ctx, index, root[:])
}

// GetBlockRootAtIndex retrieves the block root from the BeaconStore.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetBlockRootAtIndex(
	index uint64,
) (primitives.Root, error) {
	bz, err := kv.blockRoots.Get(kv.ctx, index)
	if err != nil {
		return primitives.Root{}, err
	}
	return primitives.Root(bz), nil
}

// SetLatestBlockHeader sets the latest block header in the BeaconStore.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) SetLatestBlockHeader(
	header BeaconBlockHeaderT,
) error {
	return kv.latestBlockHeader.Set(kv.ctx, header)
}

// GetLatestBlockHeader retrieves the latest block header from the BeaconStore.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetLatestBlockHeader() (
	BeaconBlockHeaderT, error,
) {
	return kv.latestBlockHeader.Get(kv.ctx)
}

// UpdateStateRootAtIndex updates the state root at the given slot.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot primitives.Root,
) error {
	return kv.stateRoots.Set(kv.ctx, idx, stateRoot[:])
}

// StateRootAtIndex returns the state root at the given slot.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) StateRootAtIndex(
	idx uint64,
) (primitives.Root, error) {
	bz, err := kv.stateRoots.Get(kv.ctx, idx)
	if err != nil {
		return primitives.Root{}, err
	}
	return primitives.Root(bz), nil
}
