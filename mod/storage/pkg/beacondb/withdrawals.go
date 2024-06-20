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

import "github.com/berachain/beacon-kit/mod/primitives/pkg/math"

// GetNextWithdrawalIndex returns the next withdrawal index.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetNextWithdrawalIndex() (uint64, error) {
	return kv.nextWithdrawalIndex.Get(kv.ctx)
}

// SetNextWithdrawalIndex sets the next withdrawal index.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) SetNextWithdrawalIndex(
	index uint64,
) error {
	return kv.nextWithdrawalIndex.Set(kv.ctx, index)
}

// GetNextWithdrawalValidatorIndex returns the next withdrawal validator index.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetNextWithdrawalValidatorIndex() (
	math.ValidatorIndex, error,
) {
	idx, err := kv.nextWithdrawalValidatorIndex.Get(kv.ctx)
	return math.ValidatorIndex(idx), err
}

// SetNextWithdrawalValidatorIndex sets the next withdrawal validator index.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) SetNextWithdrawalValidatorIndex(
	index math.ValidatorIndex,
) error {
	return kv.nextWithdrawalValidatorIndex.Set(kv.ctx, uint64(index))
}
