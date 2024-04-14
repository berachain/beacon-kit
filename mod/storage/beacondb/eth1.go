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
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
)

// GetLatestExecutionPayload retrieves the latest execution payload from the
// BeaconStore.
func (kv *KVStore) GetLatestExecutionPayload() (
	engineprimitives.ExecutionPayload, error,
) {
	return kv.latestExecutionPayload.Get(kv.ctx)
}

// UpdateLatestExecutionPayload sets the latest execution payload in the
// BeaconStore.
func (kv *KVStore) UpdateLatestExecutionPayload(
	payload engineprimitives.ExecutionPayload,
) error {
	return kv.latestExecutionPayload.Set(kv.ctx, payload)
}

// GetEth1DepositIndex retrieves the eth1 deposit index from the beacon state.
func (kv *KVStore) GetEth1DepositIndex() (uint64, error) {
	return kv.eth1DepositIndex.Get(kv.ctx)
}

// SetEth1DepositIndex sets the eth1 deposit index in the beacon state.
func (kv *KVStore) SetEth1DepositIndex(index uint64) error {
	return kv.eth1DepositIndex.Set(kv.ctx, index)
}

// GetEth1Data retrieves the eth1 data from the beacon state.
func (kv *KVStore) GetEth1Data() (*primitives.Eth1Data, error) {
	return kv.eth1Data.Get(kv.ctx)
}

// SetEth1Data sets the eth1 data in the beacon state.
func (kv *KVStore) SetEth1Data(data *primitives.Eth1Data) error {
	return kv.eth1Data.Set(kv.ctx, data)
}
