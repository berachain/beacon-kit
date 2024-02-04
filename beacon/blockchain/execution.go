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

package blockchain

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

// GetOrBuildBlock constructs the next block in the execution chain.
func (s *Service) GetOrBuildBlock(
	ctx context.Context, slot primitives.Slot,
) (interfaces.BeaconKitBlock, error) {
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	var (
		beaconState        = s.bsp.BeaconState(ctx)
		executionData      interfaces.ExecutionData
		lastFinalizedBlock = beaconState.GetFinalizedEth1BlockHash()
	)

	// Create a new empty block from the current state.
	beaconBlock, err := consensusv1.NewEmptyBeaconKitBlockFromState(beaconState)
	if err != nil {
		return nil, err
	}

	// Attempt to get a previously built payload, otherwise we will trigger a payload to be build.
	// Building a payload at this point is not ideal, as it will block the consensus client
	// from proposing a block until the payload is built. However, this is a rare case, and
	// we can optimize this later.
	//
	// TODO: 4844.
	if executionData, _, _, err = s.en.GetBuiltPayload(
		ctx, slot, lastFinalizedBlock,
	); err != nil || executionData == nil {
		// This branch represents a cache miss.
		telemetry.IncrCounter(1, MetricGetBuiltPayloadMiss)
		executionData, err = s.buildNewPayloadForBlock(ctx, beaconBlock, lastFinalizedBlock)
		if err != nil {
			return nil, err
		}
	} else {
		// This branch represents a cache hit.
		telemetry.IncrCounter(1, MetricGetBuiltPayloadHit)
	}

	// Assemble a new block with the payload.
	if err = beaconBlock.AttachExecutionData(executionData); err != nil {
		return nil, err
	}

	// Return the block.
	return beaconBlock, nil
}

// buildNewPayloadForBlock begins building a new payload ontop of `headHash` for the given block,
// and waits for the payload to be built by the execution client.
func (s *Service) buildNewPayloadForBlock(
	ctx context.Context, beaconBlock interfaces.BeaconKitBlock, headHash common.Hash,
) (interfaces.ExecutionData, error) {
	// Notify the execution client of the new finalized and safe hashes, while also
	// triggering a block to be built. We wait for the payload ID to be returned.
	fcuConfig := &execution.FCUConfig{
		HeadEth1Hash:  headHash,
		ProposingSlot: beaconBlock.GetSlot(),
		Attributes:    s.getPayloadAttribute(ctx),
	}

	err := s.en.NotifyForkchoiceUpdate(
		ctx, fcuConfig,
	)

	if err != nil {
		return nil, err
	}

	return s.waitForPayload(ctx, payloadBuildDelay*time.Second, beaconBlock.GetSlot(), headHash)
}

// waitForPayload waits for a payload to be built by the execution client.
func (s *Service) waitForPayload(
	ctx context.Context, delay time.Duration, slot primitives.Slot, headHash common.Hash,
) (interfaces.ExecutionData, error) {
	// Create a timer for the delay duration
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		// The context was cancelled before the timer expired
		s.Logger().Error(
			"Context cancelled while waiting for payload", "slot", slot, "head_hash", headHash,
		)
		return nil, ctx.Err()
	case <-timer.C:
		// Hitting this case is basically the desired behaviour.
		break
	}

	// If neither context cancellation nor timeout occurred, proceed to get the payload
	payload, _, _, err := s.en.GetBuiltPayload(ctx, slot, headHash)
	if err != nil {
		s.Logger().Error(
			"Failed to get built payload", "error", err, "slot", slot, "head_hash", headHash,
		)
		return nil, err
	}
	return payload, nil
}
