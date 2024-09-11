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

package storage

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

// UpdateBlockRootAtIndex sets a block root in the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) UpdateBlockRootAtIndex(
	index uint64,
	root common.Root,
) error {
	return kv.blockRoots.Set(kv.ctx, index, root[:])
}

// GetBlockRootAtIndex retrieves the block root from the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetBlockRootAtIndex(
	index uint64,
) (common.Root, error) {
	bz, err := kv.blockRoots.Get(kv.ctx, index)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

// SetLatestBlockHeader sets the latest block header in the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) SetLatestBlockHeader(
	header BeaconBlockHeaderT,
) error {
	return kv.latestBlockHeader.Set(kv.ctx, header)
}

// GetLatestBlockHeader retrieves the latest block header from the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetLatestBlockHeader() (
	BeaconBlockHeaderT, error,
) {
	return kv.latestBlockHeader.Get(kv.ctx)
}

// UpdateStateRootAtIndex updates the state root at the given slot.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot common.Root,
) error {
	return kv.stateRoots.Set(kv.ctx, idx, stateRoot[:])
}

// StateRootAtIndex returns the state root at the given slot.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) StateRootAtIndex(
	idx uint64,
) (common.Root, error) {
	bz, err := kv.stateRoots.Get(kv.ctx, idx)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

// GetLatestExecutionPayloadHeader retrieves the latest execution payload
// header from the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) GetLatestExecutionPayloadHeader() (
	ExecutionPayloadHeaderT, error,
) {
	forkVersion, err := kv.latestExecutionPayloadVersion.Get(kv.ctx)
	if err != nil {
		var t ExecutionPayloadHeaderT
		return t, err
	}
	kv.latestExecutionPayloadCodec.SetActiveForkVersion(forkVersion)
	return kv.latestExecutionPayloadHeader.Get(kv.ctx)
}

// SetLatestExecutionPayloadHeader sets the latest execution payload header in
// the BeaconStore.
func (kv *KVStore[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
]) SetLatestExecutionPayloadHeader(
	payloadHeader ExecutionPayloadHeaderT,
) error {
	if err := kv.latestExecutionPayloadVersion.Set(
		kv.ctx, payloadHeader.Version(),
	); err != nil {
		return err
	}
	kv.latestExecutionPayloadCodec.SetActiveForkVersion(payloadHeader.Version())
	return kv.latestExecutionPayloadHeader.Set(kv.ctx, payloadHeader)
}
