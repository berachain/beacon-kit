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
	"bytes"
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
)

// UpdateBlockRootAtIndex sets a block root in the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) UpdateBlockRootAtIndex(
	index uint64,
	root common.Root,
) error {
	err := kv.sszDB.SetBlockRootAtIndex(kv.ctx, index, root)
	if err != nil {
		return err
	}
	return kv.blockRoots.Set(kv.ctx, index, root[:])
}

// GetBlockRootAtIndex retrieves the block root from the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetBlockRootAtIndex(
	index uint64,
) (common.Root, error) {
	bz, err := kv.blockRoots.Get(kv.ctx, index)
	if err != nil {
		return common.Root{}, err
	}
	sszBz, err := kv.sszDB.GetBlockRootAtIndex(kv.ctx, index)
	if err != nil {
		return common.Root{}, err
	}
	if !bytes.Equal(bz, sszBz[:]) {
		return common.Root{}, errors.New("block root mismatch")
	}
	return common.Root(bz), nil
}

// SetLatestBlockHeader sets the latest block header in the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) SetLatestBlockHeader(
	header BeaconBlockHeaderT,
) error {
	err := kv.sszDB.SetObject(kv.ctx, "latest_block_header", header)
	if err != nil {
		return err
	}
	return kv.latestBlockHeader.Set(kv.ctx, header)
}

// GetLatestBlockHeader retrieves the latest block header from the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) GetLatestBlockHeader() (
	BeaconBlockHeaderT, error,
) {
	var header BeaconBlockHeaderT
	header = header.Empty()
	err := kv.sszDB.GetObject(kv.ctx, "latest_block_header", header)
	if err != nil {
		return header, err
	}
	return kv.latestBlockHeader.Get(kv.ctx)
}

// UpdateStateRootAtIndex updates the state root at the given slot.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot common.Root,
) error {
	err := kv.sszDB.SetStateRootAtIndex(kv.ctx, idx, stateRoot)
	if err != nil {
		return err
	}
	return kv.stateRoots.Set(kv.ctx, idx, stateRoot[:])
}

// StateRootAtIndex returns the state root at the given slot.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT,
]) StateRootAtIndex(
	idx uint64,
) (common.Root, error) {
	bz, err := kv.stateRoots.Get(kv.ctx, idx)
	if err != nil {
		return common.Root{}, err
	}
	path := fmt.Sprintf("state_roots/%d", idx)
	sszBz, err := kv.sszDB.GetPath(kv.ctx, sszdb.ObjectPath(path))
	if err != nil {
		return common.Root{}, err
	}
	if !bytes.Equal(bz, sszBz[:]) {
		return common.Root{}, errors.New("state root mismatch")
	}

	return common.Root(bz), nil
}
