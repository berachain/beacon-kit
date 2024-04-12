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
	"github.com/berachain/beacon-kit/light/mod/storage/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// GetLatestExecutionPayload retrieves the latest execution payload from the
// BeaconStore.
func (kv *KVStore) GetLatestExecutionPayload() (types.ExecutionPayload, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.latestExecutionPayload.Key(),
		0,
	)
	if err != nil {
		return &types.ExecutableDataDeneb{}, err
	}

	payload, err := kv.latestExecutionPayload.Decode(res)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

// UpdateLatestExecutionPayload sets the latest execution payload in the
// BeaconStore.
func (kv *KVStore) UpdateLatestExecutionPayload(
	payload types.ExecutionPayload,
) error {
	panic(writesNotSupported)
}

// GetEth1DepositIndex retrieves the eth1 deposit index from the beacon state.
func (kv *KVStore) GetEth1DepositIndex() (uint64, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.eth1DepositIndex.Key(),
		0,
	)
	if err != nil {
		return 0, err
	}

	deposit, err := kv.eth1DepositIndex.Decode(res)
	if err != nil {
		return deposit, err
	}

	return deposit, nil
}

// SetEth1DepositIndex sets the eth1 deposit index in the beacon state.
func (kv *KVStore) SetEth1DepositIndex(index uint64) error {
	panic(writesNotSupported)
}

// GetEth1Data retrieves the eth1 data from the beacon state.
func (kv *KVStore) GetEth1Data() (*primitives.Eth1Data, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.eth1Data.Key(),
		0,
	)
	if err != nil {
		return &primitives.Eth1Data{}, err
	}

	data, err := kv.eth1Data.Decode(res)
	if err != nil {
		return data, err
	}

	return data, nil
}

// SetEth1Data sets the eth1 data in the beacon state.
func (kv *KVStore) SetEth1Data(data *primitives.Eth1Data) error {
	panic(writesNotSupported)
}
