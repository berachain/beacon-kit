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

package types

import (
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	"github.com/berachain/beacon-kit/config/version"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

// BeaconBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconBlock(
	slot primitives.Slot,
	executionData enginetypes.ExecutionPayload,
	parentBlockRoot [32]byte,
	forkVersion uint32,
	reveal types.Reveal,
) (BeaconBlock, error) {
	var block BeaconBlock
	switch forkVersion {
	case version.Deneb:
		block = &BeaconBlockDeneb{
			Slot:            slot,
			ParentBlockRoot: parentBlockRoot,
			Body: &BeaconBlockBodyDeneb{
				RandaoReveal: reveal,
				Graffiti:     [32]byte{},
			},
		}
	default:
		return nil, ErrForkVersionNotSupported
	}

	if executionData != nil {
		if err := block.GetBody().SetExecutionData(executionData); err != nil {
			return nil, err
		}
	}
	return block, nil
}

// EmptyBeaconBlock assembles a new beacon block
// with no execution data.
func EmptyBeaconBlock(
	slot primitives.Slot,
	parentBlockRoot [32]byte,
	version uint32,
	reveal types.Reveal,
) (BeaconBlock, error) {
	return NewBeaconBlock(slot, nil, parentBlockRoot, version, reveal)
}

// BeaconBlockFromSSZ assembles a new beacon block
// from the given SSZ bytes and fork version.
func BeaconBlockFromSSZ(
	bz []byte,
	forkVersion uint32,
) (BeaconBlock, error) {
	var block BeaconBlock
	switch forkVersion {
	case version.Deneb:
		block = &BeaconBlockDeneb{}
	default:
		return nil, ErrForkVersionNotSupported
	}

	if err := block.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return block, nil
}
