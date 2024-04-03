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

package beacon2

import (
	"context"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/beacond/store/beacon/collections/encoding"
	"github.com/berachain/beacon-kit/mod/core/state"
)

// beaconStateKeyPrefix is the key prefix for the beacon state.
const beaconStateKeyPrefix = "beacon_state"

// Store is a wrapper around an sdk.Context
// that provides access to all beacon related data.
type Store struct {
	// beaconState is the store for the beacon state.
	beaconState sdkcollections.Item[*state.BeaconStateDeneb]
}

// Store creates a new instance of Store.
//

func NewStore(
	env appmodule.Environment,
) *Store {
	return &Store{
		beaconState: sdkcollections.NewItem[*state.BeaconStateDeneb](
			sdkcollections.NewSchemaBuilder(env.KVStoreService),
			sdkcollections.NewPrefix(beaconStateKeyPrefix),
			beaconStateKeyPrefix,
			encoding.SSZValueCodec[*state.BeaconStateDeneb]{},
		),
	}
}

// GetBeaconState returns the beacon state.
func (s *Store) GetBeaconState(
	ctx context.Context,
) (*state.BeaconStateDeneb, error) {
	return s.beaconState.Get(ctx)
}

// SetBeaconState sets the beacon state.
func (s *Store) SetBeaconState(
	ctx context.Context,
	st *state.BeaconStateDeneb,
) error {
	return s.beaconState.Set(ctx, st)
}
