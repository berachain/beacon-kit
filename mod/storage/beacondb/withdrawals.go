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

import "github.com/berachain/beacon-kit/mod/primitives/math"

// GetNextWithdrawalIndex returns the next withdrawal index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetNextWithdrawalIndex() (uint64, error) {
	return kv.nextWithdrawalIndex.Get(kv.ctx)
}

// SetNextWithdrawalIndex sets the next withdrawal index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) SetNextWithdrawalIndex(
	index uint64,
) error {
	return kv.nextWithdrawalIndex.Set(kv.ctx, index)
}

// GetNextWithdrawalValidatorIndex returns the next withdrawal validator index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetNextWithdrawalValidatorIndex() (
	math.ValidatorIndex, error,
) {
	idx, err := kv.nextWithdrawalValidatorIndex.Get(kv.ctx)
	return math.ValidatorIndex(idx), err
}

// SetNextWithdrawalValidatorIndex sets the next withdrawal validator index.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) SetNextWithdrawalValidatorIndex(
	index math.ValidatorIndex,
) error {
	return kv.nextWithdrawalValidatorIndex.Set(kv.ctx, uint64(index))
}
