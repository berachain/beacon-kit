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

package consensus

import (
	"errors"

	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine"
)

// BeaconKitBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconKitBlock(
	slot primitives.Slot,
	executionData engine.ExecutionPayload,
	forkVersion int,
) (interfaces.BeaconKitBlock, error) {
	switch forkVersion {
	case version.Deneb:
		return nil, errors.New("TODO: Deneb block")
	case version.Capella:
		return consensusv1.NewBeaconKitBlock(slot, executionData)
	default:
		return nil, ErrForkVersionNotSupported
	}
}

// EmptyBeaconKitBlock assembles a new beacon block
// with no execution data.
func EmptyBeaconKitBlock(
	slot primitives.Slot,
	version int,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(slot, nil, version)
}

// EmptyBeaconKitBlockFromState assembles a new beacon block
// with no execution data from the given state.
func EmptyBeaconKitBlockFromState(
	beaconState state.BeaconState,
) (interfaces.BeaconKitBlock, error) {
	return EmptyBeaconKitBlock(
		beaconState.Slot(),
		beaconState.Version(),
	)
}

// BeaconKitBlockFromState assembles a new beacon block
// from the given state and execution data.
func BeaconKitBlockFromState(
	beaconState state.ReadOnlyBeaconState,
	executionData engine.ExecutionPayload,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(
		beaconState.Slot(),
		executionData,
		beaconState.Version(),
	)
}

// BeaconKitBlockFromSSZ assembles a new beacon block
// from the given SSZ bytes and fork version.
func BeaconKitBlockFromSSZ(
	bz []byte,
	forkVersion int,
) (interfaces.BeaconKitBlock, error) {
	switch forkVersion {
	case version.Deneb:
		return nil, errors.New("TODO: Deneb block")
	case version.Capella:
		block := &consensusv1.BeaconKitBlockCapella{}
		if err := block.UnmarshalSSZ(bz); err != nil {
			return nil, err
		}
		return block, nil
	default:
		return nil, ErrForkVersionNotSupported
	}
}
