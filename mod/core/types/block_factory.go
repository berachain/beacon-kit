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
	"unsafe"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/version"
)

// BeaconBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconBlock[BeaconBlockT primitives.BeaconBlock](
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot primitives.Root,
	stateRoot primitives.Root,
	forkVersion uint32,
	reveal primitives.BLSSignature,
) (BeaconBlockT, error) {
	var block BeaconBlockT
	switch forkVersion {
	case version.Deneb:
		blockDeneb := &primitives.BeaconBlockDeneb{
			Slot:            slot,
			ProposerIndex:   proposerIndex,
			ParentBlockRoot: parentBlockRoot,
			StateRoot:       stateRoot,
			Body: &primitives.BeaconBlockBodyDeneb{
				RandaoReveal: reveal,
				Graffiti:     [32]byte{},
			},
		}
		block = *(*BeaconBlockT)(unsafe.Pointer(&blockDeneb))
	default:
		return block, ErrForkVersionNotSupported
	}
	return block, nil
}

// EmptyBeaconBlock assembles a new beacon block
// with no execution data.
func EmptyBeaconBlock[BeaconBlockT primitives.BeaconBlock](
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot primitives.Root,
	stateRoot primitives.Root,
	version uint32,
	reveal primitives.BLSSignature,
) (BeaconBlockT, error) {
	var block BeaconBlockT
	newBlock, err := NewBeaconBlock[BeaconBlockT](
		slot,
		proposerIndex,
		parentBlockRoot,
		stateRoot,
		version,
		reveal,
	)
	if err != nil {
		return block, err
	}
	return newBlock, nil
}

// BeaconBlockFromSSZ assembles a new beacon block
// from the given SSZ bytes and fork version.
func BeaconBlockFromSSZ(
	bz []byte,
	forkVersion uint32,
) (primitives.BeaconBlock, error) {
	var block primitives.BeaconBlockDeneb
	switch forkVersion {
	case version.Deneb:
		_ = block
	default:
		return &block, ErrForkVersionNotSupported
	}

	if err := block.UnmarshalSSZ(bz); err != nil {
		return &block, err
	}
	return &block, nil
}
