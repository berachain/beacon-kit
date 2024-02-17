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

package validator

import (
	"context"

	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/engine"
)

// BuildBeaconBlock builds a new beacon block.
func (s *Service) BuildBeaconBlock(
	ctx context.Context, _ primitives.Slot,
) (interfaces.BeaconKitBlock, error) {
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	var (
		beaconState   = s.BeaconState(ctx)
		executionData engine.ExecutionPayload
		slot          = beaconState.Slot()
	)

	// // TODO: SIGN UR RANDAO THINGY HERE OR SOMETHING.
	_ = s.beaconKitValKey
	// _, err := s.beaconKitValKey.Key.PrivKey.Sign([]byte("hello world"))
	// if err != nil {
	// 	return nil, err
	// }

	// Create a new empty block from the current state.
	beaconBlock, err := consensus.EmptyBeaconKitBlock(
		slot, s.BeaconCfg().ActiveForkVersion(primitives.Epoch(slot)),
	)
	if err != nil {
		return nil, err
	}

	executionData, overrideBuilder, err := s.getLocalPayload(ctx, beaconBlock, beaconState)
	if err != nil {
		return nil, err
	}

	// TODO: allow external block builders to override the payload.
	_ = overrideBuilder

	// Assemble a new block with the payload.
	if err = beaconBlock.AttachExecution(executionData); err != nil {
		return nil, err
	}

	// Return the block.
	return beaconBlock, nil
}
