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

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

// UpdateRandaoMixAtIndex sets the current RANDAO mix in the store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) UpdateRandaoMixAtIndex(
	index uint64,
	mix common.Bytes32,
) error {
	return kv.randaoMix.Set(kv.ctx, index, mix[:])
}

// GetRandaoMixAtIndex retrieves the current RANDAO mix from the store.
func (kv *KVStore[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT,
]) GetRandaoMixAtIndex(
	index uint64,
) (common.Bytes32, error) {
	bz, err := kv.randaoMix.Get(kv.ctx, index)
	if err != nil {
		return common.Bytes32{}, err
	}
	return common.Bytes32(bz), nil
}
