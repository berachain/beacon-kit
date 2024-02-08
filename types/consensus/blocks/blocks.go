// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package blocks

import (
	"errors"

	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/v1/capella"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// BeaconKitBlockFromState assembles a new beacon block
// from the given state and execution data.
func BeaconKitBlockFromState(
	beaconState state.ReadOnlyBeaconState,
	executionData interfaces.ExecutionData,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(
		beaconState.Slot(),
		executionData,
		beaconState.Version(),
	)
}

// BeaconKitBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconKitBlock(
	slot primitives.Slot,
	executionData interfaces.ExecutionData,
	forkVersion int,
) (interfaces.BeaconKitBlock, error) {
	switch forkVersion {
	case version.Deneb:
		return nil, errors.New("TODO: Deneb block")
	case version.Capella:
		return capella.NewBeaconKitBlock(slot, executionData)
	default:
		return nil, ErrForkVersionNotSupported
	}
}

// NewEmptyBeaconKitBlockFromState assembles a new beacon block
// with no execution data from the given state.
func NewEmptyBeaconKitBlockFromState(
	beaconState state.BeaconState,
) (interfaces.BeaconKitBlock, error) {
	return NewEmptyBeaconKitBlock(
		beaconState.Slot(),
		beaconState.Version(),
	)
}

// NewEmptyBeaconKitBlock assembles a new beacon block
// with no execution data.
func NewEmptyBeaconKitBlock(
	slot primitives.Slot,
	version int,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(slot, nil, version)
}

// ReadOnlyBeaconKitBlockFromABCIRequest assembles a
// new read-only beacon block by extracting a marshalled
// block out of an ABCI request.
func ReadOnlyBeaconKitBlockFromABCIRequest(
	req interfaces.ABCIRequest,
	bzIndex uint,
	forkVersion int,
) (interfaces.ReadOnlyBeaconKitBlock, error) {
	txs := req.GetTxs()

	// Ensure there are transactions in the request and
	// that the request is valid.
	if lenTxs := uint(len(txs)); lenTxs == 0 {
		return nil, ErrNoBeaconBlockInProposal
	} else if bzIndex >= lenTxs {
		return nil, ErrBzIndexOutOfBounds
	}

	// Unmarshal the block from the request.
	switch forkVersion {
	case version.Deneb:
		return nil, errors.New("TODO: Deneb block")
	case version.Capella:
		block := &capella.BeaconKitBlockCapella{}
		if err := block.UnmarshalSSZ(txs[bzIndex]); err != nil {
			return nil, err
		}
		return block, nil
	default:
		return nil, ErrForkVersionNotSupported
	}
}
