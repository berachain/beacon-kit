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
)

// UpdateRandaoMixAtIndex sets the current RANDAO mix in the store.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) UpdateRandaoMixAtIndex(
	index uint64,
	mix primitives.Bytes32,
) error {
	return kv.randaoMix.Set(kv.ctx, index, mix[:])
}

// GetRandaoMixAtIndex retrieves the current RANDAO mix from the store.
func (kv *KVStore[
	ForkT, BeaconBlockHeaderT, ExecutionPayloadT, Eth1DataT, ValidatorT,
]) GetRandaoMixAtIndex(
	index uint64,
) (primitives.Bytes32, error) {
	bz, err := kv.randaoMix.Get(kv.ctx, index)
	if err != nil {
		return primitives.Bytes32{}, err
	}
	return primitives.Bytes32(bz), nil
}
