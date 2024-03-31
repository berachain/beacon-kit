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

package beacon

import (
	"github.com/berachain/beacon-kit/beacon/core/state"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
)

// UpdateStateRootAtIndex updates the state root at the given slot.
func (s *Store) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot primitives.Root,
) error {
	return s.stateRoots.Set(s.ctx, idx, stateRoot)
}

// StateRootAtIndex returns the state root at the given slot.
func (s *Store) StateRootAtIndex(idx uint64) (primitives.Root, error) {
	return s.stateRoots.Get(s.ctx, idx)
}

// Store is the interface for the beacon store.
func (s *Store) HashTreeRoot() ([32]byte, error) {
	// TODO: Implement getting the HashTreeRoot (StateRoot)
	// We currently return at least *SOMETHING* so that we
	// can simulate having to keep track of the StateRoot of the
	// BeaconState, since this value with change every slot.
	// TODO: Actually implementation.

	slot, err := s.GetSlot()
	if err != nil {
		return [32]byte{}, err
	}

	randaoMix, err := s.GetRandaoMixAtIndex(0)
	if err != nil {
		return [32]byte{}, err
	}

	randaoMixes := make([][32]byte, 32) //nolint:gomnd // temp.
	randaoMixes[0] = randaoMix

	latestBlockHeader, err := s.GetLatestBlockHeader()
	if err != nil {
		return [32]byte{}, err
	}

	return (&state.BeaconStateDeneb{
		GenesisValidatorsRoot: primitives.Root{},
		Slot:                  slot,
		LatestBlockHeader:     latestBlockHeader,
		BlockRoots:            make([][32]byte, 32), //nolint:gomnd // temp.
		StateRoots:            make([][32]byte, 32), //nolint:gomnd // temp.
		Eth1GenesisHash:       [32]byte{},
		Eth1DepositIndex:      0,
		Validators:            []*beacontypes.Validator{},
		RandaoMixes:           randaoMixes,
	}).HashTreeRoot()
}
