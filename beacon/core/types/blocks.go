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
	beacontypesv1 "github.com/itsdevbear/bolaris/beacon/core/types/v1"
	"github.com/itsdevbear/bolaris/config/version"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	"github.com/itsdevbear/bolaris/primitives"
)

// BeaconBuoy assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconBuoy(
	slot primitives.Slot,
	executionData enginetypes.ExecutionPayload,
	parentBlockRoot [32]byte,
	forkVersion int,
) (BeaconBuoy, error) {
	var block BeaconBuoy
	switch forkVersion {
	case version.Deneb:
		block = &beacontypesv1.BeaconBuoyDeneb{
			Slot:       slot,
			ParentRoot: parentBlockRoot[:],
			Body: &beacontypesv1.BeaconBuoyBodyDeneb{
				RandaoReveal: make([]byte, 96), //nolint:gomnd
				Graffiti:     make([]byte, 32), //nolint:gomnd
			},
		}
	default:
		return nil, ErrForkVersionNotSupported
	}

	if executionData != nil {
		if err := block.AttachExecution(executionData); err != nil {
			return nil, err
		}
	}
	return block, nil
}

// EmptyBeaconBuoy assembles a new beacon block
// with no execution data.
func EmptyBeaconBuoy(
	slot primitives.Slot,
	parentBlockRoot [32]byte,
	version int,
) (BeaconBuoy, error) {
	return NewBeaconBuoy(slot, nil, parentBlockRoot, version)
}

// BeaconBuoyFromSSZ assembles a new beacon block
// from the given SSZ bytes and fork version.
func BeaconBuoyFromSSZ(
	bz []byte,
	forkVersion int,
) (BeaconBuoy, error) {
	var block BeaconBuoy
	switch forkVersion {
	case version.Deneb:
		block = &beacontypesv1.BeaconBuoyDeneb{}
	default:
		return nil, ErrForkVersionNotSupported
	}

	if err := block.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return block, nil
}
